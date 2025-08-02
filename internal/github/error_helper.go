package github

import (
	"fmt"
	"strings"
)

// EnhancedAPIError provides intelligent error handling with suggestions
type EnhancedAPIError struct {
	OriginalError error
	Command       string
	CommentID     int
	Suggestions   []string
	AutoFix       string
}

func (e *EnhancedAPIError) Error() string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("Error: %v\n\n", e.OriginalError))

	msg.WriteString("🤖 **Intelligent Analysis**:\n")
	for _, suggestion := range e.Suggestions {
		msg.WriteString(fmt.Sprintf("   • %s\n", suggestion))
	}

	if e.AutoFix != "" {
		msg.WriteString(fmt.Sprintf("\n💡 **Auto-correction suggestion**:\n   %s\n", e.AutoFix))
	}

	return msg.String()
}

// AnalyzeAndEnhanceError provides intelligent error analysis and suggestions
func AnalyzeAndEnhanceError(err error, command string, commentID int) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	enhanced := &EnhancedAPIError{
		OriginalError: err,
		Command:       command,
		CommentID:     commentID,
		Suggestions:   []string{},
	}

	// Analyze different error patterns with context-aware suggestions
	if strings.Contains(errMsg, "404") && strings.Contains(errMsg, "pulls/comments") {
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Comment ID might be for an issue comment (general PR comment), not a review comment")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Issue comments appear in the main conversation tab, not 'Files changed'")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Run 'gh comment list <PR>' to see all comment types with their IDs")
		enhanced.AutoFix = fmt.Sprintf("gh comment %s %d --type issue", command, commentID)
	}

	if strings.Contains(errMsg, "404") && strings.Contains(errMsg, "issues/comments") {
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Comment ID might be for a review comment (line-specific), not an issue comment")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Review comments show file paths like '📁 src/main.go:L42' in list output")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Issue comments don't show file paths and appear in main conversation")
		enhanced.AutoFix = fmt.Sprintf("gh comment %s %d --type review", command, commentID)
	}

	if strings.Contains(errMsg, "in_reply_to_id") && strings.Contains(errMsg, "not permitted") {
		enhanced.Suggestions = append(enhanced.Suggestions,
			"GitHub's review comment threading API has changed - direct replies not supported")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Use reactions for quick feedback: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"For longer responses, create a new comment: 'gh comment add <PR> <file> <line> \"response\"'")
		enhanced.AutoFix = fmt.Sprintf("gh comment reply %d --reaction +1  # Quick approval", commentID)
	}

	if strings.Contains(errMsg, "422") && strings.Contains(errMsg, "commitId") {
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Review creation needs valid commit IDs for each comment")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"The PR might have new commits since the review was started")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Use individual 'gh comment add' commands instead of bulk review creation")
		enhanced.AutoFix = "gh comment add <PR> <file> <line> \"comment\"  # For current commit"
	}

	if strings.Contains(errMsg, "resource not found") {
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Comment ID doesn't exist or you lack repository access")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Use 'gh comment list <PR>' to see all available comments with IDs")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"Comment IDs are shown as [1], [2], etc. in the list output")
		enhanced.AutoFix = "gh comment list <PR>  # See all comments and their IDs"
	}

	// Command-specific guidance based on your built-in help
	switch command {
	case "reply":
		enhanced.Suggestions = append(enhanced.Suggestions,
			"💡 Reply alternatives from 'gh comment reply --help':")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"   • Use reactions: --reaction +1, heart, hooray")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"   • Resolve conversations: --resolve")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"   • Add suggestions: \"[SUGGEST: improved code]\"")

	case "edit":
		enhanced.Suggestions = append(enhanced.Suggestions,
			"💡 Edit command supports both issue and review comments")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"   • Use -m flag for multi-line: -m \"line1\" -m \"line2\"")

	case "review":
		enhanced.Suggestions = append(enhanced.Suggestions,
			"💡 Review alternatives from 'gh comment review --help':")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"   • Try individual comments: 'gh comment add <PR> <file> <line> \"comment\"'")
		enhanced.Suggestions = append(enhanced.Suggestions,
			"   • Use add-review for draft reviews: 'gh comment add-review'")
	}

	// Add contextual help reference
	enhanced.Suggestions = append(enhanced.Suggestions,
		fmt.Sprintf("📖 See 'gh comment %s --help' for detailed examples and usage", command))

	return enhanced
}

// GetCommentTypeHelp provides detailed help about comment types
func GetCommentTypeHelp() string {
	return `
🔍 **Understanding GitHub Comment Types**:

**Issue Comments** (General PR Comments):
   • Appear in the main conversation tab
   • Are not tied to specific lines of code
   • Use: gh comment add <PR> "General feedback about the PR"
   • IDs typically work with --type issue

**Review Comments** (Line-specific Comments):
   • Appear in the "Files changed" tab
   • Are tied to specific lines/ranges in files
   • Use: gh comment add <PR> <file> <line> "Specific code feedback"
   • IDs typically work with --type review

**Pro Tips**:
   • Use 'gh comment list <PR>' to see all comments with their types
   • Review comments show file paths (e.g., "📁 src/main.go:L42")
   • Issue comments don't show file paths
   • When in doubt, try both --type issue and --type review
`
}
