package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var (
	reviewBody string
	reviewComments []string
	reviewEvent string
	noExpandSuggestionsReview bool
)

var addReviewCmd = &cobra.Command{
	Use:   "add-review [pr] [review-body]",
	Short: "Create a review with multiple comments in one shot",
	Long: `Create a complete review with multiple line-specific comments.

Since GitHub's API requires all review comments to be created together,
this command allows you to specify multiple comments for a single review.

The review can be submitted immediately or left as a pending draft.

Examples:
  # Create pending review with multiple comments
  gh comment add-review 123 "Overall looks good" \
    --comment "src/api.js:42:This handles rate limiting well" \
    --comment "src/auth.js:15:20:Consider using async/await here"

  # Create and submit review immediately
  gh comment add-review 123 "LGTM with minor suggestions" \
    --event APPROVE \
    --comment "src/api.js:42:Great error handling" \
    --comment "src/utils.js:10:Minor: consider extracting this constant"

  # Auto-detect PR from current branch
  gh comment add-review "Code review feedback" \
    --comment "README.md:25:Update installation instructions" \
    --comment "package.json:15:Bump version number"`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runAddReview,
}

func init() {
	rootCmd.AddCommand(addReviewCmd)
	addReviewCmd.Flags().StringVar(&reviewBody, "body", "", "Review summary body (optional)")
	addReviewCmd.Flags().StringArrayVar(&reviewComments, "comment", []string{}, "Add comment in format 'file:line:message' or 'file:start:end:message' (can be used multiple times)")
	addReviewCmd.Flags().StringVar(&reviewEvent, "event", "", "Review event: APPROVE, REQUEST_CHANGES, or COMMENT (leave empty for pending review)")
	addReviewCmd.Flags().BoolVar(&noExpandSuggestionsReview, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax")
}

func runAddReview(cmd *cobra.Command, args []string) error {
	var pr int
	var body string
	var err error

	// Parse arguments
	if len(args) == 2 {
		// PR number and body provided
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return formatValidationError("PR number", args[0], "must be a valid integer")
		}
		body = args[1]
	} else if len(args) == 1 {
		// Check if it's a PR number or review body
		if prNum, err := strconv.Atoi(args[0]); err == nil {
			// It's a PR number
			pr = prNum
		} else {
			// It's a review body, auto-detect PR
			pr, err = getCurrentPR()
			if err != nil {
				return err
			}
			body = args[0]
		}
	} else {
		// Auto-detect PR
		pr, err = getCurrentPR()
		if err != nil {
			return err
		}
	}

	// Use --body flag if provided, otherwise use positional arg
	if reviewBody != "" {
		body = reviewBody
	}

	// Validate that we have at least one comment
	if len(reviewComments) == 0 {
		return fmt.Errorf("must provide at least one --comment")
	}

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR: %d\n", pr)
		fmt.Printf("Review body: %s\n", body)
		fmt.Printf("Review event: %s\n", reviewEvent)
		fmt.Printf("Comments: %d\n", len(reviewComments))
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would create review on PR #%d:\n", pr)
		fmt.Printf("Body: %s\n", body)
		fmt.Printf("Event: %s\n", reviewEvent)
		fmt.Printf("Comments:\n")
		for i, comment := range reviewComments {
			fmt.Printf("  %d. %s\n", i+1, comment)
		}
		return nil
	}

	// Create the review with all comments
	return createReviewWithComments(repository, pr, body, reviewEvent, reviewComments)
}

func createReviewWithComments(repo string, pr int, body, event string, commentSpecs []string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	// Get PR data for commit SHA
	prData := struct {
		Head struct {
			SHA string `json:"sha"`
		} `json:"head"`
	}{}

	err = client.Get(fmt.Sprintf("repos/%s/pulls/%d", repo, pr), &prData)
	if err != nil {
		return fmt.Errorf("failed to get PR data: %w", err)
	}

	// Parse comment specifications
	var comments []map[string]interface{}
	for _, spec := range commentSpecs {
		comment, err := parseCommentSpec(spec, prData.Head.SHA)
		if err != nil {
			return fmt.Errorf("invalid comment spec '%s': %w", spec, err)
		}
		comments = append(comments, comment)
	}

	// Create review payload
	reviewPayload := map[string]interface{}{
		"commit_id": prData.Head.SHA,
		"body":      body,
		"comments":  comments,
	}

	// Add event if specified (otherwise creates pending review)
	if event != "" {
		reviewPayload["event"] = event
	}

	payloadJSON, err := json.Marshal(reviewPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal review payload: %w", err)
	}

	if verbose {
		fmt.Printf("Review payload:\n%s\n\n", string(payloadJSON))
	}

	// Create the review
	var response map[string]interface{}
	err = client.Post(fmt.Sprintf("repos/%s/pulls/%d/reviews", repo, pr), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	// Display success message
	if event == "" {
		fmt.Printf("✅ Created pending review with %d comments on PR #%d\n", len(comments), pr)
		fmt.Printf("💡 Use 'gh pr review --approve/--request-changes/--comment' to submit the review\n")
	} else {
		fmt.Printf("✅ Created and submitted %s review with %d comments on PR #%d\n", event, len(comments), pr)
	}

	if verbose {
		if htmlURL, ok := response["html_url"].(string); ok {
			fmt.Printf("Review URL: %s\n", htmlURL)
		}
	}

	return nil
}

func parseCommentSpec(spec, commitSHA string) (map[string]interface{}, error) {
	// Format: "file:line:message" or "file:start:end:message"
	parts := strings.SplitN(spec, ":", 4)
	if len(parts) < 3 {
		return nil, fmt.Errorf("format should be 'file:line:message' or 'file:start:end:message'")
	}

	file := parts[0]

	if len(parts) == 3 {
		// Single line: file:line:message
		line, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid line number: %s", parts[1])
		}

		body := parts[2]
		if !noExpandSuggestionsReview {
			body = expandSuggestions(body)
		}

		return map[string]interface{}{
			"path": file,
			"line": line,
			"body": body,
		}, nil
	} else {
		// Range: file:start:end:message
		startLine, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid start line: %s", parts[1])
		}

		endLine, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid end line: %s", parts[2])
		}

		if startLine > endLine {
			return nil, fmt.Errorf("start line (%d) cannot be greater than end line (%d)", startLine, endLine)
		}

		body := parts[3]
		if !noExpandSuggestionsReview {
			body = expandSuggestions(body)
		}

		return map[string]interface{}{
			"path":       file,
			"line":       endLine,
			"start_line": startLine,
			"start_side": "RIGHT",
			"body":       body,
		}, nil
	}
}
