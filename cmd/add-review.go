package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	reviewBody                string
	reviewComments            []string
	reviewEvent               string
	noExpandSuggestionsReview bool

	// Client for dependency injection (tests can override)
	addReviewClient github.GitHubAPI
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
	// Initialize client if not set (production use)
	if addReviewClient == nil {
		addReviewClient = &github.RealClient{}
	}

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

	// Parse owner/repo
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repository)
	}
	owner, repoName := parts[0], parts[1]

	// Create the review with all comments using the client
	return createReviewWithComments(addReviewClient, owner, repoName, pr, body, reviewEvent, reviewComments)
}

func createReviewWithComments(client github.GitHubAPI, owner, repo string, pr int, body, event string, commentSpecs []string) error {
	// Get PR details for commit SHA
	prDetails, err := client.GetPRDetails(owner, repo, pr)
	if err != nil {
		return fmt.Errorf("failed to get PR details: %w", err)
	}

	// Extract commit SHA from PR details
	var commitSHA string
	if head, ok := prDetails["head"].(map[string]interface{}); ok {
		if sha, ok := head["sha"].(string); ok {
			commitSHA = sha
		} else {
			return fmt.Errorf("could not extract commit SHA from PR details")
		}
	} else {
		return fmt.Errorf("could not extract head information from PR details")
	}

	// Parse comment specifications
	var reviewComments []github.ReviewCommentInput
	for _, spec := range commentSpecs {
		comment, err := parseCommentSpec(spec, commitSHA)
		if err != nil {
			return fmt.Errorf("invalid comment spec '%s': %w", spec, err)
		}
		reviewComments = append(reviewComments, comment)
	}

	// Create review input
	reviewInput := github.ReviewInput{
		Body:     body,
		Comments: reviewComments,
	}

	// Set event if specified (otherwise creates pending review)
	if event != "" {
		reviewInput.Event = event
	} else {
		reviewInput.Event = "COMMENT" // Default for pending review
	}

	// Create the review
	err = client.CreateReview(owner, repo, pr, reviewInput)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	// Display success message
	if event == "" {
		fmt.Printf("✅ Created pending review with %d comments on PR #%d\n", len(reviewComments), pr)
		fmt.Printf("💡 Use 'gh pr review --approve/--request-changes/--comment' to submit the review\n")
	} else {
		fmt.Printf("✅ Created and submitted %s review with %d comments on PR #%d\n", event, len(reviewComments), pr)
	}

	return nil
}

func parseCommentSpec(spec, commitSHA string) (github.ReviewCommentInput, error) {
	// Format: "file:line:message" or "file:start:end:message"
	parts := strings.Split(spec, ":")
	if len(parts) < 3 {
		return github.ReviewCommentInput{}, fmt.Errorf("format should be 'file:line:message' or 'file:start:end:message'")
	}

	file := parts[0]

	// Try to parse as range first (file:start:end:message...)
	if len(parts) >= 4 {
		startLine, err := strconv.Atoi(parts[1])
		if err == nil {
			endLine, err := strconv.Atoi(parts[2])
			if err == nil {
				// This is a valid range format
				if startLine > endLine {
					return github.ReviewCommentInput{}, fmt.Errorf("start line (%d) cannot be greater than end line (%d)", startLine, endLine)
				}

				// Join remaining parts as the message (in case message contains colons)
				body := strings.Join(parts[3:], ":")
				if !noExpandSuggestionsReview {
					body = expandSuggestions(body)
				}

				return github.ReviewCommentInput{
					Path:      file,
					Line:      endLine,
					StartLine: startLine,
					Side:      "RIGHT",
					Body:      body,
					CommitID:  commitSHA,
				}, nil
			}
		}
	}

	// Parse as single line: file:line:message...
	line, err := strconv.Atoi(parts[1])
	if err != nil {
		return github.ReviewCommentInput{}, fmt.Errorf("invalid line number: %s", parts[1])
	}

	// Join remaining parts as the message (in case message contains colons)
	body := strings.Join(parts[2:], ":")
	if !noExpandSuggestionsReview {
		body = expandSuggestions(body)
	}

	return github.ReviewCommentInput{
		Path:     file,
		Line:     line,
		Body:     body,
		CommitID: commitSHA,
	}, nil
}
