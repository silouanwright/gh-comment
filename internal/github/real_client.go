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
	if err := validateRepoParams(owner, repo); err != nil {
		return nil, err
	}
	if prNumber <= 0 {
		return nil, fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}

	endpoint := fmt.Sprintf("repos/%s/%s/issues/%d/comments?per_page=100", owner, repo, prNumber)

	var comments []Comment
	err := c.restClient.Get(endpoint, &comments)
	if err != nil {
		return nil, c.wrapAPIError(err, "fetch issue comments for PR #%d in %s/%s", prNumber, owner, repo)
	}

	// Mark as issue comments
	for i := range comments {
		comments[i].Type = "issue"
	}

	return comments, nil
}

// ListReviewComments fetches all review comments for a PR
func (c *RealClient) ListReviewComments(owner, repo string, prNumber int) ([]Comment, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return nil, err
	}
	if prNumber <= 0 {
		return nil, fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}

	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/comments?per_page=100", owner, repo, prNumber)

	var comments []Comment
	err := c.restClient.Get(endpoint, &comments)
	if err != nil {
		return nil, c.wrapAPIError(err, "fetch review comments for PR #%d in %s/%s", prNumber, owner, repo)
	}

	// Mark as review comments
	for i := range comments {
		comments[i].Type = "review"
	}

	return comments, nil
}

// CreateIssueComment adds a general comment to a PR
func (c *RealClient) CreateIssueComment(owner, repo string, prNumber int, body string) (*Comment, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return nil, err
	}
	if prNumber <= 0 {
		return nil, fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}
	if strings.TrimSpace(body) == "" {
		return nil, fmt.Errorf("comment body cannot be empty")
	}

	endpoint := fmt.Sprintf("repos/%s/%s/issues/%d/comments", owner, repo, prNumber)

	payload := map[string]string{"body": body}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comment payload: %w", err)
	}

	var comment Comment
	err = c.restClient.Post(endpoint, bytes.NewReader(bodyBytes), &comment)
	if err != nil {
		return nil, c.wrapAPIError(err, "create issue comment on PR #%d in %s/%s", prNumber, owner, repo)
	}

	comment.Type = "issue"
	return &comment, nil
}

// CreateReviewCommentReply adds a reply to a review comment
func (c *RealClient) CreateReviewCommentReply(owner, repo string, commentID int, body string) (*Comment, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return nil, err
	}
	if commentID <= 0 {
		return nil, fmt.Errorf("invalid comment ID %d: must be positive", commentID)
	}
	if strings.TrimSpace(body) == "" {
		return nil, fmt.Errorf("reply body cannot be empty")
	}

	// First, get the original comment to find PR number, path, and commit
	originalEndpoint := fmt.Sprintf("repos/%s/%s/pulls/comments/%d", owner, repo, commentID)
	var originalComment struct {
		PullRequestURL string `json:"pull_request_url"`
		Path           string `json:"path"`
		CommitID       string `json:"commit_id"`
		Line           int    `json:"line"`
	}

	err := c.restClient.Get(originalEndpoint, &originalComment)
	if err != nil {
		return nil, c.wrapAPIError(err, "get original comment #%d in %s/%s", commentID, owner, repo)
	}

	// Extract PR number from pull_request_url
	// URL format: https://api.github.com/repos/owner/repo/pulls/123
	urlParts := strings.Split(originalComment.PullRequestURL, "/")
	if len(urlParts) < 1 {
		return nil, fmt.Errorf("invalid pull request URL format: %s", originalComment.PullRequestURL)
	}
	prNumber := urlParts[len(urlParts)-1]

	// Create a new review comment as a reply (with in_reply_to_id)
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%s/comments", owner, repo, prNumber)

	payload := map[string]interface{}{
		"body":           body,
		"path":           originalComment.Path,
		"line":           originalComment.Line,
		"commit_id":      originalComment.CommitID,
		"in_reply_to_id": commentID,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal reply payload: %w", err)
	}

	var comment Comment
	err = c.restClient.Post(endpoint, bytes.NewReader(bodyBytes), &comment)
	if err != nil {
		// Provide intelligent error analysis
		enhancedErr := AnalyzeAndEnhanceError(err, "reply", commentID)
		return nil, enhancedErr
	}

	comment.Type = "review"
	return &comment, nil
}

// FindReviewThreadForComment finds the thread ID for a review comment
func (c *RealClient) FindReviewThreadForComment(owner, repo string, prNumber, commentID int) (string, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return "", err
	}
	if prNumber <= 0 {
		return "", fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}
	if commentID <= 0 {
		return "", fmt.Errorf("invalid comment ID %d: must be positive", commentID)
	}
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
		return "", c.wrapAPIError(err, "find review thread for comment #%d in PR #%d (%s/%s)", commentID, prNumber, owner, repo)
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
	if strings.TrimSpace(threadID) == "" {
		return fmt.Errorf("thread ID cannot be empty")
	}

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
		return c.wrapAPIError(err, "resolve review thread %s", threadID)
	}

	return nil
}

// Additional methods for other operations can be added here as needed...

// AddReaction adds a reaction to a comment
func (c *RealClient) AddReaction(owner, repo string, commentID int, prNumber int, reaction string) error {
	if err := validateRepoParams(owner, repo); err != nil {
		return err
	}
	if commentID <= 0 {
		return fmt.Errorf("invalid comment ID %d: must be positive", commentID)
	}
	if prNumber <= 0 {
		return fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}
	if !isValidReaction(reaction) {
		return fmt.Errorf("invalid reaction '%s': must be one of +1, -1, laugh, hooray, confused, heart, rocket, eyes", reaction)
	}

	// Detect comment type first to use the correct endpoint
	commentInfo, err := c.DetectCommentType(owner, repo, commentID, prNumber)
	if err != nil {
		return CreateSmartError(c, "add_reaction", "reply", commentID, prNumber, err)
	}

	if !commentInfo.Found {
		return CreateSmartError(c, "add_reaction", "reply", commentID, prNumber, fmt.Errorf("comment #%d not found", commentID))
	}

	// Use the correct endpoint based on detected comment type
	var endpoint string
	if commentInfo.Type == "review" {
		endpoint = fmt.Sprintf("repos/%s/%s/pulls/comments/%d/reactions", owner, repo, commentID)
	} else {
		endpoint = fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions", owner, repo, commentID)
	}

	payload := map[string]string{"content": reaction}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal reaction payload: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(body), nil)
	if err != nil {
		return CreateSmartError(c, "add_reaction", "reply", commentID, prNumber, err)
	}

	return nil
}

// RemoveReaction removes a reaction from a comment
func (c *RealClient) RemoveReaction(owner, repo string, commentID int, prNumber int, reaction string) error {
	if err := validateRepoParams(owner, repo); err != nil {
		return err
	}
	if commentID <= 0 {
		return fmt.Errorf("invalid comment ID %d: must be positive", commentID)
	}
	if prNumber <= 0 {
		return fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}
	if !isValidReaction(reaction) {
		return fmt.Errorf("invalid reaction '%s': must be one of +1, -1, laugh, hooray, confused, heart, rocket, eyes", reaction)
	}

	// Detect comment type first to use the correct endpoint
	commentInfo, err := c.DetectCommentType(owner, repo, commentID, prNumber)
	if err != nil {
		return CreateSmartError(c, "remove_reaction", "reply", commentID, prNumber, err)
	}

	if !commentInfo.Found {
		return CreateSmartError(c, "remove_reaction", "reply", commentID, prNumber, fmt.Errorf("comment #%d not found", commentID))
	}

	// Use the correct endpoint based on detected comment type
	var endpoint string
	if commentInfo.Type == "review" {
		endpoint = fmt.Sprintf("repos/%s/%s/pulls/comments/%d/reactions", owner, repo, commentID)
	} else {
		endpoint = fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions", owner, repo, commentID)
	}

	// Get reactions to find the ID to delete
	var reactions []struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
		User    User   `json:"user"`
	}

	err = c.restClient.Get(endpoint, &reactions)
	if err != nil {
		return CreateSmartError(c, "remove_reaction", "reply", commentID, prNumber, err)
	}

	// Find current user's reaction
	var currentUser struct {
		Login string `json:"login"`
	}
	err = c.restClient.Get("user", &currentUser)
	if err != nil {
		return c.wrapAPIError(err, "get current user info")
	}

	// Find and delete the reaction
	for _, r := range reactions {
		if r.Content == reaction && r.User.Login == currentUser.Login {
			// Use the correct delete endpoint based on comment type
			var deleteEndpoint string
			if commentInfo.Type == "review" {
				deleteEndpoint = fmt.Sprintf("repos/%s/%s/pulls/comments/%d/reactions/%d", owner, repo, commentID, r.ID)
			} else {
				deleteEndpoint = fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions/%d", owner, repo, commentID, r.ID)
			}

			err = c.restClient.Delete(deleteEndpoint, nil)
			if err != nil {
				return CreateSmartError(c, "remove_reaction", "reply", commentID, prNumber, err)
			}
			return nil
		}
	}

	return fmt.Errorf("'%s' reaction not found on comment #%d (you may not have reacted with this emoji)", reaction, commentID)
}

// EditComment edits an existing comment
func (c *RealClient) EditComment(owner, repo string, commentID int, prNumber int, body string) error {
	if err := validateRepoParams(owner, repo); err != nil {
		return err
	}
	if commentID <= 0 {
		return fmt.Errorf("invalid comment ID %d: must be positive", commentID)
	}
	if prNumber <= 0 {
		return fmt.Errorf("invalid PR number %d: must be positive", prNumber)
	}
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("comment body cannot be empty")
	}

	// Detect comment type first to use the correct endpoint
	commentInfo, err := c.DetectCommentType(owner, repo, commentID, prNumber)
	if err != nil {
		return CreateSmartError(c, "edit", "edit", commentID, prNumber, err)
	}

	if !commentInfo.Found {
		return CreateSmartError(c, "edit", "edit", commentID, prNumber, fmt.Errorf("comment #%d not found", commentID))
	}

	// Use the correct endpoint based on detected comment type
	var endpoint string
	if commentInfo.Type == "review" {
		endpoint = fmt.Sprintf("repos/%s/%s/pulls/comments/%d", owner, repo, commentID)
	} else {
		endpoint = fmt.Sprintf("repos/%s/%s/issues/comments/%d", owner, repo, commentID)
	}

	payload := map[string]string{"body": body}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal comment payload: %w", err)
	}

	err = c.restClient.Patch(endpoint, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return CreateSmartError(c, "edit", "edit", commentID, prNumber, err)
	}

	return nil
}

// AddReviewComment adds a line-specific comment to a PR
func (c *RealClient) AddReviewComment(owner, repo string, pr int, comment ReviewCommentInput) error {
	if err := validateRepoParams(owner, repo); err != nil {
		return err
	}
	if pr <= 0 {
		return fmt.Errorf("invalid PR number %d: must be positive", pr)
	}
	if strings.TrimSpace(comment.Body) == "" {
		return fmt.Errorf("review comment body cannot be empty")
	}
	if comment.Path == "" {
		return fmt.Errorf("review comment path cannot be empty")
	}

	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/comments", owner, repo, pr)

	body, err := json.Marshal(comment)
	if err != nil {
		return fmt.Errorf("failed to marshal review comment payload: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(body), nil)
	if err != nil {
		return c.wrapAPIError(err, "add review comment to %s:%d in PR #%d (%s/%s)", comment.Path, comment.Line, pr, owner, repo)
	}

	return nil
}

// FetchPRDiff fetches the diff for a pull request
func (c *RealClient) FetchPRDiff(owner, repo string, pr int) (*PullRequestDiff, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return nil, err
	}
	if pr <= 0 {
		return nil, fmt.Errorf("invalid PR number %d: must be positive", pr)
	}

	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, pr)

	resp, err := c.restClient.Request("GET", endpoint, nil)
	if err != nil {
		return nil, c.wrapAPIError(err, "fetch PR #%d details from %s/%s", pr, owner, repo)
	}
	defer resp.Body.Close()

	// Parse the diff URL from the response
	var prData struct {
		DiffURL string `json:"diff_url"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PR response: %w", err)
	}

	err = json.Unmarshal(body, &prData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PR data: %w", err)
	}

	if prData.DiffURL == "" {
		return nil, fmt.Errorf("PR #%d does not have a diff URL (may be empty or merged)", pr)
	}

	// Fetch the actual diff
	diffResp, err := http.Get(prData.DiffURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch diff from GitHub: %w", err)
	}
	defer diffResp.Body.Close()

	if diffResp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch diff: HTTP %d", diffResp.StatusCode)
	}

	diffContent, err := io.ReadAll(diffResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff content: %w", err)
	}

	// Parse the diff to extract file and line information
	diff := parseDiff(string(diffContent))

	return diff, nil
}

// validateRepoParams validates repository owner and name parameters
func validateRepoParams(owner, repo string) error {
	if owner == "" {
		return fmt.Errorf("repository owner cannot be empty")
	}
	if repo == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	if strings.Contains(owner, "/") || strings.Contains(repo, "/") {
		return fmt.Errorf("invalid repository format: use 'owner/repo' format")
	}
	return nil
}

// checkRateLimit provides rate limit awareness following GitHub CLI patterns
// This provides user guidance without automatic retries (which GitHub CLI team discourages)
func (c *RealClient) checkRateLimit() {
	// This is a placeholder for rate limit checking
	// In a production implementation, you might:
	// 1. Check X-RateLimit-Remaining header after requests
	// 2. Warn when approaching limits (50-75% usage)
	// 3. Provide proactive guidance to users
	//
	// GitHub CLI team prefers user awareness over automatic handling
}

// isValidReaction checks if the reaction is valid for GitHub API
func isValidReaction(reaction string) bool {
	validReactions := map[string]bool{
		"+1":       true,
		"-1":       true,
		"laugh":    true,
		"hooray":   true,
		"confused": true,
		"heart":    true,
		"rocket":   true,
		"eyes":     true,
	}
	return validReactions[reaction]
}

// wrapAPIError wraps GitHub API errors with context and rate limit information
// Following GitHub CLI philosophy: provide helpful guidance instead of automatic retries
func (c *RealClient) wrapAPIError(err error, operation string, args ...interface{}) error {
	context := fmt.Sprintf(operation, args...)

	// Check if this is a rate limit error
	if strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "403") {
		return fmt.Errorf("rate limit exceeded while trying to %s: %w\n\n💡 Tips:\n   • Wait a few minutes before retrying\n   • Check your rate limit status: gh api rate_limit\n   • Consider reducing API calls if this happens frequently", context, err)
	}

	// Check for common API errors and provide helpful messages
	if strings.Contains(err.Error(), "404") {
		return fmt.Errorf("resource not found while trying to %s: %w\n\n💡 Tips:\n   • Verify the repository exists and you have access to it\n   • Check the PR/comment ID is correct\n   • Ensure you have the right permissions", context, err)
	}

	if strings.Contains(err.Error(), "401") {
		return fmt.Errorf("authentication failed while trying to %s: %w\n\n💡 Tips:\n   • Check your GitHub CLI authentication: gh auth status\n   • Re-authenticate if needed: gh auth login\n   • Verify you have access to this repository", context, err)
	}

	if strings.Contains(err.Error(), "422") {
		return fmt.Errorf("validation error while trying to %s: %w\n\n💡 Tips:\n   • Check that your input parameters are valid\n   • Verify line numbers exist in the diff\n   • Ensure comment body is not empty", context, err)
	}

	// Check for secondary rate limits (GitHub doesn't always send proper headers)
	if strings.Contains(err.Error(), "abuse") || strings.Contains(err.Error(), "secondary") {
		return fmt.Errorf("secondary rate limit triggered while trying to %s: %w\n\n💡 Tips:\n   • This is a temporary protective measure by GitHub\n   • Wait 60 seconds before retrying\n   • Reduce the frequency of API calls", context, err)
	}

	// Generic API error
	return fmt.Errorf("GitHub API error while trying to %s: %w", context, err)
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
	if err := validateRepoParams(owner, repo); err != nil {
		return err
	}
	if pr <= 0 {
		return fmt.Errorf("invalid PR number %d: must be positive", pr)
	}
	// Validate review event if provided
	if review.Event != "" && review.Event != "APPROVE" && review.Event != "REQUEST_CHANGES" && review.Event != "COMMENT" {
		return fmt.Errorf("invalid review event '%s': must be APPROVE, REQUEST_CHANGES, or COMMENT", review.Event)
	}

	// If there are comments, we need to set commit_id for each one
	if len(review.Comments) > 0 {
		// Try to get commit ID from existing review comments first (more efficient)
		commitID, err := c.getCommitIDFromExistingComments(owner, repo, pr)
		if err != nil || commitID == "" {
			// Fallback to fetching PR details for latest commit SHA
			prDetails, err := c.GetPRDetails(owner, repo, pr)
			if err != nil {
				return fmt.Errorf("failed to get PR details for commit ID: %w", err)
			}

			// Extract the latest commit SHA
			head, ok := prDetails["head"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid PR details: missing head information")
			}

			latestCommitSHA, ok := head["sha"].(string)
			if !ok {
				return fmt.Errorf("invalid PR details: missing commit SHA")
			}
			commitID = latestCommitSHA
		}

		// Add commit_id to each comment
		for i := range review.Comments {
			review.Comments[i].CommitID = commitID
		}
	}

	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/reviews", owner, repo, pr)

	body, err := json.Marshal(review)
	if err != nil {
		return fmt.Errorf("failed to marshal review payload: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(body), nil)
	if err != nil {
		// Provide intelligent error analysis
		enhancedErr := AnalyzeAndEnhanceError(err, "review", pr)
		return enhancedErr
	}

	return nil
}

// GetPRDetails fetches basic PR information
func (c *RealClient) GetPRDetails(owner, repo string, pr int) (map[string]interface{}, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return nil, err
	}
	if pr <= 0 {
		return nil, fmt.Errorf("invalid PR number %d: must be positive", pr)
	}

	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, pr)

	var result map[string]interface{}
	err := c.restClient.Get(endpoint, &result)
	if err != nil {
		return nil, c.wrapAPIError(err, "fetch PR #%d details from %s/%s", pr, owner, repo)
	}

	return result, nil
}

// FindPendingReview finds a pending review for the current user on a PR
func (c *RealClient) FindPendingReview(owner, repo string, pr int) (int, error) {
	if err := validateRepoParams(owner, repo); err != nil {
		return 0, err
	}
	if pr <= 0 {
		return 0, fmt.Errorf("invalid PR number %d: must be positive", pr)
	}

	// Get existing reviews for this PR
	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/reviews", owner, repo, pr)

	var reviews []map[string]interface{}
	err := c.restClient.Get(endpoint, &reviews)
	if err != nil {
		return 0, c.wrapAPIError(err, "get reviews for PR #%d in %s/%s", pr, owner, repo)
	}

	// Look for an existing PENDING review
	for _, review := range reviews {
		if state, ok := review["state"].(string); ok && state == "PENDING" {
			if id, ok := review["id"].(float64); ok {
				return int(id), nil
			}
		}
	}

	return 0, nil // No pending review found
}

// SubmitReview submits a pending review with a body and event
func (c *RealClient) SubmitReview(owner, repo string, pr, reviewID int, body, event string) error {
	if err := validateRepoParams(owner, repo); err != nil {
		return err
	}
	if pr <= 0 {
		return fmt.Errorf("invalid PR number %d: must be positive", pr)
	}
	if reviewID <= 0 {
		return fmt.Errorf("invalid review ID %d: must be positive", reviewID)
	}
	if event != "APPROVE" && event != "REQUEST_CHANGES" && event != "COMMENT" {
		return fmt.Errorf("invalid review event '%s': must be APPROVE, REQUEST_CHANGES, or COMMENT", event)
	}

	endpoint := fmt.Sprintf("repos/%s/%s/pulls/%d/reviews/%d/events", owner, repo, pr, reviewID)

	payload := map[string]interface{}{
		"body":  body,
		"event": event,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal submit review payload: %w", err)
	}

	err = c.restClient.Post(endpoint, bytes.NewReader(payloadBytes), nil)
	if err != nil {
		return c.wrapAPIError(err, "submit review #%d on PR #%d in %s/%s", reviewID, pr, owner, repo)
	}

	return nil
}

// getCommitIDFromExistingComments tries to get a commit ID from existing review comments
// This is more efficient than fetching PR details since we might already have review comments
func (c *RealClient) getCommitIDFromExistingComments(owner, repo string, pr int) (string, error) {
	// Get existing review comments
	comments, err := c.ListReviewComments(owner, repo, pr)
	if err != nil {
		return "", err
	}

	// If we have any review comments, use the commit ID from the most recent one
	// This assumes that the most recent comment is likely on the latest commit
	if len(comments) > 0 {
		// Find the most recent comment with a commit ID
		for i := len(comments) - 1; i >= 0; i-- {
			if comments[i].CommitID != "" {
				return comments[i].CommitID, nil
			}
		}
	}

	return "", nil // No commit ID found
}
