package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFormatTimeAgoList tests the time formatting function for list command
func TestFormatTimeAgoList(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "just now",
			input:    now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			input:    now.Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "5 minutes ago",
			input:    now.Add(-5 * time.Minute),
			expected: "5 minutes ago",
		},
		{
			name:     "1 hour ago",
			input:    now.Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "3 hours ago",
			input:    now.Add(-3 * time.Hour),
			expected: "3 hours ago",
		},
		{
			name:     "1 day ago",
			input:    now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "3 days ago",
			input:    now.Add(-3 * 24 * time.Hour),
			expected: "3 days ago",
		},
		{
			name:     "old date",
			input:    now.Add(-365 * 24 * time.Hour),
			expected: now.Add(-365 * 24 * time.Hour).Format("Jan 2, 2006"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTimeAgo(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// TestFilterCommentsList tests comment filtering logic for list command
func TestFilterCommentsList(t *testing.T) {
	// Create test comments using the actual Comment struct
	comments := []Comment{
		{
			ID:        1,
			Author:    "alice",
			Body:      "First comment",
			Type:      "issue",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        2,
			Author:    "bob",
			Body:      "Second comment",
			Type:      "review",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        3,
			Author:    "alice",
			Body:      "Third comment",
			Type:      "review",
			CreatedAt: time.Now().Add(-3 * time.Hour),
		},
	}

	tests := []struct {
		name          string
		authorFilter  string
		expectedCount int
		expectedIDs   []int
	}{
		{
			name:          "no filter",
			authorFilter:  "",
			expectedCount: 3,
			expectedIDs:   []int{1, 2, 3},
		},
		{
			name:          "filter by alice",
			authorFilter:  "alice",
			expectedCount: 2,
			expectedIDs:   []int{1, 3},
		},
		{
			name:          "filter by bob",
			authorFilter:  "bob",
			expectedCount: 1,
			expectedIDs:   []int{2},
		},
		{
			name:          "filter by non-existent author",
			authorFilter:  "charlie",
			expectedCount: 0,
			expectedIDs:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the global author filter
			originalAuthor := author
			author = tt.authorFilter
			defer func() { author = originalAuthor }()

			filtered := filterComments(comments)

			assert.Len(t, filtered, tt.expectedCount)

			if tt.expectedCount > 0 {
				var actualIDs []int
				for _, comment := range filtered {
					actualIDs = append(actualIDs, comment.ID)
				}
				assert.Equal(t, tt.expectedIDs, actualIDs)
			}
		})
	}
}
