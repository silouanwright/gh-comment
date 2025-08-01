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
	quiet bool
	hideAuthors bool
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
  # List all comments on PR 123 (shows URLs and IDs by default)
  gh comment list 123

  # Minimal output for human reading
  gh comment list 123 --quiet

  # List comments from specific author
  gh comment list 123 --author octocat

  # Hide author names for privacy
  gh comment list 123 --hide-authors

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
	listCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output without URLs and IDs (default shows full context for AI)")
	listCmd.Flags().BoolVar(&hideAuthors, "hide-authors", false, "Hide author names for privacy")
}

func runList(cmd *cobra.Command, args []string) error {
	var pr int
	var err error

	// Parse PR argument
	if len(args) == 1 {
		pr, err = strconv.Atoi(args[0])
		if err != nil {
			return formatValidationError("PR number", args[0], "must be a valid integer")
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
		fmt.Printf("Quiet mode: %v\n", quiet)
		fmt.Printf("Hide authors: %v\n", hideAuthors)
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
	DiffHunk  string `json:"diff_hunk,omitempty"`

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

	// Fetch reviews (parent comments that group line-specific comments)
	var reviews []struct {
		ID        int       `json:"id"`
		User      struct {
			Login string `json:"login"`
		} `json:"user"`
		Body      string    `json:"body"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"submitted_at"`
		HTMLURL   string    `json:"html_url"`
	}

	err = client.Get(fmt.Sprintf("repos/%s/pulls/%d/reviews", repo, pr), &reviews)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}

	// Convert reviews (only include ones with body text)
	for _, review := range reviews {
		if strings.TrimSpace(review.Body) != "" {
			allComments = append(allComments, Comment{
				ID:        review.ID,
				Author:    review.User.Login,
				Body:      review.Body,
				CreatedAt: review.CreatedAt,
				HTMLURL:   review.HTMLURL,
				Type:      "review",
				State:     review.State,
			})
		}
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
		DiffHunk  string    `json:"diff_hunk"`
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
			DiffHunk:  comment.DiffHunk,
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
	var issueComments, reviewComments, lineComments []Comment
	for _, comment := range comments {

		if comment.Type == "issue" {
			issueComments = append(issueComments, comment)
		} else if comment.Type == "review" {
			reviewComments = append(reviewComments, comment)
		} else {
			lineComments = append(lineComments, comment)
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

	// Display review-level comments (parent comments that group line-specific ones)
	if len(reviewComments) > 0 {
		fmt.Printf("📋 Review Comments (%d)\n", len(reviewComments))
		fmt.Println(strings.Repeat("─", 50))
		for i, comment := range reviewComments {
			displayComment(comment, i+1)
		}
		fmt.Println()
	}

	// Display line-specific comments
	if len(lineComments) > 0 {
		fmt.Printf("📍 Line-Specific Comments (%d)\n", len(lineComments))
		fmt.Println(strings.Repeat("─", 50))
		for i, comment := range lineComments {
			displayComment(comment, i+1)
		}
	}
}

func displayComment(comment Comment, index int) {
	// Header with author and timestamp
	timeAgo := formatTimeAgo(comment.CreatedAt)
	if hideAuthors {
		fmt.Printf("[%d] 👤 [hidden] • %s", index, timeAgo)
	} else {
		fmt.Printf("[%d] 👤 %s • %s", index, comment.Author, timeAgo)
	}

	// Show review state for review-level comments
	if comment.Type == "review" && comment.State != "" {
		stateEmoji := "📝"
		switch comment.State {
		case "approved":
			stateEmoji = "✅"
		case "changes_requested":
			stateEmoji = "🔴"
		case "commented":
			stateEmoji = "💬"
		}
		fmt.Printf(" %s %s", stateEmoji, comment.State)
	}
	fmt.Println()

	// File and line info for line-specific comments
	if comment.Path != "" {
		lineInfo := fmt.Sprintf("L%d", comment.Line)
		if comment.StartLine > 0 && comment.StartLine != comment.Line {
			lineInfo = fmt.Sprintf("L%d-L%d", comment.StartLine, comment.Line)
		}
		fmt.Printf("📁 %s:%s\n", comment.Path, lineInfo)

		// Show the actual diff context if available
		if comment.DiffHunk != "" {
			fmt.Printf("📝 Code Context:\n")
			displayDiffHunk(comment.DiffHunk)
		}
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

	// Show URLs by default (AI-friendly), hide only in quiet mode
	if !quiet {
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

func displayDiffHunk(diffHunk string) {
	// Split diff hunk into lines
	lines := strings.Split(strings.TrimSpace(diffHunk), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Color code the diff lines
		if strings.HasPrefix(line, "@@") {
			// Diff header - show line numbers
			fmt.Printf("   🔹 %s\n", line)
		} else if strings.HasPrefix(line, "+") {
			// Added line - green
			fmt.Printf("   ➕ %s\n", line)
		} else if strings.HasPrefix(line, "-") {
			// Removed line - red
			fmt.Printf("   ➖ %s\n", line)
		} else {
			// Context line - neutral
			fmt.Printf("     %s\n", line)
		}
	}
	fmt.Println()
}
