# Integration Testing Guide for gh-comment

## Overview

This guide follows a help-text-driven approach to ensure we dogfood our own documentation. Your task is simple: run every example from every command's help text.

## Setup

### Step 1: Build Local Binary

```bash
# Build the local development version
go build

# Verify binary was created
ls -la gh-comment
```

### Step 2: Create Test PR

```bash
# Create test branch
git checkout -b integration-test-$(date +%Y%m%d-%H%M%S)

# Add test files that match help text examples
mkdir -p src tests
cat > src/api.js << 'EOF'
function rateLimit(req, res, next) {
    // TODO: Add rate limiting logic
    next();
}

function authenticate(token) {
    const apiKey = "sk-test-key"; // TODO: Move to env vars
    return token === apiKey;
}

module.exports = { rateLimit, authenticate };
EOF

cat > src/main.go << 'EOF'
package main

func main() {
    println("Hello from gh-comment")
}
EOF

cat > tests/auth_test.js << 'EOF'
describe('auth', () => {
    it('should validate tokens', () => {
        // Test implementation
    });
});
EOF

# Commit and push
git add .
git commit -m "Add test files for gh-comment help text examples"
git push -u origin $(git branch --show-current)

# Create PR
gh pr create --title "Integration Test: Dogfooding Help Text" --body "Testing all help text examples"

# Get PR number
PR_NUM=$(gh pr view --json number -q .number)
echo "Testing with PR #$PR_NUM"
```

## Testing Process

Your task:

1. Run `./gh-comment --help` to see all available commands
2. For each command listed, run `./gh-comment [command] --help`
3. Execute every example shown in the help text
   - Replace `gh comment` with `./gh-comment`
   - Replace PR numbers (usually `123`) with `$PR_NUM`
   - Use actual comment IDs from `./gh-comment list` output where needed

## Success Criteria

✅ Every example from every command's help text has been executed
✅ All commands work as documented
✅ No help text contains incorrect or non-working examples

## Cleanup

```bash
# After testing
gh pr close $PR_NUM
git checkout main
git branch -D $(git branch --show-current)
```

## Important Notes

- Always use `./gh-comment` (local binary), not `gh comment` (global extension)
- The help text is the single source of truth - if it doesn't work, the help text needs fixing
- Don't skip any examples - the goal is 100% help text coverage