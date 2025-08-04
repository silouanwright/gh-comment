package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	// Test defaults
	assert.Equal(t, "", config.Defaults.Repository)
	assert.Equal(t, "", config.Defaults.Author)
	assert.Equal(t, 0, config.Defaults.PR)

	assert.Equal(t, false, config.Behavior.DryRun)
	assert.Equal(t, false, config.Behavior.Verbose)
	assert.Equal(t, true, config.Behavior.Validate)
	assert.Equal(t, false, config.Behavior.NoExpandSuggestions)

	assert.Equal(t, "table", config.Display.Format)
	assert.Equal(t, "auto", config.Display.Color)
	assert.Equal(t, false, config.Display.Quiet)

	assert.Equal(t, "all", config.Filters.Status)
	assert.Equal(t, "all", config.Filters.Type)
	assert.Equal(t, "", config.Filters.Since)
	assert.Equal(t, "", config.Filters.Until)

	assert.Equal(t, "COMMENT", config.Review.Event)

	assert.Equal(t, 30, config.API.Timeout)
	assert.Equal(t, 3, config.API.RetryCount)
	assert.Equal(t, 10, config.API.RateLimitBuffer)

	assert.Equal(t, true, config.Suggestions.ExpandByDefault)
	assert.Equal(t, 999, config.Suggestions.MaxOffset)

	assert.NotNil(t, config.Aliases)
	assert.Equal(t, "Code review complete", config.Templates.DefaultReviewBody)
	assert.Equal(t, "LGTM! Ready to merge", config.Templates.DefaultApprovalMessage)
}

func TestLoadConfigFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
		content  string
		wantErr  bool
	}{
		{
			name:     "valid YAML config",
			filename: "config.yaml",
			content: `
defaults:
  repository: "owner/repo"
  author: "testuser"
  pr: 123
behavior:
  dry_run: true
  verbose: true
display:
  format: "json"
  color: "always"
filters:
  status: "open"
  type: "review"
`,
			wantErr: false,
		},
		{
			name:     "valid JSON config",
			filename: "config.json",
			content: `{
  "defaults": {
    "repository": "owner/repo",
    "author": "testuser",
    "pr": 456
  },
  "behavior": {
    "dry_run": false,
    "verbose": true
  },
  "display": {
    "format": "quiet",
    "color": "never"
  }
}`,
			wantErr: false,
		},
		{
			name:     "invalid YAML",
			filename: "invalid.yaml",
			content:  "invalid: yaml: content: [",
			wantErr:  true,
		},
		{
			name:     "invalid JSON",
			filename: "invalid.json",
			content:  `{"invalid": json}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(configPath, []byte(tt.content), 0644)
			require.NoError(t, err)

			config := NewDefaultConfig()
			err = loadConfigFile(config, configPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Test that some values were loaded
				if tt.name == "valid YAML config" {
					assert.Equal(t, "owner/repo", config.Defaults.Repository)
					assert.Equal(t, "testuser", config.Defaults.Author)
					assert.Equal(t, 123, config.Defaults.PR)
					assert.Equal(t, true, config.Behavior.DryRun)
					assert.Equal(t, "json", config.Display.Format)
					assert.Equal(t, "open", config.Filters.Status)
					assert.Equal(t, "review", config.Filters.Type)
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid default config",
			config:  NewDefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid repository format",
			config: &Config{
				Defaults:    DefaultsConfig{Repository: "invalid-repo"},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "invalid repository format",
		},
		{
			name: "invalid display format",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "invalid", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "invalid display format",
		},
		{
			name: "invalid color setting",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "invalid"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "invalid color setting",
		},
		{
			name: "invalid status filter",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "invalid", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "invalid status filter",
		},
		{
			name: "invalid type filter",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "invalid"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "invalid type filter",
		},
		{
			name: "invalid review event",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "INVALID"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "invalid review event",
		},
		{
			name: "invalid timeout",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 0, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "API timeout must be positive",
		},
		{
			name: "invalid retry count",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: -1},
				Suggestions: SuggestionsConfig{MaxOffset: 999},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "retry count must be non-negative",
		},
		{
			name: "invalid max offset",
			config: &Config{
				Defaults:    DefaultsConfig{},
				Behavior:    BehaviorConfig{Validate: true},
				Display:     DisplayConfig{Format: "table", Color: "auto"},
				Filters:     FiltersConfig{Status: "all", Type: "all"},
				Review:      ReviewDefaultsConfig{Event: "COMMENT"},
				API:         APIConfig{Timeout: 30, RetryCount: 3},
				Suggestions: SuggestionsConfig{MaxOffset: 10000},
				Aliases:     make(map[string]string),
			},
			wantErr: true,
			errMsg:  "max offset must be between 1 and 9999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidRepoFormat(t *testing.T) {
	tests := []struct {
		name string
		repo string
		want bool
	}{
		{"valid repo", "owner/repo", true},
		{"valid with hyphens", "my-org/my-repo", true},
		{"valid with underscores", "my_org/my_repo", true},
		{"valid with dots", "my.org/my.repo", true},
		{"valid with numbers", "org123/repo456", true},
		{"missing slash", "owner-repo", false},
		{"multiple slashes", "owner/group/repo", false},
		{"empty owner", "/repo", false},
		{"empty repo", "owner/", false},
		{"special chars", "owner@/repo", false},
		{"spaces", "owner /repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidRepoFormat(tt.repo)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestApplyEnvironmentOverrides(t *testing.T) {
	// Save original environment
	originalEnv := map[string]string{
		"GH_COMMENT_REPO":    os.Getenv("GH_COMMENT_REPO"),
		"GH_COMMENT_AUTHOR":  os.Getenv("GH_COMMENT_AUTHOR"),
		"GH_COMMENT_DRY_RUN": os.Getenv("GH_COMMENT_DRY_RUN"),
		"GH_COMMENT_VERBOSE": os.Getenv("GH_COMMENT_VERBOSE"),
		"GH_COMMENT_FORMAT":  os.Getenv("GH_COMMENT_FORMAT"),
		"GH_COMMENT_COLOR":   os.Getenv("GH_COMMENT_COLOR"),
	}

	// Restore environment after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Set test environment variables
	os.Setenv("GH_COMMENT_REPO", "test/repo")
	os.Setenv("GH_COMMENT_AUTHOR", "testuser")
	os.Setenv("GH_COMMENT_DRY_RUN", "true")
	os.Setenv("GH_COMMENT_VERBOSE", "yes")
	os.Setenv("GH_COMMENT_FORMAT", "json")
	os.Setenv("GH_COMMENT_COLOR", "always")

	config := NewDefaultConfig()
	applyEnvironmentOverrides(config)

	assert.Equal(t, "test/repo", config.Defaults.Repository)
	assert.Equal(t, "testuser", config.Defaults.Author)
	assert.Equal(t, true, config.Behavior.DryRun)
	assert.Equal(t, true, config.Behavior.Verbose)
	assert.Equal(t, "json", config.Display.Format)
	assert.Equal(t, "always", config.Display.Color)
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"yes", true},
		{"Yes", true},
		{"YES", true},
		{"1", true},
		{"on", true},
		{"On", true},
		{"ON", true},

		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"no", false},
		{"No", false},
		{"NO", false},
		{"0", false},
		{"off", false},
		{"Off", false},
		{"OFF", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseBool(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfigMarshaling(t *testing.T) {
	config := NewDefaultConfig()
	config.Defaults.Repository = "test/repo"
	config.Defaults.Author = "testuser"

	t.Run("YAML marshaling", func(t *testing.T) {
		data, err := yaml.Marshal(config)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "repository: test/repo")
		assert.Contains(t, string(data), "author: testuser")

		// Test unmarshaling
		var unmarshaled Config
		err = yaml.Unmarshal(data, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, config.Defaults.Repository, unmarshaled.Defaults.Repository)
		assert.Equal(t, config.Defaults.Author, unmarshaled.Defaults.Author)
	})

	t.Run("JSON marshaling", func(t *testing.T) {
		data, err := json.Marshal(config)
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"repository":"test/repo"`)
		assert.Contains(t, string(data), `"author":"testuser"`)

		// Test unmarshaling
		var unmarshaled Config
		err = json.Unmarshal(data, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, config.Defaults.Repository, unmarshaled.Defaults.Repository)
		assert.Equal(t, config.Defaults.Author, unmarshaled.Defaults.Author)
	})
}
