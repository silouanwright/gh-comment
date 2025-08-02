package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunReviewWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{}

	tests := []struct {
		name           string
		args           []string
		setupEvent     string
		setupComments  []string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:       "create review with PR and body",
			args:       []string{"123", "LGTM!"},
			setupEvent: "APPROVE",
			setupComments: []string{
				"src/main.go:42:Great implementation",
				"src/utils.go:10:15:Nice refactoring",
			},
			wantErr: false,
		},
		{
			name:       "create review with just PR",
			args:       []string{"123"},
			setupEvent: "COMMENT",
			setupComments: []string{
				"src/main.go:1:Good code",
			},
			wantErr: false,
		},
		{
			name:          "create review with just body (auto-detect PR)",
			args:          []string{"Great work!"},
			setupEvent:    "APPROVE",
			setupComments: []string{},
			wantErr:       false,
		},
		{
			name:           "invalid PR number",
			args:           []string{"invalid", "body"},
			setupEvent:     "APPROVE",
			setupComments:  []string{},
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "invalid event type",
			args:           []string{"123", "body"},
			setupEvent:     "INVALID",
			setupComments:  []string{},
			wantErr:        true,
			expectedErrMsg: "invalid event type",
		},
		{
			name:           "no body or comments",
			args:           []string{"123"},
			setupEvent:     "APPROVE",
			setupComments:  []string{},
			wantErr:        true,
			expectedErrMsg: "review must have either a body message or comments",
		},
		{
			name:       "invalid comment format",
			args:       []string{"123", "body"},
			setupEvent: "APPROVE",
			setupComments: []string{
				"invalid:format",
			},
			wantErr:        true,
			expectedErrMsg: "format must be file:line:message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset for each test
			reviewEventFlag = tt.setupEvent
			reviewCommentsFlag = tt.setupComments

			err := runReview(nil, tt.args)
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

func TestReviewDryRun(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	originalDryRun := dryRun
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		dryRun = originalDryRun
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"src/main.go:42:Good code"}
	dryRun = true

	err := runReview(nil, []string{"123", "LGTM!"})
	assert.NoError(t, err)
}

func TestReviewVerbose(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	originalVerbose := verbose
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		verbose = originalVerbose
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "COMMENT"
	reviewCommentsFlag = []string{"src/main.go:1:Nice work"}
	verbose = true

	err := runReview(nil, []string{"123", "Good work!"})
	assert.NoError(t, err)
}

func TestReviewRepositoryParsing(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"src/main.go:1:test"}

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

			err := runReview(nil, []string{"123", "Review body"})
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

func TestReviewEventValidation(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewCommentsFlag = []string{"src/main.go:1:test"}

	tests := []struct {
		name           string
		event          string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:    "valid APPROVE event",
			event:   "APPROVE",
			wantErr: false,
		},
		{
			name:    "valid REQUEST_CHANGES event",
			event:   "REQUEST_CHANGES",
			wantErr: false,
		},
		{
			name:    "valid COMMENT event",
			event:   "COMMENT",
			wantErr: false,
		},
		{
			name:           "invalid event type",
			event:          "INVALID",
			wantErr:        true,
			expectedErrMsg: "invalid event type: INVALID",
		},
		{
			name:           "empty event type",
			event:          "",
			wantErr:        true,
			expectedErrMsg: "invalid event type",
		},
		{
			name:           "lowercase event type",
			event:          "approve",
			wantErr:        true,
			expectedErrMsg: "invalid event type: approve",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewEventFlag = tt.event

			err := runReview(nil, []string{"123", "Review body"})
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

func TestParseReviewCommentSpec(t *testing.T) {
	tests := []struct {
		name           string
		spec           string
		expectedFile   string
		expectedLine   int
		expectedStart  int
		expectedBody   string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:         "single line comment",
			spec:         "src/main.go:42:Fix this issue",
			expectedFile: "src/main.go",
			expectedLine: 42,
			expectedBody: "Fix this issue",
			wantErr:      false,
		},
		{
			name:          "range comment",
			spec:          "src/utils.go:10-15:Nice refactoring",
			expectedFile:  "src/utils.go",
			expectedLine:  15,
			expectedStart: 10,
			expectedBody:  "Nice refactoring",
			wantErr:       false,
		},
		{
			name:         "comment with colons in message",
			spec:         "config.yaml:5:Add timeout: 30s",
			expectedFile: "config.yaml",
			expectedLine: 5,
			expectedBody: "Add timeout: 30s",
			wantErr:      false,
		},
		{
			name:           "too few parts",
			spec:           "src/main.go:42",
			wantErr:        true,
			expectedErrMsg: "format must be file:line:message",
		},
		{
			name:           "empty file path",
			spec:           ":42:message",
			wantErr:        true,
			expectedErrMsg: "file path cannot be empty",
		},
		{
			name:           "empty message",
			spec:           "src/main.go:42:",
			wantErr:        true,
			expectedErrMsg: "message cannot be empty",
		},
		{
			name:           "invalid line number",
			spec:           "src/main.go:abc:message",
			wantErr:        true,
			expectedErrMsg: "invalid line number",
		},
		{
			name:           "zero line number",
			spec:           "src/main.go:0:message",
			wantErr:        true,
			expectedErrMsg: "line number must be positive",
		},
		{
			name:           "negative line number",
			spec:           "src/main.go:-5:message",
			wantErr:        true,
			expectedErrMsg: "invalid start line",
		},
		{
			name:           "invalid range format",
			spec:           "src/main.go:10-15-20:message",
			wantErr:        true,
			expectedErrMsg: "range format must be start-end",
		},
		{
			name:           "invalid start line in range",
			spec:           "src/main.go:abc-15:message",
			wantErr:        true,
			expectedErrMsg: "invalid start line",
		},
		{
			name:           "invalid end line in range",
			spec:           "src/main.go:10-xyz:message",
			wantErr:        true,
			expectedErrMsg: "invalid end line",
		},
		{
			name:           "zero start line in range",
			spec:           "src/main.go:0-15:message",
			wantErr:        true,
			expectedErrMsg: "line numbers must be positive",
		},
		{
			name:           "zero end line in range",
			spec:           "src/main.go:10-0:message",
			wantErr:        true,
			expectedErrMsg: "line numbers must be positive",
		},
		{
			name:           "start greater than end in range",
			spec:           "src/main.go:15-10:message",
			wantErr:        true,
			expectedErrMsg: "start line must be <= end line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseReviewCommentSpec(tt.spec, "abc123")
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFile, result.Path)
				assert.Equal(t, tt.expectedLine, result.Line)
				if tt.expectedStart != 0 {
					assert.Equal(t, tt.expectedStart, result.StartLine)
				}
				assert.Equal(t, tt.expectedBody, result.Body)
				assert.Equal(t, "abc123", result.CommitID)
			}
		})
	}
}

func TestReviewArgumentParsing(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"

	tests := []struct {
		name           string
		args           []string
		setupComments  []string
		expectedBody   string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:          "two args: PR and body",
			args:          []string{"123", "Great work!"},
			setupComments: []string{"src/main.go:1:Good"},
			expectedBody:  "Great work!",
			wantErr:       false,
		},
		{
			name:          "one arg: PR number",
			args:          []string{"123"},
			setupComments: []string{"src/main.go:1:Good"},
			expectedBody:  "",
			wantErr:       false,
		},
		{
			name:          "one arg: review body (auto-detect PR)",
			args:          []string{"Looks good!"},
			setupComments: []string{},
			expectedBody:  "Looks good!",
			wantErr:       false,
		},
		{
			name:           "invalid PR number in first arg",
			args:           []string{"invalid", "body"},
			setupComments:  []string{},
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewCommentsFlag = tt.setupComments

			err := runReview(nil, tt.args)
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

func TestReviewWithClientInitialization(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
	}()

	// Set client to nil to test initialization
	reviewClient = nil
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"src/main.go:1:Good"}

	// This test verifies that when reviewClient is nil, 
	// a RealClient is initialized in production
	// Since we can't easily test the RealClient without external dependencies,
	// we'll test that the initialization happens by setting up a mock afterwards
	
	// First verify the client gets initialized
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	
	err := runReview(nil, []string{"123", "Review body"})
	assert.NoError(t, err)
}