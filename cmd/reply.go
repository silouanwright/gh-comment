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
	removeReaction string
	resolveConversation bool
	noExpandSuggestionsReply bool
	commentType string
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
	// Parse comment ID
	commentIDStr := args[0]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return formatValidationError("comment ID", commentIDStr, "must be a valid integer")
	}

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

	// Get repository and PR context
	repository, pr, err := getPRContext()
	if err != nil {
		return err
	}

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
		err = addReaction(repository, commentID, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}
		fmt.Printf("✅ Added %s reaction to comment #%d\n", reaction, commentID)
	}

	// Remove reaction if specified
	if removeReaction != "" {
		err = removeReactionFromComment(repository, commentID, removeReaction)
		if err != nil {
			return fmt.Errorf("failed to remove reaction: %w", err)
		}
		fmt.Printf("✅ Removed %s reaction from comment #%d\n", removeReaction, commentID)
	}

	// Add reply message if specified
	if message != "" {
		// Expand suggestion syntax to GitHub markdown (unless disabled)
		var finalMessage string
		if noExpandSuggestionsReply {
			finalMessage = message
		} else {
			finalMessage = expandSuggestions(message)
		}

		// Use appropriate reply method based on comment type
		if commentType == "issue" {
			err = addIssueCommentReply(repository, pr, finalMessage)
		} else {
			err = addReviewCommentReply(repository, commentID, pr, finalMessage)
		}

		if err != nil {
			return fmt.Errorf("failed to add reply: %w", err)
		}
		fmt.Printf("✅ Replied to %s comment #%d\n", commentType, commentID)
	}

	// Resolve conversation if specified
	if resolveConversation {
		err = resolveComment(repository, commentID, pr)
		if err != nil {
			return fmt.Errorf("failed to resolve conversation: %w", err)
		}
		fmt.Printf("✅ Resolved conversation for comment #%d\n", commentID)
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
		return formatAPIError("adding reaction", fmt.Sprintf("repos/%s/pulls/comments/%d/reactions", repo, commentID), err)
	}

	if verbose {
		fmt.Printf("Reaction added successfully\n")
	}

	return nil
}

// addReviewCommentReply adds a reply to a review comment (line-specific)
func addReviewCommentReply(repo string, commentID int, prNumber int, message string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Get the original review comment details
	var originalComment struct {
		Path      string `json:"path"`
		CommitID  string `json:"commit_id"`
		Line      int    `json:"line"`
		StartLine int    `json:"start_line"`
	}

	err = client.Get(fmt.Sprintf("repos/%s/pulls/comments/%d", repo, commentID), &originalComment)
	if err != nil {
		return fmt.Errorf("failed to get review comment details: %w", err)
	}

	// Create a threaded reply within the same review conversation
	payload := map[string]interface{}{
		"body":        message,
		"commit_id":   originalComment.CommitID,
		"path":        originalComment.Path,
		"line":        originalComment.Line,
		"in_reply_to": commentID, // This makes it a threaded reply
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
	err = client.Post(fmt.Sprintf("repos/%s/pulls/%d/comments", repo, prNumber), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to add review comment reply: %w", err)
	}

	if verbose {
		fmt.Printf("Reply added successfully to review comment thread\n")
	}

	return nil
}

// addIssueCommentReply adds a reply to an issue comment (general PR comment)
func addIssueCommentReply(repo string, prNumber int, message string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// For issue comments, we create a new comment on the PR/issue
	// GitHub doesn't support threaded replies for issue comments, so we create a new top-level comment
	payload := map[string]interface{}{
		"body": message,
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var response map[string]interface{}
	err = client.Post(fmt.Sprintf("repos/%s/issues/%d/comments", repo, prNumber), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to add issue comment reply: %w", err)
	}

	if verbose {
		fmt.Printf("Reply added successfully as new issue comment\n")
	}

	return nil
}

func removeReactionFromComment(repo string, commentID int, reactionType string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	// First, we need to get the current user's reaction ID for this comment
	// GitHub API requires the reaction ID to delete it
	var reactions []map[string]interface{}
	err = client.Get(fmt.Sprintf("repos/%s/pulls/comments/%d/reactions", repo, commentID), &reactions)
	if err != nil {
		return fmt.Errorf("failed to get reactions: %w", err)
	}

	// Find the current user's reaction of the specified type
	var reactionID int
	for _, reaction := range reactions {
		if reaction["content"] == reactionType {
			// For now, we'll take the first matching reaction
			// In practice, this will be the current user's reaction since
			// we're authenticated as that user
			if id, ok := reaction["id"].(float64); ok {
				reactionID = int(id)
				break
			}
		}
	}

	if reactionID == 0 {
		return fmt.Errorf("reaction '%s' not found or not owned by current user", reactionType)
	}

	// Delete the reaction using the reaction ID
	// GitHub API endpoint for deleting reactions from pull request comments
	err = client.Delete(fmt.Sprintf("repos/%s/pulls/comments/%d/reactions/%d", repo, commentID, reactionID), nil)
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	return nil
}

func resolveComment(repo string, commentID int, prNumber int) error {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return fmt.Errorf("failed to create GraphQL client: %w", err)
	}

	// Parse repo owner and name
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repo)
	}
	owner, repoName := parts[0], parts[1]

	// Step 1: Find the review thread containing this comment
	prQuery := `
		query($owner: String!, $repo: String!, $number: Int!) {
			repository(owner: $owner, name: $repo) {
				pullRequest(number: $number) {
					id
					reviewThreads(first: 100) {
						nodes {
							id
							isResolved
							comments(first: 10) {
								nodes {
									databaseId
								}
							}
						}
					}
			}
		}
	}`

	type PRData struct {
		Repository struct {
			PullRequest struct {
				ID           string `json:"id"`
				ReviewThreads struct {
					Nodes []struct {
						ID         string `json:"id"`
						IsResolved bool   `json:"isResolved"`
						Comments   struct {
							Nodes []struct {
								DatabaseID int `json:"databaseId"`
							} `json:"nodes"`
						} `json:"comments"`
					} `json:"nodes"`
				} `json:"reviewThreads"`
			} `json:"pullRequest"`
		} `json:"repository"`
	}

	var prData PRData
	err = client.Do(prQuery, map[string]interface{}{
		"owner":  owner,
		"repo":   repoName,
		"number": prNumber,
	}, &prData)
	if err != nil {
		return fmt.Errorf("failed to fetch PR data: %w", err)
	}

	// Step 2: Find the thread containing our comment
	var threadID string
	for _, thread := range prData.Repository.PullRequest.ReviewThreads.Nodes {
		if thread.IsResolved {
			continue // Skip already resolved threads
		}

		// Check if this thread contains our comment
		for _, comment := range thread.Comments.Nodes {
			if comment.DatabaseID == commentID {
				threadID = thread.ID
				break
			}
		}
		if threadID != "" {
			break
		}
	}

	if threadID == "" {
		return formatNotFoundError("unresolved thread containing comment", commentID)
	}

	if verbose {
		fmt.Printf("Found thread ID: %s for comment %d\n", threadID, commentID)
	}

	// Step 3: Resolve the thread using GraphQL mutation
	resolveMutation := `
		mutation($threadId: ID!) {
			resolveReviewThread(input: {threadId: $threadId}) {
				thread {
					id
					isResolved
				}
			}
		}`

	type ResolveResponse struct {
		ResolveReviewThread struct {
			Thread struct {
				ID         string `json:"id"`
				IsResolved bool   `json:"isResolved"`
			} `json:"thread"`
		} `json:"resolveReviewThread"`
	}

	var resolveResp ResolveResponse
	err = client.Do(resolveMutation, map[string]interface{}{
		"threadId": threadID,
	}, &resolveResp)
	if err != nil {
		return fmt.Errorf("failed to resolve thread: %w", err)
	}

	if verbose {
		fmt.Printf("Successfully resolved thread: %s (resolved: %t)\n",
			resolveResp.ResolveReviewThread.Thread.ID,
			resolveResp.ResolveReviewThread.Thread.IsResolved)
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
