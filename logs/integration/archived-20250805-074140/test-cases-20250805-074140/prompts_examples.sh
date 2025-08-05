# Examples for command: prompts
# Generated from help text on Tue Aug  5 07:42:33 CDT 2025

# Example 1 from prompts help text
echo "🧪 Testing: ./gh-comment prompts --list"
./gh-comment prompts --list
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: prompts example 1"
    echo "prompts:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: prompts example 1"
    echo "prompts:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from prompts help text
echo "🧪 Testing: ./gh-comment prompts security-audit"
./gh-comment prompts security-audit
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: prompts example 2"
    echo "prompts:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: prompts example 2"
    echo "prompts:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from prompts help text
echo "🧪 Testing: ./gh-comment prompts --category performance --list"
./gh-comment prompts --category performance --list
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: prompts example 3"
    echo "prompts:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: prompts example 3"
    echo "prompts:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from prompts help text
echo "🧪 Testing: ./gh-comment prompts performance-optimization"
./gh-comment prompts performance-optimization
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: prompts example 4"
    echo "prompts:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: prompts example 4"
    echo "prompts:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 5 from prompts help text
echo "🧪 Testing: ./gh-comment prompts security-audit"
./gh-comment prompts security-audit
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: prompts example 5"
    echo "prompts:example_5:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: prompts example 5"
    echo "prompts:example_5:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 6 from prompts help text
echo "🧪 Testing: ./gh-comment prompts architecture-review"
./gh-comment prompts architecture-review
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: prompts example 6"
    echo "prompts:example_6:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: prompts example 6"
    echo "prompts:example_6:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

