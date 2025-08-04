package cmd

import (
	"testing"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
)

func BenchmarkListComments(b *testing.B) {
	// Create mock comments
	comments := make([]Comment, MaxGraphQLResults)
	for i := 0; i < MaxGraphQLResults; i++ {
		comments[i] = Comment{
			ID:        i + 1,
			Author:    "testuser",
			Body:      "This is a test comment with some content to benchmark",
			CreatedAt: time.Now(),
			Type:      "issue",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := filterComments(comments)
		_ = filtered // Prevent optimization
	}
}

func BenchmarkFormatTimeAgo(b *testing.B) {
	now := time.Now()
	times := []time.Time{
		now.Add(-30 * time.Second),
		now.Add(-5 * time.Minute),
		now.Add(-2 * time.Hour),
		now.Add(-24 * time.Hour),
		now.Add(-7 * 24 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, t := range times {
			_ = formatTimeAgo(t)
		}
	}
}

func BenchmarkMockClientOperations(b *testing.B) {
	mockClient := github.NewMockClient()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mockClient.ListIssueComments("owner", "repo", 1)
		_, _ = mockClient.ListReviewComments("owner", "repo", 1)
	}
}
