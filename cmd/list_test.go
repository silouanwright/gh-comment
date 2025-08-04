package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllComments(t *testing.T) {
	// Create a mock client
	mockClient := &github.MockClient{
		IssueComments: []github.Comment{
			{
				ID:   1,
				Body: "Test issue comment",
				User: github.User{Login: "testuser"},
			},
		},
		ReviewComments: []github.Comment{
			{
				ID:   2,
				Body: "Test review comment",
				User: github.User{Login: "reviewer"},
				Path: "test.go",
				Line: 42,
			},
		},
	}

	comments, err := fetchAllComments(mockClient, "owner/repo", 123)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)

	// Check issue comment
	assert.Equal(t, 1, comments[0].ID)
	assert.Equal(t, "Test issue comment", comments[0].Body)
	assert.Equal(t, "testuser", comments[0].Author)
	assert.Equal(t, "issue", comments[0].Type)

	// Check review comment
	assert.Equal(t, 2, comments[1].ID)
	assert.Equal(t, "Test review comment", comments[1].Body)
	assert.Equal(t, "reviewer", comments[1].Author)
	assert.Equal(t, "review", comments[1].Type)
	assert.Equal(t, "test.go", comments[1].Path)
	assert.Equal(t, 42, comments[1].Line)
}

func TestFetchAllCommentsInvalidRepo(t *testing.T) {
	mockClient := &github.MockClient{}

	_, err := fetchAllComments(mockClient, "invalid", 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid repository format")
}

// Helper function to capture stdout output
func captureOutput(fn func()) string {
	// Save the original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to capture the output
	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		buf.ReadFrom(r)
		outputChan <- buf.String()
	}()

	// Execute the function
	fn()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Get the captured output
	return <-outputChan
}

func TestDisplayDiffHunk(t *testing.T) {
	tests := []struct {
		name          string
		diffHunk      string
		expectedLines []string // Lines that should be present in output
	}{
		{
			name: "simple diff with addition and removal",
			diffHunk: `@@ -10,7 +10,7 @@ func example() {
 	fmt.Println("hello")
-	old := "removed"
+	new := "added"
 	fmt.Println("world")`,
			expectedLines: []string{
				"🔹 @@ -10,7 +10,7 @@ func example() {",
				"fmt.Println(\"hello\")",
				"➖ -	old := \"removed\"",
				"➕ +	new := \"added\"",
				"fmt.Println(\"world\")",
			},
		},
		{
			name:     "diff header only",
			diffHunk: `@@ -1,3 +1,3 @@`,
			expectedLines: []string{
				"🔹 @@ -1,3 +1,3 @@",
			},
		},
		{
			name: "only additions",
			diffHunk: `@@ -0,0 +1,3 @@
+func newFunction() {
+	return true
+}`,
			expectedLines: []string{
				"🔹 @@ -0,0 +1,3 @@",
				"➕ +func newFunction() {",
				"➕ +	return true",
				"➕ +}",
			},
		},
		{
			name: "only removals",
			diffHunk: `@@ -1,3 +0,0 @@
-func oldFunction() {
-	return false
-}`,
			expectedLines: []string{
				"🔹 @@ -1,3 +0,0 @@",
				"➖ -func oldFunction() {",
				"➖ -	return false",
				"➖ -}",
			},
		},
		{
			name: "context lines only",
			diffHunk: `@@ -10,3 +10,3 @@ func context() {
 	line1
 	line2
 	line3`,
			expectedLines: []string{
				"🔹 @@ -10,3 +10,3 @@ func context() {",
				"line1",
				"line2",
				"line3",
			},
		},
		{
			name:          "empty diff hunk",
			diffHunk:      "",
			expectedLines: []string{}, // Should just output a newline
		},
		{
			name: "diff with empty lines",
			diffHunk: `@@ -1,5 +1,5 @@
 line1

 line3
+added line
 line5`,
			expectedLines: []string{
				"🔹 @@ -1,5 +1,5 @@",
				"line1",
				"line3",
				"➕ +added line",
				"line5",
			},
		},
		{
			name: "complex diff with multiple hunks",
			diffHunk: `@@ -1,3 +1,3 @@
 context line
-removed line
+added line
@@ -10,2 +10,3 @@
 another context
+another addition`,
			expectedLines: []string{
				"🔹 @@ -1,3 +1,3 @@",
				"context line",
				"➖ -removed line",
				"➕ +added line",
				"🔹 @@ -10,2 +10,3 @@",
				"another context",
				"➕ +another addition",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				displayDiffHunk(tt.diffHunk)
			})

			// Check that all expected lines are present in the output
			for _, expectedLine := range tt.expectedLines {
				assert.Contains(t, output, expectedLine, "Expected line not found in output: %s", expectedLine)
			}

			// For non-empty diff hunks, ensure we have proper formatting
			if tt.diffHunk != "" {
				lines := strings.Split(strings.TrimSpace(output), "\n")

				// Check that we have at least some output
				assert.Greater(t, len(lines), 0, "Should have at least one line of output")

				// Check that the output ends with a blank line (due to fmt.Println())
				assert.True(t, strings.HasSuffix(output, "\n"), "Output should end with newline")
			}
		})
	}
}

func TestDisplayDiffHunkEdgeCases(t *testing.T) {
	tests := []struct {
		name                string
		diffHunk            string
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name:             "whitespace only diff hunk",
			diffHunk:         "   \n   \n   ",
			expectedContains: []string{}, // Should handle gracefully
		},
		{
			name: "diff with special characters",
			diffHunk: `@@ -1,1 +1,1 @@
-const regex = /[.*+?^${}()|[\]\\]/g;
+const regex = /[.*+?^${}()|[\]\\]/gi;`,
			expectedContains: []string{
				"🔹 @@ -1,1 +1,1 @@",
				"➖ -const regex = /[.*+?^${}()|[\\]\\\\]/g;",
				"➕ +const regex = /[.*+?^${}()|[\\]\\\\]/gi;",
			},
		},
		{
			name: "diff with unicode characters",
			diffHunk: `@@ -1,2 +1,2 @@
-message := "Hello 世界"
+message := "Hello 🌍"`,
			expectedContains: []string{
				"➖ -message := \"Hello 世界\"",
				"➕ +message := \"Hello 🌍\"",
			},
		},
		{
			name: "very long lines in diff",
			diffHunk: `@@ -1,1 +1,1 @@
-` + strings.Repeat("a", 200) + `
+` + strings.Repeat("b", 200),
			expectedContains: []string{
				"➖ -" + strings.Repeat("a", 200),
				"➕ +" + strings.Repeat("b", 200),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				displayDiffHunk(tt.diffHunk)
			})

			// Check expected content
			for _, expected := range tt.expectedContains {
				assert.Contains(t, output, expected, "Expected content not found in output")
			}

			// Check that unwanted content is not present
			for _, notExpected := range tt.expectedNotContains {
				assert.NotContains(t, output, notExpected, "Unexpected content found in output")
			}
		})
	}
}

func TestDisplayDiffHunkFormatting(t *testing.T) {
	// Test that different line types get proper prefixes
	tests := []struct {
		name           string
		line           string
		expectedPrefix string
	}{
		{"diff header", "@@ -1,1 +1,1 @@", "🔹"},
		{"added line", "+new code", "➕"},
		{"removed line", "-old code", "➖"},
		{"context line", " unchanged", "    "}, // 4 spaces for context
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				displayDiffHunk(tt.line)
			})

			assert.Contains(t, output, tt.expectedPrefix, "Expected prefix not found")
			assert.Contains(t, output, tt.line, "Original line content not found")
		})
	}
}

// REGRESSION TEST: Ensure comment IDs are always displayed in list output
// This prevents the critical bug where users couldn't get comment IDs for replies/edits
func TestDisplayCommentAlwaysShowsID(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
	}{
		{
			name: "issue comment with ID",
			comment: Comment{
				ID:        12345678,
				Author:    "testuser",
				Body:      "Test issue comment",
				Type:      "issue",
				CreatedAt: testTime(),
			},
		},
		{
			name: "review comment with ID",
			comment: Comment{
				ID:        87654321,
				Author:    "reviewer",
				Body:      "Test review comment",
				Type:      "review",
				Path:      "test.go",
				Line:      42,
				CreatedAt: testTime(),
			},
		},
		{
			name: "comment with large ID",
			comment: Comment{
				ID:        2249431211, // Real GitHub comment ID format
				Author:    "realuser",
				Body:      "Real comment with large ID",
				Type:      "review",
				CreatedAt: testTime(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				displayComment(tt.comment, 1)
			})

			// CRITICAL: Must contain the actual GitHub comment ID
			expectedIDFormat := fmt.Sprintf("ID:%d", tt.comment.ID)
			assert.Contains(t, output, expectedIDFormat,
				"Comment ID must be displayed in format 'ID:%d' for integration testing workflows")

			// Should also contain the index number
			assert.Contains(t, output, "[1]", "Should contain index number")

			// Should contain author and body
			assert.Contains(t, output, tt.comment.Author, "Should contain author")
			assert.Contains(t, output, tt.comment.Body, "Should contain comment body")
		})
	}
}

// REGRESSION TEST: Ensure comment ID format is consistent
func TestDisplayCommentIDFormat(t *testing.T) {
	comment := Comment{
		ID:        123456789,
		Author:    "testuser",
		Body:      "Test comment",
		Type:      "issue",
		CreatedAt: testTime(),
	}

	output := captureOutput(func() {
		displayComment(comment, 1)
	})

	// Test exact format: [index] ID:actualID
	assert.Contains(t, output, "[1] ID:123456789",
		"Comment must display in exact format: [index] ID:actualID")
}

// REGRESSION TEST: Verify hideAuthors flag still shows comment IDs
func TestDisplayCommentWithHideAuthorsStillShowsID(t *testing.T) {
	// Save original state
	originalHideAuthors := hideAuthors
	defer func() { hideAuthors = originalHideAuthors }()

	// Enable hide authors flag
	hideAuthors = true

	comment := Comment{
		ID:        999888777,
		Author:    "secretuser",
		Body:      "Hidden author comment",
		Type:      "issue",
		CreatedAt: testTime(),
	}

	output := captureOutput(func() {
		displayComment(comment, 1)
	})

	// Even with hidden authors, ID must still be visible
	assert.Contains(t, output, "ID:999888777",
		"Comment ID must be visible even when authors are hidden")
	assert.Contains(t, output, "[hidden]", "Author should be hidden")
	assert.NotContains(t, output, "secretuser", "Author name should not appear")
}

func TestDisplayIDsOnly(t *testing.T) {
	tests := []struct {
		name     string
		comments []Comment
		expected []string
	}{
		{
			name: "single comment",
			comments: []Comment{
				{ID: 12345, Author: "user1", Body: "Test comment", Type: "issue"},
			},
			expected: []string{"12345"},
		},
		{
			name: "multiple comments",
			comments: []Comment{
				{ID: 111, Author: "user1", Body: "First", Type: "issue"},
				{ID: 222, Author: "user2", Body: "Second", Type: "review"},
				{ID: 333, Author: "user3", Body: "Third", Type: "issue"},
			},
			expected: []string{"111", "222", "333"},
		},
		{
			name: "large comment IDs",
			comments: []Comment{
				{ID: 2249431211, Author: "user1", Body: "Real GitHub ID", Type: "review"},
				{ID: 1987654321, Author: "user2", Body: "Another real ID", Type: "issue"},
			},
			expected: []string{"2249431211", "1987654321"},
		},
		{
			name:     "empty comments",
			comments: []Comment{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				displayIDsOnly(tt.comments)
			})

			lines := strings.Split(strings.TrimSpace(output), "\n")
			
			if len(tt.expected) == 0 {
				// Empty output should result in empty string when trimmed
				assert.Equal(t, "", strings.TrimSpace(output))
			} else {
				assert.Len(t, lines, len(tt.expected), "Should have correct number of lines")
				
				for i, expectedID := range tt.expected {
					assert.Equal(t, expectedID, lines[i], "ID should match at position %d", i)
				}
			}
		})
	}
}

func TestDisplayCommentsJSON(t *testing.T) {
	tests := []struct {
		name     string
		comments []Comment
		pr       int
		validate func(t *testing.T, output string)
	}{
		{
			name: "single comment JSON output",
			comments: []Comment{
				{
					ID:        12345,
					Author:    "testuser",
					Body:      "Test comment body",
					Type:      "issue",
					CreatedAt: testTime(),
				},
			},
			pr: 123,
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, `"pr": 123`)
				assert.Contains(t, output, `"total": 1`)
				assert.Contains(t, output, `"id": 12345`)
				assert.Contains(t, output, `"author": "testuser"`)
				assert.Contains(t, output, `"body": "Test comment body"`)
				assert.Contains(t, output, `"type": "issue"`)
				// Should be valid JSON
				assert.True(t, strings.HasPrefix(output, "{"))
				assert.True(t, strings.HasSuffix(strings.TrimSpace(output), "}"))
			},
		},
		{
			name: "multiple comments JSON output",
			comments: []Comment{
				{ID: 111, Author: "user1", Body: "First comment", Type: "issue", CreatedAt: testTime()},
				{ID: 222, Author: "user2", Body: "Second comment", Type: "review", Path: "test.go", Line: 42, CreatedAt: testTime()},
			},
			pr: 456,
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, `"pr": 456`)
				assert.Contains(t, output, `"total": 2`)
				assert.Contains(t, output, `"id": 111`)
				assert.Contains(t, output, `"id": 222`)
				assert.Contains(t, output, `"path": "test.go"`)
				assert.Contains(t, output, `"line": 42`)
				// Verify JSON structure
				assert.Contains(t, output, `"comments": [`)
			},
		},
		{
			name:     "empty comments JSON output",
			comments: []Comment{},
			pr:       789,
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, `"pr": 789`)
				assert.Contains(t, output, `"total": 0`)
				assert.Contains(t, output, `"comments": []`)
				// Should still be valid JSON
				assert.True(t, strings.HasPrefix(output, "{"))
				assert.True(t, strings.HasSuffix(strings.TrimSpace(output), "}"))
			},
		},
		{
			name: "comment with special characters in JSON",
			comments: []Comment{
				{
					ID:     999,
					Author: "test\"user",
					Body:   "Comment with \"quotes\" and \n newlines",
					Type:   "issue",
					CreatedAt: testTime(),
				},
			},
			pr: 100,
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, `"id": 999`)
				// JSON should properly escape quotes and newlines
				assert.Contains(t, output, `test\"user`)
				assert.Contains(t, output, `\"quotes\"`)
				assert.Contains(t, output, `\n`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				err := displayCommentsJSON(tt.comments, tt.pr)
				assert.NoError(t, err, "displayCommentsJSON should not return error")
			})

			tt.validate(t, output)
		})
	}
}

func TestDisplayComments(t *testing.T) {
	tests := []struct {
		name             string
		comments         []Comment
		pr               int
		expectedContains []string
	}{
		{
			name:     "empty comments",
			comments: []Comment{},
			pr:       123,
			expectedContains: []string{
				"No comments found on PR #123",
			},
		},
		{
			name: "single issue comment",
			comments: []Comment{
				{ID: 1, Author: "user1", Body: "Issue comment", Type: "issue", CreatedAt: testTime()},
			},
			pr: 456,
			expectedContains: []string{
				"📝 Comments on PR #456 (1 total)",
				"💬 General PR Comments (1)",
				"[1] ID:1",
				"Issue comment",
			},
		},
		{
			name: "single review comment",
			comments: []Comment{
				{ID: 2, Author: "reviewer", Body: "Review comment", Type: "review", Path: "test.go", Line: 42, CreatedAt: testTime()},
			},
			pr: 789,
			expectedContains: []string{
				"📝 Comments on PR #789 (1 total)",
				"📋 Review Comments (1)",
				"[1] ID:2",
				"Review comment",
			},
		},
		{
			name: "line-specific comment (other type)",
			comments: []Comment{
				{ID: 3, Author: "dev", Body: "Line comment", Type: "line", Path: "main.go", Line: 10, CreatedAt: testTime()},
			},
			pr: 111,
			expectedContains: []string{
				"📝 Comments on PR #111 (1 total)",
				"📍 Line-Specific Comments (1)",
				"[1] ID:3",
				"Line comment",
			},
		},
		{
			name: "mixed comment types",
			comments: []Comment{
				{ID: 1, Author: "user1", Body: "Issue comment", Type: "issue", CreatedAt: testTime()},
				{ID: 2, Author: "reviewer", Body: "Review comment", Type: "review", Path: "test.go", Line: 42, CreatedAt: testTime()},
				{ID: 3, Author: "dev", Body: "Line comment", Type: "line", Path: "main.go", Line: 10, CreatedAt: testTime()},
			},
			pr: 222,
			expectedContains: []string{
				"📝 Comments on PR #222 (3 total)",
				"💬 General PR Comments (1)",
				"📋 Review Comments (1)", 
				"📍 Line-Specific Comments (1)",
			},
		},
		{
			name: "comments with single commit ID",
			comments: []Comment{
				{ID: 1, Author: "user1", Body: "Comment", Type: "issue", CommitID: "abc123", CreatedAt: testTime()},
			},
			pr: 333,
			expectedContains: []string{
				"📝 Comments on PR #333 (1 total, 1 commit)",
			},
		},
		{
			name: "comments with multiple commit IDs",
			comments: []Comment{
				{ID: 1, Author: "user1", Body: "Comment 1", Type: "issue", CommitID: "abc123", CreatedAt: testTime()},
				{ID: 2, Author: "user2", Body: "Comment 2", Type: "review", CommitID: "def456", CreatedAt: testTime()},
				{ID: 3, Author: "user3", Body: "Comment 3", Type: "issue", CommitID: "abc123", CreatedAt: testTime()}, // Duplicate commit ID
			},
			pr: 444,
			expectedContains: []string{
				"📝 Comments on PR #444 (3 total, 2 commits)",
			},
		},
		{
			name: "comments with mixed commit IDs (some empty)",
			comments: []Comment{
				{ID: 1, Author: "user1", Body: "Comment 1", Type: "issue", CommitID: "abc123", CreatedAt: testTime()},
				{ID: 2, Author: "user2", Body: "Comment 2", Type: "review", CommitID: "", CreatedAt: testTime()}, // No commit ID
				{ID: 3, Author: "user3", Body: "Comment 3", Type: "issue", CommitID: "def456", CreatedAt: testTime()},
			},
			pr: 555,
			expectedContains: []string{
				"📝 Comments on PR #555 (3 total, 2 commits)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				displayComments(tt.comments, tt.pr)
			})

			for _, expected := range tt.expectedContains {
				assert.Contains(t, output, expected, "Should contain: %s", expected)
			}
		})
	}
}

func testTime() time.Time {
	// Return a fixed time for consistent testing
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}
