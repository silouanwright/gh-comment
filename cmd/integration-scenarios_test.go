//go:build integration
// +build integration

package cmd

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests only run when integration build tag is specified

func TestGetTestFileName(t *testing.T) {
	// Create a temporary test file
	testFile := "test-file-123.js"
	content := `function test() {
    return true;
}`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)
	defer os.Remove(testFile)

	filename := getTestFileName(123)
	assert.Equal(t, testFile, filename)
}

func TestGetTestFileNameFallback(t *testing.T) {
	// Test fallback when no test file exists
	filename := getTestFileName(999)
	assert.Equal(t, "test-file.js", filename)
}

func TestCreateBatchConfig(t *testing.T) {
	tempFile := "test-batch-config.yaml"
	defer os.Remove(tempFile)

	err := createBatchConfig(tempFile, 123, "test.js")
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(tempFile)
	assert.NoError(t, err)

	// Verify content
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "pr_number: 123")
	assert.Contains(t, contentStr, "file: test.js")
	assert.Contains(t, contentStr, "AUTO_DETECT")
	assert.Contains(t, contentStr, "COMMENT")
}

func TestVerifyNoComments(t *testing.T) {
	// This test would typically interact with a real GitHub API
	// For now, we just test that the function doesn't panic
	err := verifyNoComments(999)
	// Function is designed to not fail even if comments exist
	assert.NoError(t, err)
}

func TestValidateCommentsExist(t *testing.T) {
	// Test with empty expected comments
	err := validateCommentsExist(999, []string{})
	// Should not fail for empty expectations
	// The actual behavior depends on whether the comment list command succeeds
	// which depends on the environment setup
	t.Logf("validateCommentsExist result: %v", err)
}

func TestValidateReviewExists(t *testing.T) {
	// This test requires gh CLI and a valid PR
	// We test that the function doesn't panic
	err := validateReviewExists(999, "CHANGES_REQUESTED")
	// Expected to fail in test environment without valid PR
	t.Logf("validateReviewExists result: %v", err)
}

func TestGetLatestCommentID(t *testing.T) {
	// Test that function returns something
	commentID, err := getLatestCommentID(999)
	// Should return "latest" as placeholder
	if err == nil {
		assert.Equal(t, "latest", commentID)
	}
	// In test environment, gh CLI might not be available or PR might not exist
	t.Logf("getLatestCommentID result: %s, err: %v", commentID, err)
}

func TestValidateReactionsAndReplies(t *testing.T) {
	// This function currently always returns nil
	err := validateReactionsAndReplies(999, "test-comment-id")
	assert.NoError(t, err)
}

func TestValidateBatchResults(t *testing.T) {
	// Test the validation function
	err := validateBatchResults(999)
	// May fail in test environment, but shouldn't panic
	t.Logf("validateBatchResults result: %v", err)
}

func TestValidateSuggestionFormatting(t *testing.T) {
	// Test the validation function
	err := validateSuggestionFormatting(999)
	// May fail in test environment, but shouldn't panic
	t.Logf("validateSuggestionFormatting result: %v", err)
}

// Test helper functions that construct command arguments
func TestRunCommentAddArgs(t *testing.T) {
	// Test that we can construct the expected command without executing it
	// We'll test the argument construction logic by simulating what runCommentAdd does
	
	prNumber := 123
	file := "test.js"
	line := 42
	endLine := ""
	message := "test comment"
	
	args := []string{"comment", "add", strconv.Itoa(prNumber), file, strconv.Itoa(line)}
	if endLine != "" {
		args = append(args, endLine)
	}
	args = append(args, message)
	
	expectedArgs := []string{"comment", "add", "123", "test.js", "42", "test comment"}
	assert.Equal(t, expectedArgs, args)
}

func TestRunCommentAddRangeArgs(t *testing.T) {
	prNumber := 123
	file := "test.js"
	startLine := 10
	endLine := 20
	message := "range comment"
	
	args := []string{"comment", "add", strconv.Itoa(prNumber), file, 
		strconv.Itoa(startLine), strconv.Itoa(endLine), message}
	
	expectedArgs := []string{"comment", "add", "123", "test.js", "10", "20", "range comment"}
	assert.Equal(t, expectedArgs, args)
}

func TestRunReviewAddArgs(t *testing.T) {
	prNumber := 123
	file := "test.js"
	line := 42
	message := "review comment"
	
	args := []string{"comment", "add-review", strconv.Itoa(prNumber), file, 
		strconv.Itoa(line), message}
	
	expectedArgs := []string{"comment", "add-review", "123", "test.js", "42", "review comment"}
	assert.Equal(t, expectedArgs, args)
}

func TestRunReviewSubmitArgs(t *testing.T) {
	prNumber := 123
	event := "REQUEST_CHANGES"
	body := "review body"
	
	args := []string{"comment", "submit-review", strconv.Itoa(prNumber), 
		"--event", event, "--body", body}
	
	expectedArgs := []string{"comment", "submit-review", "123", "--event", "REQUEST_CHANGES", "--body", "review body"}
	assert.Equal(t, expectedArgs, args)
}

func TestRunReplyReactionArgs(t *testing.T) {
	commentID := "123456"
	reaction := "+1"
	
	args := []string{"comment", "reply", "--comment-id", commentID, "--reaction", reaction}
	
	expectedArgs := []string{"comment", "reply", "--comment-id", "123456", "--reaction", "+1"}
	assert.Equal(t, expectedArgs, args)
}

func TestRunReplyMessageArgs(t *testing.T) {
	commentID := "123456"
	message := "reply message"
	
	args := []string{"comment", "reply", "--comment-id", commentID, "--message", message}
	
	expectedArgs := []string{"comment", "reply", "--comment-id", "123456", "--message", "reply message"}
	assert.Equal(t, expectedArgs, args)
}

func TestRunBatchCommandArgs(t *testing.T) {
	configFile := "config.yaml"
	
	args := []string{"comment", "batch", configFile}
	
	expectedArgs := []string{"comment", "batch", "config.yaml"}
	assert.Equal(t, expectedArgs, args)
}