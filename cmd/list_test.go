package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

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
