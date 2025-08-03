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

## Reply Command Testing

**CRITICAL**: The reply command recently had its GitHub API bug fixed. These tests verify the fix works correctly.

### Background
- GitHub API only supports message replies for issue comments, not review comments
- Review comments only support reactions for replies, messages fail with HTTP 422
- Fixed: Changed API endpoint from incorrect in_reply_to_id approach to proper /repos/{owner}/{repo}/pulls/comments/{comment_id}/replies

### Test Sequence

```bash
# Step 1: Create a general PR comment first
./gh-comment add $PR_NUM "This is a general discussion comment for reply testing"

# Step 2: Get the comment ID
COMMENT_ID=$(./gh-comment list $PR_NUM --quiet | grep "general discussion" | grep -o 'ID:[0-9]*' | cut -d: -f2)
echo "Testing with comment ID: $COMMENT_ID"

# Step 3: Test all reply help examples (they should work now)
./gh-comment reply --help

# Test message reply to issue comment (should work - default type is 'issue')
./gh-comment reply $COMMENT_ID "Good point, I'll fix that"

# Test reactions (should work for both comment types)
./gh-comment reply $COMMENT_ID --reaction +1
./gh-comment reply $COMMENT_ID --reaction heart

# Test removing reactions 
./gh-comment reply $COMMENT_ID --remove-reaction +1

# Test invalid type validation
! ./gh-comment reply $COMMENT_ID "message" --type invalid

# Step 4: Create line-specific review comment for review testing
./gh-comment add $PR_NUM src/api.js 30 "This is a review comment for testing"

# Get review comment ID
REVIEW_COMMENT_ID=$(./gh-comment list $PR_NUM --type review --quiet | grep "review comment" | grep -o 'ID:[0-9]*' | cut -d: -f2)

# Test reactions on review comments (should work)
./gh-comment reply $REVIEW_COMMENT_ID --reaction +1 --type review
./gh-comment reply $REVIEW_COMMENT_ID --reaction rocket --type review

# Test message reply to review comment (should work with fixed API endpoint)
./gh-comment reply $REVIEW_COMMENT_ID "Fixed in commit abc123" --type review

# Test resolve functionality (review comments only)
./gh-comment reply $REVIEW_COMMENT_ID --reaction thumbs_up --resolve --type review
```

### Expected Results
- ✅ All help text examples should work without HTTP 422 errors
- ✅ Message replies to issue comments work (create new issue comments)  
- ✅ Message replies to review comments work (create threaded replies)
- ✅ Reactions work for both comment types
- ✅ Resolve works only for review comments
- ✅ Error validation works correctly

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