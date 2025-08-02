package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	prNumber     int
	repo         string
	validateDiff bool
	dryRun       bool
	verbose      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-comment",
	Short: "Strategic GitHub PR commenting and review management",
	Long: heredoc.Doc(`
		gh-comment provides professional-grade tools for GitHub pull request
		comment management, systematic code reviews, and review process automation.

		Designed for developers, code review leads, and teams who need sophisticated
		comment workflows beyond GitHub's web interface capabilities.

		Strategic GitHub PR Commenting - Beyond the Web Interface:

		gh-comment is designed for professional code review workflows that require:
		• Systematic line-by-line code analysis with context
		• Bulk comment operations for comprehensive reviews
		• Advanced filtering for comment analysis and metrics
		• Automation integration for CI/CD review processes
		• Data export for review process optimization

		Perfect for: Senior developers, code review leads, QA teams, and DevOps engineers
		who need sophisticated comment management beyond GitHub's web interface.
	`),
	Example: heredoc.Doc(`
		Commands:
		  add           Add targeted comments to specific lines
		  add-review    Create draft reviews with multiple comments
		  batch         Process comments from YAML configuration files
		  edit          Modify existing comments
		  list          List and filter comments with advanced options
		  reply         Reply to comments and manage reactions
		  resolve       Resolve conversation threads
		  review        Create reviews with comments in one operation
		  submit-review Submit pending reviews with approval/changes
		  help          Help about any command

		Global Flags:
		  -p, --pr int        PR number (auto-detect from branch if omitted)
		  -R, --repo string   Repository (owner/repo format)
		      --dry-run       Show what would be commented without executing
		  -v, --verbose       Show detailed API interactions
		      --validate      Validate line exists in diff before commenting (default true)

		Filtering Flags (list command):
		      --author string     Filter by author (supports wildcards: 'user*')
		      --since string      Show comments after date ('2024-01-01', '1 week ago')
		      --until string      Show comments before date
		      --status string     Filter by status: open, resolved, all
		      --type string       Filter by type: issue, review, all
		  -q, --quiet            Minimal output for scripts

		Review Flags:
		      --event string      Review event: APPROVE, REQUEST_CHANGES, COMMENT
		      --comment strings   Add comments: file:line:message

		Examples:
		# Basic Operations
		$ gh comment list 123                           List all comments on PR #123
		$ gh comment add 123 "Looks good overall!"     Add general PR comment

		# Strategic Line Commenting (Unique Value)
		$ gh comment add 123 src/api.js 42 "This handles the rate limiting edge case - consider moving to middleware"
		$ gh comment add 123 auth.go 15:25 "This entire auth flow needs refactoring for OAuth2 compliance"
		$ gh comment add 123 database.py 156 "This query is vulnerable to SQL injection - use parameterized queries"

		# Advanced Filtering (Power User Features)
		$ gh comment list 123 --author "senior-dev*" --status open --since "1 week ago"
		$ gh comment list 123 --type review --author "*@company.com" --since "deployment-date"
		$ gh comment list 123 --status resolved --until "2024-01-01" --format json

		# Review Workflows (Professional Code Review)
		$ gh comment review 123 "Migration review complete" \
		  --comment src/api.js:42:"Add rate limiting middleware" \
		  --comment src/auth.js:15:20:"Update to OAuth2 flow" \
		  --comment tests/api_test.go:100:"Add edge case tests" \
		  --event REQUEST_CHANGES

		$ gh comment add-review 123 "Security audit findings" \
		  --comment auth.go:67:"Use crypto.randomBytes(32) for tokens" \
		  --comment api.js:134:140:"Extract business logic to service layer"

		# Batch Operations (Systematic Reviews)
		$ gh comment batch comprehensive-review.yaml
		$ gh comment batch --dry-run security-audit.yaml
		$ gh comment batch post-deployment-checklist.yaml --pr $(gh pr view --json number -q .number)

		# Conversation Management
		$ gh comment reply 2246362251 "Fixed in commit abc123" --resolve
		$ gh comment reply 3141344022 "Great catch! This would have caused issues in production" --reaction +1
		$ gh comment resolve --thread 2246362251 --reason "Addressed in latest commit"

		# Data Export & Analysis (Automation)
		$ gh comment list 123 --format json | jq '.comments[].author' | sort | uniq -c
		$ gh comment list 123 --format csv --since "2024-01-01" --output q1-review-data.csv
		$ gh comment list 123 --author "qa-team*" --format json | analyze-feedback.py

		# Automation & CI Integration
		$ gh comment add 123 src/security.js 67 "[SUGGEST: use crypto.randomBytes(32)]"
		$ for file in $(git diff --name-only); do gh comment add 123 "$file" 1 "Auto-generated security scan results"; done
		$ gh comment list --since "deployment-date" --type review --status open | review-blocker-analysis.sh

		# Advanced Comment Management
		$ gh comment edit 2246362251 "Updated: This rate limiting logic handles concurrent requests properly"
		$ gh comment list 123 --author "bot*" --format json | jq '.comments[].id' | xargs -I {} gh comment resolve {}
		$ gh comment add 123 performance.js 89:95 "Consider caching this expensive calculation"
	`),
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
