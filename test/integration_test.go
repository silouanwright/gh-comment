package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/silouanwright/gh-comment/cmd"
)

var mockServer *cmd.MockGitHubServer

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"gh-comment": func() int {
			// Execute the root command and return exit code
			if err := cmd.Execute(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return 1
			}
			return 0
		},
	}))
}

func TestIntegration(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
		Setup: func(env *testscript.Env) error {
			// Set up test environment
			env.Setenv("GH_TOKEN", "test-token")
			env.Setenv("GH_HOST", "github.com")
			return nil
		},
		Condition: func(cond string) (bool, error) {
			switch cond {
			case "has-gh":
				// Check if gh CLI is available
				_, err := os.Stat("/usr/local/bin/gh")
				return err == nil, nil
			default:
				return false, nil
			}
		},
	})
}

func TestEnhancedIntegration(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/enhanced-scripts",
		Setup: func(env *testscript.Env) error {
			// Start mock GitHub API server
			mockServer = cmd.NewMockGitHubServer()

			// DETERMINISTIC SETUP: Always set up all scenarios
			// This ensures consistent test data regardless of condition checks
			mockServer.SetupTestScenario("basic")
			mockServer.SetupTestScenario("security-review")

			// Set up test environment to use mock server
			env.Setenv("GH_TOKEN", "test-token")
			env.Setenv("GH_HOST", strings.TrimPrefix(mockServer.URL(), "http://"))
			env.Setenv("MOCK_SERVER_URL", mockServer.URL())

			// Set up test repository context
			env.Setenv("GH_REPO", "test-owner/test-repo")

			return nil
		},
		Condition: func(cond string) (bool, error) {
			switch cond {
			case "mock-server":
				return mockServer != nil, nil
			case "scenario:basic", "scenario:security-review":
				// Scenarios are always available now
				return true, nil
			default:
				return false, nil
			}
		},
	})
}

// No longer needed - main1 function removed as we're using cmd.Execute() directly
