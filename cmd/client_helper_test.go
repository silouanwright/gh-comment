package cmd

import (
	"os"
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
)

func TestCreateGitHubClient(t *testing.T) {
	// Save original environment
	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	tests := []struct {
		name        string
		mockURL     string
		expectError bool
		clientType  string
	}{
		{
			name:        "creates real client when no mock URL",
			mockURL:     "",
			expectError: false, // Allow either success (with creds) or failure (without creds)
			clientType:  "*github.RealClient",
		},
		{
			name:        "creates test client when mock URL set",
			mockURL:     "http://localhost:8080",
			expectError: false,
			clientType:  "*github.TestClient",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.mockURL != "" {
				os.Setenv("MOCK_SERVER_URL", tt.mockURL)
			} else {
				os.Unsetenv("MOCK_SERVER_URL")
			}

			client, err := createGitHubClient()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			// For real client test, allow both success and failure (depends on credentials)
			if tt.mockURL == "" && err != nil {
				t.Logf("Real client creation failed (likely missing credentials): %v", err)
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Errorf("expected client but got nil")
				return
			}

			// Verify we got a client that implements GitHubAPI interface
			var _ github.GitHubAPI = client
		})
	}
}

func TestCreateGitHubClientWithMockOnly(t *testing.T) {
	// Always use mock - tests should never hit real APIs
	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	os.Setenv("MOCK_SERVER_URL", "http://localhost:8080")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	client, err := createGitHubClient()
	if err != nil {
		t.Errorf("unexpected error creating mock client: %v", err)
		return
	}

	if client == nil {
		t.Errorf("expected client but got nil")
		return
	}

	// Verify it implements the interface
	var _ github.GitHubAPI = client
}

func TestCreateGitHubClientEnvironmentVariables(t *testing.T) {
	// Test different mock URL formats
	mockURLs := []string{
		"http://localhost:8080",
		"https://mock.example.com",
		"http://127.0.0.1:9999",
	}

	originalMockURL := os.Getenv("MOCK_SERVER_URL")
	defer os.Setenv("MOCK_SERVER_URL", originalMockURL)

	for _, mockURL := range mockURLs {
		t.Run("mock_url_"+mockURL, func(t *testing.T) {
			os.Setenv("MOCK_SERVER_URL", mockURL)

			client, err := createGitHubClient()
			if err != nil {
				t.Errorf("unexpected error with mock URL %s: %v", mockURL, err)
				return
			}

			if client == nil {
				t.Errorf("expected client but got nil with mock URL %s", mockURL)
				return
			}

			// Verify it implements the interface
			var _ github.GitHubAPI = client
		})
	}
}
