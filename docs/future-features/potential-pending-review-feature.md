# GitHub Pending Review API Implementation

**Scenario:** If Github implements the pending review API feature 🎉

This document outlines how `gh-comment` would be enhanced to support it.

---

## 🚀 GitHub's API Implementation

### New API Endpoints

#### 1. Create Pending Review
```http
POST /repos/{owner}/{repo}/pulls/{pull_number}/reviews
```

**New Request Body:**
```json
{
  "body": "Starting my review",
  "event": "PENDING",  // 🆕 New event type
  "comments": [
    {
      "path": "src/api.js",
      "line": 42,
      "body": "Initial comment"
    }
  ]
}
```

**Response:**
```json
{
  "id": 12345,
  "state": "PENDING",  // 🆕 New state
  "body": "Starting my review",
  "user": {...},
  "submitted_at": null,  // null for pending reviews
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### 2. Add Comments to Pending Review
```http
POST /repos/{owner}/{repo}/pulls/{pull_number}/reviews/{review_id}/comments
```

**Request Body:**
```json
{
  "path": "src/api.js",
  "line": 50,
  "body": "Additional comment discovered during analysis"
}
```

**Response:**
```json
{
  "id": 67890,
  "pull_request_review_id": 12345,
  "path": "src/api.js",
  "line": 50,
  "body": "Additional comment discovered during analysis",
  "user": {...},
  "created_at": "2024-01-15T10:35:00Z"
}
```

#### 3. Submit Pending Review
```http
PATCH /repos/{owner}/{repo}/pulls/{pull_number}/reviews/{review_id}
```

**Request Body:**
```json
{
  "event": "COMMENT",  // or "APPROVE" or "REQUEST_CHANGES"
  "body": "Updated review summary"  // optional
}
```

#### 4. List Pending Reviews
```http
GET /repos/{owner}/{repo}/pulls/{pull_number}/reviews?state=PENDING
```

---

## 🛠️ gh-comment Implementation

### New Commands

#### 1. Add to Review (Smart Auto-Creation)
```bash
# First comment - automatically creates pending review
gh comment add-to-review 123 src/api.js 42 "First issue found"

# Subsequent comments - automatically adds to existing pending review
gh comment add-to-review 123 src/api.js 50 "Second issue discovered"
gh comment add-to-review 123 src/api.js 60 "Third issue found"

# Optional: specify review body when creating
gh comment add-to-review 123 src/api.js 42 "First issue" --review-body "Starting comprehensive analysis"
```

#### 2. Submit Pending Review
```bash
# Submit pending review as COMMENT
gh comment submit-review 123

# Submit with approval
gh comment submit-review 123 --approve

# Submit requesting changes
gh comment submit-review 123 --request-changes --body "Please address these issues"
```

#### 4. List Pending Reviews
```bash
# List your pending reviews
gh comment list-pending 123

# List all pending reviews (if you have permissions)
gh comment list-pending 123 --all
```

### Enhanced Existing Commands

#### Enhanced `add-review` Command
```bash
# Create pending review (new --pending flag)
gh comment add-review 123 --pending \
  --comment "src/api.js:42:First issue" \
  --comment "src/utils.js:15:Second issue" \
  --body "Starting comprehensive review"

# Add more comments to existing pending review
gh comment add-review 123 --add-to-pending \
  --comment "src/api.js:60:Third issue found"

# Submit the pending review
gh comment add-review 123 --submit-pending --approve
```

---

## 💻 Code Implementation in gh-comment

### New File: `cmd/pending-review.go`

```go
package cmd

import (
    "encoding/json"
    "fmt"
    "strconv"
    "github.com/cli/go-gh/v2/pkg/api"
    "github.com/spf13/cobra"
)

// Data structures for API responses
type PendingReview struct {
    ID          int    `json:"id"`
    State       string `json:"state"`
    Body        string `json:"body"`
    User        User   `json:"user"`
    SubmittedAt *string `json:"submitted_at"` // null for pending
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type User struct {
    Login string `json:"login"`
    ID    int    `json:"id"`
}

type ReviewComment struct {
    ID                 int    `json:"id"`
    PullRequestReviewID int   `json:"pull_request_review_id"`
    Path               string `json:"path"`
    Line               int    `json:"line"`
    Body               string `json:"body"`
    User               User   `json:"user"`
    CreatedAt          string `json:"created_at"`
}

var (
    pendingReviewBody string
    submitEvent       string
    approveReview     bool
    requestChanges    bool
)

var addToReviewCmd = &cobra.Command{
    Use:   "add-to-review <pr-number> <file> <line> <message>",
    Short: "Add comment to pending review (creates review if none exists)",
    Args:  cobra.ExactArgs(4),
    RunE:  runAddToReview,
}

var submitReviewCmd = &cobra.Command{
    Use:   "submit-review <pr-number>",
    Short: "Submit pending review",
    Args:  cobra.ExactArgs(1),
    RunE:  runSubmitReview,
}

func init() {
    rootCmd.AddCommand(addToReviewCmd)
    rootCmd.AddCommand(submitReviewCmd)
    
    addToReviewCmd.Flags().StringVar(&pendingReviewBody, "review-body", "", "Review body (only used when creating new review)")
    submitReviewCmd.Flags().StringVar(&submitEvent, "event", "COMMENT", "Review event: COMMENT, APPROVE, REQUEST_CHANGES")
    submitReviewCmd.Flags().BoolVar(&approveReview, "approve", false, "Approve the PR")
    submitReviewCmd.Flags().BoolVar(&requestChanges, "request-changes", false, "Request changes")
}

func runAddToReview(cmd *cobra.Command, args []string) error {
    prNumber := args[0]
    file := args[1]
    lineStr := args[2]
    message := args[3]
    
    // Parse line number
    line, err := strconv.Atoi(lineStr)
    if err != nil {
        return fmt.Errorf("invalid line number: %s", lineStr)
    }
    
    repository, err := getPRContext()
    if err != nil {
        return err
    }
    
    // Try to find existing pending review
    reviewID, err := findMyPendingReview(repository, prNumber)
    if err != nil {
        // No pending review found, create one
        reviewID, err = createPendingReview(repository, prNumber, pendingReviewBody)
        if err != nil {
            return fmt.Errorf("failed to create pending review: %w", err)
        }
        fmt.Printf("✅ Created pending review #%d\n", reviewID)
    }
    
    // Expand suggestion syntax
    expandedMessage := expandSuggestions(message)
    
    // Add comment to pending review
    commentData := map[string]interface{}{
        "path": file,
        "line": line,
        "body": expandedMessage,
    }
    
    client, err := api.DefaultRESTClient()
    if err != nil {
        return err
    }
    
    var response struct {
        ID int `json:"id"`
    }
    
    err = client.Post(fmt.Sprintf("repos/%s/pulls/%s/reviews/%d/comments", repository, prNumber, reviewID), commentData, &response)
    if err != nil {
        return err
    }
    
    fmt.Printf("✅ Added comment to pending review #%d\n", reviewID)
    return nil
}

func runSubmitReview(cmd *cobra.Command, args []string) error {
    prNumber := args[0]
    
    repository, err := getPRContext()
    if err != nil {
        return err
    }
    
    // Find my pending review
    reviewID, err := findMyPendingReview(repository, prNumber)
    if err != nil {
        return fmt.Errorf("no pending review found: %w", err)
    }
    
    // Determine event type
    event := submitEvent
    if approveReview {
        event = "APPROVE"
    } else if requestChanges {
        event = "REQUEST_CHANGES"
    }
    
    // Submit review
    submitData := map[string]interface{}{
        "event": event,
    }
    
    client, err := api.DefaultRESTClient()
    if err != nil {
        return err
    }
    
    err = client.Patch(fmt.Sprintf("repos/%s/pulls/%s/reviews/%d", repository, prNumber, reviewID), submitData, nil)
    if err != nil {
        return err
    }
    
    fmt.Printf("✅ Submitted review #%d as %s\n", reviewID, event)
    return nil
}

func findMyPendingReview(repository, prNumber string) (int, error) {
    client, err := api.DefaultRESTClient()
    if err != nil {
        return 0, err
    }
    
    var reviews []PendingReview
    
    // Get all reviews (GitHub will filter by state in the future)
    err = client.Get(fmt.Sprintf("repos/%s/pulls/%s/reviews", repository, prNumber), &reviews)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch reviews: %w", err)
    }
    
    // Find current user's pending review
    currentUser, err := getCurrentUser()
    if err != nil {
        return 0, fmt.Errorf("failed to get current user: %w", err)
    }
    
    for _, review := range reviews {
        if review.User.Login == currentUser && review.State == "PENDING" {
            return review.ID, nil
        }
    }
    
    return 0, fmt.Errorf("no pending review found for user %s", currentUser)
}

func getCurrentUser() (string, error) {
    client, err := api.DefaultRESTClient()
    if err != nil {
        return "", err
    }
    
    var user User
    err = client.Get("user", &user)
    if err != nil {
        return "", fmt.Errorf("failed to get current user: %w", err)
    }
    
    return user.Login, nil
}

func createPendingReview(repository, prNumber, body string) (int, error) {
    client, err := api.DefaultRESTClient()
    if err != nil {
        return 0, fmt.Errorf("failed to create API client: %w", err)
    }
    
    reviewData := map[string]interface{}{
        "body":  body,
        "event": "PENDING",
    }
    
    var response PendingReview
    
    err = client.Post(fmt.Sprintf("repos/%s/pulls/%s/reviews", repository, prNumber), reviewData, &response)
    if err != nil {
        return 0, fmt.Errorf("failed to create pending review: %w", err)
    }
    
    return response.ID, nil
}
```

### Enhanced Workflow Examples

#### AI-Assisted Review Workflow
```bash
# AI starts comprehensive review
gh comment start-review 123 "AI analysis beginning"

# AI finds issues incrementally and adds them
gh comment add-to-review 123 src/api.js 42 "Potential null pointer: [SUGGEST: if (user?.id) {]"
gh comment add-to-review 123 src/utils.js 15 "Consider using const: [SUGGEST: const result = processData(input);]"
gh comment add-to-review 123 src/auth.js 28 "Security concern: validate input before processing"

# AI completes analysis and submits review
gh comment submit-review 123 --request-changes
```

#### Interactive Human Review
```bash
# Human starts review
gh comment start-review 123 "Reviewing the new authentication feature"

# Human adds comments as they discover issues
gh comment add-to-review 123 src/auth.js 45 "Great improvement! The error handling looks solid."
gh comment add-to-review 123 src/api.js 67 "Minor: consider extracting this logic into a helper function"

# Human approves after thorough review
gh comment submit-review 123 --approve
```

---

## 🎯 Benefits Realized

### For AI Workflows
- ✅ **Natural incremental analysis** - Add comments as issues are discovered
- ✅ **Professional review presentation** - All comments grouped in a single review
- ✅ **Flexible workflow** - Can build comprehensive reviews over time
- ✅ **Suggestion syntax support** - Works seamlessly with `[SUGGEST: code]` expansion

### For Human Workflows  
- ✅ **API parity with web interface** - Same workflow available programmatically
- ✅ **CLI power** - All the benefits of command-line efficiency
- ✅ **Batch operations** - Can script complex review workflows
- ✅ **Integration friendly** - Easy to integrate with other tools

### For gh-comment Users
- ✅ **Backward compatibility** - Existing commands still work
- ✅ **Enhanced functionality** - New pending review capabilities
- ✅ **Consistent UX** - Same patterns and conventions
- ✅ **Shell agnostic** - Works identically across Fish, Bash, Zsh

---

## 🚀 Migration Guide

### Existing Users
```bash
# Old way (still works)
gh comment add-review 123 \
  --comment "file.js:42:Issue 1" \
  --comment "file.js:50:Issue 2" \
  --event COMMENT

# New way (more flexible)
gh comment start-review 123 "Comprehensive review"
gh comment add-to-review 123 file.js 42 "Issue 1"
gh comment add-to-review 123 file.js 50 "Issue 2"  
gh comment submit-review 123
```

### AI Prompt Updates
The AI prompts would be updated to leverage the new pending review workflow:

```bash
# Start a pending review for comprehensive feedback
gh comment start-review {PR_NUMBER} "AI-assisted code review"

# Add individual comments as you discover issues
gh comment add-to-review {PR_NUMBER} {FILE} {LINE} "{COMMENT_WITH_SUGGESTIONS}"

# Submit the complete review when analysis is finished
gh comment submit-review {PR_NUMBER} --request-changes  # or --approve
```

---

## 🧪 Testing Strategy

### Unit Tests
```go
// cmd/pending_review_test.go
func TestFindMyPendingReview(t *testing.T) {
    // Test finding existing pending review
    // Test no pending review found
    // Test API error handling
}

func TestCreatePendingReview(t *testing.T) {
    // Test successful creation
    // Test API error responses
    // Test empty body handling
}

func TestAddToReview(t *testing.T) {
    // Test adding to existing review
    // Test auto-creation workflow
    // Test suggestion expansion integration
}
```

### Integration Tests
```bash
# Test complete workflow
gh comment add-to-review 123 file.js 42 "First comment" --review-body "Test review"
gh comment add-to-review 123 file.js 50 "Second comment"
gh comment submit-review 123 --approve

# Test error cases
gh comment submit-review 999 # No pending review
gh comment add-to-review 123 invalid.js 42 "Comment" # Invalid file
```

## 🔄 Migration Strategy

### Phase 1: Backward Compatibility
- All existing commands continue to work unchanged
- New commands are additive, not replacing existing functionality
- Users can adopt incrementally

### Phase 2: Enhanced Workflows
```bash
# Existing workflow (still supported)
gh comment add-review 123 \
  --comment "file.js:42:Issue 1" \
  --comment "file.js:50:Issue 2"

# New workflow (more flexible)
gh comment add-to-review 123 file.js 42 "Issue 1"
gh comment add-to-review 123 file.js 50 "Issue 2"
gh comment submit-review 123
```

### Phase 3: Documentation Updates
- Update AI prompts to leverage pending review workflow
- Add examples showing both approaches
- Highlight benefits of new incremental approach

---

**This implementation would make `gh-comment` the definitive tool for programmatic GitHub PR reviews, with full parity to the web interface experience!** 🎉
