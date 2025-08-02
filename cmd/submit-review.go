package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	submitEvent   string
	submitBody    string
	submitClient  github.GitHubAPI
)

var submitReviewCmd = &cobra.Command{
	Use:   "submit-review [pr] [body]",
	Short: "Submit a pending review",
	Long: `Submit a pending review with a summary message and approval status.

This command finds your pending review on the PR and submits it with the specified
event type (APPROVE, REQUEST_CHANGES, or COMMENT) and optional summary body.

Once submitted, the pending review becomes visible to others and you can create
new reviews.

Examples:
  # Submit pending review with approval
  gh comment submit-review 123 "LGTM! Great work" --event APPROVE

  # Submit with change requests
  gh comment submit-review 123 "Please address the comments" --event REQUEST_CHANGES

  # Submit as general comment
  gh comment submit-review 123 "Thanks for the updates" --event COMMENT

  # Submit with minimal body (auto-detect PR)
  gh comment submit-review "Looks good" --event APPROVE

  # Submit without additional body
  gh comment submit-review 123 --event APPROVE`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runSubmitReview,
}

func init() {
	rootCmd.AddCommand(submitReviewCmd)
	submitReviewCmd.Flags().StringVar(&submitEvent, "event", "COMMENT", "Review event: APPROVE, REQUEST_CHANGES, or COMMENT")
	submitReviewCmd.Flags().StringVar(&submitBody, "body", "", "Review summary body (optional)")
}

func runSubmitReview(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if submitClient == nil {
		submitClient = &github.RealClient{}
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
			_, prNum, err := getPRContext()
			if err != nil {
				return err
			}
			pr = prNum
			body = args[0]
		}
	} else {
		// Auto-detect PR
		_, prNum, err := getPRContext()
		if err != nil {
			return err
		}
		pr = prNum
	}

	// Use --body flag if provided, otherwise use positional arg
	if submitBody != "" {
		body = submitBody
	}

	// Validate event type
	validEvents := []string{"APPROVE", "REQUEST_CHANGES", "COMMENT"}
	isValidEvent := false
	for _, validEvent := range validEvents {
		if submitEvent == validEvent {
			isValidEvent = true
			break
		}
	}
	if !isValidEvent {
		return fmt.Errorf("invalid event type: %s (must be APPROVE, REQUEST_CHANGES, or COMMENT)", submitEvent)
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
		fmt.Printf("Review event: %s\n", submitEvent)
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would submit pending review on PR #%d:\n", pr)
		fmt.Printf("Body: %s\n", body)
		fmt.Printf("Event: %s\n", submitEvent)
		return nil
	}

	// Find the pending review
	reviewID, err := submitClient.FindPendingReview(owner, repoName, pr)
	if err != nil {
		return fmt.Errorf("failed to find pending review: %w", err)
	}

	if reviewID == 0 {
		return fmt.Errorf("no pending review found on PR #%d", pr)
	}

	if verbose {
		fmt.Printf("Found pending review ID: %d\n", reviewID)
	}

	// Submit the review
	err = submitClient.SubmitReview(owner, repoName, pr, reviewID, body, submitEvent)
	if err != nil {
		return fmt.Errorf("failed to submit review: %w", err)
	}

	// Display success message
	eventText := ""
	switch submitEvent {
	case "APPROVE":
		eventText = "approved"
	case "REQUEST_CHANGES":
		eventText = "requested changes"
	case "COMMENT":
		eventText = "commented"
	}

	fmt.Printf("✅ Successfully submitted review and %s PR #%d\n", eventText, pr)
	return nil
}
