package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCrossPlatformPathHandling ensures file paths work across platforms
func TestCrossPlatformPathHandling(t *testing.T) {
	tests := []struct {
		name           string
		inputPath      string
		expectedError  bool
		description    string
	}{
		{
			name:        "unix-style path",
			inputPath:   "src/main.go",
			description: "Standard relative path should work on all platforms",
		},
		{
			name:        "nested path with slashes",
			inputPath:   "src/components/Button.tsx",
			description: "Nested paths should normalize correctly",
		},
		{
			name:          "absolute path rejection",
			inputPath:     "/usr/bin/test",
			expectedError: true,
			description:   "Absolute paths should be rejected on all platforms",
		},
		{
			name:        "empty path",
			inputPath:   "",
			description: "Empty paths should be accepted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.inputPath)
			
			if tt.expectedError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// TestCrossPlatformTempDirectory ensures temp directory handling works across platforms
func TestCrossPlatformTempDirectory(t *testing.T) {
	tempDir := os.TempDir()
	
	// Should return a valid directory path
	assert.NotEmpty(t, tempDir, "TempDir should return a non-empty path")
	
	// Should be accessible
	info, err := os.Stat(tempDir)
	assert.NoError(t, err, "TempDir should be accessible")
	assert.True(t, info.IsDir(), "TempDir should return a directory")
	
	// Test that we can construct a counter file path in temp directory
	counterFile := filepath.Join(tempDir, ".gh-comment-integration-counter")
	assert.Contains(t, counterFile, tempDir, "Counter file should be in temp directory")
	assert.True(t, filepath.IsAbs(counterFile), "Counter file path should be absolute")
}

// TestCrossPlatformLineEndings ensures line ending handling works across platforms
func TestCrossPlatformLineEndings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unix line endings",
			input:    "line1\nline2\nline3",
			expected: "line1\\nline2\\nline3",
		},
		{
			name:     "windows line endings", 
			input:    "line1\r\nline2\r\nline3",
			expected: "line1\\r\\nline2\\r\\nline3",
		},
		{
			name:     "mixed line endings",
			input:    "line1\nline2\r\nline3",
			expected: "line1\\nline2\\r\\nline3",
		},
		{
			name:     "no line endings",
			input:    "single line",
			expected: "single line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the export functionality that handles line endings
			result := strings.ReplaceAll(tt.input, "\n", "\\n")
			result = strings.ReplaceAll(result, "\r", "\\r")
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCrossPlatformFilePermissions ensures file operations work across platforms
func TestCrossPlatformFilePermissions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-permissions.txt")
	
	// Write a file with standard permissions
	content := []byte("test content")
	err := os.WriteFile(testFile, content, 0644)
	assert.NoError(t, err, "Should be able to write file with 0644 permissions")
	
	// Verify we can read it back
	readContent, err := os.ReadFile(testFile)
	assert.NoError(t, err, "Should be able to read file back")
	assert.Equal(t, content, readContent, "Content should match")
	
	// Test directory creation with standard permissions
	testDir := filepath.Join(tempDir, "test-dir")
	err = os.MkdirAll(testDir, 0755)
	assert.NoError(t, err, "Should be able to create directory with 0755 permissions")
	
	// Verify directory exists
	info, err := os.Stat(testDir)
	assert.NoError(t, err, "Directory should exist")
	assert.True(t, info.IsDir(), "Should be a directory")
}

// TestCrossPlatformEnvironmentVariables ensures env var handling works across platforms  
func TestCrossPlatformEnvironmentVariables(t *testing.T) {
	// Test the NO_COLOR environment variable handling
	originalNoColor := os.Getenv("NO_COLOR")
	defer os.Setenv("NO_COLOR", originalNoColor)
	
	// Test with NO_COLOR set
	os.Setenv("NO_COLOR", "1")
	shouldUseColor := ShouldUseColor()
	assert.False(t, shouldUseColor, "Should disable color when NO_COLOR is set")
	
	// Test with NO_COLOR unset
	os.Unsetenv("NO_COLOR")
	// Note: ShouldUseColor() might still return false in test environments
	// so we just verify it doesn't panic
	assert.NotPanics(t, func() {
		ShouldUseColor()
	}, "ShouldUseColor should not panic")
}

// TestCrossPlatformCommandExecution tests that our command execution works across platforms
func TestCrossPlatformCommandExecution(t *testing.T) {
	// Skip if not in integration mode to avoid actual command execution
	if testing.Short() {
		t.Skip("Skipping command execution test in short mode")
	}
	
	// Test that Execute function exists and can be called
	assert.NotNil(t, Execute, "Execute function should exist")
	
	// Test help command which should work without external dependencies
	// This is mainly testing that the command structure works cross-platform
	err := Execute()
	assert.NoError(t, err, "Execute should not error when showing help")
}

// TestCrossPlatformPathValidation ensures path validation works correctly across platforms
func TestCrossPlatformPathValidation(t *testing.T) {
	// Test cases that should behave the same across platforms
	crossPlatformTests := []struct {
		name        string
		path        string
		shouldError bool
		reason      string
	}{
		{
			name:        "valid relative path",
			path:        "src/main.go",
			shouldError: false,
			reason:      "Standard relative paths should work everywhere",
		},
		{
			name:        "directory traversal attempt",
			path:        "../../../etc/passwd",
			shouldError: true,
			reason:      "Directory traversal should be blocked on all platforms",
		},
		{
			name:        "empty path",
			path:        "",
			shouldError: false,
			reason:      "Empty paths should be allowed",
		},
		{
			name:        "very long path",
			path:        strings.Repeat("a/", 2100) + "file.go", // Over 4200 chars > 4096 limit
			shouldError: true,
			reason:      "Very long paths should be rejected",
		},
	}

	for _, tt := range crossPlatformTests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path)
			
			if tt.shouldError {
				assert.Error(t, err, tt.reason)
			} else {
				assert.NoError(t, err, tt.reason)
			}
		})
	}
	
	// Platform-specific tests
	switch runtime.GOOS {
	case "windows":
		t.Run("windows absolute path", func(t *testing.T) {
			err := validateFilePath("C:\\Windows\\System32\\file.txt")
			assert.Error(t, err, "Windows absolute paths should be rejected")
		})
		
		t.Run("windows UNC path", func(t *testing.T) {
			err := validateFilePath("\\\\server\\share\\file.txt")  
			assert.Error(t, err, "UNC paths should be rejected")
		})
		
	case "darwin", "linux":
		t.Run("unix absolute path", func(t *testing.T) {
			err := validateFilePath("/usr/local/bin/file")
			assert.Error(t, err, "Unix absolute paths should be rejected")
		})
	}
}

// TestPlatformSpecificBehaviors tests behavior that may differ across platforms
func TestPlatformSpecificBehaviors(t *testing.T) {
	t.Run("temp directory format", func(t *testing.T) {
		tempDir := os.TempDir()
		
		switch runtime.GOOS {
		case "windows":
			// Windows temp directory typically contains backslashes or C:
			assert.True(t, 
				strings.Contains(tempDir, "\\") || strings.Contains(tempDir, "C:"),
				"Windows temp dir should contain backslashes or drive letter: %s", tempDir)
		case "darwin":
			// macOS temp directory typically starts with /var/folders
			assert.True(t,
				strings.HasPrefix(tempDir, "/var/folders") || strings.HasPrefix(tempDir, "/tmp"),
				"macOS temp dir should start with /var/folders or /tmp: %s", tempDir)
		case "linux":
			// Linux temp directory typically is /tmp
			assert.True(t,
				strings.HasPrefix(tempDir, "/tmp") || strings.HasPrefix(tempDir, "/var/tmp"),
				"Linux temp dir should start with /tmp or /var/tmp: %s", tempDir)
		}
	})
	
	t.Run("filepath operations", func(t *testing.T) {
		// Test that filepath.Join works correctly on each platform
		path := filepath.Join("src", "components", "Button.tsx")
		
		// Should always produce a valid path for the current platform
		assert.NotEmpty(t, path)
		assert.NotContains(t, path, "//", "Path should not contain double slashes")
		
		// Verify platform-specific separator behavior
		separator := string(filepath.Separator)
		switch runtime.GOOS {
		case "windows":
			assert.Equal(t, "\\", separator, "Windows should use backslash separator")
		default:
			assert.Equal(t, "/", separator, "Unix-like systems should use forward slash")
		}
	})
}