package cmd

import (
	"fmt"
)

// Constants for API limits and defaults
const (
	MaxGraphQLResults = 100
	MaxCommentLength  = 65536
	DefaultPageSize   = 30
)

// getPRContext gets the repository and PR number, handling both flag and auto-detection
func getPRContext() (repo string, pr int, err error) {
	repo, err = getCurrentRepo()
	if err != nil {
		return "", 0, fmt.Errorf("failed to get repository: %w", err)
	}

	if prNumber > 0 {
		pr = prNumber
	} else {
		pr, err = getCurrentPR()
		if err != nil {
			return "", 0, fmt.Errorf("failed to detect PR number: %w (try specifying --pr)", err)
		}
	}

	return repo, pr, nil
}

// formatAPIError creates consistent error messages for API failures
func formatAPIError(operation, endpoint string, err error) error {
	return fmt.Errorf("GitHub API error during %s: %w", operation, err)
}

// formatValidationError creates consistent error messages for validation failures
func formatValidationError(field, value, expected string) error {
	return fmt.Errorf("invalid %s '%s': %s", field, value, expected)
}

// formatNotFoundError creates consistent error messages for missing resources
func formatNotFoundError(resource string, identifier interface{}) error {
	return fmt.Errorf("%s not found: %v", resource, identifier)
}

// lineRange represents a range of consecutive line numbers for display
type lineRange struct {
	start, end int
}

// groupConsecutiveLines groups consecutive line numbers into ranges for better display
func groupConsecutiveLines(lines []int) []lineRange {
	if len(lines) == 0 {
		return nil
	}

	var ranges []lineRange
	start := lines[0]
	end := lines[0]

	for i := 1; i < len(lines); i++ {
		if lines[i] == end+1 {
			// Consecutive line, extend the range
			end = lines[i]
		} else {
			// Gap found, close current range and start new one
			ranges = append(ranges, lineRange{start: start, end: end})
			start = lines[i]
			end = lines[i]
		}
	}

	// Add the final range
	ranges = append(ranges, lineRange{start: start, end: end})

	return ranges
}
