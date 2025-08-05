#!/bin/bash
set -euo pipefail

# 03_run_tests.sh - Execute parsed help text examples
# Runs all extracted test cases and collects results

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs/integration"

# Find the latest test environment
LATEST_ENV=$(ls -t "$LOG_DIR"/test-env-*.sh 2>/dev/null | head -1 || echo "")

if [[ -z "$LATEST_ENV" ]]; then
    echo "❌ ERROR: No test environment found"
    echo "   Run scripts/integration/01_setup.sh first"
    exit 1
fi

echo "🧪 Running integration tests..."
echo "🔗 Using environment: $(basename "$LATEST_ENV")"

# Source test environment
source "$LATEST_ENV"

cd "$PROJECT_ROOT"

# Verify we're on the test branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$CURRENT_BRANCH" != "$TEST_BRANCH" ]]; then
    echo "❌ ERROR: Not on test branch"
    echo "   Current: $CURRENT_BRANCH"
    echo "   Expected: $TEST_BRANCH"
    exit 1
fi

# Ensure test cases exist
if [[ ! -d "$TEST_CASES_DIR" ]]; then
    echo "❌ ERROR: Test cases not found"
    echo "   Run scripts/integration/02_parse_help.sh first"
    exit 1
fi

# Ensure binary exists and is current
if [[ ! -f "./gh-comment" ]] || [[ "$PROJECT_ROOT/main.go" -nt "./gh-comment" ]]; then
    echo "🔨 Rebuilding binary..."
    go build -o gh-comment .
fi

echo "🚀 Starting test execution..."
echo "📊 Results will be saved to: $TEST_RESULTS_FILE"
echo ""

# Initialize results file
cat > "$TEST_RESULTS_FILE" << EOF
# Integration Test Results
# Generated: $(date)
# Branch: $TEST_BRANCH
# PR: #$PR_NUM ($PR_URL)
# Binary: $(./gh-comment --version 2>/dev/null || echo "unknown version")
#
# Format: command:example_id:status:details
#
EOF

# We need to get a real comment ID for tests that require it
echo "🔍 Setting up test data (getting comment ID)..."

# Add a test comment to get a real comment ID
TEST_COMMENT_OUTPUT=$(./gh-comment add "$PR_NUM" "Integration test marker comment - $(date)" 2>&1 || echo "FAILED")
if [[ "$TEST_COMMENT_OUTPUT" =~ ID:[[:space:]]*([0-9]+) ]]; then
    COMMENT_ID="${BASH_REMATCH[1]}"
    echo "✅ Created test comment with ID: $COMMENT_ID"
else
    echo "⚠️  WARNING: Could not create test comment for ID-dependent tests"
    COMMENT_ID="999999999"  # Fallback ID that will likely fail gracefully
fi

# Export for test scripts
export COMMENT_ID

echo ""
echo "🧪 Running test cases..."

# Function to run a single test script
run_test_script() {
    local script="$1"
    local cmd_name=$(basename "$script" "_examples.sh")

    echo "▶️  Testing command: $cmd_name"

    # Create a temporary log for this command
    local cmd_log="$LOG_DIR/cmd-${cmd_name}-$TIMESTAMP.log"

    # Run the test script and capture output
    if bash "$script" > "$cmd_log" 2>&1; then
        echo "  ✅ All examples for $cmd_name completed"

        # Show any failures from the command log
        if grep -q "❌ FAILED" "$cmd_log"; then
            echo "  ⚠️  Some examples failed - check $cmd_log"
        fi
    else
        echo "  ❌ Test script failed for $cmd_name"
        echo "  📄 Log: $cmd_log"
    fi

    # Extract the last few lines for summary
    echo "  📝 $(tail -3 "$cmd_log" | head -1)"
}

# Find and run all test scripts
TEST_SCRIPTS=($(find "$TEST_CASES_DIR" -name "*_examples.sh" | sort))

if [[ ${#TEST_SCRIPTS[@]} -eq 0 ]]; then
    echo "❌ ERROR: No test scripts found in $TEST_CASES_DIR"
    exit 1
fi

echo "📋 Found ${#TEST_SCRIPTS[@]} test script(s) to run"
echo ""

# Run each test script
for script in "${TEST_SCRIPTS[@]}"; do
    run_test_script "$script"
done

echo ""
echo "📊 Test execution complete!"

# Generate summary report
echo "📋 Generating summary report..."

# Count results
TOTAL_TESTS=$(grep -E "^[^#].*:.*:" "$TEST_RESULTS_FILE" | wc -l | tr -d ' ')
SUCCESS_COUNT=$(grep -c ":SUCCESS:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
FAILED_COUNT=$(grep -c ":FAILED:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
WARNING_COUNT=$(grep -c ":WARNING:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")

# Calculate success rate
if [[ "$TOTAL_TESTS" -gt 0 ]]; then
    SUCCESS_RATE=$(( (SUCCESS_COUNT * 100) / TOTAL_TESTS ))
else
    SUCCESS_RATE=0
fi

# Create summary report
SUMMARY_FILE="$LOG_DIR/summary-$TIMESTAMP.md"
cat > "$SUMMARY_FILE" << EOF
# Integration Test Summary

**Generated:** $(date)
**Branch:** $TEST_BRANCH
**PR:** [#$PR_NUM]($PR_URL)
**Binary Version:** $(./gh-comment --version 2>/dev/null || echo "unknown")

## Results Overview

| Metric | Count |
|--------|-------|
| Total Tests | $TOTAL_TESTS |
| Successful | $SUCCESS_COUNT |
| Failed | $FAILED_COUNT |
| Warnings | $WARNING_COUNT |
| **Success Rate** | **${SUCCESS_RATE}%** |

EOF

# Add failed tests section if any
if [[ "$FAILED_COUNT" -gt 0 ]]; then
    echo "## ❌ Failed Tests" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
    grep ":FAILED:" "$TEST_RESULTS_FILE" | while IFS=':' read -r cmd example status details; do
        echo "- **$cmd**: $example - $details" >> "$SUMMARY_FILE"
    done
    echo "" >> "$SUMMARY_FILE"
fi

# Add warnings section if any
if [[ "$WARNING_COUNT" -gt 0 ]]; then
    echo "## ⚠️ Warnings" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
    grep ":WARNING:" "$TEST_RESULTS_FILE" | while IFS=':' read -r cmd example status details; do
        echo "- **$cmd**: $example - $details" >> "$SUMMARY_FILE"
    done
    echo "" >> "$SUMMARY_FILE"
fi

# Add commands tested section
echo "## 📋 Commands Tested" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
for script in "${TEST_SCRIPTS[@]}"; do
    cmd_name=$(basename "$script" "_examples.sh")
    cmd_success=$(grep "^$cmd_name:.*:SUCCESS:" "$TEST_RESULTS_FILE" | wc -l | tr -d ' ')
    cmd_failed=$(grep "^$cmd_name:.*:FAILED:" "$TEST_RESULTS_FILE" | wc -l | tr -d ' ')
    cmd_warnings=$(grep "^$cmd_name:.*:WARNING:" "$TEST_RESULTS_FILE" | wc -l | tr -d ' ')

    status_icon="✅"
    if [[ "$cmd_failed" -gt 0 ]]; then
        status_icon="❌"
    elif [[ "$cmd_warnings" -gt 0 ]]; then
        status_icon="⚠️"
    fi

    echo "- $status_icon **$cmd_name**: $cmd_success success, $cmd_failed failed, $cmd_warnings warnings" >> "$SUMMARY_FILE"
done

echo "" >> "$SUMMARY_FILE"
echo "## 📄 Detailed Results" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "- **Full Results:** \`$(basename "$TEST_RESULTS_FILE")\`" >> "$SUMMARY_FILE"
echo "- **Test Cases:** \`$(basename "$TEST_CASES_DIR")/\`" >> "$SUMMARY_FILE"
echo "- **Individual Logs:** \`logs/integration/cmd-*-$TIMESTAMP.log\`" >> "$SUMMARY_FILE"

# Display summary
echo ""
echo "🎯 INTEGRATION TEST SUMMARY"
echo "=========================="
echo "📊 Results: $SUCCESS_COUNT/$TOTAL_TESTS passed (${SUCCESS_RATE}%)"
if [[ "$FAILED_COUNT" -gt 0 ]]; then
    echo "❌ Failed: $FAILED_COUNT tests"
fi
if [[ "$WARNING_COUNT" -gt 0 ]]; then
    echo "⚠️  Warnings: $WARNING_COUNT tests"
fi
echo ""
echo "📄 Reports:"
echo "  • Summary: $SUMMARY_FILE"
echo "  • Detailed: $TEST_RESULTS_FILE"
echo "  • Test cases: $TEST_CASES_DIR"
echo ""

# Cleanup test comment
if [[ "$COMMENT_ID" != "999999999" ]]; then
    echo "🧹 Cleaning up test comment #$COMMENT_ID..."
    ./gh-comment edit "$COMMENT_ID" "Integration test completed - $(date)" >/dev/null 2>&1 || true
fi

if [[ "$SUCCESS_RATE" -lt 100 ]]; then
    echo "⚠️ Some tests failed. Check the reports above for details."
    echo ""
    echo "Next steps:"
    echo "  1. Review failed tests in: $SUMMARY_FILE"
    echo "  2. Fix help text or implementation issues"
    echo "  3. Re-run tests: $0"
    echo "  4. When ready: scripts/integration/04_cleanup.sh"
    exit 1
else
    echo "🎉 All tests passed! Help text examples are working correctly."
    echo ""
    echo "Next steps:"
    echo "  1. Review summary: $SUMMARY_FILE"
    echo "  2. Clean up: scripts/integration/04_cleanup.sh"
fi
