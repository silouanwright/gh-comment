# Examples for command: add
# Generated from help text on Tue Aug  5 07:30:41 CDT 2025

# Example 1 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "LGTM! Just a few minor suggestions below""
./gh-comment add $PR_NUM "LGTM! Just a few minor suggestions below"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 1"
    echo "add:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 1"
    echo "add:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Thanks for addressing the security concerns""
./gh-comment add $PR_NUM "Thanks for addressing the security concerns"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 2"
    echo "add:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 2"
    echo "add:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "This looks great - ready to merge after CI passes""
./gh-comment add $PR_NUM "This looks great - ready to merge after CI passes"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 3"
    echo "add:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 3"
    echo "add:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM -m "Overall this is excellent work!" -m "The architecture is clean and the tests are comprehensive""
./gh-comment add $PR_NUM -m "Overall this is excellent work!" -m "The architecture is clean and the tests are comprehensive"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 4"
    echo "add:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 4"
    echo "add:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 5 from add help text
echo "🧪 Testing: ./gh-comment add "Looks good to merge!""
./gh-comment add "Looks good to merge!"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 5"
    echo "add:example_5:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 5"
    echo "add:example_5:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 6 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Approved! The performance improvements in this PR will make a huge difference""
./gh-comment add $PR_NUM "Approved! The performance improvements in this PR will make a huge difference"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 6"
    echo "add:example_6:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 6"
    echo "add:example_6:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 7 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Could you address the failing tests? Otherwise looks good to go""
./gh-comment add $PR_NUM "Could you address the failing tests? Otherwise looks good to go"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 7"
    echo "add:example_7:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 7"
    echo "add:example_7:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 8 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Consider using async/await: [SUGGEST: const result = await fetchData();]""
./gh-comment add $PR_NUM "Consider using async/await: [SUGGEST: const result = await fetchData();]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 8"
    echo "add:example_8:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 8"
    echo "add:example_8:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 9 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Add error handling above: [SUGGEST:-1: try {]""
./gh-comment add $PR_NUM "Add error handling above: [SUGGEST:-1: try {]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 9"
    echo "add:example_9:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 9"
    echo "add:example_9:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 10 from add help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Add timeout below: [SUGGEST:+2: const timeout = 5000;]""
./gh-comment add $PR_NUM "Add timeout below: [SUGGEST:+2: const timeout = 5000;]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: add example 10"
    echo "add:example_10:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: add example 10"
    echo "add:example_10:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

