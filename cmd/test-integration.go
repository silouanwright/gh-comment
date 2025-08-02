//go:build integration
// +build integration

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

var (
	cleanup        bool
	inspect        bool
	scenario       string
	force          bool
	integrationLog *log.Logger
)

// testIntegrationCmd represents the test-integration command
var testIntegrationCmd = &cobra.Command{
	Use:   "test-integration",
	Short: "Run integration tests against real GitHub PRs",
	Long: `Run comprehensive integration tests by creating real pull requests
and testing all gh-comment functionality live. This uses a "dogfooding"
approach - testing the extension on its own repository.

Integration tests are MANUAL-ONLY by default to avoid accidental API usage:
- Default: Never run automatically
- Use --force flag to run when needed
- Control with environment variables for CI/automation

Environment Variables:
  GH_COMMENT_INTEGRATION_TESTS=always    # Always run (for CI)
  GH_COMMENT_INTEGRATION_TESTS=never     # Never run (explicit disable)

Examples:
  # MANUAL RUN: Force integration tests to run (required for all runs)
  gh comment test-integration --force

  # Run all integration tests with auto-cleanup
  gh comment test-integration --force

  # Run specific scenario and leave PR open for inspection
  gh comment test-integration --scenario=comments --inspect

  # Run tests with cleanup disabled for debugging
  gh comment test-integration --no-cleanup`,
	RunE: runTestIntegration,
}

func init() {
	rootCmd.AddCommand(testIntegrationCmd)

	testIntegrationCmd.Flags().BoolVar(&cleanup, "cleanup", true, "Auto-close PR after tests complete")
	testIntegrationCmd.Flags().BoolVar(&inspect, "inspect", false, "Leave PR open for manual inspection (implies --no-cleanup)")
	testIntegrationCmd.Flags().StringVar(&scenario, "scenario", "", "Run specific test scenario only (comments, reviews, reactions, batch, suggestions)")
	testIntegrationCmd.Flags().BoolVar(&force, "force", false, "Force integration tests to run (bypasses frequency controls)")
}

func runTestIntegration(cmd *cobra.Command, args []string) error {
	// Check if integration tests should run based on frequency controls
	if !force && !shouldRunIntegrationTests() {
		fmt.Println("⏭️  Skipping integration tests (use --force flag to override)")
		return nil
	}

	// Setup logging
	if err := setupIntegrationLogging(); err != nil {
		return fmt.Errorf("failed to setup logging: %w", err)
	}

	integrationLog.Println("🚀 Starting integration tests...")

	// If inspect is enabled, disable cleanup
	if inspect {
		cleanup = false
		integrationLog.Println("📋 Inspect mode enabled - PR will remain open")
	}

	// Get current repository
	currentRepo, err := getCurrentRepo()
	if err != nil {
		return fmt.Errorf("failed to get current repository: %w", err)
	}
	integrationLog.Printf("📂 Testing repository: %s", currentRepo)

	// Create test PR
	branchName := fmt.Sprintf("integration-test-%d", time.Now().Unix())
	prNumber, err := createTestPR(branchName)
	if err != nil {
		return fmt.Errorf("failed to create test PR: %w", err)
	}

	integrationLog.Printf("🔧 Created test PR #%d on branch %s", prNumber, branchName)

	// Defer cleanup if enabled
	if cleanup {
		defer func() {
			integrationLog.Printf("🧹 Cleaning up test PR #%d...", prNumber)
			if err := cleanupTestPR(branchName, prNumber); err != nil {
				integrationLog.Printf("⚠️  Cleanup failed: %v", err)
			} else {
				integrationLog.Printf("✅ Cleanup completed")
			}
		}()
	}

	// Run test scenarios
	if scenario != "" {
		return runSpecificScenario(scenario, prNumber)
	}

	return runAllScenarios(prNumber)
}

func setupIntegrationLogging() error {
	// Create results directory if it doesn't exist
	resultsDir := "integration-tests/results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("20060102-150405")
	logPath := filepath.Join(resultsDir, fmt.Sprintf("integration-%s.log", timestamp))

	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}

	integrationLog = log.New(logFile, "", log.LstdFlags)
	fmt.Printf("📝 Integration test log: %s\n", logPath)

	return nil
}

func createTestPR(branchName string) (int, error) {
	integrationLog.Printf("Creating test branch: %s", branchName)

	// Create and checkout new branch
	if err := runGitCommand("checkout", "-b", branchName); err != nil {
		return 0, fmt.Errorf("failed to create branch: %w", err)
	}

	// Copy template file
	templatePath := "integration-tests/templates/dummy-code.js"
	targetPath := fmt.Sprintf("test-file-%d.js", time.Now().Unix())

	if err := copyTemplateFile(templatePath, targetPath); err != nil {
		return 0, fmt.Errorf("failed to copy template: %w", err)
	}

	// Git add and commit
	if err := runGitCommand("add", targetPath); err != nil {
		return 0, fmt.Errorf("failed to git add: %w", err)
	}

	commitMsg := fmt.Sprintf("Integration test: %s", branchName)
	if err := runGitCommand("commit", "-m", commitMsg); err != nil {
		return 0, fmt.Errorf("failed to commit: %w", err)
	}

	// Push branch
	if err := runGitCommand("push", "-u", "origin", branchName); err != nil {
		return 0, fmt.Errorf("failed to push branch: %w", err)
	}

	// Create PR using gh CLI
	prTitle := fmt.Sprintf("Integration Test: %s", branchName)
	prBody := `This is an automated integration test PR created by gh-comment.

**Test Purpose**: Validate all gh-comment functionality against real GitHub APIs

**What this tests**:
- Line-specific comments
- Review comments and submissions
- Reactions and replies
- Batch operations
- Suggestion syntax

**Cleanup**: This PR will be automatically closed unless --inspect flag was used.`

	stdout, _, err := gh.Exec("pr", "create", "--title", prTitle, "--body", prBody)
	if err != nil {
		return 0, fmt.Errorf("failed to create PR: %w", err)
	}

	// Extract PR number from output
	// gh pr create typically outputs the PR URL
	output := stdout.String()
	integrationLog.Printf("PR created: %s", output)

	// Get PR number using gh pr view
	prStdout, _, err := gh.Exec("pr", "view", "--json", "number", "-q", ".number")
	if err != nil {
		return 0, fmt.Errorf("failed to get PR number: %w", err)
	}

	var prNum int
	if _, err := fmt.Sscanf(prStdout.String(), "%d", &prNum); err != nil {
		return 0, fmt.Errorf("failed to parse PR number: %w", err)
	}

	return prNum, nil
}

func cleanupTestPR(branchName string, prNumber int) error {
	// Close PR
	_, _, err := gh.Exec("pr", "close", fmt.Sprintf("%d", prNumber))
	if err != nil {
		integrationLog.Printf("Failed to close PR: %v", err)
	}

	// Switch to main branch
	if err := runGitCommand("checkout", "main"); err != nil {
		integrationLog.Printf("Failed to checkout main: %v", err)
	}

	// Delete local branch
	if err := runGitCommand("branch", "-D", branchName); err != nil {
		integrationLog.Printf("Failed to delete local branch: %v", err)
	}

	// Delete remote branch
	if err := runGitCommand("push", "origin", "--delete", branchName); err != nil {
		integrationLog.Printf("Failed to delete remote branch: %v", err)
	}

	return nil
}

func runSpecificScenario(scenarioName string, prNumber int) error {
	integrationLog.Printf("🎯 Running scenario: %s", scenarioName)

	switch scenarioName {
	case "comments":
		return runBasicCommentsScenario(prNumber)
	case "reviews":
		return runReviewWorkflowScenario(prNumber)
	case "reactions":
		return runReactionsRepliesScenario(prNumber)
	case "batch":
		return runBatchOperationsScenario(prNumber)
	case "suggestions":
		return runSuggestionsScenario(prNumber)
	default:
		return fmt.Errorf("unknown scenario: %s", scenarioName)
	}
}

func runAllScenarios(prNumber int) error {
	scenarios := []string{"comments", "reviews", "reactions", "batch", "suggestions"}

	for _, s := range scenarios {
		integrationLog.Printf("🏃 Running scenario: %s", s)
		if err := runSpecificScenario(s, prNumber); err != nil {
			integrationLog.Printf("❌ Scenario %s failed: %v", s, err)
			return fmt.Errorf("scenario %s failed: %w", s, err)
		}
		integrationLog.Printf("✅ Scenario %s completed", s)
	}

	integrationLog.Println("🎉 All integration tests completed successfully!")
	return nil
}

// Helper functions
func runGitCommand(args ...string) error {
	_, _, err := gh.Exec("api", "--method", "GET", "/user") // Test gh auth first
	if err != nil {
		return fmt.Errorf("gh CLI authentication required: %w", err)
	}

	// Use git command directly for git operations
	gitArgs := append([]string{"!", "git"}, args...)
	_, _, err = gh.Exec(gitArgs...)
	return err
}

func copyTemplateFile(src, dst string) error {
	// First ensure template exists, if not create it
	if _, err := os.Stat(src); os.IsNotExist(err) {
		if err := createDummyTemplate(src); err != nil {
			return fmt.Errorf("failed to create template: %w", err)
		}
	}

	// Copy file content (simple implementation)
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, content, 0644)
}

func createDummyTemplate(path string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	template := `// Integration Test File - Contains intentional issues for commenting
function calculateTotal(items) {
    let total = 0;
    for (let i = 0; i < items.length; i++) {
        total += items[i].price * items[i].quantity; // Potential null pointer
    }
    return total; // Missing input validation
}

// TODO: Add error handling
// FIXME: Handle empty arrays
const processOrder = (order) => {
    const total = calculateTotal(order.items);
    return { total, tax: total * 0.08 }; // Hardcoded tax rate
};

module.exports = { calculateTotal, processOrder };
`

	return os.WriteFile(path, []byte(template), 0644)
}

// shouldRunIntegrationTests checks if integration tests should run based on environment variables and frequency controls
func shouldRunIntegrationTests() bool {
	// Check environment variables for controls
	if os.Getenv("GH_COMMENT_INTEGRATION_TESTS") == "always" {
		return true
	}

	if os.Getenv("GH_COMMENT_INTEGRATION_TESTS") == "never" {
		return false
	}

	// Default behavior: NEVER run automatically (manual only)
	return false
}

// getIntegrationTestFrequency returns the frequency for running integration tests
func getIntegrationTestFrequency() int {
	freqStr := os.Getenv("GH_COMMENT_INTEGRATION_FREQUENCY")
	if freqStr == "" {
		return 10 // Default: every 10th run
	}

	var freq int
	if n, err := fmt.Sscanf(freqStr, "%d", &freq); err == nil && n == 1 && freq > 0 {
		return freq
	}

	return 10 // Fallback to default
}

// getIntegrationRunCount gets the current run count from file
func getIntegrationRunCount() int {
	counterFile := getIntegrationCounterFile()

	data, err := os.ReadFile(counterFile)
	if err != nil {
		return 1 // First run
	}

	var count int
	if _, err := fmt.Sscanf(string(data), "%d", &count); err != nil {
		return 1
	}

	return count
}

// updateIntegrationRunCount updates the run count in file
func updateIntegrationRunCount(count int) {
	counterFile := getIntegrationCounterFile()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(counterFile), 0755); err != nil {
		return // Fail silently
	}

	countStr := fmt.Sprintf("%d", count)
	os.WriteFile(counterFile, []byte(countStr), 0644)
}

// getIntegrationCounterFile returns the path to the counter file
func getIntegrationCounterFile() string {
	// Try to use a system temp directory that persists across runs
	tempDir := os.Getenv("TMPDIR")
	if tempDir == "" {
		tempDir = "/tmp"
	}

	return filepath.Join(tempDir, ".gh-comment-integration-counter")
}
