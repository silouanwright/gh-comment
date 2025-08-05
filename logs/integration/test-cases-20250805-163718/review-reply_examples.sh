# Examples for command: review-reply
# Generated from help text on Tue Aug  5 16:37:28 CDT 2025

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
echo "🧪 Testing: ./gh-comment review-reply 789012 --resolve"
./gh-comment review-reply 789012 --resolve
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 2"
    echo "review-reply:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 2"
    echo "review-reply:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from review-reply help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/main.go 42 "Fixed this issue""
./gh-comment add $PR_NUM src/main.go 42 "Fixed this issue"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 3"
    echo "review-reply:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 3"
    echo "review-reply:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from review-reply help text
echo "🧪 Testing: ./gh-comment react 789012 +1"
./gh-comment react 789012 +1
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 4"
    echo "review-reply:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 4"
    echo "review-reply:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 5 from review-reply help text
echo "🧪 Testing: ./gh-comment react 789012 heart"
./gh-comment react 789012 heart
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 5"
    echo "review-reply:example_5:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 5"
    echo "review-reply:example_5:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 6 from review-reply help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Thanks for the review feedback!""
./gh-comment add $PR_NUM "Thanks for the review feedback!"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: review-reply example 6"
    echo "review-reply:example_6:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: review-reply example 6"
    echo "review-reply:example_6:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""
