#!/bin/bash
# Integration test environment variables - sourced by other scripts

export TEST_BRANCH="integration-test-20250805-072949"
export PR_NUM="14"
export PR_URL="https://github.com/silouanwright/gh-comment/pull/14"
export TIMESTAMP="20250805-072949"
export LOG_FILE="/Users/silouan.wright/repos/gh-comment/logs/integration/setup-20250805-072949.log"
export PROJECT_ROOT="/Users/silouan.wright/repos/gh-comment"

# Test files created
export TEST_FILES=(
    "src/api.js"
    "src/main.go"
    "tests/auth_test.js"
    "README.md"
)
# Test cases information
export TEST_CASES_DIR="/Users/silouan.wright/repos/gh-comment/logs/integration/test-cases-20250805-072949"
export MASTER_TEST_SCRIPT="/Users/silouan.wright/repos/gh-comment/logs/integration/test-cases-20250805-072949/run_all_examples.sh"
export TEST_RESULTS_FILE="/Users/silouan.wright/repos/gh-comment/logs/integration/test-results-20250805-072949.txt"
export TOTAL_EXAMPLES="113"
