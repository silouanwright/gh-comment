// Package help provides utilities for building consistent, professional help text
// following GitHub CLI patterns.
package help

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
)

// Example represents a single command example with description and command.
type Example struct {
	Description string
	Command     string
}

// BuildLongHelp creates a Long help string from a heredoc string.
// This matches the GitHub CLI pattern of using heredoc.Doc() for multi-line help.
func BuildLongHelp(content string) string {
	return heredoc.Doc(content)
}

// BuildExamples creates an Example section string from a slice of examples.
// This follows the GitHub CLI pattern for formatting examples consistently.
func BuildExamples(examples []Example) string {
	if len(examples) == 0 {
		return ""
	}

	var lines []string
	for _, ex := range examples {
		// Format: "# Description"
		//         "$ command"
		lines = append(lines, fmt.Sprintf("# %s", ex.Description))
		lines = append(lines, fmt.Sprintf("$ %s", ex.Command))
		lines = append(lines, "") // Empty line between examples
	}

	// Remove the last empty line
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

// BuildSectionHelp creates a help section with title and content.
func BuildSectionHelp(title, content string) string {
	return fmt.Sprintf("\n%s\n%s\n%s",
		title,
		strings.Repeat("-", len(title)),
		content)
}

// FormatCommandUsage formats a command usage string with proper syntax highlighting.
func FormatCommandUsage(usage string) string {
	// Add backticks around command parts for better visibility
	return strings.ReplaceAll(usage, "<", "`<") + "`"
}

// BuildWorkflowHelp creates workflow-oriented help text.
func BuildWorkflowHelp(workflow string) string {
	return BuildLongHelp(fmt.Sprintf(`
		%s

		This command supports both interactive and automated workflows,
		making it ideal for both manual code review and CI/CD integration.
	`, workflow))
}
