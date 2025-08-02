package github

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMockClient(t *testing.T) {
	client := NewMockClient()

	assert.NotNil(t, client)
	assert.Len(t, client.IssueComments, 1)
	assert.Len(t, client.ReviewComments, 1)

	// Verify default issue comment
	issueComment := client.IssueComments[0]
	assert.Equal(t, 123456, issueComment.ID)
	assert.Equal(t, "LGTM! Great work on this PR.", issueComment.Body)
	assert.Equal(t, "issue", issueComment.Type)
	assert.Equal(t, "reviewer1", issueComment.User.Login)

	// Verify default review comment
	reviewComment := client.ReviewComments[0]
	assert.Equal(t, 654321, reviewComment.ID)
	assert.Equal(t, "Consider using a more descriptive variable name here.", reviewComment.Body)
	assert.Equal(t, "review", reviewComment.Type)
	assert.Equal(t, "reviewer2", reviewComment.User.Login)
	assert.Equal(t, "main.go", reviewComment.Path)
	assert.Equal(t, 42, reviewComment.Line)
}

func TestMockClientListIssueComments(t *testing.T) {
	client := NewMockClient()

	comments, err := client.ListIssueComments("owner", "repo", 123)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "issue", comments[0].Type)
}

func TestMockClientListIssueCommentsError(t *testing.T) {
	client := NewMockClient()
	client.ListIssueCommentsError = assert.AnError

	comments, err := client.ListIssueComments("owner", "repo", 123)
	assert.Error(t, err)
	assert.Nil(t, comments)
}

func TestMockClientListReviewComments(t *testing.T) {
	client := NewMockClient()

	comments, err := client.ListReviewComments("owner", "repo", 123)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, "review", comments[0].Type)
}

func TestMockClientListReviewCommentsError(t *testing.T) {
	client := NewMockClient()
	client.ListReviewCommentsError = assert.AnError

	comments, err := client.ListReviewComments("owner", "repo", 123)
	assert.Error(t, err)
	assert.Nil(t, comments)
}

func TestMockClientCreateIssueComment(t *testing.T) {
	client := NewMockClient()

	comment, err := client.CreateIssueComment("owner", "repo", 123, "Test comment")
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, 789012, comment.ID)
	assert.Equal(t, "Test comment", comment.Body)
	assert.Equal(t, "issue", comment.Type)
	assert.Equal(t, "testuser", comment.User.Login)
	assert.Equal(t, comment, client.CreatedComment)
}

func TestMockClientCreateIssueCommentError(t *testing.T) {
	client := NewMockClient()
	client.CreateCommentError = assert.AnError

	comment, err := client.CreateIssueComment("owner", "repo", 123, "Test comment")
	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestMockClientCreateReviewCommentReply(t *testing.T) {
	client := NewMockClient()

	comment, err := client.CreateReviewCommentReply("owner", "repo", 123456, "Reply comment")
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, 345678, comment.ID)
	assert.Equal(t, "Reply comment", comment.Body)
	assert.Equal(t, "review", comment.Type)
	assert.Equal(t, "testuser", comment.User.Login)
	assert.Equal(t, comment, client.CreatedComment)
}

func TestMockClientCreateReviewCommentReplyError(t *testing.T) {
	client := NewMockClient()
	client.CreateCommentError = assert.AnError

	comment, err := client.CreateReviewCommentReply("owner", "repo", 123456, "Reply comment")
	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestMockClientFindReviewThreadForComment(t *testing.T) {
	client := NewMockClient()

	threadID, err := client.FindReviewThreadForComment("owner", "repo", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, "RT_123", threadID)
}

func TestMockClientResolveReviewThread(t *testing.T) {
	client := NewMockClient()

	err := client.ResolveReviewThread("RT_123")
	assert.NoError(t, err)
	assert.Equal(t, "RT_123", client.ResolvedThread)
}

func TestMockClientResolveReviewThreadError(t *testing.T) {
	client := NewMockClient()
	client.ResolveThreadError = assert.AnError

	err := client.ResolveReviewThread("RT_123")
	assert.Error(t, err)
}

func TestMockClientAddReaction(t *testing.T) {
	client := NewMockClient()

	err := client.AddReaction("owner", "repo", 123456, 123, "+1")
	assert.NoError(t, err)
}

func TestMockClientRemoveReaction(t *testing.T) {
	client := NewMockClient()

	err := client.RemoveReaction("owner", "repo", 123456, 123, "+1")
	assert.NoError(t, err)
}

func TestMockClientEditComment(t *testing.T) {
	client := NewMockClient()

	err := client.EditComment("owner", "repo", 123456, 123, "Updated comment")
	assert.NoError(t, err)
}

func TestMockClientAddReviewComment(t *testing.T) {
	client := NewMockClient()

	reviewComment := ReviewCommentInput{
		Body:     "Test review comment",
		Path:     "main.go",
		Line:     42,
		CommitID: "abc123",
	}

	err := client.AddReviewComment("owner", "repo", 123, reviewComment)
	assert.NoError(t, err)
}

func TestMockClientFetchPRDiff(t *testing.T) {
	client := NewMockClient()

	diff, err := client.FetchPRDiff("owner", "repo", 123)
	assert.NoError(t, err)
	assert.NotNil(t, diff)
	assert.Len(t, diff.Files, 1)
	assert.Equal(t, "test.go", diff.Files[0].Filename)
	assert.True(t, diff.Files[0].Lines[42])
	assert.True(t, diff.Files[0].Lines[43])
}

func TestMockClientGetPRDetails(t *testing.T) {
	client := NewMockClient()

	details, err := client.GetPRDetails("owner", "repo", 123)
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, 123, details["number"])
	assert.Equal(t, "open", details["state"])
	assert.Equal(t, "Test PR", details["title"])

	// Check head structure
	head, ok := details["head"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "abc123def456", head["sha"])
}

func TestMockClientCreateReview(t *testing.T) {
	client := NewMockClient()

	review := ReviewInput{
		Body:  "LGTM",
		Event: "APPROVE",
	}

	err := client.CreateReview("owner", "repo", 123, review)
	assert.NoError(t, err)
}

func TestMockClientFindPendingReview(t *testing.T) {
	client := NewMockClient()

	// Test successful case
	client.PendingReviewID = 456
	reviewID, err := client.FindPendingReview("owner", "repo", 123)
	assert.NoError(t, err)
	assert.Equal(t, 456, reviewID)

	// Test error case
	client.FindPendingReviewError = assert.AnError
	reviewID, err = client.FindPendingReview("owner", "repo", 123)
	assert.Error(t, err)
	assert.Equal(t, 0, reviewID)
}

func TestMockClientSubmitReview(t *testing.T) {
	client := NewMockClient()

	// Test successful case
	err := client.SubmitReview("owner", "repo", 123, 456, "LGTM", "APPROVE")
	assert.NoError(t, err)
	assert.Equal(t, 456, client.SubmittedReviewID)

	// Test error case
	client.SubmitReviewError = assert.AnError
	err = client.SubmitReview("owner", "repo", 123, 789, "LGTM", "APPROVE")
	assert.Error(t, err)
	// SubmittedReviewID should still be 456 from previous call
	assert.Equal(t, 456, client.SubmittedReviewID)
}

func TestCommentStruct(t *testing.T) {
	now := time.Now()
	comment := Comment{
		ID:        123456,
		Body:      "Test comment",
		User:      User{Login: "testuser", ID: 789, AvatarURL: "https://example.com/avatar.jpg"},
		CreatedAt: now,
		UpdatedAt: now,
		Path:      "main.go",
		Line:      42,
		Position:  10,
		Type:      "review",
	}

	assert.Equal(t, 123456, comment.ID)
	assert.Equal(t, "Test comment", comment.Body)
	assert.Equal(t, "testuser", comment.User.Login)
	assert.Equal(t, 789, comment.User.ID)
	assert.Equal(t, "https://example.com/avatar.jpg", comment.User.AvatarURL)
	assert.Equal(t, now, comment.CreatedAt)
	assert.Equal(t, now, comment.UpdatedAt)
	assert.Equal(t, "main.go", comment.Path)
	assert.Equal(t, 42, comment.Line)
	assert.Equal(t, 10, comment.Position)
	assert.Equal(t, "review", comment.Type)
}

func TestReviewCommentInput(t *testing.T) {
	input := ReviewCommentInput{
		Body:      "Review comment",
		Path:      "src/main.go",
		Line:      42,
		StartLine: 40,
		Side:      "RIGHT",
		CommitID:  "abc123def456",
	}

	assert.Equal(t, "Review comment", input.Body)
	assert.Equal(t, "src/main.go", input.Path)
	assert.Equal(t, 42, input.Line)
	assert.Equal(t, 40, input.StartLine)
	assert.Equal(t, "RIGHT", input.Side)
	assert.Equal(t, "abc123def456", input.CommitID)
}

func TestReviewInput(t *testing.T) {
	review := ReviewInput{
		Body:  "Overall looks good",
		Event: "APPROVE",
		Comments: []ReviewCommentInput{
			{
				Body:     "Nice work here",
				Path:     "main.go",
				Line:     42,
				CommitID: "abc123",
			},
		},
	}

	assert.Equal(t, "Overall looks good", review.Body)
	assert.Equal(t, "APPROVE", review.Event)
	assert.Len(t, review.Comments, 1)
	assert.Equal(t, "Nice work here", review.Comments[0].Body)
}

func TestPullRequestDiff(t *testing.T) {
	diff := PullRequestDiff{
		Files: []DiffFile{
			{
				Filename: "main.go",
				Lines:    map[int]bool{42: true, 43: true, 44: false},
			},
			{
				Filename: "test.go",
				Lines:    map[int]bool{10: true, 11: true},
			},
		},
	}

	assert.Len(t, diff.Files, 2)
	assert.Equal(t, "main.go", diff.Files[0].Filename)
	assert.True(t, diff.Files[0].Lines[42])
	assert.True(t, diff.Files[0].Lines[43])
	assert.False(t, diff.Files[0].Lines[44])
	assert.Equal(t, "test.go", diff.Files[1].Filename)
	assert.True(t, diff.Files[1].Lines[10])
	assert.True(t, diff.Files[1].Lines[11])
}
