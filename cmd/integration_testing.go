package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MockGitHubServer represents a mock GitHub API server for integration testing
type MockGitHubServer struct {
	server   *httptest.Server
	mu       sync.RWMutex
	comments map[string][]MockComment
	reviews  map[string][]MockReview
	users    map[string]MockUser
}

// MockComment represents a GitHub comment for testing
type MockComment struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	User      MockUser  `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	HTMLURL   string    `json:"html_url"`
	Path      string    `json:"path,omitempty"`
	Line      int       `json:"line,omitempty"`
	StartLine int       `json:"start_line,omitempty"`
}

// MockReview represents a GitHub review for testing
type MockReview struct {
	ID       int           `json:"id"`
	Body     string        `json:"body"`
	User     MockUser      `json:"user"`
	State    string        `json:"state"`
	Comments []MockComment `json:"comments,omitempty"`
}

// MockUser represents a GitHub user for testing
type MockUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

// MockPRDetails represents PR details for testing
type MockPRDetails struct {
	Number int `json:"number"`
	Head   struct {
		SHA string `json:"sha"`
	} `json:"head"`
}

// NewMockGitHubServer creates a new mock GitHub API server
func NewMockGitHubServer() *MockGitHubServer {
	s := &MockGitHubServer{
		comments: make(map[string][]MockComment),
		reviews:  make(map[string][]MockReview),
		users: map[string]MockUser{
			"test-user":    {Login: "test-user", ID: 1},
			"reviewer":     {Login: "reviewer", ID: 2},
			"senior-dev":   {Login: "senior-dev", ID: 3},
			"security-bot": {Login: "security-bot", ID: 4},
		},
	}

	mux := http.NewServeMux()

	// GET /repos/{owner}/{repo}/pulls/{pr}/comments - List review comments
	mux.HandleFunc("/repos/", s.handleRepoRequests)

	s.server = httptest.NewServer(mux)
	return s
}

// URL returns the mock server URL
func (s *MockGitHubServer) URL() string {
	return s.server.URL
}

// Close shuts down the mock server
func (s *MockGitHubServer) Close() {
	s.server.Close()
}

// AddComment adds a mock comment to the server state
func (s *MockGitHubServer) AddComment(repo string, pr int, comment MockComment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s/%d", repo, pr)
	if comment.ID == 0 {
		comment.ID = len(s.comments[key]) + 1000
	}
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}
	if comment.UpdatedAt.IsZero() {
		comment.UpdatedAt = comment.CreatedAt
	}
	if comment.HTMLURL == "" {
		comment.HTMLURL = fmt.Sprintf("%s/repos/%s/pulls/%d#issuecomment-%d", s.URL(), repo, pr, comment.ID)
	}

	s.comments[key] = append(s.comments[key], comment)
}

// GetComments returns comments for a PR
func (s *MockGitHubServer) GetComments(repo string, pr int) []MockComment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := fmt.Sprintf("%s/%d", repo, pr)
	return s.comments[key]
}

// SetupTestScenario sets up predefined test data
func (s *MockGitHubServer) SetupTestScenario(scenario string) {
	switch scenario {
	case "basic":
		s.AddComment("test-owner/test-repo", 123, MockComment{
			Body: "This looks good to me!",
			User: s.users["test-user"],
		})
		s.AddComment("test-owner/test-repo", 123, MockComment{
			Body: "Please fix the typo in line 42",
			User: s.users["reviewer"],
			Path: "src/main.go",
			Line: 42,
		})
	case "security-review":
		s.AddComment("test-owner/test-repo", 456, MockComment{
			Body: "Security scan detected potential SQL injection vulnerability",
			User: s.users["security-bot"],
			Path: "database.py",
			Line: 156,
		})
		s.AddComment("test-owner/test-repo", 456, MockComment{
			Body: "Use crypto.randomBytes(32) instead of Math.random() for token generation",
			User: s.users["senior-dev"],
			Path: "auth.go",
			Line: 67,
		})
	}
}

// handleRepoRequests handles all repository-related API requests
func (s *MockGitHubServer) handleRepoRequests(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/repos/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		http.Error(w, "Invalid repository path", http.StatusBadRequest)
		return
	}

	owner, repo := parts[0], parts[1]
	repoKey := fmt.Sprintf("%s/%s", owner, repo)

	// Handle different API endpoints
	if len(parts) >= 4 && parts[2] == "pulls" {
		prStr := parts[3]
		pr, err := strconv.Atoi(prStr)
		if err != nil {
			http.Error(w, "Invalid PR number", http.StatusBadRequest)
			return
		}

		if len(parts) == 4 {
			// GET /repos/{owner}/{repo}/pulls/{pr} - Get PR details
			if r.Method == "GET" {
				s.handleGetPRDetails(w, r, repoKey, pr)
				return
			}
		}

		if len(parts) >= 5 && parts[4] == "comments" {
			switch r.Method {
			case "GET":
				// GET /repos/{owner}/{repo}/pulls/{pr}/comments - List comments
				s.handleListComments(w, r, repoKey, pr)
			case "POST":
				// POST /repos/{owner}/{repo}/pulls/{pr}/comments - Create comment
				s.handleCreateComment(w, r, repoKey, pr)
			}
		}

		if len(parts) >= 5 && parts[4] == "reviews" {
			if r.Method == "POST" {
				// POST /repos/{owner}/{repo}/pulls/{pr}/reviews - Create review
				s.handleCreateReview(w, r, repoKey, pr)
			}
		}
	}

	if len(parts) >= 4 && parts[2] == "issues" {
		prStr := parts[3]
		pr, err := strconv.Atoi(prStr)
		if err != nil {
			http.Error(w, "Invalid issue number", http.StatusBadRequest)
			return
		}

		if len(parts) >= 5 && parts[4] == "comments" {
			switch r.Method {
			case "GET":
				// GET /repos/{owner}/{repo}/issues/{pr}/comments - List issue comments
				s.handleListIssueComments(w, r, repoKey, pr)
			case "POST":
				// POST /repos/{owner}/{repo}/issues/{pr}/comments - Create issue comment
				s.handleCreateIssueComment(w, r, repoKey, pr)
			}
		}
	}
}

// handleGetPRDetails handles GET /repos/{owner}/{repo}/pulls/{pr}
func (s *MockGitHubServer) handleGetPRDetails(w http.ResponseWriter, r *http.Request, repo string, pr int) {
	details := MockPRDetails{
		Number: pr,
	}
	details.Head.SHA = "abc123def456" // Mock commit SHA

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// handleListComments handles GET /repos/{owner}/{repo}/pulls/{pr}/comments
func (s *MockGitHubServer) handleListComments(w http.ResponseWriter, r *http.Request, repo string, pr int) {
	comments := s.GetComments(repo, pr)

	// Filter for review comments only (have path/line)
	var reviewComments []MockComment
	for _, comment := range comments {
		if comment.Path != "" {
			reviewComments = append(reviewComments, comment)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviewComments)
}

// handleListIssueComments handles GET /repos/{owner}/{repo}/issues/{pr}/comments
func (s *MockGitHubServer) handleListIssueComments(w http.ResponseWriter, r *http.Request, repo string, pr int) {
	comments := s.GetComments(repo, pr)

	// Filter for issue comments only (no path/line)
	var issueComments []MockComment
	for _, comment := range comments {
		if comment.Path == "" {
			issueComments = append(issueComments, comment)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issueComments)
}

// handleCreateComment handles POST /repos/{owner}/{repo}/pulls/{pr}/comments
func (s *MockGitHubServer) handleCreateComment(w http.ResponseWriter, r *http.Request, repo string, pr int) {
	var comment MockComment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set defaults
	comment.User = s.users["test-user"] // Default test user
	s.AddComment(repo, pr, comment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// handleCreateIssueComment handles POST /repos/{owner}/{repo}/issues/{pr}/comments
func (s *MockGitHubServer) handleCreateIssueComment(w http.ResponseWriter, r *http.Request, repo string, pr int) {
	var comment MockComment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Issue comments don't have path/line
	comment.Path = ""
	comment.Line = 0
	comment.StartLine = 0
	comment.User = s.users["test-user"]

	s.AddComment(repo, pr, comment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// handleCreateReview handles POST /repos/{owner}/{repo}/pulls/{pr}/reviews
func (s *MockGitHubServer) handleCreateReview(w http.ResponseWriter, r *http.Request, repo string, pr int) {
	var review MockReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s/%d", repo, pr)
	review.ID = len(s.reviews[key]) + 2000
	review.User = s.users["test-user"]

	// Add review comments to the comment list
	for _, comment := range review.Comments {
		comment.User = s.users["test-user"]
		s.AddComment(repo, pr, comment)
	}

	s.reviews[key] = append(s.reviews[key], review)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}
