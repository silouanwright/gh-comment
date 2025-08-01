package cmd

import (
	"regexp"
	"strings"
)

// expandSuggestions processes inline suggestion syntax and converts to GitHub markdown
func expandSuggestions(message string) string {
	// First handle multi-line suggestions: <<<SUGGEST\ncode\nSUGGEST>>>
	message = expandMultilineSuggestions(message)

	// Then handle inline suggestions: [SUGGEST: code]
	message = expandInlineSuggestions(message)

	return message
}

// expandInlineSuggestions handles [SUGGEST: code] syntax
func expandInlineSuggestions(message string) string {
	// Regex to match [SUGGEST: code] patterns
	re := regexp.MustCompile(`\[SUGGEST:\s*([^\]]+)\]`)

	return re.ReplaceAllStringFunc(message, func(match string) string {
		// Extract the code part
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match // Return original if parsing fails
		}

		code := strings.TrimSpace(submatches[1])
		return "\n\n```suggestion\n" + code + "\n```\n\n"
	})
}

// expandMultilineSuggestions handles <<<SUGGEST\ncode\nSUGGEST>>> syntax
func expandMultilineSuggestions(message string) string {
	// Regex to match <<<SUGGEST...code...SUGGEST>>> blocks (with flexible whitespace)
	re := regexp.MustCompile(`(?s)<<<SUGGEST\s*\n(.*?)\nSUGGEST>>>`)

	return re.ReplaceAllStringFunc(message, func(match string) string {
		// Extract the code part
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match // Return original if parsing fails
		}

		code := strings.TrimSpace(submatches[1])
		return "\n\n```suggestion\n" + code + "\n```\n\n"
	})
}
