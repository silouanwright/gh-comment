# Examples for command: edit
# Generated from help text on Tue Aug  5 07:30:35 CDT 2025

# Example 1 from edit help text
echo "🧪 Testing: ./gh-comment edit 2246362251 "Updated comment with better explanation""
./gh-comment edit 2246362251 "Updated comment with better explanation"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: edit example 1"
    echo "edit:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: edit example 1"
    echo "edit:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from edit help text
echo "🧪 Testing: ./gh-comment edit 2246362251 --message "First paragraph" --message "Second paragraph""
./gh-comment edit 2246362251 --message "First paragraph" --message "Second paragraph"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: edit example 2"
    echo "edit:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: edit example 2"
    echo "edit:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from edit help text
echo "🧪 Testing: ./gh-comment edit 2246362251 "Line 1"
./gh-comment edit 2246362251 "Line 1
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: edit example 3"
    echo "edit:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: edit example 3"
    echo "edit:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

