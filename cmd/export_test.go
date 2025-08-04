package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunExport(t *testing.T) {
	// Save original state
	originalClient := exportClient
	originalRepo := repo
	originalPRNumber := prNumber
	defer func() {
		exportClient = originalClient
		repo = originalRepo
		prNumber = originalPRNumber
	}()

	tests := []struct {
		name           string
		args           []string
		setupRepo      string
		setupPR        int
		expectedFormat string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:           "export with PR number",
			args:           []string{"123"},
			setupRepo:      "owner/repo",
			expectedFormat: "json",
			wantErr:        false,
		},
		{
			name:           "export with auto-detect PR",
			args:           []string{},
			setupRepo:      "owner/repo",
			setupPR:        456,
			expectedFormat: "json",
			wantErr:        false,
		},
		{
			name:           "invalid PR number",
			args:           []string{"invalid"},
			setupRepo:      "owner/repo",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock client
			mockClient := github.NewMockClient()
			exportClient = mockClient

			// Set up repository and PR context
			repo = tt.setupRepo
			prNumber = tt.setupPR

			// Setup mock data
			mockClient.IssueComments = []github.Comment{
				{ID: 1, User: github.User{Login: "alice"}, Body: "Test comment", Type: "issue"},
			}
			mockClient.ReviewComments = []github.Comment{
				{ID: 2, User: github.User{Login: "bob"}, Body: "Test review", Path: "test.go", Line: 10},
			}

			// Capture output
			exportOutput = ""
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := runExport(exportCmd, tt.args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.expectedFormat == "json" {
					// Should be valid JSON
					var result []map[string]interface{}
					err := json.Unmarshal([]byte(output), &result)
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestExportFormats(t *testing.T) {
	comments := []ExportComment{
		{
			ID:        123,
			Type:      "issue",
			Author:    "alice",
			Body:      "Test issue comment",
			CreatedAt: time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC),
			URL:       "https://github.com/owner/repo/issues/123#issuecomment-123",
		},
		{
			ID:        456,
			Type:      "review",
			Author:    "bob",
			Body:      "Test review comment",
			File:      "main.go",
			Line:      42,
			CreatedAt: time.Date(2023, 1, 15, 11, 0, 0, 0, time.UTC),
			URL:       "https://github.com/owner/repo/pull/123#discussion_r456",
			DiffHunk:  "@@ -40,4 +40,4 @@\n func main() {\n-    fmt.Println(\"old\")\n+    fmt.Println(\"new\")\n }",
			Resolved:  true,
		},
	}

	t.Run("JSON export", func(t *testing.T) {
		var buf bytes.Buffer
		err := exportJSON(&buf, comments)
		assert.NoError(t, err)

		var result []ExportComment
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 123, result[0].ID)
		assert.Equal(t, "alice", result[0].Author)
	})

	t.Run("JSON with field filter", func(t *testing.T) {
		exportInclude = []string{"id", "author", "body"}
		defer func() { exportInclude = []string{} }()

		var buf bytes.Buffer
		err := exportJSON(&buf, comments)
		assert.NoError(t, err)

		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Contains(t, result[0], "id")
		assert.Contains(t, result[0], "author")
		assert.Contains(t, result[0], "body")
		assert.NotContains(t, result[0], "type")
	})

	t.Run("JSON with all field types", func(t *testing.T) {
		// Test all possible field types in the filter
		exportInclude = []string{"id", "type", "author", "body", "file", "line", "created_at", "updated_at", "url", "diff_hunk", "commit_id", "in_reply_to", "resolved"}
		defer func() { exportInclude = []string{} }()
		
		var buf bytes.Buffer
		err := exportJSON(&buf, comments)
		assert.NoError(t, err)
		
		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		
		// Check first comment (issue comment)
		comment := result[0]
		assert.Equal(t, float64(123), comment["id"])
		assert.Equal(t, "issue", comment["type"]) 
		assert.Equal(t, "alice", comment["author"])
		assert.Equal(t, "Test issue comment", comment["body"])
		assert.Contains(t, comment, "created_at")
		assert.Contains(t, comment, "url")
		assert.Equal(t, false, comment["resolved"])
		// File and line should not be present for issue comment
		assert.NotContains(t, comment, "file")
		assert.NotContains(t, comment, "line")
		
		// Check second comment (review comment)
		reviewComment := result[1]
		assert.Equal(t, float64(456), reviewComment["id"])
		assert.Equal(t, "review", reviewComment["type"])
		assert.Equal(t, "bob", reviewComment["author"])
		assert.Equal(t, "main.go", reviewComment["file"])
		assert.Equal(t, float64(42), reviewComment["line"])
		assert.Equal(t, "@@ -40,4 +40,4 @@\n func main() {\n-    fmt.Println(\"old\")\n+    fmt.Println(\"new\")\n }", reviewComment["diff_hunk"])
		assert.Equal(t, true, reviewComment["resolved"])
	})

	t.Run("JSON with empty field values", func(t *testing.T) {
		exportInclude = []string{"file", "line", "diff_hunk", "commit_id", "in_reply_to"}
		defer func() { exportInclude = []string{} }()
		
		// Use the issue comment which has empty file/line/diff_hunk
		var buf bytes.Buffer
		err := exportJSON(&buf, []ExportComment{comments[0]}) // Only issue comment
		assert.NoError(t, err)
		
		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		
		// Empty/zero fields should not be included
		comment := result[0]
		assert.NotContains(t, comment, "file")
		assert.NotContains(t, comment, "line")
		assert.NotContains(t, comment, "diff_hunk")
		assert.NotContains(t, comment, "commit_id")
		assert.NotContains(t, comment, "in_reply_to")
	})

	t.Run("JSON with single field filter", func(t *testing.T) {
		exportInclude = []string{"id"}
		defer func() { exportInclude = []string{} }()
		
		var buf bytes.Buffer
		err := exportJSON(&buf, comments)
		assert.NoError(t, err)
		
		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		
		// Should only contain ID field
		comment := result[0]
		assert.Len(t, comment, 1)
		assert.Contains(t, comment, "id")
		assert.Equal(t, float64(123), comment["id"])
	})

	t.Run("JSON with mixed field types", func(t *testing.T) {
		exportInclude = []string{"author", "resolved", "line"}
		defer func() { exportInclude = []string{} }()
		
		var buf bytes.Buffer
		err := exportJSON(&buf, comments)
		assert.NoError(t, err)
		
		var result []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		
		// First comment (issue) - no line field
		comment := result[0]
		assert.Contains(t, comment, "author")
		assert.Contains(t, comment, "resolved")
		assert.NotContains(t, comment, "line") // Issue comment has no line
		
		// Second comment (review) - has line field
		reviewComment := result[1]
		assert.Contains(t, reviewComment, "author")
		assert.Contains(t, reviewComment, "resolved")
		assert.Contains(t, reviewComment, "line")
		assert.Equal(t, float64(42), reviewComment["line"])
	})

	t.Run("CSV export", func(t *testing.T) {
		var buf bytes.Buffer
		err := exportCSV(&buf, comments)
		assert.NoError(t, err)

		output := buf.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Len(t, lines, 3) // header + 2 data rows
		
		// Check header
		assert.Contains(t, lines[0], "ID,Type,Author")
		
		// Check data
		assert.Contains(t, lines[1], "123,issue,alice")
		assert.Contains(t, lines[2], "456,review,bob")
	})

	t.Run("Markdown export", func(t *testing.T) {
		var buf bytes.Buffer
		err := exportMarkdown(&buf, comments, "owner/repo", 123)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "# PR Comments Export")
		assert.Contains(t, output, "**Repository:** owner/repo")
		assert.Contains(t, output, "**PR:** #123")
		assert.Contains(t, output, "## General PR Comments (1)")
		assert.Contains(t, output, "## Review Comments (1)")
		assert.Contains(t, output, "**Author:** @alice")
		assert.Contains(t, output, "**Author:** @bob")
		assert.Contains(t, output, "**File:** `main.go:42`")
		assert.Contains(t, output, "**Status:** ✅ Resolved")
	})

	t.Run("HTML export", func(t *testing.T) {
		var buf bytes.Buffer
		err := exportHTML(&buf, comments, "owner/repo", 123)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "<!DOCTYPE html>")
		assert.Contains(t, output, "<title>PR #123 Comments - owner/repo</title>")
		assert.Contains(t, output, "<h1>PR #123 Comments Export</h1>")
		assert.Contains(t, output, "<h2>General PR Comments (1)</h2>")
		assert.Contains(t, output, "<h2>Review Comments (1)</h2>")
		assert.Contains(t, output, "class=\"author\">@alice")
		assert.Contains(t, output, "class=\"author\">@bob")
		assert.Contains(t, output, "class=\"resolved\">✅ Resolved")
	})
}

func TestExportFormatValidation(t *testing.T) {
	// Save original state
	originalClient := exportClient
	originalRepo := repo
	defer func() {
		exportClient = originalClient
		repo = originalRepo
	}()

	// Setup
	mockClient := github.NewMockClient()
	exportClient = mockClient
	repo = "owner/repo"

	tests := []struct {
		name           string
		format         string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:    "valid json format",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "valid csv format",
			format:  "csv",
			wantErr: false,
		},
		{
			name:    "valid markdown format",
			format:  "markdown",
			wantErr: false,
		},
		{
			name:    "valid html format",
			format:  "html",
			wantErr: false,
		},
		{
			name:           "invalid format",
			format:         "xml",
			wantErr:        true,
			expectedErrMsg: "invalid format: xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exportFormat = tt.format
			
			// Redirect stdout to avoid output during tests
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := runExport(exportCmd, []string{"123"})

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout
			buf := make([]byte, 1024)
			r.Read(buf) // drain the pipe

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

func TestFetchAllCommentsForExport(t *testing.T) {
	mockClient := github.NewMockClient()

	// Setup mock data  
	issueComments := []github.Comment{
		{ID: 1, User: github.User{Login: "alice"}, Body: "Issue comment 1", Type: "issue"},
		{ID: 2, User: github.User{Login: "bob"}, Body: "Issue comment 2", Type: "issue"},
	}
	reviewComments := []github.Comment{
		{ID: 3, User: github.User{Login: "charlie"}, Body: "Review comment 1", Path: "test.go", Line: 10},
		{ID: 4, User: github.User{Login: "dave"}, Body: "Review comment 2", Path: "main.go", Line: 20},
	}

	mockClient.IssueComments = issueComments
	mockClient.ReviewComments = reviewComments

	comments, err := fetchAllCommentsForExport(mockClient, "owner", "repo", 123)
	assert.NoError(t, err)
	assert.Len(t, comments, 4)

	// Check issue comments
	assert.Equal(t, "issue", comments[0].Type)
	assert.Equal(t, "alice", comments[0].Author)
	assert.Equal(t, "issue", comments[1].Type)
	assert.Equal(t, "bob", comments[1].Author)

	// Check review comments
	assert.Equal(t, "review", comments[2].Type)
	assert.Equal(t, "charlie", comments[2].Author)
	assert.Equal(t, "test.go", comments[2].File)
	assert.Equal(t, 10, comments[2].Line)
	
	assert.Equal(t, "review", comments[3].Type)
	assert.Equal(t, "dave", comments[3].Author)
	assert.Equal(t, "main.go", comments[3].File)
	assert.Equal(t, 20, comments[3].Line)
	// Note: Resolved field is not available in current Comment struct
}

func TestExportWithResolvedFilter(t *testing.T) {
	// Save original state
	originalClient := exportClient
	originalRepo := repo
	originalIncludeResolved := includeResolved
	originalFormat := exportFormat
	defer func() {
		exportClient = originalClient
		repo = originalRepo
		includeResolved = originalIncludeResolved
		exportFormat = originalFormat
	}()

	// Setup
	mockClient := github.NewMockClient()
	exportClient = mockClient
	repo = "owner/repo"
	exportFormat = "json"

	// Setup mock data (note: resolved status not available in current Comment struct)
	reviewComments := []github.Comment{
		{ID: 1, User: github.User{Login: "alice"}, Body: "Comment 1", Path: "test.go", Line: 10},
		{ID: 2, User: github.User{Login: "bob"}, Body: "Comment 2", Path: "main.go", Line: 20},
	}
	// Clear default comments and set only our test data
	mockClient.IssueComments = []github.Comment{}
	mockClient.ReviewComments = reviewComments

	// Note: Since resolved status is not available in current Comment struct,
	// this test just verifies basic export functionality
	t.Run("basic export functionality", func(t *testing.T) {
		includeResolved = true
		
		// Capture output
		exportOutput = ""
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := runExport(exportCmd, []string{"123"})
		assert.NoError(t, err)

		// Restore stdout and read output
		w.Close()
		os.Stdout = oldStdout
		outputBytes := make([]byte, 1024)
		n, _ := r.Read(outputBytes)
		output := string(outputBytes[:n])

		// Parse JSON output
		var result []ExportComment
		err = json.Unmarshal([]byte(output), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2) // Both comments
		assert.Equal(t, "alice", result[0].Author)
		assert.Equal(t, "bob", result[1].Author)
	})
}