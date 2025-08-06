package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"

	"github.com/spf13/cobra"

	"github.com/silouanwright/gh-comment/internal/github"
)

var (
	closePendingEvent  string
	closePendingBody   string
	closePendingClient github.GitHubAPI
)

var closePendingReviewCmd = &cobra.Command{
	Use:   "close-pending-review [pr] [body]",
	Short: "Close/submit a pending review created in GitHub's web interface",
	Long: heredoc.Doc(`
		Close and submit a pending review that was created in GitHub's web interface.

		IMPORTANT: This command only works with pending reviews created through GitHub's
		web interface. The GitHub API cannot create pending reviews - only the web UI can.

		This command finds your existing pending review on the PR and submits it with
		the specified event type (APPROVE, REQUEST_CHANGES, or COMMENT) and optional
		summary body.

		Once submitted, the pending review becomes visible to others and you can create
		new reviews.

		Note: This does NOT work with reviews created via 'gh comment review' commands,
		as those create submitted reviews immediately.
	`),
	Example: heredoc.Doc(`
		# Submit GUI-created pending review with approval
		$ gh comment close-pending-review 123 "LGTM! Great work" --event APPROVE

		# Submit with change requests
		$ gh comment close-pending-review 123 "Please address the comments" --event REQUEST_CHANGES

		# Submit as general comment
		$ gh comment close-pending-review 123 "Thanks for the updates" --event COMMENT

		# Submit with minimal body (auto-detect PR)
		$ gh comment close-pending-review "Looks good" --event APPROVE
	`),
	Args: cobra.RangeArgs(0, 2),
	RunE: runClosePendingReview,
}

func init() {
	rootCmd.AddCommand(closePendingReviewCmd)
	closePendingReviewCmd.Flags().StringVar(&closePendingEvent, "event", "COMMENT", "Review event (APPROVE|REQUEST_CHANGES|COMMENT) (default: COMMENT)")
	closePendingReviewCmd.Flags().StringVar(&closePendingBody, "body", "", "Review summary body (default: empty)")
}

func runClosePendingReview(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if closePendingClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		closePendingClient = client
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
	if closePendingBody != "" {
		body = closePendingBody
	}

	// Validate event type
	validEvents := []string{"APPROVE", "REQUEST_CHANGES", "COMMENT"}
	isValidEvent := false
	for _, validEvent := range validEvents {
		if closePendingEvent == validEvent {
			isValidEvent = true
			break
		}
	}
	if !isValidEvent {
		return fmt.Errorf("invalid event type: %s (must be APPROVE, REQUEST_CHANGES, or COMMENT)", closePendingEvent)
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
		fmt.Printf("Review event: %s\n", closePendingEvent)
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would close pending review on PR #%d:\n", pr)
		fmt.Printf("Body: %s\n", body)
		fmt.Printf("Event: %s\n", closePendingEvent)
		return nil
	}

	// Find the pending review
	reviewID, err := closePendingClient.FindPendingReview(owner, repoName, pr)
	if err != nil {
		return fmt.Errorf("failed to find pending review: %w", err)
	}

	if reviewID == 0 {
		return fmt.Errorf("no pending review found on PR #%d\n\nNote: This command only works with pending reviews created in GitHub's web interface.\nUse 'gh comment review' to create and submit reviews via CLI.", pr)
	}

	if verbose {
		fmt.Printf("Found pending review ID: %d\n", reviewID)
	}

	// Submit the review
	err = closePendingClient.SubmitReview(owner, repoName, pr, reviewID, body, closePendingEvent)
	if err != nil {
		return fmt.Errorf("failed to submit review: %w", err)
	}

	// Display success message
	eventText := ""
	switch closePendingEvent {
	case "APPROVE":
		eventText = "approved"
	case "REQUEST_CHANGES":
		eventText = "requested changes"
	case "COMMENT":
		eventText = "commented"
	}

	fmt.Printf("✅ Successfully submitted pending review and %s PR #%d\n", eventText, pr)
	return nil
}
