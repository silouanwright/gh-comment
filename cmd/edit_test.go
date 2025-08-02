package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunEditWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := editClient
	originalRepo := repo
	defer func() {
		editClient = originalClient
		repo = originalRepo
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	editClient = mockClient
	repo = "owner/repo"

	tests := []struct {
		name           string
		args           []string
		setupMessages  []string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:    "edit comment with message as positional arg",
			args:    []string{"123456", "Updated comment message"},
			wantErr: false,
		},
		{
			name:          "edit comment with --message flag",
			args:          []string{"123456"},
			setupMessages: []string{"Updated via flag"},
			wantErr:       false,
		},
		{
			name:          "edit comment with multiple --message flags",
			args:          []string{"123456"},
			setupMessages: []string{"Line 1", "Line 2", "Line 3"},
			wantErr:       false,
		},
		{
			name:          "edit comment with both positional and flag (positional wins)",
			args:          []string{"123456", "Positional message"},
			setupMessages: []string{"Flag message"},
			wantErr:       false,
		},
		{
			name:           "invalid comment ID",
			args:           []string{"invalid", "Message"},
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "missing message",
			args:           []string{"123456"},
			setupMessages:  []string{}, // No messages
			wantErr:        true,
			expectedErrMsg: "must provide either a message argument or --message flags",
		},
		{
			name:           "missing comment ID",
			args:           []string{},
			wantErr:        true,
			expectedErrMsg: "accepts between 1 and 2 arg(s), received 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			editMessages = tt.setupMessages

			// Handle cases with wrong number of args
			if len(tt.args) < 1 || len(tt.args) > 2 {
				// This would be caught by cobra before runEdit is called
				err := editCmd.Args(nil, tt.args)
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				return
			}

			err := runEdit(nil, tt.args)
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

func TestRunEditDryRun(t *testing.T) {
	// Save original values
	originalClient := editClient
	originalRepo := repo
	originalDryRun := dryRun
	defer func() {
		editClient = originalClient
		repo = originalRepo
		dryRun = originalDryRun
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	editClient = mockClient
	repo = "owner/repo"
	dryRun = true

	err := runEdit(nil, []string{"123456", "Test message"})
	assert.NoError(t, err)
}

func TestRunEditVerbose(t *testing.T) {
	// Save original values
	originalClient := editClient
	originalRepo := repo
	originalVerbose := verbose
	defer func() {
		editClient = originalClient
		repo = originalRepo
		verbose = originalVerbose
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	editClient = mockClient
	repo = "owner/repo"
	verbose = true

	err := runEdit(nil, []string{"123456", "Test message"})
	assert.NoError(t, err)
}

func TestEditMessageHandling(t *testing.T) {
	// Save original client and environment
	originalClient := editClient
	originalRepo := repo
	defer func() {
		editClient = originalClient
		repo = originalRepo
		editMessages = []string{}
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	editClient = mockClient
	repo = "owner/repo"

	tests := []struct {
		name              string
		args              []string
		setupMessages     []string
		expectedMessage   string
		shouldCallEdit    bool
	}{
		{
			name:            "single line positional",
			args:            []string{"123456", "Simple message"},
			expectedMessage: "Simple message",
			shouldCallEdit:  true,
		},
		{
			name:            "multi-line positional",
			args:            []string{"123456", "Line 1\nLine 2\nLine 3"},
			expectedMessage: "Line 1\nLine 2\nLine 3",
			shouldCallEdit:  true,
		},
		{
			name:            "single --message flag",
			args:            []string{"123456"},
			setupMessages:   []string{"Flag message"},
			expectedMessage: "Flag message",
			shouldCallEdit:  true,
		},
		{
			name:            "multiple --message flags joined with newlines",
			args:            []string{"123456"},
			setupMessages:   []string{"First line", "Second line", "Third line"},
			expectedMessage: "First line\nSecond line\nThird line",
			shouldCallEdit:  true,
		},
		{
			name:            "empty message in flag",
			args:            []string{"123456"},
			setupMessages:   []string{""},
			expectedMessage: "",
			shouldCallEdit:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset messages
			editMessages = tt.setupMessages

			// Run the command
			err := runEdit(nil, tt.args)
			assert.NoError(t, err)

			// In a real test with a spy client, we could verify the exact message
			// For now, we just verify no error occurred
		})
	}
}

func TestEditRepositoryParsing(t *testing.T) {
	// Save original values
	originalClient := editClient
	originalRepo := repo
	defer func() {
		editClient = originalClient
		repo = originalRepo
	}()

	mockClient := github.NewMockClient()
	editClient = mockClient

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
		// Note: Testing empty repository requires external gh CLI calls,
		// which are better tested in integration tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo = tt.setupRepo

			err := runEdit(nil, []string{"123456", "Test message"})
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