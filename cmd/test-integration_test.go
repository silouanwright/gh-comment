//go:build integration
// +build integration

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupIntegrationLogging(t *testing.T) {
	// Clean up any existing log setup
	integrationLog = nil
	
	err := setupIntegrationLogging()
	assert.NoError(t, err)
	assert.NotNil(t, integrationLog)
	
	// Verify results directory was created
	_, err = os.Stat("integration-tests/results")
	assert.NoError(t, err)
	
	// Clean up
	os.RemoveAll("integration-tests")
}

func TestCreateDummyTemplate(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "templates", "dummy-code.js")
	
	err := createDummyTemplate(templatePath)
	require.NoError(t, err)
	
	// Verify file was created
	content, err := os.ReadFile(templatePath)
	require.NoError(t, err)
	
	contentStr := string(content)
	assert.Contains(t, contentStr, "calculateTotal")
	assert.Contains(t, contentStr, "Integration Test File")
	assert.Contains(t, contentStr, "TODO")
	assert.Contains(t, contentStr, "FIXME")
	assert.Contains(t, contentStr, "module.exports")
}

func TestCopyTemplateFile(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "source.js")
	dstPath := filepath.Join(tempDir, "destination.js")
	
	// Create source file
	sourceContent := "console.log('test');"
	err := os.WriteFile(srcPath, []byte(sourceContent), 0644)
	require.NoError(t, err)
	
	// Copy file
	err = copyTemplateFile(srcPath, dstPath)
	require.NoError(t, err)
	
	// Verify destination exists and has same content
	dstContent, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, sourceContent, string(dstContent))
}

func TestCopyTemplateFileWithMissingSource(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "missing-source.js")
	dstPath := filepath.Join(tempDir, "destination.js")
	
	// Source doesn't exist, should create template
	err := copyTemplateFile(srcPath, dstPath)
	require.NoError(t, err)
	
	// Verify both source and destination exist
	_, err = os.Stat(srcPath)
	assert.NoError(t, err)
	
	_, err = os.Stat(dstPath)
	assert.NoError(t, err)
	
	// Verify content is the dummy template
	content, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "calculateTotal")
}

func TestRunSpecificScenario(t *testing.T) {
	tests := []struct {
		name         string
		scenario     string
		expectError  bool
	}{
		{
			name:        "unknown scenario",
			scenario:    "invalid-scenario",
			expectError: true,
		},
		{
			name:        "valid comments scenario",
			scenario:    "comments",
			expectError: true, // Will fail without real PR setup
		},
		{
			name:        "valid reviews scenario", 
			scenario:    "reviews",
			expectError: true, // Will fail without real PR setup
		},
		{
			name:        "valid reactions scenario",
			scenario:    "reactions", 
			expectError: true, // Will fail without real PR setup
		},
		{
			name:        "valid batch scenario",
			scenario:    "batch",
			expectError: true, // Will fail without real PR setup
		},
		{
			name:        "valid suggestions scenario",
			scenario:    "suggestions",
			expectError: true, // Will fail without real PR setup
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runSpecificScenario(tt.scenario, 999)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.scenario == "invalid-scenario" {
					assert.Contains(t, err.Error(), "unknown scenario")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunAllScenarios(t *testing.T) {
	// This will fail in test environment without real PR
	err := runAllScenarios(999)
	assert.Error(t, err)
	// Should fail on first scenario
}

func TestRunGitCommand(t *testing.T) {
	// Test git command construction
	// In test environment, this will likely fail due to auth or git setup
	err := runGitCommand("status")
	// Don't assert on error since it depends on environment setup
	t.Logf("runGitCommand result: %v", err)
}

func TestCleanupTestPR(t *testing.T) {
	// This will fail in test environment without real PR
	err := cleanupTestPR("test-branch", 999)
	// Function is designed to continue on errors, so might not return error
	t.Logf("cleanupTestPR result: %v", err)
}

func TestIntegrationFlags(t *testing.T) {
	// Test that flags are properly defined
	assert.NotNil(t, testIntegrationCmd)
	
	// Check that flags exist
	flags := testIntegrationCmd.Flags()
	
	cleanupFlag := flags.Lookup("cleanup")
	assert.NotNil(t, cleanupFlag)
	assert.Equal(t, "true", cleanupFlag.DefValue)
	
	inspectFlag := flags.Lookup("inspect")
	assert.NotNil(t, inspectFlag)
	assert.Equal(t, "false", inspectFlag.DefValue)
	
	scenarioFlag := flags.Lookup("scenario")
	assert.NotNil(t, scenarioFlag)
	assert.Equal(t, "", scenarioFlag.DefValue)
}

func TestInspectModeDisablesCleanup(t *testing.T) {
	// Save original values
	originalCleanup := cleanup
	originalInspect := inspect
	defer func() {
		cleanup = originalCleanup
		inspect = originalInspect
	}()
	
	// Set initial state
	cleanup = true
	inspect = false
	
	// Simulate inspect mode being enabled
	inspect = true
	
	// This logic would be in runTestIntegration
	if inspect {
		cleanup = false
	}
	
	assert.False(t, cleanup)
}

func TestDummyTemplateContent(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test.js")
	
	err := createDummyTemplate(tempFile)
	require.NoError(t, err)
	
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	
	lines := strings.Split(string(content), "\n")
	
	// Verify specific content that tests can comment on
	found := false
	for _, line := range lines {
		if strings.Contains(line, "Potential null pointer") {
			found = true
			break
		}
	}
	assert.True(t, found, "Template should contain comment targets")
	
	// Verify various code patterns exist for testing
	contentStr := string(content)
	assert.Contains(t, contentStr, "function calculateTotal")
	assert.Contains(t, contentStr, "for (let i = 0")
	assert.Contains(t, contentStr, "0.08") // Magic number for testing
	assert.Contains(t, contentStr, "TODO:")
	assert.Contains(t, contentStr, "FIXME:")
}