package github

import (
	"fmt"
	"strings"
)

// CommentInfo contains detected information about a comment
type CommentInfo struct {
	ID       int
	Type     string // "issue" or "review"
	FilePath string // empty for issue comments
	Line     int    // 0 for issue comments
	Found    bool
}

// DetectCommentType intelligently detects comment type by looking it up
func (c *RealClient) DetectCommentType(owner, repo string, commentID int, prNumber int) (*CommentInfo, error) {
	if prNumber <= 0 {
		return nil, fmt.Errorf("PR number required for comment type detection")
	}

	// Get all comments to find the specific one
	issueComments, err := c.ListIssueComments(owner, repo, prNumber)
	if err == nil {
		for _, comment := range issueComments {
			if comment.ID == commentID {
				return &CommentInfo{
					ID:    commentID,
					Type:  "issue",
					Found: true,
				}, nil
			}
		}
	}

	// Check review comments
	reviewComments, err := c.ListReviewComments(owner, repo, prNumber)
	if err == nil {
		for _, comment := range reviewComments {
			if comment.ID == commentID {
				return &CommentInfo{
					ID:       commentID,
					Type:     "review",
					FilePath: comment.Path,
					Line:     comment.Line,
					Found:    true,
				}, nil
			}
		}
	}

	return &CommentInfo{
		ID:    commentID,
		Found: false,
	}, nil
}

// GetHelpfulErrorMessage provides context-aware error messaging
func GetHelpfulErrorMessage(operation, command string, commentID int, commentInfo *CommentInfo) string {
	var msg strings.Builder

	if !commentInfo.Found {
		msg.WriteString(fmt.Sprintf("❌ **Comment #%d not found**\n\n", commentID))
		msg.WriteString("🔍 **Troubleshooting**:\n")
		msg.WriteString(fmt.Sprintf("   • Run `gh comment list <PR>` to see all available comments\n"))
		msg.WriteString(fmt.Sprintf("   • Check that comment #%d exists and you have access\n", commentID))
		msg.WriteString(fmt.Sprintf("   • Comment might have been deleted or you're using wrong PR number\n\n"))
		msg.WriteString(fmt.Sprintf("💡 **Next step**: `gh comment list <PR>`\n"))
		return msg.String()
	}

	// Provide specific guidance based on comment type
	if commentInfo.Type == "issue" {
		msg.WriteString(fmt.Sprintf("✅ **Issue Comment #%d** (General PR comment)\n\n", commentID))
		msg.WriteString("📍 **Location**: Main conversation tab\n")
		msg.WriteString("🎯 **Best for**: General feedback, overall PR discussions\n\n")

		switch operation {
		case "reply":
			msg.WriteString("💬 **Reply Options**:\n")
			msg.WriteString(fmt.Sprintf("   • Reactions: `gh comment reply %d --reaction +1`\n", commentID))
			msg.WriteString(fmt.Sprintf("   • Note: Direct text replies to issue comments create new comments\n"))
			msg.WriteString(fmt.Sprintf("   • Alternative: `gh comment add <PR> \"Reply message\"`\n"))
		case "edit":
			msg.WriteString("✏️ **Edit Command**:\n")
			msg.WriteString(fmt.Sprintf("   • `gh comment edit %d \"Updated message\"`\n", commentID))
		}
	} else {
		msg.WriteString(fmt.Sprintf("✅ **Review Comment #%d** (Line-specific comment)\n\n", commentID))
		if commentInfo.FilePath != "" {
			msg.WriteString(fmt.Sprintf("📁 **File**: %s:%d\n", commentInfo.FilePath, commentInfo.Line))
		}
		msg.WriteString("📍 **Location**: 'Files changed' tab\n")
		msg.WriteString("🎯 **Best for**: Specific code feedback, suggestions\n\n")

		switch operation {
		case "reply":
			msg.WriteString("💬 **Reply Options**:\n")
			msg.WriteString(fmt.Sprintf("   • Reactions: `gh comment reply %d --reaction +1`\n", commentID))
			msg.WriteString(fmt.Sprintf("   • New comment on same line: `gh comment add <PR> %s %d \"Reply\"`\n",
				commentInfo.FilePath, commentInfo.Line))
			msg.WriteString("   • Note: GitHub doesn't support threaded replies to review comments\n")
		case "edit":
			msg.WriteString("✏️ **Edit Command**:\n")
			msg.WriteString(fmt.Sprintf("   • `gh comment edit %d \"Updated message\"`\n", commentID))
		}
	}

	msg.WriteString(fmt.Sprintf("\n📖 **Full help**: `gh comment %s --help`\n", command))
	return msg.String()
}

// CreateSmartError creates an intelligent error with detection
func CreateSmartError(client GitHubAPI, operation, command string, commentID int, prNumber int, originalErr error) error {
	// For now, we'll create a simpler version without full detection
	// since we don't have owner/repo context here
	commentInfo := &CommentInfo{
		ID:    commentID,
		Found: false, // Assume failure for smart error guidance
	}

	// Create enhanced error message
	helpMsg := GetHelpfulErrorMessage(operation, command, commentID, commentInfo)

	return fmt.Errorf("Operation failed: %v\n\n%s", originalErr, helpMsg)
}
