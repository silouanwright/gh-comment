# Examples for command: react
# Generated from help text on Tue Aug  5 07:42:34 CDT 2025

# Example 1 from react help text
echo "🧪 Testing: ./gh-comment react $PR_NUM456 +1"
./gh-comment react $PR_NUM456 +1
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 1"
    echo "react:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 1"
    echo "react:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from react help text
echo "🧪 Testing: ./gh-comment react 789012 heart"
./gh-comment react 789012 heart
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 2"
    echo "react:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 2"
    echo "react:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from react help text
echo "🧪 Testing: ./gh-comment react 999999 rocket"
./gh-comment react 999999 rocket
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 3"
    echo "react:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 3"
    echo "react:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from react help text
echo "🧪 Testing: ./gh-comment react $PR_NUM456 +1 --remove"
./gh-comment react $PR_NUM456 +1 --remove
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 4"
    echo "react:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 4"
    echo "react:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 5 from react help text
echo "🧪 Testing: ./gh-comment react 789012 heart --remove"
./gh-comment react 789012 heart --remove
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 5"
    echo "react:example_5:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 5"
    echo "react:example_5:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 6 from react help text
echo "🧪 Testing: ./gh-comment react $PR_NUM456 eyes    # Works for issue comments"
./gh-comment react $PR_NUM456 eyes    # Works for issue comments
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 6"
    echo "react:example_6:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 6"
    echo "react:example_6:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 7 from react help text
echo "🧪 Testing: ./gh-comment react 789012 hooray  # Works for review comments"
./gh-comment react 789012 hooray  # Works for review comments
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: react example 7"
    echo "react:example_7:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: react example 7"
    echo "react:example_7:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

