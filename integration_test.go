package main

import (
	"os"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/silouanwright/gh-comment/cmd"
)

var mockServer *cmd.MockGitHubServer

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"gh-comment": main1,
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

// main1 is a wrapper around main that returns an exit code
func main1() int {
	// Call the actual main function and handle panics
	defer func() {
		if r := recover(); r != nil {
			// If main panics, return error code
			os.Exit(1)
		}
	}()

	// Call the real main function
	main()
	return 0
}
