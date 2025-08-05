# Examples for command: batch
# Generated from help text on Tue Aug  5 07:42:39 CDT 2025

# Example 1 from batch help text
echo "🧪 Testing: ./gh-comment batch $PR_NUM review-config.yaml"
./gh-comment batch $PR_NUM review-config.yaml
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: batch example 1"
    echo "batch:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: batch example 1"
    echo "batch:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from batch help text
echo "🧪 Testing: ./gh-comment batch $PR_NUM review-config.yaml --dry-run"
./gh-comment batch $PR_NUM review-config.yaml --dry-run
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: batch example 2"
    echo "batch:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: batch example 2"
    echo "batch:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from batch help text
echo "🧪 Testing: ./gh-comment batch $PR_NUM review-config.yaml --verbose"
./gh-comment batch $PR_NUM review-config.yaml --verbose
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: batch example 3"
    echo "batch:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: batch example 3"
    echo "batch:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

