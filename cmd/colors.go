package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

var (
	// Color settings
	colorEnabled = true

	// Test override for terminal detection (used in testing)
	testTerminalOverride *bool

	// Color definitions for different types of output
	ColorCommentID    *color.Color
	ColorAuthor       *color.Color
	ColorTimestamp    *color.Color
	ColorFilePath     *color.Color
	ColorLineNumber   *color.Color
	ColorURL          *color.Color
	ColorSuccess      *color.Color
	ColorError        *color.Color
	ColorWarning      *color.Color
	ColorHeader       *color.Color
	ColorReviewState  *color.Color
	ColorCommitSHA    *color.Color
	ColorIssueComment *color.Color
	ColorReviewComment *color.Color
)

// InitColors initializes color settings and creates color functions
func InitColors() {
	// Detect if colors should be enabled
	colorEnabled = ShouldUseColor()

	if colorEnabled {
		// Enable colors globally
		color.NoColor = false

		// Comment-related colors
		ColorCommentID = color.New(color.FgCyan, color.Bold)
		ColorAuthor = color.New(color.FgGreen, color.Bold)
		ColorTimestamp = color.New(color.FgBlue)
		ColorFilePath = color.New(color.FgYellow)
		ColorLineNumber = color.New(color.FgYellow, color.Bold)
		ColorURL = color.New(color.FgBlue, color.Underline)

		// Status colors
		ColorSuccess = color.New(color.FgGreen, color.Bold)
		ColorError = color.New(color.FgRed, color.Bold)
		ColorWarning = color.New(color.FgYellow, color.Bold)

		// Section headers
		ColorHeader = color.New(color.FgMagenta, color.Bold)
		ColorReviewState = color.New(color.FgCyan)
		ColorCommitSHA = color.New(color.FgBlue)

		// Comment type colors
		ColorIssueComment = color.New(color.FgGreen)
		ColorReviewComment = color.New(color.FgYellow)
	} else {
		// Create disabled color functions that just return the text as-is
		ColorCommentID = color.New()
		ColorAuthor = color.New()
		ColorTimestamp = color.New()
		ColorFilePath = color.New()
		ColorLineNumber = color.New()
		ColorURL = color.New()
		ColorSuccess = color.New()
		ColorError = color.New()
		ColorWarning = color.New()
		ColorHeader = color.New()
		ColorReviewState = color.New()
		ColorCommitSHA = color.New()
		ColorIssueComment = color.New()
		ColorReviewComment = color.New()

		// Disable all colors
		color.NoColor = true
	}
}

// ShouldUseColor determines if colors should be used based on environment and flags
func ShouldUseColor() bool {
	// Check if --no-color flag was set
	if noColor {
		return false
	}

	// Check if output is being piped or redirected
	if !isTerminal() {
		return false
	}

	// Check NO_COLOR environment variable (http://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}

	return true
}

// isTerminal checks if output is going to a terminal
func isTerminal() bool {
	// Use test override if set (for testing purposes)
	if testTerminalOverride != nil {
		return *testTerminalOverride
	}

	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// DisableColors forcibly disables all colors (useful for testing)
func DisableColors() {
	colorEnabled = false
	color.NoColor = true
	InitColors()
}

// EnableColors forcibly enables colors (useful for testing)
func EnableColors() {
	colorEnabled = true
	color.NoColor = false
	noColor = false // Reset the flag too

	// Force initialize colors without checking ShouldUseColor()
	// Comment-related colors
	ColorCommentID = color.New(color.FgCyan, color.Bold)
	ColorAuthor = color.New(color.FgGreen, color.Bold)
	ColorTimestamp = color.New(color.FgBlue)
	ColorFilePath = color.New(color.FgYellow)
	ColorLineNumber = color.New(color.FgYellow, color.Bold)
	ColorURL = color.New(color.FgBlue, color.Underline)

	// Status colors
	ColorSuccess = color.New(color.FgGreen, color.Bold)
	ColorError = color.New(color.FgRed, color.Bold)
	ColorWarning = color.New(color.FgYellow, color.Bold)

	// Section headers
	ColorHeader = color.New(color.FgMagenta, color.Bold)
	ColorReviewState = color.New(color.FgCyan)
	ColorCommitSHA = color.New(color.FgBlue)

	// Comment type colors
	ColorIssueComment = color.New(color.FgGreen)
	ColorReviewComment = color.New(color.FgYellow)
}

// Helper functions for common color patterns

// ColorizeCommentType returns colored text based on comment type
func ColorizeCommentType(commentType string) string {
	switch commentType {
	case "issue":
		if ColorIssueComment != nil {
			return ColorIssueComment.Sprintf("💬 General PR Comments")
		}
		return "💬 General PR Comments"
	case "review":
		if ColorReviewComment != nil {
			return ColorReviewComment.Sprintf("📋 Review Comments")
		}
		return "📋 Review Comments"
	default:
		return commentType
	}
}

// ColorizeReviewState returns colored text based on review state
func ColorizeReviewState(state string) string {
	switch state {
	case "approved":
		if ColorSuccess != nil {
			return ColorSuccess.Sprintf("✅ %s", state)
		}
		return fmt.Sprintf("✅ %s", state)
	case "changes_requested":
		if ColorError != nil {
			return ColorError.Sprintf("🔴 %s", state)
		}
		return fmt.Sprintf("🔴 %s", state)
	case "commented":
		if ColorReviewState != nil {
			return ColorReviewState.Sprintf("💬 %s", state)
		}
		return fmt.Sprintf("💬 %s", state)
	default:
		if ColorReviewState != nil {
			return ColorReviewState.Sprintf("📝 %s", state)
		}
		return fmt.Sprintf("📝 %s", state)
	}
}

// ColorizeSuccess returns green colored success text
func ColorizeSuccess(text string) string {
	if ColorSuccess == nil {
		return fmt.Sprintf("✅ %s", text)
	}
	return ColorSuccess.Sprintf("✅ %s", text)
}

// ColorizeError returns red colored error text
func ColorizeError(text string) string {
	if ColorError == nil {
		return fmt.Sprintf("❌ %s", text)
	}
	return ColorError.Sprintf("❌ %s", text)
}

// ColorizeWarning returns yellow colored warning text
func ColorizeWarning(text string) string {
	if ColorWarning == nil {
		return fmt.Sprintf("⚠️ %s", text)
	}
	return ColorWarning.Sprintf("⚠️ %s", text)
}
