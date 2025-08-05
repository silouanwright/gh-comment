# Examples for command: close-pending-review
# Generated from help text on Tue Aug  5 07:42:39 CDT 2025

# Example 1 from close-pending-review help text
echo "🧪 Testing: ./gh-comment close-pending-review $PR_NUM "LGTM! Great work" --event APPROVE"
./gh-comment close-pending-review $PR_NUM "LGTM! Great work" --event APPROVE
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: close-pending-review example 1"
    echo "close-pending-review:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: close-pending-review example 1"
    echo "close-pending-review:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from close-pending-review help text
echo "🧪 Testing: ./gh-comment close-pending-review $PR_NUM "Please address the comments" --event REQUEST_CHANGES"
./gh-comment close-pending-review $PR_NUM "Please address the comments" --event REQUEST_CHANGES
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: close-pending-review example 2"
    echo "close-pending-review:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: close-pending-review example 2"
    echo "close-pending-review:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from close-pending-review help text
echo "🧪 Testing: ./gh-comment close-pending-review $PR_NUM "Thanks for the updates" --event COMMENT"
./gh-comment close-pending-review $PR_NUM "Thanks for the updates" --event COMMENT
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: close-pending-review example 3"
    echo "close-pending-review:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: close-pending-review example 3"
    echo "close-pending-review:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from close-pending-review help text
echo "🧪 Testing: ./gh-comment close-pending-review "Looks good" --event APPROVE"
./gh-comment close-pending-review "Looks good" --event APPROVE
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: close-pending-review example 4"
    echo "close-pending-review:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: close-pending-review example 4"
    echo "close-pending-review:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

