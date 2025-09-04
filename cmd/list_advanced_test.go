package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/silouanwright/gh-comment/internal/github"
)

func TestAdvancedFiltering(t *testing.T) {
	// Save original client
	originalClient := listClient
	defer func() { listClient = originalClient }()

	// Reset global variables before each test
	resetListFlags := func() {
		author = ""
		filter = ""
		showRecent = false
		since = ""
		until = ""
		listType = ""
		sinceTime = nil
		untilTime = nil
		quiet = false
		hideAuthors = false
		outputFormat = "default"
		idsOnly = false
	}

	t.Run("Date Filter Parsing", func(t *testing.T) {
		tests := []struct {
			name        string
			dateStr     string
			expectError bool
		}{
			{"Valid YYYY-MM-DD", "2024-01-15", false},
			{"Valid relative time", "3 days ago", false},
			{"Valid relative time plural", "2 weeks ago", false},
			{"Valid MM/DD/YYYY", "01/15/2024", false},
			{"Valid ISO 8601", "2024-01-15T10:30:00Z", false},
			{"Invalid format", "not-a-date", true},
			{"Invalid relative", "invalid ago", true},
			{"Invalid relative unit", "3 fortnights ago", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := parseFlexibleDate(tt.dateStr)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Relative Time Parsing with dateparse", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expectError bool
		}{
			{"3 days ago", "3 days ago", false},
			{"1 week ago", "1 week ago", false},
			{"2 months ago", "2 months ago", false},
			{"1 year ago", "1 year ago", false},
			{"5 hours ago", "5 hours ago", false},
			{"invalid relative", "invalid ago", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := parseFlexibleDate(tt.input)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Author Filter Matching", func(t *testing.T) {
		tests := []struct {
			name     string
			author   string
			filter   string
			expected bool
		}{
			{"Exact match", "octocat", "octocat", true},
			{"Case insensitive", "OctoCat", "octocat", true},
			{"Partial match", "octocat-dev", "octo", true},
			{"Wildcard prefix", "octocat", "octo*", true},
			{"Wildcard suffix", "octocat", "*cat", true},
			{"Wildcard middle", "octocat-dev", "octo*dev", true},
			{"Email filter", "user@company.com", "*@company.com", true},
			{"No match", "different", "octocat", false},
			{"Wildcard no match", "octocat", "dog*", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := matchesAuthorFilter(tt.author, tt.filter)
				assert.Equal(t, tt.expected, result, "author: %s, filter: %s", tt.author, tt.filter)
			})
		}
	})

	t.Run("Filter Validation", func(t *testing.T) {
		resetListFlags()

		tests := []struct {
			name        string
			setupFunc   func()
			expectError bool
			errorMsg    string
		}{
			{
				name: "Valid filter",
				setupFunc: func() {
					filter = "today"
				},
				expectError: false,
			},
			{
				name: "Invalid filter",
				setupFunc: func() {
					filter = "invalid"
				},
				expectError: true,
				errorMsg:    "invalid filter",
			},
			{
				name: "Valid comment type",
				setupFunc: func() {
					listType = "review"
				},
				expectError: false,
			},
			{
				name: "Invalid comment type",
				setupFunc: func() {
					listType = "invalid"
				},
				expectError: true,
				errorMsg:    "invalid type",
			},
			{
				name: "Valid date range",
				setupFunc: func() {
					since = "2024-01-01"
					until = "2024-12-31"
				},
				expectError: false,
			},
			{
				name: "Invalid date range",
				setupFunc: func() {
					since = "2024-12-31"
					until = "2024-01-01"
				},
				expectError: true,
				errorMsg:    "since date",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resetListFlags()
				tt.setupFunc()

				err := validateAndParseFilters()
				if tt.expectError {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.errorMsg)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Comment Filtering", func(t *testing.T) {
		// Create test comments with different properties
		baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
		testComments := []Comment{
			{
				ID:        1,
				Author:    "octocat",
				Body:      "First comment",
				CreatedAt: baseTime.AddDate(0, 0, -5), // 5 days ago
				Type:      "issue",
			},
			{
				ID:        2,
				Author:    "developer",
				Body:      "Review comment",
				CreatedAt: baseTime.AddDate(0, 0, -2), // 2 days ago
				Type:      "review",
				Path:      "src/main.go",
				Line:      42,
			},
			{
				ID:        3,
				Author:    "user@company.com",
				Body:      "Another comment",
				CreatedAt: baseTime.AddDate(0, 0, -1), // 1 day ago
				Type:      "issue",
			},
		}

		tests := []struct {
			name          string
			setupFunc     func()
			expectedCount int
			expectedIDs   []int
		}{
			{
				name: "No filters - all comments",
				setupFunc: func() {
					resetListFlags()
				},
				expectedCount: 3,
				expectedIDs:   []int{1, 2, 3},
			},
			{
				name: "Filter by author exact",
				setupFunc: func() {
					resetListFlags()
					author = "octocat"
				},
				expectedCount: 1,
				expectedIDs:   []int{1},
			},
			{
				name: "Filter by author wildcard",
				setupFunc: func() {
					resetListFlags()
					author = "*@company.com"
				},
				expectedCount: 1,
				expectedIDs:   []int{3},
			},
			{
				name: "Filter by comment type",
				setupFunc: func() {
					resetListFlags()
					listType = "review"
				},
				expectedCount: 1,
				expectedIDs:   []int{2},
			},
			{
				name: "Filter by since date",
				setupFunc: func() {
					resetListFlags()
					since = baseTime.AddDate(0, 0, -3).Format("2006-01-02") // 3 days ago
					validateAndParseFilters()                               // Parse the date
				},
				expectedCount: 2,
				expectedIDs:   []int{2, 3}, // Comments from 2 days ago and 1 day ago
			},
			{
				name: "Filter by until date",
				setupFunc: func() {
					resetListFlags()
					until = baseTime.AddDate(0, 0, -3).Format("2006-01-02") // 3 days ago
					validateAndParseFilters()                               // Parse the date
				},
				expectedCount: 1,
				expectedIDs:   []int{1}, // Comment from 5 days ago
			},
			{
				name: "Combined filters",
				setupFunc: func() {
					resetListFlags()
					author = "developer"
					listType = "review"
				},
				expectedCount: 1,
				expectedIDs:   []int{2},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setupFunc()

				filtered := filterComments(testComments)
				assert.Equal(t, tt.expectedCount, len(filtered), "Expected %d comments, got %d", tt.expectedCount, len(filtered))

				// Check that the right comments were returned
				for i, expectedID := range tt.expectedIDs {
					if i < len(filtered) {
						assert.Equal(t, expectedID, filtered[i].ID, "Expected comment ID %d at position %d", expectedID, i)
					}
				}
			})
		}
	})

	t.Run("Integration Test - List Command with Filters", func(t *testing.T) {
		resetListFlags()

		// Set up mock client
		mockClient := github.NewMockClient()
		listClient = mockClient

		// Configure mock responses - set the mock data directly
		mockClient.IssueComments = []github.Comment{
			{
				ID:        123,
				User:      github.User{Login: "octocat"},
				Body:      "General comment",
				CreatedAt: time.Now().AddDate(0, 0, -1),
				Type:      "issue",
			},
		}

		mockClient.ReviewComments = []github.Comment{
			{
				ID:        456,
				User:      github.User{Login: "developer"},
				Body:      "Line comment",
				CreatedAt: time.Now().AddDate(0, 0, -2),
				Path:      "main.go",
				Line:      42,
				Type:      "review",
			},
		}

		// Test with author filter
		author = "octocat"

		// This would normally run the full command, but we're testing the filtering logic
		comments, err := fetchAllComments(mockClient, "owner/repo", 123)
		require.NoError(t, err)
		assert.Equal(t, 2, len(comments), "Should fetch both issue and review comments")

		filtered := filterComments(comments)
		assert.Equal(t, 1, len(filtered), "Should filter to only octocat's comments")
		assert.Equal(t, "octocat", filtered[0].Author)
	})
}

func TestContainsHelper(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	assert.True(t, containsString(slice, "banana"))
	assert.False(t, containsString(slice, "orange"))
	assert.False(t, containsString([]string{}, "anything"))
}
