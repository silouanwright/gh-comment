package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateReaction(t *testing.T) {
	tests := []struct {
		reaction string
		valid    bool
	}{
		{"+1", true},
		{"-1", true},
		{"laugh", true},
		{"confused", true},
		{"heart", true},
		{"hooray", true},
		{"rocket", true},
		{"eyes", true},
		{"invalid", false},
		{"thumbsup", false},   // This is different from +1
		{"thumbsdown", false}, // This is different from -1
		{"", false},
		{"LAUGH", false}, // Case sensitive
		{"Heart", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.reaction, func(t *testing.T) {
			result := validateReaction(tt.reaction)
			assert.Equal(t, tt.valid, result, "validateReaction(%q) should be %v", tt.reaction, tt.valid)
		})
	}
}

// These functions are already tested in helpers_test.go, so we skip them here

func TestParseFlexibleDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		checkFunc func(*testing.T, time.Time) // Custom validation function
	}{
		{
			name:    "YYYY-MM-DD format",
			input:   "2024-01-15",
			wantErr: false,
			checkFunc: func(t *testing.T, result time.Time) {
				assert.Equal(t, 2024, result.Year())
				assert.Equal(t, time.January, result.Month())
				assert.Equal(t, 15, result.Day())
			},
		},
		{
			name:    "relative time - days ago",
			input:   "3 days ago",
			wantErr: false,
			checkFunc: func(t *testing.T, result time.Time) {
				// Should be approximately 3 days ago
				expected := now.AddDate(0, 0, -3)
				assert.WithinDuration(t, expected, result, 24*time.Hour)
			},
		},
		{
			name:    "relative time - weeks ago",
			input:   "2 weeks ago",
			wantErr: false,
			checkFunc: func(t *testing.T, result time.Time) {
				// Should be approximately 2 weeks ago
				expected := now.AddDate(0, 0, -14)
				assert.WithinDuration(t, expected, result, 24*time.Hour)
			},
		},
		{
			name:    "relative time - months ago",
			input:   "1 month ago",
			wantErr: false,
			checkFunc: func(t *testing.T, result time.Time) {
				// Should be approximately 1 month ago
				expected := now.AddDate(0, -1, 0)
				assert.WithinDuration(t, expected, result, 48*time.Hour) // Allow more tolerance for months
			},
		},
		{
			name:    "ISO 8601 format",
			input:   "2024-01-15T10:30:00Z",
			wantErr: false,
			checkFunc: func(t *testing.T, result time.Time) {
				assert.Equal(t, 2024, result.Year())
				assert.Equal(t, time.January, result.Month())
				assert.Equal(t, 15, result.Day())
				assert.Equal(t, 10, result.Hour())
				assert.Equal(t, 30, result.Minute())
			},
		},
		{
			name:      "invalid format",
			input:     "not-a-date",
			wantErr:   true,
			checkFunc: nil,
		},
		{
			name:    "yesterday format",
			input:   "yesterday",
			wantErr: false,
			checkFunc: func(t *testing.T, result time.Time) {
				// Should be approximately yesterday
				expected := now.AddDate(0, 0, -1)
				assert.WithinDuration(t, expected, result, 24*time.Hour)
			},
		},
		{
			name:      "empty string",
			input:     "",
			wantErr:   true,
			checkFunc: nil,
		},
		{
			name:      "invalid relative number",
			input:     "abc days ago",
			wantErr:   true,
			checkFunc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFlexibleDate(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.checkFunc != nil {
					tt.checkFunc(t, result)
				}
			}
		})
	}
}

func TestMatchesAuthorFilter(t *testing.T) {
	tests := []struct {
		name     string
		author   string
		filter   string
		expected bool
	}{
		{
			name:     "exact match",
			author:   "john",
			filter:   "john",
			expected: true,
		},
		{
			name:     "no match",
			author:   "john",
			filter:   "jane",
			expected: false,
		},
		{
			name:     "wildcard prefix",
			author:   "john-doe",
			filter:   "john*",
			expected: true,
		},
		{
			name:     "wildcard suffix",
			author:   "john-doe",
			filter:   "*doe",
			expected: true,
		},
		{
			name:     "wildcard middle",
			author:   "john-doe-smith",
			filter:   "john*smith",
			expected: true,
		},
		{
			name:     "case insensitive - should match",
			author:   "John",
			filter:   "john",
			expected: true, // Function is case-insensitive
		},
		{
			name:     "wildcard no match",
			author:   "alice",
			filter:   "john*",
			expected: false,
		},
		{
			name:     "empty filter matches all",
			author:   "anyone",
			filter:   "",
			expected: true,
		},
		{
			name:     "empty author with filter",
			author:   "",
			filter:   "john",
			expected: false,
		},
		{
			name:     "both empty",
			author:   "",
			filter:   "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesAuthorFilter(tt.author, tt.filter)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item found in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "item not found in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "orange",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "apple",
			expected: false,
		},
		{
			name:     "empty item in slice",
			slice:    []string{"", "banana", "cherry"},
			item:     "",
			expected: true,
		},
		{
			name:     "empty item not in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "",
			expected: false,
		},
		{
			name:     "case sensitive match",
			slice:    []string{"Apple", "banana", "Cherry"},
			item:     "apple",
			expected: false,
		},
		{
			name:     "exact case match",
			slice:    []string{"Apple", "banana", "Cherry"},
			item:     "Apple",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsString(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}
