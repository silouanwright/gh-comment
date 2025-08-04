package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCompiledAuthorFilter(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		author   string
		expected bool
	}{
		{
			name:     "exact match",
			pattern:  "john-doe",
			author:   "john-doe",
			expected: true,
		},
		{
			name:     "case insensitive match",
			pattern:  "John-Doe",
			author:   "john-doe",
			expected: true,
		},
		{
			name:     "wildcard prefix",
			pattern:  "senior-*",
			author:   "senior-dev",
			expected: true,
		},
		{
			name:     "wildcard suffix",
			pattern:  "*@company.com",
			author:   "user@company.com",
			expected: true,
		},
		{
			name:     "wildcard middle",
			pattern:  "team-*-lead",
			author:   "team-frontend-lead",
			expected: true,
		},
		{
			name:     "no match",
			pattern:  "admin*",
			author:   "regular-user",
			expected: false,
		},
		{
			name:     "partial literal match",
			pattern:  "dev",
			author:   "senior-dev",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := compileAuthorFilter(tt.pattern)
			result := filter.Matches(tt.author)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthorFilterCaching(t *testing.T) {
	// Clear cache
	authorFilterMutex.Lock()
	authorFilterCache = make(map[string]*CompiledAuthorFilter)
	authorFilterMutex.Unlock()

	pattern := "test-pattern*"
	
	// First call should create and cache
	filter1 := compileAuthorFilter(pattern)
	assert.NotNil(t, filter1)
	
	// Second call should return cached version
	filter2 := compileAuthorFilter(pattern)
	assert.Equal(t, filter1, filter2, "Should return same cached filter instance")
	
	// Verify cache contains the pattern
	authorFilterMutex.RLock()
	_, exists := authorFilterCache[pattern]
	authorFilterMutex.RUnlock()
	assert.True(t, exists, "Pattern should be cached")
}

func TestOptimizedFilterContext(t *testing.T) {
	// Save original values
	originalAuthor := author
	originalListType := listType
	originalSinceTime := sinceTime
	originalUntilTime := untilTime
	originalStatus := status
	
	defer func() {
		author = originalAuthor
		listType = originalListType
		sinceTime = originalSinceTime
		untilTime = originalUntilTime
		status = originalStatus
	}()

	// Set up test conditions
	author = "test-user*"
	listType = "review"
	now := time.Now()
	sinceTime = &now
	status = "open"

	ctx := CreateOptimizedFilterContext()
	
	assert.True(t, ctx.HasAuthorFilter, "Should detect author filter")
	assert.True(t, ctx.HasTypeFilter, "Should detect type filter")  
	assert.True(t, ctx.HasDateFilter, "Should detect date filter")
	assert.Equal(t, "review", ctx.TargetType, "Should set target type")
	assert.Equal(t, "open", ctx.StatusFilter, "Should set status filter")
	assert.NotNil(t, ctx.AuthorFilter, "Should create author filter")
}

func TestFastFilterComments(t *testing.T) {
	// Create test comments
	comments := []Comment{
		{
			ID:        1,
			Author:    "test-user-1",
			Type:      "issue",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        2,
			Author:    "test-user-2", 
			Type:      "review",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        3,
			Author:    "other-user",
			Type:      "issue",
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
		{
			ID:        4,
			Author:    "test-user-3",
			Type:      "review",
			CreatedAt: time.Now().Add(-3 * time.Hour),
		},
	}

	// Test author filter
	author = "test-user*"
	listType = ""
	sinceTime = nil
	untilTime = nil
	status = ""
	
	ctx := CreateOptimizedFilterContext()
	filtered := ctx.FastFilterComments(comments)
	
	assert.Len(t, filtered, 3, "Should filter to test-user* authors")
	for _, comment := range filtered {
		assert.True(t, strings.HasPrefix(comment.Author, "test-user"), 
			"All results should match author filter")
	}

	// Test type filter
	author = ""
	listType = "review"
	
	ctx = CreateOptimizedFilterContext()
	filtered = ctx.FastFilterComments(comments)
	
	assert.Len(t, filtered, 2, "Should filter to review comments only")
	for _, comment := range filtered {
		assert.Equal(t, "review", comment.Type, "All results should be review type")
	}

	// Test combined filters
	author = "test-user*" 
	listType = "review"
	
	ctx = CreateOptimizedFilterContext()
	filtered = ctx.FastFilterComments(comments)
	
	assert.Len(t, filtered, 2, "Should match both author and type filters")
	for _, comment := range filtered {
		assert.True(t, strings.HasPrefix(comment.Author, "test-user"))
		assert.Equal(t, "review", comment.Type)
	}
}

func TestStringPool(t *testing.T) {
	pool := NewStringPool()
	
	// Test basic interning
	s1 := pool.Intern("test-string")
	s2 := pool.Intern("test-string")
	
	// Should return same string value
	assert.Equal(t, s1, s2, "Should return same string value for identical strings")
	
	// Test different strings
	s3 := pool.Intern("different-string")
	assert.NotEqual(t, s1, s3, "Should return different values for different strings")
	
	// Test empty string
	empty1 := pool.Intern("")
	empty2 := pool.Intern("")
	assert.Equal(t, "", empty1)
	assert.Equal(t, "", empty2)
	
	// Test that the pool actually stores strings
	pool.mu.RLock()
	_, exists := pool.pool["test-string"]
	pool.mu.RUnlock()
	assert.True(t, exists, "String should be stored in pool")
}

func TestOptimizeCommentsMemory(t *testing.T) {
	// Create comments with duplicate strings
	comments := []Comment{
		{
			ID:     1,
			Author: "john-doe",
			Type:   "issue",
			Path:   "src/main.go",
			Body:   "LGTM",
		},
		{
			ID:     2,
			Author: "john-doe", // Duplicate
			Type:   "review",   // Different
			Path:   "src/main.go", // Duplicate
			Body:   "LGTM",        // Duplicate
		},
		{
			ID:     3,
			Author: "jane-smith",
			Type:   "issue", // Duplicate
			Path:   "src/utils.go",
			Body:   "Please fix this issue",
		},
	}

	optimized := OptimizeCommentsMemory(comments)
	
	// Should return same slice length
	assert.Len(t, optimized, 3)
	
	// Verify string interning worked (same authors should have same pointer)
	// Note: This is hard to test directly, but we can verify the function runs without error
	assert.Equal(t, "john-doe", optimized[0].Author)
	assert.Equal(t, "john-doe", optimized[1].Author)
	assert.Equal(t, "issue", optimized[0].Type)  
	assert.Equal(t, "issue", optimized[2].Type)
}

func TestPerformanceMetricsLogging(t *testing.T) {
	// Save original verbose state
	originalVerbose := verbose
	defer func() { verbose = originalVerbose }()

	metrics := PerformanceMetrics{
		FetchDuration:    100 * time.Millisecond,
		FilterDuration:   10 * time.Millisecond,
		DisplayDuration:  5 * time.Millisecond,
		TotalComments:    100,
		FilteredComments: 25,
		MemoryOptimized:  true,
	}

	// Test with verbose disabled
	verbose = false
	assert.NotPanics(t, func() {
		LogPerformanceMetrics(metrics)
	}, "Should not panic when verbose is disabled")

	// Test with verbose enabled  
	verbose = true
	assert.NotPanics(t, func() {
		LogPerformanceMetrics(metrics)
	}, "Should not panic when verbose is enabled")
}

// Benchmark tests to measure performance improvements
func BenchmarkAuthorFilterCompilation(b *testing.B) {
	patterns := []string{
		"john-doe",
		"senior-*",
		"*@company.com", 
		"team-*-lead",
		"user*",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pattern := patterns[i%len(patterns)]
		compileAuthorFilter(pattern)
	}
}

func BenchmarkAuthorFilterMatching(b *testing.B) {
	filter := compileAuthorFilter("senior-*")
	authors := []string{
		"senior-dev",
		"senior-architect", 
		"junior-dev",
		"lead-engineer",
		"senior-qa",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		author := authors[i%len(authors)]
		filter.Matches(author)
	}
}

func BenchmarkFilterComments(b *testing.B) {
	// Create large comment set for benchmarking
	comments := make([]Comment, 1000)
	for i := 0; i < 1000; i++ {
		comments[i] = Comment{
			ID:        i,
			Author:    []string{"senior-dev", "junior-dev", "architect", "qa-lead", "designer"}[i%5],
			Type:      []string{"issue", "review"}[i%2],
			CreatedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	// Set up filter context
	author = "senior-*"
	listType = "review"
	sinceTime = nil
	untilTime = nil
	
	ctx := CreateOptimizedFilterContext()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.FastFilterComments(comments)
	}
}

func BenchmarkStringPoolInterning(b *testing.B) {
	pool := NewStringPool()
	strings := []string{
		"john-doe",
		"jane-smith", 
		"senior-dev",
		"junior-dev",
		"architect",
		"john-doe", // Repeat
		"jane-smith", // Repeat
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings[i%len(strings)]
		pool.Intern(s)
	}
}