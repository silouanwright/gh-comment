# Examples for command: -R,
# Generated from help text on Tue Aug  5 07:30:39 CDT 2025

# Example 1 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM                           List all comments on PR #$PR_NUM"
./gh-comment list $PR_NUM                           List all comments on PR #$PR_NUM
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 1"
    echo "-R,:example_1:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 1"
    echo "-R,:example_1:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 2 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Looks good overall!"     Add general discussion comment"
./gh-comment add $PR_NUM "Looks good overall!"     Add general discussion comment
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 2"
    echo "-R,:example_2:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 2"
    echo "-R,:example_2:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 3 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "LGTM! Just waiting for CI to pass""
./gh-comment add $PR_NUM "LGTM! Just waiting for CI to pass"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 3"
    echo "-R,:example_3:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 3"
    echo "-R,:example_3:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 4 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM "Thanks for addressing the security concerns""
./gh-comment add $PR_NUM "Thanks for addressing the security concerns"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 4"
    echo "-R,:example_4:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 4"
    echo "-R,:example_4:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 5 from -R, help text
echo "🧪 Testing: ./gh-comment review-reply $PR_NUM45 "Good point, I'll make those changes""
./gh-comment review-reply $PR_NUM45 "Good point, I'll make those changes"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 5"
    echo "-R,:example_5:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 5"
    echo "-R,:example_5:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 6 from -R, help text
echo "🧪 Testing: ./gh-comment review $PR_NUM "Code review complete" \"
./gh-comment review $PR_NUM "Code review complete" \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 6"
    echo "-R,:example_6:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 6"
    echo "-R,:example_6:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 7 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "senior-dev*" --status open --since "1 week ago""
./gh-comment list $PR_NUM --author "senior-dev*" --status open --since "1 week ago"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 7"
    echo "-R,:example_7:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 7"
    echo "-R,:example_7:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 8 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --type review --author "*@company.com" --since "2024-01-01""
./gh-comment list $PR_NUM --type review --author "*@company.com" --since "2024-01-01"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 8"
    echo "-R,:example_8:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 8"
    echo "-R,:example_8:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 9 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --status resolved --until "2024-01-01" --quiet"
./gh-comment list $PR_NUM --status resolved --until "2024-01-01" --quiet
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 9"
    echo "-R,:example_9:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 9"
    echo "-R,:example_9:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 10 from -R, help text
echo "🧪 Testing: ./gh-comment review $PR_NUM "Migration review complete" \"
./gh-comment review $PR_NUM "Migration review complete" \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 10"
    echo "-R,:example_10:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 10"
    echo "-R,:example_10:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 11 from -R, help text
echo "🧪 Testing: ./gh-comment review $PR_NUM "Security audit findings" \"
./gh-comment review $PR_NUM "Security audit findings" \
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 11"
    echo "-R,:example_11:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 11"
    echo "-R,:example_11:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 12 from -R, help text
echo "🧪 Testing: ./gh-comment batch $PR_NUM review-config.yaml"
./gh-comment batch $PR_NUM review-config.yaml
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 12"
    echo "-R,:example_12:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 12"
    echo "-R,:example_12:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 13 from -R, help text
echo "🧪 Testing: ./gh-comment batch 456 security-checklist.yaml --dry-run"
./gh-comment batch 456 security-checklist.yaml --dry-run
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 13"
    echo "-R,:example_13:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 13"
    echo "-R,:example_13:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 14 from -R, help text
echo "🧪 Testing: ./gh-comment batch 789 bulk-comments.yaml --verbose"
./gh-comment batch 789 bulk-comments.yaml --verbose
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 14"
    echo "-R,:example_14:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 14"
    echo "-R,:example_14:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 15 from -R, help text
echo "🧪 Testing: ./gh-comment review-reply 2246362251 "Fixed in commit abc$PR_NUM" --resolve"
./gh-comment review-reply 2246362251 "Fixed in commit abc$PR_NUM" --resolve
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 15"
    echo "-R,:example_15:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 15"
    echo "-R,:example_15:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 16 from -R, help text
echo "🧪 Testing: ./gh-comment react 3141344022 +1"
./gh-comment react 3141344022 +1
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 16"
    echo "-R,:example_16:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 16"
    echo "-R,:example_16:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 17 from -R, help text
echo "🧪 Testing: ./gh-comment react 2246362251 rocket"
./gh-comment react 2246362251 rocket
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 17"
    echo "-R,:example_17:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 17"
    echo "-R,:example_17:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 18 from -R, help text
echo "🧪 Testing: ./gh-comment react 3141344022 heart --remove"
./gh-comment react 3141344022 heart --remove
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 18"
    echo "-R,:example_18:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 18"
    echo "-R,:example_18:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 19 from -R, help text
echo "🧪 Testing: ./gh-comment resolve 2246362251"
./gh-comment resolve 2246362251
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 19"
    echo "-R,:example_19:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 19"
    echo "-R,:example_19:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 20 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --quiet | grep "👤" | cut -d' ' -f2 | sort | uniq -c"
./gh-comment list $PR_NUM --quiet | grep "👤" | cut -d' ' -f2 | sort | uniq -c
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 20"
    echo "-R,:example_20:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 20"
    echo "-R,:example_20:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 21 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --since "2024-01-01" --quiet | tee q1-review-data.txt"
./gh-comment list $PR_NUM --since "2024-01-01" --quiet | tee q1-review-data.txt
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 21"
    echo "-R,:example_21:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 21"
    echo "-R,:example_21:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 22 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "qa-team*" --quiet | analyze-feedback.py"
./gh-comment list $PR_NUM --author "qa-team*" --quiet | analyze-feedback.py
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 22"
    echo "-R,:example_22:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 22"
    echo "-R,:example_22:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 23 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/security.js 67 "[SUGGEST: use crypto.randomBytes(32)]""
./gh-comment add $PR_NUM src/security.js 67 "[SUGGEST: use crypto.randomBytes(32)]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 23"
    echo "-R,:example_23:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 23"
    echo "-R,:example_23:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 24 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/src/api.js 42 "[SUGGEST:+2: const timeout = 5000;]""
./gh-comment add $PR_NUM src/src/api.js 42 "[SUGGEST:+2: const timeout = 5000;]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 24"
    echo "-R,:example_24:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 24"
    echo "-R,:example_24:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 25 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/utils.js 15 "[SUGGEST:-1: import { validateInput } from './validators';]""
./gh-comment add $PR_NUM src/utils.js 15 "[SUGGEST:-1: import { validateInput } from './validators';]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 25"
    echo "-R,:example_25:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 25"
    echo "-R,:example_25:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 26 from -R, help text
echo "🧪 Testing: ./gh-comment list --since "1 week ago" --type review --status open | review-blocker-analysis.sh"
./gh-comment list --since "1 week ago" --type review --status open | review-blocker-analysis.sh
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 26"
    echo "-R,:example_26:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 26"
    echo "-R,:example_26:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 27 from -R, help text
echo "🧪 Testing: ./gh-comment edit 2246362251 "Updated: This rate limiting logic handles concurrent requests properly""
./gh-comment edit 2246362251 "Updated: This rate limiting logic handles concurrent requests properly"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 27"
    echo "-R,:example_27:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 27"
    echo "-R,:example_27:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 28 from -R, help text
echo "🧪 Testing: ./gh-comment list $PR_NUM --author "bot*" --quiet | grep "ID:" | cut -d':' -f2 | xargs -I {} ./gh-comment resolve {}"
./gh-comment list $PR_NUM --author "bot*" --quiet | grep "ID:" | cut -d':' -f2 | xargs -I {} ./gh-comment resolve {}
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 28"
    echo "-R,:example_28:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 28"
    echo "-R,:example_28:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 29 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM performance.js 89:95 "Consider caching this expensive calculation""
./gh-comment add $PR_NUM performance.js 89:95 "Consider caching this expensive calculation"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 29"
    echo "-R,:example_29:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 29"
    echo "-R,:example_29:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 30 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/src/api.js 42 "[SUGGEST: const timeout = 5000;]""
./gh-comment add $PR_NUM src/src/api.js 42 "[SUGGEST: const timeout = 5000;]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 30"
    echo "-R,:example_30:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 30"
    echo "-R,:example_30:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 31 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/src/api.js 40 "[SUGGEST:+2: // Add error handling]""
./gh-comment add $PR_NUM src/src/api.js 40 "[SUGGEST:+2: // Add error handling]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 31"
    echo "-R,:example_31:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 31"
    echo "-R,:example_31:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

# Example 32 from -R, help text
echo "🧪 Testing: ./gh-comment add $PR_NUM src/src/api.js 45 "[SUGGEST:-1: import { logger } from './utils';]""
./gh-comment add $PR_NUM src/src/api.js 45 "[SUGGEST:-1: import { logger } from './utils';]"
if [[ $? -eq 0 ]]; then
    echo "✅ SUCCESS: -R, example 32"
    echo "-R,:example_32:SUCCESS" >> "$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: -R, example 32"
    echo "-R,:example_32:FAILED" >> "$TEST_RESULTS_FILE"
fi
echo ""

