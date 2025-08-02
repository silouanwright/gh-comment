package cmd

import (
	"os"

	"github.com/silouanwright/gh-comment/internal/github"
)

// createGitHubClient creates the appropriate GitHub client based on environment
func createGitHubClient() (github.GitHubAPI, error) {
	// Check if we're in a test environment with mock server
	if mockURL := os.Getenv("MOCK_SERVER_URL"); mockURL != "" {
		return github.NewTestClient()
	}

	// Use real client for production
	return github.NewRealClient()
}
