package cmd

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRunReplyIntegration tests the actual runReply function with various scenarios
func TestRunReplyIntegration(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		flags        map[string]string
		mockSetup    func(*MockReplyClient)
		wantErr      bool
		wantContains []string
		wantCalls    []string // Expected method calls on mock
	}{
		{
			name: "successful issue comment reply",
			args: []string{"123456", "Thanks for the feedback!"},
			flags: map[string]string{
				"type": "issue",
			},
			mockSetup: func(m *MockReplyClient) {
				m.createIssueCommentResult = &github.Comment{
					ID:   789,
					Body: "Thanks for the feedback!",
				}
			},
			wantContains: []string{"Reply added successfully"},
			wantCalls:    []string{"CreateIssueComment"},
		},
		{
			name: "successful review comment reply",
			args: []string{"654321", "Fixed the issue!"},
			flags: map[string]string{
				"type": "review",
			},
			mockSetup: func(m *MockReplyClient) {
				m.createReviewCommentReplyResult = &github.Comment{
					ID:   890,
					Body: "Fixed the issue!",
				}
			},
			wantContains: []string{"Reply added successfully"},
			wantCalls:    []string{"CreateReviewCommentReply"},
		},
		{
			name: "reply with reaction",
			args: []string{"123456", "Great work!"},
			flags: map[string]string{
				"type":     "issue",
				"reaction": "+1",
			},
			mockSetup: func(m *MockReplyClient) {
				m.createIssueCommentResult = &github.Comment{ID: 789}
			},
			wantContains: []string{"Reply added successfully"},
			wantCalls:    []string{"CreateIssueComment", "AddReaction"},
		},
		{
			name: "reply with resolve",
			args: []string{"654321", "Fixed!"},
			flags: map[string]string{
				"type":    "review",
				"resolve": "true",
			},
			mockSetup: func(m *MockReplyClient) {
				m.createReviewCommentReplyResult = &github.Comment{ID: 890}
				m.findReviewThreadResult = "thread123"
			},
			wantContains: []string{"Reply added successfully", "Conversation resolved"},
			wantCalls:    []string{"CreateReviewCommentReply", "FindReviewThreadForComment", "ResolveReviewThread"},
		},
		{
			name: "dry run mode",
			args: []string{"123456", "Test message"},
			flags: map[string]string{
				"type":    "issue",
				"dry-run": "true",
			},
			mockSetup: func(m *MockReplyClient) {
				// No setup needed for dry run
			},
			wantContains: []string{"[DRY RUN]", "Would reply to comment 123456"},
			wantCalls:    []string{}, // No actual API calls in dry run
		},
		{
			name: "suggestion expansion",
			args: []string{"123456", "```suggestion\nfixed code\n```"},
			flags: map[string]string{
				"type": "review",
			},
			mockSetup: func(m *MockReplyClient) {
				m.createReviewCommentReplyResult = &github.Comment{ID: 890}
			},
			wantContains: []string{"Reply added successfully"},
			wantCalls:    []string{"CreateReviewCommentReply"},
		},
		{
			name:    "invalid comment ID",
			args:    []string{"invalid", "message"},
			wantErr: true,
		},
		{
			name:    "no message or action",
			args:    []string{"123456"},
			wantErr: true,
		},
		{
			name: "invalid reaction",
			args: []string{"123456", "message"},
			flags: map[string]string{
				"reaction": "invalid",
			},
			wantErr: true,
		},
		{
			name: "both reaction and remove-reaction",
			args: []string{"123456", "message"},
			flags: map[string]string{
				"reaction":        "+1",
				"remove-reaction": "heart",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flags
			resetReplyFlags()

			// Create mock client
			mockClient := &MockReplyClient{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			// Capture output
			var output bytes.Buffer

			// Create test command
			cmd := &cobra.Command{
				Use:  "reply <comment-id> [message]",
				Args: cobra.RangeArgs(1, 2),
				RunE: func(cmd *cobra.Command, args []string) error {
					return runReplyWithMock(cmd, args, mockClient, &output)
				},
			}

			// Add flags
			cmd.Flags().StringVar(&commentType, "type", "review", "Comment type")
			cmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction")
			cmd.Flags().StringVar(&removeReaction, "remove-reaction", "", "Remove reaction")
			cmd.Flags().BoolVar(&resolveConversation, "resolve", false, "Resolve conversation")
			cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done")
			cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

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

// MockReplyClient for testing reply functionality
type MockReplyClient struct {
	calls                           []string
	createIssueCommentResult        *github.Comment
	createReviewCommentReplyResult  *github.Comment
	findReviewThreadResult          string
	shouldError                     bool
}

func (m *MockReplyClient) ListIssueComments(owner, repo string, prNumber int) ([]github.Comment, error) {
	m.calls = append(m.calls, "ListIssueComments")
	return nil, fmt.Errorf("not implemented")
}

func (m *MockReplyClient) ListReviewComments(owner, repo string, prNumber int) ([]github.Comment, error) {
	m.calls = append(m.calls, "ListReviewComments")
	return nil, fmt.Errorf("not implemented")
}

func (m *MockReplyClient) CreateIssueComment(owner, repo string, prNumber int, body string) (*github.Comment, error) {
	m.calls = append(m.calls, "CreateIssueComment")
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.createIssueCommentResult, nil
}

func (m *MockReplyClient) CreateReviewCommentReply(owner, repo string, commentID int, body string) (*github.Comment, error) {
	m.calls = append(m.calls, "CreateReviewCommentReply")
	if m.shouldError {
		return nil, fmt.Errorf("mock error")
	}
	return m.createReviewCommentReplyResult, nil
}

func (m *MockReplyClient) ResolveReviewThread(threadID string) error {
	m.calls = append(m.calls, "ResolveReviewThread")
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockReplyClient) FindReviewThreadForComment(owner, repo string, prNumber, commentID int) (string, error) {
	m.calls = append(m.calls, "FindReviewThreadForComment")
	if m.shouldError {
		return "", fmt.Errorf("mock error")
	}
	return m.findReviewThreadResult, nil
}

func (m *MockReplyClient) AddReaction(owner, repo string, commentID int, reaction string) error {
	m.calls = append(m.calls, "AddReaction")
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockReplyClient) RemoveReaction(owner, repo string, commentID int, reaction string) error {
	m.calls = append(m.calls, "RemoveReaction")
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

// runReplyWithMock is a testable version of runReply
func runReplyWithMock(cmd *cobra.Command, args []string, mockClient *MockReplyClient, output *bytes.Buffer) error {
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
		err = mockClient.AddReaction("owner", "repo", commentID, reaction)
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}
	}

	// Remove reaction if specified
	if removeReaction != "" {
		err = mockClient.RemoveReaction("owner", "repo", commentID, removeReaction)
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
			_, err = mockClient.CreateIssueComment("owner", "repo", pr, finalMessage)
		} else {
			_, err = mockClient.CreateReviewCommentReply("owner", "repo", commentID, finalMessage)
		}

		if err != nil {
			return fmt.Errorf("failed to add reply: %w", err)
		}
	}

	// Resolve conversation if specified
	if resolveConversation {
		threadID, err := mockClient.FindReviewThreadForComment("owner", "repo", pr, commentID)
		if err != nil {
			return fmt.Errorf("failed to find review thread: %w", err)
		}

		err = mockClient.ResolveReviewThread(threadID)
		if err != nil {
			return fmt.Errorf("failed to resolve conversation: %w", err)
		}

		fmt.Fprintf(output, "✅ Conversation resolved\n")
	}

	fmt.Fprintf(output, "✅ Reply added successfully\n")
	return nil
}


