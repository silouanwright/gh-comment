package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunReplyWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := replyClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		replyClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	replyClient = mockClient
	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name           string
		args           []string
		setupReaction  string
		setupRemove    string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:          "add reaction to comment",
			args:          []string{"123456"},
			setupReaction: "+1",
			wantErr:       false,
		},
		{
			name:        "remove reaction from comment",
			args:        []string{"123456"},
			setupRemove: "heart",
			wantErr:     false,
		},
		{
			name:           "invalid comment ID",
			args:           []string{"invalid"},
			setupReaction:  "+1",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:    "message reply to review comment",
			args:    []string{"123456", "Great point!"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			reaction = tt.setupReaction
			removeReaction = tt.setupRemove
			resolveConversation = false
			commentType = "review"

			err := runReply(nil, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// INTEGRATION TEST: Document GitHub API limitation for review comment threading
func TestReplyToReviewCommentKnownLimitations(t *testing.T) {
	// Save original client
	originalClient := replyClient
	defer func() { replyClient = originalClient }()

	// Set up mock client (mock client allows operations that real GitHub would reject)
	mockClient := github.NewMockClient()
	replyClient = mockClient

	// Save original globals
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		repo = originalRepo
		prNumber = originalPR
	}()

	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name        string
		commentID   string
		message     string
		commentType string
		description string
	}{
		{
			name:        "review comment threading limitation",
			commentID:   "123456",
			message:     "This would fail in real GitHub",
			commentType: "review",
			description: "GitHub API doesn't support direct threading for review comments",
		},
		{
			name:        "issue comment should work",
			commentID:   "123456",
			message:     "This should work fine",
			commentType: "issue",
			description: "Issue comment threading is supported by GitHub",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset comment type
			commentType = tt.commentType

			// This test documents the limitation - real implementation would handle the error
			err := runReply(nil, []string{tt.commentID, tt.message})

			// For now, mock client allows it, but real GitHub would reject review comment threading
			if tt.commentType == "review" {
				// In a real scenario, this would fail with HTTP 422
				// We document this limitation in the integration guide
				t.Logf("LIMITATION: %s", tt.description)
			}

			// Mock client succeeds, but real API has limitations
			assert.NoError(t, err, "Mock client allows all operations, real GitHub has limitations")
		})
	}
}
