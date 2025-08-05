package cmd

import (
	"fmt"
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
	Use:   "add [pr] <comment>",
	Short: "Add a general discussion comment to a pull request",
	Long: heredoc.Doc(`
		Add a general discussion comment to a pull request conversation.

		This command creates issue-style comments that appear in the main PR conversation,
		not attached to specific code lines. These comments support threaded replies
		and are perfect for general discussion, approval, or high-level feedback.

		⚠️  IMPORTANT: This command does NOT create line-specific comments.
		    For line-specific code review comments, use: 'gh comment review'

		    Wrong: gh comment add 123 file.js 42 "comment"  ❌
		    Right: gh comment review 123 --comment file.js:42:"comment" ✅

		The comment message supports GitHub markdown formatting and can include
		code suggestions using the [SUGGEST: code] syntax. Use offset syntax
		[SUGGEST:+N: code] for lines below or [SUGGEST:-N: code] for lines above.
	`),
	Example: heredoc.Doc(`
		# General PR discussion comments
		$ gh comment add 123 "LGTM! Just a few minor suggestions below"
		$ gh comment add 123 "Thanks for addressing the security concerns"
		$ gh comment add 123 "This looks great - ready to merge after CI passes"

		# Multi-line comments with --message flags
		$ gh comment add 123 -m "Overall this is excellent work!" -m "The architecture is clean and the tests are comprehensive"

		# Auto-detect PR from current branch
		$ gh comment add "Looks good to merge!"

		# Approval with context
		$ gh comment add 123 "Approved! The performance improvements in this PR will make a huge difference"

		# Request for changes with discussion
		$ gh comment add 123 "Could you address the failing tests? Otherwise looks good to go"

		# Comments with suggestion syntax (auto-expands to GitHub suggestion blocks)
		$ gh comment add 123 "Consider using async/await: [SUGGEST: const result = await fetchData();]"
		$ gh comment add 123 "Add error handling above: [SUGGEST:-1: try {]"
		$ gh comment add 123 "Add timeout below: [SUGGEST:+2: const timeout = 5000;]"
	`),
	Args: cobra.RangeArgs(1, 2),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringArrayVarP(&messages, "message", "m", []string{}, "Add message (can be used multiple times for multi-line comments) (default: empty)")
	addCmd.Flags().BoolVar(&noExpandSuggestions, "no-expand-suggestions", false, "Disable automatic expansion of [SUGGEST:] and <<<SUGGEST>>> syntax (default: false)")
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
	var repository string
	var comment string
	var err error

	// Parse arguments for general PR comments
	if len(messages) > 0 {
		// Using --message flags
		if len(args) == 1 {
			// PR provided + --message flags
			pr, err = parsePositiveInt(args[0], "PR number")
			if err != nil {
				return err
			}
			// Get repository for explicitly provided PR
			repository, err = getCurrentRepo()
			if err != nil {
				return err
			}

			// Validate repository name
			if err := validateRepositoryName(repository); err != nil {
				return err
			}
		} else if len(args) == 0 {
			// Auto-detect PR + --message flags using centralized function
			repository, pr, err = getPRContext()
			if err != nil {
				return err
			}

			// Validate repository name
			if err := validateRepositoryName(repository); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid arguments when using --message flags: expected [pr] or no args")
		}
		comment = strings.Join(messages, "\n")
	} else if len(args) == 2 {
		// PR number provided + comment
		pr, err = parsePositiveInt(args[0], "PR number")
		if err != nil {
			return err
		}
		comment = args[1]
		// Get repository for explicitly provided PR
		repository, err = getCurrentRepo()
		if err != nil {
			return err
		}

		// Validate repository name
		if err := validateRepositoryName(repository); err != nil {
			return err
		}
	} else if len(args) == 1 {
		// Auto-detect PR from current branch + comment using centralized function
		repository, pr, err = getPRContext()
		if err != nil {
			return err
		}

		// Validate repository name
		if err := validateRepositoryName(repository); err != nil {
			return err
		}
		comment = args[0]
	} else {
		return fmt.Errorf("invalid arguments. Use: gh comment add [pr] <comment> OR gh comment add [pr] --message \"line1\" --message \"line2\"")
	}

	// Validate comment
	if strings.TrimSpace(comment) == "" {
		return fmt.Errorf("comment cannot be empty")
	}

	// Validate comment body length
	if err := validateCommentBody(comment); err != nil {
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
		fmt.Printf("Comment type: General discussion\n")
		fmt.Printf("Original comment: %s\n", comment)
		fmt.Printf("Transformed comment: %s\n", transformedComment)
	}

	if dryRun {
		fmt.Printf("Would add general comment to PR #%d:\n%s\n", pr, transformedComment)
		return nil
	}

	// Parse owner/repo
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid repository format: %s (expected owner/repo)", repository)
	}
	owner, repoName := parts[0], parts[1]

	// Create general issue comment (not review comment)
	createdComment, err := addClient.CreateIssueComment(owner, repoName, pr, transformedComment)
	if err != nil {
		return formatActionableError("comment creation", err)
	}

	// Success message
	fmt.Printf("✓ Added comment to PR #%d", pr)
	if createdComment != nil && createdComment.ID != 0 {
		fmt.Printf(" (ID: %d)", createdComment.ID)
	}
	fmt.Println()

	return nil
}
