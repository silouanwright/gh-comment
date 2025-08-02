package github

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestClient(t *testing.T) {
	// Save original environment
	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	tests := []struct {
		name        string
		mockURL     string
		expectError bool
	}{
		{
			name:        "with mock server URL",
			mockURL:     "http://localhost:8080",
			expectError: false,
		},
		{
			name:        "with https mock server URL",
			mockURL:     "https://mock.example.com",
			expectError: false,
		},
		{
			name:        "without mock server URL",
			mockURL:     "",
			expectError: false, // Should succeed and default to api.github.com
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockURL != "" {
				os.Setenv("MOCK_SERVER_URL", tt.mockURL)
			} else {
				os.Unsetenv("MOCK_SERVER_URL")
			}

			client, err := NewTestClient()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				if tt.mockURL != "" {
					assert.Equal(t, tt.mockURL, client.baseURL)
				} else {
					assert.Equal(t, "https://api.github.com", client.baseURL)
				}
			}
		})
	}
}

func TestTestClientWithMockServer(t *testing.T) {
	// Set up mock server URL for testing
	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	os.Setenv("MOCK_SERVER_URL", "http://localhost:8080")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	client, err := NewTestClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("ListIssueComments", func(t *testing.T) {
		// This will try to make an HTTP request to the mock server
		// In a real test environment, this would fail since no server is running
		// But we can verify the method is callable
		comments, err := client.ListIssueComments("owner", "repo", 123)
		
		// We expect an error since no real mock server is running
		assert.Error(t, err)
		assert.Nil(t, comments)
		
		// Verify error is related to connection (not validation)
		assert.Contains(t, err.Error(), "connection refused")
	})

	t.Run("ListReviewComments", func(t *testing.T) {
		comments, err := client.ListReviewComments("owner", "repo", 123)
		
		// We expect an error since no real mock server is running
		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Contains(t, err.Error(), "connection refused")
	})

	t.Run("CreateIssueComment", func(t *testing.T) {
		comment, err := client.CreateIssueComment("owner", "repo", 123, "test comment")
		
		// We expect an error since no real mock server is running
		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Contains(t, err.Error(), "connection refused")
	})

	t.Run("AddReviewComment", func(t *testing.T) {
		reviewComment := ReviewCommentInput{
			Body: "Test review comment",
			Path: "test.go",
			Line: 42,
		}
		
		err := client.AddReviewComment("owner", "repo", 123, reviewComment)
		
		// We expect an error since no real mock server is running
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection refused")
	})

	t.Run("CreateReview", func(t *testing.T) {
		review := ReviewInput{
			Body:  "Test review",
			Event: "APPROVE",
		}
		
		err := client.CreateReview("owner", "repo", 123, review)
		
		// We expect an error since no real mock server is running
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection refused")
	})

	t.Run("GetPRDetails", func(t *testing.T) {
		details, err := client.GetPRDetails("owner", "repo", 123)
		
		// We expect an error since no real mock server is running
		assert.Error(t, err)
		assert.Nil(t, details)
		assert.Contains(t, err.Error(), "connection refused")
	})

	// Test methods that are not implemented (should return appropriate errors)
	t.Run("CreateReviewCommentReply", func(t *testing.T) {
		comment, err := client.CreateReviewCommentReply("owner", "repo", 123, "reply")
		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("FindReviewThreadForComment", func(t *testing.T) {
		threadID, err := client.FindReviewThreadForComment("owner", "repo", 123, 456)
		assert.Error(t, err)
		assert.Empty(t, threadID)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("ResolveReviewThread", func(t *testing.T) {
		err := client.ResolveReviewThread("thread123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("AddReaction", func(t *testing.T) {
		err := client.AddReaction("owner", "repo", 123, "+1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("RemoveReaction", func(t *testing.T) {
		err := client.RemoveReaction("owner", "repo", 123, "+1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("EditComment", func(t *testing.T) {
		err := client.EditComment("owner", "repo", 123, "edited")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("FetchPRDiff", func(t *testing.T) {
		diff, err := client.FetchPRDiff("owner", "repo", 123)
		assert.Error(t, err)
		assert.Nil(t, diff)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("FindPendingReview", func(t *testing.T) {
		reviewID, err := client.FindPendingReview("owner", "repo", 123)
		assert.Error(t, err)
		assert.Zero(t, reviewID)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("SubmitReview", func(t *testing.T) {
		err := client.SubmitReview("owner", "repo", 123, 456, "body", "APPROVE")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})
}

func TestTestClientDoRequest(t *testing.T) {
	// Set up mock server URL for testing
	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	os.Setenv("MOCK_SERVER_URL", "http://localhost:8080")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	client, err := NewTestClient()
	require.NoError(t, err)

	// Test that doRequest constructs URLs correctly
	// We can't actually make requests without a real server, but we can verify
	// the method is callable and handles errors appropriately
	
	resp, err := client.doRequest("GET", "/test/endpoint", nil)
	
	// Should fail with connection error since no server is running
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "connection refused")
}

func TestTestClientErrorHandling(t *testing.T) {
	// Test with invalid URL to trigger URL parsing errors
	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	os.Setenv("MOCK_SERVER_URL", "invalid-url-format")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	client, err := NewTestClient()
	require.NoError(t, err) // NewTestClient doesn't validate URL format

	// Test doRequest with invalid base URL
	resp, err := client.doRequest("GET", "/test", nil)
	
	// Should fail with URL error
	assert.Error(t, err)
	assert.Nil(t, resp)
	// Error message will vary depending on what exactly fails first
	assert.NotEmpty(t, err.Error())
}