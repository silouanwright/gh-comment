package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
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

// BenchmarkExpandSuggestions benchmarks the suggestion parsing performance
func BenchmarkExpandSuggestions(b *testing.B) {
	testCases := []string{
		"Simple comment without suggestions",
		"[SUGGEST: optimized code here]",
		"Multiple suggestions: [SUGGEST: first suggestion] and [SUGGEST: second suggestion]",
		"<<<SUGGEST>>>old code>>>new code<<<SUGGEST>>>",
		"Mixed suggestions: [SUGGEST: simple] and <<<SUGGEST>>>old>>>new<<<SUGGEST>>>",
		"Very long suggestion: [SUGGEST: " + generateLongString(500) + "]",
		"No suggestions but contains brackets [not a suggestion] and <<<not a suggestion>>> text",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			expandSuggestions(testCase)
		}
	}
}

// BenchmarkValidateCommentBody benchmarks comment body validation performance
func BenchmarkValidateCommentBody(b *testing.B) {
	testCases := []string{
		"Short comment",
		"Medium length comment with some details and explanations",
		generateLongString(1000),  // 1KB comment
		generateLongString(10000), // 10KB comment
		generateLongString(50000), // 50KB comment (near GitHub limit)
		"Comment with special characters: áéíóú çñ 中文 🚀 🎉",
		"Comment with HTML tags: <script>alert('test')</script>",
		"Comment with markdown: **bold** *italic* `code` [link](url)",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			validateCommentBody(testCase)
		}
	}
}

// BenchmarkParseBatchConfig benchmarks YAML parsing performance for batch configurations
func BenchmarkParseBatchConfig(b *testing.B) {
	// Create test YAML files with different sizes
	testConfigs := []string{
		// Small config
		`pr: 123
comments:
  - file: test.go
    line: 10
    message: "Small test comment"`,
		// Medium config
		`pr: 123
review:
  body: "Review body"
  event: "COMMENT"
comments:
  - file: test1.go
    line: 10
    message: "Comment 1"
  - file: test2.go
    line: 20
    message: "Comment 2"
  - file: test3.go
    line: 30
    message: "Comment 3"`,
		// Large config with many comments
		generateLargeBatchConfig(50),
	}

	// Create temporary files for benchmarking
	var tempFiles []string
	for _, config := range testConfigs {
		tmpFile, err := ioutil.TempFile("", "benchmark_config_*.yaml")
		if err != nil {
			b.Fatalf("Failed to create temp file: %v", err)
		}
		if _, err := tmpFile.WriteString(config); err != nil {
			b.Fatalf("Failed to write config: %v", err)
		}
		tmpFile.Close()
		tempFiles = append(tempFiles, tmpFile.Name())
	}

	// Clean up temp files after benchmark
	defer func() {
		for _, file := range tempFiles {
			os.Remove(file)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, configFile := range tempFiles {
			readBatchConfig(configFile)
		}
	}
}

// BenchmarkFormatActionableError benchmarks error formatting performance
func BenchmarkFormatActionableError(b *testing.B) {
	testErrors := []error{
		fmt.Errorf("404 Not Found"),
		fmt.Errorf("422 Validation Failed"),
		fmt.Errorf("403 Forbidden"),
		fmt.Errorf("rate limit exceeded"),
		fmt.Errorf("network timeout"),
		fmt.Errorf("authentication failed"),
		fmt.Errorf("API rate limit exceeded. Try again in 60 minutes"),
		fmt.Errorf("Resource not accessible by integration"),
		fmt.Errorf("GraphQL: Field 'invalidField' doesn't exist on type 'Comment'"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, err := range testErrors {
			formatActionableError("test-operation", err)
		}
	}
}

// BenchmarkParsePositiveInt benchmarks integer parsing performance
func BenchmarkParsePositiveInt(b *testing.B) {
	testCases := []string{
		"1",
		"123",
		"999999",
		"0",
		"-1",
		"invalid",
		"123abc",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			parsePositiveInt(testCase, "test")
		}
	}
}

// BenchmarkColorizeSuccess benchmarks colorization performance
func BenchmarkColorizeSuccess(b *testing.B) {
	testMessages := []string{
		"Success",
		"Successfully created comment",
		"Operation completed successfully with detailed information",
		generateLongString(200),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, msg := range testMessages {
			ColorizeSuccess(msg)
		}
	}
}

// Helper functions for benchmark data generation

func generateLongString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte('a' + (i % 26))
	}
	return string(result)
}

func generateLargeBatchConfig(commentCount int) string {
	config := `pr: 123
review:
  body: "Large batch review"
  event: "COMMENT"
comments:
`
	for i := 0; i < commentCount; i++ {
		config += `  - file: test` + string(rune('0'+i%10)) + `.go
    line: ` + string(rune('1'+i%9)) + `0
    message: "Comment ` + string(rune('0'+i%10)) + ` with some detailed explanation"
`
	}
	return config
}
