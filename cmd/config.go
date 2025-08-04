package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration structure
type Config struct {
	Defaults    DefaultsConfig       `yaml:"defaults" json:"defaults"`
	Behavior    BehaviorConfig       `yaml:"behavior" json:"behavior"`
	Display     DisplayConfig        `yaml:"display" json:"display"`
	Filters     FiltersConfig        `yaml:"filters" json:"filters"`
	Review      ReviewDefaultsConfig `yaml:"review" json:"review"`
	API         APIConfig            `yaml:"api" json:"api"`
	Suggestions SuggestionsConfig    `yaml:"suggestions" json:"suggestions"`
	Aliases     map[string]string    `yaml:"aliases" json:"aliases"`
	Templates   TemplatesConfig      `yaml:"templates" json:"templates"`
}

// DefaultsConfig holds default values for common flags
type DefaultsConfig struct {
	Repository string `yaml:"repository" json:"repository"`
	Author     string `yaml:"author" json:"author"`
	PR         int    `yaml:"pr" json:"pr"`
}

// BehaviorConfig holds default behavior settings
type BehaviorConfig struct {
	DryRun              bool `yaml:"dry_run" json:"dry_run"`
	Verbose             bool `yaml:"verbose" json:"verbose"`
	Validate            bool `yaml:"validate" json:"validate"`
	NoExpandSuggestions bool `yaml:"no_expand_suggestions" json:"no_expand_suggestions"`
}

// DisplayConfig holds display preference settings
type DisplayConfig struct {
	Format string `yaml:"format" json:"format"`
	Color  string `yaml:"color" json:"color"`
	Quiet  bool   `yaml:"quiet" json:"quiet"`
}

// FiltersConfig holds default filter settings
type FiltersConfig struct {
	Status string `yaml:"status" json:"status"`
	Type   string `yaml:"type" json:"type"`
	Since  string `yaml:"since" json:"since"`
	Until  string `yaml:"until" json:"until"`
}

// ReviewDefaultsConfig holds review-specific settings
type ReviewDefaultsConfig struct {
	Event string `yaml:"event" json:"event"`
}

// APIConfig holds API-related settings
type APIConfig struct {
	Timeout         int `yaml:"timeout" json:"timeout"`
	RetryCount      int `yaml:"retry_count" json:"retry_count"`
	RateLimitBuffer int `yaml:"rate_limit_buffer" json:"rate_limit_buffer"`
}

// SuggestionsConfig holds suggestion syntax settings
type SuggestionsConfig struct {
	ExpandByDefault bool `yaml:"expand_by_default" json:"expand_by_default"`
	MaxOffset       int  `yaml:"max_offset" json:"max_offset"`
}

// TemplatesConfig holds template settings
type TemplatesConfig struct {
	DefaultReviewBody      string `yaml:"default_review_body" json:"default_review_body"`
	DefaultApprovalMessage string `yaml:"default_approval_message" json:"default_approval_message"`
}

// NewDefaultConfig returns a configuration with sensible defaults
func NewDefaultConfig() *Config {
	return &Config{
		Defaults: DefaultsConfig{
			Repository: "",
			Author:     "",
			PR:         0,
		},
		Behavior: BehaviorConfig{
			DryRun:              false,
			Verbose:             false,
			Validate:            true,
			NoExpandSuggestions: false,
		},
		Display: DisplayConfig{
			Format: "table",
			Color:  "auto",
			Quiet:  false,
		},
		Filters: FiltersConfig{
			Status: "all",
			Type:   "all",
			Since:  "",
			Until:  "",
		},
		Review: ReviewDefaultsConfig{
			Event: "COMMENT",
		},
		API: APIConfig{
			Timeout:         30,
			RetryCount:      3,
			RateLimitBuffer: 10,
		},
		Suggestions: SuggestionsConfig{
			ExpandByDefault: true,
			MaxOffset:       999,
		},
		Aliases: make(map[string]string),
		Templates: TemplatesConfig{
			DefaultReviewBody:      "Code review complete",
			DefaultApprovalMessage: "LGTM! Ready to merge",
		},
	}
}

// LoadConfig loads configuration from various sources in priority order
func LoadConfig(configPath string) (*Config, error) {
	config := NewDefaultConfig()

	// Find config file if not explicitly provided
	if configPath == "" {
		var err error
		configPath, err = findConfigFile()
		if err != nil {
			// No config file found, use defaults
			return config, nil
		}
	}

	// Load and parse config file
	if configPath != "" {
		err := loadConfigFile(config, configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
		}
	}

	// Apply environment variable overrides
	applyEnvironmentOverrides(config)

	// Validate configuration
	err := validateConfig(config)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// findConfigFile searches for config files in priority order
func findConfigFile() (string, error) {
	// Search locations in priority order
	searchPaths := []string{
		".gh-comment.yaml",
		".gh-comment.yml",
		".gh-comment.json",
	}

	// Add repository root if in git repo
	if gitRoot := findGitRoot(); gitRoot != "" {
		for _, name := range []string{".gh-comment.yaml", ".gh-comment.yml", ".gh-comment.json"} {
			searchPaths = append(searchPaths, filepath.Join(gitRoot, name))
		}
	}

	// Add user config locations
	if homeDir, err := os.UserHomeDir(); err == nil {
		searchPaths = append(searchPaths,
			filepath.Join(homeDir, ".config", "gh-comment", "config.yaml"),
			filepath.Join(homeDir, ".config", "gh-comment", "config.yml"),
			filepath.Join(homeDir, ".config", "gh-comment", "config.json"),
			filepath.Join(homeDir, ".gh-comment.yaml"),
			filepath.Join(homeDir, ".gh-comment.yml"),
			filepath.Join(homeDir, ".gh-comment.json"),
		)
	}

	// Find first existing file
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no config file found")
}

// findGitRoot finds the root of the current git repository
func findGitRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		gitDir := filepath.Join(currentDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return currentDir
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached filesystem root
			break
		}
		currentDir = parentDir
	}

	return ""
}

// loadConfigFile loads and parses a configuration file
func loadConfigFile(config *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Determine format by extension
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, config)
	case ".json":
		err = json.Unmarshal(data, config)
	default:
		// Try YAML first, then JSON
		err = yaml.Unmarshal(data, config)
		if err != nil {
			err = json.Unmarshal(data, config)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// applyEnvironmentOverrides applies environment variable overrides
func applyEnvironmentOverrides(config *Config) {
	if val := os.Getenv("GH_COMMENT_REPO"); val != "" {
		config.Defaults.Repository = val
	}
	if val := os.Getenv("GH_COMMENT_AUTHOR"); val != "" {
		config.Defaults.Author = val
	}
	if val := os.Getenv("GH_COMMENT_DRY_RUN"); val != "" {
		config.Behavior.DryRun = parseBool(val)
	}
	if val := os.Getenv("GH_COMMENT_VERBOSE"); val != "" {
		config.Behavior.Verbose = parseBool(val)
	}
	if val := os.Getenv("GH_COMMENT_FORMAT"); val != "" {
		config.Display.Format = val
	}
	if val := os.Getenv("GH_COMMENT_COLOR"); val != "" {
		config.Display.Color = val
	}
}

// parseBool parses a string as a boolean (true/false, yes/no, 1/0)
func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "yes" || s == "1" || s == "on"
}

// validateConfig validates the configuration values
func validateConfig(config *Config) error {
	// Validate repository format
	if config.Defaults.Repository != "" {
		if !isValidRepoFormat(config.Defaults.Repository) {
			return fmt.Errorf("invalid repository format: %s (must be owner/repo)", config.Defaults.Repository)
		}
	}

	// Validate enum values
	validFormats := map[string]bool{"table": true, "json": true, "quiet": true}
	if !validFormats[config.Display.Format] {
		return fmt.Errorf("invalid display format: %s (must be table, json, or quiet)", config.Display.Format)
	}

	validColors := map[string]bool{"auto": true, "always": true, "never": true}
	if !validColors[config.Display.Color] {
		return fmt.Errorf("invalid color setting: %s (must be auto, always, or never)", config.Display.Color)
	}

	validStatuses := map[string]bool{"all": true, "open": true, "resolved": true}
	if !validStatuses[config.Filters.Status] {
		return fmt.Errorf("invalid status filter: %s (must be all, open, or resolved)", config.Filters.Status)
	}

	validTypes := map[string]bool{"all": true, "issue": true, "review": true}
	if !validTypes[config.Filters.Type] {
		return fmt.Errorf("invalid type filter: %s (must be all, issue, or review)", config.Filters.Type)
	}

	validEvents := map[string]bool{"APPROVE": true, "REQUEST_CHANGES": true, "COMMENT": true}
	if !validEvents[config.Review.Event] {
		return fmt.Errorf("invalid review event: %s (must be APPROVE, REQUEST_CHANGES, or COMMENT)", config.Review.Event)
	}

	// Validate numeric ranges
	if config.API.Timeout <= 0 {
		return fmt.Errorf("API timeout must be positive: %d", config.API.Timeout)
	}
	if config.API.RetryCount < 0 {
		return fmt.Errorf("retry count must be non-negative: %d", config.API.RetryCount)
	}
	if config.Suggestions.MaxOffset <= 0 || config.Suggestions.MaxOffset > 9999 {
		return fmt.Errorf("max offset must be between 1 and 9999: %d", config.Suggestions.MaxOffset)
	}

	return nil
}

// isValidRepoFormat validates repository format (owner/repo)
func isValidRepoFormat(repo string) bool {
	// Must be owner/repo format
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
	return pattern.MatchString(repo)
}

// Global config instance
var globalConfig *Config

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	if globalConfig == nil {
		// Load default config if none loaded
		globalConfig = NewDefaultConfig()
	}
	return globalConfig
}

// LoadGlobalConfig loads the global configuration
func LoadGlobalConfig(configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}
	globalConfig = config
	return nil
}
