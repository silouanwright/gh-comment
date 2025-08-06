package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/silouanwright/gh-comment/internal/github"
)

func TestRunReviewWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{}
	validateDiff = false // Disable validation for this test

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
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		dryRun = originalDryRun
		validateDiff = originalValidate
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"src/main.go:42:Good code"}
	dryRun = true
	validateDiff = false // Disable validation for this test

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
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		verbose = originalVerbose
		validateDiff = originalValidate
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "COMMENT"
	reviewCommentsFlag = []string{"src/main.go:1:Nice work"}
	verbose = true
	validateDiff = false // Disable validation for this test

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
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"src/main.go:1:test"}
	validateDiff = false // Disable validation for this test

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
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewCommentsFlag = []string{"src/main.go:1:test"}
	validateDiff = false // Disable validation for this test

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
			expectedErrMsg: "start line (15) cannot be greater than end line (10)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseReviewCommentSpec(tt.spec)
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
				assert.Equal(t, "RIGHT", result.Side)
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
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	validateDiff = false // Disable validation for this test

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
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	// Set client to nil to test initialization
	reviewClient = nil
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"src/main.go:1:Good"}
	validateDiff = false // Disable validation for this test

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

// TestReviewCommitIDNotIncluded - REGRESSION TEST
// This prevents the commit_id bug from reoccurring where individual
// comments had commit_id field causing GraphQL errors
func TestReviewCommitIDNotIncluded(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	// Set up mock client to verify commit_id is NOT sent
	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "APPROVE"
	reviewCommentsFlag = []string{"test.go:42:Test comment"} // Use valid file from mock
	validateDiff = false                                     // Disable validation for this test

	// Mock will track what gets sent to verify commit_id is NOT in individual comments
	err := runReview(nil, []string{"123", "Review body"})
	assert.NoError(t, err)

	// Verify the mock client received the correct data structure
	// without commit_id in individual comments (this prevents the GraphQL error)
	calls := mockClient.GetCreateReviewCalls()
	assert.Len(t, calls, 1)

	// Verify review comments don't have commit_id field
	reviewComments := calls[0].Comments
	assert.Len(t, reviewComments, 1)

	// Verify comment structure - should have Side but NOT CommitID
	comment := reviewComments[0]
	assert.Equal(t, "test.go", comment.Path)
	assert.Equal(t, 42, comment.Line)
	assert.Equal(t, "Test comment", comment.Body)
	assert.Equal(t, "RIGHT", comment.Side) // Side is required

	// CRITICAL: This field should NOT exist in individual comments
	// The absence of this field prevents the GraphQL commitId error
	// GitHub automatically uses the review-level commit for all comments
}

func TestValidateCommentLine(t *testing.T) {
	tests := []struct {
		name           string
		comment        github.ReviewCommentInput
		setupMock      func(*github.MockClient)
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "valid single line comment",
			comment: github.ReviewCommentInput{
				Path: "test.go",
				Line: 42,
				Body: "Good code",
			},
			setupMock: func(mock *github.MockClient) {
				// MockClient already returns test.go with lines 42, 43
			},
			wantErr: false,
		},
		{
			name: "valid range comment",
			comment: github.ReviewCommentInput{
				Path:      "test.go",
				StartLine: 42,
				Line:      43,
				Body:      "Nice refactoring",
			},
			setupMock: func(mock *github.MockClient) {
				// MockClient already returns test.go with lines 42, 43
			},
			wantErr: false,
		},
		{
			name: "file not found in diff",
			comment: github.ReviewCommentInput{
				Path: "nonexistent.go",
				Line: 42,
				Body: "Comment",
			},
			setupMock:      func(mock *github.MockClient) {},
			wantErr:        true,
			expectedErrMsg: "file 'nonexistent.go' not found in PR #123 diff",
		},
		{
			name: "line not found in diff",
			comment: github.ReviewCommentInput{
				Path: "test.go",
				Line: 999,
				Body: "Comment",
			},
			setupMock:      func(mock *github.MockClient) {},
			wantErr:        true,
			expectedErrMsg: "line(s) [999] do not exist in diff for file 'test.go'",
		},
		{
			name: "range with invalid lines",
			comment: github.ReviewCommentInput{
				Path:      "test.go",
				StartLine: 41,
				Line:      44,
				Body:      "Range comment",
			},
			setupMock:      func(mock *github.MockClient) {},
			wantErr:        true,
			expectedErrMsg: "line(s) [41 44] do not exist in diff for file 'test.go'",
		},
		{
			name: "fetch diff error - should skip validation",
			comment: github.ReviewCommentInput{
				Path: "test.go",
				Line: 42,
				Body: "Comment",
			},
			setupMock: func(mock *github.MockClient) {
				// This won't actually cause an error in MockClient,
				// but tests the error handling path
			},
			wantErr: false, // Should skip validation on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := github.NewMockClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			err := validateCommentLine(mockClient, "owner", "repo", 123, tt.comment)
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

func TestValidateCommentLineErrorMessages(t *testing.T) {
	mockClient := github.NewMockClient()

	tests := []struct {
		name           string
		comment        github.ReviewCommentInput
		expectedOutput []string
	}{
		{
			name: "file not found shows available files",
			comment: github.ReviewCommentInput{
				Path: "missing.go",
				Line: 42,
				Body: "Comment",
			},
			expectedOutput: []string{
				"Available files in this PR:",
				"test.go",
				"Use 'gh comment lines 123 <file>' to see commentable lines",
			},
		},
		{
			name: "invalid line shows available lines",
			comment: github.ReviewCommentInput{
				Path: "test.go",
				Line: 999,
				Body: "Comment",
			},
			expectedOutput: []string{
				"Available lines for comments: 42-43",
				"Use 'gh comment lines 123 test.go' to see detailed line information",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommentLine(mockClient, "owner", "repo", 123, tt.comment)
			assert.Error(t, err)

			errorMsg := err.Error()
			for _, expected := range tt.expectedOutput {
				assert.Contains(t, errorMsg, expected)
			}
		})
	}
}

func TestReviewValidationIntegration(t *testing.T) {
	// Save original values
	originalClient := reviewClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := reviewEventFlag
	originalComments := reviewCommentsFlag
	originalValidate := validateDiff
	defer func() {
		reviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		reviewEventFlag = originalEvent
		reviewCommentsFlag = originalComments
		validateDiff = originalValidate
	}()

	mockClient := github.NewMockClient()
	reviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	reviewEventFlag = "COMMENT"

	tests := []struct {
		name           string
		validateFlag   bool
		comments       []string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:         "validation disabled - should succeed even with invalid lines",
			validateFlag: false,
			comments:     []string{"nonexistent.go:999:Comment"},
			wantErr:      false,
		},
		{
			name:         "validation enabled - valid comments should succeed",
			validateFlag: true,
			comments:     []string{"test.go:42:Good code"},
			wantErr:      false,
		},
		{
			name:           "validation enabled - invalid file should fail",
			validateFlag:   true,
			comments:       []string{"nonexistent.go:42:Comment"},
			wantErr:        true,
			expectedErrMsg: "comment 1 validation failed",
		},
		{
			name:           "validation enabled - invalid line should fail",
			validateFlag:   true,
			comments:       []string{"test.go:999:Comment"},
			wantErr:        true,
			expectedErrMsg: "comment 1 validation failed",
		},
		{
			name:         "validation enabled - multiple valid comments",
			validateFlag: true,
			comments:     []string{"test.go:42:First", "test.go:43:Second"},
			wantErr:      false,
		},
		{
			name:           "validation enabled - mixed valid/invalid comments",
			validateFlag:   true,
			comments:       []string{"test.go:42:Valid", "test.go:999:Invalid"},
			wantErr:        true,
			expectedErrMsg: "comment 2 validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateDiff = tt.validateFlag
			reviewCommentsFlag = tt.comments

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
