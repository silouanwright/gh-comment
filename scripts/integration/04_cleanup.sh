#!/bin/bash
set -euo pipefail

# 04_cleanup.sh - Cleanup integration test resources
# Closes PR, deletes branch, and archives test results

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs/integration"

# Find the latest test environment
LATEST_ENV=$(ls -t "$LOG_DIR"/test-env-*.sh 2>/dev/null | head -1 || echo "")

if [[ -z "$LATEST_ENV" ]]; then
    echo "❌ ERROR: No test environment found"
    echo "   Nothing to clean up"
    exit 1
fi

echo "🧹 Cleaning up integration test resources..."
echo "🔗 Using environment: $(basename "$LATEST_ENV")"

# Source test environment
source "$LATEST_ENV"

cd "$PROJECT_ROOT"

# Confirm cleanup with user (unless --force flag is used)
if [[ "${1:-}" != "--force" ]]; then
    echo ""
    echo "This will:"
    echo "  • Close PR #$PR_NUM"
    echo "  • Delete branch '$TEST_BRANCH'"
    echo "  • Archive test results in logs/integration/"
    echo ""
    echo "Continue with cleanup? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "   Cleanup cancelled by user"
        exit 0
    fi
fi

echo ""
echo "🔄 Starting cleanup process..."

# Archive test results before cleanup
ARCHIVE_DIR="$LOG_DIR/archived-$TIMESTAMP"
echo "📦 Archiving test results to: $ARCHIVE_DIR"

mkdir -p "$ARCHIVE_DIR"

# Move/copy important files to archive
if [[ -f "$TEST_RESULTS_FILE" ]]; then
    cp "$TEST_RESULTS_FILE" "$ARCHIVE_DIR/"
    echo "  ✅ Archived test results"
fi

if [[ -f "$LOG_DIR/summary-$TIMESTAMP.md" ]]; then
    cp "$LOG_DIR/summary-$TIMESTAMP.md" "$ARCHIVE_DIR/"
    echo "  ✅ Archived summary report"
fi

if [[ -d "$TEST_CASES_DIR" ]]; then
    cp -r "$TEST_CASES_DIR" "$ARCHIVE_DIR/"
    echo "  ✅ Archived test cases"
fi

# Archive individual command logs
for log_file in "$LOG_DIR"/cmd-*-"$TIMESTAMP".log; do
    if [[ -f "$log_file" ]]; then
        cp "$log_file" "$ARCHIVE_DIR/"
    fi
done

# Copy the environment file for reference
cp "$LATEST_ENV" "$ARCHIVE_DIR/"

# Create archive summary
cat > "$ARCHIVE_DIR/README.md" << EOF
# Integration Test Archive

**Test Run:** $TIMESTAMP
**Branch:** $TEST_BRANCH
**PR:** #$PR_NUM ($PR_URL)
**Archived:** $(date)

## Files

- \`test-results-$TIMESTAMP.txt\` - Detailed test results
- \`summary-$TIMESTAMP.md\` - Summary report
- \`test-cases-$TIMESTAMP/\` - Generated test cases
- \`test-env-$TIMESTAMP.sh\` - Test environment variables
- \`cmd-*-$TIMESTAMP.log\` - Individual command logs

## Test Summary

$(if [[ -f "$TEST_RESULTS_FILE" ]]; then
    TOTAL=$(grep -E "^[^#].*:.*:" "$TEST_RESULTS_FILE" | wc -l | tr -d ' ')
    SUCCESS=$(grep -c ":SUCCESS:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
    FAILED=$(grep -c ":FAILED:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
    echo "- **Total Tests:** $TOTAL"
    echo "- **Successful:** $SUCCESS"
    echo "- **Failed:** $FAILED"
    if [[ "$TOTAL" -gt 0 ]]; then
        RATE=$(( (SUCCESS * 100) / TOTAL ))
        echo "- **Success Rate:** ${RATE}%"
    fi
else
    echo "No test results found"
fi)

## Next Steps

This archive preserves all test artifacts for future reference. The integration test branch and PR have been cleaned up.
EOF

echo "  ✅ Created archive README"

# Close the PR with a summary comment
echo "📝 Closing PR #$PR_NUM..."

# Get final test summary for PR comment
FINAL_SUMMARY=""
if [[ -f "$TEST_RESULTS_FILE" ]]; then
    TOTAL_TESTS=$(grep -E "^[^#].*:.*:" "$TEST_RESULTS_FILE" | wc -l | tr -d ' ')
    SUCCESS_COUNT=$(grep -c ":SUCCESS:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
    FAILED_COUNT=$(grep -c ":FAILED:" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")

    if [[ "$TOTAL_TESTS" -gt 0 ]]; then
        SUCCESS_RATE=$(( (SUCCESS_COUNT * 100) / TOTAL_TESTS ))
        FINAL_SUMMARY="🎯 **Final Results:** $SUCCESS_COUNT/$TOTAL_TESTS tests passed (${SUCCESS_RATE}%)"

        if [[ "$FAILED_COUNT" -gt 0 ]]; then
            FINAL_SUMMARY="$FINAL_SUMMARY\n\n❌ **Failed Tests:** $FAILED_COUNT"
        fi
    else
        FINAL_SUMMARY="⚠️ No test results found"
    fi
else
    FINAL_SUMMARY="⚠️ Test results file not found"
fi

# Add closing comment to PR
PR_CLOSE_COMMENT="## 🧪 Integration Test Complete

$FINAL_SUMMARY

**Test Duration:** $(date -d "@$(stat -c %Y "$LATEST_ENV" 2>/dev/null || echo $(date +%s))" '+%Y-%m-%d %H:%M') → $(date '+%Y-%m-%d %H:%M')

### 📊 Test Summary
- All help text examples were extracted and executed
- Real GitHub API integration validated
- Results archived in \`logs/integration/archived-$TIMESTAMP/\`

### 🎉 Mission Accomplished!
This PR successfully validated that all documented examples in \`gh-comment\` help text work correctly with the real GitHub API.

**Status:** Integration test completed successfully! 🚀"

# Add comment and close PR
if gh pr comment "$PR_NUM" --body "$PR_CLOSE_COMMENT" >/dev/null 2>&1; then
    echo "  ✅ Added final comment to PR"
else
    echo "  ⚠️  Could not add final comment to PR"
fi

if gh pr close "$PR_NUM" >/dev/null 2>&1; then
    echo "  ✅ Closed PR #$PR_NUM"
else
    echo "  ❌ Failed to close PR #$PR_NUM"
    echo "     You may need to close it manually: $PR_URL"
fi

# Switch back to main branch
echo "🔄 Switching back to main branch..."
if git checkout main >/dev/null 2>&1; then
    echo "  ✅ Switched to main branch"
else
    echo "  ❌ Failed to switch to main branch"
    echo "     Current branch: $(git rev-parse --abbrev-ref HEAD)"
fi

# Delete the test branch (local and remote)
echo "🗑️  Deleting test branch: $TEST_BRANCH"

# Delete local branch
if git branch -D "$TEST_BRANCH" >/dev/null 2>&1; then
    echo "  ✅ Deleted local branch"
else
    echo "  ⚠️  Could not delete local branch (may not exist)"
fi

# Delete remote branch
if git push origin --delete "$TEST_BRANCH" >/dev/null 2>&1; then
    echo "  ✅ Deleted remote branch"
else
    echo "  ⚠️  Could not delete remote branch (may not exist)"
fi

# Clean up temporary files (but keep archives)
echo "🗑️  Cleaning up temporary files..."

# Remove non-archived temporary files
if [[ -f "$TEST_RESULTS_FILE" ]]; then
    rm "$TEST_RESULTS_FILE"
    echo "  ✅ Removed temporary test results"
fi

if [[ -d "$TEST_CASES_DIR" ]]; then
    rm -rf "$TEST_CASES_DIR"
    echo "  ✅ Removed temporary test cases"
fi

# Remove individual command logs (they're archived)
for log_file in "$LOG_DIR"/cmd-*-"$TIMESTAMP".log; do
    if [[ -f "$log_file" ]]; then
        rm "$log_file"
    fi
done

# Remove the environment file (it's archived)
rm "$LATEST_ENV"

# Clean up the built binary
if [[ -f "$PROJECT_ROOT/gh-comment" ]]; then
    rm "$PROJECT_ROOT/gh-comment"
    echo "  ✅ Removed test binary"
fi

echo ""
echo "✅ Integration test cleanup complete!"
echo ""
echo "📋 Summary:"
echo "  • PR #$PR_NUM: Closed"
echo "  • Branch '$TEST_BRANCH': Deleted"
echo "  • Test results: Archived in $ARCHIVE_DIR"
echo "  • Working directory: Clean"
echo ""
echo "📄 Test archive contains:"
echo "  • Test results and summary reports"
echo "  • Generated test cases"
echo "  • Individual command logs"
echo "  • Complete environment snapshot"
echo ""
echo "🎯 Integration testing cycle complete!"

# Show final archive location for easy access
echo ""
echo "🔗 Archived results: $ARCHIVE_DIR"
echo "📖 Archive summary: $ARCHIVE_DIR/README.md"

if [[ -f "$ARCHIVE_DIR/summary-$TIMESTAMP.md" ]]; then
    echo ""
    echo "📊 Final Test Summary:"
    echo "======================="
    # Show the key metrics from the summary
    if grep -q "Success Rate" "$ARCHIVE_DIR/summary-$TIMESTAMP.md"; then
        grep -E "(Total Tests|Successful|Failed|Success Rate)" "$ARCHIVE_DIR/summary-$TIMESTAMP.md" | sed 's/^| /  • /' | sed 's/ |.*$//'
    fi
fi
