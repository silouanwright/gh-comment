package cmd

import (
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestInitColors(t *testing.T) {
	// Save original state
	originalColorEnabled := colorEnabled
	originalNoColor := noColor
	defer func() {
		colorEnabled = originalColorEnabled
		noColor = originalNoColor
		// Reset colors to original state
		InitColors()
	}()

	t.Run("colors enabled", func(t *testing.T) {
		// Set up environment for colors to be enabled
		testTerminalOverride = &[]bool{true}[0] // Simulate terminal
		noColor = false
		
		// Clear environment variables that would disable colors
		originalTerm := os.Getenv("TERM")
		originalNoColor := os.Getenv("NO_COLOR")
		os.Setenv("TERM", "xterm-256color")
		os.Unsetenv("NO_COLOR")
		defer func() {
			if originalTerm != "" {
				os.Setenv("TERM", originalTerm)
			} else {
				os.Unsetenv("TERM")
			}
			if originalNoColor != "" {
				os.Setenv("NO_COLOR", originalNoColor)
			} else {
				os.Unsetenv("NO_COLOR")
			}
		}()
		
		InitColors()

		// Check that all color objects are created
		assert.NotNil(t, ColorCommentID)
		assert.NotNil(t, ColorAuthor)
		assert.NotNil(t, ColorTimestamp)
		assert.NotNil(t, ColorFilePath)
		assert.NotNil(t, ColorLineNumber)
		assert.NotNil(t, ColorURL)
		assert.NotNil(t, ColorSuccess)
		assert.NotNil(t, ColorError)
		assert.NotNil(t, ColorWarning)
		assert.NotNil(t, ColorHeader)
		assert.NotNil(t, ColorReviewState)
		assert.NotNil(t, ColorCommitSHA)
		assert.NotNil(t, ColorIssueComment)
		assert.NotNil(t, ColorReviewComment)
		
		// Verify that colors are actually enabled (not NoColor)
		assert.False(t, color.NoColor)
	})

	t.Run("colors disabled", func(t *testing.T) {
		noColor = true
		colorEnabled = false
		InitColors()

		// Check that all color objects are still created (but disabled)
		assert.NotNil(t, ColorCommentID)
		assert.NotNil(t, ColorAuthor)
		assert.NotNil(t, ColorTimestamp)
		assert.NotNil(t, ColorFilePath)
		assert.NotNil(t, ColorLineNumber)
		assert.NotNil(t, ColorURL)
		assert.NotNil(t, ColorSuccess)
		assert.NotNil(t, ColorError)
		assert.NotNil(t, ColorWarning)
		assert.NotNil(t, ColorHeader)
		assert.NotNil(t, ColorReviewState)
		assert.NotNil(t, ColorCommitSHA)
		assert.NotNil(t, ColorIssueComment)
		assert.NotNil(t, ColorReviewComment)
		
		// Global color disable should be set
		assert.True(t, color.NoColor)
	})

	t.Run("colors enabled through ShouldUseColor detection", func(t *testing.T) {
		// Test the automatic detection path
		testTerminalOverride = &[]bool{true}[0] // Simulate terminal
		noColor = false
		// Don't preset colorEnabled - let InitColors() call ShouldUseColor()
		
		// Clear environment variables that would disable colors
		originalTerm := os.Getenv("TERM")
		originalNoColor := os.Getenv("NO_COLOR")
		os.Setenv("TERM", "xterm-256color")
		os.Unsetenv("NO_COLOR")
		defer func() {
			if originalTerm != "" {
				os.Setenv("TERM", originalTerm)
			} else {
				os.Unsetenv("TERM")
			}
			if originalNoColor != "" {
				os.Setenv("NO_COLOR", originalNoColor)
			} else {
				os.Unsetenv("NO_COLOR")
			}
		}()
		
		InitColors()

		// Colors should be enabled
		assert.True(t, colorEnabled)
		assert.False(t, color.NoColor)
		assert.NotNil(t, ColorSuccess)
	})

	t.Run("colors disabled through ShouldUseColor detection", func(t *testing.T) {
		// Test the automatic detection path with colors disabled
		testTerminalOverride = &[]bool{false}[0] // Simulate non-terminal
		noColor = false
		
		InitColors()

		// Colors should be disabled
		assert.False(t, colorEnabled)
		assert.True(t, color.NoColor)
		assert.NotNil(t, ColorError) // Still created, just disabled
	})
}

func TestIsTerminal(t *testing.T) {
	// Save original state
	originalTestTerminalOverride := testTerminalOverride
	defer func() {
		testTerminalOverride = originalTestTerminalOverride
	}()

	t.Run("returns override value when set to true", func(t *testing.T) {
		testTerminalOverride = &[]bool{true}[0]
		result := isTerminal()
		assert.True(t, result)
	})

	t.Run("returns override value when set to false", func(t *testing.T) {
		testTerminalOverride = &[]bool{false}[0]
		result := isTerminal()
		assert.False(t, result)
	})

	t.Run("uses actual terminal detection when override not set", func(t *testing.T) {
		testTerminalOverride = nil
		result := isTerminal()
		// We can't predict the actual result since it depends on the test environment,
		// but we can verify the function runs without error
		assert.IsType(t, true, result) // Just check it returns a bool
	})
}

func TestShouldUseColor(t *testing.T) {
	// Save original state
	originalNoColor := noColor
	originalTestTerminalOverride := testTerminalOverride
	defer func() {
		noColor = originalNoColor
		testTerminalOverride = originalTestTerminalOverride
	}()

	tests := []struct {
		name            string
		noColorFlag     bool
		noColorEnv      string
		termEnv         string
		terminalOverride *bool
		expected        bool
	}{
		{
			name:        "no color flag set",
			noColorFlag: true,
			terminalOverride: &[]bool{true}[0], // Terminal is available but --no-color overrides
			expected:    false,
		},
		{
			name:       "NO_COLOR environment variable set",
			noColorEnv: "1",
			terminalOverride: &[]bool{true}[0], // Terminal is available but NO_COLOR overrides
			expected:   false,
		},
		{
			name:     "TERM environment variable dumb",
			termEnv:  "dumb",
			terminalOverride: &[]bool{true}[0], // Terminal is available but TERM=dumb overrides
			expected: false,
		},
		{
			name:     "TERM environment variable empty",
			termEnv:  "",
			terminalOverride: &[]bool{true}[0], // Terminal is available but TERM="" overrides
			expected: false,
		},
		{
			name:        "colors should be enabled with normal terminal",
			noColorFlag: false,
			noColorEnv:  "",
			termEnv:     "xterm-256color",
			terminalOverride: &[]bool{true}[0], // Simulate terminal is available
			expected:    true,
		},
		{
			name:        "colors enabled with screen terminal",
			noColorFlag: false,
			noColorEnv:  "",
			termEnv:     "screen",
			terminalOverride: &[]bool{true}[0], // Simulate terminal is available
			expected:    true,
		},
		{
			name:        "colors enabled with tmux terminal",
			noColorFlag: false,
			noColorEnv:  "",
			termEnv:     "tmux-256color",
			terminalOverride: &[]bool{true}[0], // Simulate terminal is available
			expected:    true,
		},
		{
			name:        "colors disabled when not a terminal",
			noColorFlag: false,
			noColorEnv:  "",
			termEnv:     "xterm-256color",
			terminalOverride: &[]bool{false}[0], // Simulate not a terminal (pipe/redirect)
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			originalNoColorEnv := os.Getenv("NO_COLOR")
			originalTermEnv := os.Getenv("TERM")
			defer func() {
				if originalNoColorEnv != "" {
					os.Setenv("NO_COLOR", originalNoColorEnv)
				} else {
					os.Unsetenv("NO_COLOR")
				}
				if originalTermEnv != "" {
					os.Setenv("TERM", originalTermEnv)
				} else {
					os.Unsetenv("TERM")
				}
			}()

			// Set up test environment
			noColor = tt.noColorFlag
			testTerminalOverride = tt.terminalOverride
			
			// Handle NO_COLOR environment variable
			if tt.noColorEnv != "" {
				os.Setenv("NO_COLOR", tt.noColorEnv)
			} else {
				os.Unsetenv("NO_COLOR")
			}
			
			// Handle TERM environment variable  
			if tt.termEnv != "" {
				os.Setenv("TERM", tt.termEnv)
			} else if tt.name == "TERM environment variable empty" {
				os.Unsetenv("TERM")
			}

			result := ShouldUseColor()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDisableColors(t *testing.T) {
	// Save original state
	originalColorEnabled := colorEnabled
	originalNoColor := color.NoColor
	defer func() {
		colorEnabled = originalColorEnabled
		color.NoColor = originalNoColor
		InitColors()
	}()

	DisableColors()

	assert.False(t, colorEnabled)
	assert.True(t, color.NoColor)
}

func TestEnableColors(t *testing.T) {
	// Save original state
	originalColorEnabled := colorEnabled
	originalNoColor := color.NoColor
	originalNoColorFlag := noColor
	defer func() {
		colorEnabled = originalColorEnabled
		color.NoColor = originalNoColor
		noColor = originalNoColorFlag
		InitColors()
	}()

	// First disable colors
	DisableColors()
	assert.False(t, colorEnabled)
	assert.True(t, color.NoColor)

	// Then enable them
	EnableColors()
	assert.True(t, colorEnabled)
	assert.False(t, color.NoColor)
	assert.False(t, noColor) // noColor flag should also be false
}

func TestColorizeCommentType(t *testing.T) {
	tests := []struct {
		name         string
		commentType  string
		expectedText string
	}{
		{
			name:         "issue comment",
			commentType:  "issue",
			expectedText: "💬 General PR Comments",
		},
		{
			name:         "review comment",
			commentType:  "review",
			expectedText: "📋 Review Comments",
		},
		{
			name:         "unknown comment type",
			commentType:  "unknown",
			expectedText: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizeCommentType(tt.commentType)
			assert.Contains(t, result, tt.expectedText)
		})
	}
}

func TestColorizeReviewState(t *testing.T) {
	tests := []struct {
		name         string
		state        string
		expectedText string
	}{
		{
			name:         "approved state",
			state:        "approved",
			expectedText: "✅ approved",
		},
		{
			name:         "changes requested state",
			state:        "changes_requested",
			expectedText: "🔴 changes_requested",
		},
		{
			name:         "commented state",
			state:        "commented",
			expectedText: "💬 commented",
		},
		{
			name:         "unknown state",
			state:        "unknown",
			expectedText: "📝 unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizeReviewState(tt.state)
			assert.Contains(t, result, tt.expectedText)
		})
	}
}

func TestColorizeSuccess(t *testing.T) {
	text := "Test success message"
	result := ColorizeSuccess(text)
	assert.Contains(t, result, "✅")
	assert.Contains(t, result, text)
}

func TestColorizeError(t *testing.T) {
	text := "Test error message"
	result := ColorizeError(text)
	assert.Contains(t, result, "❌")
	assert.Contains(t, result, text)
}

func TestColorizeWarning(t *testing.T) {
	text := "Test warning message"  
	result := ColorizeWarning(text)
	assert.Contains(t, result, "⚠️")
	assert.Contains(t, result, text)
}

func TestColorFunctionsWithNilColors(t *testing.T) {
	// Test that color functions work even when colors aren't initialized
	
	// Save original colors
	originalColors := []interface{}{
		ColorSuccess, ColorError, ColorWarning, ColorReviewState,
		ColorIssueComment, ColorReviewComment,
	}
	
	// Set colors to nil
	ColorSuccess = nil
	ColorError = nil
	ColorWarning = nil
	ColorReviewState = nil
	ColorIssueComment = nil
	ColorReviewComment = nil
	
	defer func() {
		// Restore original colors
		ColorSuccess = originalColors[0].(*color.Color)
		ColorError = originalColors[1].(*color.Color)
		ColorWarning = originalColors[2].(*color.Color)
		ColorReviewState = originalColors[3].(*color.Color)
		ColorIssueComment = originalColors[4].(*color.Color)
		ColorReviewComment = originalColors[5].(*color.Color)
	}()

	t.Run("ColorizeSuccess with nil colors", func(t *testing.T) {
		result := ColorizeSuccess("test")
		assert.Contains(t, result, "✅ test")
	})

	t.Run("ColorizeError with nil colors", func(t *testing.T) {
		result := ColorizeError("test")
		assert.Contains(t, result, "❌ test")
	})

	t.Run("ColorizeWarning with nil colors", func(t *testing.T) {
		result := ColorizeWarning("test")
		assert.Contains(t, result, "⚠️ test")
	})

	t.Run("ColorizeReviewState with nil colors", func(t *testing.T) {
		result := ColorizeReviewState("approved")
		assert.Contains(t, result, "✅ approved")
	})

	t.Run("ColorizeCommentType with nil colors", func(t *testing.T) {
		result := ColorizeCommentType("issue")
		// Should fall back to uncolored text when colors are nil
		assert.Equal(t, "💬 General PR Comments", result)
		
		result2 := ColorizeCommentType("review")
		assert.Equal(t, "📋 Review Comments", result2)
		
		result3 := ColorizeCommentType("unknown")
		assert.Equal(t, "unknown", result3)
	})
}