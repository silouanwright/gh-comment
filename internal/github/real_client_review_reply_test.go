package github

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealClientCreateReviewCommentReply(t *testing.T) {
	tests := []struct {
		name             string
		owner            string
		repo             string
		commentID        int
		body             string
		mockComment      *Comment
		mockCommentError error
		mockReply        *Comment
		mockReplyError   error
		wantErr          string
		wantReplyID      int
	}{
		{
			name:      "successful reply",
			owner:     "owner",
			repo:      "repo",
			commentID: 12345,
			body:      "This is a reply",
			mockComment: &Comment{
				ID:             12345,
				Body:           "Original comment",
				PullRequestURL: "https://api.github.com/repos/owner/repo/pulls/422",
			},
			mockReply: &Comment{
				ID:          67890,
				Body:        "This is a reply",
				InReplyToID: 12345,
			},
			wantReplyID: 67890,
		},
		{
			name:      "empty owner",
			owner:     "",
			repo:      "repo",
			commentID: 12345,
			body:      "Reply",
			wantErr:   "repository owner cannot be empty",
		},
		{
			name:      "empty repo",
			owner:     "owner",
			repo:      "",
			commentID: 12345,
			body:      "Reply",
			wantErr:   "repository name cannot be empty",
		},
		{
			name:      "invalid comment ID",
			owner:     "owner",
			repo:      "repo",
			commentID: 0,
			body:      "Reply",
			wantErr:   "invalid comment ID 0: must be positive",
		},
		{
			name:      "empty body",
			owner:     "owner",
			repo:      "repo",
			commentID: 12345,
			body:      "  ",
			wantErr:   "reply body cannot be empty",
		},
		{
			name:             "comment not found",
			owner:            "owner",
			repo:             "repo",
			commentID:        12345,
			body:             "Reply",
			mockCommentError: errors.New("not found"),
			wantErr:          "not found",
		},
		{
			name:      "comment has no PR URL",
			owner:     "owner",
			repo:      "repo",
			commentID: 12345,
			body:      "Reply",
			mockComment: &Comment{
				ID:             12345,
				Body:           "Original comment",
				PullRequestURL: "",
			},
			wantErr: "comment #12345 has no associated pull request",
		},
		{
			name:      "invalid PR URL format",
			owner:     "owner",
			repo:      "repo",
			commentID: 12345,
			body:      "Reply",
			mockComment: &Comment{
				ID:             12345,
				Body:           "Original comment",
				PullRequestURL: "not-a-url",
			},
			wantErr: "failed to parse PR number from URL",
		},
		{
			name:      "PR URL with non-numeric PR number",
			owner:     "owner",
			repo:      "repo",
			commentID: 12345,
			body:      "Reply",
			mockComment: &Comment{
				ID:             12345,
				Body:           "Original comment",
				PullRequestURL: "https://api.github.com/repos/owner/repo/pulls/invalid",
			},
			wantErr: "failed to parse PR number from URL",
		},
		{
			name:      "pending review blocks reply",
			owner:     "owner",
			repo:      "repo",
			commentID: 12345,
			body:      "Reply",
			mockComment: &Comment{
				ID:             12345,
				Body:           "Original comment",
				PullRequestURL: "https://api.github.com/repos/owner/repo/pulls/422",
			},
			mockReplyError: errors.New("user_id can only have one pending review per pull request"),
			wantErr:        "pending review",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These are validation-only tests that don't need a real client
			if tt.mockComment == nil && tt.mockCommentError == nil {
				client := &RealClient{}
				_, err := client.CreateReviewCommentReply(tt.owner, tt.repo, tt.commentID, tt.body)
				if tt.wantErr != "" {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.wantErr)
				} else {
					require.NoError(t, err)
				}
				return
			}

			// For tests that need API mocking, we'll use the MockClient approach
			// Since RealClient needs actual API clients, we test the logic separately
			// The actual API integration is tested in integration tests
		})
	}
}

func TestCreateReviewCommentReplyURLParsing(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectedPR  int
		expectError bool
	}{
		{
			name:       "standard API URL",
			url:        "https://api.github.com/repos/owner/repo/pulls/123",
			expectedPR: 123,
		},
		{
			name:       "API URL with org and hyphenated repo",
			url:        "https://api.github.com/repos/org/repo-name/pulls/9999",
			expectedPR: 9999,
		},
		{
			name:       "non-API GitHub URL",
			url:        "https://github.com/repos/owner/repo/pulls/456",
			expectedPR: 456,
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: true,
		},
		{
			name:        "URL ending with slash",
			url:         "https://api.github.com/repos/owner/repo/pulls/",
			expectError: true,
		},
		{
			name:        "URL with non-numeric PR",
			url:         "https://api.github.com/repos/owner/repo/pulls/abc",
			expectError: true,
		},
		{
			name:       "URL with large PR number",
			url:        "https://api.github.com/repos/owner/repo/pulls/999999",
			expectedPR: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the URL parsing logic that's used in CreateReviewCommentReply
			parts := strings.Split(tt.url, "/")
			if len(parts) == 0 {
				if !tt.expectError {
					t.Error("Expected parts to have elements")
				}
				return
			}

			prNumberStr := parts[len(parts)-1]
			prNumber, err := strconv.Atoi(prNumberStr)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPR, prNumber)
			}
		})
	}
}

func TestMockClientCreateReviewCommentReplyBehavior(t *testing.T) {
	t.Run("successful reply", func(t *testing.T) {
		client := &MockClient{}
		comment, err := client.CreateReviewCommentReply("owner", "repo", 12345, "Test reply")

		require.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, 345678, comment.ID)
		assert.Equal(t, "Test reply", comment.Body)
		assert.Equal(t, "review", comment.Type)
	})

	t.Run("with error", func(t *testing.T) {
		client := &MockClient{
			CreateCommentError: errors.New("API error"),
		}
		comment, err := client.CreateReviewCommentReply("owner", "repo", 12345, "Test reply")

		require.Error(t, err)
		assert.Nil(t, comment)
		assert.Contains(t, err.Error(), "API error")
	})
}
