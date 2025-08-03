package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestAddCommand(t *testing.T) {
	// Save original client
	originalClient := addClient
	defer func() { addClient = originalClient }()

	// Test cases for the new add command (general PR comments)
	tests := []struct {
		name         string
		args         []string
		setupMock    func(*github.MockClient)
		setupGlobals func()
		wantErr      bool
		expectedErr  string
	}{
		{
			name: "add general comment with PR number",
			args: []string{"123", "Great work on this PR!"},
			setupMock: func(m *github.MockClient) {
				m.CreatedComment = &github.Comment{ID: 12345}
			},
			setupGlobals: func() {
				repo = "test/repo"
			},
			wantErr: false,
		},
		{
			name: "add general comment with auto-detect PR",
			args: []string{"LGTM!"},
			setupMock: func(m *github.MockClient) {
				m.CreatedComment = &github.Comment{ID: 67890}
			},
			setupGlobals: func() {
				prNumber = 456
				repo = "test/repo"
			},
			wantErr: false,
		},
		{
			name:        "invalid PR number",
			args:        []string{"abc", "test comment"},
			wantErr:     true,
			expectedErr: "must be a valid integer",
		},
		{
			name:        "empty comment",
			args:        []string{"123", ""},
			wantErr:     true,
			expectedErr: "comment cannot be empty",
		},
		{
			name:        "too many arguments",
			args:        []string{"123", "file.js", "42", "comment"},
			wantErr:     true,
			expectedErr: "invalid arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock client
			mockClient := github.NewMockClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			addClient = mockClient

			// Setup globals
			if tt.setupGlobals != nil {
				tt.setupGlobals()
			}

			// Run command
			err := runAdd(nil, tt.args)

			// Check results
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddCommandWithMessages(t *testing.T) {
	// Save original state
	originalClient := addClient
	originalMessages := messages
	defer func() {
		addClient = originalClient
		messages = originalMessages
	}()

	// Setup mock
	mockClient := github.NewMockClient()
	mockClient.CreatedComment = &github.Comment{ID: 11111}
	addClient = mockClient

	// Setup globals
	repo = "test/repo"

	// Test multi-line with --message flags
	messages = []string{"First line", "Second line"}
	err := runAdd(nil, []string{"123"})
	assert.NoError(t, err)

	// Verify the comment was created with joined content
	// Note: This would need access to the mock's internal state
	// For now, just verify no error occurred
}
