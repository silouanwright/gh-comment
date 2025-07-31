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
	reaction string
)

var replyCmd = &cobra.Command{
	Use:   "reply <comment-id> [message]",
	Short: "Reply to a specific comment on a PR",
	Long: `Reply to a specific comment on a pull request review.

You can reply with a message, add a reaction, or both. Use the comment ID from the URL 
shown in 'gh comment list --verbose' output.

Common use cases:
- Acknowledge feedback: "Good point, thanks!"
- Ask for clarification: "What do you mean by this?"
- Confirm fix: "Fixed in latest commit"
- Show appreciation with reactions

Examples:
  # Reply with a message
  gh comment reply 2246362251 "Good catch, fixed this!"
  
  # Add a thumbs up reaction
  gh comment reply 2246362251 --reaction +1
  
  # Reply with message and reaction
  gh comment reply 2246362251 "Thanks for the feedback!" --reaction heart
  
  # Quick acknowledgment
  gh comment reply 2246362251 "👍 Fixed"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runReply,
}

func init() {
	rootCmd.AddCommand(replyCmd)
	
	replyCmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
}

func runReply(cmd *cobra.Command, args []string) error {
	// Parse comment ID
	commentIDStr := args[0]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %s", commentIDStr)
	}

	// Get message if provided
	var message string
	if len(args) > 1 {
		message = args[1]
	}

	// Validate that we have either message or reaction
	if message == "" && reaction == "" {
		return fmt.Errorf("must provide either a message or --reaction")
	}

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("Comment ID: %d\n", commentID)
		if message != "" {
			fmt.Printf("Message: %s\n", message)
		}
		if reaction != "" {
			fmt.Printf("Reaction: %s\n", reaction)
		}
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would reply to comment #%d:\n", commentID)
		if message != "" {
			fmt.Printf("Message: %s\n", message)
		}
		if reaction != "" {
			fmt.Printf("Reaction: %s\n", reaction)
		}
		return nil
	}

	// Add reaction if specified
	if reaction != "" {
		err = addReaction(repository, commentID, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}
		fmt.Printf("✅ Added %s reaction to comment #%d\n", reaction, commentID)
	}

	// Add reply message if specified
	if message != "" {
		err = addReply(repository, commentID, message)
		if err != nil {
			return fmt.Errorf("failed to add reply: %w", err)
		}
		fmt.Printf("✅ Replied to comment #%d\n", commentID)
	}

	return nil
}

func addReaction(repo string, commentID int, reactionType string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// GitHub API endpoint for adding reactions to pull request review comments
	payload := map[string]interface{}{
		"content": reactionType,
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var response map[string]interface{}
	err = client.Post(fmt.Sprintf("repos/%s/pulls/comments/%d/reactions", repo, commentID), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	if verbose {
		fmt.Printf("Reaction added successfully\n")
	}

	return nil
}

func addReply(repo string, commentID int, message string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Get the original comment to extract PR info and commit SHA
	var originalComment struct {
		PullRequestURL string `json:"pull_request_url"`
		Path           string `json:"path"`
		CommitID       string `json:"commit_id"`
		Line           int    `json:"line"`
		StartLine      int    `json:"start_line"`
	}

	err = client.Get(fmt.Sprintf("repos/%s/pulls/comments/%d", repo, commentID), &originalComment)
	if err != nil {
		return fmt.Errorf("failed to get comment details: %w", err)
	}

	// Extract PR number from URL
	urlParts := strings.Split(originalComment.PullRequestURL, "/")
	prNumber := urlParts[len(urlParts)-1]

	// Create a reply using the pull request review comments API
	// This will create a threaded reply within the same review conversation
	payload := map[string]interface{}{
		"body":         message,
		"commit_id":    originalComment.CommitID,
		"path":         originalComment.Path,
		"line":         originalComment.Line,
		"in_reply_to":  commentID, // This makes it a threaded reply
	}

	// Add start_line if it's a range comment
	if originalComment.StartLine > 0 && originalComment.StartLine != originalComment.Line {
		payload["start_line"] = originalComment.StartLine
		payload["start_side"] = "RIGHT"
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var response map[string]interface{}
	err = client.Post(fmt.Sprintf("repos/%s/pulls/%s/comments", repo, prNumber), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to add reply: %w", err)
	}

	if verbose {
		fmt.Printf("Reply added successfully to review thread\n")
	}

	return nil
}

func validateReaction(reaction string) bool {
	validReactions := []string{"+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", "eyes"}
	for _, valid := range validReactions {
		if reaction == valid {
			return true
		}
	}
	return false
}
