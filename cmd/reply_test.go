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
			name:           "message replies not supported yet",
			args:           []string{"123456", "Great point!"},
			wantErr:        true,
			expectedErrMsg: "message replies not yet refactored",
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
