package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunAddWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := addClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		addClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	addClient = mockClient
	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name           string
		args           []string
		setupMessages  []string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:    "add single line comment with PR specified",
			args:    []string{"123", "main.go", "42", "This needs fixing"},
			wantErr: false,
		},
		{
			name:    "add range comment with PR specified",
			args:    []string{"123", "main.go", "42:45", "This whole block needs review"},
			wantErr: false,
		},
		{
			name:    "add comment with auto-detected PR",
			args:    []string{"main.go", "42", "Auto-detected PR comment"},
			wantErr: false,
		},
		{
			name:          "add comment with message flags",
			args:          []string{"main.go", "42"},
			setupMessages: []string{"First line", "Second line"},
			wantErr:       false,
		},
		{
			name:           "invalid PR number",
			args:           []string{"invalid", "main.go", "42", "Comment"},
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "invalid line number",
			args:           []string{"123", "main.go", "invalid", "Comment"},
			wantErr:        true,
			expectedErrMsg: "invalid line number",
		},
		{
			name:           "invalid line range",
			args:           []string{"123", "main.go", "45:42", "Comment"},
			wantErr:        true,
			expectedErrMsg: "start line (45) cannot be greater than end line (42)",
		},
		{
			name:           "invalid arguments",
			args:           []string{"only-one-arg"},
			wantErr:        true,
			expectedErrMsg: "invalid arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			messages = tt.setupMessages
			noExpandSuggestions = false

			err := runAdd(nil, tt.args)
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

func TestParseLineSpec(t *testing.T) {
	tests := []struct {
		name         string
		lineSpec     string
		wantStart    int
		wantEnd      int
		wantErr      bool
		expectedErr  string
	}{
		{
			name:      "single line",
			lineSpec:  "42",
			wantStart: 42,
			wantEnd:   42,
			wantErr:   false,
		},
		{
			name:      "valid range",
			lineSpec:  "42:45",
			wantStart: 42,
			wantEnd:   45,
			wantErr:   false,
		},
		{
			name:        "invalid single line",
			lineSpec:    "invalid",
			wantErr:     true,
			expectedErr: "invalid line number",
		},
		{
			name:        "invalid range format",
			lineSpec:    "42:45:50",
			wantErr:     true,
			expectedErr: "invalid line range format",
		},
		{
			name:        "invalid start line in range",
			lineSpec:    "invalid:45",
			wantErr:     true,
			expectedErr: "invalid start line",
		},
		{
			name:        "invalid end line in range",
			lineSpec:    "42:invalid",
			wantErr:     true,
			expectedErr: "invalid end line",
		},
		{
			name:        "start greater than end",
			lineSpec:    "45:42",
			wantErr:     true,
			expectedErr: "start line (45) cannot be greater than end line (42)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := parseLineSpec(tt.lineSpec)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStart, start)
				assert.Equal(t, tt.wantEnd, end)
			}
		})
	}
}