package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandInlineSuggestions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single inline suggestion",
			input:    "Here's a fix: [SUGGEST: const MAX = 100]",
			expected: "Here's a fix: \n\n```suggestion\nconst MAX = 100\n```\n\n",
		},
		{
			name:     "multiple inline suggestions",
			input:    "Fix 1: [SUGGEST: var x = 1] and Fix 2: [SUGGEST: var y = 2]",
			expected: "Fix 1: \n\n```suggestion\nvar x = 1\n```\n\n and Fix 2: \n\n```suggestion\nvar y = 2\n```\n\n",
		},
		{
			name:     "suggestion with extra spaces",
			input:    "[SUGGEST:   spaced code   ]",
			expected: "\n\n```suggestion\nspaced code\n```\n\n",
		},
		{
			name:     "nested brackets in suggestion",
			input:    "[SUGGEST: arr[0] = func() { return true }]",
			expected: "\n\n```suggestion\narr[0] = func() { return true }\n```\n\n", // Fixed: Now handles nested brackets correctly
		},
		{
			name:     "empty suggestion",
			input:    "[SUGGEST: ]",
			expected: "\n\n```suggestion\n\n```\n\n",
		},
		{
			name:     "no suggestions",
			input:    "Regular comment without suggestions",
			expected: "Regular comment without suggestions",
		},
		{
			name:     "incomplete suggestion syntax",
			input:    "[SUGGEST: incomplete",
			expected: "[SUGGEST: incomplete",
		},
		{
			name:     "suggestion at start",
			input:    "[SUGGEST: fix] is the solution",
			expected: "\n\n```suggestion\nfix\n```\n\n is the solution",
		},
		{
			name:     "suggestion at end",
			input:    "The solution is [SUGGEST: fix]",
			expected: "The solution is \n\n```suggestion\nfix\n```\n\n",
		},
		{
			name:     "suggestion with newlines inside",
			input:    "[SUGGEST: line1\nline2]",
			expected: "\n\n```suggestion\nline1\nline2\n```\n\n", // Current regex does match newlines
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandInlineSuggestions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandMultilineSuggestions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple multiline suggestion",
			input: `Here's a fix:
<<<SUGGEST
const MAX = 100
const MIN = 0
SUGGEST>>>`,
			expected: `Here's a fix:
` + "\n\n```suggestion\nconst MAX = 100\nconst MIN = 0\n```\n\n",
		},
		{
			name: "multiple multiline suggestions",
			input: `Fix 1:
<<<SUGGEST
var x = 1
SUGGEST>>>
Fix 2:
<<<SUGGEST
var y = 2
SUGGEST>>>`,
			expected: `Fix 1:
` + "\n\n```suggestion\nvar x = 1\n```\n\n" + `
Fix 2:
` + "\n\n```suggestion\nvar y = 2\n```\n\n",
		},
		{
			name: "multiline with extra whitespace",
			input: `<<<SUGGEST   
code with spaces
SUGGEST>>>`,
			expected: "\n\n```suggestion\ncode with spaces\n```\n\n",
		},
		{
			name: "empty multiline suggestion",
			input: `<<<SUGGEST

SUGGEST>>>`,
			expected: "\n\n```suggestion\n\n```\n\n",
		},
		{
			name:     "no multiline suggestions",
			input:    "Regular comment without suggestions",
			expected: "Regular comment without suggestions",
		},
		{
			name: "incomplete multiline syntax - missing end",
			input: `<<<SUGGEST
some code`,
			expected: `<<<SUGGEST
some code`,
		},
		{
			name: "incomplete multiline syntax - missing start",
			input: `some code
SUGGEST>>>`,
			expected: `some code
SUGGEST>>>`,
		},
		{
			name: "nested code blocks",
			input: `<<<SUGGEST
if (true) {
    console.log("nested");
}
SUGGEST>>>`,
			expected: "\n\n```suggestion\nif (true) {\n    console.log(\"nested\");\n}\n```\n\n",
		},
		{
			name: "suggestion with special characters",
			input: `<<<SUGGEST
const regex = /[A-Z]+/g;
const str = "Hello $USER!";
SUGGEST>>>`,
			expected: "\n\n```suggestion\nconst regex = /[A-Z]+/g;\nconst str = \"Hello $USER!\";\n```\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandMultilineSuggestions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandSuggestions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "both inline and multiline suggestions",
			input: `Here's an inline fix: [SUGGEST: var x = 1]
And a multiline fix:
<<<SUGGEST
function foo() {
    return true;
}
SUGGEST>>>`,
			expected: `Here's an inline fix: ` + "\n\n```suggestion\nvar x = 1\n```\n\n" + `
And a multiline fix:
` + "\n\n```suggestion\nfunction foo() {\n    return true;\n}\n```\n\n",
		},
		{
			name: "multiline processed before inline",
			input: `<<<SUGGEST
[SUGGEST: inner]
SUGGEST>>>`,
			expected: "\n\n```suggestion\n\n\n```suggestion\ninner\n```\n\n\n```\n\n", // Inner suggestion gets expanded too
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name: "complex mixed example",
			input: `Review feedback:
1. Update constant: [SUGGEST: const MAX_RETRIES = 5]
2. Refactor function:
<<<SUGGEST
async function fetchData(url) {
    try {
        const response = await fetch(url);
        return await response.json();
    } catch (error) {
        console.error('Fetch failed:', error);
        return null;
    }
}
SUGGEST>>>
3. Fix typo: [SUGGEST: // Fixed typo in comment]`,
			expected: `Review feedback:
1. Update constant: ` + "\n\n```suggestion\nconst MAX_RETRIES = 5\n```\n\n" + `
2. Refactor function:
` + "\n\n```suggestion\nasync function fetchData(url) {\n    try {\n        const response = await fetch(url);\n        return await response.json();\n    } catch (error) {\n        console.error('Fetch failed:', error);\n        return null;\n    }\n}\n```\n\n" + `
3. Fix typo: ` + "\n\n```suggestion\n// Fixed typo in comment\n```\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandSuggestions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test edge cases that might cause regex issues
func TestExpandSuggestionsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "suggestion with regex special chars",
			input:    "[SUGGEST: regex = /.*+?[]{}()|\\^$/]",
			expected: "\n\n```suggestion\nregex = /.*+?[]{}()|\\^$/\n```\n\n", // Fixed: Now handles special chars correctly
		},
		{
			name:     "very long suggestion",
			input:    "[SUGGEST: " + string(make([]byte, 1000, 1000)) + "]",
			expected: "\n\n```suggestion\n" + string(make([]byte, 1000, 1000)) + "\n```\n\n",
		},
		{
			name:     "unicode in suggestions",
			input:    "[SUGGEST: const message = '你好世界 🌍']",
			expected: "\n\n```suggestion\nconst message = '你好世界 🌍'\n```\n\n",
		},
		{
			name: "malformed nesting attempts",
			input: `[SUGGEST: <<<SUGGEST
nested
SUGGEST>>>]`,
			expected: "\n\n```suggestion\n```suggestion\nnested\n```\n```\n\n", // Both expansions happen
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandSuggestions(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestExpandInlineSuggestionsRegressionCases - REGRESSION TESTS
// These tests prevent the bracket parsing bugs from reoccurring
func TestExpandInlineSuggestionsRegressionCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "deeply nested brackets - REGRESSION TEST for bracket counting",
			input:    "[SUGGEST: obj[key[nested[deep]]]]",
			expected: "\n\n```suggestion\nobj[key[nested[deep]]]\n```\n\n",
		},
		{
			name:     "mixed bracket types - REGRESSION TEST",
			input:    "[SUGGEST: arr[obj.prop] = func() { return data[idx]; }]",
			expected: "\n\n```suggestion\narr[obj.prop] = func() { return data[idx]; }\n```\n\n",
		},
		{
			name:     "json-like structure - REGRESSION TEST",
			input:    "[SUGGEST: const config = { items: [1, 2, [3, 4]], nested: { deep: [5] } }]",
			expected: "\n\n```suggestion\nconst config = { items: [1, 2, [3, 4]], nested: { deep: [5] } }\n```\n\n",
		},
		{
			name:     "regex with brackets - REGRESSION TEST for special chars",
			input:    "[SUGGEST: const pattern = /[a-zA-Z0-9\\[\\]]+/g]",
			expected: "\n\n```suggestion\nconst pattern = /[a-zA-Z0-9\\[\\]]+/g\n```\n\n",
		},
		{
			name:     "escaped brackets in strings - REGRESSION TEST",
			input:    `[SUGGEST: const msg = "Use brackets [like this\] in text"]`,
			expected: "\n\n```suggestion\nconst msg = \"Use brackets [like this\\] in text\"\n```\n\n",
		},
		{
			name:     "multiple suggestions with complex brackets - REGRESSION TEST",
			input:    "First: [SUGGEST: arr[0][1]] and Second: [SUGGEST: obj[key][prop]]",
			expected: "First: \n\n```suggestion\narr[0][1]\n```\n\n and Second: \n\n```suggestion\nobj[key][prop]\n```\n\n",
		},
		{
			name:     "array destructuring - REGRESSION TEST for complex syntax",
			input:    "[SUGGEST: const [first, ...rest] = items[index][subIndex]]",
			expected: "\n\n```suggestion\nconst [first, ...rest] = items[index][subIndex]\n```\n\n",
		},
		{
			name:     "typescript generic with brackets - REGRESSION TEST",
			input:    "[SUGGEST: const result: Array<Map<string, object[]>> = data[key]]",
			expected: "\n\n```suggestion\nconst result: Array<Map<string, object[]>> = data[key]\n```\n\n",
		},
		{
			name:     "sql-like syntax with brackets - REGRESSION TEST",
			input:    "[SUGGEST: SELECT * FROM table WHERE col IN [1,2,3] AND other_col = data[idx]]",
			expected: "\n\n```suggestion\nSELECT * FROM table WHERE col IN [1,2,3] AND other_col = data[idx]\n```\n\n",
		},
		{
			name:     "unterminated suggestion - REGRESSION TEST for incomplete syntax",
			input:    "Start [SUGGEST: incomplete and [SUGGEST: complete]",
			expected: "Start [SUGGEST: incomplete and [SUGGEST: complete]", // Correctly leaves incomplete suggestions unchanged
		},
		{
			name:     "suggestion with line breaks - REGRESSION TEST for multiline content",
			input:    "[SUGGEST: const multiline = [\n  item1,\n  item2[index]\n]]",
			expected: "\n\n```suggestion\nconst multiline = [\n  item1,\n  item2[index]\n]\n```\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandInlineSuggestions(tt.input)
			assert.Equal(t, tt.expected, result, "Regression test failed - bracket counting parser may be broken")
		})
	}
}
