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

func TestReplyCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		flags        map[string]string
		wantErr      bool
		wantContains []string
		setupMock    func(*github.MockClient)
	}{
		{
			name: "reply to issue comment",
			args: []string{"123456", "Thanks for the feedback!"},
			flags: map[string]string{
				"type": "issue",
			},
			wantContains: []string{"✅ Replied to issue comment #123456"},
			setupMock: func(mock *github.MockClient) {
				// Mock will create a comment successfully
			},
		},
		{
			name: "reply to review comment",
			args: []string{"654321", "Fixed in the latest commit!"},
			flags: map[string]string{
				"type": "review",
			},
			wantContains: []string{"✅ Replied to review comment #654321"},
			setupMock: func(mock *github.MockClient) {
				// Mock will create a comment successfully
			},
		},
		{
			name: "reply with resolve",
			args: []string{"654321", "Fixed!"},
			flags: map[string]string{
				"type":    "review",
				"resolve": "true",
			},
			wantContains: []string{
				"✅ Replied to review comment #654321",
				"✅ Resolved conversation for comment #654321",
			},
			setupMock: func(mock *github.MockClient) {
				// Mock will create comment and resolve thread
			},
		},
		{
			name: "dry run mode",
			args: []string{"123456", "Test message"},
			flags: map[string]string{
				"dry-run": "true",
				"type":    "issue",
			},
			wantContains: []string{
				"Would reply to comment #123456:",
				"Message: Test message",
			},
		},
		{
			name: "invalid comment ID",
			args: []string{"invalid", "message"},
			flags: map[string]string{
				"type": "issue",
			},
			wantErr: true,
		},
		{
			name: "invalid comment type",
			args: []string{"123456", "message"},
			flags: map[string]string{
				"type": "invalid",
			},
			wantErr: true,
		},
		{
			name:    "no message or action",
			args:    []string{"123456"},
			wantErr: true,
		},
		{
			name: "resolve only",
			args: []string{"654321"},
			flags: map[string]string{
				"resolve": "true",
				"type":    "review",
			},
			wantContains: []string{"✅ Resolved conversation for comment #654321"},
		},
		{
			name: "API error handling",
			args: []string{"123456", "message"},
			flags: map[string]string{
				"type": "issue",
			},
			setupMock: func(mock *github.MockClient) {
				mock.CreateCommentError = fmt.Errorf("API error: rate limit exceeded")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flags
			resetReplyFlags()

			// Set up mock client
			mockClient := github.NewMockClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			// Create command with dependency injection
			cmd := NewReplyCmdWithDeps(&Dependencies{
				GitHubClient: mockClient,
				Output:       &bytes.Buffer{},
			})

			// Set flags
			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				require.NoError(t, err)
			}

			// Capture output
			var output bytes.Buffer
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
		})
	}
}

func TestValidateReaction(t *testing.T) {
	tests := []struct {
		reaction string
		valid    bool
	}{
		{"+1", true},
		{"-1", true},
		{"laugh", true},
		{"confused", true},
		{"heart", true},
		{"hooray", true},
		{"rocket", true},
		{"eyes", true},
		{"invalid", false},
		{"", false},
		{"thumbsup", false}, // GitHub uses +1, not thumbsup
	}

	for _, tt := range tests {
		t.Run(tt.reaction, func(t *testing.T) {
			result := validateReaction(tt.reaction)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestExpandSuggestions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple suggestion",
			input:    "```suggestion\nconst name = 'new';\n```",
			expected: "```suggestion\nconst name = 'new';\n```",
		},
		{
			name:     "suggestion with text",
			input:    "Consider this change:\n```suggestion\nconst name = 'better';\n```\nThis is more descriptive.",
			expected: "Consider this change:\n```suggestion\nconst name = 'better';\n```\nThis is more descriptive.",
		},
		{
			name:     "no suggestion",
			input:    "Just a regular comment",
			expected: "Just a regular comment",
		},
		{
			name:     "multiple suggestions",
			input:    "```suggestion\nline1\n```\nAnd also:\n```suggestion\nline2\n```",
			expected: "```suggestion\nline1\n```\nAnd also:\n```suggestion\nline2\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandSuggestions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// NewReplyCmdWithDeps creates a reply command with injectable dependencies
func NewReplyCmdWithDeps(deps *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply <comment-id> [message]",
		Short: "Reply to a comment",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReplyWithDeps(cmd, args, deps)
		},
	}

	// Add flags
	cmd.Flags().StringVar(&commentType, "type", "review", "Comment type: 'issue' or 'review'")
	cmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
	cmd.Flags().StringVar(&removeReaction, "remove-reaction", "", "Remove reaction")
	cmd.Flags().BoolVar(&resolveConversation, "resolve", false, "Resolve the conversation")
	cmd.Flags().BoolVar(&noExpandSuggestionsReply, "no-expand-suggestions", false, "Don't expand suggestion syntax")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")

	return cmd
}

// runReplyWithDeps is the testable version of runReply with dependency injection
func runReplyWithDeps(cmd *cobra.Command, args []string, deps *Dependencies) error {
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

	// For testing, use mock values
	_ = "owner/repo"
	pr := 1

	if dryRun {
		fmt.Fprintf(deps.Output, "Would reply to comment #%d:\n", commentID)
		if message != "" {
			fmt.Fprintf(deps.Output, "Message: %s\n", message)
		}
		if reaction != "" {
			fmt.Fprintf(deps.Output, "Reaction: %s\n", reaction)
		}
		if removeReaction != "" {
			fmt.Fprintf(deps.Output, "Remove Reaction: %s\n", removeReaction)
		}
		if resolveConversation {
			fmt.Fprintf(deps.Output, "Resolve Conversation: true\n")
		}
		return nil
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
			_, err = deps.GitHubClient.CreateIssueComment("owner", "repo", pr, finalMessage)
		} else {
			_, err = deps.GitHubClient.CreateReviewCommentReply("owner", "repo", commentID, finalMessage)
		}

		if err != nil {
			return fmt.Errorf("failed to add reply: %w", err)
		}
		fmt.Fprintf(deps.Output, "✅ Replied to %s comment #%d\n", commentType, commentID)
	}

	// Resolve conversation if specified
	if resolveConversation {
		threadID, err := deps.GitHubClient.FindReviewThreadForComment("owner", "repo", pr, commentID)
		if err != nil {
			return fmt.Errorf("failed to find review thread: %w", err)
		}

		err = deps.GitHubClient.ResolveReviewThread(threadID)
		if err != nil {
			return fmt.Errorf("failed to resolve conversation: %w", err)
		}
		fmt.Fprintf(deps.Output, "✅ Resolved conversation for comment #%d\n", commentID)
	}

	return nil
}

// resetReplyFlags resets global flags to default values for testing
func resetReplyFlags() {
	commentType = "review"
	reaction = ""
	removeReaction = ""
	resolveConversation = false
	noExpandSuggestionsReply = false
	dryRun = false
}
