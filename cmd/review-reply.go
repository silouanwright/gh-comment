package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	resolveConversationReviewReply bool
	noExpandSuggestionsReviewReply bool

	// Client for dependency injection (tests can override)
	reviewReplyClient github.GitHubAPI
)

var reviewReplyCmd = &cobra.Command{
	Use:   "review-reply <comment-id> [message]",
	Short: "Reply to a review comment with a text message",
	Long: heredoc.Doc(`
		Reply to an existing review comment with a text message.

		This command is specifically for review comment threading - replying to
		line-specific comments that appear in the "Files Changed" tab. Review
		comment threading has limited GitHub API support.

		For general PR discussion, use 'gh comment add' instead.
		For emoji reactions, use 'gh comment react' command.

		Comment IDs can be found in the output of 'gh comment list'.
	`),
	Example: heredoc.Doc(`
		# Reply to review comment with message
		$ gh comment review-reply 789012 "Fixed this issue"

		# Reply and resolve conversation in one operation
		$ gh comment review-reply 789012 "Addressed your feedback" --resolve

		# Just resolve conversation without adding message
		$ gh comment review-reply 789012 --resolve

		# For general PR discussion, use add command instead:
		# $ gh comment add 123 "Thanks for the review!"

		# For emoji reactions, use react command instead:
		# $ gh comment react 789012 +1
	`),
	Args: cobra.RangeArgs(1, 2),
	RunE: runReviewReply,
}

func init() {
	rootCmd.AddCommand(reviewReplyCmd)

	reviewReplyCmd.Flags().BoolVar(&resolveConversationReviewReply, "resolve", false, "Resolve the conversation after replying")
	reviewReplyCmd.Flags().BoolVar(&noExpandSuggestionsReviewReply, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax")
}

func runReviewReply(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if reviewReplyClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		reviewReplyClient = client
	}

	// Parse comment ID
	commentIDStr := args[0]
	commentID, err := parsePositiveInt(commentIDStr, "comment ID")
	if err != nil {
		return err
	}

	// Get message if provided
	var message string
	if len(args) > 1 {
		message = args[1]
	}

	// Validate that we have either message or resolve
	if message == "" && !resolveConversationReviewReply {
		return fmt.Errorf("must provide either a message or --resolve")
	}

	// Validate comment body length if provided
	if message != "" {
		if err := validateCommentBody(message); err != nil {
			return err
		}
	}

	// Get repository context
	repository, prNumber, err := getPRContext()
	if err != nil {
		return err
	}

	// Validate repository name
	if err := validateRepositoryName(repository); err != nil {
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
		fmt.Printf("Comment ID: %d\n", commentID)
		fmt.Printf("Comment Type: review\n")
		if message != "" {
			fmt.Printf("Message: %s\n", message)
		}
		if resolveConversationReviewReply {
			fmt.Printf("Resolve Conversation: true\n")
		}
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would reply to review comment #%d:\n", commentID)
		if message != "" {
			fmt.Printf("Message: %s\n", message)
		}
		if resolveConversationReviewReply {
			fmt.Printf("Resolve Conversation: true\n")
		}
		return nil
	}

	// Handle message reply
	if message != "" {
		// Expand suggestions if enabled
		if !noExpandSuggestionsReviewReply {
			message = expandSuggestions(message)
		}

		// Reply to review comment
		_, err = reviewReplyClient.CreateReviewCommentReply(owner, repoName, commentID, message)
		if err != nil {
			return fmt.Errorf("failed to create reply: %w", err)
		}
		fmt.Printf("✅ Replied to review comment #%d: %s\n", commentID, message)
	}

	// Handle resolve conversation
	if resolveConversationReviewReply {
		// Find the thread ID for this comment
		threadID, err := reviewReplyClient.FindReviewThreadForComment(owner, repoName, prNumber, commentID)
		if err != nil {
			return fmt.Errorf("failed to find review thread: %w", err)
		}

		// Resolve the thread
		err = reviewReplyClient.ResolveReviewThread(threadID)
		if err != nil {
			return fmt.Errorf("failed to resolve conversation: %w", err)
		}
		fmt.Printf("✅ Resolved conversation for review comment #%d\n", commentID)
	}

	return nil
}
