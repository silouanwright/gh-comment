package cmd

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	// Client for dependency injection (tests can override)
	batchClient github.GitHubAPI
)

// BatchConfig represents the structure of a batch comment configuration file
type BatchConfig struct {
	PR       int             `yaml:"pr,omitempty"`
	Repo     string          `yaml:"repo,omitempty"`
	Review   *ReviewConfig   `yaml:"review,omitempty"`
	Comments []CommentConfig `yaml:"comments,omitempty"`
}

// ReviewConfig represents review-level configuration
type ReviewConfig struct {
	Body  string `yaml:"body,omitempty"`
	Event string `yaml:"event,omitempty"` // APPROVE, REQUEST_CHANGES, COMMENT
}

// CommentConfig represents individual comment configuration
type CommentConfig struct {
	File    string `yaml:"file"`
	Line    int    `yaml:"line,omitempty"`
	Range   string `yaml:"range,omitempty"` // e.g., "10-15"
	Message string `yaml:"message"`
	Type    string `yaml:"type,omitempty"` // "review" or "issue", defaults to "review"
}

var batchCmd = &cobra.Command{
	Use:   "batch <config-file>",
	Short: "Process multiple comments from a YAML configuration file",
	Long: heredoc.Doc(`
		Process multiple comments, reactions, and reviews from a YAML configuration file.

		This is ideal for bulk operations, automated workflows, or complex review
		scenarios. The config file can specify mixed comment types, create reviews
		with multiple comments, and set up entire review workflows.
	`),
	Example: heredoc.Doc(`
		# Process comments from config
		$ gh comment batch 123 review-config.yaml

		# Validate config without executing
		$ gh comment batch 123 review-config.yaml --dry-run

		# Use verbose output
		$ gh comment batch 123 review-config.yaml --verbose
	`),
	Args: cobra.ExactArgs(2),
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)
}

func runBatch(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if batchClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		batchClient = client
	}

	// Parse PR number
	prArg := args[0]
	pr, err := strconv.Atoi(prArg)
	if err != nil {
		return formatValidationError("PR number", prArg, "must be a valid integer")
	}

	// Read and parse configuration file
	configFile := args[1]
	config, err := readBatchConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Override PR and repo from config if specified
	if config.PR != 0 {
		pr = config.PR
	}

	repository := repo
	if config.Repo != "" {
		repository = config.Repo
	}

	// Get repository and PR context if not specified
	if repository == "" {
		repository, pr, err = getPRContext()
		if err != nil {
			return err
		}
	}

	// Parse owner/repo
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repository)
	}
	owner, repoName := parts[0], parts[1]

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR: %d\n", pr)
		fmt.Printf("Config file: %s\n", configFile)
		fmt.Printf("Comments to process: %d\n", len(config.Comments))
		if config.Review != nil {
			fmt.Printf("Review event: %s\n", config.Review.Event)
		}
		fmt.Println()
	}

	if dryRun {
		fmt.Printf("Would process %d comments from %s on PR #%d:\n", len(config.Comments), configFile, pr)
		for i, comment := range config.Comments {
			fmt.Printf("  %d. %s:%s - %s\n", i+1, comment.File, formatLineOrRange(comment), truncateMessage(comment.Message, 50))
		}
		if config.Review != nil {
			fmt.Printf("Would create review with event: %s\n", config.Review.Event)
		}
		return nil
	}

	// Process comments
	return processBatchComments(batchClient, owner, repoName, pr, config)
}

func readBatchConfig(configFile string) (*BatchConfig, error) {
	// Read file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", configFile, err)
	}

	// Parse YAML
	var config BatchConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate configuration
	if len(config.Comments) == 0 && config.Review == nil {
		return nil, fmt.Errorf("configuration must contain either comments or review")
	}

	// Validate comments
	for i, comment := range config.Comments {
		if comment.File == "" {
			return nil, fmt.Errorf("comment %d: file is required", i+1)
		}
		if comment.Message == "" {
			return nil, fmt.Errorf("comment %d: message is required", i+1)
		}
		if comment.Line == 0 && comment.Range == "" {
			return nil, fmt.Errorf("comment %d: either line or range is required", i+1)
		}
		if comment.Line != 0 && comment.Range != "" {
			return nil, fmt.Errorf("comment %d: cannot specify both line and range", i+1)
		}
		if comment.Type != "" && comment.Type != "review" && comment.Type != "issue" {
			return nil, fmt.Errorf("comment %d: type must be 'review' or 'issue'", i+1)
		}
	}

	// Validate review if present
	if config.Review != nil {
		if config.Review.Event != "" {
			validEvents := []string{"APPROVE", "REQUEST_CHANGES", "COMMENT"}
			isValid := false
			for _, validEvent := range validEvents {
				if config.Review.Event == validEvent {
					isValid = true
					break
				}
			}
			if !isValid {
				return nil, fmt.Errorf("review event must be one of: %s", strings.Join(validEvents, ", "))
			}
		}
	}

	return &config, nil
}

func processBatchComments(client github.GitHubAPI, owner, repo string, pr int, config *BatchConfig) error {
	// If we have a review configuration, create the review with comments
	if config.Review != nil {
		return processAsReview(client, owner, repo, pr, config)
	}

	// Otherwise, process comments individually
	return processIndividualComments(client, owner, repo, pr, config.Comments)
}

func processAsReview(client github.GitHubAPI, owner, repo string, pr int, config *BatchConfig) error {
	// Convert comments to review comment format
	var reviewComments []github.ReviewCommentInput

	for _, comment := range config.Comments {
		if comment.Type == "issue" {
			// Issue comments can't be part of a review, process separately
			if verbose {
				fmt.Printf("Processing issue comment separately: %s:%s\n", comment.File, formatLineOrRange(comment))
			}
			_, err := client.CreateIssueComment(owner, repo, pr, comment.Message)
			if err != nil {
				return fmt.Errorf("failed to create issue comment: %w", err)
			}
			continue
		}

		// Note: GitHub automatically uses the latest commit SHA for review comments

		// Create review comment input
		reviewComment := github.ReviewCommentInput{
			Body: expandSuggestions(comment.Message),
			Path: comment.File,
			Side: "RIGHT", // Default to RIGHT side (additions/new lines)
		}

		// Set line or range
		if comment.Range != "" {
			startLine, endLine, err := parseRange(comment.Range)
			if err != nil {
				return fmt.Errorf("invalid range %s: %w", comment.Range, err)
			}
			reviewComment.StartLine = startLine
			reviewComment.Line = endLine
		} else {
			reviewComment.Line = comment.Line
		}

		// Validate line exists in diff if validation is enabled
		if validateDiff {
			if err := validateCommentLine(client, owner, repo, pr, reviewComment); err != nil {
				return fmt.Errorf("review comment validation failed: %w", err)
			}
		}

		reviewComments = append(reviewComments, reviewComment)
	}

	// Create the review
	reviewInput := github.ReviewInput{
		Body:     config.Review.Body,
		Event:    config.Review.Event,
		Comments: reviewComments,
	}

	if reviewInput.Event == "" {
		reviewInput.Event = "COMMENT"
	}

	err := client.CreateReview(owner, repo, pr, reviewInput)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	fmt.Printf("✅ Successfully created review with %d comments\n", len(reviewComments))
	return nil
}

func processIndividualComments(client github.GitHubAPI, owner, repo string, pr int, comments []CommentConfig) error {
	successCount := 0

	for i, comment := range comments {
		if verbose {
			fmt.Printf("Processing comment %d/%d: %s:%s\n", i+1, len(comments), comment.File, formatLineOrRange(comment))
		}

		commentType := comment.Type
		if commentType == "" {
			commentType = "review" // Default to review comments
		}

		var err error
		if commentType == "issue" {
			_, err = client.CreateIssueComment(owner, repo, pr, expandSuggestions(comment.Message))
		} else {
			// For review comments, we need the commit SHA
			// Note: GitHub automatically uses the latest commit SHA for review comments

			reviewComment := github.ReviewCommentInput{
				Body: expandSuggestions(comment.Message),
				Path: comment.File,
				Side: "RIGHT", // Default to RIGHT side (additions/new lines)
			}

			// Set line or range
			if comment.Range != "" {
				startLine, endLine, err := parseRange(comment.Range)
				if err != nil {
					return fmt.Errorf("invalid range %s: %w", comment.Range, err)
				}
				reviewComment.StartLine = startLine
				reviewComment.Line = endLine
			} else {
				reviewComment.Line = comment.Line
			}

			// Validate line exists in diff if validation is enabled
			if validateDiff {
				if err := validateCommentLine(client, owner, repo, pr, reviewComment); err != nil {
					return fmt.Errorf("comment %d validation failed: %w", i+1, err)
				}
			}

			err = client.AddReviewComment(owner, repo, pr, reviewComment)
		}

		if err != nil {
			return fmt.Errorf("failed to create comment %d: %w", i+1, err)
		}

		successCount++
	}

	fmt.Printf("✅ Successfully created %d comments\n", successCount)
	return nil
}

// Helper functions
func formatLineOrRange(comment CommentConfig) string {
	if comment.Range != "" {
		return comment.Range
	}
	return fmt.Sprintf("%d", comment.Line)
}

func truncateMessage(message string, maxLen int) string {
	if len(message) <= maxLen {
		return message
	}
	return message[:maxLen-3] + "..."
}

func parseRange(rangeStr string) (startLine, endLine int, err error) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("range must be in format 'start-end'")
	}

	startLine, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start line: %w", err)
	}

	endLine, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end line: %w", err)
	}

	if startLine <= 0 || endLine <= 0 {
		return 0, 0, fmt.Errorf("line numbers must be positive")
	}

	if startLine > endLine {
		return 0, 0, fmt.Errorf("start line (%d) cannot be greater than end line (%d)", startLine, endLine)
	}

	return startLine, endLine, nil
}
