# Integration Testing Strategy - Dogfooding Approach

## Overview

This document outlines the integration testing strategy for `gh-comment` using a "dogfooding" approach - testing the extension on its own GitHub repository by creating real pull requests and exercising all functionality live.

## Core Concept

The beauty of `gh-comment` being hosted on GitHub is that we can use the tool to test itself:

1. **Create Test PR**: Automatically create a branch with dummy changes and open a PR
2. **Exercise All Commands**: Run the full suite of `gh-comment` commands against the test PR
3. **Validate Results**: Use both `gh-comment` commands and GitHub API calls to verify functionality
4. **Cleanup or Inspect**: Either auto-close the PR or leave it open for manual inspection

## Test Architecture

### Test Runner Command
```bash
# Local integration test runner
go run . test-integration [--cleanup] [--inspect] [--scenario=<name>]
```

**Flags:**
- `--cleanup`: Auto-close PR after tests complete (default)
- `--inspect`: Leave PR open for manual inspection
- `--scenario=<name>`: Run specific test scenario only

### Test Template Files

**Template File: `test-templates/dummy-code.js`**
```javascript
// Integration Test File - Contains intentional issues for commenting
function calculateTotal(items) {
    let total = 0;
    for (let i = 0; i < items.length; i++) {
        total += items[i].price * items[i].quantity; // Potential null pointer
    }
    return total; // Missing input validation
}

// TODO: Add error handling
// FIXME: Handle empty arrays
const processOrder = (order) => {
    const total = calculateTotal(order.items);
    return { total, tax: total * 0.08 }; // Hardcoded tax rate
};
```

This template provides multiple commenting opportunities:
- Line-specific issues (null pointer, hardcoded values)
- Suggestions for improvements
- Multi-line comment opportunities
- Range comments for entire functions

## Test Scenarios

### Scenario 1: Basic Comment Workflow
```bash
# 1. Create test PR
./create-test-pr.sh "integration-test-comments-$(date +%s)"

# 2. Verify no comments exist
go run . list

# 3. Add line comment
go run . add --line 4 --message "Add null check for items array"

# 4. Add range comment  
go run . add --start-line 8 --end-line 12 --message "Consider extracting tax calculation to constant"

# 5. Validate comments exist
go run . list | grep "Add null check"
go run . list | grep "tax calculation"

# 6. Cleanup
./cleanup-test-pr.sh
```

### Scenario 2: Review Comment Workflow
```bash
# 1. Create test PR
./create-test-pr.sh "integration-test-review-$(date +%s)"

# 2. Add review comments
go run . add-review --line 4 --message "Needs input validation"
go run . add-review --line 8 --message "Magic number should be configurable"

# 3. Submit review
go run . submit-review --event REQUEST_CHANGES --body "Please address these issues"

# 4. Validate review exists
gh pr view --json reviews | jq '.reviews[0].state' | grep "CHANGES_REQUESTED"

# 5. Cleanup
./cleanup-test-pr.sh
```

### Scenario 3: Reaction & Reply Workflow
```bash
# 1. Create test PR with existing comment
./create-test-pr.sh "integration-test-reactions-$(date +%s)"
go run . add --line 4 --message "This needs refactoring"

# 2. Add reaction to comment
COMMENT_ID=$(go run . list --json | jq -r '.[0].id')
go run . reply --comment-id $COMMENT_ID --reaction "+1"

# 3. Reply to comment
go run . reply --comment-id $COMMENT_ID --message "I agree, let's create a separate function"

# 4. Validate reactions and replies
go run . list | grep "👍 1" # Reaction count
go run . list | grep "I agree" # Reply message

# 5. Cleanup
./cleanup-test-pr.sh
```

### Scenario 4: Batch Operations
```bash
# 1. Create test PR
./create-test-pr.sh "integration-test-batch-$(date +%s)"

# 2. Create batch config
cat > test-batch.yaml << EOF
repository: silouan.wright/gh-comment
pr_number: AUTO_DETECT
comments:
  - type: issue
    message: "Overall code quality looks good"
  - type: review
    file: test-templates/dummy-code.js
    line: 4
    message: "Add input validation here"
  - type: review
    file: test-templates/dummy-code.js
    start_line: 8
    end_line: 12
    message: "Extract tax calculation logic"
review:
  event: COMMENT
  body: "Automated review from integration tests"
EOF

# 3. Execute batch operations
go run . batch test-batch.yaml

# 4. Validate all comments created
go run . list | wc -l # Should show 3 comments

# 5. Cleanup
./cleanup-test-pr.sh
rm test-batch.yaml
```

### Scenario 5: Suggestion Syntax Testing
```bash
# 1. Create test PR
./create-test-pr.sh "integration-test-suggestions-$(date +%s)"

# 2. Add suggestion comments
go run . add --line 4 --message "[SUGGEST: if (!items || items.length === 0) throw new Error('Invalid items')]"
go run . add --line 12 --message "<<<SUGGEST>>>
const TAX_RATE = 0.08;
return { total, tax: total * TAX_RATE };
<<<SUGGEST>>>"

# 3. Validate suggestions rendered correctly
go run . list --format json | jq '.[0].body' | grep "suggestion"

# 4. Cleanup or inspect
./cleanup-test-pr.sh
```

## Validation Methods

### 1. Command Output Validation
```bash
# Verify command success
if go run . add --line 4 --message "test"; then
    echo "✅ Add command succeeded"
else
    echo "❌ Add command failed"
    exit 1
fi
```

### 2. GitHub API Verification
```bash
# Verify comment via GitHub API
COMMENT_COUNT=$(gh api "repos/silouan.wright/gh-comment/pulls/$PR_NUMBER/comments" | jq length)
if [[ $COMMENT_COUNT -gt 0 ]]; then
    echo "✅ Comments verified via API"
else
    echo "❌ No comments found via API"
    exit 1
fi
```

### 3. Cross-Command Validation
```bash
# Add comment, then verify with list
go run . add --line 4 --message "Test comment"
if go run . list | grep "Test comment"; then
    echo "✅ Comment verified via list command"
else
    echo "❌ Comment not found in list output"
    exit 1
fi
```

## Implementation Structure

### Directory Structure
```
integration-tests/
├── runner.go              # Main test runner
├── scenarios/              # Individual test scenarios
│   ├── basic-comments.go
│   ├── review-workflow.go
│   ├── reactions-replies.go
│   ├── batch-operations.go
│   └── suggestions.go
├── scripts/                # Shell utilities
│   ├── create-test-pr.sh
│   ├── cleanup-test-pr.sh
│   └── validate-results.sh
├── templates/              # Test file templates
│   ├── dummy-code.js
│   ├── sample-config.yaml
│   └── batch-examples/
└── results/               # Test execution logs
    └── integration-YYYYMMDD-HHMMSS.log
```

### Test Runner Features

**Automatic PR Management:**
- Generate unique branch names with timestamps
- Create minimal, realistic code changes
- Auto-detect PR numbers for commands
- Cleanup branches and PRs after tests

**Result Validation:**
- Cross-validate results using multiple methods
- Log all commands and outputs
- Generate test reports with pass/fail status
- Screenshot/export functionality for manual inspection

**Safety Features:**
- Confirmation prompts for destructive operations
- Dry-run mode for testing without side effects
- Rate limiting to respect GitHub API limits
- Rollback capabilities for failed tests

## Advanced Features

### Conditional Execution
```bash
# Run integration tests every 10th execution
EXECUTION_COUNT=$(cat .execution-count 2>/dev/null || echo 0)
if (( (EXECUTION_COUNT + 1) % 10 == 0 )); then
    go run . test-integration --cleanup
fi
echo $((EXECUTION_COUNT + 1)) > .execution-count
```

### Environment Configuration
```bash
# Environment variables for test configuration
export GH_COMMENT_INTEGRATION_REPO="silouan.wright/gh-comment"
export GH_COMMENT_INTEGRATION_BRANCH_PREFIX="integration-test"
export GH_COMMENT_INTEGRATION_CLEANUP="true"
export GH_COMMENT_INTEGRATION_LOG_LEVEL="debug"
```

### Parallel Test Execution
```bash
# Run multiple scenarios in parallel
./runner.go --scenario=comments &
./runner.go --scenario=reviews &  
./runner.go --scenario=reactions &
wait
echo "All integration tests completed"
```

## Success Criteria

**Functional Validation:**
- ✅ All commands execute without errors
- ✅ Comments appear correctly in GitHub UI
- ✅ List command shows all created comments
- ✅ Reactions and replies function properly
- ✅ Batch operations process correctly

**Integration Validation:**
- ✅ Real GitHub API interactions work
- ✅ Authentication flows properly
- ✅ Repository and PR detection works
- ✅ Error handling for API failures
- ✅ Rate limiting respected

**Cleanup Validation:**
- ✅ Test PRs are properly closed/deleted
- ✅ No orphaned branches remain
- ✅ No test artifacts pollute repository
- ✅ Clean test environment for next run

## Implementation Timeline

**Phase 1: Basic Framework (1-2 hours)**
- [ ] Create test runner skeleton
- [ ] Implement PR creation/cleanup scripts
- [ ] Add basic comment scenario
- [ ] Test with manual cleanup

**Phase 2: Full Scenarios (2-3 hours)**
- [ ] Implement all 5 test scenarios
- [ ] Add comprehensive validation
- [ ] Create template files and configs
- [ ] Test all scenarios end-to-end

**Phase 3: Automation & Polish (1-2 hours)**
- [ ] Add conditional execution logic
- [ ] Implement parallel test support
- [ ] Create detailed logging and reporting
- [ ] Add safety features and error handling

**Total Estimated Time: 4-7 hours**

This dogfooding approach will provide the highest confidence that `gh-comment` works correctly with real GitHub APIs while leveraging the project's own infrastructure for testing.