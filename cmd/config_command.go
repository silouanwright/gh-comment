package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configFormat string
	globalFlag   bool
	showSource   bool
	showEffective bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration files",
	Long: heredoc.Doc(`
		Manage gh-comment configuration files and settings.

		Configuration files are searched in this order:
		1. --config flag specified file
		2. .gh-comment.yaml in current directory
		3. .gh-comment.yaml in repository root
		4. ~/.config/gh-comment/config.yaml (user config)
		5. ~/.gh-comment.yaml (legacy user config)
	`),
	Example: heredoc.Doc(`
		# Generate a new config file
		$ gh comment config init
		$ gh comment config init --global --format json

		# Show current configuration
		$ gh comment config show
		$ gh comment config show --source
		$ gh comment config show --effective

		# Validate configuration
		$ gh comment config validate
		$ gh comment config validate ~/.gh-comment.yaml
	`),
}

// configInitCmd initializes a new config file
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a default configuration file",
	Long: heredoc.Doc(`
		Generate a default configuration file with common settings.

		By default, creates .gh-comment.yaml in the current directory.
		Use --global to create a user-wide config file.
	`),
	Example: heredoc.Doc(`
		# Create local config file
		$ gh comment config init

		# Create global user config file
		$ gh comment config init --global

		# Create config in JSON format
		$ gh comment config init --format json
	`),
	RunE: runConfigInit,
}

// configShowCmd shows current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: heredoc.Doc(`
		Show the current configuration, including merged values from all sources.

		Use --source to show which file each setting comes from.
		Use --effective to show only the final merged configuration.
	`),
	Example: heredoc.Doc(`
		# Show current config
		$ gh comment config show

		# Show config with source information
		$ gh comment config show --source

		# Show only effective (merged) config
		$ gh comment config show --effective
	`),
	RunE: runConfigShow,
}

// configValidateCmd validates configuration files
var configValidateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate configuration file",
	Long: heredoc.Doc(`
		Validate a configuration file for syntax and semantic errors.

		If no file is specified, validates the current effective configuration.
	`),
	Example: heredoc.Doc(`
		# Validate current config
		$ gh comment config validate

		# Validate specific file
		$ gh comment config validate ~/.gh-comment.yaml
		$ gh comment config validate .gh-comment.json
	`),
	Args: cobra.MaximumNArgs(1),
	RunE: runConfigValidate,
}

func init() {
	// Add subcommands
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)

	// Init command flags
	configInitCmd.Flags().StringVar(&configFormat, "format", "yaml", "Config file format (yaml or json)")
	configInitCmd.Flags().BoolVar(&globalFlag, "global", false, "Create global user config file")

	// Show command flags
	configShowCmd.Flags().BoolVar(&showSource, "source", false, "Show which file each setting comes from")
	configShowCmd.Flags().BoolVar(&showEffective, "effective", false, "Show only effective (merged) configuration")

	// Add to root command  
	rootCmd.AddCommand(configCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	config := NewDefaultConfig()

	// Determine output file path
	var outputPath string
	if globalFlag {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		
		configDir := filepath.Join(homeDir, ".config", "gh-comment")
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		
		ext := "yaml"
		if configFormat == "json" {
			ext = "json"
		}
		outputPath = filepath.Join(configDir, "config."+ext)
	} else {
		ext := "yaml"
		if configFormat == "json" {
			ext = "json"
		}
		outputPath = ".gh-comment." + ext
	}

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("config file already exists: %s", outputPath)
	}

	// Generate config content
	var content []byte
	var err error
	
	if configFormat == "json" {
		content, err = json.MarshalIndent(config, "", "  ")
	} else {
		content, err = yaml.Marshal(config)
	}
	
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	err = os.WriteFile(outputPath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✅ Created config file: %s\n", outputPath)
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	if showEffective {
		// Show only the effective merged config
		config := GetConfig()
		content, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		fmt.Print(string(content))
		return nil
	}

	if showSource {
		// Show config with source information
		fmt.Println("# Configuration sources (in priority order):")
		fmt.Println("# 1. Command-line flags (highest priority)")
		fmt.Println("# 2. Environment variables")
		fmt.Println("# 3. Configuration files")
		fmt.Println("# 4. Built-in defaults (lowest priority)")
		fmt.Println()
		
		// Find and show config file sources
		configPath, _ := findConfigFile()
		if configPath != "" {
			fmt.Printf("# Active config file: %s\n", configPath)
		} else {
			fmt.Println("# No config file found - using defaults")
		}
		fmt.Println()
	}

	// Show current effective config
	config := GetConfig()
	content, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	fmt.Print(string(content))
	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	var configPath string
	if len(args) > 0 {
		configPath = args[0]
	}

	// Load and validate config
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("❌ Configuration validation failed: %w", err)
	}

	// Show which file was validated
	if configPath != "" {
		fmt.Printf("✅ Configuration file is valid: %s\n", configPath)
	} else {
		// Show which files were found and merged
		foundPath, _ := findConfigFile()
		if foundPath != "" {
			fmt.Printf("✅ Configuration is valid: %s\n", foundPath)
		} else {
			fmt.Println("✅ Using default configuration (no config file found)")
		}
	}

	// Show any warnings
	showConfigWarnings(config)
	
	return nil
}

func showConfigWarnings(config *Config) {
	warnings := []string{}

	// Check for potentially problematic values
	if config.API.Timeout > 300 {
		warnings = append(warnings, fmt.Sprintf("Large API timeout: %d seconds", config.API.Timeout))
	}
	if config.API.RetryCount > 10 {
		warnings = append(warnings, fmt.Sprintf("High retry count: %d", config.API.RetryCount))
	}
	if config.Suggestions.MaxOffset > 100 {
		warnings = append(warnings, fmt.Sprintf("Large max offset: %d", config.Suggestions.MaxOffset))
	}

	// Check repository format if set
	if config.Defaults.Repository != "" && !isValidRepoFormat(config.Defaults.Repository) {
		warnings = append(warnings, fmt.Sprintf("Invalid repository format: %s", config.Defaults.Repository))
	}

	if len(warnings) > 0 {
		fmt.Println("\n⚠️  Warnings:")
		for _, warning := range warnings {
			fmt.Printf("  • %s\n", warning)
		}
	}
}