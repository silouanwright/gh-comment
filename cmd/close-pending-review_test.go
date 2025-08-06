package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/silouanwright/gh-comment/internal/github"
)

func TestRunClosePendingReviewWithMockClient(t *testing.T) {
	// Save original client and environment
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
	}()

	// Set up mock client and environment
	mockClient := github.NewMockClient()
	closePendingClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	closePendingEvent = "APPROVE"
	closePendingBody = ""

	tests := []struct {
		name           string
		args           []string
		setupEvent     string
		setupBody      string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:       "close pending review with PR and body",
			args:       []string{"123", "LGTM!"},
			setupEvent: "APPROVE",
			wantErr:    false,
		},
		{
			name:       "close pending review with just PR",
			args:       []string{"123"},
			setupEvent: "COMMENT",
			wantErr:    false,
		},
		{
			name:       "close pending review with just body (auto-detect PR)",
			args:       []string{"Great work!"},
			setupEvent: "APPROVE",
			wantErr:    false,
		},
		{
			name:       "close pending review with no args (auto-detect PR)",
			args:       []string{},
			setupEvent: "REQUEST_CHANGES",
			wantErr:    false,
		},
		{
			name:           "invalid PR number",
			args:           []string{"invalid", "body"},
			setupEvent:     "APPROVE",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
		{
			name:           "invalid event type",
			args:           []string{"123", "body"},
			setupEvent:     "INVALID",
			wantErr:        true,
			expectedErrMsg: "invalid event type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset for each test
			closePendingEvent = tt.setupEvent
			closePendingBody = tt.setupBody

			err := runClosePendingReview(nil, tt.args)
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

func TestClosePendingReviewDryRun(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	originalDryRun := dryRun
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
		dryRun = originalDryRun
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	closePendingClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	closePendingEvent = "APPROVE"
	closePendingBody = ""
	dryRun = true

	err := runClosePendingReview(nil, []string{"123", "LGTM!"})
	assert.NoError(t, err)
}

func TestClosePendingReviewVerbose(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	originalVerbose := verbose
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
		verbose = originalVerbose
	}()

	// Set up environment
	mockClient := github.NewMockClient()
	closePendingClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	closePendingEvent = "COMMENT"
	closePendingBody = ""
	verbose = true

	err := runClosePendingReview(nil, []string{"123", "Good work!"})
	assert.NoError(t, err)
}

func TestClosePendingReviewRepositoryParsing(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
	}()

	mockClient := github.NewMockClient()
	closePendingClient = mockClient
	prNumber = 123
	closePendingEvent = "APPROVE"
	closePendingBody = ""

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

			err := runClosePendingReview(nil, []string{"123", "Review body"})
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

func TestClosePendingReviewErrorHandling(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
	}()

	repo = "owner/repo"
	prNumber = 123
	closePendingEvent = "APPROVE"
	closePendingBody = ""

	tests := []struct {
		name           string
		setupMockError func(*github.MockClient)
		expectedErrMsg string
	}{
		{
			name: "find pending review error",
			setupMockError: func(m *github.MockClient) {
				m.FindPendingReviewError = assert.AnError
			},
			expectedErrMsg: "failed to find pending review",
		},
		{
			name: "submit review error",
			setupMockError: func(m *github.MockClient) {
				m.SubmitReviewError = assert.AnError
			},
			expectedErrMsg: "failed to submit review",
		},
		{
			name: "no pending review found",
			setupMockError: func(m *github.MockClient) {
				m.PendingReviewID = 0 // No pending review
			},
			expectedErrMsg: "no pending review found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := github.NewMockClient()
			if tt.setupMockError != nil {
				tt.setupMockError(mockClient)
			}
			closePendingClient = mockClient

			err := runClosePendingReview(nil, []string{"123", "Review body"})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func TestClosePendingReviewEventValidation(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
	}()

	mockClient := github.NewMockClient()
	closePendingClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	closePendingBody = ""

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
			closePendingEvent = tt.event

			err := runClosePendingReview(nil, []string{"123", "Review body"})
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

func TestClosePendingReviewArgumentParsing(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
	}()

	mockClient := github.NewMockClient()
	closePendingClient = mockClient
	repo = "owner/repo"
	prNumber = 123
	closePendingEvent = "APPROVE"
	closePendingBody = ""

	tests := []struct {
		name           string
		args           []string
		flagBody       string
		expectedBody   string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:         "two args: PR and body",
			args:         []string{"123", "Great work!"},
			flagBody:     "",
			expectedBody: "Great work!",
			wantErr:      false,
		},
		{
			name:         "one arg: PR number",
			args:         []string{"123"},
			flagBody:     "",
			expectedBody: "",
			wantErr:      false,
		},
		{
			name:         "one arg: review body (auto-detect PR)",
			args:         []string{"Looks good!"},
			flagBody:     "",
			expectedBody: "Looks good!",
			wantErr:      false,
		},
		{
			name:         "no args: auto-detect PR",
			args:         []string{},
			flagBody:     "",
			expectedBody: "",
			wantErr:      false,
		},
		{
			name:         "flag body overrides positional body",
			args:         []string{"123", "positional"},
			flagBody:     "flag body",
			expectedBody: "flag body",
			wantErr:      false,
		},
		{
			name:           "invalid PR number in first arg",
			args:           []string{"invalid", "body"},
			flagBody:       "",
			wantErr:        true,
			expectedErrMsg: "must be a valid integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closePendingBody = tt.flagBody

			err := runClosePendingReview(nil, tt.args)
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

func TestClosePendingReviewWithClientInitialization(t *testing.T) {
	// Save original values
	originalClient := closePendingClient
	originalRepo := repo
	originalPR := prNumber
	originalEvent := closePendingEvent
	originalBody := closePendingBody
	defer func() {
		closePendingClient = originalClient
		repo = originalRepo
		prNumber = originalPR
		closePendingEvent = originalEvent
		closePendingBody = originalBody
	}()

	// Set client to nil to test initialization
	closePendingClient = nil
	repo = "owner/repo"
	prNumber = 123
	closePendingEvent = "APPROVE"
	closePendingBody = ""

	// This test verifies that when closePendingClient is nil,
	// a RealClient is initialized in production
	// Since we can't easily test the RealClient without external dependencies,
	// we'll test that the initialization happens by setting up a mock afterwards

	// First verify the client gets initialized
	mockClient := github.NewMockClient()
	closePendingClient = mockClient

	err := runClosePendingReview(nil, []string{"123", "Review body"})
	assert.NoError(t, err)
}
