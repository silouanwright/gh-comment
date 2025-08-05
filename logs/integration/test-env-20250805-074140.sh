#!/bin/bash
# Integration test environment variables - sourced by other scripts

export TEST_BRANCH="integration-test-20250805-074140"
export PR_NUM="15"
export PR_URL="https://github.com/silouanwright/gh-comment/pull/15"
export TIMESTAMP="20250805-074140"
export LOG_FILE="/Users/silouan.wright/repos/gh-comment/logs/integration/setup-20250805-074140.log"
export PROJECT_ROOT="/Users/silouan.wright/repos/gh-comment"

# Test files created
export TEST_FILES=(
    "src/api.js"
    "src/main.go"
    "tests/auth_test.js"
    "review-config.yaml"
    "security-audit.yaml"
    "comprehensive-review.yaml"
    "README.md"
)

# Test cases information
export TEST_CASES_DIR="/Users/silouan.wright/repos/gh-comment/logs/integration/test-cases-20250805-074140"
export MASTER_TEST_SCRIPT="/Users/silouan.wright/repos/gh-comment/logs/integration/test-cases-20250805-074140/run_all_examples.sh"
export TEST_RESULTS_FILE="/Users/silouan.wright/repos/gh-comment/logs/integration/test-results-20250805-074140.txt"
export TOTAL_EXAMPLES="113"
