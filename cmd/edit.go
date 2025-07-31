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
	editMessages []string
)

var editCmd = &cobra.Command{
	Use:   "edit <comment-id> [message]",
	Short: "Edit an existing comment on a PR",
	Long: `Edit an existing comment on a pull request.

You can edit with a new message using either positional argument or --message flags.
Use the comment ID from the URL shown in 'gh comment list' output.

Common use cases:
- Fix typos in comments: "Fixed typo in previous comment"
- Add more context: "Adding more details about the implementation"
- Refine AI-generated comments: "Updating comment based on new analysis"
- Correct mistakes: "Correcting the suggested approach"

Examples:
  # Edit with new message
  gh comment edit 2246362251 "Updated comment with better explanation"
  
  # Edit with multi-line content using --message flags (AI-friendly)
  gh comment edit 2246362251 --message "First paragraph" --message "Second paragraph"
  
  # Edit with multi-line content (shell native)
  gh comment edit 2246362251 "Line 1
Line 2
Line 3"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringArrayVarP(&editMessages, "message", "m", []string{}, "Edit message (can be used multiple times for multi-line comments)")
}

func runEdit(cmd *cobra.Command, args []string) error {
	// Parse comment ID
	commentID, err := strconv.Atoi(args[0])
	if err != nil {
		return formatValidationError("comment ID", args[0], "must be a valid integer")
	}

	var message string
	
	// Handle message from positional arg or --message flags
	if len(args) == 2 {
		message = args[1]
	} else if len(editMessages) > 0 {
		message = strings.Join(editMessages, "\n")
	} else {
		return fmt.Errorf("must provide either a message argument or --message flags")
	}

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("Comment ID: %d\n", commentID)
		fmt.Printf("New message: %s\n", message)
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would edit comment #%d:\n", commentID)
		fmt.Printf("New message: %s\n", message)
		return nil
	}

	// Edit the comment
	err = editComment(repository, commentID, message)
	if err != nil {
		return fmt.Errorf("failed to edit comment: %w", err)
	}

	fmt.Printf("✅ Edited comment #%d\n", commentID)
	return nil
}

func editComment(repo string, commentID int, newMessage string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	// GitHub API endpoint for editing pull request review comments
	payload := map[string]interface{}{
		"body": newMessage,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var response map[string]interface{}
	err = client.Patch(fmt.Sprintf("repos/%s/pulls/comments/%d", repo, commentID), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to edit comment: %w", err)
	}

	return nil
}
