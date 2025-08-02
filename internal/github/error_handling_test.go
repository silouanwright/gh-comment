package github

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateRepoParams(t *testing.T) {
	tests := []struct {
		name        string
		owner       string
		repo        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid params",
			owner:       "octocat",
			repo:        "hello-world",
			expectError: false,
		},
		{
			name:        "empty owner",
			owner:       "",
			repo:        "hello-world",
			expectError: true,
			errorMsg:    "repository owner cannot be empty",
		},
		{
			name:        "empty repo",
			owner:       "octocat",
			repo:        "",
			expectError: true,
			errorMsg:    "repository name cannot be empty",
		},
		{
			name:        "owner with slash",
			owner:       "octo/cat",
			repo:        "hello-world",
			expectError: true,
			errorMsg:    "invalid repository format",
		},
		{
			name:        "repo with slash",
			owner:       "octocat",
			repo:        "hello/world",
			expectError: true,
			errorMsg:    "invalid repository format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRepoParams(tt.owner, tt.repo)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestIsValidReaction(t *testing.T) {
	tests := []struct {
		reaction string
		valid    bool
	}{
		{"+1", true},
		{"-1", true},
		{"laugh", true},
		{"hooray", true},
		{"confused", true},
		{"heart", true},
		{"rocket", true},
		{"eyes", true},
		{"invalid", false},
		{"thumbsup", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.reaction, func(t *testing.T) {
			result := isValidReaction(tt.reaction)
			if result != tt.valid {
				t.Errorf("isValidReaction(%q) = %v, want %v", tt.reaction, result, tt.valid)
			}
		})
	}
}

func TestWrapAPIError(t *testing.T) {
	client := &RealClient{}

	tests := []struct {
		name      string
		err       error
		operation string
		args      []interface{}
		wantMsg   string
	}{
		{
			name:      "rate limit error",
			err:       errors.New("rate limit exceeded"),
			operation: "fetch comments for PR #%d in %s/%s",
			args:      []interface{}{123, "owner", "repo"},
			wantMsg:   "rate limit exceeded while trying to fetch comments for PR #123 in owner/repo",
		},
		{
			name:      "403 error (likely rate limit)",
			err:       errors.New("HTTP 403 Forbidden"),
			operation: "create comment on PR #%d",
			args:      []interface{}{456},
			wantMsg:   "rate limit exceeded while trying to create comment on PR #456",
		},
		{
			name:      "404 error",
			err:       errors.New("HTTP 404 Not Found"),
			operation: "fetch PR #%d details",
			args:      []interface{}{789},
			wantMsg:   "resource not found while trying to fetch PR #789 details",
		},
		{
			name:      "401 error",
			err:       errors.New("HTTP 401 Unauthorized"),
			operation: "list comments",
			args:      []interface{}{},
			wantMsg:   "authentication failed while trying to list comments",
		},
		{
			name:      "422 error",
			err:       errors.New("HTTP 422 Unprocessable Entity"),
			operation: "create review comment",
			args:      []interface{}{},
			wantMsg:   "validation error while trying to create review comment",
		},
		{
			name:      "secondary rate limit",
			err:       errors.New("abuse detection triggered"),
			operation: "create multiple comments",
			args:      []interface{}{},
			wantMsg:   "secondary rate limit triggered while trying to create multiple comments",
		},
		{
			name:      "generic error",
			err:       errors.New("network error"),
			operation: "perform operation",
			args:      []interface{}{},
			wantMsg:   "GitHub API error while trying to perform operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.wrapAPIError(tt.err, tt.operation, tt.args...)

			if !strings.Contains(result.Error(), tt.wantMsg) {
				t.Errorf("wrapAPIError() error = %v, want to contain %v", result.Error(), tt.wantMsg)
			}

			// Verify the original error is wrapped
			if !errors.Is(result, tt.err) {
				t.Errorf("wrapAPIError() should wrap the original error")
			}

			// Check for helpful tips in specific error types
			if strings.Contains(tt.err.Error(), "rate limit") {
				if !strings.Contains(result.Error(), "gh api rate_limit") {
					t.Errorf("rate limit error should include tip about checking rate limit status")
				}
			}

			if strings.Contains(tt.err.Error(), "404") {
				if !strings.Contains(result.Error(), "Verify the repository exists") {
					t.Errorf("404 error should include tip about repository access")
				}
			}

			if strings.Contains(tt.err.Error(), "401") {
				if !strings.Contains(result.Error(), "gh auth status") {
					t.Errorf("401 error should include tip about checking authentication")
				}
			}
		})
	}
}

func TestErrorHandlingInMethods(t *testing.T) {
	// Test that our methods properly validate parameters
	client := &RealClient{}

	t.Run("CreateIssueComment with invalid params", func(t *testing.T) {
		tests := []struct {
			name     string
			owner    string
			repo     string
			prNumber int
			body     string
			wantErr  string
		}{
			{
				name:     "empty owner",
				owner:    "",
				repo:     "repo",
				prNumber: 123,
				body:     "test",
				wantErr:  "repository owner cannot be empty",
			},
			{
				name:     "empty repo",
				owner:    "owner",
				repo:     "",
				prNumber: 123,
				body:     "test",
				wantErr:  "repository name cannot be empty",
			},
			{
				name:     "invalid PR number",
				owner:    "owner",
				repo:     "repo",
				prNumber: -1,
				body:     "test",
				wantErr:  "invalid PR number -1: must be positive",
			},
			{
				name:     "empty body",
				owner:    "owner",
				repo:     "repo",
				prNumber: 123,
				body:     "   ",
				wantErr:  "comment body cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := client.CreateIssueComment(tt.owner, tt.repo, tt.prNumber, tt.body)
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.wantErr, err.Error())
				}
			})
		}
	})

	t.Run("AddReaction with invalid params", func(t *testing.T) {
		tests := []struct {
			name      string
			owner     string
			repo      string
			commentID int
			reaction  string
			wantErr   string
		}{
			{
				name:      "invalid comment ID",
				owner:     "owner",
				repo:      "repo",
				commentID: 0,
				reaction:  "+1",
				wantErr:   "invalid comment ID 0: must be positive",
			},
			{
				name:      "invalid reaction",
				owner:     "owner",
				repo:      "repo",
				commentID: 123,
				reaction:  "invalid",
				wantErr:   "invalid reaction 'invalid'",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := client.AddReaction(tt.owner, tt.repo, tt.commentID, 123, tt.reaction)
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.wantErr, err.Error())
				}
			})
		}
	})
}
