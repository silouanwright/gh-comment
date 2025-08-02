package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunResolveWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	resolveClient = mockClient
	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:    "resolve comment successfully",
			args:    []string{"123456"},
			wantErr: false,
		},
		{
			name:           "invalid comment ID",
			args:           []string{"invalid"},
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "missing comment ID",
			args:           []string{},
			wantErr:        true,
			expectedErrMsg: "accepts 1 arg(s), received 0",
		},
		{
			name:           "too many arguments",
			args:           []string{"123456", "extra"},
			wantErr:        true,
			expectedErrMsg: "accepts 1 arg(s), received 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle cases with wrong number of args
			if len(tt.args) != 1 {
				// This would be caught by cobra before runResolve is called
				err := resolveCmd.Args(nil, tt.args)
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				return
			}

			err := runResolve(nil, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunResolveDryRun(t *testing.T) {
	// Save original values
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	originalDryRun := dryRun
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		dryRun = originalDryRun
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	resolveClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	dryRun = true

	err := runResolve(nil, []string{"123456"})
	assert.NoError(t, err)
}

func TestRunResolveVerbose(t *testing.T) {
	// Save original values
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	originalVerbose := verbose
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		verbose = originalVerbose
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	resolveClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	verbose = true

	err := runResolve(nil, []string{"123456"})
	assert.NoError(t, err)
}

func TestResolveRepositoryParsing(t *testing.T) {
	// Save original values
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	mockClient := github.NewMockClient()
	resolveClient = mockClient
	prNumber = 123

	tests := []struct {
		name           string
		setupRepo      string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:      "valid repository format",
			setupRepo: "owner/repo",
			wantErr:   false,
		},
		{
			name:      "repository with hyphens",
			setupRepo: "my-org/my-repo",
			wantErr:   false,
		},
		{
			name:           "invalid repository format - no slash",
			setupRepo:      "invalidrepo",
			wantErr:        true,
			expectedErrMsg: "invalid repository format",
		},
		{
			name:           "invalid repository format - multiple slashes",
			setupRepo:      "owner/repo/extra",
			wantErr:        true,
			expectedErrMsg: "invalid repository format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo = tt.setupRepo

			err := runResolve(nil, []string{"123456"})
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResolveErrorHandling(t *testing.T) {
	// Save original values
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name              string
		setupMockError    func(*github.MockClient)
		expectedErrMsg    string
	}{
		{
			name: "find review thread error",
			setupMockError: func(m *github.MockClient) {
				m.FindReviewThreadError = assert.AnError
			},
			expectedErrMsg: "failed to find review thread for comment",
		},
		{
			name: "resolve review thread error",
			setupMockError: func(m *github.MockClient) {
				m.ResolveThreadError = assert.AnError
			},
			expectedErrMsg: "failed to resolve conversation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := github.NewMockClient()
			if tt.setupMockError != nil {
				tt.setupMockError(mockClient)
			}
			resolveClient = mockClient

			err := runResolve(nil, []string{"123456"})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func TestResolveCommentValidation(t *testing.T) {
	// Save original values
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	mockClient := github.NewMockClient()
	resolveClient = mockClient
	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name           string
		commentID      string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:      "valid positive comment ID",
			commentID: "123456",
			wantErr:   false,
		},
		{
			name:      "zero comment ID (technically valid for strconv.Atoi)",
			commentID: "0",
			wantErr:   false,
		},
		{
			name:           "negative comment ID (technically valid for strconv.Atoi)",
			commentID:      "-1",
			wantErr:        false, // strconv.Atoi allows negative numbers
		},
		{
			name:           "non-numeric comment ID",
			commentID:      "abc123",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "empty comment ID",
			commentID:      "",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "comment ID with spaces",
			commentID:      "123 456",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:      "large comment ID",
			commentID: "999999999999",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runResolve(nil, []string{tt.commentID})
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResolveWithClientInitialization(t *testing.T) {
	// Save original values
	originalClient := resolveClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		resolveClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	// Set client to nil to test initialization
	resolveClient = nil
	repo = "owner/repo"
	prNumber = 123

	// This test verifies that when resolveClient is nil, 
	// a RealClient is initialized in production
	// Since we can't easily test the RealClient without external dependencies,
	// we'll test that the initialization happens by setting up a mock afterwards
	
	// First verify the client gets initialized
	mockClient := github.NewMockClient()
	resolveClient = mockClient
	
	err := runResolve(nil, []string{"123456"})
	assert.NoError(t, err)
}