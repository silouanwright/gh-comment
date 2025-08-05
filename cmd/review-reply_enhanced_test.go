package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleReviewReplyError(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		commentID    int
		message      string
		expectedMsg  string
	}{
		{
			name:        "404 GitHub API limitation",
			inputError:  fmt.Errorf("404 Not Found: pulls/comments/123/replies"),
			commentID:   123,
			message:     "Fixed issue",
			expectedMsg: "review comment threading not supported",
		},
		{
			name:        "Issue comment mismatch",
			inputError:  fmt.Errorf("404 Not Found: issues/comments/456"),
			commentID:   456,
			message:     "Response",
			expectedMsg: "is an issue comment",
		},
		{
			name:        "Generic error",
			inputError:  fmt.Errorf("network timeout"),
			commentID:   789,
			message:     "Test",
			expectedMsg: "failed to create reply",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleReviewReplyError(tt.inputError, tt.commentID, tt.message, "owner", "repo", 123)
			
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedMsg)
		})
	}
}

func TestReviewReplyIntelligentErrorHandling(t *testing.T) {
	// Test the specific error patterns that the enhanced error handler recognizes
	
	t.Run("GitHub API 404 threading limitation", func(t *testing.T) {
		err := fmt.Errorf("404 Not Found: pulls/comments/12345/replies")
		result := handleReviewReplyError(err, 12345, "Test message", "owner", "repo", 123)
		
		assert.Error(t, result)
		assert.Contains(t, result.Error(), "review comment threading not supported")
	})
	
	t.Run("Issue comment type mismatch", func(t *testing.T) {
		err := fmt.Errorf("404 Not Found: issues/comments/67890")
		result := handleReviewReplyError(err, 67890, "Test message", "owner", "repo", 123)
		
		assert.Error(t, result)
		assert.Contains(t, result.Error(), "is an issue comment")
	})
	
	t.Run("Unrecognized error passes through", func(t *testing.T) {
		err := fmt.Errorf("some other API error")
		result := handleReviewReplyError(err, 111, "Test message", "owner", "repo", 123)
		
		assert.Error(t, result)
		assert.Contains(t, result.Error(), "failed to create reply")
	})
}