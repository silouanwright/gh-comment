# Examples for command: lines
# Generated from help text on Tue Aug  5 07:42:32 CDT 2025

# Example 1 from lines help text
echo "🧪 Testing: ./gh-comment lines $PR_NUM src/src/main.go"
./gh-comment lines $PR_NUM src/src/main.go
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: lines example 1"
    echo "lines:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: lines example 1"
    echo "lines:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from lines help text
echo "🧪 Testing: ./gh-comment lines $PR_NUM src/src/main.go --show-code"
./gh-comment lines $PR_NUM src/src/main.go --show-code
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: lines example 2"
    echo "lines:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: lines example 2"
    echo "lines:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from lines help text
echo "🧪 Testing: ./gh-comment lines $PR_NUM src/src/main.go | grep "^42:""
./gh-comment lines $PR_NUM src/src/main.go | grep "^42:"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: lines example 3"
    echo "lines:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: lines example 3"
    echo "lines:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from lines help text
echo "🧪 Testing: ./gh-comment lines $PR_NUM src/src/main.go | grep -o "^[0-9]*" | head -5"
./gh-comment lines $PR_NUM src/src/main.go | grep -o "^[0-9]*" | head -5
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: lines example 4"
    echo "lines:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: lines example 4"
    echo "lines:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

