package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants for API limits and defaults
const (
	MaxGraphQLResults = 100
	MaxCommentLength  = 65536 // GitHub's actual limit for comment body
	MaxFilePathLength = 4096  // Reasonable file path limit
	MaxAuthorLength   = 39    // GitHub username max length
	MaxRepoNameLength = 100   // GitHub repository name max length
	MaxBranchLength   = 255   // Git branch name max length
	DefaultPageSize   = 30

	// Display constants
	MaxDisplayBodyLength   = 200 // Max length for comment body display
	TruncationSuffix       = "..."
	TruncationReserve      = 3    // Length of "..."
	SeparatorLength        = 50   // Length of separator lines
	MessageTruncateLength  = 50   // Length for message truncation in batch dry-run
	CommitSHADisplayLength = 8    // Number of characters to show from commit SHA
	DefaultBufferSize      = 4096 // Default buffer size for I/O operations
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
	// Check more specific patterns first before generic ones
	case containsAny(errStr, []string{"No subschema in oneOf matched", "invalid request"}):
		return fmt.Errorf("invalid request format during %s: %w\n\n💡 Suggestions:\n  • Check the command syntax in --help\n  • Verify all required arguments are provided\n  • For line comments, ensure the line exists in the PR diff\n  • Documentation: https://docs.github.com/en/rest/reference", operation, err)

	// Command-specific error patterns - check these before generic HTTP codes
	case containsAny(errStr, []string{"line not in diff", "line not part of", "comment line"}):
		return fmt.Errorf("line not commentable during %s: %w\n\n💡 Suggestions:\n  • Use 'gh comment lines <pr> <file>' to see which lines can have comments\n  • Only lines that were modified in the PR can have review comments\n  • For general comments, use 'gh comment add' without line numbers\n  • Documentation: https://docs.github.com/en/rest/pulls/comments#create-a-review-comment-for-a-pull-request", operation, err)

	case containsAny(errStr, []string{"file not found", "file not in PR", "path not found"}):
		return fmt.Errorf("file not found in PR during %s: %w\n\n💡 Suggestions:\n  • Verify the file path is correct and exists in the PR\n  • Use relative paths from repository root\n  • Check 'gh pr diff' to see modified files\n  • File paths are case-sensitive\n  • Documentation: https://docs.github.com/en/rest/pulls/files", operation, err)

	case containsAny(errStr, []string{"pending review", "review state", "already reviewing"}):
		return fmt.Errorf("review state conflict during %s: %w\n\n💡 Suggestions:\n  • You may have a pending review that needs to be submitted\n  • Use 'gh comment close-pending-review' to discard pending review\n  • Submit your pending review before creating a new one\n  • Documentation: https://docs.github.com/en/rest/pulls/reviews", operation, err)

	case containsAny(errStr, []string{"yaml", "json", "parse", "unmarshal", "decode"}):
		return fmt.Errorf("configuration parsing error during %s: %w\n\n💡 Suggestions:\n  • Check YAML/JSON syntax for errors\n  • Validate indentation (YAML is sensitive to spaces)\n  • Use a YAML/JSON validator to check your file\n  • Ensure all required fields are present\n  • Documentation: https://yaml.org/spec/", operation, err)

	// Generic HTTP error codes
	case containsAny(errStr, []string{"422", "Unprocessable Entity", "validation failed"}):
		return fmt.Errorf("validation error during %s: %w\n\n💡 Suggestions:\n  • Check if the line number exists in the PR diff\n  • Use 'gh comment lines <pr> <file>' to see commentable lines\n  • Verify the file path is correct in the PR\n  • For line-specific comments, ensure the line was modified in this PR\n  • Documentation: https://docs.github.com/en/rest/pulls/comments", operation, err)

	case containsAny(errStr, []string{"404", "Not Found"}):
		return fmt.Errorf("resource not found during %s: %w\n\n💡 Suggestions:\n  • Verify the PR number exists and is accessible\n  • Check if the comment ID is valid\n  • Ensure you have permission to access this repository\n  • Verify repository exists and PR number is correct\n  • Documentation: https://docs.github.com/en/rest/reference/pulls", operation, err)

	case containsAny(errStr, []string{"403", "Forbidden"}):
		return fmt.Errorf("permission denied during %s: %w\n\n💡 Suggestions:\n  • Check if you have write access to the repository\n  • Verify your GitHub authentication with 'gh auth status'\n  • You cannot approve your own PR or comment on private repos without access\n  • Check repository access permissions\n  • Documentation: https://docs.github.com/en/rest/overview/permissions-required-for-github-apps", operation, err)

	case containsAny(errStr, []string{"401", "Unauthorized"}):
		return fmt.Errorf("authentication failed during %s: %w\n\n💡 Suggestions:\n  • Run 'gh auth login' to authenticate with GitHub\n  • Check if your token has expired with 'gh auth status'\n  • Verify you're authenticated with the correct GitHub account\n  • Run 'gh auth status' to verify login\n  • Documentation: https://cli.github.com/manual/gh_auth_login", operation, err)

	case containsAny(errStr, []string{"rate limit", "rate_limit", "too many requests", "429"}):
		return fmt.Errorf("rate limit exceeded during %s: %w\n\n💡 Suggestions:\n  • Wait a few minutes before trying again\n  • Use authenticated requests (ensure 'gh auth status' shows logged in)\n  • Consider reducing the frequency of API calls\n  • Use smaller batch sizes for bulk operations\n  • Check rate limit status: 'gh api rate_limit'\n  • Documentation: https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting", operation, err)

	case containsAny(errStr, []string{"500", "502", "503", "Internal Server Error", "Bad Gateway", "Service Unavailable"}):
		return fmt.Errorf("GitHub server error during %s: %w\n\n💡 Suggestions:\n  • This is a temporary GitHub server issue\n  • Try again in a few minutes\n  • Check GitHub's status page at https://status.github.com\n  • Documentation: https://www.githubstatus.com", operation, err)

	case containsAny(errStr, []string{"network", "timeout", "connection", "EOF", "broken pipe"}):
		return fmt.Errorf("network error during %s: %w\n\n💡 Suggestions:\n  • Check your internet connection\n  • Try again in a moment\n  • Verify GitHub is accessible from your network\n  • Try with --verbose flag for more details\n  • Check proxy settings if behind corporate firewall\n  • Documentation: https://docs.github.com/en/rest/overview/troubleshooting", operation, err)
	}

	// Default case - provide general guidance
	return fmt.Errorf("error during %s: %w\n\n💡 Suggestions:\n  • Check 'gh comment --help' for correct usage\n  • Verify PR number and file paths are correct\n  • Run with --verbose for more details\n  • Documentation: https://github.com/your-org/gh-comment#readme", operation, err)
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

// parsePositiveInt parses a string to a positive integer with consistent validation
func parsePositiveInt(s, fieldName string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, formatValidationError(fieldName, s, "must be a valid integer")
	}
	if val <= 0 {
		return 0, formatValidationError(fieldName, s, "must be a positive integer")
	}
	return val, nil
}

// validateCommentBody validates comment body length and content
func validateCommentBody(body string) error {
	if len(body) > MaxCommentLength {
		return fmt.Errorf("comment too long: %d characters (maximum %d allowed)", len(body), MaxCommentLength)
	}
	return nil
}

// validateFilePath validates file path length and format to prevent directory traversal
func validateFilePath(path string) error {
	if len(path) > MaxFilePathLength {
		return fmt.Errorf("file path too long: %d characters (maximum %d allowed)", len(path), MaxFilePathLength)
	}

	// Check for directory traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid file path: directory traversal not allowed")
	}

	// Check for absolute paths (should be relative to repo root)
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid file path: absolute paths not allowed, use relative paths from repository root")
	}

	return nil
}

// validateRepositoryName validates GitHub repository name format
func validateRepositoryName(repo string) error {
	if len(repo) > MaxRepoNameLength {
		return fmt.Errorf("repository name too long: %d characters (maximum %d allowed)", len(repo), MaxRepoNameLength)
	}

	// Check for basic owner/repo format
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: must be 'owner/repo'")
	}

	owner, repoName := parts[0], parts[1]
	if len(owner) == 0 || len(repoName) == 0 {
		return fmt.Errorf("invalid repository format: owner and repository name cannot be empty")
	}

	if len(owner) > MaxAuthorLength {
		return fmt.Errorf("repository owner too long: %d characters (maximum %d allowed)", len(owner), MaxAuthorLength)
	}

	return nil
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
