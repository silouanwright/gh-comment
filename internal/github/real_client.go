package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

// RealClient implements GitHubAPI using actual GitHub API calls
type RealClient struct {
	restClient    *api.RESTClient
	graphqlClient *api.GraphQLClient
}

// NewRealClient creates a new GitHub API client
func NewRealClient() (*RealClient, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}

	graphqlClient, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL client: %w", err)
	}

	return &RealClient{
		restClient:    restClient,
		graphqlClient: graphqlClient,
	}, nil
}

// ListIssueComments fetches all issue comments for a PR
func (c *RealClient) ListIssueComments(owner, repo string, prNumber int) ([]Comment, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/issues/%d/comments?per_page=100", owner, repo, prNumber)

	var comments []Comment
	err := c.restClient.Get(endpoint, &comments)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue comments: %w", err)
	}

	// Mark as issue comments
	for i := range comments {
		comments[i].Type = "issue"
	}

	return comments, nil
}

// ListReviewComments fetches all review comments for a PR
func (c *RealClient) ListReviewComments(owner, repo string, prNumber int) ([]Comment, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/comments?per_page=100", owner, repo, prNumber)

	var comments []Comment
	err := c.restClient.Get(endpoint, &comments)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review comments: %w", err)
	}

	// Mark as review comments
	for i := range comments {
		comments[i].Type = "review"
	}

	return comments, nil
}

// CreateIssueComment adds a general comment to a PR
func (c *RealClient) CreateIssueComment(owner, repo string, prNumber int, body string) (*Comment, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, prNumber)

	payload := map[string]string{"body": body}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comment: %w", err)
	}

	var comment Comment
	err = c.restClient.Post(endpoint, bytes.NewReader(bodyBytes), &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to add issue comment: %w", err)
	}

	comment.Type = "issue"
	return &comment, nil
}

// CreateReviewCommentReply adds a reply to a review comment
func (c *RealClient) CreateReviewCommentReply(owner, repo string, commentID int, body string) (*Comment, error) {
	// For review comments, we need to get the PR number first
	// This is a simplified version - in production, you'd want to cache or pass this info
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/comments/%d/replies", owner, repo, commentID)

	payload := map[string]string{"body": body}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reply: %w", err)
	}

	var comment Comment
	err = c.restClient.Post(endpoint, bytes.NewReader(bodyBytes), &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to add review comment reply: %w", err)
	}

	comment.Type = "review"
	return &comment, nil
}

// FindReviewThreadForComment finds the thread ID for a review comment
func (c *RealClient) FindReviewThreadForComment(owner, repo string, prNumber, commentID int) (string, error) {
	query := `
		query($owner: String!, $name: String!, $number: Int!) {
			repository(owner: $owner, name: $name) {
				pullRequest(number: $number) {
					reviewThreads(first: 100) {
						nodes {
							id
							comments(first: 10) {
								nodes {
									databaseId
								}
							}
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"owner":  owner,
		"name":   repo,
		"number": prNumber,
	}

	var result struct {
		Repository struct {
			PullRequest struct {
				ReviewThreads struct {
					Nodes []struct {
						ID       string `json:"id"`
						Comments struct {
							Nodes []struct {
								DatabaseID int `json:"databaseId"`
							} `json:"nodes"`
						} `json:"comments"`
					} `json:"nodes"`
				} `json:"reviewThreads"`
			} `json:"pullRequest"`
		} `json:"repository"`
	}

	err := c.graphqlClient.Do(query, variables, &result)
	if err != nil {
		return "", fmt.Errorf("failed to find thread: %w", err)
	}

	// Find the thread containing our comment
	for _, thread := range result.Repository.PullRequest.ReviewThreads.Nodes {
		for _, comment := range thread.Comments.Nodes {
			if comment.DatabaseID == commentID {
				return thread.ID, nil
			}
		}
	}

	return "", fmt.Errorf("thread not found for comment %d", commentID)
}

// ResolveReviewThread resolves a review thread
func (c *RealClient) ResolveReviewThread(threadID string) error {
	mutation := `
		mutation($threadId: ID!) {
			resolveReviewThread(input: {threadId: $threadId}) {
				thread {
					id
				}
			}
		}`

	variables := map[string]interface{}{
		"threadId": threadID,
	}

	err := c.graphqlClient.Do(mutation, variables, nil)
	if err != nil {
		return fmt.Errorf("failed to resolve thread: %w", err)
	}

	return nil
}

// Additional methods for other operations can be added here as needed...

// AddReaction adds a reaction to a comment
func (c *RealClient) AddReaction(owner, repo string, commentID int, reaction string) error {
	endpoint := fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions", owner, repo, commentID)

	payload := map[string]string{"content": reaction}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal reaction: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

// RemoveReaction removes a reaction from a comment
func (c *RealClient) RemoveReaction(owner, repo string, commentID int, reaction string) error {
	// First, get reactions to find the ID to delete
	endpoint := fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions", owner, repo, commentID)

	var reactions []struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
		User    User   `json:"user"`
	}

	err := c.restClient.Get(endpoint, &reactions)
	if err != nil {
		return fmt.Errorf("failed to fetch reactions: %w", err)
	}

	// Find current user's reaction
	var currentUser struct {
		Login string `json:"login"`
	}
	err = c.restClient.Get("user", &currentUser)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Find and delete the reaction
	for _, r := range reactions {
		if r.Content == reaction && r.User.Login == currentUser.Login {
			deleteEndpoint := fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions/%d", owner, repo, commentID, r.ID)
			err = c.restClient.Delete(deleteEndpoint, nil)
			if err != nil {
				return fmt.Errorf("failed to remove reaction: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("reaction not found")
}

// EditComment edits an existing comment
func (c *RealClient) EditComment(owner, repo string, commentID int, body string) error {
	endpoint := fmt.Sprintf("repos/%s/%s/issues/comments/%d", owner, repo, commentID)

	payload := map[string]string{"body": body}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}

	err = c.restClient.Patch(endpoint, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to edit comment: %w", err)
	}

	return nil
}

// AddReviewComment adds a line-specific comment to a PR
func (c *RealClient) AddReviewComment(owner, repo string, pr int, comment ReviewCommentInput) error {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/comments", owner, repo, pr)

	body, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("failed to marshal comment: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to add review comment: %w", err)
	}

	return nil
}

// FetchPRDiff fetches the diff for a pull request
func (c *RealClient) FetchPRDiff(owner, repo string, pr int) (*PullRequestDiff, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, pr)

	resp, err := c.restClient.Request("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR: %w", err)
	}
	defer resp.Body.Close()

	// Parse the diff URL from the response
	var prData struct {
		DiffURL string `json:"diff_url"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	err = json.Unmarshal(body, &prData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PR data: %w", err)
	}

	// Fetch the actual diff
	diffResp, err := http.Get(prData.DiffURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch diff: %w", err)
	}
	defer diffResp.Body.Close()

	diffContent, err := io.ReadAll(diffResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff: %w", err)
	}

	// Parse the diff to extract file and line information
	diff := parseDiff(string(diffContent))

	return diff, nil
}

// Helper function to parse diff content
func parseDiff(diffContent string) *PullRequestDiff {
	// This is a simplified diff parser
	// In a real implementation, you'd want to use a proper diff parsing library
	diff := &PullRequestDiff{
		Files: []DiffFile{},
	}

	// Basic parsing logic - this would need to be more sophisticated
	lines := strings.Split(diffContent, "\n")
	var currentFile *DiffFile

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// Extract filename
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				filename := strings.TrimPrefix(parts[3], "b/")
				currentFile = &DiffFile{
					Filename: filename,
					Lines:    make(map[int]bool),
				}
				diff.Files = append(diff.Files, *currentFile)
			}
		} else if strings.HasPrefix(line, "@@") && currentFile != nil {
			// Parse line numbers from hunk header
			// This is simplified - real implementation would parse the hunk header properly
		}
	}

	return diff
}

// CreateReview creates a new review with comments
func (c *RealClient) CreateReview(owner, repo string, pr int, review ReviewInput) error {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/reviews", owner, repo, pr)

	body, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("failed to marshal review: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

// GetPRDetails fetches basic PR information
func (c *RealClient) GetPRDetails(owner, repo string, pr int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, pr)

	var result map[string]interface{}
	err := c.restClient.Get(endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR details: %w", err)
	}

	return result, nil
}
