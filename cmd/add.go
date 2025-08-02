package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	messages            []string
	noExpandSuggestions bool

	// Client for dependency injection (tests can override)
	addClient github.GitHubAPI
)

var addCmd = &cobra.Command{
	Use:   "add [pr] <file> <line> <comment>",
	Short: "Add a comment to a pull request",
	Long: heredoc.Doc(`
		Add a comment to a pull request. Comments can be general PR comments
		or line-specific review comments.

		For line-specific comments, use file and line arguments to target
		specific code locations. Supports both single-line and range comments.

		The comment message supports GitHub markdown formatting and can include
		code suggestions using the [SUGGEST: code] syntax.
	`),
	Example: heredoc.Doc(`
		# Strategic line-specific commenting
		$ gh comment add 123 src/api.js 42 "This rate limiting logic needs edge case handling for concurrent requests"
		$ gh comment add 123 auth.go 15:25 "Consider OAuth2 PKCE flow for mobile clients - current implementation has security gaps"

		# Security-focused reviews
		$ gh comment add 123 database.py 156 "This query is vulnerable to SQL injection - use parameterized queries"
		$ gh comment add 123 crypto.js 67 "[SUGGEST: use crypto.randomBytes(32) instead of Math.random()]"

		# Performance optimization suggestions
		$ gh comment add 123 performance.js 89:95 "Extract this expensive calculation to a cached service - it's called on every render"
		$ gh comment add 123 db/migrations.sql 23 "Add index on user_id column for faster lookups: CREATE INDEX idx_user_id ON orders(user_id)"

		# Architecture and design feedback
		$ gh comment add 123 service.go 134:150 "This business logic should be extracted to a domain service layer"
		$ gh comment add 123 component.tsx 45 "Consider using React.memo() to prevent unnecessary re-renders"

		# Multi-line strategic feedback
		$ gh comment add 123 error-handler.js 78 -m "**Critical:** This error handling is incomplete" -m "Missing: rate limit errors, network timeouts, auth failures"

		# Auto-detect PR from current branch
		$ gh comment add src/validation.js 156 "Add input sanitization before database operations"
	`),
	Args: cobra.RangeArgs(2, 4),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringArrayVarP(&messages, "message", "m", []string{}, "Add message (can be used multiple times for multi-line comments)")
	addCmd.Flags().BoolVar(&noExpandSuggestions, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax")
}

func runAdd(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if addClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to initialize GitHub client: %w", err)
		}
		addClient = client
	}

	var pr int
	var file, lineSpec, comment string
	var err error

	// Parse arguments - handle both 3 and 4 arg cases, plus --message flags
	if len(args) == 4 {
		// PR number provided
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return formatValidationError("PR number", args[0], "must be a valid integer")
		}
		file = args[1]
		lineSpec = args[2]
		comment = args[3]
	} else if len(args) == 3 {
		// PR number not provided, auto-detect
		pr, err = getCurrentPR()
		if err != nil {
			return err
		}
		file = args[0]
		lineSpec = args[1]
		comment = args[2]
	} else if len(args) == 2 && len(messages) > 0 {
		// Using --message flags instead of positional comment
		pr, err = getCurrentPR()
		if err != nil {
			return err
		}
		file = args[0]
		lineSpec = args[1]
		comment = strings.Join(messages, "\n")
	} else if len(args) == 3 && len(messages) > 0 {
		// PR provided + --message flags
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return formatValidationError("PR number", args[0], "must be a valid integer")
		}
		file = args[1]
		lineSpec = args[2]
		comment = strings.Join(messages, "\n")
	} else {
		return fmt.Errorf("invalid arguments. Use: gh comment add [pr] <file> <line> <comment> OR gh comment add [pr] <file> <line> --message \"line1\" --message \"line2\"")
	}

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

	// Parse line specification
	startLine, endLine, err := parseLineSpec(lineSpec)
	if err != nil {
		return err
	}

	// Expand suggestion syntax to GitHub markdown (unless disabled)
	var transformedComment string
	if noExpandSuggestions {
		transformedComment = comment
	} else {
		transformedComment = expandSuggestions(comment)
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR: %d\n", pr)
		fmt.Printf("File: %s\n", file)
		fmt.Printf("Line(s): %d", startLine)
		if endLine != startLine {
			fmt.Printf("-%d", endLine)
		}
		fmt.Printf("\nOriginal comment: %s\n", comment)
		fmt.Printf("Transformed comment: %s\n", transformedComment)
	}

	if dryRun {
		fmt.Printf("Would add comment to %s:%d", file, startLine)
		if endLine != startLine {
			fmt.Printf("-%d", endLine)
		}
		fmt.Printf(" in PR #%d:\n%s\n", pr, transformedComment)
		return nil
	}

	// Parse owner/repo
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repository)
	}
	owner, repoName := parts[0], parts[1]

	// Create review comment input
	reviewComment := github.ReviewCommentInput{
		Body: transformedComment,
		Path: file,
		Line: endLine, // GitHub API uses the end line for ranges
	}

	// If it's a range, add start_line
	if startLine != endLine {
		reviewComment.StartLine = startLine
		reviewComment.Side = "RIGHT"
	}

	// Note: GitHub automatically uses the latest commit SHA for review comments
	// No need to manually set commit_id

	// Add the comment via GitHub API
	err = addClient.AddReviewComment(owner, repoName, pr, reviewComment)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	// Success message
	fmt.Printf("✓ Added comment to %s:%d", file, startLine)
	if endLine != startLine {
		fmt.Printf("-%d", endLine)
	}
	fmt.Printf(" in PR #%d\n", pr)

	return nil
}

func parseLineSpec(lineSpec string) (int, int, error) {
	if strings.Contains(lineSpec, ":") {
		// Range specification
		parts := strings.Split(lineSpec, ":")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid line range format: %s (use start:end)", lineSpec)
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid start line: %s", parts[0])
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid end line: %s", parts[1])
		}

		if start <= 0 || end <= 0 {
			return 0, 0, fmt.Errorf("line numbers must be positive")
		}

		if start > end {
			return 0, 0, fmt.Errorf("start line (%d) cannot be greater than end line (%d)", start, end)
		}

		return start, end, nil
	} else {
		// Single line
		line, err := strconv.Atoi(lineSpec)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid line number: %s", lineSpec)
		}
		if line <= 0 {
			return 0, 0, fmt.Errorf("line numbers must be positive")
		}
		return line, line, nil
	}
}
