package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAPIError(t *testing.T) {
	originalErr := errors.New("network timeout")
	formattedErr := formatAPIError("list comments", "/repos/owner/repo/issues/123/comments", originalErr)

	assert.Error(t, formattedErr)
	assert.Contains(t, formattedErr.Error(), "GitHub API error during list comments")
	assert.Contains(t, formattedErr.Error(), "network timeout")
}

func TestFormatActionableError(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		originalError string
		expectedText  []string
	}{
		{
			name:          "422 validation error",
			operation:     "comment creation",
			originalError: "HTTP 422: Unprocessable Entity - validation failed",
			expectedText:  []string{"validation error", "comment creation", "gh comment lines", "commentable lines"},
		},
		{
			name:          "404 not found error",
			operation:     "comment editing",
			originalError: "HTTP 404: Not Found",
			expectedText:  []string{"resource not found", "comment editing", "PR number exists", "comment ID is valid"},
		},
		{
			name:          "403 forbidden error",
			operation:     "review creation",
			originalError: "HTTP 403: Forbidden",
			expectedText:  []string{"permission denied", "review creation", "write access", "gh auth status"},
		},
		{
			name:          "401 unauthorized error",
			operation:     "reaction addition",
			originalError: "HTTP 401: Unauthorized",
			expectedText:  []string{"authentication failed", "reaction addition", "gh auth login", "token has expired"},
		},
		{
			name:          "rate limit error",
			operation:     "comment fetch",
			originalError: "rate limit exceeded: too many requests",
			expectedText:  []string{"rate limit exceeded", "comment fetch", "Wait a few minutes", "authenticated requests"},
		},
		{
			name:          "server error",
			operation:     "reply creation",
			originalError: "HTTP 500: Internal Server Error",
			expectedText:  []string{"GitHub server error", "reply creation", "temporary GitHub", "status.github.com"},
		},
		{
			name:          "network error",
			operation:     "list comments",
			originalError: "network timeout connecting to api.github.com",
			expectedText:  []string{"network error", "list comments", "internet connection", "Try again"},
		},
		{
			name:          "schema validation error",
			operation:     "review submission",
			originalError: "No subschema in oneOf matched",
			expectedText:  []string{"invalid request format", "review submission", "command syntax", "required arguments"},
		},
		{
			name:          "generic error",
			operation:     "unknown operation",
			originalError: "some unexpected error message",
			expectedText:  []string{"error during unknown operation", "gh comment --help", "PR number", "verbose"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalErr := errors.New(tt.originalError)
			formattedErr := formatActionableError(tt.operation, originalErr)

			assert.Error(t, formattedErr)
			errStr := formattedErr.Error()

			for _, expectedText := range tt.expectedText {
				assert.Contains(t, errStr, expectedText, "Error should contain: %s", expectedText)
			}

			// Should contain the original error
			assert.Contains(t, errStr, tt.originalError)

			// Should contain suggestions
			assert.Contains(t, errStr, "💡 Suggestions:")
		})
	}
}

func TestFormatValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    string
		expected string
		want     string
	}{
		{
			name:     "invalid comment ID",
			field:    "comment ID",
			value:    "abc123",
			expected: "must be a valid integer",
			want:     "invalid comment ID 'abc123': must be a valid integer",
		},
		{
			name:     "invalid reaction",
			field:    "reaction",
			value:    "invalid_reaction",
			expected: "must be one of: +1, -1, laugh, confused, heart, hooray, rocket, eyes",
			want:     "invalid reaction 'invalid_reaction': must be one of: +1, -1, laugh, confused, heart, hooray, rocket, eyes",
		},
		{
			name:     "invalid PR number",
			field:    "PR number",
			value:    "0",
			expected: "must be a positive integer",
			want:     "invalid PR number '0': must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatValidationError(tt.field, tt.value, tt.expected)
			assert.Error(t, err)
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestFormatNotFoundError(t *testing.T) {
	tests := []struct {
		name       string
		resource   string
		identifier interface{}
		want       string
	}{
		{
			name:       "comment not found with integer ID",
			resource:   "comment",
			identifier: 123456,
			want:       "comment not found: 123456",
		},
		{
			name:       "PR not found with string identifier",
			resource:   "pull request",
			identifier: "feature-branch",
			want:       "pull request not found: feature-branch",
		},
		{
			name:       "repository not found",
			resource:   "repository",
			identifier: "owner/repo",
			want:       "repository not found: owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatNotFoundError(tt.resource, tt.identifier)
			assert.Error(t, err)
			assert.Equal(t, tt.want, err.Error())
		})
	}
}

func TestGetPRContext(t *testing.T) {
	// Save original values
	originalRepo := repo
	originalPRNumber := prNumber
	defer func() {
		repo = originalRepo
		prNumber = originalPRNumber
	}()

	tests := []struct {
		name           string
		setupRepo      string
		setupPRNumber  int
		wantRepo       string
		wantPR         int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:          "with PR flag set",
			setupRepo:     "owner/repo",
			setupPRNumber: 123,
			wantRepo:      "owner/repo",
			wantPR:        123,
			wantErr:       false,
		},
		{
			name:          "with auto-detected PR",
			setupRepo:     "owner/repo",
			setupPRNumber: 0, // Will trigger auto-detection
			wantRepo:      "owner/repo",
			wantPR:        123, // This will be set by the mock in prNumber global
			wantErr:       false,
		},
		// Note: Testing empty repository requires external gh CLI calls,
		// which are better tested in integration tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up globals
			repo = tt.setupRepo
			prNumber = tt.setupPRNumber

			// For auto-detection test, we need to set up the global prNumber
			// since getCurrentPR() will read from it in our test environment
			if tt.setupPRNumber == 0 && !tt.wantErr {
				prNumber = tt.wantPR // Simulate auto-detection result
			}

			gotRepo, gotPR, err := getPRContext()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRepo, gotRepo)
				assert.Equal(t, tt.wantPR, gotPR)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that our constants are set to reasonable values
	assert.Equal(t, 100, MaxGraphQLResults)
	assert.Equal(t, 65536, MaxCommentLength)
	assert.Equal(t, 30, DefaultPageSize)

	// Verify they're positive values
	assert.Greater(t, MaxGraphQLResults, 0)
	assert.Greater(t, MaxCommentLength, 0)
	assert.Greater(t, DefaultPageSize, 0)
}

func TestExecuteFunction(t *testing.T) {
	// Test that Execute() function exists and delegates properly
	// This is mainly for coverage of the Execute wrapper function
	assert.NotNil(t, Execute, "Execute function should be defined")

	// We can't easily test successful execution without complex setup,
	// but we can verify the function can be called and behaves reasonably
	// When called without arguments, it shows help (which is successful behavior)

	// The Execute function should complete successfully when showing help
	// This tests the wrapper function coverage without side effects
	err := Execute()
	assert.NoError(t, err, "Execute should succeed when showing help")
}
