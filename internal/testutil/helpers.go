package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Constants for testing
const (
	MaxGraphQLResults = 100
	MaxCommentLength  = 65536
	DefaultPageSize   = 30
)

// MockGitHubAPI creates a mock HTTP server for GitHub API testing
func MockGitHubAPI(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default mock responses based on endpoint
		switch {
		case r.URL.Path == "/repos/owner/repo/issues/1/comments":
			mockIssueComments(w, r)
		case r.URL.Path == "/repos/owner/repo/pulls/1/comments":
			mockReviewComments(w, r)
		case r.URL.Path == "/graphql":
			mockGraphQL(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}

// CaptureOutput captures stdout and stderr from a function
func CaptureOutput(fn func()) (stdout, stderr string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	os.Stdout = wOut
	os.Stderr = wErr
	
	outC := make(chan string)
	errC := make(chan string)
	
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		outC <- buf.String()
	}()
	
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rErr)
		errC <- buf.String()
	}()
	
	fn()
	
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	
	stdout = <-outC
	stderr = <-errC
	
	return
}

// LoadGoldenFile loads a golden file for comparison
func LoadGoldenFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", "golden", name)
	data, err := os.ReadFile(path)
	require.NoError(t, err, "failed to load golden file %s", path)
	return data
}

// WriteGoldenFile writes a golden file (for updating test expectations)
func WriteGoldenFile(t *testing.T, name string, data []byte) {
	path := filepath.Join("testdata", "golden", name)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	require.NoError(t, err)
	err = os.WriteFile(path, data, 0644)
	require.NoError(t, err)
}

// Mock response helpers
func mockIssueComments(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		comments := []map[string]interface{}{
			{
				"id":   123456,
				"body": "This is a general PR comment",
				"user": map[string]interface{}{
					"login": "testuser",
				},
				"created_at": "2024-01-01T12:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(comments)
	} else if r.Method == "POST" {
		// Mock creating a new comment
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   789012,
			"body": "New comment",
		})
	}
}

func mockReviewComments(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		comments := []map[string]interface{}{
			{
				"id":   654321,
				"body": "This is a line-specific review comment",
				"user": map[string]interface{}{
					"login": "reviewer",
				},
				"created_at": "2024-01-01T13:00:00Z",
				"path":       "main.go",
				"line":       42,
			},
		}
		json.NewEncoder(w).Encode(comments)
	} else if r.Method == "POST" {
		// Mock creating a new review comment
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   345678,
			"body": "New review comment",
		})
	}
}

func mockGraphQL(w http.ResponseWriter, r *http.Request) {
	// Mock GraphQL responses for resolve functionality
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"repository": map[string]interface{}{
				"pullRequest": map[string]interface{}{
					"reviewThreads": map[string]interface{}{
						"nodes": []map[string]interface{}{
							{
								"id": "RT_123",
								"comments": map[string]interface{}{
									"nodes": []map[string]interface{}{
										{
											"id": "654321",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

// TestComment represents a test comment structure
type TestComment struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Type      string `json:"type"` // "issue" or "review"
	User      string `json:"user"`
	CreatedAt string `json:"created_at"`
	Path      string `json:"path,omitempty"`
	Line      int    `json:"line,omitempty"`
}

// CreateTestComments creates a set of test comments for testing
func CreateTestComments() []TestComment {
	return []TestComment{
		{
			ID:        123456,
			Body:      "LGTM! Great work on this PR.",
			Type:      "issue",
			User:      "reviewer1",
			CreatedAt: "2024-01-01T12:00:00Z",
		},
		{
			ID:        654321,
			Body:      "Consider using a more descriptive variable name here.",
			Type:      "review",
			User:      "reviewer2",
			CreatedAt: "2024-01-01T13:00:00Z",
			Path:      "main.go",
			Line:      42,
		},
		{
			ID:        789012,
			Body:      "Thanks for the feedback! I'll address this.",
			Type:      "issue",
			User:      "author",
			CreatedAt: "2024-01-01T14:00:00Z",
		},
	}
}

// AssertGoldenMatch compares output with golden file
func AssertGoldenMatch(t *testing.T, goldenFile string, actual string) {
	goldenPath := filepath.Join("testdata", "golden", goldenFile)
	
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		// Update golden file
		err := os.MkdirAll(filepath.Dir(goldenPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(goldenPath, []byte(actual), 0644)
		require.NoError(t, err)
		return
	}
	
	expected, err := os.ReadFile(goldenPath)
	if os.IsNotExist(err) {
		t.Fatalf("Golden file %s does not exist. Run with UPDATE_GOLDEN=1 to create it.", goldenPath)
	}
	require.NoError(t, err)
	
	require.Equal(t, string(expected), actual, "Output doesn't match golden file %s", goldenFile)
}
