# Examples for command: resolve
# Generated from help text on Tue Aug  5 07:42:34 CDT 2025

# Example 1 from resolve help text
echo "🧪 Testing: ./gh-comment resolve 2246362251"
./gh-comment resolve 2246362251
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: resolve example 1"
    echo "resolve:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: resolve example 1"
    echo "resolve:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from resolve help text
echo "🧪 Testing: ./gh-comment resolve --dry-run 2246362251"
./gh-comment resolve --dry-run 2246362251
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: resolve example 2"
    echo "resolve:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: resolve example 2"
    echo "resolve:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from resolve help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --status open --ids-only | xargs -I {} ./gh-comment resolve {}"
./gh-comment list $PR_NUM --status open --ids-only | xargs -I {} ./gh-comment resolve {}
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: resolve example 3"
    echo "resolve:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: resolve example 3"
    echo "resolve:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

