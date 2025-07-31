package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var (
	showResolved bool
	onlyUnresolved bool
	author string
)

var listCmd = &cobra.Command{
	Use:   "list [pr]",
	Short: "List all comments on a PR",
	Long: `List all comments on a pull request, including both general PR comments and line-specific review comments.

This is perfect for the workflow where you:
1. Create a PR
2. Someone reviews it and adds comments
3. You run 'gh comment list' to see all feedback
4. You address the comments and push fixes

Shows key information for each comment:
- Author and timestamp
- File and line (for line-specific comments)
- Comment body
- Resolution status

Examples:
  # List all comments on PR 123
  gh comment list 123
  
  # List only unresolved comments
  gh comment list 123 --unresolved
  
  # List comments from specific author
  gh comment list 123 --author octocat
  
  # Auto-detect PR from current branch
  gh comment list`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	
	listCmd.Flags().BoolVar(&showResolved, "resolved", false, "Include resolved comments")
	listCmd.Flags().BoolVar(&onlyUnresolved, "unresolved", false, "Show only unresolved comments")
	listCmd.Flags().StringVar(&author, "author", "", "Filter comments by author")
}

func runList(cmd *cobra.Command, args []string) error {
	var pr int
	var err error

	// Parse PR argument
	if len(args) == 1 {
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid PR number: %s", args[0])
		}
	} else {
		// Auto-detect PR
		pr, err = getCurrentPR()
		if err != nil {
			return err
		}
	}

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR: %d\n", pr)
		fmt.Printf("Show resolved: %v\n", showResolved)
		fmt.Printf("Only unresolved: %v\n", onlyUnresolved)
		if author != "" {
			fmt.Printf("Filter by author: %s\n", author)
		}
		fmt.Println()
	}

	// Fetch comments
	comments, err := fetchAllComments(repository, pr)
	if err != nil {
		return err
	}

	// Filter comments
	filteredComments := filterComments(comments)

	// Display comments
	displayComments(filteredComments, pr)

	return nil
}

type Comment struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	HTMLURL   string    `json:"html_url"`
	
	// For line-specific comments
	Path      string `json:"path,omitempty"`
	Line      int    `json:"line,omitempty"`
	StartLine int    `json:"start_line,omitempty"`
	
	// Comment type
	Type string `json:"type"` // "issue" or "review"
	
	// Resolution status (for review comments)
	State string `json:"state,omitempty"` // "pending", "submitted", etc.
}

func fetchAllComments(repo string, pr int) ([]Comment, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	var allComments []Comment

	// Fetch general PR comments (issue comments)
	var issueComments []struct {
		ID        int       `json:"id"`
		User      struct {
			Login string `json:"login"`
		} `json:"user"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		HTMLURL   string    `json:"html_url"`
	}

	err = client.Get(fmt.Sprintf("repos/%s/issues/%d/comments", repo, pr), &issueComments)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue comments: %w", err)
	}

	// Convert issue comments
	for _, comment := range issueComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			HTMLURL:   comment.HTMLURL,
			Type:      "issue",
		})
	}

	// Fetch review comments (line-specific)
	var reviewComments []struct {
		ID        int       `json:"id"`
		User      struct {
			Login string `json:"login"`
		} `json:"user"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		HTMLURL   string    `json:"html_url"`
		Path      string    `json:"path"`
		Line      int       `json:"line"`
		StartLine int       `json:"start_line"`
	}

	err = client.Get(fmt.Sprintf("repos/%s/pulls/%d/comments", repo, pr), &reviewComments)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review comments: %w", err)
	}

	// Convert review comments
	for _, comment := range reviewComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			HTMLURL:   comment.HTMLURL,
			Path:      comment.Path,
			Line:      comment.Line,
			StartLine: comment.StartLine,
			Type:      "review",
		})
	}

	return allComments, nil
}

func filterComments(comments []Comment) []Comment {
	var filtered []Comment

	for _, comment := range comments {
		// Filter by author if specified
		if author != "" && comment.Author != author {
			continue
		}

		// TODO: Add resolution status filtering when we implement it
		// For now, we don't have resolution status from the API
		
		filtered = append(filtered, comment)
	}

	return filtered
}

func displayComments(comments []Comment, pr int) {
	if len(comments) == 0 {
		fmt.Printf("No comments found on PR #%d\n", pr)
		return
	}

	fmt.Printf("📝 Comments on PR #%d (%d total)\n\n", pr, len(comments))

	// Group comments by type
	var issueComments, reviewComments []Comment
	for _, comment := range comments {
		if comment.Type == "issue" {
			issueComments = append(issueComments, comment)
		} else {
			reviewComments = append(reviewComments, comment)
		}
	}

	// Display general PR comments
	if len(issueComments) > 0 {
		fmt.Printf("💬 General PR Comments (%d)\n", len(issueComments))
		fmt.Println(strings.Repeat("─", 50))
		for i, comment := range issueComments {
			displayComment(comment, i+1)
		}
		fmt.Println()
	}

	// Display line-specific comments
	if len(reviewComments) > 0 {
		fmt.Printf("📍 Line-Specific Comments (%d)\n", len(reviewComments))
		fmt.Println(strings.Repeat("─", 50))
		for i, comment := range reviewComments {
			displayComment(comment, i+1)
		}
	}
}

func displayComment(comment Comment, index int) {
	// Header with author and timestamp
	timeAgo := formatTimeAgo(comment.CreatedAt)
	fmt.Printf("[%d] 👤 %s • %s\n", index, comment.Author, timeAgo)

	// File and line info for review comments
	if comment.Type == "review" && comment.Path != "" {
		lineInfo := fmt.Sprintf("L%d", comment.Line)
		if comment.StartLine > 0 && comment.StartLine != comment.Line {
			lineInfo = fmt.Sprintf("L%d-L%d", comment.StartLine, comment.Line)
		}
		fmt.Printf("📁 %s:%s\n", comment.Path, lineInfo)
	}

	// Comment body (truncate if too long)
	body := strings.TrimSpace(comment.Body)
	if len(body) > 200 {
		body = body[:197] + "..."
	}
	
	// Indent the comment body
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		fmt.Printf("   %s\n", line)
	}

	if verbose {
		fmt.Printf("   🔗 %s\n", comment.HTMLURL)
	}

	fmt.Println()
}

func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		} else if days < 7 {
			return fmt.Sprintf("%d days ago", days)
		} else {
			return t.Format("Jan 2, 2006")
		}
	}
}
