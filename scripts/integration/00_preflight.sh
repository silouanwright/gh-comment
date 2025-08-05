#!/bin/bash
set -euo pipefail

# 00_preflight.sh - Pre-flight checks for integration testing
# Ensures environment is ready for integration tests

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "🔍 Running pre-flight checks..."

# Check if we're in the project root
if [[ ! -f "$PROJECT_ROOT/go.mod" ]] || [[ ! -f "$PROJECT_ROOT/main.go" ]]; then
    echo "❌ ERROR: Not in gh-comment project root"
    echo "   Expected go.mod and main.go files"
    exit 1
fi

# Check working tree is clean
if ! git diff-index --quiet HEAD --; then
    echo "❌ ERROR: Working tree has uncommitted changes"
    echo "   Commit or stash changes before running integration tests"
    echo ""
    echo "   Uncommitted files:"
    git status --porcelain | head -10
    if [[ $(git status --porcelain | wc -l) -gt 10 ]]; then
        echo "   ... and $(( $(git status --porcelain | wc -l) - 10 )) more files"
    fi
    exit 1
fi

# Check if on main branch (recommended)
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$CURRENT_BRANCH" != "main" ]]; then
    echo "⚠️  WARNING: Not on main branch (currently on: $CURRENT_BRANCH)"
    echo "   Integration tests typically run from main branch"
    echo "   Continue anyway? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "   Aborted by user"
        exit 1
    fi
fi

# Check required tools are available
REQUIRED_TOOLS=("go" "git" "gh")
for tool in "${REQUIRED_TOOLS[@]}"; do
    if ! command -v "$tool" &> /dev/null; then
        echo "❌ ERROR: Required tool '$tool' not found in PATH"
        exit 1
    fi
done

# Check GitHub CLI authentication
if ! gh auth status &> /dev/null; then
    echo "❌ ERROR: GitHub CLI not authenticated"
    echo "   Run: gh auth login"
    exit 1
fi

# Check if binary can be built
echo "🔨 Testing build..."
cd "$PROJECT_ROOT"
if ! go build -o /tmp/gh-comment-test .; then
    echo "❌ ERROR: Failed to build gh-comment binary"
    exit 1
fi
rm -f /tmp/gh-comment-test

# Check test coverage (optional warning)
if command -v bc &> /dev/null; then
    COVERAGE=$(go test ./cmd -cover 2>/dev/null | grep "coverage:" | grep -o '[0-9.]*%' | sed 's/%//' || echo "0")
    if (( $(echo "$COVERAGE < 80" | bc -l) )); then
        echo "⚠️  WARNING: Test coverage is ${COVERAGE}% (below 80%)"
        echo "   Consider running tests first: go test ./cmd"
    else
        echo "✅ Test coverage: ${COVERAGE}%"
    fi
fi

# Create scripts/integration directory if it doesn't exist
mkdir -p "$PROJECT_ROOT/scripts/integration"

# Create logs directory for integration test results
mkdir -p "$PROJECT_ROOT/logs/integration"

echo "✅ Pre-flight checks passed!"
echo ""
echo "Environment ready for integration testing:"
echo "  • Project root: $PROJECT_ROOT"
echo "  • Current branch: $CURRENT_BRANCH"
echo "  • GitHub CLI authenticated as: $(gh api user --jq .login 2>/dev/null || echo 'unknown')"
echo "  • Build test: ✅"
echo ""
