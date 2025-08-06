package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"

	"github.com/spf13/cobra"

	"github.com/silouanwright/gh-comment/internal/github"
)

var (
	exportFormat    string
	exportOutput    string
	exportInclude   []string
	includeResolved bool

	// Client for dependency injection (tests can override)
	exportClient github.GitHubAPI
)

var exportCmd = &cobra.Command{
	Use:   "export <pr>",
	Short: "Export PR comments to various formats",
	Long: heredoc.Doc(`
		Export PR comments to JSON, CSV, Markdown, or HTML format.

		This command fetches all comments from a PR (both issue and review comments)
		and exports them in the specified format. You can filter which fields to
		include and save to a file or output to stdout.

		Supported formats:
		- json: Machine-readable JSON format
		- csv: Spreadsheet-compatible CSV format
		- markdown: Documentation-friendly Markdown format
		- html: Presentation-ready HTML format
	`),
	Example: heredoc.Doc(`
		# Export to JSON (default)
		$ gh comment export 123

		# Export to CSV file
		$ gh comment export 123 --format csv --output comments.csv

		# Export only specific fields to JSON
		$ gh comment export 123 --include id,author,body

		# Export to Markdown including resolved comments
		$ gh comment export 123 --format markdown --include-resolved

		# Export to HTML for presentation
		$ gh comment export 123 --format html --output pr-123-review.html

		# Export with auto-detected PR
		$ gh comment export --format csv
	`),
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Export format (json|csv|markdown|html) (default: json)")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file (default: stdout)")
	exportCmd.Flags().StringSliceVar(&exportInclude, "include", []string{}, "Fields to include (default: all)")
	exportCmd.Flags().BoolVar(&includeResolved, "include-resolved", false, "Include resolved comments (default: false)")
}

type ExportComment struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	URL       string    `json:"url"`
	DiffHunk  string    `json:"diff_hunk,omitempty"`
	CommitID  string    `json:"commit_id,omitempty"`
	InReplyTo int       `json:"in_reply_to,omitempty"`
	Resolved  bool      `json:"resolved,omitempty"`
}

func runExport(cmd *cobra.Command, args []string) error {
	// Initialize client if not set (production use)
	if exportClient == nil {
		client, err := createGitHubClient()
		if err != nil {
			return fmt.Errorf("failed to create GitHub client: %w", err)
		}
		exportClient = client
	}

	var pr int
	var err error

	// Parse PR number
	if len(args) == 1 {
		pr, err = parsePositiveInt(args[0], "PR number")
		if err != nil {
			return err
		}
	} else {
		// Auto-detect PR
		_, pr, err = getPRContext()
		if err != nil {
			return err
		}
	}

	// Get repository
	repository, err := getCurrentRepo()
	if err != nil {
		return err
	}

	// Validate repository name
	if err := validateRepositoryName(repository); err != nil {
		return err
	}

	// Parse owner/repo
	parts := strings.Split(repository, "/")
	owner, repoName := parts[0], parts[1]

	// Validate format
	validFormats := []string{"json", "csv", "markdown", "html"}
	isValidFormat := false
	for _, valid := range validFormats {
		if exportFormat == valid {
			isValidFormat = true
			break
		}
	}
	if !isValidFormat {
		return fmt.Errorf("invalid format: %s (must be json, csv, markdown, or html)", exportFormat)
	}

	if verbose {
		fmt.Printf("Repository: %s\n", repository)
		fmt.Printf("PR: %d\n", pr)
		fmt.Printf("Format: %s\n", exportFormat)
		if exportOutput != "" {
			fmt.Printf("Output: %s\n", exportOutput)
		}
		fmt.Println()
	}

	// Fetch all comments
	comments, err := fetchAllCommentsForExport(exportClient, owner, repoName, pr)
	if err != nil {
		return fmt.Errorf("failed to fetch comments: %w", err)
	}

	// Filter resolved comments if needed
	if !includeResolved {
		var filtered []ExportComment
		for _, comment := range comments {
			if !comment.Resolved {
				filtered = append(filtered, comment)
			}
		}
		comments = filtered
	}

	// Setup output writer
	var writer io.Writer
	if exportOutput != "" {
		file, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}

	// Export based on format
	switch exportFormat {
	case "json":
		return exportJSON(writer, comments)
	case "csv":
		return exportCSV(writer, comments)
	case "markdown":
		return exportMarkdown(writer, comments, repository, pr)
	case "html":
		return exportHTML(writer, comments, repository, pr)
	}

	return nil
}

func fetchAllCommentsForExport(client github.GitHubAPI, owner, repo string, pr int) ([]ExportComment, error) {
	var allComments []ExportComment

	// Fetch issue comments
	issueComments, err := client.ListIssueComments(owner, repo, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue comments: %w", err)
	}

	for _, comment := range issueComments {
		allComments = append(allComments, ExportComment{
			ID:        comment.ID,
			Type:      "issue",
			Author:    comment.User.Login,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			URL:       fmt.Sprintf("https://github.com/%s/%s/issues/%d#issuecomment-%d", owner, repo, pr, comment.ID),
		})
	}

	// Fetch review comments
	reviewComments, err := client.ListReviewComments(owner, repo, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review comments: %w", err)
	}

	for _, comment := range reviewComments {
		allComments = append(allComments, ExportComment{
			ID:        comment.ID,
			Type:      "review",
			Author:    comment.User.Login,
			Body:      comment.Body,
			File:      comment.Path,
			Line:      comment.Line,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			URL:       fmt.Sprintf("https://github.com/%s/%s/pull/%d#discussion_r%d", owner, repo, pr, comment.ID),
			CommitID:  comment.CommitID,
			// Note: DiffHunk, InReplyTo, and Resolved are not available in the current Comment struct
			// These would need to be fetched separately if needed
		})
	}

	return allComments, nil
}

func exportJSON(w io.Writer, comments []ExportComment) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	// If specific fields requested, filter them
	if len(exportInclude) > 0 {
		var filtered []map[string]interface{}
		for _, comment := range comments {
			item := make(map[string]interface{})
			for _, field := range exportInclude {
				switch field {
				case "id":
					item["id"] = comment.ID
				case "type":
					item["type"] = comment.Type
				case "author":
					item["author"] = comment.Author
				case "body":
					item["body"] = comment.Body
				case "file":
					if comment.File != "" {
						item["file"] = comment.File
					}
				case "line":
					if comment.Line != 0 {
						item["line"] = comment.Line
					}
				case "created_at":
					item["created_at"] = comment.CreatedAt
				case "updated_at":
					item["updated_at"] = comment.UpdatedAt
				case "url":
					item["url"] = comment.URL
				case "diff_hunk":
					if comment.DiffHunk != "" {
						item["diff_hunk"] = comment.DiffHunk
					}
				case "commit_id":
					if comment.CommitID != "" {
						item["commit_id"] = comment.CommitID
					}
				case "in_reply_to":
					if comment.InReplyTo != 0 {
						item["in_reply_to"] = comment.InReplyTo
					}
				case "resolved":
					item["resolved"] = comment.Resolved
				}
			}
			filtered = append(filtered, item)
		}
		return encoder.Encode(filtered)
	}

	return encoder.Encode(comments)
}

func exportCSV(w io.Writer, comments []ExportComment) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Define headers
	headers := []string{"ID", "Type", "Author", "Body", "File", "Line", "Created At", "Updated At", "URL", "Resolved"}
	if len(exportInclude) > 0 {
		headers = exportInclude
	}

	// Write headers
	if err := csvWriter.Write(headers); err != nil {
		return err
	}

	// Write data
	for _, comment := range comments {
		var row []string
		for _, header := range headers {
			switch strings.ToLower(header) {
			case "id":
				row = append(row, strconv.Itoa(comment.ID))
			case "type":
				row = append(row, comment.Type)
			case "author":
				row = append(row, comment.Author)
			case "body":
				// Escape newlines in CSV
				body := strings.ReplaceAll(comment.Body, "\n", "\\n")
				row = append(row, body)
			case "file":
				row = append(row, comment.File)
			case "line":
				if comment.Line != 0 {
					row = append(row, strconv.Itoa(comment.Line))
				} else {
					row = append(row, "")
				}
			case "created at", "created_at":
				row = append(row, comment.CreatedAt.Format(time.RFC3339))
			case "updated at", "updated_at":
				if !comment.UpdatedAt.IsZero() {
					row = append(row, comment.UpdatedAt.Format(time.RFC3339))
				} else {
					row = append(row, "")
				}
			case "url":
				row = append(row, comment.URL)
			case "resolved":
				row = append(row, strconv.FormatBool(comment.Resolved))
			default:
				row = append(row, "")
			}
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func exportMarkdown(w io.Writer, comments []ExportComment, repo string, pr int) error {
	fmt.Fprintf(w, "# PR Comments Export\n\n")
	fmt.Fprintf(w, "**Repository:** %s  \n", repo)
	fmt.Fprintf(w, "**PR:** #%d  \n", pr)
	fmt.Fprintf(w, "**Total Comments:** %d  \n", len(comments))
	fmt.Fprintf(w, "**Exported:** %s  \n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Group by type
	var issueComments, reviewComments []ExportComment
	for _, comment := range comments {
		if comment.Type == "issue" {
			issueComments = append(issueComments, comment)
		} else {
			reviewComments = append(reviewComments, comment)
		}
	}

	// Export issue comments
	if len(issueComments) > 0 {
		fmt.Fprintf(w, "## General PR Comments (%d)\n\n", len(issueComments))
		for _, comment := range issueComments {
			fmt.Fprintf(w, "### Comment #%d\n", comment.ID)
			fmt.Fprintf(w, "**Author:** @%s  \n", comment.Author)
			fmt.Fprintf(w, "**Created:** %s  \n", comment.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Fprintf(w, "\n%s\n\n", comment.Body)
			_, _ = fmt.Fprintln(w, "---") // Export output
		}
	}

	// Export review comments
	if len(reviewComments) > 0 {
		fmt.Fprintf(w, "## Review Comments (%d)\n\n", len(reviewComments))
		for _, comment := range reviewComments {
			fmt.Fprintf(w, "### Comment #%d\n", comment.ID)
			fmt.Fprintf(w, "**Author:** @%s  \n", comment.Author)
			fmt.Fprintf(w, "**File:** `%s:%d`  \n", comment.File, comment.Line)
			fmt.Fprintf(w, "**Created:** %s  \n", comment.CreatedAt.Format("2006-01-02 15:04"))
			if comment.Resolved {
				_, _ = fmt.Fprintln(w, "**Status:** ✅ Resolved") // Export output
			}
			fmt.Fprintf(w, "\n%s\n\n", comment.Body)
			if comment.DiffHunk != "" {
				_, _ = fmt.Fprintln(w, "```diff")        // Export output
				_, _ = fmt.Fprintln(w, comment.DiffHunk) // Export output
				_, _ = fmt.Fprintln(w, "```")            // Export output
			}
			_, _ = fmt.Fprintln(w, "---") // Export output
		}
	}

	return nil
}

func exportHTML(w io.Writer, comments []ExportComment, repo string, pr int) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>PR #%d Comments - %s</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif; margin: 20px; line-height: 1.6; }
        h1, h2, h3 { color: #24292e; }
        .comment { background: #f6f8fa; border: 1px solid #d1d5da; border-radius: 6px; padding: 16px; margin-bottom: 16px; }
        .comment-header { margin-bottom: 8px; }
        .author { font-weight: bold; color: #0366d6; }
        .timestamp { color: #586069; font-size: 14px; }
        .file-info { background: #f1f8ff; padding: 4px 8px; border-radius: 3px; font-family: monospace; font-size: 12px; }
        .body { margin-top: 8px; }
        .diff { background: #fafbfc; border: 1px solid #e1e4e8; border-radius: 3px; padding: 8px; margin-top: 8px; overflow-x: auto; }
        pre { margin: 0; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace; font-size: 12px; }
        .resolved { color: #28a745; font-size: 14px; }
        .stats { background: #f6f8fa; padding: 12px; border-radius: 6px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <h1>PR #%d Comments Export</h1>
    <div class="stats">
        <strong>Repository:</strong> %s<br>
        <strong>Total Comments:</strong> %d<br>
        <strong>Exported:</strong> %s
    </div>
`, pr, repo, pr, repo, len(comments), time.Now().Format("2006-01-02 15:04:05"))

	// Group by type
	var issueComments, reviewComments []ExportComment
	for _, comment := range comments {
		if comment.Type == "issue" {
			issueComments = append(issueComments, comment)
		} else {
			reviewComments = append(reviewComments, comment)
		}
	}

	// Export issue comments
	if len(issueComments) > 0 {
		html += fmt.Sprintf("\n    <h2>General PR Comments (%d)</h2>\n", len(issueComments))
		for _, comment := range issueComments {
			html += fmt.Sprintf(`    <div class="comment">
        <div class="comment-header">
            <span class="author">@%s</span>
            <span class="timestamp">%s</span>
        </div>
        <div class="body">%s</div>
    </div>
`, comment.Author, comment.CreatedAt.Format("Jan 2, 2006 15:04"), strings.ReplaceAll(comment.Body, "\n", "<br>"))
		}
	}

	// Export review comments
	if len(reviewComments) > 0 {
		html += fmt.Sprintf("\n    <h2>Review Comments (%d)</h2>\n", len(reviewComments))
		for _, comment := range reviewComments {
			resolvedTag := ""
			if comment.Resolved {
				resolvedTag = ` <span class="resolved">✅ Resolved</span>`
			}

			diffSection := ""
			if comment.DiffHunk != "" {
				diffSection = fmt.Sprintf(`        <div class="diff"><pre>%s</pre></div>`, comment.DiffHunk)
			}

			html += fmt.Sprintf(`    <div class="comment">
        <div class="comment-header">
            <span class="author">@%s</span>
            <span class="file-info">%s:%d</span>
            <span class="timestamp">%s</span>%s
        </div>
        <div class="body">%s</div>
%s
    </div>
`, comment.Author, comment.File, comment.Line, comment.CreatedAt.Format("Jan 2, 2006 15:04"),
				resolvedTag, strings.ReplaceAll(comment.Body, "\n", "<br>"), diffSection)
		}
	}

	html += "\n</body>\n</html>"

	_, err := fmt.Fprint(w, html)
	return err
}
