package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	// Client for dependency injection (tests can override)
	reviewClient github.GitHubAPI

	// Review-specific flags
	reviewEventFlag    string
	reviewCommentsFlag []string
)

var reviewCmd = &cobra.Command{
	Use:   "review <pr> <body>",
	Short: "Create a review with multiple comments",
	Long: heredoc.Doc(`
		Create a review with multiple comments using a streamlined interface.

		This command provides a simplified way to create reviews with multiple comments
		using command-line flags. Perfect for comprehensive code reviews where you
		want to add several comments and submit a review decision in one operation.
	`),
	Example: heredoc.Doc(`
		# Security-focused comprehensive review
		$ gh comment review 123 "Security audit complete - critical issues found" \
		  --comment auth.go:67:"Use crypto.randomBytes(32) instead of Math.random() for token generation" \
		  --comment api.js:134:140:"This endpoint lacks rate limiting - vulnerable to DoS attacks" \
		  --comment validation.js:25:"Input sanitization missing - SQL injection risk" \
		  --event REQUEST_CHANGES

		# Performance optimization review
		$ gh comment review 123 "Performance review - optimization opportunities identified" \
		  --comment database.py:89:95:"Extract this N+1 query to a single batch operation" \
		  --comment cache.js:156:"Consider Redis clustering for this high-traffic endpoint" \
		  --comment monitoring.go:78:"Add performance metrics for this critical path" \
		  --event COMMENT

		# Architecture migration approval
		$ gh comment review 123 "Migration to microservices architecture approved" \
		  --comment service-layer.js:45:"Excellent separation of concerns in the new service layer" \
		  --comment api-gateway.go:123:130:"API gateway implementation follows best practices" \
		  --comment docker-compose.yml:67:"Container orchestration setup looks solid" \
		  --event APPROVE

		# Code quality and maintainability review
		$ gh comment review 123 "Code quality review - refactoring needed" \
		  --comment legacy-handler.js:200:250:"This function is doing too much - extract into separate services" \
		  --comment utils.go:45:"Consider using dependency injection pattern here" \
		  --comment test-helpers.js:89:"Add integration tests for this critical business logic" \
		  --event REQUEST_CHANGES
	`),
	Args: cobra.RangeArgs(1, 2),
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.Flags().StringVar(&reviewEventFlag, "event", "COMMENT", "Review event: APPROVE, REQUEST_CHANGES, or COMMENT")
	reviewCmd.Flags().StringArrayVar(&reviewCommentsFlag, "comment", []string{}, "Add comment in format file:line:message or file:start-end:message")
}

func runReview(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if reviewClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		reviewClient = client
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
			_, pr, err = getPRContext()
			if err != nil {
				return err
			}
			body = args[0]
		}
	}

	// Validate event type
	validEvents := []string{"APPROVE", "REQUEST_CHANGES", "COMMENT"}
	isValidEvent := false
	for _, validEvent := range validEvents {
		if reviewEventFlag == validEvent {
			isValidEvent = true
			break
		}
	}
	if !isValidEvent {
		return fmt.Errorf("invalid event type: %s (must be APPROVE, REQUEST_CHANGES, or COMMENT)", reviewEventFlag)
	}

	// Validate that we have either a body or comments
	if body == "" && len(reviewCommentsFlag) == 0 {
		return fmt.Errorf("review must have either a body message or comments (use --comment)")
	}

	// Get repository and PR context
	repository, pr, err := getPRContext()
	if err != nil {
		return err
	}

	// Parse owner/repo
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repository)
	}
	owner, repoName := parts[0], parts[1]

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR: %d\n", pr)
		fmt.Printf("Review body: %s\n", body)
		fmt.Printf("Review event: %s\n", reviewEventFlag)
		fmt.Printf("Comments: %d\n", len(reviewCommentsFlag))
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would create review on PR #%d:\n", pr)
		fmt.Printf("Body: %s\n", body)
		fmt.Printf("Event: %s\n", reviewEventFlag)
		fmt.Printf("Comments: %d\n", len(reviewCommentsFlag))
		for i, comment := range reviewCommentsFlag {
			fmt.Printf("  %d. %s\n", i+1, comment)
		}
		return nil
	}

	// Get PR details for commit SHA
	prDetails, err := reviewClient.GetPRDetails(owner, repoName, pr)
	if err != nil {
		return fmt.Errorf("failed to get PR details: %w", err)
	}

	headSHA, ok := prDetails["head"].(map[string]interface{})["sha"].(string)
	if !ok {
		return fmt.Errorf("failed to get commit SHA from PR details")
	}

	// Parse and create review comments
	var reviewCommentInputs []github.ReviewCommentInput
	for i, commentSpec := range reviewCommentsFlag {
		commentInput, err := parseReviewCommentSpec(commentSpec, headSHA)
		if err != nil {
			return fmt.Errorf("invalid comment %d (%s): %w", i+1, commentSpec, err)
		}
		reviewCommentInputs = append(reviewCommentInputs, commentInput)
	}

	// Create the review
	review := github.ReviewInput{
		Body:     body,
		Event:    reviewEventFlag,
		Comments: reviewCommentInputs,
	}

	err = reviewClient.CreateReview(owner, repoName, pr, review)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	// Display success message
	eventText := ""
	switch reviewEventFlag {
	case "APPROVE":
		eventText = "approved"
	case "REQUEST_CHANGES":
		eventText = "requested changes"
	case "COMMENT":
		eventText = "commented on"
	}

	fmt.Printf("✅ Successfully created review and %s PR #%d", eventText, pr)
	if len(reviewCommentInputs) > 0 {
		fmt.Printf(" with %d comments", len(reviewCommentInputs))
	}
	fmt.Println()

	return nil
}

// parseReviewCommentSpec parses a comment specification in the format:
// file:line:message or file:start-end:message
func parseReviewCommentSpec(spec, commitSHA string) (github.ReviewCommentInput, error) {
	parts := strings.SplitN(spec, ":", 3)
	if len(parts) < 3 {
		return github.ReviewCommentInput{}, fmt.Errorf("format must be file:line:message or file:start-end:message")
	}

	filePath := parts[0]
	lineSpec := parts[1]
	message := strings.Join(parts[2:], ":") // Rejoin in case message contains colons

	if filePath == "" {
		return github.ReviewCommentInput{}, fmt.Errorf("file path cannot be empty")
	}
	if message == "" {
		return github.ReviewCommentInput{}, fmt.Errorf("message cannot be empty")
	}

	comment := github.ReviewCommentInput{
		Body:     expandSuggestions(message),
		Path:     filePath,
		CommitID: commitSHA,
	}

	// Parse line specification (single line or range)
	if strings.Contains(lineSpec, "-") {
		// Range format: start-end
		rangeParts := strings.Split(lineSpec, "-")
		if len(rangeParts) != 2 {
			return github.ReviewCommentInput{}, fmt.Errorf("range format must be start-end")
		}

		startLine, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
		if err != nil {
			return github.ReviewCommentInput{}, fmt.Errorf("invalid start line: %w", err)
		}

		endLine, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
		if err != nil {
			return github.ReviewCommentInput{}, fmt.Errorf("invalid end line: %w", err)
		}

		if startLine <= 0 || endLine <= 0 {
			return github.ReviewCommentInput{}, fmt.Errorf("line numbers must be positive")
		}

		if startLine > endLine {
			return github.ReviewCommentInput{}, fmt.Errorf("start line must be <= end line")
		}

		comment.StartLine = startLine
		comment.Line = endLine
	} else {
		// Single line format
		line, err := strconv.Atoi(lineSpec)
		if err != nil {
			return github.ReviewCommentInput{}, fmt.Errorf("invalid line number: %w", err)
		}

		if line <= 0 {
			return github.ReviewCommentInput{}, fmt.Errorf("line number must be positive")
		}

		comment.Line = line
	}

	return comment, nil
}
