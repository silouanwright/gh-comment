package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/silouanwright/gh-comment/internal/testutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		flags          map[string]string
		mockComments   []github.Comment
		wantErr        bool
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "list basic comments",
			args: []string{"123"},
			mockComments: []github.Comment{
				{
					ID:        123456,
					Body:      "LGTM! Great work on this PR.",
					Type:      "issue",
					User:      github.User{Login: "reviewer1"},
					CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				},
				{
					ID:        654321,
					Body:      "Consider using a more descriptive variable name here.",
					Type:      "review",
					User:      github.User{Login: "reviewer2"},
					CreatedAt: time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
					Path:      "main.go",
					Line:      42,
				},
			},
			wantContains: []string{
				"💬 General PR Comments (1)",
				"📋 Review Comments (1)",
				"LGTM! Great work on this PR.",
				"Consider using a more descriptive variable name here.",
				"reviewer1",
				"reviewer2",
				"main.go:L42",
			},
		},
		{
			name: "quiet mode",
			args: []string{"123"},
			flags: map[string]string{
				"quiet": "true",
			},
			mockComments: []github.Comment{
				{
					ID:        123456,
					Body:      "Test comment",
					Type:      "issue",
					User:      github.User{Login: "testuser"},
					CreatedAt: time.Now(),
				},
			},
			wantContains:   []string{"Test comment", "testuser"},
			wantNotContain: []string{"🔗"}, // URLs should be hidden in quiet mode
		},
		{
			name: "hide authors",
			args: []string{"123"},
			flags: map[string]string{
				"hide-authors": "true",
			},
			mockComments: []github.Comment{
				{
					ID:        123456,
					Body:      "Test comment",
					Type:      "issue",
					User:      github.User{Login: "testuser"},
					CreatedAt: time.Now(),
				},
			},
			wantContains:   []string{"Test comment", "[hidden]"},
			wantNotContain: []string{"testuser"},
		},
		{
			name: "filter by author",
			args: []string{"123"},
			flags: map[string]string{
				"author": "reviewer1",
			},
			mockComments: []github.Comment{
				{
					ID:        123456,
					Body:      "Comment from reviewer1",
					Type:      "issue",
					User:      github.User{Login: "reviewer1"},
					CreatedAt: time.Now(),
				},
				{
					ID:        654321,
					Body:      "Comment from reviewer2",
					Type:      "issue",
					User:      github.User{Login: "reviewer2"},
					CreatedAt: time.Now(),
				},
			},
			wantContains:   []string{"Comment from reviewer1", "reviewer1"},
			wantNotContain: []string{"Comment from reviewer2", "reviewer2"},
		},
		{
			name:    "invalid PR number",
			args:    []string{"invalid"},
			wantErr: true,
		},
		{
			name:         "no comments",
			args:         []string{"123"},
			mockComments: []github.Comment{},
			wantContains: []string{"No comments found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flags
			resetListFlags()

			// Set up mock client
			mockClient := github.NewMockClient()
			if tt.mockComments != nil {
				// Separate comments by type
				var issueComments, reviewComments []github.Comment
				for _, comment := range tt.mockComments {
					if comment.Type == "issue" {
						issueComments = append(issueComments, comment)
					} else {
						reviewComments = append(reviewComments, comment)
					}
				}
				mockClient.IssueComments = issueComments
				mockClient.ReviewComments = reviewComments
			}

			// Create command with dependency injection
			cmd := NewListCmdWithDeps(&Dependencies{
				GitHubClient: mockClient,
				Output:       &bytes.Buffer{},
			})

			// Set flags
			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				require.NoError(t, err)
			}

			// Capture output using testutil
			var outputStr string
			stdout, stderr := testutil.CaptureOutput(func() {
				cmd.SetArgs(tt.args)
				err := cmd.Execute()
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
			})
			
			if tt.wantErr {
				return
			}
			
			outputStr = stdout + stderr
			for _, want := range tt.wantContains {
				assert.Contains(t, outputStr, want, "Output should contain: %s", want)
			}
			for _, notWant := range tt.wantNotContain {
				assert.NotContains(t, outputStr, notWant, "Output should not contain: %s", notWant)
			}
		})
	}
}

func TestFilterComments(t *testing.T) {
	comments := []Comment{
		{
			ID:     1,
			Author: "user1",
			Body:   "Comment 1",
			Type:   "issue",
		},
		{
			ID:     2,
			Author: "user2",
			Body:   "Comment 2",
			Type:   "review",
		},
		{
			ID:     3,
			Author: "user1",
			Body:   "Comment 3",
			Type:   "issue",
		},
	}

	tests := []struct {
		name     string
		author   string
		expected int
	}{
		{
			name:     "no filter",
			author:   "",
			expected: 3,
		},
		{
			name:     "filter by user1",
			author:   "user1",
			expected: 2,
		},
		{
			name:     "filter by user2",
			author:   "user2",
			expected: 1,
		},
		{
			name:     "filter by nonexistent user",
			author:   "nonexistent",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global author filter
			author = tt.author
			
			filtered := filterComments(comments)
			assert.Len(t, filtered, tt.expected)
			
			// Verify all filtered comments match the author
			if tt.author != "" {
				for _, comment := range filtered {
					assert.Equal(t, tt.author, comment.Author)
				}
			}
		})
	}
}

func TestFormatTimeAgo(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "just now",
			time:     now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "5 minutes ago",
			time:     now.Add(-5 * time.Minute),
			expected: "5 minutes ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "3 hours ago",
			time:     now.Add(-3 * time.Hour),
			expected: "3 hours ago",
		},
		{
			name:     "1 day ago",
			time:     now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "3 days ago",
			time:     now.Add(-3 * 24 * time.Hour),
			expected: "3 days ago",
		},
		{
			name:     "old date",
			time:     time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: "Jan 15, 2023",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeAgo(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDisplayDiffHunk(t *testing.T) {
	diffHunk := `@@ -10,7 +10,7 @@ func main() {
 	fmt.Println("Hello")
-	name := "old"
+	name := "new"
 	fmt.Printf("Name: %s\n", name)
 }`

	// Capture output
	stdout, _ := testutil.CaptureOutput(func() {
		displayDiffHunk(diffHunk)
	})

	assert.Contains(t, stdout, "🔹 @@ -10,7 +10,7 @@ func main() {")
	assert.Contains(t, stdout, "➖ \tname := \"old\"")
	assert.Contains(t, stdout, "➕ \tname := \"new\"")
	assert.Contains(t, stdout, "fmt.Println(\"Hello\")")
}

// Dependencies represents injectable dependencies for testing
type Dependencies struct {
	GitHubClient github.GitHubAPI
	Output       io.Writer
}

// NewListCmdWithDeps creates a list command with injectable dependencies
func NewListCmdWithDeps(deps *Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [pr]",
		Short: "List all comments on a PR",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListWithDeps(cmd, args, deps)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&showResolved, "resolved", false, "Include resolved comments")
	cmd.Flags().BoolVar(&onlyUnresolved, "unresolved", false, "Show only unresolved comments")
	cmd.Flags().StringVar(&author, "author", "", "Filter comments by author")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output without URLs and IDs")
	cmd.Flags().BoolVar(&hideAuthors, "hide-authors", false, "Hide author names for privacy")

	return cmd
}

// runListWithDeps is the testable version of runList with dependency injection
func runListWithDeps(cmd *cobra.Command, args []string, deps *Dependencies) error {
	var pr int
	var err error

	// Parse PR argument
	if len(args) == 1 {
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return formatValidationError("PR number", args[0], "must be a valid integer")
		}
	} else {
		// For testing, use a default PR number
		pr = 1
	}

	// For testing, use a mock repository
	_ = "owner/repo"

	// Fetch comments using injected client
	issueComments, err := deps.GitHubClient.ListIssueComments("owner", "repo", pr)
	if err != nil {
		return fmt.Errorf("failed to fetch issue comments: %w", err)
	}

	reviewComments, err := deps.GitHubClient.ListReviewComments("owner", "repo", pr)
	if err != nil {
		return fmt.Errorf("failed to fetch review comments: %w", err)
	}

	// Convert to internal Comment format
	var allComments []Comment
	for _, comment := range issueComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Type:      comment.Type,
		})
	}
	for _, comment := range reviewComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Path:      comment.Path,
			Line:      comment.Line,
			Type:      comment.Type,
		})
	}

	// Filter comments
	filteredComments := filterComments(allComments)

	// Display comments
	if len(filteredComments) == 0 {
		fmt.Fprintln(deps.Output, "No comments found")
		return nil
	}

	// Redirect output for testing
	oldStdout := os.Stdout
	if deps.Output != nil {
		if buf, ok := deps.Output.(*bytes.Buffer); ok {
			r, w, _ := os.Pipe()
			os.Stdout = w
			go func() {
				io.Copy(buf, r)
			}()
			displayComments(filteredComments, pr)
			w.Close()
			os.Stdout = oldStdout
		} else {
			displayComments(filteredComments, pr)
		}
	} else {
		displayComments(filteredComments, pr)
	}
	return nil
}

// resetListFlags resets global flags to default values for testing
func resetListFlags() {
	showResolved = false
	onlyUnresolved = false
	author = ""
	quiet = false
	hideAuthors = false
}
