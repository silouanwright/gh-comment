package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	prNumber     int
	repo         string
	tone         string
	validateDiff bool
	dryRun       bool
	verbose      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-comment",
	Short: "Strategic line-specific PR commenting for GitHub CLI",
	Long: `gh-comment is the first GitHub CLI extension for strategic, line-specific PR commenting workflows.

Add targeted comments to specific lines in pull requests, create professional reviews,
and streamline your code review process with batch operations and smart tone transformation.

Examples:
  # Add a single line comment
  gh comment add 123 src/api.js 42 "this handles the rate limiting edge case"
  
  # Create a review with multiple comments
  gh comment review 123 "Migration review" \
    --comment src/api.js:42:"rate limiting fix" \
    --comment src/auth.js:15:20:"updated auth flow"
  
  # Process comments from a config file
  gh comment batch 123 comments.yaml`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().IntVarP(&prNumber, "pr", "p", 0, "PR number (auto-detect from branch if omitted)")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "R", "", "Repository (owner/repo format)")
	rootCmd.PersistentFlags().StringVar(&tone, "tone", "casual", "Comment tone: casual|formal|technical")
	rootCmd.PersistentFlags().BoolVar(&validateDiff, "validate", true, "Validate line exists in diff before commenting")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be commented without executing")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed API interactions")
}

// Helper function to get current repository if not specified
func getCurrentRepo() (string, error) {
	if repo != "" {
		return repo, nil
	}
	
	// Use gh CLI to get current repository
	stdout, _, err := gh.Exec("repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner")
	if err != nil {
		return "", fmt.Errorf("failed to get current repository: %w", err)
	}
	
	return strings.TrimSpace(stdout.String()), nil
}

// Helper function to get current PR number if not specified
func getCurrentPR() (int, error) {
	if prNumber != 0 {
		return prNumber, nil
	}
	
	// Use gh CLI to get PR for current branch
	stdout, _, err := gh.Exec("pr", "view", "--json", "number", "-q", ".number")
	if err != nil {
		return 0, fmt.Errorf("failed to get current PR: %w (try specifying --pr)", err)
	}
	
	prStr := strings.TrimSpace(stdout.String())
	pr, err := strconv.Atoi(prStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PR number: %s", prStr)
	}
	
	return pr, nil
}
