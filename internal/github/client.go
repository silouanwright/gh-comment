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

// MockClient implements GitHubAPI for testing
type MockClient struct {
	IssueComments  []Comment
	ReviewComments []Comment
	CreatedComment *Comment
	ResolvedThread string

	// Error simulation
	ListIssueCommentsError  error
	ListReviewCommentsError error
	CreateCommentError      error
	ResolveThreadError      error
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
	return "RT_123", nil
}

func (m *MockClient) ResolveReviewThread(threadID string) error {
	if m.ResolveThreadError != nil {
		return m.ResolveThreadError
	}
	m.ResolvedThread = threadID
	return nil
}
