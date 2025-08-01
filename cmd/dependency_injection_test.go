package cmd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDependencyInjection tests commands with proper dependency injection
func TestDependencyInjection(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*FullMockClient)
		command      string
		args         []string
		flags        map[string]string
		wantErr      bool
		wantContains []string
		wantCalls    []string
	}{
		{
			name: "list command with dependency injection",
			setupMock: func(m *FullMockClient) {
				m.issueComments = []github.Comment{
					{
						ID:        123456,
						Body:      "Great work on this feature!",
						Type:      "issue",
						User:      github.User{Login: "reviewer1"},
						CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					},
				}
				m.reviewComments = []github.Comment{
					{
						ID:        654321,
						Body:      "Consider adding error handling here.",
						Type:      "review",
						User:      github.User{Login: "reviewer2"},
						CreatedAt: time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
						Path:      "main.go",
						Line:      42,
					},
				}
			},
			command: "list",
			args:    []string{"123"},
			wantContains: []string{
				"Comments on PR #123",
				"General PR Comments",
				"Review Comments",
				"Great work on this feature!",
				"Consider adding error handling here.",
				"reviewer1",
				"reviewer2",
				"main.go:L42",
			},
			wantCalls: []string{"ListIssueComments", "ListReviewComments"},
		},
		{
			name: "reply command with dependency injection",
			setupMock: func(m *FullMockClient) {
				m.createIssueCommentResult = &github.Comment{
					ID:   789,
					Body: "Thanks for the review!",
				}
			},
			command: "reply",
			args:    []string{"123456", "Thanks for the review!"},
			flags: map[string]string{
				"type": "issue",
			},
			wantContains: []string{"Reply added successfully"},
			wantCalls:    []string{"CreateIssueComment"},
		},
		{
			name: "reply with reaction and resolve",
			setupMock: func(m *FullMockClient) {
				m.createReviewCommentReplyResult = &github.Comment{
					ID:   890,
					Body: "Fixed the issue!",
				}
				m.findReviewThreadResult = "thread123"
			},
			command: "reply",
			args:    []string{"654321", "Fixed the issue!"},
			flags: map[string]string{
				"type":     "review",
				"reaction": "+1",
				"resolve":  "true",
			},
			wantContains: []string{
				"Reply added successfully",
				"Conversation resolved",
			},
			wantCalls: []string{
				"AddReaction",
				"CreateReviewCommentReply",
				"FindReviewThreadForComment",
				"ResolveReviewThread",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &FullMockClient{}
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			// Reset global flags
			resetAllFlags()

			// Capture output
			var output bytes.Buffer

			// Create command with dependency injection
			var cmd *cobra.Command
			switch tt.command {
			case "list":
				cmd = createListCommandWithDI(mockClient, &output)
			case "reply":
				cmd = createReplyCommandWithDI(mockClient, &output)
			default:
				t.Fatalf("Unknown command: %s", tt.command)
			}

			// Set flags
			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				require.NoError(t, err)
			}

			// Set output
			cmd.SetOut(&output)
			cmd.SetErr(&output)

			// Run command
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			// Check error expectation
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check output
			outputStr := output.String()
			for _, want := range tt.wantContains {
				assert.Contains(t, outputStr, want, "Output should contain: %s", want)
			}

			// Check expected method calls
			for _, expectedCall := range tt.wantCalls {
				assert.Contains(t, mockClient.calls, expectedCall, "Expected method call: %s", expectedCall)
			}
		})
	}
}

// FullMockClient implements all GitHub API methods for comprehensive testing
type FullMockClient struct {
	calls                          []string
	issueComments                  []github.Comment
	reviewComments                 []github.Comment
	createIssueCommentResult       *github.Comment
	createReviewCommentReplyResult *github.Comment
	findReviewThreadResult         string
	shouldError                    bool
}

func (m *FullMockClient) ListIssueComments(owner, repo string, prNumber int) ([]github.Comment, error) {
	m.calls = append(m.calls, "ListIssueComments")
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.issueComments, nil
}

func (m *FullMockClient) ListReviewComments(owner, repo string, prNumber int) ([]github.Comment, error) {
	m.calls = append(m.calls, "ListReviewComments")
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.reviewComments, nil
}

func (m *FullMockClient) CreateIssueComment(owner, repo string, prNumber int, body string) (*github.Comment, error) {
	m.calls = append(m.calls, "CreateIssueComment")
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.createIssueCommentResult, nil
}

func (m *FullMockClient) CreateReviewCommentReply(owner, repo string, commentID int, body string) (*github.Comment, error) {
	m.calls = append(m.calls, "CreateReviewCommentReply")
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.createReviewCommentReplyResult, nil
}

func (m *FullMockClient) ResolveReviewThread(threadID string) error {
	m.calls = append(m.calls, "ResolveReviewThread")
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *FullMockClient) FindReviewThreadForComment(owner, repo string, prNumber, commentID int) (string, error) {
	m.calls = append(m.calls, "FindReviewThreadForComment")
	if m.shouldError {
		return "", fmt.Errorf("mock error")
	}
	return m.findReviewThreadResult, nil
}

func (m *FullMockClient) AddReaction(owner, repo string, commentID int, reaction string) error {
	m.calls = append(m.calls, "AddReaction")
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *FullMockClient) RemoveReaction(owner, repo string, commentID int, reaction string) error {
	m.calls = append(m.calls, "RemoveReaction")
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

// createListCommandWithDI creates a list command with dependency injection
func createListCommandWithDI(client *FullMockClient, output *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [pr]",
		Short: "List comments on a PR",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListWithDI(cmd, args, client, output)
		},
	}

	cmd.Flags().BoolVar(&showResolved, "resolved", false, "Include resolved comments")
	cmd.Flags().BoolVar(&onlyUnresolved, "unresolved", false, "Show only unresolved comments")
	cmd.Flags().StringVar(&author, "author", "", "Filter comments by author")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output without URLs and IDs")
	cmd.Flags().BoolVar(&hideAuthors, "hide-authors", false, "Hide author names for privacy")

	return cmd
}

// createReplyCommandWithDI creates a reply command with dependency injection
func createReplyCommandWithDI(client *FullMockClient, output *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply <comment-id> [message]",
		Short: "Reply to a comment",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReplyWithDI(cmd, args, client, output)
		},
	}

	cmd.Flags().StringVar(&commentType, "type", "review", "Comment type (issue or review)")
	cmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction")
	cmd.Flags().StringVar(&removeReaction, "remove-reaction", "", "Remove reaction")
	cmd.Flags().BoolVar(&resolveConversation, "resolve", false, "Resolve conversation")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done")

	return cmd
}

// runListWithDI is a dependency-injected version of runList
func runListWithDI(cmd *cobra.Command, args []string, client *FullMockClient, output *bytes.Buffer) error {
	var pr int
	var err error

	// Parse PR argument
	if len(args) == 1 {
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return formatValidationError("PR number", args[0], "must be a valid integer")
		}
	} else {
		pr = 1 // Default for testing
	}

	// Mock repository
	repository := "owner/repo"

	if verbose {
		fmt.Fprintf(output, "Repository: %s\n", repository)
		fmt.Fprintf(output, "PR: %d\n", pr)
		fmt.Fprintf(output, "\n")
	}

	// Fetch comments using injected client
	issueComments, err := client.ListIssueComments("owner", "repo", pr)
	if err != nil {
		return fmt.Errorf("failed to fetch issue comments: %w", err)
	}

	reviewComments, err := client.ListReviewComments("owner", "repo", pr)
	if err != nil {
		return fmt.Errorf("failed to fetch review comments: %w", err)
	}

	// Convert to internal Comment format and combine
	var allComments []Comment
	for _, comment := range issueComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Type:      comment.Type,
		})
	}
	for _, comment := range reviewComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Path:      comment.Path,
			Line:      comment.Line,
			Type:      comment.Type,
		})
	}

	// Filter comments
	filteredComments := filterComments(allComments)

	// Display comments
	if len(filteredComments) == 0 {
		fmt.Fprintf(output, "No comments found on PR #%d\n", pr)
		return nil
	}

	// Display header
	fmt.Fprintf(output, "📝 Comments on PR #%d (%d total)\n\n", pr, len(filteredComments))

	// Group comments by type
	var issueCommentsFiltered, reviewCommentsFiltered []Comment
	for _, comment := range filteredComments {
		if comment.Type == "issue" {
			issueCommentsFiltered = append(issueCommentsFiltered, comment)
		} else {
			reviewCommentsFiltered = append(reviewCommentsFiltered, comment)
		}
	}

	// Display general PR comments
	if len(issueCommentsFiltered) > 0 {
		fmt.Fprintf(output, "💬 General PR Comments (%d)\n", len(issueCommentsFiltered))
		fmt.Fprintf(output, "%s\n", strings.Repeat("─", 50))
		for i, comment := range issueCommentsFiltered {
			displayCommentToBuffer(comment, i+1, output)
		}
		fmt.Fprintf(output, "\n")
	}

	// Display review comments
	if len(reviewCommentsFiltered) > 0 {
		fmt.Fprintf(output, "📋 Review Comments (%d)\n", len(reviewCommentsFiltered))
		fmt.Fprintf(output, "%s\n", strings.Repeat("─", 50))
		for i, comment := range reviewCommentsFiltered {
			displayCommentToBuffer(comment, i+1, output)
		}
	}

	return nil
}

// runReplyWithDI is a dependency-injected version of runReply
func runReplyWithDI(cmd *cobra.Command, args []string, client *FullMockClient, output *bytes.Buffer) error {
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

	// Mock repository and PR
	repository := "owner/repo"
	pr := 1

	if verbose {
		fmt.Fprintf(output, "Repository: %s\n", repository)
		fmt.Fprintf(output, "PR: %d\n", pr)
		fmt.Fprintf(output, "Comment ID: %d\n", commentID)
		fmt.Fprintf(output, "Comment Type: %s\n", commentType)
		if message != "" {
			fmt.Fprintf(output, "Message: %s\n", message)
		}
		if reaction != "" {
			fmt.Fprintf(output, "Reaction: %s\n", reaction)
		}
		if removeReaction != "" {
			fmt.Fprintf(output, "Remove Reaction: %s\n", removeReaction)
		}
		fmt.Fprintf(output, "Resolve: %v\n", resolveConversation)
		fmt.Fprintf(output, "Dry Run: %v\n", dryRun)
		fmt.Fprintf(output, "\n")
	}

	if dryRun {
		fmt.Fprintf(output, "🔍 [DRY RUN] Would reply to comment %d\n", commentID)
		if message != "" {
			fmt.Fprintf(output, "   Message: %s\n", message)
		}
		if reaction != "" {
			fmt.Fprintf(output, "   Add reaction: %s\n", reaction)
		}
		if removeReaction != "" {
			fmt.Fprintf(output, "   Remove reaction: %s\n", removeReaction)
		}
		if resolveConversation {
			fmt.Fprintf(output, "   Resolve conversation: yes\n")
		}
		return nil
	}

	// Add reaction if specified
	if reaction != "" {
		err = client.AddReaction("owner", "repo", commentID, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}
	}

	// Remove reaction if specified
	if removeReaction != "" {
		err = client.RemoveReaction("owner", "repo", commentID, removeReaction)
		if err != nil {
			return fmt.Errorf("failed to remove reaction: %w", err)
		}
	}

	// Add reply message if specified
	if message != "" {
		// Expand suggestion syntax
		finalMessage := expandSuggestions(message)

		// Use appropriate reply method based on comment type
		if commentType == "issue" {
			_, err = client.CreateIssueComment("owner", "repo", pr, finalMessage)
		} else {
			_, err = client.CreateReviewCommentReply("owner", "repo", commentID, finalMessage)
		}

		if err != nil {
			return fmt.Errorf("failed to add reply: %w", err)
		}
	}

	// Resolve conversation if specified
	if resolveConversation {
		threadID, err := client.FindReviewThreadForComment("owner", "repo", pr, commentID)
		if err != nil {
			return fmt.Errorf("failed to find review thread: %w", err)
		}

		err = client.ResolveReviewThread(threadID)
		if err != nil {
			return fmt.Errorf("failed to resolve conversation: %w", err)
		}

		fmt.Fprintf(output, "✅ Conversation resolved\n")
	}

	fmt.Fprintf(output, "✅ Reply added successfully\n")
	return nil
}

// resetAllFlags resets all command flags to their default values
func resetAllFlags() {
	// List flags
	showResolved = false
	onlyUnresolved = false
	author = ""
	quiet = false
	hideAuthors = false

	// Reply flags
	commentType = "review"
	reaction = ""
	removeReaction = ""
	resolveConversation = false
	dryRun = false

	// Global flags
	verbose = false
}
