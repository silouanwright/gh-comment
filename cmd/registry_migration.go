package cmd

import (
	"github.com/spf13/cobra"
)

// MigrateToRegistry helps convert existing init()-based commands to use the registry pattern
// This file provides helper functions to ease the transition from init() to explicit registration

// CommandCategories defines standard command categories for consistent organization
const (
	CategoryCore    = "core"      // Primary commenting operations (add, review, list)
	CategoryManage  = "manage"    // Comment management (edit, resolve, react)
	CategoryAdmin   = "admin"     // Administrative tasks (config, export)
	CategoryUtility = "utility"   // Helper commands (lines, prompts)
	CategoryTest    = "test"      // Testing and integration commands
)

// CommandPriorities defines standard priority levels for help display ordering
const (
	PriorityHigh   = 1  // Most important commands (add, review, list)
	PriorityMedium = 5  // Regular commands
	PriorityLow    = 10 // Less frequently used commands
)

// RegisterCoreCommands registers all core commenting commands
func RegisterCoreCommands(registry CommandRegistry) error {
	coreCommands := []CommandInfo{
		{
			Name:        "add",
			Category:    CategoryCore,
			Description: "Add issue comments for general PR discussion",
			Priority:    PriorityHigh,
			Builder:     buildAddCommand,
		},
		{
			Name:        "review",
			Category:    CategoryCore,
			Description: "Create code reviews with line-specific comments",
			Priority:    PriorityHigh,
			Builder:     buildReviewCommand,
		},
		{
			Name:        "list",
			Category:    CategoryCore,
			Description: "List and filter PR comments with advanced options",
			Priority:    PriorityHigh,
			Builder:     buildListCommand,
		},
		{
			Name:        "batch",
			Category:    CategoryCore,
			Description: "Process multiple comments from YAML configuration",
			Priority:    PriorityMedium,
			Builder:     buildBatchCommand,
		},
	}
	
	for _, info := range coreCommands {
		if err := registry.Register(info); err != nil {
			return err
		}
	}
	
	return nil
}

// RegisterManagementCommands registers comment management commands
func RegisterManagementCommands(registry CommandRegistry) error {
	managementCommands := []CommandInfo{
		{
			Name:        "edit",
			Category:    CategoryManage,
			Description: "Edit existing comments on a PR",
			Priority:    PriorityMedium,
			Builder:     buildEditCommand,
		},
		{
			Name:        "resolve",
			Category:    CategoryManage,
			Description: "Mark conversation threads as resolved",
			Priority:    PriorityMedium,
			Builder:     buildResolveCommand,
		},
		{
			Name:        "react",
			Category:    CategoryManage,
			Description: "Add or remove emoji reactions to comments",
			Priority:    PriorityMedium,
			Builder:     buildReactCommand,
		},
		{
			Name:        "review-reply",
			Category:    CategoryManage,
			Description: "Reply to review comments with text messages",
			Priority:    PriorityMedium,
			Builder:     buildReviewReplyCommand,
		},
		{
			Name:        "close-pending-review",
			Category:    CategoryManage,
			Description: "Submit pending review comments",
			Priority:    PriorityLow,
			Builder:     buildClosePendingReviewCommand,
		},
	}
	
	for _, info := range managementCommands {
		if err := registry.Register(info); err != nil {
			return err
		}
	}
	
	return nil
}

// RegisterAdminCommands registers administrative commands
func RegisterAdminCommands(registry CommandRegistry) error {
	adminCommands := []CommandInfo{
		{
			Name:        "config",
			Category:    CategoryAdmin,
			Description: "Manage configuration files and settings",
			Priority:    PriorityLow,
			Builder:     buildConfigCommand,
		},
		{
			Name:        "export",
			Category:    CategoryAdmin,
			Description: "Export PR comments to various formats",
			Priority:    PriorityLow,
			Builder:     buildExportCommand,
		},
	}
	
	for _, info := range adminCommands {
		if err := registry.Register(info); err != nil {
			return err
		}
	}
	
	return nil
}

// RegisterUtilityCommands registers utility and helper commands
func RegisterUtilityCommands(registry CommandRegistry) error {
	utilityCommands := []CommandInfo{
		{
			Name:        "lines",
			Category:    CategoryUtility,
			Description: "Show commentable lines in PR files",
			Priority:    PriorityMedium,
			Builder:     buildLinesCommand,
		},
		{
			Name:        "prompts",
			Category:    CategoryUtility,
			Description: "Get AI-powered code review prompts and best practices",
			Priority:    PriorityMedium,
			Builder:     buildPromptsCommand,
		},
	}
	
	for _, info := range utilityCommands {
		if err := registry.Register(info); err != nil {
			return err
		}
	}
	
	return nil
}

// RegisterTestCommands registers testing and integration commands
func RegisterTestCommands(registry CommandRegistry) error {
	testCommands := []CommandInfo{
		{
			Name:        "test-integration",
			Category:    CategoryTest,
			Description: "Run integration tests with live GitHub API",
			Priority:    PriorityLow,
			Builder:     buildTestIntegrationCommand,
		},
	}
	
	for _, info := range testCommands {
		if err := registry.Register(info); err != nil {
			return err
		}
	}
	
	return nil
}

// RegisterAllCommands registers all available commands with the registry
func RegisterAllCommands(registry CommandRegistry) error {
	registrationFunctions := []func(CommandRegistry) error{
		RegisterCoreCommands,
		RegisterManagementCommands,
		RegisterAdminCommands,
		RegisterUtilityCommands,
		RegisterTestCommands,
	}
	
	for _, registerFunc := range registrationFunctions {
		if err := registerFunc(registry); err != nil {
			return err
		}
	}
	
	return nil
}

// Command builders - these replace the init() functions
// Each builder creates the command without automatically registering it

func buildAddCommand() *cobra.Command {
	// This would contain the actual command creation logic from add.go
	// For now, returning a placeholder to satisfy the interface
	return &cobra.Command{
		Use:   "add",
		Short: "Add issue comments for general PR discussion",
		RunE:  runAdd, // The actual run function would remain the same
	}
}

func buildReviewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "Create code reviews with line-specific comments",
		RunE:  runReview,
	}
}

func buildListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List and filter PR comments with advanced options",
		RunE:  runList,
	}
}

func buildBatchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "batch",
		Short: "Process multiple comments from YAML configuration",
		RunE:  runBatch,
	}
}

func buildEditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit existing comments on a PR",
		RunE:  runEdit,
	}
}

func buildResolveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "resolve",
		Short: "Mark conversation threads as resolved",
		RunE:  runResolve,
	}
}

func buildReactCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "react",
		Short: "Add or remove emoji reactions to comments",
		RunE:  runReact,
	}
}

func buildReviewReplyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "review-reply",
		Short: "Reply to review comments with text messages",
		RunE:  runReviewReply,
	}
}

func buildClosePendingReviewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "close-pending-review",
		Short: "Submit pending review comments",
		RunE:  runClosePendingReview,
	}
}

func buildConfigCommand() *cobra.Command {
	// Config command is a parent command with subcommands, no RunE needed
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration files and settings",
	}
	
	// Add subcommands (in a real implementation, these would be built from the original)
	// For now, returning the parent command structure
	return cmd
}

func buildExportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Export PR comments to various formats",
		RunE:  runExport,
	}
}

func buildLinesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "lines",
		Short: "Show commentable lines in PR files",
		RunE:  runLines,
	}
}

func buildPromptsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "prompts",
		Short: "Get AI-powered code review prompts and best practices",
		RunE:  runPrompts,
	}
}

func buildTestIntegrationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "test-integration",
		Short: "Run integration tests with live GitHub API",
		RunE: func(cmd *cobra.Command, args []string) error {
			// In a real implementation, this would call the actual runTestIntegration function
			// For demonstration purposes, using a placeholder
			return nil
		},
	}
}

// Legacy compatibility functions for gradual migration

// InitializeCommandRegistry sets up the registry with all commands
// This can be called from root.go instead of relying on init() functions
func InitializeCommandRegistry() error {
	registry := GetRegistry()
	return RegisterAllCommands(registry)
}

// BuildAndRegisterCommands initializes the registry and builds all commands
// This is the main function to replace the init() pattern
func BuildAndRegisterCommands(rootCmd *cobra.Command) error {
	if err := InitializeCommandRegistry(); err != nil {
		return err
	}
	return BuildAllCommands(rootCmd)
}