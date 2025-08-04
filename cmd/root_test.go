package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentRepo(t *testing.T) {
	// Save original state
	originalRepo := repo
	defer func() { repo = originalRepo }()

	// Test the main branch where repo is already set (most common case)
	t.Run("repo already set in global variable", func(t *testing.T) {
		repo = "owner/repo"
		
		gotRepo, err := getCurrentRepo()
		
		assert.NoError(t, err)
		assert.Equal(t, "owner/repo", gotRepo)
	})

	// Test with different repo formats
	testRepos := []string{
		"myorg/myrepo",
		"user/project",
		"org-name/repo-name",
	}

	for _, testRepo := range testRepos {
		t.Run("repo set to "+testRepo, func(t *testing.T) {
			repo = testRepo
			
			gotRepo, err := getCurrentRepo()
			
			assert.NoError(t, err)
			assert.Equal(t, testRepo, gotRepo)
		})
	}
}


func TestGetCurrentRepo_ValidationIntegration(t *testing.T) {
	// Save original state
	originalRepo := repo
	defer func() { repo = originalRepo }()

	// Test integration with validateRepositoryName
	tests := []struct {
		name      string
		setupRepo string
		wantErr   bool
	}{
		{
			name:      "valid repository format",
			setupRepo: "owner/repo",
			wantErr:   false,
		},
		{
			name:      "invalid repository format",
			setupRepo: "invalid-repo-format",
			wantErr:   false, // getCurrentRepo doesn't validate format, just returns what's set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo = tt.setupRepo

			gotRepo, err := getCurrentRepo()

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.setupRepo, gotRepo)
			}

			// Test that we could validate the result if needed
			if gotRepo != "" {
				validationErr := validateRepositoryName(gotRepo)
				if tt.setupRepo == "owner/repo" {
					assert.NoError(t, validationErr, "Valid repo should pass validation")
				} else {
					// Note: validateRepositoryName would catch invalid formats
				}
			}
		})
	}
}

func TestGetCurrentRepo_EmptyGlobalVariable(t *testing.T) {
	// Save original state
	originalRepo := repo
	defer func() { repo = originalRepo }()

	// Test the branch where repo is empty (would trigger gh CLI call)
	t.Run("empty repo global variable", func(t *testing.T) {
		repo = "" // This will trigger the gh CLI execution path

		gotRepo, err := getCurrentRepo()

		// In test environment, this may succeed or fail depending on setup
		// We're mainly testing that the function handles the empty repo case
		// and doesn't panic or have unexpected behavior
		if err != nil {
			// Expected case in test environment - gh CLI call may fail
			assert.Contains(t, err.Error(), "failed to get current repository",
				"Error should be descriptive when gh CLI fails")
		} else {
			// If it succeeds, we should get a valid-looking repo name
			assert.NotEmpty(t, gotRepo, "If no error, repo should not be empty")
			// The repo should be trimmed of whitespace
			assert.Equal(t, gotRepo, strings.TrimSpace(gotRepo), "Repo should be trimmed")
		}
	})
}

func TestGetCurrentPRZeroCase(t *testing.T) {
	// Save original state
	originalPRNumber := prNumber
	defer func() { prNumber = originalPRNumber }()

	// Test the branch where prNumber is 0 (would trigger gh CLI call)
	// This is the missing coverage case not tested in current_pr_test.go
	t.Run("zero prNumber triggers gh CLI path", func(t *testing.T) {
		prNumber = 0 // This will trigger the gh CLI execution path

		gotPR, err := getCurrentPR()

		// In test environment, this will likely fail since gh CLI call may not work
		// We're mainly testing that the function handles the zero prNumber case
		// and provides descriptive error messages
		if err != nil {
			// Expected case in test environment - gh CLI call may fail
			assert.Contains(t, err.Error(), "failed to get current PR",
				"Error should be descriptive when gh CLI fails")
			assert.Contains(t, err.Error(), "try specifying --pr",
				"Error should suggest using --pr flag")
			assert.Equal(t, 0, gotPR, "PR should be 0 when error occurs")
		} else {
			// If it succeeds, we should get a positive PR number
			assert.Greater(t, gotPR, 0, "If no error, PR should be positive")
		}
	})
}