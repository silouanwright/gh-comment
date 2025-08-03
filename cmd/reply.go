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
	reaction                 string
	removeReaction           string
	resolveConversation      bool
	noExpandSuggestionsReply bool
	commentType              string

	// Client for dependency injection (tests can override)
	replyClient github.GitHubAPI
)

var replyCmd = &cobra.Command{
	Use:   "reply <comment-id> [message]",
	Short: "Reply to a comment or add reactions",
	Long: heredoc.Doc(`
		Reply to an existing comment with a message or reaction.

		Comment Types:
		- Issue comments: General PR discussion, support message replies
		- Review comments: Line-specific feedback, only reactions work for replies

		Note: GitHub API only supports message replies for issue comments.
		For review comments, use reactions instead of messages.

		Comment IDs can be found in the output of 'gh comment list'.
	`),
	Example: heredoc.Doc(`
		# Reply to general PR discussion with message (default)
		$ gh comment reply 123456 "Good point, I'll fix that"

		# Add reactions to any comment type
		$ gh comment reply 123456 --reaction +1
		$ gh comment reply 789012 --reaction heart

		# For review comments, use reactions (messages not supported)
		$ gh comment reply 789012 --reaction +1 --type review
		$ gh comment reply 789012 --reaction rocket --resolve --type review

		# Remove a reaction
		$ gh comment reply 123456 --remove-reaction +1
	`),
	Args: cobra.RangeArgs(1, 2),
	RunE: runReply,
}

func init() {
	rootCmd.AddCommand(replyCmd)

	replyCmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	replyCmd.Flags().StringVar(&removeReaction, "remove-reaction", "", "Remove reaction: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	replyCmd.Flags().BoolVar(&resolveConversation, "resolve", false, "Resolve the conversation after replying")
	replyCmd.Flags().BoolVar(&noExpandSuggestionsReply, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax")
	replyCmd.Flags().StringVar(&commentType, "type", "issue", "Comment type: 'issue' for general PR comments, 'review' for line-specific comments (default: issue)")
}

func runReply(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if replyClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		replyClient = client
	}

	// Parse comment ID
	commentIDStr := args[0]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return formatValidationError("comment ID", commentIDStr, "must be a valid integer")
	}

	// We'll get repository from getPRContext below

	// Get message if provided
	var message string
	if len(args) > 1 {
		message = args[1]
	}

	// Validate that we have either message, reaction, remove-reaction, or resolve
	if message == "" && reaction == "" && removeReaction == "" && !resolveConversation {
		return fmt.Errorf("must provide either a message, --reaction, --remove-reaction, or --resolve")
	}

	// Validate that we don't have both reaction and remove-reaction
	if reaction != "" && removeReaction != "" {
		return fmt.Errorf("cannot use both --reaction and --remove-reaction at the same time")
	}

	// Validate reaction if provided
	if reaction != "" && !validateReaction(reaction) {
		return formatValidationError("reaction", reaction, "must be one of: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	}

	// Validate remove-reaction if provided
	if removeReaction != "" && !validateReaction(removeReaction) {
		return formatValidationError("remove-reaction", removeReaction, "must be one of: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	}

	// Validate comment type
	if commentType != "issue" && commentType != "review" {
		return formatValidationError("type", commentType, "must be either 'issue' or 'review'")
	}

	// Get repository context
	repository, prNumber, err := getPRContext()
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
		fmt.Printf("Comment ID: %d\n", commentID)
		fmt.Printf("Comment Type: %s\n", commentType)
		if message != "" {
			fmt.Printf("Message: %s\n", message)
		}
		if reaction != "" {
			fmt.Printf("Reaction: %s\n", reaction)
		}
		if removeReaction != "" {
			fmt.Printf("Remove Reaction: %s\n", removeReaction)
		}
		if resolveConversation {
			fmt.Printf("Resolve Conversation: true\n")
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
		if removeReaction != "" {
			fmt.Printf("Remove Reaction: %s\n", removeReaction)
		}
		if resolveConversation {
			fmt.Printf("Resolve Conversation: true\n")
		}
		return nil
	}

	// Add reaction if specified
	if reaction != "" {
		err = replyClient.AddReaction(owner, repoName, commentID, prNumber, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}
		fmt.Printf("✅ Added %s reaction to comment #%d\n", reaction, commentID)
	}

	// Remove reaction if specified
	if removeReaction != "" {
		err = replyClient.RemoveReaction(owner, repoName, commentID, prNumber, removeReaction)
		if err != nil {
			return fmt.Errorf("failed to remove reaction: %w", err)
		}
		fmt.Printf("✅ Removed %s reaction from comment #%d\n", removeReaction, commentID)
	}

	// Handle message reply
	if message != "" {
		// Expand suggestions if enabled
		if !noExpandSuggestionsReply {
			message = expandSuggestions(message)
		}

		var err error
		if commentType == "issue" {
			// For issue comments, create a new issue comment (GitHub API doesn't support direct replies to issue comments)
			// We need the PR number for this - let's get it from context
			_, prNum, err := getPRContext()
			if err != nil {
				return fmt.Errorf("failed to get PR context: %w", err)
			}
			_, err = replyClient.CreateIssueComment(owner, repoName, prNum, message)
		} else {
			// Reply to review comment
			_, err = replyClient.CreateReviewCommentReply(owner, repoName, commentID, message)
		}

		if err != nil {
			return fmt.Errorf("failed to create reply: %w", err)
		}
		fmt.Printf("✅ Replied to comment #%d: %s\n", commentID, message)
	}

	// Handle resolve conversation
	if resolveConversation {
		if commentType == "issue" {
			return fmt.Errorf("cannot resolve issue comments - only review comments can be resolved")
		}

		// Find the thread ID for this comment
		threadID, err := replyClient.FindReviewThreadForComment(owner, repoName, 0, commentID) // PR number not needed for this operation
		if err != nil {
			return fmt.Errorf("failed to find review thread: %w", err)
		}

		// Resolve the thread
		err = replyClient.ResolveReviewThread(threadID)
		if err != nil {
			return fmt.Errorf("failed to resolve conversation: %w", err)
		}
		fmt.Printf("✅ Resolved conversation for comment #%d\n", commentID)
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
