package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/silouanwright/gh-comment/internal/github"
)

// Performance optimization utilities for gh-comment

// CompiledAuthorFilter holds pre-compiled regex patterns for author matching
type CompiledAuthorFilter struct {
	pattern *regexp.Regexp
	literal string
	isRegex bool
}

// authorFilterCache caches compiled author patterns to avoid recompilation
var (
	authorFilterCache = make(map[string]*CompiledAuthorFilter)
	authorFilterMutex sync.RWMutex
)

// compileAuthorFilter creates or retrieves a cached author filter
func compileAuthorFilter(pattern string) *CompiledAuthorFilter {
	authorFilterMutex.RLock()
	if filter, exists := authorFilterCache[pattern]; exists {
		authorFilterMutex.RUnlock()
		return filter
	}
	authorFilterMutex.RUnlock()

	// Create new filter
	filter := &CompiledAuthorFilter{}

	// Check if pattern contains wildcards
	if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		// Convert shell-style wildcards to regex
		regexPattern := strings.ReplaceAll(pattern, "*", ".*")
		regexPattern = strings.ReplaceAll(regexPattern, "?", ".")
		regexPattern = "^" + regexPattern + "$"

		if compiled, err := regexp.Compile("(?i)" + regexPattern); err == nil {
			filter.pattern = compiled
			filter.isRegex = true
		} else {
			// Fallback to literal matching
			filter.literal = strings.ToLower(pattern)
			filter.isRegex = false
		}
	} else {
		// Simple literal match (case-insensitive)
		filter.literal = strings.ToLower(pattern)
		filter.isRegex = false
	}

	// Cache the filter
	authorFilterMutex.Lock()
	authorFilterCache[pattern] = filter
	authorFilterMutex.Unlock()

	return filter
}

// Matches checks if an author matches the compiled filter
func (f *CompiledAuthorFilter) Matches(author string) bool {
	if f.isRegex {
		return f.pattern.MatchString(author)
	}
	return strings.Contains(strings.ToLower(author), f.literal)
}

// OptimizedFilterContext holds pre-processed filter conditions for efficient comment filtering
type OptimizedFilterContext struct {
	// Pre-compiled filters
	AuthorFilter *CompiledAuthorFilter

	// Pre-parsed conditions
	HasAuthorFilter bool
	HasTypeFilter   bool
	HasDateFilter   bool

	// Filter values (pre-processed)
	TargetType   string
	SinceTime    *time.Time
	UntilTime    *time.Time
	StatusFilter string

	// Performance flags
	EarlyExitOnAuthor bool
	EarlyExitOnType   bool
	EarlyExitOnDate   bool
}

// CreateOptimizedFilterContext pre-processes all filter conditions for maximum performance
func CreateOptimizedFilterContext() *OptimizedFilterContext {
	ctx := &OptimizedFilterContext{}

	// Process author filter
	if author != "" {
		ctx.AuthorFilter = compileAuthorFilter(author)
		ctx.HasAuthorFilter = true
		// Author filtering is usually highly selective
		ctx.EarlyExitOnAuthor = true
	}

	// Process type filter
	if listType != "" && listType != "all" {
		ctx.TargetType = listType
		ctx.HasTypeFilter = true
		// Type filtering is very fast and selective
		ctx.EarlyExitOnType = true
	}

	// Process date filters
	if sinceTime != nil || untilTime != nil {
		ctx.SinceTime = sinceTime
		ctx.UntilTime = untilTime
		ctx.HasDateFilter = true
		// Date filtering is moderately selective
		ctx.EarlyExitOnDate = true
	}

	// Process status filter
	if status != "" && status != "all" {
		ctx.StatusFilter = status
	}

	return ctx
}

// FastFilterComments provides optimized comment filtering with early exits and pre-compiled patterns
func (ctx *OptimizedFilterContext) FastFilterComments(comments []Comment) []Comment {
	if len(comments) == 0 {
		return comments
	}

	// Pre-allocate result slice with reasonable capacity
	// Use len/2 as a reasonable estimate for filtered results
	estimatedSize := len(comments) / 2
	if estimatedSize < 10 {
		estimatedSize = len(comments) // For small lists, allocate full size
	}
	filtered := make([]Comment, 0, estimatedSize)

	for i := range comments {
		comment := &comments[i] // Use pointer to avoid copying

		// Apply filters in order of selectivity (most selective first)

		// 1. Type filter (very fast, highly selective)
		if ctx.HasTypeFilter && ctx.EarlyExitOnType {
			if comment.Type != ctx.TargetType {
				continue
			}
		}

		// 2. Author filter (pre-compiled pattern, moderately selective)
		if ctx.HasAuthorFilter && ctx.EarlyExitOnAuthor {
			if !ctx.AuthorFilter.Matches(comment.Author) {
				continue
			}
		}

		// 3. Date filters (time comparison, moderately selective)
		if ctx.HasDateFilter && ctx.EarlyExitOnDate {
			if ctx.SinceTime != nil && comment.CreatedAt.Before(*ctx.SinceTime) {
				continue
			}
			if ctx.UntilTime != nil && comment.CreatedAt.After(*ctx.UntilTime) {
				continue
			}
		}

		// 4. Status filter (less selective, applied last)
		if ctx.StatusFilter != "" {
			// This is a placeholder for actual resolution status filtering
			// In a real implementation, you'd check comment.ResolvedAt or similar
			if ctx.StatusFilter == "resolved" {
				// Skip for now as we don't have resolution data
				continue
			}
		}

		// Comment passed all filters - add to results
		filtered = append(filtered, *comment)
	}

	return filtered
}

// ConcurrentCommentFetcher enables concurrent fetching of issue and review comments
type ConcurrentCommentFetcher struct {
	IssueComments  []Comment
	ReviewComments []Comment
	IssueError     error
	ReviewError    error
	wg             sync.WaitGroup
}

// FetchCommentsAsync fetches issue and review comments concurrently for better performance
func FetchCommentsAsync(client github.GitHubAPI, owner, repoName string, pr int, timeout time.Duration) ([]Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fetcher := &ConcurrentCommentFetcher{}

	// Fetch issue comments in goroutine
	fetcher.wg.Add(1)
	go func() {
		defer fetcher.wg.Done()

		select {
		case <-ctx.Done():
			fetcher.IssueError = ctx.Err()
			return
		default:
		}

		issueComments, err := client.ListIssueComments(owner, repoName, pr)
		if err != nil {
			fetcher.IssueError = err
			return
		}

		// Convert to Comment structs
		fetcher.IssueComments = make([]Comment, len(issueComments))
		for i, comment := range issueComments {
			fetcher.IssueComments[i] = Comment{
				ID:        comment.ID,
				Author:    comment.User.Login,
				Body:      comment.Body,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
				Type:      "issue",
			}
		}
	}()

	// Fetch review comments in goroutine
	fetcher.wg.Add(1)
	go func() {
		defer fetcher.wg.Done()

		select {
		case <-ctx.Done():
			fetcher.ReviewError = ctx.Err()
			return
		default:
		}

		reviewComments, err := client.ListReviewComments(owner, repoName, pr)
		if err != nil {
			fetcher.ReviewError = err
			return
		}

		// Convert to Comment structs
		fetcher.ReviewComments = make([]Comment, len(reviewComments))
		for i, comment := range reviewComments {
			fetcher.ReviewComments[i] = Comment{
				ID:        comment.ID,
				Author:    comment.User.Login,
				Body:      comment.Body,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
				Path:      comment.Path,
				Line:      comment.Line,
				CommitID:  comment.CommitID,
				Type:      "review",
			}
		}
	}()

	// Wait for both goroutines to complete
	fetcher.wg.Wait()

	// Check for errors
	if fetcher.IssueError != nil {
		return nil, formatActionableError("issue comments fetch", fetcher.IssueError)
	}
	if fetcher.ReviewError != nil {
		return nil, formatActionableError("review comments fetch", fetcher.ReviewError)
	}

	// Combine results
	totalComments := len(fetcher.IssueComments) + len(fetcher.ReviewComments)
	allComments := make([]Comment, 0, totalComments)
	allComments = append(allComments, fetcher.IssueComments...)
	allComments = append(allComments, fetcher.ReviewComments...)

	return allComments, nil
}

// StringPool provides string deduplication to reduce memory usage for repeated values
type StringPool struct {
	pool map[string]string
	mu   sync.RWMutex
}

// NewStringPool creates a new string pool for deduplication
func NewStringPool() *StringPool {
	return &StringPool{
		pool: make(map[string]string),
	}
}

// Intern returns a canonical version of the string, deduplicating memory usage
func (sp *StringPool) Intern(s string) string {
	if s == "" {
		return s
	}

	sp.mu.RLock()
	if canonical, exists := sp.pool[s]; exists {
		sp.mu.RUnlock()
		return canonical
	}
	sp.mu.RUnlock()

	sp.mu.Lock()
	// Double-check after acquiring write lock
	if canonical, exists := sp.pool[s]; exists {
		sp.mu.Unlock()
		return canonical
	}

	// Create canonical version
	canonical := s
	sp.pool[s] = canonical
	sp.mu.Unlock()

	return canonical
}

// OptimizeCommentsMemory reduces memory usage by deduplicating common strings
func OptimizeCommentsMemory(comments []Comment) []Comment {
	if len(comments) <= 1 {
		return comments
	}

	pool := NewStringPool()

	for i := range comments {
		// Intern common strings that are likely to be repeated
		comments[i].Author = pool.Intern(comments[i].Author)
		comments[i].Type = pool.Intern(comments[i].Type)
		comments[i].Path = pool.Intern(comments[i].Path)
		comments[i].State = pool.Intern(comments[i].State)

		// For bodies, only intern if they're short (likely to be repeated)
		if len(comments[i].Body) < 100 {
			comments[i].Body = pool.Intern(comments[i].Body)
		}
	}

	return comments
}

// PerformanceMetrics tracks performance statistics for optimization monitoring
type PerformanceMetrics struct {
	FetchDuration    time.Duration
	FilterDuration   time.Duration
	DisplayDuration  time.Duration
	TotalComments    int
	FilteredComments int
	MemoryOptimized  bool
}

// LogPerformanceMetrics logs performance statistics if verbose mode is enabled
func LogPerformanceMetrics(metrics PerformanceMetrics) {
	if !verbose {
		return
	}

	fmt.Printf("Performance metrics:\n")
	fmt.Printf("  Fetch time: %v\n", metrics.FetchDuration)
	fmt.Printf("  Filter time: %v\n", metrics.FilterDuration)
	fmt.Printf("  Display time: %v\n", metrics.DisplayDuration)
	fmt.Printf("  Total comments: %d\n", metrics.TotalComments)
	fmt.Printf("  Filtered to: %d\n", metrics.FilteredComments)
	fmt.Printf("  Memory optimized: %v\n", metrics.MemoryOptimized)
	fmt.Printf("  Filter efficiency: %.1f%%\n",
		float64(metrics.FilteredComments)/float64(metrics.TotalComments)*100)
}
