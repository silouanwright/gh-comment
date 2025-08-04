package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrompt(t *testing.T) {
	tests := []struct {
		name             string
		promptName       string
		expectExists     bool
		expectedTitle    string
		expectedCategory string
	}{
		{
			name:             "existing security prompt",
			promptName:       "security-audit",
			expectExists:     true,
			expectedTitle:    "Comprehensive Security Audit Review",
			expectedCategory: "security",
		},
		{
			name:             "existing AI prompt",
			promptName:       "ai-assistant",
			expectExists:     true,
			expectedTitle:    "AI Assistant Code Review Template",
			expectedCategory: "ai",
		},
		{
			name:             "existing performance prompt",
			promptName:       "performance-optimization",
			expectExists:     true,
			expectedTitle:    "Performance Optimization Review",
			expectedCategory: "performance",
		},
		{
			name:         "non-existing prompt",
			promptName:   "non-existent",
			expectExists: false,
		},
		{
			name:         "empty prompt name",
			promptName:   "",
			expectExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, exists := getPrompt(tt.promptName)

			assert.Equal(t, tt.expectExists, exists)

			if tt.expectExists {
				assert.Equal(t, tt.expectedTitle, prompt.Title)
				assert.Equal(t, tt.expectedCategory, prompt.Category)
				assert.NotEmpty(t, prompt.Content)
				assert.NotEmpty(t, prompt.EstimatedTime)
				assert.NotEmpty(t, prompt.Tags)
			}
		})
	}
}

func TestGetAllPrompts(t *testing.T) {
	prompts := getAllPrompts()

	// Verify we have the expected number of prompts
	assert.GreaterOrEqual(t, len(prompts), 6, "Should have at least 6 prompts")

	// Verify each prompt has required fields
	for name, prompt := range prompts {
		assert.NotEmpty(t, prompt.Name, "Prompt name should not be empty")
		assert.Equal(t, name, prompt.Name, "Map key should match prompt name")
		assert.NotEmpty(t, prompt.Title, "Prompt title should not be empty")
		assert.NotEmpty(t, prompt.Category, "Prompt category should not be empty")
		assert.NotEmpty(t, prompt.Content, "Prompt content should not be empty")
		assert.NotEmpty(t, prompt.EstimatedTime, "Prompt estimated time should not be empty")
		assert.NotEmpty(t, prompt.Tags, "Prompt tags should not be empty")

		// Verify content contains framework sections
		assert.Contains(t, prompt.Content, "##", "Prompt should have section headers")

		// Verify tags are reasonable
		assert.LessOrEqual(t, len(prompt.Tags), 6, "Should not have too many tags")
	}

	// Verify specific prompts exist
	expectedPrompts := []string{
		"security-audit",
		"performance-optimization",
		"architecture-review",
		"code-quality",
		"ai-assistant",
		"migration-review",
	}

	for _, expectedName := range expectedPrompts {
		_, exists := prompts[expectedName]
		assert.True(t, exists, "Expected prompt '%s' should exist", expectedName)
	}
}

func TestPromptCategories(t *testing.T) {
	prompts := getAllPrompts()
	categories := make(map[string]int)

	for _, prompt := range prompts {
		categories[prompt.Category]++
	}

	// Verify we have the expected categories
	expectedCategories := []string{"security", "performance", "architecture", "quality", "ai"}
	for _, category := range expectedCategories {
		count, exists := categories[category]
		assert.True(t, exists, "Category '%s' should exist", category)
		assert.Greater(t, count, 0, "Category '%s' should have at least one prompt", category)
	}
}

func TestPromptContentQuality(t *testing.T) {
	prompts := getAllPrompts()

	for name, prompt := range prompts {
		t.Run(name, func(t *testing.T) {
			// Verify emojis are used for communication style (at least one of the CREG system)
			hasEmoji := strings.Contains(prompt.Content, "🔧") ||
				strings.Contains(prompt.Content, "🤔") ||
				strings.Contains(prompt.Content, "♻️") ||
				strings.Contains(prompt.Content, "📝") ||
				strings.Contains(prompt.Content, "😃")
			assert.True(t, hasEmoji, "Should contain at least one CREG emoji (🔧🤔♻️📝😃)")

			// Verify framework structure
			assert.Contains(t, prompt.Content, "##", "Should have framework sections")

			// Verify it mentions [SUGGEST:] syntax (except for specific prompts)
			if name != "ai-assistant" && name != "code-quality" {
				// Some prompts are more general and don't need to mention suggestions specifically
				assert.Contains(t, prompt.Content, "[SUGGEST:", "Should mention suggestion syntax")
			}

			// Verify reasonable content length (not too short)
			assert.Greater(t, len(prompt.Content), 500, "Content should be substantial")

			// Verify estimated time format
			assert.Contains(t, prompt.EstimatedTime, "minute", "Estimated time should mention minutes")
		})
	}
}

func TestPromptConsistency(t *testing.T) {
	prompts := getAllPrompts()

	for name, prompt := range prompts {
		t.Run(name, func(t *testing.T) {
			// Verify naming consistency
			assert.Equal(t, name, prompt.Name, "Map key should match prompt.Name")

			// Verify title format (should be title case)
			assert.NotContains(t, prompt.Title, "  ", "Title should not have double spaces")
			assert.Greater(t, len(prompt.Title), 10, "Title should be descriptive")

			// Verify category is lowercase
			assert.Equal(t, prompt.Category, strings.ToLower(prompt.Category), "Category should be lowercase")

			// Verify tags are reasonable
			for _, tag := range prompt.Tags {
				assert.NotEmpty(t, tag, "Tags should not be empty")
				assert.Equal(t, tag, strings.ToLower(tag), "Tags should be lowercase")
			}
		})
	}
}

func TestListAvailablePrompts(t *testing.T) {
	// Save original state
	originalCategory := promptCategory
	defer func() { promptCategory = originalCategory }()

	tests := []struct {
		name                 string
		category            string
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name:             "list all prompts",
			category:         "",
			expectedContains: []string{
				"📋 **Available Code Review Prompts**",
				"security-audit",
				"performance-optimization",
				"architecture-review",
				"## Security",
				"## Performance",
				"## Architecture",
			},
		},
		{
			name:             "filter by security category",
			category:         "security",
			expectedContains: []string{
				"📋 **Code Review Prompts - Security Category**",
				"security-audit",
				"## Security",
			},
			expectedNotContains: []string{
				"performance-optimization",
				"## Performance",
			},
		},
		{
			name:             "filter by performance category",
			category:         "performance",
			expectedContains: []string{
				"📋 **Code Review Prompts - Performance Category**",
				"performance-optimization",
				"## Performance",
			},
			expectedNotContains: []string{
				"security-audit",
				"## Security",
			},
		},
		{
			name:             "filter by ai category",
			category:         "ai",
			expectedContains: []string{
				"📋 **Code Review Prompts - Ai Category**",
				"ai-assistant",
				"## Ai",
			},
			expectedNotContains: []string{
				"security-audit",
				"performance-optimization",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set category filter
			promptCategory = tt.category

			// Capture output
			output := captureOutputFromPrompts(func() {
				err := listAvailablePrompts()
				assert.NoError(t, err)
			})

			// Check expected content
			for _, expected := range tt.expectedContains {
				assert.Contains(t, output, expected, "Should contain: %s", expected)
			}

			// Check unwanted content is not present
			for _, notExpected := range tt.expectedNotContains {
				assert.NotContains(t, output, notExpected, "Should not contain: %s", notExpected)
			}
		})
	}
}

func TestRunPrompts(t *testing.T) {
	// Save original state
	originalListPrompts := listPrompts
	originalCategory := promptCategory
	defer func() {
		listPrompts = originalListPrompts
		promptCategory = originalCategory
	}()

	tests := []struct {
		name             string
		args             []string
		listFlag         bool
		expectedContains []string
		wantError        bool
	}{
		{
			name:     "list prompts with flag",
			args:     []string{},
			listFlag: true,
			expectedContains: []string{
				"📋 **Available Code Review Prompts**",
				"security-audit",
			},
		},
		{
			name: "list prompts with no args",
			args: []string{},
			expectedContains: []string{
				"📋 **Available Code Review Prompts**",
				"security-audit",
			},
		},
		{
			name: "get specific prompt",
			args: []string{"security-audit"},
			expectedContains: []string{
				"📋 **Comprehensive Security Audit Review**",
				"🎯 **Category**: security",
				"⏱️  **Estimated Time**:",
				"📝 **Prompt:**",
				"```",
			},
		},
		{
			name: "get specific AI prompt",
			args: []string{"ai-assistant"},
			expectedContains: []string{
				"📋 **AI Assistant Code Review Template**",
				"🎯 **Category**: ai",
				"⏱️  **Estimated Time**:",
				"📝 **Prompt:**",
			},
		},
		{
			name: "non-existent prompt shows error and lists available",
			args: []string{"non-existent-prompt"},
			expectedContains: []string{
				"❌ Prompt 'non-existent-prompt' not found.",
				"Available prompts:",
				"📋 **Available Code Review Prompts**",
				"security-audit",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			listPrompts = tt.listFlag
			promptCategory = ""

			// Capture output
			output := captureOutputFromPrompts(func() {
				err := runPrompts(nil, tt.args)
				if tt.wantError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})

			// Check expected content
			for _, expected := range tt.expectedContains {
				assert.Contains(t, output, expected, "Should contain: %s", expected)
			}
		})
	}
}

// Helper function to capture stdout output specifically for prompts tests
func captureOutputFromPrompts(fn func()) string {
	// Save the original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a buffer to capture the output
	outputChan := make(chan string)
	go func() {
		var buf strings.Builder
		buffer := make([]byte, DefaultBufferSize)
		for {
			n, err := r.Read(buffer)
			if n > 0 {
				buf.Write(buffer[:n])
			}
			if err != nil {
				break
			}
		}
		outputChan <- buf.String()
	}()

	// Execute the function
	fn()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Get the captured output
	return <-outputChan
}
