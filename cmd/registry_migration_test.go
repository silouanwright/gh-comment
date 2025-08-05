package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandCategories(t *testing.T) {
	// Test that category constants are defined correctly
	categories := []string{CategoryCore, CategoryManage, CategoryAdmin, CategoryUtility, CategoryTest}
	expectedCategories := []string{"core", "manage", "admin", "utility", "test"}
	
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
	
	for i, expected := range expectedCategories {
		if i < len(categories) && categories[i] != expected {
			t.Errorf("Category[%d] = %s, want %s", i, categories[i], expected)
		}
	}
}

func TestCommandPriorities(t *testing.T) {
	// Test that priority constants have expected values
	if PriorityHigh != 1 {
		t.Errorf("PriorityHigh = %d, want 1", PriorityHigh)
	}
	if PriorityMedium != 5 {
		t.Errorf("PriorityMedium = %d, want 5", PriorityMedium)
	}
	if PriorityLow != 10 {
		t.Errorf("PriorityLow = %d, want 10", PriorityLow)
	}
}

func TestRegisterCoreCommands(t *testing.T) {
	registry := NewCommandRegistry()
	
	err := RegisterCoreCommands(registry)
	if err != nil {
		t.Errorf("RegisterCoreCommands() unexpected error = %v", err)
	}
	
	// Check that core commands were registered
	expectedCommands := []string{"add", "review", "list", "batch"}
	for _, cmdName := range expectedCommands {
		_, err := registry.GetCommandInfo(cmdName)
		if err != nil {
			t.Errorf("Core command %s not registered: %v", cmdName, err)
		}
	}
	
	// Verify command count
	if registry.GetRegisteredCount() != len(expectedCommands) {
		t.Errorf("RegisterCoreCommands() registered %d commands, want %d", 
			registry.GetRegisteredCount(), len(expectedCommands))
	}
	
	// Verify commands are in correct category
	categories := registry.GetCommandsByCategory()
	coreCommands := categories[CategoryCore]
	if len(coreCommands) != len(expectedCommands) {
		t.Errorf("Core category has %d commands, want %d", len(coreCommands), len(expectedCommands))
	}
}

func TestRegisterManagementCommands(t *testing.T) {
	registry := NewCommandRegistry()
	
	err := RegisterManagementCommands(registry)
	if err != nil {
		t.Errorf("RegisterManagementCommands() unexpected error = %v", err)
	}
	
	expectedCommands := []string{"edit", "resolve", "react", "review-reply", "close-pending-review"}
	
	// Verify all management commands were registered
	for _, cmdName := range expectedCommands {
		info, err := registry.GetCommandInfo(cmdName)
		if err != nil {
			t.Errorf("Management command %s not registered: %v", cmdName, err)
		} else if info.Category != CategoryManage {
			t.Errorf("Command %s has category %s, want %s", cmdName, info.Category, CategoryManage)
		}
	}
	
	if registry.GetRegisteredCount() != len(expectedCommands) {
		t.Errorf("RegisterManagementCommands() registered %d commands, want %d", 
			registry.GetRegisteredCount(), len(expectedCommands))
	}
}

func TestRegisterAdminCommands(t *testing.T) {
	registry := NewCommandRegistry()
	
	err := RegisterAdminCommands(registry)
	if err != nil {
		t.Errorf("RegisterAdminCommands() unexpected error = %v", err)
	}
	
	expectedCommands := []string{"config", "export"}
	
	for _, cmdName := range expectedCommands {
		info, err := registry.GetCommandInfo(cmdName)
		if err != nil {
			t.Errorf("Admin command %s not registered: %v", cmdName, err)
		} else if info.Category != CategoryAdmin {
			t.Errorf("Command %s has category %s, want %s", cmdName, info.Category, CategoryAdmin)
		}
	}
}

func TestRegisterUtilityCommands(t *testing.T) {
	registry := NewCommandRegistry()
	
	err := RegisterUtilityCommands(registry)
	if err != nil {
		t.Errorf("RegisterUtilityCommands() unexpected error = %v", err)
	}
	
	expectedCommands := []string{"lines", "prompts"}
	
	for _, cmdName := range expectedCommands {
		info, err := registry.GetCommandInfo(cmdName)
		if err != nil {
			t.Errorf("Utility command %s not registered: %v", cmdName, err)
		} else if info.Category != CategoryUtility {
			t.Errorf("Command %s has category %s, want %s", cmdName, info.Category, CategoryUtility)
		}
	}
}

func TestRegisterTestCommands(t *testing.T) {
	registry := NewCommandRegistry()
	
	err := RegisterTestCommands(registry)
	if err != nil {
		t.Errorf("RegisterTestCommands() unexpected error = %v", err)
	}
	
	expectedCommands := []string{"test-integration"}
	
	for _, cmdName := range expectedCommands {
		info, err := registry.GetCommandInfo(cmdName)
		if err != nil {
			t.Errorf("Test command %s not registered: %v", cmdName, err)
		} else if info.Category != CategoryTest {
			t.Errorf("Command %s has category %s, want %s", cmdName, info.Category, CategoryTest)
		}
	}
}

func TestRegisterAllCommands(t *testing.T) {
	registry := NewCommandRegistry()
	
	err := RegisterAllCommands(registry)
	if err != nil {
		t.Errorf("RegisterAllCommands() unexpected error = %v", err)
	}
	
	// Count expected total commands
	expectedTotal := 4 + 5 + 2 + 2 + 1 // core + manage + admin + utility + test
	if registry.GetRegisteredCount() != expectedTotal {
		t.Errorf("RegisterAllCommands() registered %d commands, want %d", 
			registry.GetRegisteredCount(), expectedTotal)
	}
	
	// Verify all categories are present
	categories := registry.GetCommandsByCategory()
	expectedCategories := []string{CategoryCore, CategoryManage, CategoryAdmin, CategoryUtility, CategoryTest}
	for _, expectedCategory := range expectedCategories {
		if _, exists := categories[expectedCategory]; !exists {
			t.Errorf("RegisterAllCommands() missing category %s", expectedCategory)
		}
	}
	
	// Verify specific high-priority commands exist
	highPriorityCommands := []string{"add", "review", "list"}
	for _, cmdName := range highPriorityCommands {
		info, err := registry.GetCommandInfo(cmdName)
		if err != nil {
			t.Errorf("High priority command %s not found: %v", cmdName, err)
		} else if info.Priority != PriorityHigh {
			t.Errorf("Command %s has priority %d, want %d", cmdName, info.Priority, PriorityHigh)
		}
	}
}

func TestCommandBuilders(t *testing.T) {
	// Test that each command builder returns a valid command
	builders := map[string]func() *cobra.Command{
		"add":                   buildAddCommand,
		"review":               buildReviewCommand,
		"list":                 buildListCommand,
		"batch":                buildBatchCommand,
		"edit":                 buildEditCommand,
		"resolve":              buildResolveCommand,
		"react":                buildReactCommand,
		"review-reply":         buildReviewReplyCommand,
		"close-pending-review": buildClosePendingReviewCommand,
		"config":               buildConfigCommand,
		"export":               buildExportCommand,
		"lines":                buildLinesCommand,
		"prompts":              buildPromptsCommand,
		"test-integration":     buildTestIntegrationCommand,
	}
	
	for cmdName, builder := range builders {
		t.Run(cmdName, func(t *testing.T) {
			cmd := builder()
			if cmd == nil {
				t.Errorf("Builder for %s returned nil", cmdName)
				return
			}
			
			if cmd.Use == "" {
				t.Errorf("Builder for %s created command with empty Use", cmdName)
			}
			
			// Verify the Use field matches expected command name
			if !strings.HasPrefix(cmd.Use, cmdName) {
				t.Errorf("Builder for %s created command with Use=%s, expected to start with %s", 
					cmdName, cmd.Use, cmdName)
			}
			
			if cmd.Short == "" {
				t.Errorf("Builder for %s created command with empty Short description", cmdName)
			}
			
			// Note: We can't test RunE functions without significant setup,
			// but we can verify they're not nil for most commands
			if cmd.RunE == nil {
				// Some commands might legitimately have nil RunE if they're just containers
				t.Logf("Builder for %s created command with nil RunE (might be intentional)", cmdName)
			}
		})
	}
}

func TestInitializeCommandRegistry(t *testing.T) {
	// Save original registry
	originalRegistry := defaultRegistry
	defer func() { defaultRegistry = originalRegistry }()
	
	// Reset to clean state
	defaultRegistry = nil
	
	err := InitializeCommandRegistry()
	if err != nil {
		t.Errorf("InitializeCommandRegistry() unexpected error = %v", err)
	}
	
	registry := GetRegistry()
	if registry.GetRegisteredCount() == 0 {
		t.Error("InitializeCommandRegistry() should have registered commands")
	}
	
	// Verify we can get some core commands
	coreCommands := []string{"add", "review", "list"}
	for _, cmdName := range coreCommands {
		_, err := registry.GetCommand(cmdName)
		if err != nil {
			t.Errorf("InitializeCommandRegistry() did not register %s: %v", cmdName, err)
		}
	}
}

func TestBuildAndRegisterCommands(t *testing.T) {
	// Save original registry
	originalRegistry := defaultRegistry
	defer func() { defaultRegistry = originalRegistry }()
	
	// Reset to clean state
	defaultRegistry = nil
	
	rootCmd := &cobra.Command{Use: "root"}
	
	err := BuildAndRegisterCommands(rootCmd)
	if err != nil {
		t.Errorf("BuildAndRegisterCommands() unexpected error = %v", err)
	}
	
	// Verify commands were added to root
	if len(rootCmd.Commands()) == 0 {
		t.Error("BuildAndRegisterCommands() should have added commands to root")
	}
	
	// Verify specific commands exist in root
	coreCommands := []string{"add", "review", "list"}
	for _, expectedCmd := range coreCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if strings.HasPrefix(cmd.Use, expectedCmd) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BuildAndRegisterCommands() did not add %s to root command", expectedCmd)
		}
	}
}

func TestRegistryMigrationIntegration(t *testing.T) {
	// Test the complete migration workflow
	registry := NewCommandRegistry()
	
	// Register all commands
	err := RegisterAllCommands(registry)
	if err != nil {
		t.Errorf("RegisterAllCommands() failed: %v", err)
	}
	
	// Create root command
	rootCmd := &cobra.Command{Use: "gh-comment"}
	
	// Build all commands and add to root
	err = registry.BuildAll(rootCmd)
	if err != nil {
		t.Errorf("BuildAll() failed: %v", err)
	}
	
	// Verify the integration worked
	if len(rootCmd.Commands()) == 0 {
		t.Error("Integration test failed: no commands added to root")
	}
	
	// Test command discoverability
	commandNames := registry.ListCommands()
	if len(commandNames) == 0 {
		t.Error("Integration test failed: no commands discoverable")
	}
	
	// Test category organization
	categories := registry.GetCommandsByCategory()
	if len(categories) == 0 {
		t.Error("Integration test failed: no command categories")
	}
	
	// Verify core commands are prioritized correctly
	coreCommands := categories[CategoryCore]
	if len(coreCommands) == 0 {
		t.Error("Integration test failed: no core commands")
	}
	
	// Check priority ordering in core category
	for i := 1; i < len(coreCommands); i++ {
		if coreCommands[i-1].Priority > coreCommands[i].Priority {
			t.Errorf("Integration test failed: core commands not sorted by priority")
		}
	}
}

func TestRegistryVsInitPatternComparison(t *testing.T) {
	// This test demonstrates the benefits of the registry pattern over init()
	
	// 1. Test discoverability - registry pattern allows introspection
	registry := NewCommandRegistry()
	err := RegisterAllCommands(registry)
	if err != nil {
		t.Errorf("Failed to register commands: %v", err)
	}
	
	// We can now discover commands programmatically
	allCommands := registry.ListCommands()
	if len(allCommands) == 0 {
		t.Error("Registry pattern should allow command discovery")
	}
	
	// 2. Test categorization - registry pattern allows organization
	categories := registry.GetCommandsByCategory()
	if len(categories) < 2 { // Should have at least core and one other category
		t.Error("Registry pattern should support command categorization")
	}
	
	// 3. Test selective registration - registry pattern allows control
	selectiveRegistry := NewCommandRegistry()
	err = RegisterCoreCommands(selectiveRegistry)
	if err != nil {
		t.Errorf("Failed to register core commands: %v", err)
	}
	
	coreCount := selectiveRegistry.GetRegisteredCount()
	fullCount := registry.GetRegisteredCount()
	
	if coreCount >= fullCount {
		t.Error("Selective registration should register fewer commands than full registration")
	}
	
	// 4. Test metadata access - registry pattern provides rich information
	info, err := registry.GetCommandInfo("add")
	if err != nil {
		t.Errorf("Failed to get command info: %v", err)
	}
	
	if info.Category == "" || info.Description == "" {
		t.Error("Registry pattern should provide rich command metadata")
	}
	
	// These capabilities are difficult or impossible with the init() pattern
	t.Logf("Registry pattern successfully provides:")
	t.Logf("- Command discovery: %d commands", len(allCommands))
	t.Logf("- Categorization: %d categories", len(categories))
	t.Logf("- Selective registration: %d/%d commands", coreCount, fullCount)
	t.Logf("- Rich metadata: category=%s, priority=%d", info.Category, info.Priority)
}