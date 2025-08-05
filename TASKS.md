# Tasks & Development Progress

This file tracks ongoing development tasks, features, and improvements for `gh-comment`. Tasks are organized by priority and status.

## 🚧 In Progress

### 🚨 URGENT BLOCKERS

#### **Phase 3: Test Coverage Boost (3 hours)**
- [ ] Import 20+ missing test files from integration branch
- [ ] Verify coverage increases from 73.3% → 85%+
- [ ] Ensure all imported tests pass without modification
- [ ] Update coverage tracking in CLAUDE.md

### 🐛 **Critical Help Text Issues** - Fix non-working examples discovered during integration testing
**Context**: Integration testing on 2025-08-04 revealed that many help text examples don't work when copy/pasted
**Impact**: Poor user experience, frustration when following documentation
**Test PR**: #12 (https://github.com/silouanwright/gh-comment/pull/12) - Contains test files: `src/api.js`, `src/main.go`, `tests/auth_test.js`
**Test Branch**: `integration-test-20250804-213410`

#### **1. Fix `list` Command Date Placeholder Issues**
- [ ] **Locate invalid examples**: Find all instances of placeholder dates in help text
  - Current bad examples: `"deployment-date"`, `"sprint-start"`, `"release-date"`
  - Files to check: `cmd/list.go`, `cmd/root.go` (global help)
- [ ] **Replace with valid date formats**:
  - Use actual dates: `"2024-01-01"`, `"2024-12-31"`
  - Use relative dates: `"1 week ago"`, `"yesterday"`, `"last month"`
  - Use ISO format examples: `"2024-01-15T09:00:00Z"`
- [ ] **Add date format documentation**:
  - Create a comment block explaining supported date formats
  - Reference Go's time parsing capabilities
  - Include timezone handling examples
- [ ] **Run integration test on PR #12**:
  ```bash
  # Test all list command date examples
  ./gh-comment list 12 --since "1 week ago"
  ./gh-comment list 12 --since "2024-01-01" --until "2024-12-31"
  ./gh-comment list 12 --since "yesterday"
  # Verify no parsing errors occur
  ```

#### **2. Fix `review` Command File Path Examples**
- [ ] **Audit all file references in help text**:
  - Current non-existent files: `auth.go`, `api.js`, `validation.js`, `database.py`
  - These files rarely exist in typical PRs
- [ ] **Replace with commonly existing files**:
  - Use files from PR #12: `src/api.js`, `src/main.go`, `tests/auth_test.js`
  - Use generic names: `README.md`, `main.go`, `index.js`
- [ ] **Add file existence note**:
  - Add comment: "Note: Replace these file names with actual files from your PR"
  - Consider adding a `--validate=false` example for non-existent files
- [ ] **Run integration test on PR #12**:
  ```bash
  # Test review command with actual files from PR #12
  ./gh-comment review 12 "Test review" \
    --comment src/api.js:6:"Good use of middleware" \
    --comment src/main.go:4:"Consider adding error handling" \
    --comment tests/auth_test.js:2:"Add more test cases"
  ```

#### **3. Fix `batch` Command Usage Syntax**
- [ ] **Fix usage line inconsistency**:
  - Current usage: `gh-comment batch <config-file>`
  - Examples show: `gh comment batch 123 review-config.yaml`
  - PR number requirement is unclear
- [ ] **Update usage syntax to**:
  - `gh-comment batch [pr] <config-file>`
  - Or make PR number come from config file
- [ ] **Update all examples to match**:
  - Ensure consistency between usage line and examples
  - Add note about PR number source (CLI vs config file)
- [ ] **Run integration test on PR #12**:
  ```bash
  # Create test batch file
  cat > test-batch.yaml << 'EOF'
  comments:
    - file: src/api.js
      line: 2
      message: "Test batch comment"
  EOF
  
  # Test batch command
  ./gh-comment batch 12 test-batch.yaml
  ```

#### **4. Fix `batch` Command YAML Field Documentation**
- [ ] **Document correct field names**:
  - Clarify that comments use `message` field, not `body`
  - This causes validation errors when users follow incorrect examples
- [ ] **Create comprehensive YAML schema**:
  - Document all supported fields
  - Show type requirements (string, int, array)
  - Include validation rules
- [ ] **Run integration test on PR #12**:
  ```bash
  # Test with correct field name (message)
  cat > correct-batch.yaml << 'EOF'
  pr: 12
  comments:
    - file: src/api.js
      line: 6
      message: "Correct field name test"
  EOF
  ./gh-comment batch correct-batch.yaml
  
  # Test with incorrect field name (body) to verify error message
  cat > incorrect-batch.yaml << 'EOF'
  pr: 12
  comments:
    - file: src/api.js
      line: 6
      body: "This should fail with helpful error"
  EOF
  ./gh-comment batch incorrect-batch.yaml
  # Should show clear error about using 'message' not 'body'
  ```

#### **5. Investigate `review-reply` Command API Issues**
- [ ] **Debug 404 errors**:
  - Test with various review comment IDs
  - Check if issue is with comment ID format or API endpoint
  - Verify GitHub API documentation for correct endpoint
- [ ] **Test API endpoint directly**:
  - Use `gh api` to test the underlying endpoint
  - Compare with GitHub's REST API documentation
- [ ] **Run integration test on PR #12**:
  ```bash
  # First create a review comment to reply to
  ./gh-comment review 12 "Creating review for reply test" \
    --comment src/api.js:2:"This needs a reply"
  
  # List to get the comment ID
  ./gh-comment list 12 --type review
  # Note the comment ID (e.g., 1234567890)
  
  # Test review-reply with the actual comment ID
  ./gh-comment review-reply [COMMENT_ID] "Testing reply functionality"
  
  # Test with --resolve flag (this works)
  ./gh-comment review-reply [COMMENT_ID] --resolve
  ```

#### **6. Fix `lines` Command for New Files**
- [ ] **Investigate new file behavior**:
  - Test why new files show "No commentable lines found"
  - Check if this is GitHub API limitation or our code
- [ ] **Add new file support if possible**:
  - Research GitHub API capabilities for new files
  - Implement support if API allows
- [ ] **Run integration test on PR #12**:
  ```bash
  # Test with new files (these were added in PR #12)
  ./gh-comment lines 12 src/api.js
  ./gh-comment lines 12 src/main.go
  ./gh-comment lines 12 tests/auth_test.js
  
  # Also test with modified files if any exist
  # Expected: New files may show "No commentable lines found"
  # Document whether this is API limitation or bug
  ```

#### **7. Fix `prompts` Command Invalid Example**
- [ ] **Fix incorrect prompt name**:
  - Current example: `gh comment prompts security-comprehensive`
  - Actual name: `security-audit`
- [ ] **Audit all prompt examples**:
  - List actual available prompts with `prompts list`
  - Ensure all examples use valid prompt names
- [ ] **Run integration test on PR #12**:
  ```bash
  # List all available prompts
  ./gh-comment prompts list
  
  # Test the incorrect example (should fail)
  ./gh-comment prompts security-comprehensive
  # Expected: Error message listing available prompts
  
  # Test the correct prompt name
  ./gh-comment prompts security-audit
  
  # Test all available prompts
  ./gh-comment prompts performance
  ./gh-comment prompts architecture
  ./gh-comment prompts code-quality
  ./gh-comment prompts ai-assisted
  ./gh-comment prompts migration
  ```

### 📊 **Help Text Fix Testing Methodology**
**Goal**: Ensure 100% of help text examples work when copy/pasted
**Test Environment**: Use PR #12 (https://github.com/silouanwright/gh-comment/pull/12)

#### **Testing Process for Each Fix**:
1. **Use existing test PR #12**:
   - Already contains: `src/api.js`, `src/main.go`, `tests/auth_test.js`
   - Branch: `integration-test-20250804-213410`
   - Keep PR open for ongoing testing

2. **Test each example after fixes**:
   - Build binary: `go build`
   - Replace `gh comment` → `./gh-comment`
   - Replace PR `123` → `12`
   - Execute command exactly as shown in help
   - Document success/failure

3. **Regression testing checklist**:
   ```bash
   # After each fix, run this test suite on PR #12
   ./gh-comment list 12
   ./gh-comment add 12 "Test comment"
   ./gh-comment review 12 "Test review" --comment src/api.js:2:"Test"
   ./gh-comment react [COMMENT_ID] +1
   ./gh-comment prompts list
   ./gh-comment export 12
   ```

#### **Success Criteria**:
- ✅ All examples in help text execute without errors on PR #12
- ✅ Error messages clearly explain what went wrong
- ✅ No regression in working commands
- ✅ Integration test guide can be followed verbatim

#### **Files to Update**:
- `cmd/*.go` - Fix help text in command files
- `cmd/root.go` - Fix global examples
- `README.md` - Ensure consistency with fixed help text
- `docs/testing/INTEGRATION_TESTING_GUIDE.md` - Update test instructions

#### **Quick Verification Script**:
```bash
# Extract and test all examples from help text
./gh-comment --help | grep -E '^\$ gh comment' | while read -r line; do
  cmd=$(echo "$line" | sed 's/\$ gh comment/\.\/gh-comment/g' | sed 's/123/12/g')
  echo "Testing: $cmd"
  eval "$cmd"
done
```

---

## 🎯 HIGH PRIORITY

### **Real GitHub Integration Tests** - End-to-end workflow testing with actual GitHub PRs
- **Context**: Current testing uses mocks, but we need to verify the extension works with real GitHub APIs
- **Strategy**: Create integration tests that open actual PRs, perform command workflows, verify results, then cleanup
- **Two Test Types**: Automated (full cycle with cleanup) and Manual Verification (leave open for inspection)
- **Conditional Execution**: Run periodically (e.g., every 10th execution) to avoid API rate limits

#### **Phase 1: Basic Integration Test Framework**
- [ ] Create integration test repository or use existing test repo
- [ ] Design test PR template (simple file changes for testing)
- [ ] Create script to programmatically open test PRs via GitHub API
- [ ] Implement basic test runner that can conditionally execute integration tests
- [ ] Add cleanup mechanism to close/delete test PRs after completion

#### **Phase 2: Automated Full-Cycle Tests**
- [ ] **Test Scenario 1: Comment Workflow**
  - Open PR → Verify no comments (`gh comment list`) → Add line comment (`gh comment add`) → Verify comment exists → Close PR
- [ ] **Test Scenario 2: Review Workflow**
  - Open PR → Add review comments (`gh comment add-review`) → Submit review (`gh comment submit-review`) → Verify review exists → Close PR
- [ ] **Test Scenario 3: Reaction Workflow**
  - Open PR with existing comment → Add reaction (`gh comment reply --reaction`) → Verify reaction → Remove reaction → Close PR
- [ ] **Test Scenario 4: Reply Workflow**
  - Open PR with existing comment → Reply to comment (`gh comment reply`) → Verify reply chain → Close PR
- [ ] **Test Scenario 5: Full Interaction Chain**
  - Open PR → Add review comment → Add reaction → Reply to comment → List all (`gh comment list`) → Verify all interactions → Close PR

#### **Phase 3: Manual Verification Tests**
- [ ] **Test Scenario 1: Visual Inspection Workflow**
  - Open PR → Perform various commands → Leave PR open for human verification → Document expected vs actual results
- [ ] **Test Scenario 2: Suggestion Syntax Testing**
  - Open PR → Test `[SUGGEST: code]` expansion → Test `<<<SUGGEST>>>` syntax → Leave open for verification
- [ ] **Test Scenario 3: Edge Case Testing**
  - Test multi-line comments, special characters, long messages, etc. → Leave open for verification

#### **Phase 4: Advanced Integration Features**
- [ ] Implement programmatic PR creation with realistic code changes
- [ ] Add support for testing against different repository types (public/private)
- [ ] Create test data generator for realistic comment scenarios
- [ ] Add integration test reporting and result comparison
- [ ] Implement test result persistence for regression detection

#### **Phase 5: Conditional Execution & CI Integration**
- [ ] Implement "every Nth run" logic for integration tests
- [ ] Add environment variable controls for integration test execution
- [ ] Create separate integration test command (`gh comment test-integration`)
- [ ] Add integration test results to CI/CD pipeline (optional/manual trigger)
- [ ] Create integration test dashboard for tracking results over time

**Technical Requirements**
- Must work with real GitHub API (not mocks)
- Must handle API rate limiting gracefully
- Must clean up test artifacts (PRs, comments, reactions)
- Must be configurable (target repo, test frequency, cleanup behavior)
- Must provide clear success/failure reporting
- Must be runnable both locally and in CI environments

**Success Criteria**
- All refactored commands work correctly with real GitHub APIs
- Integration tests can run automatically and report results
- Manual verification tests provide clear visual confirmation
- Test suite can be run periodically without manual intervention
- Zero false positives/negatives in test results

### **Push Test Coverage to 85%+** - 🚧 IN PROGRESS (Current: 77.0%)
- [x] Generate HTML coverage report: `go test ./cmd -coverprofile=coverage.out && go tool cover -html=coverage.out` ✅
- [x] Identify uncovered code paths ✅
- [x] Add comprehensive tests for parsePositiveInt helper function ✅
- [x] Add edge case tests for add command (TestAddCommandEdgeCases) ✅
- [ ] Continue adding tests for error conditions in low-coverage functions
- [ ] Test edge cases in suggestion parsing logic
- [ ] Add boundary condition tests for YAML batch processing
- **Progress**: Coverage improved from 73.3% → 77.0% (+3.7 percentage points)
- **Recent additions**: 7 parsePositiveInt tests + 5 add command edge case tests
- **Next targets**: Functions with <80% coverage (runAdd: 64.1%, getPRContext: 55.6%, processAsReview: 62.1%)

### **Add Input Length Validation**
- [ ] Define constants for GitHub API limits
- [ ] Add comment body length validation (GitHub max: 65,536 chars)
- [ ] Add file path validation to prevent directory traversal
```go
// TODO: Add to cmd/helpers.go
const (
    MaxCommentLength = 65536 // GitHub's actual limit
    MaxFilePathLength = 4096 // Reasonable file path limit
)

func validateCommentBody(body string) error {
    if len(body) > MaxCommentLength {
        return fmt.Errorf("comment too long: %d chars (max %d)", len(body), MaxCommentLength)
    }
    return nil
}
```

### **Add More Comprehensive Error Context**
- [ ] Enhance API error messages with suggested actions
- [ ] Add help hints for common error scenarios
- [ ] Include relevant documentation links in error messages
```go
// TODO: Enhance error messages
func formatAPIErrorWithHint(operation string, err error) error {
    hint := getHintForOperation(operation)
    return fmt.Errorf("GitHub API error during %s: %w\n💡 Hint: %s", operation, err, hint)
}
```

---

## 🔄 MEDIUM PRIORITY

### **Cross-Platform Testing** - Ensure compatibility across all platforms
- [ ] Add Windows-specific test scenarios (path separators, line endings)
- [ ] Test shell compatibility (bash, zsh, fish, PowerShell)
- [ ] Verify testscript behavior on different operating systems
- [ ] Add platform-specific golden files if needed
- [ ] Test GitHub CLI integration across platforms

### **Automated Test Data Cleanup** - Implement cleanup routines for E2E tests
- [ ] Add test repository cleanup after E2E test runs
- [ ] Implement comment cleanup for failed test scenarios
- [ ] Add test data isolation to prevent cross-test contamination
- [ ] Create test data lifecycle management

### **Enhancement: Create Separate Integration Test Workflow**
- [ ] **Create** `.github/workflows/integration.yml` for manual integration testing
- [ ] **Trigger**: Manual dispatch only (`workflow_dispatch`)
- [ ] **Environment**: Separate environment with proper secrets and permissions
- [ ] **Priority**: Medium - Nice to have for organized testing

### **Enhancement: Update Integration Test Documentation**
- [ ] **Update** `docs/testing/INTEGRATION_TESTING.md`
- [ ] **Add**: Best practices from recent successful integration testing
- [ ] **Document**: How to test functionality changes like we just did
- [ ] **Priority**: Medium - Helps future development

---

## 📋 PLANNED FEATURES

### Core Features
- [ ] **Configuration file support** - Default flags and repository settings
  - [ ] Design configuration file format (YAML/JSON)
  - [ ] Implement config file parsing
  - [ ] Add `--config` flag support
  - [ ] Create default config generation command
  - [ ] Support default author, format, color settings
  - [ ] Add table style configuration

- [ ] **Template system** - Reusable comment patterns and workflows
  - [ ] Design template file format
  - [ ] Implement template loading and substitution
  - [ ] Add built-in templates for common scenarios
  - [ ] Create template sharing mechanism

- [ ] **Enhanced Help System** - Better help text following GitHub CLI patterns
  - [ ] Add structured examples with descriptions
  - [ ] Improve long-form help documentation
  - [ ] Add contextual help for errors
  - [ ] Create help builder utilities

### Quality & Performance
- [ ] **Performance Optimizations**
  - [ ] Optimize comment fetching with pagination
  - [ ] Add caching for frequently accessed data
  - [ ] Implement parallel API calls where possible
  - [ ] Monitor and optimize memory usage

### User Experience
- [ ] **Professional Table Output** - Replace manual string formatting with `olekukonko/tablewriter`
  - [ ] Add table output for `list` command
  - [ ] Support auto-wrapping and formatting
  - [ ] Add configurable table styles
  - [ ] Used by 500+ CLI tools including Kubernetes tools

- [ ] **Color Support** - Add color output with `fatih/color`
  - [ ] Add color coding for different comment types
  - [ ] Color code authors, timestamps, and status
  - [ ] Add `--no-color` flag for compatibility
  - [ ] Respect terminal color capabilities

- [ ] **Progress Indicators** - Add progress bars for long operations with `schollz/progressbar`
  - [ ] Show progress when fetching many comments
  - [ ] Add progress for batch operations
  - [ ] Display ETA for long-running commands

- [ ] **Batch operations** - Apply operations to multiple comments at once
  - [ ] Design batch operation syntax
  - [ ] Implement batch comment creation
  - [ ] Add batch reaction management
  - [ ] Create batch editing capabilities

- [ ] **Export functionality** - Export comments to various formats
  - [ ] JSON export format
  - [ ] CSV export for spreadsheet analysis
  - [ ] Markdown export for documentation
  - [ ] HTML export for presentations
  - [ ] Add `export` subcommand

---

## 🔧 LOW PRIORITY

### **Code Organization Improvements**
- [ ] Group related functions in files (parsing, validation, display)
- [ ] Consider extracting large functions (>50 lines) into smaller units
- [ ] Add more granular unit tests for helper functions

### **Performance Optimizations**
- [ ] Add benchmarks for suggestion parsing
- [ ] Profile memory usage during large comment listings
- [ ] Consider pagination for very large PRs

### **Developer Experience**
- [ ] Add more debug logging in verbose mode
- [ ] Create troubleshooting guide for common issues
- [ ] Add shell completion improvements

### **Testing Enhancements**
- [ ] Add fuzz testing for suggestion syntax parsing
- [ ] Test Unicode handling in comments and file paths
- [ ] Add tests for very large PRs (100+ comments)
- [ ] Test rate limiting scenarios

### **Integration Test Improvements**
- [ ] Add automated integration test runner
- [ ] Create test data fixtures for consistent testing
- [ ] Add performance benchmarks for integration tests

### **Security Hardening**
- [ ] Add rate limiting protection for API calls
- [ ] Implement request timeouts for all HTTP operations
- [ ] Add input sanitization for file paths
- [ ] Consider adding audit logging for sensitive operations

### **Optional CI/CD Improvements**
- [ ] **Fix golangci-lint Configuration**
  - [ ] **Issue**: Lint may fail due to deprecated config options (not currently blocking)
  - [ ] **Location**: `.golangci.yml` lines 7 and 10
  - [ ] **Fix**: Remove deprecated `check-shadowing` and `maligned` settings if they cause issues
  - [ ] **Priority**: Low - only address if linting actually fails

- [ ] **Fix Benchmark PR Commenting Permissions**
  - [ ] **Issue**: Benchmark step may fail with "Resource not accessible by integration"
  - [ ] **Location**: `.github/workflows/test.yml` lines 172-187
  - [ ] **Fix**: Add proper permissions or make commenting optional if issues arise
  - [ ] **Priority**: Low - only address if benchmarking actually fails

---

## 🔍 INTEGRATION BRANCH AUDIT - Missing Features & Functionality

### **CRITICAL DISCOVERY: Major Architectural Changes in Integration Branch**

The integration branch (integration-test-20250802-224635) contains significant architectural improvements and features that were developed but never merged to main.

### **1. Command Architecture Restructuring - PARTIALLY COMPLETE**
- **COMPLETED**: ✅ `react` command extracted for emoji reactions
- **COMPLETED**: ✅ `review-reply` command created for review comment threading
- **MISSING**: ❌ `reply` command still exists on main (should be removed)
- **ACTION**: Delete `reply.go` and its tests, as functionality is now split between:
  - `add` → Issue comments (general discussion)
  - `review-reply` → Review comment replies (line-specific)
  - `react` → Emoji reactions

### **2. Enhanced Commands & Features - Files that differ:**
```
cmd/add.go               - Enhanced validation and error handling
cmd/batch.go            - Improved YAML processing and validation
cmd/close-pending-review.go - Better documentation and examples
cmd/edit.go             - Enhanced message handling
cmd/helpers.go          - New helper functions for validation
cmd/lines.go            - Better line grouping and display
cmd/list.go             - Improved filtering and output formatting
cmd/review.go           - Enhanced review creation workflow
cmd/root.go             - Updated help text and examples
```

### **3. Missing Test Files & Coverage:**
```
cmd/batch_test.go            - Enhanced batch command tests
cmd/close-pending-review_test.go - Comprehensive pending review tests
cmd/command_execution_test.go - Integration command execution tests
cmd/dependency_injection_test.go - DI pattern tests
cmd/e2e_test.go              - End-to-end test scenarios
cmd/helpers_test.go          - Helper function tests
cmd/integration-scenarios_test.go - Complex workflow tests
cmd/lines_test.go            - Lines command tests
cmd/list_comprehensive_test.go - Comprehensive list tests
cmd/react_test.go            - React command tests (copied)
cmd/reply_integration_test.go - Reply integration tests
cmd/reply_targeted_test.go   - Targeted reply tests
cmd/review_test.go           - Review command tests
cmd/review-reply_test.go     - Review-reply tests (copied)
cmd/review-reply_targeted_test.go - Targeted review-reply tests
cmd/utility_functions_test.go - Utility function tests
```

### **4. Documentation & Examples - Missing from main:**
```
docs/testing/INTEGRATION_TESTING_GUIDE.md - Comprehensive testing guide
examples/comprehensive-review.yaml - Review workflow example
examples/performance-review.yaml   - Performance review template
examples/security-audit.yaml       - Security audit template
src/api.js                        - Example API file for testing
src/main.go                       - Example Go file for testing
tests/auth_test.js                - Example test file
```

### **5. Integration Strategy Summary**
1. Complete command restructuring first
2. Port all tests to ensure safety
3. Cherry-pick feature improvements
4. Update documentation last

**Without these changes:**
- Missing critical bug fixes
- Lower test coverage (73.3% vs likely >80%)
- Incomplete command architecture
- Missing user-friendly features
- Documentation gaps

---

## 📝 Task Management Notes

### How to Use This File
1. **Add new tasks** under appropriate sections
2. **Move tasks** between sections as they progress
3. **Check off subtasks** using `- [x]` syntax
4. **Archive completed tasks** to the "Recently Completed" section
5. **Update dates** when moving tasks to completed

### Task Priorities
- **High**: Critical functionality, bug fixes, security issues
- **Medium**: Important features, performance improvements
- **Low**: Nice-to-have features, documentation improvements

### Status Indicators
- `🚧` In Progress
- `📋` Planned
- `✅` Completed
- `🎯` Success Metrics
- `⚠️` Blocked/Issues
- `🔄` Under Review

*This project is already **exceptional (A- grade)** and production-ready. These tasks will polish it to **industry-leading (A+ grade)** quality.*

Last updated: August 2025