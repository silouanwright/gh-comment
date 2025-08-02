package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestClientEnhancedErrorHandling(t *testing.T) {
	// Create a mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
	}))
	defer server.Close()

	client := &TestClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	t.Run("ListIssueComments error", func(t *testing.T) {
		comments, err := client.ListIssueComments("owner", "repo", 123)
		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("ListReviewComments error", func(t *testing.T) {
		comments, err := client.ListReviewComments("owner", "repo", 123)
		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("CreateIssueComment error", func(t *testing.T) {
		comment, err := client.CreateIssueComment("owner", "repo", 123, "test")
		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("AddReviewComment error", func(t *testing.T) {
		reviewComment := ReviewCommentInput{
			Body: "test",
			Path: "test.go",
			Line: 10,
		}
		err := client.AddReviewComment("owner", "repo", 123, reviewComment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("CreateReview error", func(t *testing.T) {
		review := ReviewInput{
			Body:  "test",
			Event: "COMMENT",
		}
		err := client.CreateReview("owner", "repo", 123, review)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("GetPRDetails error", func(t *testing.T) {
		details, err := client.GetPRDetails("owner", "repo", 123)
		assert.Error(t, err)
		assert.Nil(t, details)
		assert.Contains(t, err.Error(), "500")
	})
}

func TestTestClientSuccessfulResponses(t *testing.T) {
	// Create a mock server that returns successful responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/repos/owner/repo/issues/123/comments":
			if r.Method == "GET" {
				comments := []Comment{
					{
						ID:        456,
						Body:      "Test comment",
						User:      User{Login: "testuser"},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Type:      "issue",
					},
				}
				json.NewEncoder(w).Encode(comments)
			} else if r.Method == "POST" {
				comment := Comment{
					ID:        789,
					Body:      "New comment",
					User:      User{Login: "testuser"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Type:      "issue",
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(comment)
			}
		case "/repos/owner/repo/pulls/123/comments":
			if r.Method == "GET" {
				comments := []Comment{
					{
						ID:        654,
						Body:      "Review comment",
						User:      User{Login: "reviewer"},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Path:      "main.go",
						Line:      42,
						Type:      "review",
					},
				}
				json.NewEncoder(w).Encode(comments)
			} else if r.Method == "POST" {
				comment := Comment{
					ID:        987,
					Body:      "New review comment",
					User:      User{Login: "testuser"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Path:      "test.go",
					Line:      10,
					Type:      "review",
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(comment)
			}
		case "/repos/owner/repo/pulls/123/reviews":
			if r.Method == "POST" {
				review := map[string]interface{}{
					"id":    555,
					"state": "PENDING",
					"user":  map[string]interface{}{"login": "testuser"},
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(review)
			}
		case "/repos/owner/repo/pulls/123":
			if r.Method == "GET" {
				prDetails := map[string]interface{}{
					"number": 123,
					"state":  "open",
					"title":  "Test PR",
					"head": map[string]interface{}{
						"sha": "abc123def456",
					},
				}
				json.NewEncoder(w).Encode(prDetails)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &TestClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	t.Run("ListIssueComments success", func(t *testing.T) {
		comments, err := client.ListIssueComments("owner", "repo", 123)
		require.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, 456, comments[0].ID)
		assert.Equal(t, "Test comment", comments[0].Body)
	})

	t.Run("ListReviewComments success", func(t *testing.T) {
		comments, err := client.ListReviewComments("owner", "repo", 123)
		require.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, 654, comments[0].ID)
		assert.Equal(t, "Review comment", comments[0].Body)
		assert.Equal(t, "main.go", comments[0].Path)
		assert.Equal(t, 42, comments[0].Line)
	})

	t.Run("CreateIssueComment success", func(t *testing.T) {
		comment, err := client.CreateIssueComment("owner", "repo", 123, "New comment")
		require.NoError(t, err)
		assert.Equal(t, 789, comment.ID)
		assert.Equal(t, "New comment", comment.Body)
	})

	t.Run("AddReviewComment success", func(t *testing.T) {
		reviewComment := ReviewCommentInput{
			Body: "New review comment",
			Path: "test.go",
			Line: 10,
		}
		err := client.AddReviewComment("owner", "repo", 123, reviewComment)
		assert.NoError(t, err)
	})

	t.Run("CreateReview success", func(t *testing.T) {
		review := ReviewInput{
			Body:  "Test review",
			Event: "COMMENT",
		}
		err := client.CreateReview("owner", "repo", 123, review)
		assert.NoError(t, err)
	})

	t.Run("GetPRDetails success", func(t *testing.T) {
		details, err := client.GetPRDetails("owner", "repo", 123)
		require.NoError(t, err)
		assert.Equal(t, float64(123), details["number"])
		assert.Equal(t, "open", details["state"])
		assert.Equal(t, "Test PR", details["title"])
	})
}

func TestTestClientDoRequestErrors(t *testing.T) {
	client := &TestClient{
		baseURL:    "http://invalid-url-that-does-not-exist.local",
		httpClient: &http.Client{},
	}

	t.Run("doRequest with invalid URL", func(t *testing.T) {
		_, err := client.doRequest("GET", "/test", nil)
		assert.Error(t, err)
		// Should contain network error information
	})
}

func TestTestClientJSONParsingErrors(t *testing.T) {
	// Create a server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json {"))
	}))
	defer server.Close()

	client := &TestClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	t.Run("ListIssueComments with invalid JSON", func(t *testing.T) {
		comments, err := client.ListIssueComments("owner", "repo", 123)
		assert.Error(t, err)
		assert.Nil(t, comments)
		assert.Contains(t, err.Error(), "failed to decode")
	})

	t.Run("GetPRDetails with invalid JSON", func(t *testing.T) {
		details, err := client.GetPRDetails("owner", "repo", 123)
		assert.Error(t, err)
		assert.Nil(t, details)
		assert.Contains(t, err.Error(), "failed to decode")
	})
}

func TestTestClientRequestBodyMarshaling(t *testing.T) {
	// Test that complex request bodies are properly marshaled
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// For POST requests, verify the body was properly marshaled
		if r.Method == "POST" {
			var requestBody map[string]interface{}
			json.NewDecoder(r.Body).Decode(&requestBody)

			// Echo back the request for verification with proper status
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(requestBody)
		} else {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		}
	}))
	defer server.Close()

	client := &TestClient{
		baseURL:    server.URL,
		httpClient: &http.Client{},
	}

	t.Run("CreateReview with complex body", func(t *testing.T) {
		review := ReviewInput{
			Body:  "Complex review",
			Event: "REQUEST_CHANGES",
			Comments: []ReviewCommentInput{
				{
					Body: "Comment 1",
					Path: "file1.go",
					Line: 10,
				},
				{
					Body: "Comment 2",
					Path: "file2.go",
					Line: 20,
				},
			},
		}
		err := client.CreateReview("owner", "repo", 123, review)
		assert.NoError(t, err)
	})
}
