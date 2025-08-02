# Integration Testing Guide for gh-comment

## Overview

This guide provides comprehensive instructions for running integration tests that demonstrate ALL gh-comment functionality against real GitHub APIs using the **local development version**. Integration tests create actual PRs, add various types of comments, and showcase the complete feature set.

## Prerequisites

**IMPORTANT**: Always test with the local development version, not the global extension.

### Step 1: Build Local Binary

```bash
# Build the local development version
go build

# Verify binary was created
ls -la gh-comment
```

### Step 2: Create Test PR

Create a branch and PR for testing:

```bash
# Create test branch
git checkout -b integration-test-$(date +%Y%m%d-%H%M%S)

# Add test file with intentional issues
cat > integration-test-example.js << 'EOF'
// Integration test file with intentional issues for commenting
function processUserData(users) {
    // Missing input validation
    let results = [];
    
    // SQL injection vulnerability
    const query = "SELECT * FROM users WHERE id = " + users[0].id;
    
    // Hardcoded secrets
    const apiKey = "sk-1234567890abcdef";
    
    // Performance issue - nested loops
    for (let i = 0; i < users.length; i++) {
        for (let j = 0; j < users.length; j++) {
            if (users[i].status === users[j].status) {
                results.push(users[i]);
            }
        }
    }
    
    return results;
}

module.exports = { processUserData };
EOF

# Commit and push
git add integration-test-example.js
git commit -m "Add integration test file for gh-comment demonstration"
git push -u origin $(git branch --show-current)

# Create PR
gh pr create --title "Integration Test: gh-comment Feature Demonstration" --body "This PR demonstrates all gh-comment functionality using the local development version."
```

## Quick Start

Integration tests are **MANUAL-ONLY** by default to prevent accidental API usage:

```bash
# Run all integration tests (FORCE FLAG REQUIRED)
go run -tags integration . test-integration --force

# Run specific scenario only  
go run -tags integration . test-integration --force --scenario=comments

# Leave PR open for inspection (useful for demos)
go run -tags integration . test-integration --force --inspect
```

## Environment Controls

```bash
# For CI/automation (always run)
export GH_COMMENT_INTEGRATION_TESTS=always
go run -tags integration . test-integration

# Explicitly disable
export GH_COMMENT_INTEGRATION_TESTS=never
go run -tags integration . test-integration  # Will skip
```

## Complete Feature Demonstration Checklist

The integration tests should demonstrate **ALL** these features following the code review best practices:

### ✅ 1. General PR Comments (Issue Comments)
- [ ] Add opening comment explaining the demonstration
- [ ] Add follow-up discussion comments
- [ ] Use professional, collaborative language

### ✅ 2. Line-Specific Review Comments
- [ ] Add security-focused comments with 🔧 emoji (must fix)
- [ ] Add performance suggestions with 🤔 emoji (questions)  
- [ ] Add architecture feedback with ♻️ emoji (refactor)
- [ ] Add educational notes with 📝 emoji (no action needed)
- [ ] Add praise comments with 😃 emoji (highlight good work)

### ✅ 3. Range Comments (Multi-line)
- [ ] Comment on code blocks (e.g., entire functions)
- [ ] Comment on documentation sections
- [ ] Use range syntax: `file.js:15:25`

### ✅ 4. Suggestion Syntax (Both Types)
- [ ] Simple inline suggestions: `[SUGGEST: const result = input?.value || 'default']`
- [ ] Multi-line suggestions using `<<<SUGGEST>>>` blocks
- [ ] Show automatic GitHub markdown formatting
- [ ] Include explanatory text with suggestions

### ✅ 5. Review Creation and Submission
- [ ] Create draft reviews with `add-review` command
- [ ] Add multiple review comments in one operation
- [ ] Submit reviews with different events (APPROVE, REQUEST_CHANGES, COMMENT)
- [ ] Show professional review summaries

### ✅ 6. Reactions and Replies
- [ ] Add emoji reactions to comments (+1, heart, hooray, etc.)
- [ ] Reply to review comments with threaded responses
- [ ] Reply to issue comments with follow-ups
- [ ] Remove reactions to show full functionality

### ✅ 7. Batch Operations
- [ ] Create YAML configuration files
- [ ] Process multiple comments at once
- [ ] Mix different comment types in batch
- [ ] Show review creation via batch

### ✅ 8. Comment Management
- [ ] List all comments with full context
- [ ] Edit existing comments
- [ ] Resolve conversation threads
- [ ] Filter comments by author, date, type

## Code Review Best Practices Integration

**🚨 IMPORTANT: READ THE STYLE GUIDE FIRST! 🚨**

**BEFORE starting integration tests, you MUST read and follow the patterns in:**
`research/code-review-best-practices.md`

This file contains essential communication patterns including:
- CREG emoji system (🔧🤔♻️📝😃📌)
- Question-based communication style
- Psychological safety principles
- Specific phrasing examples

**All integration test comments MUST follow these research-backed best practices.**

Based on `research/code-review-best-practices.md`, all comments should follow these patterns:

## Manual Integration Testing Process

**ALWAYS use the local development version**: `./gh-comment`

### Step 3: Get PR Number

```bash
# Get the PR number for testing
PR_NUM=$(gh pr view --json number -q .number)
echo "Testing with PR #$PR_NUM"
```

### Communication Style
```bash
# ❌ Bad: Command style
./gh-comment add $PR_NUM file.js 42 "Use a Map here"

# ✅ Good: Question style  
./gh-comment add $PR_NUM file.js 42 "🤔 What do you think about using a Map here? Better lookup performance."
```

### CREG Emoji System
- 🔧 **Must fix** - Required changes (security, bugs)
- ⛏️ **Nitpick** - Minor style issues, not blocking  
- 😃 **Praise** - Highlight good work
- 🤔 **Question** - Need clarification or thinking out loud
- 📝 **Note** - Educational info, no action needed
- ♻️ **Refactor** - Structural improvements
- 📌 **Future** - Out of scope, note for later

### Example Comments Following Best Practices

**IMPORTANT**: Replace `$PR_NUM` with your actual PR number.

```bash
# Security issue (must fix)
./gh-comment add $PR_NUM integration-test-example.js 9 "🔧 API key exposed in code! What do you think about using environment variables instead? This prevents secrets from being committed to git."

# Performance suggestion (question)
./gh-comment add $PR_NUM integration-test-example.js 13:19 "🤔 This nested loop creates O(n²) complexity. Could we consider using a Set or Map for better performance?"

# Architecture feedback (refactor)  
./gh-comment add $PR_NUM integration-test-example.js 2:4 "♻️ What do you think about adding input validation here? It might prevent runtime errors."

# Educational note
./gh-comment add $PR_NUM integration-test-example.js 7 "📝 I noticed this pattern is vulnerable to SQL injection - thought the context might be helpful!"

# Praise for good work (when you fix something)
./gh-comment add $PR_NUM integration-test-example.js 21 "😃 Good module export structure here! Clean and testable."

# Future improvement
./gh-comment add $PR_NUM integration-test-example.js 1 "📌 I'm wondering if we could add TypeScript definitions here in a future iteration?"
```

### Suggestion Syntax Examples

**IMPORTANT**: Only use `[SUGGEST:]` format. The `<<<SUGGEST>>>` multi-line syntax is currently broken.

```bash
# Simple inline suggestion
./gh-comment add $PR_NUM integration-test-example.js 4 "🔧 What about adding null checking? [SUGGEST: if (!users || users.length === 0) return [];]"

# Multi-line suggestion using [SUGGEST:] format (WORKS CORRECTLY)
./gh-comment add $PR_NUM integration-test-example.js 9 "🔧 Let's use environment variables for security: [SUGGEST: const apiKey = process.env.API_KEY || '';
if (!apiKey) throw new Error('API_KEY environment variable required');]"

# SQL injection fix with suggestion
./gh-comment add $PR_NUM integration-test-example.js 7 "🔧 This query is vulnerable to SQL injection. What about using parameterized queries? [SUGGEST: const query = 'SELECT * FROM users WHERE id = ?';
const result = db.query(query, [users[0].id]);]"

# Performance improvement suggestion
./gh-comment add $PR_NUM integration-test-example.js 13:19 "🤔 Could we optimize this nested loop? [SUGGEST: const statusGroups = new Map();
users.forEach(user => {
  if (!statusGroups.has(user.status)) {
    statusGroups.set(user.status, []);
  }
  statusGroups.get(user.status).push(user);
});]"
```

**Note**: All `[SUGGEST:]` content will automatically convert to proper GitHub suggestion blocks that can be applied with one click.

## Integration Test Scenarios

### Scenario 1: Basic Comments Workflow
```bash
# Step 1: Add opening comment explaining the test
gh pr comment $PR_NUM --body "🎯 **Integration Test Demonstration**

This PR demonstrates gh-comment functionality following research-backed code review best practices from \`research/code-review-best-practices.md\`.

Watch as we add various comment types using professional communication patterns."

# Step 2: Add security comment with 🔧 emoji (must fix)
./gh-comment add $PR_NUM integration-test-example.js 9 "🔧 API key exposed in code! What do you think about using environment variables instead? This prevents secrets from being committed to git."

# Step 3: Add performance question with 🤔 emoji
./gh-comment add $PR_NUM integration-test-example.js 13:19 "🤔 This nested loop creates O(n²) complexity. Could we consider using a Set or Map for better performance?"

# Step 4: Add educational note with 📝 emoji
./gh-comment add $PR_NUM integration-test-example.js 7 "📝 I noticed this pattern is vulnerable to SQL injection - thought the security context might be helpful!"

# Step 5: Verify comments
./gh-comment list $PR_NUM
```

### Scenario 2: Suggestion Syntax Testing
```bash
# Step 1: Add simple suggestion (WORKING FORMAT)
./gh-comment add $PR_NUM integration-test-example.js 4 "🔧 What about adding input validation? [SUGGEST: if (!users || !Array.isArray(users)) return [];]"

# Step 2: Add multi-line suggestion (WORKING FORMAT)
./gh-comment add $PR_NUM integration-test-example.js 7 "🔧 Let's use parameterized queries for security: [SUGGEST: const query = 'SELECT * FROM users WHERE id = ?';
const result = db.query(query, [users[0].id]);]"

# Step 3: Add environment variable suggestion
./gh-comment add $PR_NUM integration-test-example.js 9 "🔧 Let's use environment variables: [SUGGEST: const apiKey = process.env.API_KEY;
if (!apiKey) throw new Error('API_KEY required');]"

# Step 4: Verify suggestions render as GitHub suggestion blocks
gh pr view $PR_NUM  # Check in browser that suggestions show properly
```

### Scenario 3: Interactive Features  
```bash
# Step 1: Get comment IDs for reactions/replies
./gh-comment list $PR_NUM

# Step 2: Add reactions to comments (use actual comment IDs)
./gh-comment reply COMMENT_ID --reaction +1
./gh-comment reply COMMENT_ID --reaction heart

# Step 3: Reply to review comments  
./gh-comment reply COMMENT_ID "Great point! I'll implement this fix right away."

# Step 4: Remove reactions to show full functionality
./gh-comment reply COMMENT_ID --remove-reaction +1
```

### Scenario 4: Comment Management
```bash
# Step 1: List all comments with filtering
./gh-comment list $PR_NUM --author $(gh api user -q .login)
./gh-comment list $PR_NUM --type review

# Step 2: Edit an existing comment (use actual comment ID)
./gh-comment edit COMMENT_ID "🔧 **Updated**: API key exposed in code! What do you think about using environment variables instead? This prevents secrets from being committed to git and improves security posture."

# Step 3: Resolve conversation threads
./gh-comment resolve COMMENT_ID --reason "Addressed in latest commit"
```

## Test File Requirements

The test file should contain realistic code issues for demonstration:

```javascript
// Example: showcase-example.js
function calculateUserMetrics(users) {
    // Security issue: Missing input validation
    let totalRevenue = 0;
    
    // Performance issue: Inefficient loops
    for (let i = 0; i < users.length; i++) {
        for (let j = 0; j < users.length; j++) {
            // Logic issues for commenting
        }
    }
    
    // Hardcoded values for refactoring suggestions
    return { total: totalRevenue, tax: totalRevenue * 0.08 };
}

// Database function with security vulnerability
function getUserData(userId) {
    const query = "SELECT * FROM users WHERE id = " + userId;  // SQL injection
    return database.query(query);
}
```

## Verification Steps

After running integration tests, verify:

1. **GitHub PR Interface**:
   - [ ] All comment types display correctly
   - [ ] `[SUGGEST:]` syntax renders as GitHub suggestion blocks
   - [ ] CREG emojis display properly in comments (🔧🤔♻️📝😃📌)
   - [ ] Question-based communication style used throughout
   - [ ] Reviews show up in "Files changed" tab

2. **Command Functionality**:
   - [ ] `./gh-comment list` shows all comment types with context
   - [ ] Filtering works (author, date, type): `./gh-comment list $PR_NUM --author username`
   - [ ] Reactions appear in comment threads
   - [ ] Edit functionality works: `./gh-comment edit COMMENT_ID "new text"`
   - [ ] Reply functionality works: `./gh-comment reply COMMENT_ID "reply message"`

3. **Integration Quality**:
   - [ ] No API errors or rate limiting issues
   - [ ] All `./gh-comment` commands execute successfully  
   - [ ] Professional communication style following `research/code-review-best-practices.md`
   - [ ] All suggestions use working `[SUGGEST:]` format (not broken `<<<SUGGEST>>>`)

## Cleanup Process

```bash
# After testing is complete, clean up the test PR
gh pr close $PR_NUM
git checkout main
git branch -D $(git branch --show-current)
```

**OR** 

```bash
# Leave PR open for demonstration purposes
echo "PR #$PR_NUM left open for showcase: $(gh pr view $PR_NUM --json url -q .url)"
```

## Troubleshooting

### Common Issues

**Rate Limiting**: 
- Use `--inspect` to leave PR open and avoid rapid API calls
- Space out test runs if hitting limits
- Check GitHub API rate limit status

**Authentication**:
- Ensure `gh auth status` shows valid authentication
- Check repository permissions for comment creation
- Verify GitHub CLI is configured correctly

**Comment Validation Errors**:
- Ensure target lines exist in the PR diff
- Check file paths are correct relative to repository root
- Verify PR is open and accepts comments  

### Debug Mode

```bash
# Run with verbose logging for local binary
./gh-comment add $PR_NUM file.js 42 "test comment" --verbose

# Run integration tests with verbose logging  
go run -tags integration . test-integration --force --verbose

# Check integration test logs
ls -la integration-tests/results/
cat integration-tests/results/integration-*.log
```

## Success Criteria

A successful integration test run demonstrates:

- ✅ **Local Development Version**: All tests use `./gh-comment` (not global extension)
- ✅ **Code Review Best Practices**: All comments follow patterns from `research/code-review-best-practices.md`
- ✅ **CREG Emoji System**: Proper use of 🔧🤔♻️📝😃📌 emojis with clear intent
- ✅ **Question-Based Communication**: "What do you think about..." instead of commands
- ✅ **Working Suggestion Syntax**: `[SUGGEST:]` format converts to GitHub suggestion blocks
- ✅ **Complete Comment Lifecycle**: Create, react, reply, edit, resolve
- ✅ **Professional Tone**: Collaborative, educational, and constructive throughout
- ✅ **Clean GitHub Interface**: All features display correctly in browser

## Quick Reference

```bash
# Essential commands for integration testing
go build                                    # Build local version
PR_NUM=$(gh pr view --json number -q .number)  # Get PR number
./gh-comment add $PR_NUM file.js 42 "🔧 comment"  # Add comment
./gh-comment list $PR_NUM                   # List all comments  
./gh-comment reply COMMENT_ID "reply"       # Reply to comment
./gh-comment edit COMMENT_ID "new text"     # Edit comment
```

**Remember**: Always read `research/code-review-best-practices.md` first and follow those communication patterns throughout the integration test process.