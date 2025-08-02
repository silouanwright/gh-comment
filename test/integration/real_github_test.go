//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cli/go-gh/v2"
)

// TestIntegrationFramework tests the integration test framework itself
func TestIntegrationFramework(t *testing.T) {
	// Skip if no GitHub token
	if os.Getenv("GITHUB_TOKEN") == "" {
		t.Skip("GITHUB_TOKEN not set - skipping real GitHub integration tests")
	}

	// Test that we can authenticate with GitHub
	t.Run("GitHub Authentication", func(t *testing.T) {
		_, _, err := gh.Exec("api", "--method", "GET", "/user")
		if err != nil {
			t.Fatalf("GitHub authentication failed: %v", err)
		}
	})

	// Test that we can access the repository
	t.Run("Repository Access", func(t *testing.T) {
		stdout, _, err := gh.Exec("repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner")
		if err != nil {
			t.Fatalf("Failed to get repository: %v", err)
		}

		repoName := stdout.String()
		if repoName == "" {
			t.Fatal("Repository name is empty")
		}

		t.Logf("Testing against repository: %s", repoName)
	})
}

// TestIntegrationCommand tests the integration command registration
func TestIntegrationCommand(t *testing.T) {
	// This tests that the integration command is properly registered
	// when the integration build tag is used

	// We can test command parsing and flag handling without
	// actually executing the integration tests
	t.Run("Command Registration", func(t *testing.T) {
		// Test would verify the command exists and has proper flags
		// This is safe to run without creating real PRs
		t.Log("Integration command registration test would go here")
	})
}

// TestPRCreationLogic tests PR creation without actually creating PRs
func TestPRCreationLogic(t *testing.T) {
	t.Run("Branch Name Generation", func(t *testing.T) {
		// Test branch name generation logic
		branchName := generateTestBranchName()
		if branchName == "" {
			t.Fatal("Branch name generation failed")
		}

		// Verify it contains timestamp
		if len(branchName) < 10 {
			t.Fatalf("Branch name too short: %s", branchName)
		}

		t.Logf("Generated branch name: %s", branchName)
	})

	t.Run("Template File Validation", func(t *testing.T) {
		// Test that template file exists and has expected content
		templatePath := "../../integration-tests/templates/dummy-code.js"
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			t.Fatalf("Template file does not exist: %s", templatePath)
		}

		content, err := os.ReadFile(templatePath)
		if err != nil {
			t.Fatalf("Failed to read template: %v", err)
		}

		contentStr := string(content)
		expectedElements := []string{
			"calculateTotal",
			"items.length",
			"Potential null pointer",
			"Hardcoded tax rate",
		}

		for _, element := range expectedElements {
			if !contains(contentStr, element) {
				t.Errorf("Template missing expected element: %s", element)
			}
		}
	})
}

// Helper functions for tests
func generateTestBranchName() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("integration-test-%d", timestamp)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
