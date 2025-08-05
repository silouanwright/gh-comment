#!/bin/bash
set -euo pipefail

# 02_parse_help.sh - Parse help text and extract examples for testing
# Generates executable test cases from command help text

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs/integration"

# Find the latest test environment
LATEST_ENV=$(ls -t "$LOG_DIR"/test-env-*.sh 2>/dev/null | head -1 || echo "")

if [[ -z "$LATEST_ENV" ]]; then
    echo "❌ ERROR: No test environment found"
    echo "   Run scripts/integration/01_setup.sh first"
    exit 1
fi

echo "📖 Parsing help text and extracting examples..."
echo "🔗 Using environment: $(basename "$LATEST_ENV")"

# Source test environment
source "$LATEST_ENV"

cd "$PROJECT_ROOT"

# Ensure binary exists
if [[ ! -f "./gh-comment" ]]; then
    echo "🔨 Building binary..."
    go build -o gh-comment .
fi

# Create test cases directory
TEST_CASES_DIR="$LOG_DIR/test-cases-$TIMESTAMP"
mkdir -p "$TEST_CASES_DIR"

echo "📝 Extracting examples from help text..."

# Get list of all commands
echo "🔍 Discovering commands..."
COMMANDS=$(./gh-comment --help 2>/dev/null | grep -E "^  [a-z-]+" | awk '{print $1}' | grep -v "^help$" || echo "")

if [[ -z "$COMMANDS" ]]; then
    echo "❌ ERROR: Could not discover commands from help text"
    exit 1
fi

echo "📋 Found commands: $(echo $COMMANDS | tr '\n' ' ')"
echo ""

# Function to extract examples from help text
extract_examples() {
    local cmd="$1"
    local help_output="$2"
    local examples_file="$3"

    echo "# Examples for command: $cmd" > "$examples_file"
    echo "# Generated from help text on $(date)" >> "$examples_file"
    echo "" >> "$examples_file"

    # Look for Examples section in help text
    local in_examples=false
    local example_count=0

    while IFS= read -r line; do
        # Start of examples section
        if [[ "$line" =~ ^[[:space:]]*Examples?:[[:space:]]*$ ]]; then
            in_examples=true
            continue
        fi

        # End of examples section (next section or end of help)
        if [[ "$in_examples" == true ]] && [[ "$line" =~ ^[[:space:]]*[A-Z][a-z]+:[[:space:]]*$ ]]; then
            in_examples=false
            continue
        fi

        # If we're in examples section, look for command examples
        if [[ "$in_examples" == true ]]; then
            # Look for lines that start with gh comment or ./gh-comment
            if [[ "$line" =~ ^[[:space:]]*(\$[[:space:]]+)?(gh[[:space:]]+comment|\.\/gh-comment) ]]; then
                ((example_count++))

                # Clean up the example
                local example=$(echo "$line" | sed -E 's/^[[:space:]]*\$?[[:space:]]*//' | sed 's/gh comment/\.\/gh-comment/g')

                # Replace common placeholders
                example=$(echo "$example" | sed "s/123/\$PR_NUM/g")
                example=$(echo "$example" | sed "s/<pr-number>/\$PR_NUM/g")
                example=$(echo "$example" | sed "s/<comment-id>/\$COMMENT_ID/g")

                # Replace file placeholders with our test files (order matters!)
                # First replace specific patterns, then generic ones
                example=$(echo "$example" | sed "s/auth\.go/src\/main.go/g")
                example=$(echo "$example" | sed "s/validation\.js/tests\/auth_test.js/g")
                example=$(echo "$example" | sed "s/database\.py/src\/main.go/g")
                # Only replace bare main.go (not already prefixed with src/)
                example=$(echo "$example" | sed "s/\([^/]\)main\.go/\1src\/main.go/g")
                # Handle main.go at the beginning of a path
                example=$(echo "$example" | sed "s/^main\.go/src\/main.go/g")
                # Handle api.js in any location (already prefixed paths won't match)
                example=$(echo "$example" | sed "s/\([^/]\)api\.js/\1src\/api.js/g")
                example=$(echo "$example" | sed "s/^api\.js/src\/api.js/g")

                # Write test case
                cat >> "$examples_file" << EOF
# Example $example_count from $cmd help text
echo "🧪 Testing: $example"
$example
if [[ \$? -eq 0 ]]; then
    echo "✅ SUCCESS: $cmd example $example_count"
    echo "$cmd:example_$example_count:SUCCESS" >> "\$TEST_RESULTS_FILE"
else
    echo "❌ FAILED: $cmd example $example_count"
    echo "$cmd:example_$example_count:FAILED" >> "\$TEST_RESULTS_FILE"
fi
echo ""

EOF
            fi
        fi
    done <<< "$help_output"

    if [[ "$example_count" -eq 0 ]]; then
        echo "# No examples found in help text for $cmd" >> "$examples_file"
        echo "echo \"⚠️  WARNING: No examples found for command '$cmd'\"" >> "$examples_file"
        echo "echo \"$cmd:no_examples:WARNING\" >> \"\$TEST_RESULTS_FILE\"" >> "$examples_file"
    else
        echo "  📝 Found $example_count examples for $cmd"
    fi
}

# Parse help text for each command
TOTAL_EXAMPLES=0
for cmd in $COMMANDS; do
    echo "📖 Parsing help for: $cmd"

    # Get help text
    HELP_OUTPUT=$(./gh-comment "$cmd" --help 2>/dev/null || echo "ERROR: Could not get help for $cmd")

    if [[ "$HELP_OUTPUT" == "ERROR:"* ]]; then
        echo "  ⚠️  WARNING: Could not get help for $cmd"
        continue
    fi

    # Extract examples
    EXAMPLES_FILE="$TEST_CASES_DIR/${cmd}_examples.sh"
    extract_examples "$cmd" "$HELP_OUTPUT" "$EXAMPLES_FILE"

    # Count examples
    EXAMPLE_COUNT=$(grep -c "# Example [0-9]" "$EXAMPLES_FILE" 2>/dev/null || echo "0")
    TOTAL_EXAMPLES=$((TOTAL_EXAMPLES + EXAMPLE_COUNT))
done

echo ""
echo "📊 Help text parsing complete!"
echo "  • Commands processed: $(echo $COMMANDS | wc -w)"
echo "  • Total examples found: $TOTAL_EXAMPLES"
echo "  • Test cases saved to: $TEST_CASES_DIR"
echo ""

# Create master test runner script
MASTER_TEST_SCRIPT="$TEST_CASES_DIR/run_all_examples.sh"
cat > "$MASTER_TEST_SCRIPT" << 'EOF'
#!/bin/bash
set -euo pipefail

# Master test runner for all help text examples
# Generated by 02_parse_help.sh

echo "🚀 Running all help text examples..."
echo "📊 Test results will be saved to: $TEST_RESULTS_FILE"
echo ""

# Initialize results file
echo "# Integration test results - $(date)" > "$TEST_RESULTS_FILE"
echo "# Format: command:example:status" >> "$TEST_RESULTS_FILE"
echo "" >> "$TEST_RESULTS_FILE"

# Source environment
if [[ -n "${1:-}" ]] && [[ -f "$1" ]]; then
    source "$1"
    echo "✅ Loaded environment: $1"
else
    echo "❌ ERROR: Please provide test environment file as first argument"
    echo "   Usage: $0 path/to/test-env-TIMESTAMP.sh"
    exit 1
fi

# Run all example scripts
TEST_SCRIPTS=($(find "$(dirname "$0")" -name "*_examples.sh" | sort))

if [[ ${#TEST_SCRIPTS[@]} -eq 0 ]]; then
    echo "❌ ERROR: No test scripts found"
    exit 1
fi

echo "📋 Found ${#TEST_SCRIPTS[@]} test script(s)"
echo ""

for script in "${TEST_SCRIPTS[@]}"; do
    echo "▶️  Running: $(basename "$script")"
    bash "$script"
done

echo ""
echo "📊 Integration test complete!"
echo "📋 Results summary:"

# Count results
TOTAL_TESTS=$(grep -c ":" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
SUCCESS_COUNT=$(grep -c ":SUCCESS" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
FAILED_COUNT=$(grep -c ":FAILED" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")
WARNING_COUNT=$(grep -c ":WARNING" "$TEST_RESULTS_FILE" 2>/dev/null || echo "0")

echo "  • Total tests: $TOTAL_TESTS"
echo "  • Successful: $SUCCESS_COUNT"
echo "  • Failed: $FAILED_COUNT"
echo "  • Warnings: $WARNING_COUNT"

if [[ "$FAILED_COUNT" -gt 0 ]]; then
    echo ""
    echo "❌ Failed tests:"
    grep ":FAILED" "$TEST_RESULTS_FILE" | sed 's/:FAILED//' | sed 's/^/  • /'
fi

if [[ "$WARNING_COUNT" -gt 0 ]]; then
    echo ""
    echo "⚠️  Warnings:"
    grep ":WARNING" "$TEST_RESULTS_FILE" | sed 's/:WARNING//' | sed 's/^/  • /'
fi

echo ""
echo "📄 Detailed results: $TEST_RESULTS_FILE"
EOF

chmod +x "$MASTER_TEST_SCRIPT"

# Update environment file with test cases info
cat >> "$LOG_DIR/test-env-$TIMESTAMP.sh" << EOF

# Test cases information
export TEST_CASES_DIR="$TEST_CASES_DIR"
export MASTER_TEST_SCRIPT="$MASTER_TEST_SCRIPT"
export TEST_RESULTS_FILE="$LOG_DIR/test-results-$TIMESTAMP.txt"
export TOTAL_EXAMPLES="$TOTAL_EXAMPLES"
EOF

echo "✅ Help text parsing complete!"
echo ""
echo "📋 Generated test files:"
echo "  • Test cases directory: $TEST_CASES_DIR"
echo "  • Master test script: $MASTER_TEST_SCRIPT"
echo "  • Total examples to test: $TOTAL_EXAMPLES"
echo ""
echo "🎯 Next steps:"
echo "  1. Review test cases: ls $TEST_CASES_DIR"
echo "  2. Run tests: $MASTER_TEST_SCRIPT $LATEST_ENV"
echo "  3. Or use: scripts/integration/03_run_tests.sh"
echo ""
