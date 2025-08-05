# Examples for command: review
# Generated from help text on Tue Aug  5 16:37:27 CDT 2025

# Example 1 from review help text
echo "🧪 Testing: ./gh-comment review $PR_NUM "Security audit complete - critical issues found" \"
./gh-comment review $PR_NUM "Security audit complete - critical issues found" \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review example 1"
    echo "review:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review example 1"
    echo "review:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from review help text
echo "🧪 Testing: ./gh-comment review $PR_NUM \"
./gh-comment review $PR_NUM \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review example 2"
    echo "review:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review example 2"
    echo "review:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from review help text
echo "🧪 Testing: ./gh-comment review $PR_NUM "Migration to microservices architecture approved" \"
./gh-comment review $PR_NUM "Migration to microservices architecture approved" \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review example 3"
    echo "review:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review example 3"
    echo "review:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from review help text
echo "🧪 Testing: ./gh-comment review $PR_NUM "Code quality review - refactoring needed" \"
./gh-comment review $PR_NUM "Code quality review - refactoring needed" \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review example 4"
    echo "review:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review example 4"
    echo "review:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""
