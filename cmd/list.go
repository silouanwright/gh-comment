package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/markusmobius/go-dateparser"
	"github.com/silouanwright/gh-comment/internal/github"
	"github.com/spf13/cobra"
)

var (
	showResolved   bool
	onlyUnresolved bool
	author         string
	quiet          bool
	hideAuthors    bool

	// Advanced filtering flags
	status   string
	since    string
	until    string
	resolved string
	listType string

	// Output format flags
	outputFormat string
	idsOnly      bool

	// Parsed time values
	sinceTime *time.Time
	untilTime *time.Time

	// Client for dependency injection (tests can override)
	listClient github.GitHubAPI
)

var listCmd = &cobra.Command{
	Use:   "list [pr]",
	Short: "List comments with advanced filtering and formatting options",
	Long: heredoc.Doc(`
		List all comments on a pull request with powerful filtering capabilities.

		Comment Types:
		- Issue comments: General PR discussion, appear in main conversation tab
		- Review comments: Line-specific feedback, appear in "Files Changed" tab

		Comments can be filtered by type, author, date range, resolution status, and more.
		Output can be formatted as tables, JSON, or plain text with color coding.
		Perfect for code review workflows, comment analysis, and automation.
	`),
	Example: heredoc.Doc(`
		# Review team analysis and metrics
		$ gh comment list 123 --author "senior-dev*" --status open --since "1 week ago"
		$ gh comment list 123 --type review --author "*@company.com" --since "deployment-date"

		# Security audit and compliance tracking
		$ gh comment list 123 --author "security-team*" --since "2024-01-01" --type review
		$ gh comment list 123 --author "bot*" --since "3 days ago" --quiet

		# Structured output for automation
		$ gh comment list 123 --format json | jq '.comments[].id'
		$ gh comment list 123 --ids-only | xargs -I {} gh comment resolve {}
		$ gh comment list 123 --format json --author "security*" > security-comments.json

		# Code review workflow optimization
		$ gh comment list 123 --status open --since "sprint-start" --author "lead*"
		$ gh comment list 123 --until "release-date" --type issue --status resolved

		# Team communication patterns
		$ gh comment list 123 --author "qa*" --since "last-deployment" --type review
		$ gh comment list 123 --author "*@contractor.com" --status open --since "1 month ago"

		# Blocker identification and resolution tracking
		$ gh comment list 123 --author "architect*" --status open --type review
		$ gh comment list 123 --since "critical-bug-report" --author "oncall*" --status resolved

		# Performance review analysis
		$ gh comment list 123 --author "performance-team" --since "load-test-date" --type review
		$ gh comment list 123 --status open --author "*perf*" --since "1 week ago"

		# Export for further analysis and automation
		$ gh comment list 123 --author "all-reviewers*" --since "quarter-start" --quiet | process-review-data.sh
		$ gh comment list 123 --ids-only --type review --status open | review-metrics.sh
	`),
	Args:   cobra.MaximumNArgs(1),
	PreRun: applyListConfigDefaults,
	RunE:   runList,
}

// applyListConfigDefaults applies configuration defaults to list command flags
func applyListConfigDefaults(cmd *cobra.Command, args []string) {
	config := GetConfig()

	// Apply filter defaults if flags weren't explicitly set
	if !cmd.Flags().Changed("author") && config.Defaults.Author != "" {
		author = config.Defaults.Author
	}
	if !cmd.Flags().Changed("status") && config.Filters.Status != "" {
		status = config.Filters.Status
	}
	if !cmd.Flags().Changed("type") && config.Filters.Type != "" {
		listType = config.Filters.Type
	}
	if !cmd.Flags().Changed("since") && config.Filters.Since != "" {
		since = config.Filters.Since
	}
	if !cmd.Flags().Changed("until") && config.Filters.Until != "" {
		until = config.Filters.Until
	}

	// Apply display defaults
	if !cmd.Flags().Changed("format") && config.Display.Format != "table" {
		// Map config format to list command format
		switch config.Display.Format {
		case "json":
			outputFormat = "json"
		case "quiet":
			quiet = true
		}
	}
	if !cmd.Flags().Changed("quiet") && config.Display.Quiet {
		quiet = config.Display.Quiet
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Legacy flags (kept for backward compatibility)
	listCmd.Flags().BoolVar(&showResolved, "resolved", false, "Include resolved comments (legacy, use --status instead) (default: false)")
	listCmd.Flags().BoolVar(&onlyUnresolved, "unresolved", false, "Show only unresolved comments (legacy, use --status instead) (default: false)")

	// Enhanced filtering flags
	listCmd.Flags().StringVar(&author, "author", "", "Filter comments by author (supports wildcards: 'user*', '*@company.com') (default: all authors)")
	listCmd.Flags().StringVar(&status, "status", "all", "Filter by comment status (open|resolved|all) (default: all)")
	listCmd.Flags().StringVar(&since, "since", "", "Show comments created after date (e.g., '2024-01-01', '1 week ago', '3 days ago') (default: all dates)")
	listCmd.Flags().StringVar(&until, "until", "", "Show comments created before date (e.g., '2024-12-31', '1 day ago') (default: all dates)")
	listCmd.Flags().StringVar(&resolved, "resolved-status", "", "Filter by resolution status (pending|resolved|dismissed) (default: all)")
	listCmd.Flags().StringVar(&listType, "type", "all", "Filter by comment type (issue|review|all) (default: all)")

	// Display options
	listCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output without URLs and IDs (default: false, shows full context)")
	listCmd.Flags().BoolVar(&hideAuthors, "hide-authors", false, "Hide author names for privacy (default: false)")

	// Output format options
	listCmd.Flags().StringVar(&outputFormat, "format", "default", "Output format (default|json) (default: default)")
	listCmd.Flags().BoolVar(&idsOnly, "ids-only", false, "Output only comment IDs (one per line) (default: false)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if listClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to initialize GitHub client: %w", err)
		}
		listClient = client
	}

	// Validate and parse filtering flags
	if err := validateAndParseFilters(); err != nil {
		return err
	}

	var pr int
	var repository string
	var err error

	// Parse PR argument using centralized function
	if len(args) == 1 {
		pr, err = parsePositiveInt(args[0], "PR number")
		if err != nil {
			return err
		}
		// Get repository for explicitly provided PR
		repository, err = getCurrentRepo()
		if err != nil {
			return err
		}
	} else {
		// Auto-detect PR and repository using centralized function
		repository, pr, err = getPRContext()
		if err != nil {
			return err
		}
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
	comments, err := fetchAllComments(listClient, repository, pr)
	if err != nil {
		return err
	}

	// Filter comments
	filteredComments := filterComments(comments)

	// Handle different output formats
	if idsOnly {
		displayIDsOnly(filteredComments)
	} else if outputFormat == "json" {
		if err := displayCommentsJSON(filteredComments, pr); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	} else {
		displayComments(filteredComments, pr)
	}

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
	CommitID  string `json:"commit_id,omitempty"`

	// Comment type
	Type string `json:"type"` // "issue" or "review"

	// Resolution status (for review comments)
	State string `json:"state,omitempty"` // "pending", "submitted", etc.
}

func fetchAllComments(client github.GitHubAPI, repo string, pr int) ([]Comment, error) {
	// Parse owner/repo
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format: %s (expected owner/repo)", repo)
	}
	owner, repoName := parts[0], parts[1]

	var allComments []Comment

	// Fetch general PR comments (issue comments)
	issueComments, err := client.ListIssueComments(owner, repoName, pr)
	if err != nil {
		return nil, formatActionableError("issue comments fetch", err)
	}

	// Convert issue comments
	for _, comment := range issueComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			HTMLURL:   "", // TODO: Add HTMLURL to github.Comment
			Type:      "issue",
		})
	}

	// Fetch review comments (line-specific)
	reviewComments, err := client.ListReviewComments(owner, repoName, pr)
	if err != nil {
		return nil, formatActionableError("review comments fetch", err)
	}

	// Convert review comments
	for _, comment := range reviewComments {
		allComments = append(allComments, Comment{
			ID:        comment.ID,
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			HTMLURL:   "", // TODO: Add HTMLURL to github.Comment
			Path:      comment.Path,
			Line:      comment.Line,
			CommitID:  comment.CommitID,
			Type:      "review",
		})
	}

	return allComments, nil
}

func validateAndParseFilters() error {
	// Validate status flag
	validStatuses := []string{"all", "open", "resolved"}
	if status != "" && !containsString(validStatuses, status) {
		return fmt.Errorf("invalid status '%s'. Must be one of: %s", status, strings.Join(validStatuses, ", "))
	}

	// Validate comment type flag
	validTypes := []string{"all", "issue", "review"}
	if listType != "" && !containsString(validTypes, listType) {
		return fmt.Errorf("invalid type '%s'. Must be one of: %s", listType, strings.Join(validTypes, ", "))
	}

	// Validate output format flag
	validFormats := []string{"default", "json"}
	if outputFormat != "" && !containsString(validFormats, outputFormat) {
		return fmt.Errorf("invalid format '%s'. Must be one of: %s", outputFormat, strings.Join(validFormats, ", "))
	}

	// Handle conflicting output options
	if idsOnly && outputFormat == "json" {
		return fmt.Errorf("cannot use --ids-only with --format json (use --format json to get structured data including IDs)")
	}

	// Parse since date
	if since != "" {
		parsedTime, err := parseFlexibleDate(since)
		if err != nil {
			return fmt.Errorf("invalid since date '%s': %w", since, err)
		}
		sinceTime = &parsedTime
	}

	// Parse until date
	if until != "" {
		parsedTime, err := parseFlexibleDate(until)
		if err != nil {
			return fmt.Errorf("invalid until date '%s': %w", until, err)
		}
		untilTime = &parsedTime
	}

	// Validate date range
	if sinceTime != nil && untilTime != nil && sinceTime.After(*untilTime) {
		return fmt.Errorf("since date (%s) cannot be after until date (%s)", since, until)
	}

	return nil
}

func parseFlexibleDate(dateStr string) (time.Time, error) {
	parsed, err := dateparser.Parse(nil, strings.TrimSpace(dateStr))
	if err != nil {
		return time.Time{}, err
	}
	return parsed.Time, nil
}

func filterComments(comments []Comment) []Comment {
	var filtered []Comment

	for _, comment := range comments {
		// Filter by author (supports wildcards)
		if author != "" && !matchesAuthorFilter(comment.Author, author) {
			continue
		}

		// Filter by comment type
		if listType != "all" && comment.Type != listType {
			continue
		}

		// Filter by status (legacy support)
		if showResolved && onlyUnresolved {
			// Conflicting flags - show all
		} else if onlyUnresolved {
			// Only show unresolved comments (this is a placeholder - actual resolution status would come from API)
			// For now, we'll consider all comments as "open" since we don't have resolution data
			if status == "resolved" {
				continue
			}
		} else if !showResolved {
			// Default behavior - don't show resolved comments
			if status == "resolved" {
				continue
			}
		}

		// Filter by new status flag
		if status != "all" {
			// This is a placeholder for actual resolution status filtering
			// In a real implementation, you'd check comment.ResolvedAt or similar
			// For now, we'll treat all comments as "open"
			if status == "resolved" {
				continue // Skip since we don't have resolution data yet
			}
		}

		// Filter by date range
		if sinceTime != nil && comment.CreatedAt.Before(*sinceTime) {
			continue
		}
		if untilTime != nil && comment.CreatedAt.After(*untilTime) {
			continue
		}

		filtered = append(filtered, comment)
	}

	return filtered
}

func matchesAuthorFilter(author, filter string) bool {
	// Exact match
	if author == filter {
		return true
	}

	// Wildcard matching
	if strings.Contains(filter, "*") {
		// Convert wildcard pattern to regex
		pattern := strings.ReplaceAll(regexp.QuoteMeta(filter), `\*`, `.*`)
		pattern = "^" + pattern + "$"

		if matched, err := regexp.MatchString(pattern, author); err == nil && matched {
			return true
		}
	}

	// Case-insensitive partial match
	return strings.Contains(strings.ToLower(author), strings.ToLower(filter))
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func displayComments(comments []Comment, pr int) {
	if len(comments) == 0 {
		fmt.Printf("No comments found on PR #%d\n", pr)
		return
	}

	// Collect unique commit IDs for summary
	commitIDs := make(map[string]bool)
	for _, comment := range comments {
		if comment.CommitID != "" {
			commitIDs[comment.CommitID] = true
		}
	}

	fmt.Printf("📝 Comments on PR #%d (%d total", pr, len(comments))
	if len(commitIDs) > 0 {
		fmt.Printf(", %d commit", len(commitIDs))
		if len(commitIDs) != 1 {
			fmt.Printf("s")
		}
	}
	fmt.Printf(")\n\n")

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

	// Safe color functions
	colorIssue := func(text string) string {
		if ColorIssueComment != nil {
			return ColorIssueComment.Sprint(text)
		}
		return text
	}
	colorReview := func(text string) string {
		if ColorReviewComment != nil {
			return ColorReviewComment.Sprint(text)
		}
		return text
	}
	colorHeader := func(text string) string {
		if ColorHeader != nil {
			return ColorHeader.Sprint(text)
		}
		return text
	}

	// Display general PR comments
	if len(issueComments) > 0 {
		fmt.Printf("%s (%d)\n", colorIssue("💬 General PR Comments"), len(issueComments))
		fmt.Println(strings.Repeat("─", SeparatorLength))
		for i, comment := range issueComments {
			displayComment(comment, i+1)
		}
		fmt.Println()
	}

	// Display review-level comments (parent comments that group line-specific ones)
	if len(reviewComments) > 0 {
		fmt.Printf("%s (%d)\n", colorReview("📋 Review Comments"), len(reviewComments))
		fmt.Println(strings.Repeat("─", SeparatorLength))
		for i, comment := range reviewComments {
			displayComment(comment, i+1)
		}
		fmt.Println()
	}

	// Display line-specific comments
	if len(lineComments) > 0 {
		fmt.Printf("%s (%d)\n", colorHeader("📍 Line-Specific Comments"), len(lineComments))
		fmt.Println(strings.Repeat("─", SeparatorLength))
		for i, comment := range lineComments {
			displayComment(comment, i+1)
		}
	}
}

func displayComment(comment Comment, index int) {
	// Header with author and timestamp
	timeAgo := formatTimeAgo(comment.CreatedAt)

	// Safe color functions that work even if colors aren't initialized
	colorID := func(text string) string {
		if ColorCommentID != nil {
			return ColorCommentID.Sprint(text)
		}
		return text
	}
	colorAuthor := func(text string) string {
		if ColorAuthor != nil {
			return ColorAuthor.Sprint(text)
		}
		return text
	}
	colorTime := func(text string) string {
		if ColorTimestamp != nil {
			return ColorTimestamp.Sprint(text)
		}
		return text
	}

	if hideAuthors {
		fmt.Printf("[%d] %s 👤 [hidden] • %s", index, colorID(fmt.Sprintf("ID:%d", comment.ID)), colorTime(timeAgo))
	} else {
		fmt.Printf("[%d] %s 👤 %s • %s", index, colorID(fmt.Sprintf("ID:%d", comment.ID)), colorAuthor(comment.Author), colorTime(timeAgo))
	}

	// Show review state for review-level comments
	if comment.Type == "review" && comment.State != "" {
		fmt.Printf(" %s", ColorizeReviewState(comment.State))
	}
	fmt.Println()

	// File and line info for line-specific comments
	if comment.Path != "" {
		lineInfo := fmt.Sprintf("L%d", comment.Line)
		if comment.StartLine > 0 && comment.StartLine != comment.Line {
			lineInfo = fmt.Sprintf("L%d-L%d", comment.StartLine, comment.Line)
		}
		colorFile := func(text string) string {
			if ColorFilePath != nil {
				return ColorFilePath.Sprint(text)
			}
			return text
		}
		colorLine := func(text string) string {
			if ColorLineNumber != nil {
				return ColorLineNumber.Sprint(text)
			}
			return text
		}
		colorSHA := func(text string) string {
			if ColorCommitSHA != nil {
				return ColorCommitSHA.Sprint(text)
			}
			return text
		}

		fmt.Printf("📁 %s:%s", colorFile(comment.Path), colorLine(lineInfo))

		// Show commit ID for review comments (helps with debugging and understanding)
		if comment.CommitID != "" {
			fmt.Printf(" • 📊 %s", colorSHA(comment.CommitID[:CommitSHADisplayLength])) // Show first chars of commit SHA
		}
		fmt.Println()

		// Show the actual diff context if available
		if comment.DiffHunk != "" {
			fmt.Printf("📝 Code Context:\n")
			displayDiffHunk(comment.DiffHunk)
		}
	}

	// Comment body (truncate if too long)
	body := strings.TrimSpace(comment.Body)
	if len(body) > MaxDisplayBodyLength {
		body = body[:MaxDisplayBodyLength-TruncationReserve] + TruncationSuffix
	}

	// Indent the comment body
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		fmt.Printf("   %s\n", line)
	}

	// Show URLs by default (AI-friendly), hide only in quiet mode
	if !quiet {
		colorURL := func(text string) string {
			if ColorURL != nil {
				return ColorURL.Sprint(text)
			}
			return text
		}
		fmt.Printf("   🔗 %s\n", colorURL(comment.HTMLURL))
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

// displayIDsOnly outputs only comment IDs, one per line
func displayIDsOnly(comments []Comment) {
	for _, comment := range comments {
		fmt.Printf("%d\n", comment.ID)
	}
}

// displayCommentsJSON outputs comments as JSON
func displayCommentsJSON(comments []Comment, pr int) error {
	output := struct {
		PR       int       `json:"pr"`
		Total    int       `json:"total"`
		Comments []Comment `json:"comments"`
	}{
		PR:       pr,
		Total:    len(comments),
		Comments: comments,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
