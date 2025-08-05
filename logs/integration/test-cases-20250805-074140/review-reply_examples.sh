# Examples for command: review-reply
# Generated from help text on Tue Aug  5 07:42:35 CDT 2025

# Example 1 from review-reply help text
echo "🧪 Testing: ./gh-comment review-reply 789012 "Fixed this issue""
./gh-comment review-reply 789012 "Fixed this issue"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 1"
    echo "review-reply:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 1"
    echo "review-reply:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from review-reply help text
echo "🧪 Testing: ./gh-comment review-reply 789012 "Addressed your feedback" --resolve"
./gh-comment review-reply 789012 "Addressed your feedback" --resolve
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 2"
    echo "review-reply:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 2"
    echo "review-reply:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from review-reply help text
echo "🧪 Testing: ./gh-comment review-reply 789012 --resolve"
./gh-comment review-reply 789012 --resolve
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 3"
    echo "review-reply:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 3"
    echo "review-reply:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

