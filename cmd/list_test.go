package cmd

import (
	"testing"

	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllComments(t *testing.T) {
	// Create a mock client
	mockClient := &github.MockClient{
		IssueComments: []github.Comment{
			{
				ID:   1,
				Body: "Test issue comment",
				User: github.User{Login: "testuser"},
			},
		},
		ReviewComments: []github.Comment{
			{
				ID:   2,
				Body: "Test review comment",
				User: github.User{Login: "reviewer"},
				Path: "test.go",
				Line: 42,
			},
		},
	}

	comments, err := fetchAllComments(mockClient, "owner/repo", 123)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)

	// Check issue comment
	assert.Equal(t, 1, comments[0].ID)
	assert.Equal(t, "Test issue comment", comments[0].Body)
	assert.Equal(t, "testuser", comments[0].Author)
	assert.Equal(t, "issue", comments[0].Type)

	// Check review comment
	assert.Equal(t, 2, comments[1].ID)
	assert.Equal(t, "Test review comment", comments[1].Body)
	assert.Equal(t, "reviewer", comments[1].Author)
	assert.Equal(t, "review", comments[1].Type)
	assert.Equal(t, "test.go", comments[1].Path)
	assert.Equal(t, 42, comments[1].Line)
}

func TestFetchAllCommentsInvalidRepo(t *testing.T) {
	mockClient := &github.MockClient{}

	_, err := fetchAllComments(mockClient, "invalid", 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid repository format")
}
