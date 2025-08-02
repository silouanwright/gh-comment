package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockGitHubAPI(t *testing.T) {
	server := MockGitHubAPI(t)
	defer server.Close()

	t.Run("mock issue comments GET", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/repos/owner/repo/issues/1/comments")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var comments []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&comments)
		require.NoError(t, err)

		assert.Len(t, comments, 1)
		assert.Equal(t, float64(123456), comments[0]["id"])
		assert.Equal(t, "This is a general PR comment", comments[0]["body"])
	})

	t.Run("mock issue comments POST", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/repos/owner/repo/issues/1/comments", "application/json", strings.NewReader(`{"body":"test"}`))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var comment map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&comment)
		require.NoError(t, err)

		assert.Equal(t, float64(789012), comment["id"])
		assert.Equal(t, "New comment", comment["body"])
	})

	t.Run("mock review comments GET", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/repos/owner/repo/pulls/1/comments")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var comments []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&comments)
		require.NoError(t, err)

		assert.Len(t, comments, 1)
		assert.Equal(t, float64(654321), comments[0]["id"])
		assert.Equal(t, "This is a line-specific review comment", comments[0]["body"])
		assert.Equal(t, "main.go", comments[0]["path"])
		assert.Equal(t, float64(42), comments[0]["line"])
	})

	t.Run("mock review comments POST", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/repos/owner/repo/pulls/1/comments", "application/json", strings.NewReader(`{"body":"test"}`))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var comment map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&comment)
		require.NoError(t, err)

		assert.Equal(t, float64(345678), comment["id"])
		assert.Equal(t, "New review comment", comment["body"])
	})

	t.Run("mock GraphQL endpoint", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/graphql", "application/json", strings.NewReader(`{"query":"test"}`))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		repo, ok := data["repository"].(map[string]interface{})
		assert.True(t, ok)
		assert.NotNil(t, repo["pullRequest"])
	})

	t.Run("unknown endpoint returns 404", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/unknown/endpoint")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestCaptureOutput(t *testing.T) {
	t.Run("capture stdout", func(t *testing.T) {
		stdout, stderr := CaptureOutput(func() {
			fmt.Print("hello stdout")
		})

		assert.Equal(t, "hello stdout", stdout)
		assert.Empty(t, stderr)
	})

	t.Run("capture stderr", func(t *testing.T) {
		stdout, stderr := CaptureOutput(func() {
			fmt.Fprint(os.Stderr, "hello stderr")
		})

		assert.Empty(t, stdout)
		assert.Equal(t, "hello stderr", stderr)
	})

	t.Run("capture both stdout and stderr", func(t *testing.T) {
		stdout, stderr := CaptureOutput(func() {
			fmt.Print("stdout message")
			fmt.Fprint(os.Stderr, "stderr message")
		})

		assert.Equal(t, "stdout message", stdout)
		assert.Equal(t, "stderr message", stderr)
	})

	t.Run("capture nothing", func(t *testing.T) {
		stdout, stderr := CaptureOutput(func() {
			// Do nothing
		})

		assert.Empty(t, stdout)
		assert.Empty(t, stderr)
	})
}

func TestLoadGoldenFile(t *testing.T) {
	// Create a temporary golden file for testing
	tempDir := t.TempDir()
	goldenDir := filepath.Join(tempDir, "testdata", "golden")
	err := os.MkdirAll(goldenDir, 0755)
	require.NoError(t, err)

	testData := []byte("test golden content")
	goldenFile := filepath.Join(goldenDir, "test.golden")
	err = os.WriteFile(goldenFile, testData, 0644)
	require.NoError(t, err)

	// Change to temp directory for test
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	t.Run("load existing golden file", func(t *testing.T) {
		data := LoadGoldenFile(t, "test.golden")
		assert.Equal(t, testData, data)
	})

	// Note: Testing the failure case of LoadGoldenFile would cause the test to fail
	// since it uses require.NoError internally. This is expected behavior.
}

func TestWriteGoldenFile(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	testData := []byte("test golden content")

	t.Run("write golden file creates directories", func(t *testing.T) {
		WriteGoldenFile(t, "test.golden", testData)

		// Verify file was created
		goldenPath := filepath.Join("testdata", "golden", "test.golden")
		data, err := os.ReadFile(goldenPath)
		require.NoError(t, err)
		assert.Equal(t, testData, data)
	})

	t.Run("write golden file in nested directory", func(t *testing.T) {
		WriteGoldenFile(t, "subdir/nested.golden", testData)

		// Verify file was created in nested directory
		goldenPath := filepath.Join("testdata", "golden", "subdir", "nested.golden")
		data, err := os.ReadFile(goldenPath)
		require.NoError(t, err)
		assert.Equal(t, testData, data)
	})
}

func TestCreateTestComments(t *testing.T) {
	comments := CreateTestComments()

	assert.Len(t, comments, 3)

	// Test first comment (issue)
	assert.Equal(t, 123456, comments[0].ID)
	assert.Equal(t, "LGTM! Great work on this PR.", comments[0].Body)
	assert.Equal(t, "issue", comments[0].Type)
	assert.Equal(t, "reviewer1", comments[0].User)
	assert.Equal(t, "2024-01-01T12:00:00Z", comments[0].CreatedAt)
	assert.Empty(t, comments[0].Path)
	assert.Zero(t, comments[0].Line)

	// Test second comment (review)
	assert.Equal(t, 654321, comments[1].ID)
	assert.Equal(t, "Consider using a more descriptive variable name here.", comments[1].Body)
	assert.Equal(t, "review", comments[1].Type)
	assert.Equal(t, "reviewer2", comments[1].User)
	assert.Equal(t, "2024-01-01T13:00:00Z", comments[1].CreatedAt)
	assert.Equal(t, "main.go", comments[1].Path)
	assert.Equal(t, 42, comments[1].Line)

	// Test third comment (issue)
	assert.Equal(t, 789012, comments[2].ID)
	assert.Equal(t, "Thanks for the feedback! I'll address this.", comments[2].Body)
	assert.Equal(t, "issue", comments[2].Type)
	assert.Equal(t, "author", comments[2].User)
	assert.Equal(t, "2024-01-01T14:00:00Z", comments[2].CreatedAt)
	assert.Empty(t, comments[2].Path)
	assert.Zero(t, comments[2].Line)
}

func TestAssertGoldenMatch(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWd)
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	testContent := "expected output content"

	t.Run("with UPDATE_GOLDEN=1 creates golden file", func(t *testing.T) {
		// Set environment variable
		originalEnv := os.Getenv("UPDATE_GOLDEN")
		os.Setenv("UPDATE_GOLDEN", "1")
		defer func() {
			if originalEnv != "" {
				os.Setenv("UPDATE_GOLDEN", originalEnv)
			} else {
				os.Unsetenv("UPDATE_GOLDEN")
			}
		}()

		AssertGoldenMatch(t, "test.golden", testContent)

		// Verify file was created
		goldenPath := filepath.Join("testdata", "golden", "test.golden")
		data, err := os.ReadFile(goldenPath)
		require.NoError(t, err)
		assert.Equal(t, testContent, string(data))
	})

	t.Run("matches existing golden file", func(t *testing.T) {
		// First create the golden file
		goldenDir := filepath.Join("testdata", "golden")
		err := os.MkdirAll(goldenDir, 0755)
		require.NoError(t, err)
		goldenPath := filepath.Join(goldenDir, "match.golden")
		err = os.WriteFile(goldenPath, []byte(testContent), 0644)
		require.NoError(t, err)

		// Now test matching
		AssertGoldenMatch(t, "match.golden", testContent)
		// If we get here, the assertion passed
	})

	// Note: Testing the failure case would cause the test to fail,
	// so we skip testing mismatched content scenarios
}

func TestConstants(t *testing.T) {
	// Test that constants are set to expected values
	assert.Equal(t, 100, MaxGraphQLResults)
	assert.Equal(t, 65536, MaxCommentLength)
	assert.Equal(t, 30, DefaultPageSize)

	// Verify they're positive values
	assert.Greater(t, MaxGraphQLResults, 0)
	assert.Greater(t, MaxCommentLength, 0)
	assert.Greater(t, DefaultPageSize, 0)
}

func TestTestCommentStruct(t *testing.T) {
	comment := TestComment{
		ID:        123456,
		Body:      "Test comment body",
		Type:      "review",
		User:      "testuser",
		CreatedAt: "2024-01-01T12:00:00Z",
		Path:      "src/main.go",
		Line:      42,
	}

	assert.Equal(t, 123456, comment.ID)
	assert.Equal(t, "Test comment body", comment.Body)
	assert.Equal(t, "review", comment.Type)
	assert.Equal(t, "testuser", comment.User)
	assert.Equal(t, "2024-01-01T12:00:00Z", comment.CreatedAt)
	assert.Equal(t, "src/main.go", comment.Path)
	assert.Equal(t, 42, comment.Line)

	// Test JSON marshaling/unmarshaling
	jsonData, err := json.Marshal(comment)
	require.NoError(t, err)

	var unmarshaled TestComment
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, comment, unmarshaled)
}