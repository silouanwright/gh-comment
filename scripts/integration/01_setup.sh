#!/bin/bash
set -euo pipefail

# 01_setup.sh - Setup integration test environment
# Creates test branch, test files, and GitHub PR

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
LOG_DIR="$PROJECT_ROOT/logs/integration"

# Generate unique branch name
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
TEST_BRANCH="integration-test-$TIMESTAMP"
LOG_FILE="$LOG_DIR/setup-$TIMESTAMP.log"

echo "🚀 Setting up integration test environment..."
echo "📝 Logging to: $LOG_FILE"

# Create log file
exec 1> >(tee -a "$LOG_FILE")
exec 2> >(tee -a "$LOG_FILE" >&2)

cd "$PROJECT_ROOT"

echo "🔧 Building local binary..."
if ! go build -o gh-comment .; then
    echo "❌ ERROR: Failed to build binary"
    exit 1
fi

echo "✅ Binary built successfully"
echo ""

echo "🌿 Creating test branch: $TEST_BRANCH"
git checkout -b "$TEST_BRANCH"

echo "📁 Creating test files for help text examples..."

# Create directory structure
mkdir -p src tests

# Create src/api.js (commonly referenced in help text)
cat > src/api.js << 'EOF'
const express = require('express');
const rateLimit = require('express-rate-limit');

// TODO: Add comprehensive rate limiting
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100, // limit each IP to 100 requests per windowMs
    message: 'Too many requests from this IP'
});

function authenticate(req, res, next) {
    const token = req.headers.authorization;
    // TODO: Implement proper JWT validation
    if (!token) {
        return res.status(401).json({ error: 'No token provided' });
    }
    next();
}

// Middleware with potential security issue for testing
function processUserInput(input) {
    // SECURITY: This should be sanitized
    return input.toLowerCase();
}

module.exports = { limiter, authenticate, processUserInput };
EOF

# Create src/main.go (commonly referenced in help text)
cat > src/main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
)

func main() {
    // TODO: Add proper error handling
    if len(os.Args) < 2 {
        log.Fatal("Usage: main <command>")
    }

    command := os.Args[1]
    fmt.Printf("Executing command: %s\n", command)

    // TODO: Implement actual command processing
    switch command {
    case "start":
        startServer()
    case "stop":
        stopServer()
    default:
        fmt.Printf("Unknown command: %s\n", command)
    }
}

// Function that could use better error handling
func startServer() {
    fmt.Println("Starting server...")
    // TODO: Add proper server initialization
}

func stopServer() {
    fmt.Println("Stopping server...")
}
EOF

# Create tests/auth_test.js (commonly referenced in help text)
cat > tests/auth_test.js << 'EOF'
const { authenticate } = require('../src/api');

describe('Authentication Tests', () => {
    it('should reject requests without token', async () => {
        // TODO: Add proper test implementation
        const mockReq = { headers: {} };
        const mockRes = {
            status: jest.fn().mockReturnThis(),
            json: jest.fn()
        };

        authenticate(mockReq, mockRes, () => {});
        expect(mockRes.status).toHaveBeenCalledWith(401);
    });

    it('should validate JWT tokens', () => {
        // TODO: Implement JWT validation test
        expect(true).toBe(true); // Placeholder
    });

    // Test case that needs review
    it('should handle edge cases', () => {
        // This test needs more specific assertions
        expect(1 + 1).toBe(2);
    });
});
EOF

# Create README.md for additional testing
cat > README.md << 'EOF'
# Integration Test Repository

This is a temporary repository created for integration testing `gh-comment`.

## Test Files

- `src/api.js` - Express.js API with authentication middleware
- `src/main.go` - Go command-line application
- `tests/auth_test.js` - Jest test suite for authentication

## Purpose

These files contain realistic code examples that match the help text examples in `gh-comment`. They include:
- Security issues for testing security-focused prompts
- TODO comments for improvement suggestions
- Various file types (.js, .go, .js) for comprehensive testing

## Usage

This repository is automatically created and cleaned up by integration test scripts.
EOF

echo "📄 Test files created:"
echo "  • src/api.js (Express.js middleware)"
echo "  • src/main.go (Go CLI application)"
echo "  • tests/auth_test.js (Jest test suite)"
echo "  • README.md (Documentation)"
echo ""

echo "💾 Committing test files..."
git add .
git commit --no-verify -m "feat: add test files for gh-comment integration testing

- src/api.js: Express middleware with auth and rate limiting
- src/main.go: Go CLI application with command processing
- tests/auth_test.js: Jest test suite for authentication
- README.md: Documentation for test repository

These files provide realistic examples that match gh-comment help text,
including security issues and improvement opportunities for testing
various comment and review scenarios."

echo "🚀 Pushing test branch..."
git push -u origin "$TEST_BRANCH"

echo "📝 Creating GitHub PR..."
PR_TITLE="Integration Test: Help Text Validation ($TIMESTAMP)"
PR_BODY="## Integration Test PR

**Purpose:** Validate all help text examples work correctly with real GitHub API

**Test Branch:** \`$TEST_BRANCH\`
**Timestamp:** $TIMESTAMP

### Test Files
- \`src/api.js\` - Express.js middleware (authentication, rate limiting)
- \`src/main.go\` - Go CLI application
- \`tests/auth_test.js\` - Jest authentication tests
- \`README.md\` - Test documentation

### Testing Process
1. Build local binary: \`go build\`
2. Extract help text examples from all commands
3. Execute every example with this PR number
4. Validate all commands work as documented

### Expected Outcome
- ✅ All help text examples execute successfully
- ✅ No incorrect or non-working examples in documentation
- ✅ Real GitHub API integration confirmed working

**Note:** This PR will be automatically closed after testing."

# Create PR and capture number
PR_OUTPUT=$(gh pr create --title "$PR_TITLE" --body "$PR_BODY" 2>&1)
PR_URL=$(echo "$PR_OUTPUT" | grep -o 'https://github.com/[^/]*/[^/]*/pull/[0-9]*' || echo "")

if [[ -z "$PR_URL" ]]; then
    echo "❌ ERROR: Failed to create PR"
    echo "Output: $PR_OUTPUT"
    exit 1
fi

# Extract PR number from URL
PR_NUM=$(echo "$PR_URL" | grep -o '[0-9]*$')

echo "✅ PR created successfully!"
echo "  • PR #$PR_NUM: $PR_URL"
echo ""

# Save environment info for other scripts
cat > "$LOG_DIR/test-env-$TIMESTAMP.sh" << EOF
#!/bin/bash
# Integration test environment variables - sourced by other scripts

export TEST_BRANCH="$TEST_BRANCH"
export PR_NUM="$PR_NUM"
export PR_URL="$PR_URL"
export TIMESTAMP="$TIMESTAMP"
export LOG_FILE="$LOG_FILE"
export PROJECT_ROOT="$PROJECT_ROOT"

# Test files created
export TEST_FILES=(
    "src/api.js"
    "src/main.go"
    "tests/auth_test.js"
    "README.md"
)
EOF

echo "📋 Environment saved to: $LOG_DIR/test-env-$TIMESTAMP.sh"
echo ""
echo "🎯 Integration test setup complete!"
echo ""
echo "Next steps:"
echo "  1. Run: source $LOG_DIR/test-env-$TIMESTAMP.sh"
echo "  2. Execute: scripts/integration/02_parse_help.sh"
echo "  3. Or run full test suite: scripts/integration/run_integration.sh"
echo ""
echo "PR Details:"
echo "  • Number: #$PR_NUM"
echo "  • URL: $PR_URL"
echo "  • Branch: $TEST_BRANCH"
echo ""
