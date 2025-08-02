package github

import (
	"time"
)

// GitHubAPI defines the interface for GitHub API operations
type GitHubAPI interface {
	// Comment operations
	ListIssueComments(owner, repo string, prNumber int) ([]Comment, error)
	ListReviewComments(owner, repo string, prNumber int) ([]Comment, error)
	CreateIssueComment(owner, repo string, prNumber int, body string) (*Comment, error)
	CreateReviewCommentReply(owner, repo string, commentID int, body string) (*Comment, error)

	// Reaction operations
	AddReaction(owner, repo string, commentID int, reaction string) error
	RemoveReaction(owner, repo string, commentID int, reaction string) error

	// Comment operations
	EditComment(owner, repo string, commentID int, body string) error
	AddReviewComment(owner, repo string, pr int, comment ReviewCommentInput) error

	// PR operations
	FetchPRDiff(owner, repo string, pr int) (*PullRequestDiff, error)
	GetPRDetails(owner, repo string, pr int) (map[string]interface{}, error)

	// Review operations
	CreateReview(owner, repo string, pr int, review ReviewInput) error
	FindPendingReview(owner, repo string, pr int) (int, error)
	SubmitReview(owner, repo string, pr, reviewID int, body, event string) error

	// GraphQL operations
	ResolveReviewThread(threadID string) error
	FindReviewThreadForComment(owner, repo string, prNumber, commentID int) (string, error)
}

// Comment represents a GitHub comment (issue or review)
type Comment struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Review comment specific fields
	Path     string `json:"path,omitempty"`
	Line     int    `json:"line,omitempty"`
	Position int    `json:"position,omitempty"`

	// Computed fields
	Type string `json:"-"` // "issue" or "review"
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

// ReviewCommentInput represents input for creating a review comment
type ReviewCommentInput struct {
	Body      string `json:"body"`
	Path      string `json:"path"`
	Line      int    `json:"line,omitempty"`
	StartLine int    `json:"start_line,omitempty"`
	Side      string `json:"side,omitempty"`
	CommitID  string `json:"commit_id"`
}

// ReviewInput represents input for creating a review
type ReviewInput struct {
	Body     string               `json:"body,omitempty"`
	Event    string               `json:"event"`
	Comments []ReviewCommentInput `json:"comments,omitempty"`
}

// PullRequestDiff represents PR diff information
type PullRequestDiff struct {
	Files []DiffFile
}

// DiffFile represents a file in a PR diff
type DiffFile struct {
	Filename string
	Lines    map[int]bool // line numbers that exist in the diff
}

// MockClient implements GitHubAPI for testing
type MockClient struct {
	IssueComments  []Comment
	ReviewComments []Comment
	CreatedComment *Comment
	ResolvedThread string
	PendingReviewID int
	SubmittedReviewID int

	// Error simulation
	ListIssueCommentsError    error
	ListReviewCommentsError   error
	CreateCommentError        error
	ResolveThreadError        error
	FindReviewThreadError     error
	FindPendingReviewError    error
	SubmitReviewError         error
}

// NewMockClient creates a new mock client for testing
func NewMockClient() *MockClient {
	return &MockClient{
		IssueComments: []Comment{
			{
				ID:        123456,
				Body:      "LGTM! Great work on this PR.",
				Type:      "issue",
				User:      User{Login: "reviewer1"},
				CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		ReviewComments: []Comment{
			{
				ID:        654321,
				Body:      "Consider using a more descriptive variable name here.",
				Type:      "review",
				User:      User{Login: "reviewer2"},
				CreatedAt: time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
				Path:      "main.go",
				Line:      42,
			},
		},
		PendingReviewID: 987654, // Mock pending review ID
	}
}

// Mock implementations
func (m *MockClient) ListIssueComments(owner, repo string, prNumber int) ([]Comment, error) {
	if m.ListIssueCommentsError != nil {
		return nil, m.ListIssueCommentsError
	}
	return m.IssueComments, nil
}

func (m *MockClient) ListReviewComments(owner, repo string, prNumber int) ([]Comment, error) {
	if m.ListReviewCommentsError != nil {
		return nil, m.ListReviewCommentsError
	}
	return m.ReviewComments, nil
}

func (m *MockClient) CreateIssueComment(owner, repo string, prNumber int, body string) (*Comment, error) {
	if m.CreateCommentError != nil {
		return nil, m.CreateCommentError
	}

	comment := &Comment{
		ID:        789012,
		Body:      body,
		Type:      "issue",
		User:      User{Login: "testuser"},
		CreatedAt: time.Now(),
	}
	m.CreatedComment = comment
	return comment, nil
}

func (m *MockClient) CreateReviewCommentReply(owner, repo string, commentID int, body string) (*Comment, error) {
	if m.CreateCommentError != nil {
		return nil, m.CreateCommentError
	}

	comment := &Comment{
		ID:        345678,
		Body:      body,
		Type:      "review",
		User:      User{Login: "testuser"},
		CreatedAt: time.Now(),
	}
	m.CreatedComment = comment
	return comment, nil
}

func (m *MockClient) FindReviewThreadForComment(owner, repo string, prNumber, commentID int) (string, error) {
	if m.FindReviewThreadError != nil {
		return "", m.FindReviewThreadError
	}
	return "RT_123", nil
}

func (m *MockClient) ResolveReviewThread(threadID string) error {
	if m.ResolveThreadError != nil {
		return m.ResolveThreadError
	}
	m.ResolvedThread = threadID
	return nil
}

func (m *MockClient) AddReaction(owner, repo string, commentID int, reaction string) error {
	return nil
}

func (m *MockClient) RemoveReaction(owner, repo string, commentID int, reaction string) error {
	return nil
}

func (m *MockClient) EditComment(owner, repo string, commentID int, body string) error {
	return nil
}

func (m *MockClient) AddReviewComment(owner, repo string, pr int, comment ReviewCommentInput) error {
	return nil
}

func (m *MockClient) FetchPRDiff(owner, repo string, pr int) (*PullRequestDiff, error) {
	return &PullRequestDiff{
		Files: []DiffFile{
			{
				Filename: "test.go",
				Lines:    map[int]bool{42: true, 43: true},
			},
		},
	}, nil
}

func (m *MockClient) GetPRDetails(owner, repo string, pr int) (map[string]interface{}, error) {
	return map[string]interface{}{
		"number": pr,
		"state":  "open",
		"title":  "Test PR",
		"head": map[string]interface{}{
			"sha": "abc123def456",
		},
	}, nil
}

func (m *MockClient) CreateReview(owner, repo string, pr int, review ReviewInput) error {
	return nil
}

func (m *MockClient) FindPendingReview(owner, repo string, pr int) (int, error) {
	if m.FindPendingReviewError != nil {
		return 0, m.FindPendingReviewError
	}
	return m.PendingReviewID, nil
}

func (m *MockClient) SubmitReview(owner, repo string, pr, reviewID int, body, event string) error {
	if m.SubmitReviewError != nil {
		return m.SubmitReviewError
	}
	m.SubmittedReviewID = reviewID
	return nil
}
