package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestRunAddReviewWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := addReviewClient
	originalRepo := repo
	originalPR := prNumber
	defer func() {
		addReviewClient = originalClient
		repo = originalRepo
		prNumber = originalPR
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	addReviewClient = mockClient
	repo = "owner/repo"
	prNumber = 123

	tests := []struct {
		name            string
		args            []string
		setupComments   []string
		setupBody       string
		setupEvent      string
		wantErr         bool
		expectedErrMsg  string
		resetGlobals    bool
	}{
		{
			name:          "create review with PR and body specified",
			args:          []string{"123", "Overall looks good"},
			setupComments: []string{"main.go:42:Great work here"},
			wantErr:       false,
		},
		{
			name:          "create review with PR only",
			args:          []string{"123"},
			setupComments: []string{"main.go:42:Great work here"},
			wantErr:       false,
		},
		{
			name:          "create review with body only (auto-detect PR)",
			args:          []string{"Overall looks good"},
			setupComments: []string{"main.go:42:Great work here"},
			wantErr:       false,
		},
		{
			name:          "create review with no args (auto-detect PR)",
			args:          []string{},
			setupComments: []string{"main.go:42:Great work here"},
			wantErr:       false,
		},
		{
			name:          "create review with --body flag",
			args:          []string{"123"},
			setupComments: []string{"main.go:42:Great work here"},
			setupBody:     "Review body from flag",
			wantErr:       false,
		},
		{
			name:          "create review with event",
			args:          []string{"123", "LGTM!"},
			setupComments: []string{"main.go:42:Excellent"},
			setupEvent:    "APPROVE",
			wantErr:       false,
		},
		{
			name:          "create review with range comment",
			args:          []string{"123", "Review feedback"},
			setupComments: []string{"main.go:40:45:This whole block is good"},
			wantErr:       false,
		},
		{
			name:          "create review with multiple comments",
			args:          []string{"123", "Multi-comment review"},
			setupComments: []string{
				"main.go:42:Single line comment",
				"test.go:10:15:Range comment",
				"docs.md:5:Documentation looks good",
			},
			wantErr: false,
		},
		{
			name:           "invalid PR number",
			args:           []string{"invalid", "Review body"},
			setupComments:  []string{"main.go:42:Comment"},
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "no comments provided",
			args:           []string{"123", "Review body"},
			setupComments:  []string{}, // Empty comments
			wantErr:        true,
			expectedErrMsg: "must provide at least one --comment",
		},
		{
			name:           "invalid comment format",
			args:           []string{"123", "Review body"},
			setupComments:  []string{"invalid-format"},
			wantErr:        true,
			expectedErrMsg: "format should be 'file:line:message'",
		},
		{
			name:           "invalid line number in comment",
			args:           []string{"123", "Review body"},
			setupComments:  []string{"main.go:invalid:message"},
			wantErr:        true,
			expectedErrMsg: "invalid line number",
		},
		{
			name:           "invalid range in comment",
			args:           []string{"123", "Review body"},
			setupComments:  []string{"main.go:50:40:message"}, // start > end
			wantErr:        true,
			expectedErrMsg: "start line (50) cannot be greater than end line (40)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global variables
			reviewComments = tt.setupComments
			reviewBody = tt.setupBody
			reviewEvent = tt.setupEvent
			noExpandSuggestionsReview = false

			// Reset globals at the end if requested
			if tt.resetGlobals {
				defer func() {
					reviewComments = []string{}
					reviewBody = ""
					reviewEvent = ""
				}()
			}

			err := runAddReview(nil, tt.args)
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

func TestParseCommentSpec(t *testing.T) {
	commitSHA := "abc123def456"

	tests := []struct {
		name           string
		spec           string
		wantPath       string
		wantLine       int
		wantStartLine  int
		wantSide       string
		wantBody       string
		wantErr        bool
		expectedErr    string
	}{
		{
			name:     "single line comment",
			spec:     "main.go:42:This is a comment",
			wantPath: "main.go",
			wantLine: 42,
			wantBody: "This is a comment",
			wantErr:  false,
		},
		{
			name:          "range comment",
			spec:          "test.go:10:15:This is a range comment",
			wantPath:      "test.go",
			wantLine:      15,
			wantStartLine: 10,
			wantSide:      "RIGHT",
			wantBody:      "This is a range comment",
			wantErr:       false,
		},
		{
			name:     "comment with colons in message",
			spec:     "config.yaml:5:Fix this: use https://example.com",
			wantPath: "config.yaml",
			wantLine: 5,
			wantBody: "Fix this: use https://example.com",
			wantErr:  false,
		},
		{
			name:        "too few parts",
			spec:        "main.go:42",
			wantErr:     true,
			expectedErr: "format should be 'file:line:message'",
		},
		{
			name:        "invalid line number",
			spec:        "main.go:invalid:message",
			wantErr:     true,
			expectedErr: "invalid line number",
		},
		{
			name:        "invalid start line in range",
			spec:        "main.go:invalid:42:message",
			wantErr:     true,
			expectedErr: "invalid line number", // Will parse as single line and fail
		},
		{
			name:        "invalid end line in range",
			spec:        "main.go:10:invalid:message",
			wantErr:     false, // Will parse as single line with "invalid:message" as body
			wantPath:    "main.go",
			wantLine:    10,
			wantBody:    "invalid:message",
		},
		{
			name:        "start line greater than end line",
			spec:        "main.go:50:40:message",
			wantErr:     true,
			expectedErr: "start line (50) cannot be greater than end line (40)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset suggestion expansion setting
			noExpandSuggestionsReview = false

			result, err := parseCommentSpec(tt.spec, commitSHA)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPath, result.Path)
				assert.Equal(t, tt.wantLine, result.Line)
				assert.Equal(t, tt.wantStartLine, result.StartLine)
				assert.Equal(t, tt.wantSide, result.Side)
				assert.Equal(t, tt.wantBody, result.Body)
				assert.Equal(t, commitSHA, result.CommitID)
			}
		})
	}
}

func TestParseCommentSpecWithSuggestionExpansion(t *testing.T) {
	commitSHA := "abc123def456"

	tests := []struct {
		name                      string
		spec                      string
		noExpandSuggestionsReview bool
		wantBody                  string
	}{
		{
			name:                      "expand suggestions enabled",
			spec:                      "main.go:42:[SUGGEST: fixed code]",
			noExpandSuggestionsReview: false,
			wantBody:                  "\n\n```suggestion\nfixed code\n```\n\n",
		},
		{
			name:                      "expand suggestions disabled",
			spec:                      "main.go:42:[SUGGEST: fixed code]",
			noExpandSuggestionsReview: true,
			wantBody:                  "[SUGGEST: fixed code]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global flag
			noExpandSuggestionsReview = tt.noExpandSuggestionsReview

			result, err := parseCommentSpec(tt.spec, commitSHA)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, result.Body)
		})
	}
}

func TestCreateReviewWithComments(t *testing.T) {
	mockClient := github.NewMockClient()

	tests := []struct {
		name         string
		owner        string
		repo         string
		pr           int
		body         string
		event        string
		commentSpecs []string
		wantErr      bool
		expectedErr  string
	}{
		{
			name:         "create pending review",
			owner:        "owner",
			repo:         "repo",
			pr:           123,
			body:         "Review body",
			event:        "",
			commentSpecs: []string{"main.go:42:Good work"},
			wantErr:      false,
		},
		{
			name:         "create approved review",
			owner:        "owner",
			repo:         "repo",
			pr:           123,
			body:         "LGTM",
			event:        "APPROVE",
			commentSpecs: []string{"main.go:42:Excellent"},
			wantErr:      false,
		},
		{
			name:         "create review with multiple comments",
			owner:        "owner",
			repo:         "repo",
			pr:           123,
			body:         "Mixed feedback",
			event:        "COMMENT",
			commentSpecs: []string{
				"main.go:42:Good",
				"test.go:10:15:Range comment",
			},
			wantErr: false,
		},
		{
			name:         "invalid comment spec",
			owner:        "owner",
			repo:         "repo",
			pr:           123,
			body:         "Review",
			event:        "",
			commentSpecs: []string{"invalid-spec"},
			wantErr:      true,
			expectedErr:  "invalid comment spec 'invalid-spec'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset suggestion expansion
			noExpandSuggestionsReview = false

			err := createReviewWithComments(mockClient, tt.owner, tt.repo, tt.pr, tt.body, tt.event, tt.commentSpecs)
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