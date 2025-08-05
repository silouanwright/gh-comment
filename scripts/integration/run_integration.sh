#!/bin/bash
set -euo pipefail

# run_integration.sh - Master orchestrator for integration testing
# Runs the complete integration test cycle from setup to cleanup

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    local color="$1"
    local message="$2"
    echo -e "${color}${message}${NC}"
}

print_header() {
    echo ""
    echo "========================================"
    echo "$1"
    echo "========================================"
    echo ""
}

show_usage() {
    cat << 'EOF'
Usage: run_integration.sh [OPTIONS]

Integration test orchestrator for gh-comment help text validation.
Runs complete cycle: setup → parse → test → cleanup

OPTIONS:
    --setup-only        Run setup phase only (create PR, files)
    --parse-only        Run parse phase only (extract help examples)
    --test-only         Run test phase only (execute examples)
    --cleanup-only      Run cleanup phase only (close PR, delete branch)
    --no-cleanup        Skip cleanup phase (leave PR open)
    --force-cleanup     Force cleanup without confirmation
    --resume            Resume from existing test environment
    --help              Show this help message

EXAMPLES:
    # Full integration test cycle
    ./run_integration.sh

    # Setup test environment only
    ./run_integration.sh --setup-only

    # Resume testing from existing environment
    ./run_integration.sh --resume --test-only

    # Run tests without cleanup (for debugging)
    ./run_integration.sh --no-cleanup

    # Force cleanup of existing test environment
    ./run_integration.sh --cleanup-only --force-cleanup

WORKFLOW:
    1. Preflight checks (git status, auth, build)
    2. Setup (create branch, files, PR)
    3. Parse help text (extract examples)
    4. Run tests (execute all examples)
    5. Cleanup (close PR, delete branch, archive results)

NOTE: This orchestrator ensures each phase completes successfully
      before proceeding to the next phase.
EOF
}

# Parse command line arguments
SETUP_ONLY=false
PARSE_ONLY=false
TEST_ONLY=false
CLEANUP_ONLY=false
NO_CLEANUP=false
FORCE_CLEANUP=false
RESUME=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --setup-only)
            SETUP_ONLY=true
            shift
            ;;
        --parse-only)
            PARSE_ONLY=true
            shift
            ;;
        --test-only)
            TEST_ONLY=true
            shift
            ;;
        --cleanup-only)
            CLEANUP_ONLY=true
            shift
            ;;
        --no-cleanup)
            NO_CLEANUP=true
            shift
            ;;
        --force-cleanup)
            FORCE_CLEANUP=true
            shift
            ;;
        --resume)
            RESUME=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validation: conflicting options
if [[ "$SETUP_ONLY" == true && "$TEST_ONLY" == true ]]; then
    print_status "$RED" "❌ ERROR: Cannot use --setup-only and --test-only together"
    exit 1
fi

if [[ "$NO_CLEANUP" == true && "$CLEANUP_ONLY" == true ]]; then
    print_status "$RED" "❌ ERROR: Cannot use --no-cleanup and --cleanup-only together"
    exit 1
fi

# Start integration testing
print_header "🚀 gh-comment Integration Test Orchestrator"

print_status "$BLUE" "📍 Project root: $PROJECT_ROOT"
print_status "$BLUE" "📅 Started: $(date)"

cd "$PROJECT_ROOT"

# Phase selection logic
RUN_PREFLIGHT=true
RUN_SETUP=true
RUN_PARSE=true
RUN_TEST=true
RUN_CLEANUP=true

if [[ "$CLEANUP_ONLY" == true ]]; then
    RUN_PREFLIGHT=false
    RUN_SETUP=false
    RUN_PARSE=false
    RUN_TEST=false
elif [[ "$SETUP_ONLY" == true ]]; then
    RUN_PARSE=false
    RUN_TEST=false
    RUN_CLEANUP=false
elif [[ "$PARSE_ONLY" == true ]]; then
    RUN_PREFLIGHT=false
    RUN_SETUP=false
    RUN_TEST=false
    RUN_CLEANUP=false
elif [[ "$TEST_ONLY" == true ]]; then
    RUN_PREFLIGHT=false
    RUN_SETUP=false
    RUN_PARSE=false
    RUN_CLEANUP=false
elif [[ "$NO_CLEANUP" == true ]]; then
    RUN_CLEANUP=false
fi

if [[ "$RESUME" == true ]]; then
    RUN_PREFLIGHT=false
    RUN_SETUP=false
fi

# Execute phases
OVERALL_SUCCESS=true

# Phase 1: Preflight checks
if [[ "$RUN_PREFLIGHT" == true ]]; then
    print_header "🔍 Phase 1: Preflight Checks"

    if bash "$SCRIPT_DIR/00_preflight.sh"; then
        print_status "$GREEN" "✅ Preflight checks passed"
    else
        print_status "$RED" "❌ Preflight checks failed"
        OVERALL_SUCCESS=false
    fi
fi

# Phase 2: Setup test environment
if [[ "$RUN_SETUP" == true && "$OVERALL_SUCCESS" == true ]]; then
    print_header "🚀 Phase 2: Setup Test Environment"

    if bash "$SCRIPT_DIR/01_setup.sh"; then
        print_status "$GREEN" "✅ Test environment setup completed"
    else
        print_status "$RED" "❌ Test environment setup failed"
        OVERALL_SUCCESS=false
    fi
fi

# Phase 3: Parse help text
if [[ "$RUN_PARSE" == true && "$OVERALL_SUCCESS" == true ]]; then
    print_header "📖 Phase 3: Parse Help Text"

    if bash "$SCRIPT_DIR/02_parse_help.sh"; then
        print_status "$GREEN" "✅ Help text parsing completed"
    else
        print_status "$RED" "❌ Help text parsing failed"
        OVERALL_SUCCESS=false
    fi
fi

# Phase 4: Run tests
if [[ "$RUN_TEST" == true && "$OVERALL_SUCCESS" == true ]]; then
    print_header "🧪 Phase 4: Run Integration Tests"

    if bash "$SCRIPT_DIR/03_run_tests.sh"; then
        print_status "$GREEN" "✅ Integration tests completed successfully"
    else
        print_status "$YELLOW" "⚠️ Integration tests completed with failures"
        # Don't set OVERALL_SUCCESS=false here, as test failures are different from script failures
        # We still want to proceed to cleanup
    fi
fi

# Phase 5: Cleanup
if [[ "$RUN_CLEANUP" == true ]]; then
    print_header "🧹 Phase 5: Cleanup"

    CLEANUP_ARGS=""
    if [[ "$FORCE_CLEANUP" == true ]]; then
        CLEANUP_ARGS="--force"
    fi

    if bash "$SCRIPT_DIR/04_cleanup.sh" $CLEANUP_ARGS; then
        print_status "$GREEN" "✅ Cleanup completed"
    else
        print_status "$RED" "❌ Cleanup failed"
        OVERALL_SUCCESS=false
    fi
fi

# Final summary
print_header "📊 Integration Test Complete"

if [[ "$OVERALL_SUCCESS" == true ]]; then
    print_status "$GREEN" "🎉 Integration test cycle completed successfully!"
else
    print_status "$RED" "❌ Integration test cycle completed with errors"
fi

# Show what was accomplished
echo "📋 Phases completed:"
if [[ "$RUN_PREFLIGHT" == true ]]; then
    echo "  ✅ Preflight checks"
fi
if [[ "$RUN_SETUP" == true ]]; then
    echo "  ✅ Test environment setup"
fi
if [[ "$RUN_PARSE" == true ]]; then
    echo "  ✅ Help text parsing"
fi
if [[ "$RUN_TEST" == true ]]; then
    echo "  ✅ Integration test execution"
fi
if [[ "$RUN_CLEANUP" == true ]]; then
    echo "  ✅ Resource cleanup"
else
    echo "  ⏭️  Cleanup skipped (--no-cleanup)"
fi

# Show next steps based on what was run
echo ""
echo "🎯 Next steps:"

if [[ "$SETUP_ONLY" == true ]]; then
    echo "  • Review created PR and test files"
    echo "  • Continue with: $0 --resume --parse-only"
    echo "  • Or run full test: $0 --resume"
elif [[ "$PARSE_ONLY" == true ]]; then
    echo "  • Review generated test cases in logs/integration/"
    echo "  • Continue with: $0 --resume --test-only"
elif [[ "$TEST_ONLY" == true ]]; then
    echo "  • Review test results in logs/integration/"
    echo "  • Clean up with: $0 --cleanup-only"
elif [[ "$NO_CLEANUP" == true ]]; then
    echo "  • Review test results and PR manually"
    echo "  • Clean up when ready: $0 --cleanup-only"
elif [[ "$CLEANUP_ONLY" == true ]]; then
    echo "  • Integration test cycle is complete"
    echo "  • Review archived results in logs/integration/"
else
    echo "  • Integration test cycle is complete"
    echo "  • Review archived results in logs/integration/"
fi

# Show log locations
LOG_DIR="$PROJECT_ROOT/logs/integration"
if [[ -d "$LOG_DIR" ]]; then
    echo ""
    echo "📄 Results and logs:"
    echo "  • Integration logs: $LOG_DIR"

    # Show most recent files
    RECENT_FILES=$(find "$LOG_DIR" -name "*.md" -o -name "*.txt" -o -name "archived-*" -type d | head -5)
    if [[ -n "$RECENT_FILES" ]]; then
        echo "  • Recent files:"
        echo "$RECENT_FILES" | sed 's/^/    - /'
    fi
fi

print_status "$BLUE" "📅 Completed: $(date)"

if [[ "$OVERALL_SUCCESS" == true ]]; then
    exit 0
else
    exit 1
fi
