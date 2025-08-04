package cmd

import (
	"fmt"
	"strings"
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

// formatActionableError creates user-friendly error messages with actionable suggestions
func formatActionableError(operation string, err error) error {
	errStr := err.Error()

	// Handle common GitHub API error patterns
	switch {
	case containsAny(errStr, []string{"422", "Unprocessable Entity", "validation failed"}):
		return fmt.Errorf("validation error during %s: %w\n\n💡 Suggestions:\n  • Check if the line number exists in the PR diff\n  • Use 'gh comment lines <pr> <file>' to see commentable lines\n  • Verify the file path is correct in the PR", operation, err)

	case containsAny(errStr, []string{"404", "Not Found"}):
		return fmt.Errorf("resource not found during %s: %w\n\n💡 Suggestions:\n  • Verify the PR number exists and is accessible\n  • Check if the comment ID is valid\n  • Ensure you have permission to access this repository", operation, err)

	case containsAny(errStr, []string{"403", "Forbidden"}):
		return fmt.Errorf("permission denied during %s: %w\n\n💡 Suggestions:\n  • Check if you have write access to the repository\n  • Verify your GitHub authentication with 'gh auth status'\n  • You cannot approve your own PR or comment on private repos without access", operation, err)

	case containsAny(errStr, []string{"401", "Unauthorized"}):
		return fmt.Errorf("authentication failed during %s: %w\n\n💡 Suggestions:\n  • Run 'gh auth login' to authenticate with GitHub\n  • Check if your token has expired\n  • Verify you're authenticated with the correct GitHub account", operation, err)

	case containsAny(errStr, []string{"rate limit", "rate_limit", "too many requests"}):
		return fmt.Errorf("rate limit exceeded during %s: %w\n\n💡 Suggestions:\n  • Wait a few minutes before trying again\n  • Use authenticated requests (ensure 'gh auth status' shows logged in)\n  • Consider reducing the frequency of API calls", operation, err)

	case containsAny(errStr, []string{"500", "502", "503", "Internal Server Error", "Bad Gateway", "Service Unavailable"}):
		return fmt.Errorf("GitHub server error during %s: %w\n\n💡 Suggestions:\n  • This is a temporary GitHub server issue\n  • Try again in a few minutes\n  • Check GitHub's status page at https://status.github.com", operation, err)

	case containsAny(errStr, []string{"network", "timeout", "connection"}):
		return fmt.Errorf("network error during %s: %w\n\n💡 Suggestions:\n  • Check your internet connection\n  • Try again in a moment\n  • Verify GitHub is accessible from your network", operation, err)

	case containsAny(errStr, []string{"No subschema in oneOf matched", "invalid request"}):
		return fmt.Errorf("invalid request format during %s: %w\n\n💡 Suggestions:\n  • Check the command syntax in --help\n  • Verify all required arguments are provided\n  • For line comments, ensure the line exists in the PR diff", operation, err)
	}

	// Default case - provide general guidance
	return fmt.Errorf("error during %s: %w\n\n💡 Suggestions:\n  • Check 'gh comment --help' for correct usage\n  • Verify PR number and file paths are correct\n  • Run with --verbose for more details", operation, err)
}

// containsAny checks if a string contains any of the provided substrings (case-insensitive)
func containsAny(str string, substrings []string) bool {
	lowerStr := strings.ToLower(str)
	for _, substr := range substrings {
		if strings.Contains(lowerStr, strings.ToLower(substr)) {
			return true
		}
	}
	return false
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
