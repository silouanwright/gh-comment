package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var (
	messages            []string
	noExpandSuggestions bool
)

var addCmd = &cobra.Command{
	Use:   "add [pr] <file> <line> <comment>",
	Short: "Add a single line comment to a PR",
	Long: `Add a targeted comment to a specific line in a pull request.

The line can be specified as a single line number or a range (start:end).
Supports both inline comments and multi-line comments using --message flags.

Comments are posted immediately to the PR.

Examples:
  # Add single-line comment (posts immediately)
  gh comment add 123 src/api.js 42 "this handles the rate limiting edge case"

  # Add range comment
  gh comment add 123 src/api.js 42:45 "this entire block needs review"

  # Add multi-line comment using --message flags (AI-friendly)
  gh comment add 123 src/api.js 42 --message "First paragraph" --message "Second paragraph"

  # Auto-detect PR with --message flags
  gh comment add src/api.js 42 -m "Line 1" -m "Line 2"`,
	Args: cobra.RangeArgs(2, 4),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringArrayVarP(&messages, "message", "m", []string{}, "Add message (can be used multiple times for multi-line comments)")
	addCmd.Flags().BoolVar(&noExpandSuggestions, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax")
}

func runAdd(cmd *cobra.Command, args []string) error {
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

	// Add the comment via GitHub API
	return addLineComment(repository, pr, file, startLine, endLine, transformedComment)
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
		return line, line, nil
	}
}

func addLineComment(repo string, pr int, file string, startLine, endLine int, comment string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// First, get the PR to find the commit SHA
	prData := struct {
		Head struct {
			SHA string `json:"sha"`
		} `json:"head"`
	}{}

	err = client.Get(fmt.Sprintf("repos/%s/pulls/%d", repo, pr), &prData)
	if err != nil {
		return fmt.Errorf("failed to get PR data: %w", err)
	}

	// Create the comment payload
	payload := map[string]interface{}{
		"body":      comment,
		"commit_id": prData.Head.SHA,
		"path":      file,
		"line":      endLine, // GitHub API uses the end line for ranges
	}

	// If it's a range, add start_line
	if startLine != endLine {
		payload["start_line"] = startLine
		payload["start_side"] = "RIGHT"
	}

	if verbose {
		payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Printf("API payload:\n%s\n", payloadJSON)
	}

	// Marshal payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make the immediate API call
	var response map[string]interface{}
	err = client.Post(fmt.Sprintf("repos/%s/pulls/%d/comments", repo, pr), bytes.NewReader(payloadJSON), &response)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	fmt.Printf("✓ Added comment to %s:%d", file, startLine)
	if endLine != startLine {
		fmt.Printf("-%d", endLine)
	}
	fmt.Printf(" in PR #%d\n", pr)

	if verbose {
		if htmlURL, ok := response["html_url"].(string); ok {
			fmt.Printf("Comment URL: %s\n", htmlURL)
		}
	}

	return nil
}
