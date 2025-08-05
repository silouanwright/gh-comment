# Examples for command: list
# Generated from help text on Tue Aug  5 07:42:32 CDT 2025

# Example 1 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "senior-dev*" --status open --since "1 week ago""
./gh-comment list $PR_NUM --author "senior-dev*" --status open --since "1 week ago"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 1"
    echo "list:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 1"
    echo "list:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --type review --author "*@company.com" --since "2024-01-01""
./gh-comment list $PR_NUM --type review --author "*@company.com" --since "2024-01-01"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 2"
    echo "list:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 2"
    echo "list:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "security-team*" --since "2024-01-01" --type review"
./gh-comment list $PR_NUM --author "security-team*" --since "2024-01-01" --type review
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 3"
    echo "list:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 3"
    echo "list:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "bot*" --since "3 days ago" --quiet"
./gh-comment list $PR_NUM --author "bot*" --since "3 days ago" --quiet
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 4"
    echo "list:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 4"
    echo "list:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 5 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --format json | jq '.comments[].id'"
./gh-comment list $PR_NUM --format json | jq '.comments[].id'
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 5"
    echo "list:example_5:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 5"
    echo "list:example_5:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 6 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --ids-only | xargs -I {} ./gh-comment resolve {}"
./gh-comment list $PR_NUM --ids-only | xargs -I {} ./gh-comment resolve {}
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 6"
    echo "list:example_6:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 6"
    echo "list:example_6:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 7 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --format json --author "security*" > security-comments.json"
./gh-comment list $PR_NUM --format json --author "security*" > security-comments.json
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 7"
    echo "list:example_7:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 7"
    echo "list:example_7:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 8 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --status open --since "1 month ago" --author "lead*""
./gh-comment list $PR_NUM --status open --since "1 month ago" --author "lead*"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 8"
    echo "list:example_8:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 8"
    echo "list:example_8:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 9 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --until "2024-12-31" --type issue --status resolved"
./gh-comment list $PR_NUM --until "2024-12-31" --type issue --status resolved
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 9"
    echo "list:example_9:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 9"
    echo "list:example_9:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 10 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "qa*" --since "3 days ago" --type review"
./gh-comment list $PR_NUM --author "qa*" --since "3 days ago" --type review
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 10"
    echo "list:example_10:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 10"
    echo "list:example_10:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 11 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "*@contractor.com" --status open --since "1 month ago""
./gh-comment list $PR_NUM --author "*@contractor.com" --status open --since "1 month ago"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 11"
    echo "list:example_11:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 11"
    echo "list:example_11:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 12 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "architect*" --status open --type review"
./gh-comment list $PR_NUM --author "architect*" --status open --type review
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 12"
    echo "list:example_12:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 12"
    echo "list:example_12:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 13 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --since "critical-bug-report" --author "oncall*" --status resolved"
./gh-comment list $PR_NUM --since "critical-bug-report" --author "oncall*" --status resolved
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 13"
    echo "list:example_13:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 13"
    echo "list:example_13:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 14 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "performance-team" --since "load-test-date" --type review"
./gh-comment list $PR_NUM --author "performance-team" --since "load-test-date" --type review
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 14"
    echo "list:example_14:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 14"
    echo "list:example_14:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 15 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --status open --author "*perf*" --since "1 week ago""
./gh-comment list $PR_NUM --status open --author "*perf*" --since "1 week ago"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 15"
    echo "list:example_15:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 15"
    echo "list:example_15:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 16 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "all-reviewers*" --since "quarter-start" --quiet | process-review-data.sh"
./gh-comment list $PR_NUM --author "all-reviewers*" --since "quarter-start" --quiet | process-review-data.sh
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 16"
    echo "list:example_16:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 16"
    echo "list:example_16:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 17 from list help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --ids-only --type review --status open | review-metrics.sh"
./gh-comment list $PR_NUM --ids-only --type review --status open | review-metrics.sh
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: list example 17"
    echo "list:example_17:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: list example 17"
    echo "list:example_17:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

