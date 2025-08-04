package cmd

import (
	"regexp"
	"strconv"
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

// expandInlineSuggestions handles [SUGGEST: code] and [SUGGEST:<offset>: code] syntax
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

		// Extract the content and parse for offset syntax
		content := strings.TrimSpace(result[suggestStart:end])
		offset, code := parseOffsetSuggestion(content)

		// Replace the [SUGGEST: code] or [SUGGEST:<offset>: code] with GitHub suggestion syntax
		var replacement string
		if offset != 0 {
			// Format with offset for GitHub
			replacement = "\n\n```suggestion:" + formatOffset(offset) + "\n" + code + "\n```\n\n"
		} else {
			// Standard suggestion format
			replacement = "\n\n```suggestion\n" + code + "\n```\n\n"
		}

		// Replace this occurrence and continue searching
		result = result[:start] + replacement + result[end+1:]
	}

	return result
}

// parseOffsetSuggestion parses content for offset syntax: "<offset>: code" or just "code"
func parseOffsetSuggestion(content string) (int, string) {
	// Check for offset pattern: +N: or -N: or N:
	offsetPattern := regexp.MustCompile(`^([+-]?\d+):\s*(.*)$`)
	matches := offsetPattern.FindStringSubmatch(content)

	if len(matches) == 3 {
		// Parse the offset
		offsetStr := matches[1]
		code := strings.TrimSpace(matches[2])

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < -999 || offset > 999 {
			// Invalid offset, treat as regular suggestion
			return 0, content
		}

		return offset, code
	}

	// No offset found, return as regular suggestion
	return 0, content
}

// formatOffset formats the offset for GitHub suggestion syntax
func formatOffset(offset int) string {
	if offset > 0 {
		return "+" + strconv.Itoa(offset)
	}
	return strconv.Itoa(offset)
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
