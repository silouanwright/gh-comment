package cmd

import (
	"fmt"
	"strings"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	// Client for dependency injection (tests can override)
	resolveClient github.GitHubAPI
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
	// Initialize client if not set (production use)
	if resolveClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		resolveClient = client
	}

	// Parse comment ID
	commentID, err := parsePositiveInt(args[0], "comment ID")
	if err != nil {
		return err
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
		fmt.Printf("PR Number: %d\n", pr)
		fmt.Printf("Comment ID: %d\n", commentID)
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would resolve conversation for comment #%d in PR #%d\n", commentID, pr)
		return nil
	}

	// Find the review thread for this comment
	threadID, err := resolveClient.FindReviewThreadForComment(owner, repoName, pr, commentID)
	if err != nil {
		return fmt.Errorf("failed to find review thread for comment: %w", err)
	}

	// Resolve the review thread
	err = resolveClient.ResolveReviewThread(threadID)
	if err != nil {
		return fmt.Errorf("failed to resolve conversation: %w", err)
	}

	fmt.Printf("✅ Resolved conversation for comment #%d\n", commentID)
	return nil
}
