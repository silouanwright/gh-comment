package cmd

import (
	"testing"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGitHubMockClient tests the mock client functionality
func TestGitHubMockClient(t *testing.T) {
	mockClient := github.NewMockClient()

	t.Run("list issue comments", func(t *testing.T) {
		comments, err := mockClient.ListIssueComments("owner", "repo", 1)
		require.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, "issue", comments[0].Type)
		assert.Equal(t, "LGTM! Great work on this PR.", comments[0].Body)
	})

	t.Run("list review comments", func(t *testing.T) {
		comments, err := mockClient.ListReviewComments("owner", "repo", 1)
		require.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, "review", comments[0].Type)
		assert.Equal(t, "main.go", comments[0].Path)
		assert.Equal(t, 42, comments[0].Line)
	})

	t.Run("create issue comment", func(t *testing.T) {
		comment, err := mockClient.CreateIssueComment("owner", "repo", 1, "Test comment")
		require.NoError(t, err)
		assert.Equal(t, "Test comment", comment.Body)
		assert.Equal(t, "issue", comment.Type)
		assert.Equal(t, "testuser", comment.User.Login)
	})

	t.Run("create review comment reply", func(t *testing.T) {
		comment, err := mockClient.CreateReviewCommentReply("owner", "repo", 123, "Test reply")
		require.NoError(t, err)
		assert.Equal(t, "Test reply", comment.Body)
		assert.Equal(t, "review", comment.Type)
	})

	t.Run("resolve review thread", func(t *testing.T) {
		err := mockClient.ResolveReviewThread("RT_123")
		require.NoError(t, err)
		assert.Equal(t, "RT_123", mockClient.ResolvedThread)
	})
}

// TestCommentFiltering tests the comment filtering logic
func TestCommentFiltering(t *testing.T) {
	comments := []Comment{
		{
			ID:     1,
			Author: "alice",
			Body:   "First comment",
			Type:   "issue",
		},
		{
			ID:     2,
			Author: "bob",
			Body:   "Second comment",
			Type:   "review",
		},
		{
			ID:     3,
			Author: "alice",
			Body:   "Third comment",
			Type:   "issue",
		},
	}

	tests := []struct {
		name           string
		authorFilter   string
		expectedCount  int
		expectedAuthor string
	}{
		{
			name:          "no filter",
			authorFilter:  "",
			expectedCount: 3,
		},
		{
			name:           "filter by alice",
			authorFilter:   "alice",
			expectedCount:  2,
			expectedAuthor: "alice",
		},
		{
			name:           "filter by bob",
			authorFilter:   "bob",
			expectedCount:  1,
			expectedAuthor: "bob",
		},
		{
			name:          "filter by nonexistent",
			authorFilter:  "charlie",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global author filter
			originalAuthor := author
			author = tt.authorFilter
			defer func() { author = originalAuthor }()

			filtered := filterComments(comments)
			assert.Len(t, filtered, tt.expectedCount)

			if tt.expectedAuthor != "" {
				for _, comment := range filtered {
					assert.Equal(t, tt.expectedAuthor, comment.Author)
				}
			}
		})
	}
}

// TestTimeFormatting tests the time formatting function
func TestTimeFormatting(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeAgo(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestValidationHelpers tests validation functions
func TestValidationHelpers(t *testing.T) {
	t.Run("validate reaction", func(t *testing.T) {
		validReactions := []string{"+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", "eyes"}
		invalidReactions := []string{"", "invalid", "thumbsup", "smile"}

		for _, reaction := range validReactions {
			assert.True(t, validateReaction(reaction), "Expected %s to be valid", reaction)
		}

		for _, reaction := range invalidReactions {
			assert.False(t, validateReaction(reaction), "Expected %s to be invalid", reaction)
		}
	})
}

// TestErrorFormatting tests error formatting helpers
func TestErrorFormatting(t *testing.T) {
	t.Run("format validation error", func(t *testing.T) {
		err := formatValidationError("comment ID", "invalid", "must be a valid integer")
		expected := "invalid comment ID 'invalid': must be a valid integer"
		assert.Equal(t, expected, err.Error())
	})
}

// TestSuggestionExpansion tests suggestion syntax expansion
func TestSuggestionExpansion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no suggestion",
			input:    "Just a regular comment",
			expected: "Just a regular comment",
		},
		{
			name:     "simple suggestion",
			input:    "```suggestion\nconst name = 'new';\n```",
			expected: "```suggestion\nconst name = 'new';\n```",
		},
		{
			name:     "suggestion with text",
			input:    "Consider this:\n```suggestion\nconst name = 'better';\n```\nThis is clearer.",
			expected: "Consider this:\n```suggestion\nconst name = 'better';\n```\nThis is clearer.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandSuggestions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
