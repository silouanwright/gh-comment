package cmd

import (
	"fmt"
	"strconv"
	"strings"

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
	Short: "Reply to a specific comment on a PR",
	Long: `Reply to a specific comment on a pull request.

You can reply with a message, add/remove reactions, or both. Use the comment ID from the URL
shown in 'gh comment list' output. Specify the comment type using --type flag:
- 'review' for line-specific code review comments (default)
- 'issue' for general PR discussion comments

Common use cases:
- Acknowledge feedback: "Good point, thanks!"
- Ask for clarification: "What do you mean by this?"
- Confirm fix: "Fixed in latest commit"
- Show appreciation with reactions
- Remove accidental or outdated reactions
- Resolve conversations after addressing feedback

Examples:
  # Reply to a review comment (line-specific)
  gh comment reply 2246362251 "Good catch, fixed this!"

  # Reply to an issue comment (general PR comment)
  gh comment reply 3141344022 "Thanks for the feedback!" --type issue

  # Reply and resolve conversation (review comments only)
  gh comment reply 2246362251 "Fixed in latest commit" --resolve

  # Add a thumbs up reaction
  gh comment reply 2246362251 --reaction +1

  # Remove a reaction
  gh comment reply 2246362251 --remove-reaction +1

  # Reply with message, reaction, and resolve
  gh comment reply 2246362251 "Thanks for the feedback!" --reaction heart --resolve

  # Quick acknowledgment
  gh comment reply 2246362251 "👍 Fixed"`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runReply,
}

func init() {
	rootCmd.AddCommand(replyCmd)

	replyCmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	replyCmd.Flags().StringVar(&removeReaction, "remove-reaction", "", "Remove reaction: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	replyCmd.Flags().BoolVar(&resolveConversation, "resolve", false, "Resolve the conversation after replying")
	replyCmd.Flags().BoolVar(&noExpandSuggestionsReply, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax")
	replyCmd.Flags().StringVar(&commentType, "type", "review", "Comment type: 'issue' for general PR comments, 'review' for line-specific comments (default: review)")
}

func runReply(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if replyClient == nil {
		replyClient = &github.RealClient{}
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
	repository, _, err := getPRContext()
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
		err = replyClient.AddReaction(owner, repoName, commentID, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}
		fmt.Printf("✅ Added %s reaction to comment #%d\n", reaction, commentID)
	}

	// Remove reaction if specified
	if removeReaction != "" {
		err = replyClient.RemoveReaction(owner, repoName, commentID, removeReaction)
		if err != nil {
			return fmt.Errorf("failed to remove reaction: %w", err)
		}
		fmt.Printf("✅ Removed %s reaction from comment #%d\n", removeReaction, commentID)
	}

	// TODO: Refactor message reply functionality
	if message != "" {
		return fmt.Errorf("message replies not yet refactored - use reactions for now")
	}

	// TODO: Refactor resolve conversation functionality
	if resolveConversation {
		return fmt.Errorf("resolve conversation not yet refactored - use reactions for now")
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
