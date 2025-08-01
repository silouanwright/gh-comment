package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve <comment-id>",
	Short: "Resolve a conversation thread",
	Long: `Resolve a conversation thread for a pull request review comment.

This marks the conversation as resolved, indicating that the feedback
has been addressed. Use the comment ID from 'gh comment list' output.

Examples:
  # Resolve a conversation
  gh comment resolve 2246362251

  # Resolve with dry-run
  gh comment resolve --dry-run 2246362251`,
	Args: cobra.ExactArgs(1),
	RunE: runResolve,
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}

func runResolve(cmd *cobra.Command, args []string) error {
	// Parse comment ID
	commentID, err := strconv.Atoi(args[0])
	if err != nil {
		return formatValidationError("comment ID", args[0], "must be a valid integer")
	}

	// Get repository and PR context
	repository, pr, err := getPRContext()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR Number: %d\n", pr)
		fmt.Printf("Comment ID: %d\n", commentID)
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would resolve conversation for comment #%d in PR #%d\n", commentID, pr)
		return nil
	}

	// Resolve the conversation using GraphQL
	err = resolveComment(repository, commentID, pr)
	if err != nil {
		return fmt.Errorf("failed to resolve conversation: %w", err)
	}

	fmt.Printf("✅ Resolved conversation for comment #%d\n", commentID)
	return nil
}
