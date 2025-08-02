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
	// Use a more sophisticated approach to handle nested brackets
	result := message
	for {
		start := strings.Index(result, "[SUGGEST:")
		if start == -1 {
			break // No more suggestions found
		}

		// Find the matching closing bracket by counting bracket depth
		bracketCount := 0
		suggestStart := start + len("[SUGGEST:")
		end := -1

		for i := suggestStart; i < len(result); i++ {
			char := result[i]
			if char == '[' {
				bracketCount++
			} else if char == ']' {
				if bracketCount == 0 {
					end = i
					break
				}
				bracketCount--
			}
		}

		if end == -1 {
			break // No matching closing bracket found
		}

		// Extract the code part and trim whitespace
		code := strings.TrimSpace(result[suggestStart:end])

		// Replace the [SUGGEST: code] with GitHub suggestion syntax
		replacement := "\n\n```suggestion\n" + code + "\n```\n\n"

		// Replace this occurrence and continue searching
		result = result[:start] + replacement + result[end+1:]
	}

	return result
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
