package cmd

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockGitHubServer(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	assert.NotNil(t, server)
	assert.NotEmpty(t, server.URL())
	assert.NotNil(t, server.users)
	assert.Equal(t, 4, len(server.users))
	
	// Verify users are properly set up
	testUser, exists := server.users["test-user"]
	assert.True(t, exists)
	assert.Equal(t, "test-user", testUser.Login)
	assert.Equal(t, 1, testUser.ID)
}

func TestMockGitHubServer_AddComment(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	comment := MockComment{
		Body: "Test comment",
		User: server.users["test-user"],
	}

	server.AddComment("owner/repo", 123, comment)

	comments := server.GetComments("owner/repo", 123)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Test comment", comments[0].Body)
	assert.Equal(t, 1000, comments[0].ID) // Auto-assigned ID
	assert.False(t, comments[0].CreatedAt.IsZero())
	assert.False(t, comments[0].UpdatedAt.IsZero())
	assert.NotEmpty(t, comments[0].HTMLURL)
}

func TestMockGitHubServer_SetupTestScenario(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	tests := []struct {
		name     string
		scenario string
		repo     string
		pr       int
		expected int
	}{
		{
			name:     "basic scenario",
			scenario: "basic",
			repo:     "test-owner/test-repo",
			pr:       123,
			expected: 2,
		},
		{
			name:     "security review scenario",
			scenario: "security-review",
			repo:     "test-owner/test-repo",
			pr:       456,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.SetupTestScenario(tt.scenario)
			comments := server.GetComments(tt.repo, tt.pr)
			assert.Len(t, comments, tt.expected)
		})
	}
}

func TestMockGitHubServer_HandleGetPRDetails(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	resp, err := http.Get(server.URL() + "/repos/owner/repo/pulls/123")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var details MockPRDetails
	err = json.NewDecoder(resp.Body).Decode(&details)
	require.NoError(t, err)

	assert.Equal(t, 123, details.Number)
	assert.Equal(t, "abc123def456", details.Head.SHA)
}

func TestMockGitHubServer_HandleListComments(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	// Add test comments
	server.AddComment("owner/repo", 123, MockComment{
		Body: "General comment",
		User: server.users["test-user"],
	})
	server.AddComment("owner/repo", 123, MockComment{
		Body: "Line comment",
		User: server.users["reviewer"],
		Path: "src/main.go",
		Line: 42,
	})

	tests := []struct {
		name         string
		url          string
		expectedType string
		expectedLen  int
	}{
		{
			name:         "review comments",
			url:          "/repos/owner/repo/pulls/123/comments",
			expectedType: "review",
			expectedLen:  1, // Only line comment
		},
		{
			name:         "issue comments",
			url:          "/repos/owner/repo/issues/123/comments",
			expectedType: "issue",
			expectedLen:  1, // Only general comment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL() + tt.url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var comments []MockComment
			err = json.NewDecoder(resp.Body).Decode(&comments)
			require.NoError(t, err)

			assert.Len(t, comments, tt.expectedLen)

			if tt.expectedType == "review" && len(comments) > 0 {
				assert.NotEmpty(t, comments[0].Path)
			}
			if tt.expectedType == "issue" && len(comments) > 0 {
				assert.Empty(t, comments[0].Path)
			}
		})
	}
}

func TestMockGitHubServer_HandleInvalidRequests(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "invalid repository path",
			url:            "/repos/",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid PR number",
			url:            "/repos/owner/repo/pulls/not-a-number",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid issue number",
			url:            "/repos/owner/repo/issues/not-a-number",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL() + tt.url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestMockComment_DefaultValues(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	// Add comment with minimal data
	comment := MockComment{
		Body: "Test",
		User: server.users["test-user"],
	}

	server.AddComment("owner/repo", 123, comment)
	comments := server.GetComments("owner/repo", 123)

	require.Len(t, comments, 1)
	addedComment := comments[0]

	// Verify defaults are set
	assert.Greater(t, addedComment.ID, 0)
	assert.False(t, addedComment.CreatedAt.IsZero())
	assert.False(t, addedComment.UpdatedAt.IsZero())
	assert.NotEmpty(t, addedComment.HTMLURL)
	assert.Contains(t, addedComment.HTMLURL, server.URL())
}

func TestMockComment_PreserveExistingValues(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	comment := MockComment{
		ID:        999,
		Body:      "Test",
		User:      server.users["test-user"],
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
		HTMLURL:   "https://example.com/custom",
	}

	server.AddComment("owner/repo", 123, comment)
	comments := server.GetComments("owner/repo", 123)

	require.Len(t, comments, 1)
	addedComment := comments[0]

	// Verify existing values are preserved
	assert.Equal(t, 999, addedComment.ID)
	assert.Equal(t, fixedTime, addedComment.CreatedAt)
	assert.Equal(t, fixedTime, addedComment.UpdatedAt)
	assert.Equal(t, "https://example.com/custom", addedComment.HTMLURL)
}

func TestMockGitHubServer_ThreadSafety(t *testing.T) {
	server := NewMockGitHubServer()
	defer server.Close()

	// Test concurrent access
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(i int) {
			comment := MockComment{
				Body: "Concurrent comment",
				User: server.users["test-user"],
			}
			server.AddComment("owner/repo", 123, comment)
			_ = server.GetComments("owner/repo", 123)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	comments := server.GetComments("owner/repo", 123)
	assert.Len(t, comments, 10)
}