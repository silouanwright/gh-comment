package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// CommandBuilder creates and configures a cobra command
type CommandBuilder func() *cobra.Command

// CommandInfo contains metadata about a registered command
type CommandInfo struct {
	Name        string
	Category    string
	Description string
	Priority    int // Lower number = higher priority for help display
	Builder     CommandBuilder
	Command     *cobra.Command // Cached command instance
}

// CommandRegistry manages command registration and discovery
type CommandRegistry interface {
	// Register adds a command to the registry
	Register(info CommandInfo) error

	// GetCommand returns a specific command by name
	GetCommand(name string) (*cobra.Command, error)

	// GetAllCommands returns all registered commands
	GetAllCommands() map[string]*cobra.Command

	// GetCommandsByCategory returns commands grouped by category
	GetCommandsByCategory() map[string][]*CommandInfo

	// ListCommands returns a sorted list of command names
	ListCommands() []string

	// BuildAll builds all registered commands and adds them to the root command
	BuildAll(rootCmd *cobra.Command) error

	// GetCommandInfo returns metadata for a command
	GetCommandInfo(name string) (*CommandInfo, error)

	// GetRegisteredCount returns the number of registered commands
	GetRegisteredCount() int
}

// DefaultCommandRegistry implements CommandRegistry
type DefaultCommandRegistry struct {
	commands map[string]*CommandInfo
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() CommandRegistry {
	return &DefaultCommandRegistry{
		commands: make(map[string]*CommandInfo),
	}
}

// Register adds a command to the registry
func (r *DefaultCommandRegistry) Register(info CommandInfo) error {
	if info.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	if info.Builder == nil {
		return fmt.Errorf("command builder cannot be nil for command %s", info.Name)
	}

	if _, exists := r.commands[info.Name]; exists {
		return fmt.Errorf("command %s is already registered", info.Name)
	}

	// Set defaults
	if info.Category == "" {
		info.Category = "general"
	}
	if info.Description == "" {
		info.Description = "No description available"
	}

	r.commands[info.Name] = &info
	return nil
}

// GetCommand returns a specific command by name
func (r *DefaultCommandRegistry) GetCommand(name string) (*cobra.Command, error) {
	info, exists := r.commands[name]
	if !exists {
		return nil, fmt.Errorf("command %s not found", name)
	}

	// Build command if not already cached
	if info.Command == nil {
		info.Command = info.Builder()
		if info.Command == nil {
			return nil, fmt.Errorf("command builder for %s returned nil", name)
		}
	}

	return info.Command, nil
}

// GetAllCommands returns all registered commands
func (r *DefaultCommandRegistry) GetAllCommands() map[string]*cobra.Command {
	result := make(map[string]*cobra.Command)

	for name, info := range r.commands {
		if info.Command == nil {
			info.Command = info.Builder()
		}
		if info.Command != nil {
			result[name] = info.Command
		}
	}

	return result
}

// GetCommandsByCategory returns commands grouped by category
func (r *DefaultCommandRegistry) GetCommandsByCategory() map[string][]*CommandInfo {
	result := make(map[string][]*CommandInfo)

	for _, info := range r.commands {
		result[info.Category] = append(result[info.Category], info)
	}

	// Sort each category by priority, then by name
	for category := range result {
		sort.Slice(result[category], func(i, j int) bool {
			if result[category][i].Priority != result[category][j].Priority {
				return result[category][i].Priority < result[category][j].Priority
			}
			return result[category][i].Name < result[category][j].Name
		})
	}

	return result
}

// ListCommands returns a sorted list of command names
func (r *DefaultCommandRegistry) ListCommands() []string {
	var names []string
	for name := range r.commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// BuildAll builds all registered commands and adds them to the root command
func (r *DefaultCommandRegistry) BuildAll(rootCmd *cobra.Command) error {
	var errors []string

	for name, info := range r.commands {
		if info.Command == nil {
			info.Command = info.Builder()
		}

		if info.Command == nil {
			errors = append(errors, fmt.Sprintf("command %s builder returned nil", name))
			continue
		}

		rootCmd.AddCommand(info.Command)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to build commands: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetCommandInfo returns metadata for a command
func (r *DefaultCommandRegistry) GetCommandInfo(name string) (*CommandInfo, error) {
	info, exists := r.commands[name]
	if !exists {
		return nil, fmt.Errorf("command %s not found", name)
	}

	// Return a copy to prevent modification
	infoCopy := *info
	return &infoCopy, nil
}

// GetRegisteredCount returns the number of registered commands
func (r *DefaultCommandRegistry) GetRegisteredCount() int {
	return len(r.commands)
}

// Global registry instance
var defaultRegistry CommandRegistry

// GetRegistry returns the global command registry
func GetRegistry() CommandRegistry {
	if defaultRegistry == nil {
		defaultRegistry = NewCommandRegistry()
	}
	return defaultRegistry
}

// SetRegistry allows setting a custom registry (useful for testing)
func SetRegistry(registry CommandRegistry) {
	defaultRegistry = registry
}

// RegisterCommand is a convenience function to register a command with the global registry
func RegisterCommand(info CommandInfo) error {
	return GetRegistry().Register(info)
}

// BuildAllCommands builds all registered commands and adds them to the root command
func BuildAllCommands(rootCmd *cobra.Command) error {
	return GetRegistry().BuildAll(rootCmd)
}
