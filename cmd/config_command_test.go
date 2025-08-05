package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunConfigInit(t *testing.T) {
	// Save original state
	originalGlobalFlag := globalFlag
	originalConfigFormat := configFormat
	defer func() {
		globalFlag = originalGlobalFlag
		configFormat = originalConfigFormat
	}()

	t.Run("creates local YAML config file", func(t *testing.T) {
		// Setup
		globalFlag = false
		configFormat = "yaml"

		// Create temporary directory
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		_ = os.Chdir(tempDir)                       // Test helper
		defer func() { _ = os.Chdir(originalWd) }() // Test cleanup

		// Run the command
		err := runConfigInit(configInitCmd, []string{})
		assert.NoError(t, err)

		// Check that file was created
		configFile := ".gh-comment.yaml"
		_, err = os.Stat(configFile)
		assert.NoError(t, err)

		// Check file content is valid YAML
		content, err := os.ReadFile(configFile)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "defaults:")
		assert.Contains(t, string(content), "behavior:")
	})

	t.Run("creates local JSON config file", func(t *testing.T) {
		// Setup
		globalFlag = false
		configFormat = "json"

		// Create temporary directory
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		_ = os.Chdir(tempDir)                       // Test helper
		defer func() { _ = os.Chdir(originalWd) }() // Test cleanup

		// Run the command
		err := runConfigInit(configInitCmd, []string{})
		assert.NoError(t, err)

		// Check that file was created
		configFile := ".gh-comment.json"
		_, err = os.Stat(configFile)
		assert.NoError(t, err)

		// Check file content is valid JSON
		content, err := os.ReadFile(configFile)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "\"defaults\":")
		assert.Contains(t, string(content), "\"behavior\":")
	})

	t.Run("creates global config file", func(t *testing.T) {
		// Setup
		globalFlag = true
		configFormat = "yaml"

		// Create temporary home directory
		tempDir := t.TempDir()
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tempDir)
		defer os.Setenv("HOME", originalHome)

		// Run the command
		err := runConfigInit(configInitCmd, []string{})
		assert.NoError(t, err)

		// Check that global file was created
		expectedPath := filepath.Join(tempDir, ".config", "gh-comment", "config.yaml")
		_, err = os.Stat(expectedPath)
		assert.NoError(t, err)
	})

	t.Run("fails if config file already exists", func(t *testing.T) {
		// Setup
		globalFlag = false
		configFormat = "yaml"

		// Create temporary directory with existing config file
		tempDir := t.TempDir()
		originalWd, _ := os.Getwd()
		_ = os.Chdir(tempDir)                       // Test helper
		defer func() { _ = os.Chdir(originalWd) }() // Test cleanup

		// Create existing config file
		configFile := ".gh-comment.yaml"
		err := os.WriteFile(configFile, []byte("existing"), 0644)
		assert.NoError(t, err)

		// Run the command
		err = runConfigInit(configInitCmd, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config file already exists")
	})
}

func TestRunConfigShow(t *testing.T) {
	// Save original state
	originalGlobalConfig := globalConfig
	originalShowEffective := showEffective
	originalShowSource := showSource
	defer func() {
		globalConfig = originalGlobalConfig
		showEffective = originalShowEffective
		showSource = originalShowSource
	}()

	t.Run("shows current config (default behavior)", func(t *testing.T) {
		// Reset flags
		showEffective = false
		showSource = false

		// Set up test config
		testConfig := &Config{
			Defaults: DefaultsConfig{
				Author: "test-user",
			},
			Display: DisplayConfig{
				Format: "json",
			},
		}
		globalConfig = testConfig

		// Capture output (this would normally print to stdout)
		// We just test that it doesn't error
		err := runConfigShow(configShowCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("shows effective config with --effective flag", func(t *testing.T) {
		// Set showEffective flag
		showEffective = true
		showSource = false

		// Set up test config
		testConfig := &Config{
			Defaults: DefaultsConfig{
				Author: "effective-user",
			},
			Display: DisplayConfig{
				Format: "table",
			},
		}
		globalConfig = testConfig

		// Run the command
		err := runConfigShow(configShowCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("shows source information with --source flag", func(t *testing.T) {
		// Set showSource flag
		showEffective = false
		showSource = true

		// Set up test config
		testConfig := &Config{
			Defaults: DefaultsConfig{
				Author: "source-user",
			},
		}
		globalConfig = testConfig

		// Run the command
		err := runConfigShow(configShowCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("shows source info when config file exists", func(t *testing.T) {
		// Set showSource flag
		showEffective = false
		showSource = true

		// Create temporary config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".gh-comment.yaml")
		configContent := `defaults:
  author: "file-user"`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Change to temp directory so findConfigFile() finds the file
		originalWd, _ := os.Getwd()
		_ = os.Chdir(tempDir)                       // Test helper
		defer func() { _ = os.Chdir(originalWd) }() // Test cleanup

		// Clear global config so it loads from file
		globalConfig = nil

		// Run the command
		err = runConfigShow(configShowCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("shows config from file when no global config", func(t *testing.T) {
		// Reset flags
		showEffective = false
		showSource = false

		// Clear global config
		globalConfig = nil

		// Create temporary config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".gh-comment.yaml")
		configContent := `defaults:
  author: "file-user"
display:
  format: "table"`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Change to temp directory
		originalWd, _ := os.Getwd()
		_ = os.Chdir(tempDir)                       // Test helper
		defer func() { _ = os.Chdir(originalWd) }() // Test cleanup

		// Run the command
		err = runConfigShow(configShowCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("handles both --effective and --source flags", func(t *testing.T) {
		// Set both flags (effective takes precedence)
		showEffective = true
		showSource = true

		testConfig := &Config{
			Defaults: DefaultsConfig{
				Author: "both-flags-user",
			},
		}
		globalConfig = testConfig

		// Run the command - should only show effective config
		err := runConfigShow(configShowCmd, []string{})
		assert.NoError(t, err)
	})
}

func TestRunConfigValidate(t *testing.T) {
	t.Run("validates existing config file", func(t *testing.T) {
		// Create temporary valid config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".gh-comment.yaml")
		configContent := `defaults:
  author: "test-user"
  repository: "owner/repo"
behavior:
  dry_run: false
  verbose: false
display:
  format: "table"
  quiet: false
filters:
  status: "all"
  type: "all"`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Change to temp directory
		originalWd, _ := os.Getwd()
		_ = os.Chdir(tempDir)                       // Test helper
		defer func() { _ = os.Chdir(originalWd) }() // Test cleanup

		// Run validation
		err = runConfigValidate(configValidateCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("validates specific config file", func(t *testing.T) {
		// Create temporary valid config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test-config.yaml")
		configContent := `defaults:
  author: "test-user"`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Run validation with specific file
		err = runConfigValidate(configValidateCmd, []string{configFile})
		assert.NoError(t, err)
	})

	t.Run("fails on invalid config file", func(t *testing.T) {
		// Create temporary invalid config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "invalid-config.yaml")
		configContent := `invalid: [unclosed bracket`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		assert.NoError(t, err)

		// Run validation
		err = runConfigValidate(configValidateCmd, []string{configFile})
		assert.Error(t, err)
	})
}

func TestShowConfigWarnings(t *testing.T) {
	t.Run("shows warnings for config issues", func(t *testing.T) {
		testConfig := &Config{
			Defaults: DefaultsConfig{
				Repository: "invalid-repo-format", // Should trigger warning
			},
		}

		// This function prints warnings, so we just test it doesn't panic
		assert.NotPanics(t, func() {
			showConfigWarnings(testConfig)
		})
	})

	t.Run("handles nil config gracefully", func(t *testing.T) {
		assert.NotPanics(t, func() {
			showConfigWarnings(nil)
		})
	})

	t.Run("handles config with valid repository", func(t *testing.T) {
		testConfig := &Config{
			Defaults: DefaultsConfig{
				Repository: "owner/repo", // Valid format
			},
		}

		assert.NotPanics(t, func() {
			showConfigWarnings(testConfig)
		})
	})
}
