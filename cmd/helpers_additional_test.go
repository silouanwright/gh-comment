package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
)

func TestFormatActionableErrorNewPatterns(t *testing.T) {
	// Test the newly added error patterns that weren't covered in the original tests
	tests := []struct {
		name         string
		operation    string
		err          error
		wantContains []string
	}{
		{
			name:         "GraphQL API error",
			operation:    "create-comment",
			err:          fmt.Errorf("GraphQL: Field 'invalidField' doesn't exist on type 'Comment'"),
			wantContains: []string{"GraphQL API error", "GraphQL queries", "tool maintainers"},
		},
		{
			name:         "network timeout error",
			operation:    "fetch-comments",
			err:          fmt.Errorf("context deadline exceeded"),
			wantContains: []string{"network timeout", "stable network", "smaller chunks"},
		},
		{
			name:         "connection refused error",
			operation:    "api-call",
			err:          fmt.Errorf("connection refused"),
			wantContains: []string{"network connection error", "internet connection", "firewall"},
		},
		{
			name:         "abuse detection error",
			operation:    "bulk-comment",
			err:          fmt.Errorf("abuse detection mechanism triggered"),
			wantContains: []string{"GitHub abuse detection triggered", "1 minute", "Reduce concurrent operations"},
		},
		{
			name:         "archived repository error",
			operation:    "comment-create",
			err:          fmt.Errorf("repository archived and read-only"),
			wantContains: []string{"repository is archived", "cannot modify archived", "unarchive"},
		},
		{
			name:         "insufficient scope error",
			operation:    "create-review",
			err:          fmt.Errorf("token does not have required scope"),
			wantContains: []string{"insufficient token permissions", "broader scopes", "repo"},
		},
		{
			name:         "branch protection error",
			operation:    "submit-review",
			err:          fmt.Errorf("branch protection rules prevent this action"),
			wantContains: []string{"branch protection rules", "status checks", "administrator"},
		},
		{
			name:         "closed PR error",
			operation:    "add-comment",
			err:          fmt.Errorf("pull request closed"),
			wantContains: []string{"closed or locked", "cannot comment", "reopen"},
		},
		{
			name:         "duplicate review error",
			operation:    "submit-review",
			err:          fmt.Errorf("review already submitted"),
			wantContains: []string{"duplicate operation", "already exists", "edit operations"},
		},
		{
			name:         "Resource not accessible error",
			operation:    "create-comment",
			err:          fmt.Errorf("Resource not accessible by integration"),
			wantContains: []string{"permission denied", "token has the required scopes"},
		},
		{
			name:         "Bad credentials error",
			operation:    "auth-test",
			err:          fmt.Errorf("Bad credentials provided"),
			wantContains: []string{"authentication failed", "haven't been revoked"},
		},
		{
			name:         "Gateway timeout error",
			operation:    "large-fetch",
			err:          fmt.Errorf("504 Gateway Timeout"),
			wantContains: []string{"GitHub server error", "smaller batch sizes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatActionableError(tt.operation, tt.err)
			resultStr := result.Error()

			for _, want := range tt.wantContains {
				if !strings.Contains(resultStr, want) {
					t.Errorf("formatActionableError() result missing expected text %q in:\n%s", want, resultStr)
				}
			}

			// Ensure operation is mentioned
			if !strings.Contains(resultStr, tt.operation) {
				t.Errorf("formatActionableError() result should mention operation %q", tt.operation)
			}
		})
	}
}

func TestValidateReviewComments(t *testing.T) {
	originalClient := batchClient
	defer func() { batchClient = originalClient }()

	mockClient := github.NewMockClient()
	batchClient = mockClient

	tests := []struct {
		name    string
		config  *BatchConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid review comments",
			config: &BatchConfig{
				Review: &ReviewConfig{Body: "Good work!", Event: "APPROVE"},
				Comments: []CommentConfig{
					{File: "test.go", Line: 42, Message: "Good code", Type: "review"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid review body",
			config: &BatchConfig{
				Review: &ReviewConfig{Body: "", Event: "APPROVE"},
				Comments: []CommentConfig{
					{File: "test.go", Line: 42, Message: "Good code", Type: "review"},
				},
			},
			wantErr: false, // Empty body is allowed
		},
		{
			name: "empty comment message is allowed",
			config: &BatchConfig{
				Review: &ReviewConfig{Body: "Review", Event: "COMMENT"},
				Comments: []CommentConfig{
					{File: "test.go", Line: 42, Message: "", Type: "review"},
				},
			},
			wantErr: false, // Empty messages are allowed
		},
		{
			name: "mix of issue and review comments",
			config: &BatchConfig{
				Review: &ReviewConfig{Body: "Mixed review", Event: "COMMENT"},
				Comments: []CommentConfig{
					{File: "", Line: 0, Message: "General comment", Type: "issue"},
					{File: "test.go", Line: 42, Message: "Line comment", Type: "review"},
				},
			},
			wantErr: false,
		},
		{
			name: "comment with range",
			config: &BatchConfig{
				Review: &ReviewConfig{Body: "Range review", Event: "COMMENT"},
				Comments: []CommentConfig{
					{File: "test.go", Range: "42-43", Message: "Range comment", Type: "review"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateReviewComments(mockClient, "owner", "repo", 123, tt.config)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateReviewComments() expected error but got nil")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateReviewComments() error = %v, want to contain %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateReviewComments() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestBuildReviewInput(t *testing.T) {
	tests := []struct {
		name     string
		config   *BatchConfig
		comments []github.ReviewCommentInput
		want     github.ReviewInput
	}{
		{
			name: "complete review input",
			config: &BatchConfig{
				Review: &ReviewConfig{
					Body:  "Comprehensive review",
					Event: "REQUEST_CHANGES",
				},
			},
			comments: []github.ReviewCommentInput{
				{Body: "Fix this", Path: "test.go", Line: 10},
			},
			want: github.ReviewInput{
				Body:  "Comprehensive review",
				Event: "REQUEST_CHANGES",
				Comments: []github.ReviewCommentInput{
					{Body: "Fix this", Path: "test.go", Line: 10},
				},
			},
		},
		{
			name: "default event when empty",
			config: &BatchConfig{
				Review: &ReviewConfig{
					Body:  "Default event review",
					Event: "",
				},
			},
			comments: []github.ReviewCommentInput{},
			want: github.ReviewInput{
				Body:     "Default event review",
				Event:    "COMMENT",
				Comments: []github.ReviewCommentInput{},
			},
		},
		{
			name: "empty body allowed",
			config: &BatchConfig{
				Review: &ReviewConfig{
					Body:  "",
					Event: "APPROVE",
				},
			},
			comments: []github.ReviewCommentInput{
				{Body: "LGTM", Path: "main.go", Line: 5},
			},
			want: github.ReviewInput{
				Body:  "",
				Event: "APPROVE",
				Comments: []github.ReviewCommentInput{
					{Body: "LGTM", Path: "main.go", Line: 5},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildReviewInput(tt.config, tt.comments)

			if got.Body != tt.want.Body {
				t.Errorf("buildReviewInput() Body = %v, want %v", got.Body, tt.want.Body)
			}
			if got.Event != tt.want.Event {
				t.Errorf("buildReviewInput() Event = %v, want %v", got.Event, tt.want.Event)
			}
			if len(got.Comments) != len(tt.want.Comments) {
				t.Errorf("buildReviewInput() Comments length = %v, want %v", len(got.Comments), len(tt.want.Comments))
			}
		})
	}
}

func TestFormatListOutput(t *testing.T) {
	// Save original global state
	originalFormat := outputFormat
	originalIDsOnly := idsOnly
	defer func() {
		outputFormat = originalFormat
		idsOnly = originalIDsOnly
	}()

	comments := []Comment{
		{ID: 1, Author: "user1", Body: "First comment"},
		{ID: 2, Author: "user2", Body: "Second comment"},
	}

	tests := []struct {
		name        string
		setupFormat func()
		wantErr     bool
	}{
		{
			name: "default format",
			setupFormat: func() {
				outputFormat = "default"
				idsOnly = false
			},
			wantErr: false,
		},
		{
			name: "JSON format",
			setupFormat: func() {
				outputFormat = "json"
				idsOnly = false
			},
			wantErr: false,
		},
		{
			name: "IDs only format",
			setupFormat: func() {
				outputFormat = "default"
				idsOnly = true
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFormat()

			err := formatListOutput(comments, 123)

			if tt.wantErr {
				if err == nil {
					t.Errorf("formatListOutput() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("formatListOutput() unexpected error = %v", err)
				}
			}
		})
	}
}
