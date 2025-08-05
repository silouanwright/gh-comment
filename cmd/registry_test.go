package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCommandRegistry(t *testing.T) {
	registry := NewCommandRegistry()
	if registry == nil {
		t.Error("NewCommandRegistry() returned nil")
	}

	count := registry.GetRegisteredCount()
	if count != 0 {
		t.Errorf("NewCommandRegistry() should start with 0 commands, got %d", count)
	}
}

func TestCommandRegistryRegister(t *testing.T) {
	tests := []struct {
		name        string
		info        CommandInfo
		wantErr     bool
		errContains string
	}{
		{
			name: "valid command registration",
			info: CommandInfo{
				Name:        "test-cmd",
				Category:    "testing",
				Description: "Test command",
				Priority:    1,
				Builder: func() *cobra.Command {
					return &cobra.Command{Use: "test-cmd"}
				},
			},
			wantErr: false,
		},
		{
			name: "empty name should fail",
			info: CommandInfo{
				Name: "",
				Builder: func() *cobra.Command {
					return &cobra.Command{Use: "test"}
				},
			},
			wantErr:     true,
			errContains: "command name cannot be empty",
		},
		{
			name: "nil builder should fail",
			info: CommandInfo{
				Name:    "test-cmd",
				Builder: nil,
			},
			wantErr:     true,
			errContains: "command builder cannot be nil",
		},
		{
			name: "defaults should be applied",
			info: CommandInfo{
				Name: "test-cmd-defaults",
				Builder: func() *cobra.Command {
					return &cobra.Command{Use: "test-cmd-defaults"}
				},
				// Category and Description omitted to test defaults
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewCommandRegistry()
			err := registry.Register(tt.info)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Register() expected error but got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Register() error = %v, want to contain %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Register() unexpected error = %v", err)
				} else {
					// Verify registration worked
					if registry.GetRegisteredCount() != 1 {
						t.Errorf("Register() should have 1 command, got %d", registry.GetRegisteredCount())
					}

					// Test defaults were applied
					if tt.info.Name == "test-cmd-defaults" {
						info, err := registry.GetCommandInfo(tt.info.Name)
						if err != nil {
							t.Errorf("GetCommandInfo() error = %v", err)
						} else {
							if info.Category != "general" {
								t.Errorf("Default category should be 'general', got '%s'", info.Category)
							}
							if info.Description != "No description available" {
								t.Errorf("Default description should be 'No description available', got '%s'", info.Description)
							}
						}
					}
				}
			}
		})
	}
}

func TestCommandRegistryDuplicateRegistration(t *testing.T) {
	registry := NewCommandRegistry()

	info := CommandInfo{
		Name: "duplicate-test",
		Builder: func() *cobra.Command {
			return &cobra.Command{Use: "duplicate-test"}
		},
	}

	// First registration should succeed
	err := registry.Register(info)
	if err != nil {
		t.Errorf("First Register() unexpected error = %v", err)
	}

	// Second registration should fail
	err = registry.Register(info)
	if err == nil {
		t.Error("Second Register() expected error but got nil")
	} else if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("Register() error = %v, want to contain 'already registered'", err)
	}
}

func TestCommandRegistryGetCommand(t *testing.T) {
	registry := NewCommandRegistry()

	info := CommandInfo{
		Name: "get-test",
		Builder: func() *cobra.Command {
			return &cobra.Command{
				Use:   "get-test",
				Short: "Test command for GetCommand",
			}
		},
	}

	// Register the command
	err := registry.Register(info)
	if err != nil {
		t.Errorf("Register() unexpected error = %v", err)
	}

	// Get existing command
	cmd, err := registry.GetCommand("get-test")
	if err != nil {
		t.Errorf("GetCommand() unexpected error = %v", err)
	}
	if cmd == nil {
		t.Error("GetCommand() returned nil command")
	} else if cmd.Use != "get-test" {
		t.Errorf("GetCommand() returned command with wrong Use: got %s, want get-test", cmd.Use)
	}

	// Test command caching - second call should return same instance
	cmd2, err := registry.GetCommand("get-test")
	if err != nil {
		t.Errorf("Second GetCommand() unexpected error = %v", err)
	}
	if cmd != cmd2 {
		t.Error("GetCommand() should return cached instance on second call")
	}

	// Get non-existent command
	_, err = registry.GetCommand("non-existent")
	if err == nil {
		t.Error("GetCommand() expected error for non-existent command")
	} else if !strings.Contains(err.Error(), "not found") {
		t.Errorf("GetCommand() error = %v, want to contain 'not found'", err)
	}
}

func TestCommandRegistryGetAllCommands(t *testing.T) {
	registry := NewCommandRegistry()

	// Register multiple commands
	commands := []string{"cmd1", "cmd2", "cmd3"}
	for _, name := range commands {
		info := CommandInfo{
			Name: name,
			Builder: func(cmdName string) func() *cobra.Command {
				return func() *cobra.Command {
					return &cobra.Command{Use: cmdName}
				}
			}(name), // Capture name in closure
		}
		err := registry.Register(info)
		if err != nil {
			t.Errorf("Register(%s) unexpected error = %v", name, err)
		}
	}

	allCommands := registry.GetAllCommands()
	if len(allCommands) != len(commands) {
		t.Errorf("GetAllCommands() returned %d commands, want %d", len(allCommands), len(commands))
	}

	for _, name := range commands {
		if cmd, exists := allCommands[name]; !exists {
			t.Errorf("GetAllCommands() missing command %s", name)
		} else if cmd.Use != name {
			t.Errorf("GetAllCommands() command %s has wrong Use: got %s", name, cmd.Use)
		}
	}
}

func TestCommandRegistryGetCommandsByCategory(t *testing.T) {
	registry := NewCommandRegistry()

	// Register commands in different categories
	testCases := []struct {
		name     string
		category string
		priority int
	}{
		{"review", "core", 1},
		{"add", "core", 2},
		{"list", "core", 3},
		{"config", "admin", 1},
		{"export", "utility", 1},
		{"lines", "utility", 2},
	}

	for _, tc := range testCases {
		info := CommandInfo{
			Name:     tc.name,
			Category: tc.category,
			Priority: tc.priority,
			Builder: func(cmdName string) func() *cobra.Command {
				return func() *cobra.Command {
					return &cobra.Command{Use: cmdName}
				}
			}(tc.name),
		}
		err := registry.Register(info)
		if err != nil {
			t.Errorf("Register(%s) unexpected error = %v", tc.name, err)
		}
	}

	categories := registry.GetCommandsByCategory()

	// Check expected categories exist
	expectedCategories := []string{"core", "admin", "utility"}
	for _, category := range expectedCategories {
		if _, exists := categories[category]; !exists {
			t.Errorf("GetCommandsByCategory() missing category %s", category)
		}
	}

	// Check core category has correct commands in priority order
	coreCommands := categories["core"]
	if len(coreCommands) != 3 {
		t.Errorf("Core category should have 3 commands, got %d", len(coreCommands))
	} else {
		expectedOrder := []string{"review", "add", "list"} // Sorted by priority
		for i, expected := range expectedOrder {
			if coreCommands[i].Name != expected {
				t.Errorf("Core commands[%d] = %s, want %s", i, coreCommands[i].Name, expected)
			}
		}
	}
}

func TestCommandRegistryListCommands(t *testing.T) {
	registry := NewCommandRegistry()

	// Register commands in random order
	commands := []string{"zebra", "alpha", "beta", "gamma"}
	for _, name := range commands {
		info := CommandInfo{
			Name: name,
			Builder: func(cmdName string) func() *cobra.Command {
				return func() *cobra.Command {
					return &cobra.Command{Use: cmdName}
				}
			}(name),
		}
		err := registry.Register(info)
		if err != nil {
			t.Errorf("Register(%s) unexpected error = %v", name, err)
		}
	}

	commandList := registry.ListCommands()
	expectedOrder := []string{"alpha", "beta", "gamma", "zebra"} // Alphabetical order

	if len(commandList) != len(expectedOrder) {
		t.Errorf("ListCommands() returned %d commands, want %d", len(commandList), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(commandList) || commandList[i] != expected {
			t.Errorf("ListCommands()[%d] = %s, want %s", i, commandList[i], expected)
		}
	}
}

func TestCommandRegistryBuildAll(t *testing.T) {
	registry := NewCommandRegistry()
	rootCmd := &cobra.Command{Use: "root"}

	// Register test commands
	commands := []string{"build1", "build2", "build3"}
	for _, name := range commands {
		info := CommandInfo{
			Name: name,
			Builder: func(cmdName string) func() *cobra.Command {
				return func() *cobra.Command {
					return &cobra.Command{
						Use:   cmdName,
						Short: "Test command " + cmdName,
					}
				}
			}(name),
		}
		err := registry.Register(info)
		if err != nil {
			t.Errorf("Register(%s) unexpected error = %v", name, err)
		}
	}

	// Test BuildAll
	err := registry.BuildAll(rootCmd)
	if err != nil {
		t.Errorf("BuildAll() unexpected error = %v", err)
	}

	// Verify commands were added to root
	if len(rootCmd.Commands()) != len(commands) {
		t.Errorf("BuildAll() added %d commands to root, want %d", len(rootCmd.Commands()), len(commands))
	}

	// Verify each command was added correctly
	for _, expectedName := range commands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BuildAll() did not add command %s to root", expectedName)
		}
	}
}

func TestCommandRegistryBuildAllWithFailingBuilder(t *testing.T) {
	registry := NewCommandRegistry()
	rootCmd := &cobra.Command{Use: "root"}

	// Register a command with a failing builder
	info := CommandInfo{
		Name: "failing-cmd",
		Builder: func() *cobra.Command {
			return nil // This should cause BuildAll to fail
		},
	}
	err := registry.Register(info)
	if err != nil {
		t.Errorf("Register() unexpected error = %v", err)
	}

	// BuildAll should fail
	err = registry.BuildAll(rootCmd)
	if err == nil {
		t.Error("BuildAll() expected error for failing builder")
	} else if !strings.Contains(err.Error(), "builder returned nil") {
		t.Errorf("BuildAll() error = %v, want to contain 'builder returned nil'", err)
	}
}

func TestCommandRegistryGetCommandInfo(t *testing.T) {
	registry := NewCommandRegistry()

	originalInfo := CommandInfo{
		Name:        "info-test",
		Category:    "testing",
		Description: "Test for GetCommandInfo",
		Priority:    5,
		Builder: func() *cobra.Command {
			return &cobra.Command{Use: "info-test"}
		},
	}

	err := registry.Register(originalInfo)
	if err != nil {
		t.Errorf("Register() unexpected error = %v", err)
	}

	// Get command info
	info, err := registry.GetCommandInfo("info-test")
	if err != nil {
		t.Errorf("GetCommandInfo() unexpected error = %v", err)
	}
	if info == nil {
		t.Error("GetCommandInfo() returned nil info")
	} else {
		if info.Name != originalInfo.Name {
			t.Errorf("GetCommandInfo() Name = %s, want %s", info.Name, originalInfo.Name)
		}
		if info.Category != originalInfo.Category {
			t.Errorf("GetCommandInfo() Category = %s, want %s", info.Category, originalInfo.Category)
		}
		if info.Description != originalInfo.Description {
			t.Errorf("GetCommandInfo() Description = %s, want %s", info.Description, originalInfo.Description)
		}
		if info.Priority != originalInfo.Priority {
			t.Errorf("GetCommandInfo() Priority = %d, want %d", info.Priority, originalInfo.Priority)
		}
	}

	// Test non-existent command
	_, err = registry.GetCommandInfo("non-existent")
	if err == nil {
		t.Error("GetCommandInfo() expected error for non-existent command")
	} else if !strings.Contains(err.Error(), "not found") {
		t.Errorf("GetCommandInfo() error = %v, want to contain 'not found'", err)
	}
}

func TestGlobalRegistryFunctions(t *testing.T) {
	// Save original registry
	originalRegistry := defaultRegistry
	defer func() { defaultRegistry = originalRegistry }()

	// Reset to ensure clean state
	defaultRegistry = nil

	// Test GetRegistry creates a new registry
	registry := GetRegistry()
	if registry == nil {
		t.Error("GetRegistry() returned nil")
	}

	// Second call should return same instance
	registry2 := GetRegistry()
	if registry != registry2 {
		t.Error("GetRegistry() should return same instance on subsequent calls")
	}

	// Test SetRegistry
	customRegistry := NewCommandRegistry()
	SetRegistry(customRegistry)

	registry3 := GetRegistry()
	if registry3 != customRegistry {
		t.Error("GetRegistry() should return custom registry after SetRegistry()")
	}

	// Test RegisterCommand convenience function
	info := CommandInfo{
		Name: "global-test",
		Builder: func() *cobra.Command {
			return &cobra.Command{Use: "global-test"}
		},
	}

	err := RegisterCommand(info)
	if err != nil {
		t.Errorf("RegisterCommand() unexpected error = %v", err)
	}

	// Verify command was registered
	if customRegistry.GetRegisteredCount() != 1 {
		t.Errorf("RegisterCommand() should register 1 command, got %d", customRegistry.GetRegisteredCount())
	}

	// Test BuildAllCommands convenience function
	rootCmd := &cobra.Command{Use: "root"}
	err = BuildAllCommands(rootCmd)
	if err != nil {
		t.Errorf("BuildAllCommands() unexpected error = %v", err)
	}

	if len(rootCmd.Commands()) != 1 {
		t.Errorf("BuildAllCommands() should add 1 command to root, got %d", len(rootCmd.Commands()))
	}
}
