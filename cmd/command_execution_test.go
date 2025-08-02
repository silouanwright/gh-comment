package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommandExecution tests the actual command execution paths to increase coverage
func TestCommandExecution(t *testing.T) {
	// Save original environment and client
	originalRepo := os.Getenv("GH_REPO")
	originalPR := os.Getenv("GH_PR")
	originalListClient := listClient

	// Set test environment
	os.Setenv("GH_REPO", "owner/repo")
	os.Setenv("GH_PR", "123")

	// Set up mock client for list commands
	listClient = github.NewMockClient()

	// Restore environment after test
	defer func() {
		if originalRepo != "" {
			os.Setenv("GH_REPO", originalRepo)
		} else {
			os.Unsetenv("GH_REPO")
		}
		if originalPR != "" {
			os.Setenv("GH_PR", originalPR)
		} else {
			os.Unsetenv("GH_PR")
		}
		listClient = originalListClient
	}()

	tests := []struct {
		name         string
		command      string
		args         []string
		flags        map[string]string
		wantErr      bool
		wantContains []string
	}{
		{
			name:    "list command with PR argument",
			command: "list",
			args:    []string{"123"},
			wantErr: false, // Should succeed with mock client
		},
		{
			name:         "list command with invalid PR",
			command:      "list",
			args:         []string{"invalid"},
			wantErr:      true,
			wantContains: []string{"invalid PR number 'invalid': must be a valid integer"},
		},
		{
			name:         "reply command with invalid comment ID",
			command:      "reply",
			args:         []string{"invalid", "message"},
			wantErr:      true,
			wantContains: []string{"invalid comment ID 'invalid': must be a valid integer"},
		},
		{
			name:         "reply command with no message or action",
			command:      "reply",
			args:         []string{"123456"},
			wantErr:      true,
			wantContains: []string{"must provide either a message, --reaction, --remove-reaction, or --resolve"},
		},
		{
			name:         "reply command with invalid reaction",
			command:      "reply",
			args:         []string{"123456", "message"},
			flags:        map[string]string{"reaction": "invalid"},
			wantErr:      true,
			wantContains: []string{"invalid reaction 'invalid': must be one of"},
		},
		{
			name:    "reply command with both reaction and remove-reaction",
			command: "reply",
			args:    []string{"123456", "message"},
			flags: map[string]string{
				"reaction":        "+1",
				"remove-reaction": "heart",
			},
			wantErr:      true,
			wantContains: []string{"cannot use both --reaction and --remove-reaction"},
		},
		{
			name:         "reply command with invalid type",
			command:      "reply",
			args:         []string{"123456", "message"},
			flags:        map[string]string{"type": "invalid"},
			wantErr:      true,
			wantContains: []string{"invalid type 'invalid': must be either 'issue' or 'review'"},
		},
		{
			name:         "resolve command with invalid comment ID",
			command:      "resolve",
			args:         []string{"invalid"},
			wantErr:      true,
			wantContains: []string{"invalid comment ID 'invalid': must be a valid integer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to capture output
			var output bytes.Buffer

			// Create the root command
			rootCmd := &cobra.Command{Use: "gh-comment"}

			// Add the specific command being tested
			var cmd *cobra.Command
			switch tt.command {
			case "list":
				cmd = createListCommand()
			case "reply":
				cmd = createReplyCommand()
			case "resolve":
				cmd = createResolveCommand()
			default:
				t.Fatalf("Unknown command: %s", tt.command)
			}

			rootCmd.AddCommand(cmd)

			// Set flags if provided
			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				require.NoError(t, err)
			}

			// Set output
			rootCmd.SetOut(&output)
			rootCmd.SetErr(&output)

			// Prepare arguments
			args := []string{tt.command}
			args = append(args, tt.args...)
			rootCmd.SetArgs(args)

			// Execute command
			err := rootCmd.Execute()

			// Check error expectation
			if tt.wantErr {
				assert.Error(t, err)

				// Check error message contains expected text
				if len(tt.wantContains) > 0 {
					errorMsg := err.Error()
					for _, want := range tt.wantContains {
						assert.Contains(t, errorMsg, want, "Error should contain: %s", want)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// createListCommand creates a list command for testing
func createListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [pr]",
		Short: "List comments on a PR",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runList,
	}

	cmd.Flags().BoolVar(&showResolved, "resolved", false, "Include resolved comments")
	cmd.Flags().BoolVar(&onlyUnresolved, "unresolved", false, "Show only unresolved comments")
	cmd.Flags().StringVar(&author, "author", "", "Filter comments by author")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output without URLs and IDs")
	cmd.Flags().BoolVar(&hideAuthors, "hide-authors", false, "Hide author names for privacy")

	return cmd
}

// createReplyCommand creates a reply command for testing
func createReplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply <comment-id> [message]",
		Short: "Reply to a comment",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runReply,
	}

	cmd.Flags().StringVar(&commentType, "type", "review", "Comment type (issue or review)")
	cmd.Flags().StringVar(&reaction, "reaction", "", "Add reaction")
	cmd.Flags().StringVar(&removeReaction, "remove-reaction", "", "Remove reaction")
	cmd.Flags().BoolVar(&resolveConversation, "resolve", false, "Resolve conversation")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done")

	return cmd
}

// createResolveCommand creates a resolve command for testing
func createResolveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve <comment-id>",
		Short: "Resolve a conversation",
		Args:  cobra.ExactArgs(1),
		RunE:  runResolve,
	}

	cmd.Flags().IntVarP(&prNumber, "pr", "p", 0, "PR number")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done")

	return cmd
}

// TestHelperFunctions tests utility functions to increase coverage
func TestHelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func() error
		wantErr  bool
	}{
		{
			name: "getCurrentRepo with flag set",
			function: func() error {
				// Test that flag takes precedence
				originalRepo := repo
				repo = "test/repo"
				defer func() {
					repo = originalRepo
				}()

				result, err := getCurrentRepo()
				if err != nil {
					return err
				}
				if result != "test/repo" {
					return fmt.Errorf("expected test/repo, got %s", result)
				}
				return nil
			},
			wantErr: false,
		},
		// Removed test that calls real gh CLI - covered by integration tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.function()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPRContext tests the getPRContext function
func TestPRContext(t *testing.T) {
	// Save original global variables
	originalRepo := repo
	originalPR := prNumber

	defer func() {
		repo = originalRepo
		prNumber = originalPR
	}()

	t.Run("successful PR context with flags", func(t *testing.T) {
		// Set global variables as if they were set by flags
		repo = "owner/repo"
		prNumber = 123

		gotRepo, gotPR, err := getPRContext()
		assert.NoError(t, err)
		assert.Equal(t, "owner/repo", gotRepo)
		assert.Equal(t, 123, gotPR)
	})

	t.Run("auto-detection handles errors gracefully", func(t *testing.T) {
		// Clear global variables to test auto-detection
		repo = ""
		prNumber = 0

		// This will use gh repo view and gh pr view
		// In CI/different environments, this will likely fail
		// We just verify it doesn't panic and gives a reasonable error
		gotRepo, gotPR, err := getPRContext()

		// In most CI environments, this will fail - that's expected
		if err == nil {
			// If it succeeds (local dev), verify the results
			assert.NotEmpty(t, gotRepo)
			assert.Greater(t, gotPR, 0)
		} else {
			// Should contain helpful error message about repo or PR
			errorMsg := err.Error()
			assert.True(t,
				strings.Contains(errorMsg, "PR") ||
					strings.Contains(errorMsg, "repository") ||
					strings.Contains(errorMsg, "gh execution failed"),
				"Error should mention PR, repository, or gh execution: %s", errorMsg)
		}
	})

	// Removed test that calls real gh CLI - covered by integration tests
}
