package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var (
	submitEvent string
	submitBody  string
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

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

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

	// Find and submit the pending review
	return submitPendingReview(repository, pr, body, submitEvent)
}

func submitPendingReview(repo string, pr int, body, event string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	// Find the pending review
	reviewID, err := findPendingReviewID(repo, pr)
	if err != nil {
		return err
	}

	if reviewID == 0 {
		return fmt.Errorf("no pending review found on PR #%d", pr)
	}

	if verbose {
		fmt.Printf("Found pending review ID: %d\n", reviewID)
	}

	// Submit the review
	submitPayload := map[string]interface{}{
		"body":  body,
		"event": event,
	}

	payloadJSON, err := json.Marshal(submitPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal submit payload: %w", err)
	}

	if verbose {
		fmt.Printf("Submit payload:\n%s\n\n", string(payloadJSON))
	}

	// Submit the review using the events endpoint
	var response map[string]interface{}
	err = client.Post(fmt.Sprintf("repos/%s/pulls/%d/reviews/%d/events", repo, pr, reviewID), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to submit review: %w", err)
	}

	// Display success message
	eventText := ""
	switch event {
	case "APPROVE":
		eventText = "approved"
	case "REQUEST_CHANGES":
		eventText = "requested changes"
	case "COMMENT":
		eventText = "commented"
	}

	fmt.Printf("✅ Successfully submitted review and %s PR #%d\n", eventText, pr)

	if verbose {
		if htmlURL, ok := response["html_url"].(string); ok {
			fmt.Printf("Review URL: %s\n", htmlURL)
		}
	}

	return nil
}

func findPendingReviewID(repo string, pr int) (int, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return 0, err
	}

	// Get existing reviews for this PR
	var reviews []map[string]interface{}
	err = client.Get(fmt.Sprintf("repos/%s/pulls/%d/reviews", repo, pr), &reviews)
	if err != nil {
		return 0, fmt.Errorf("failed to get reviews: %w", err)
	}

	if verbose {
		fmt.Printf("Found %d reviews:\n", len(reviews))
		for i, review := range reviews {
			state := "unknown"
			id := "unknown"
			if s, ok := review["state"].(string); ok {
				state = s
			}
			if reviewID, ok := review["id"].(float64); ok {
				id = fmt.Sprintf("%.0f", reviewID)
			}
			fmt.Printf("  Review %d: ID=%s, State=%s\n", i+1, id, state)
		}
	}

	// Look for an existing PENDING review
	for _, review := range reviews {
		if state, ok := review["state"].(string); ok && state == "PENDING" {
			if id, ok := review["id"].(float64); ok {
				return int(id), nil
			}
		}
	}

	return 0, nil // No pending review found
}
