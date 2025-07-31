# gh-comment

Strategic line-specific PR commenting for GitHub CLI

## Overview

`gh-comment` is the first GitHub CLI extension designed for strategic, line-specific PR commenting workflows. It fills a genuine gap in the GitHub CLI ecosystem by enabling targeted comments on specific lines in pull requests, with smart tone transformation and batch operations.

## Features

- **Line-specific comments**: Add comments to individual lines or line ranges
- **List all comments**: View all PR comments with diff context and author info
- **Reply to comments**: Respond to specific comments with threaded replies
- **Emoji reactions**: Quick acknowledgments with GitHub reactions
- **Dry-run mode**: Preview comments before posting
- **Auto-detection**: Automatically detect current repository and PR
- **Verbose mode**: Detailed API interaction logging

## Installation

```bash
gh extension install silouanwright/gh-comment
```

## Usage

### Basic Line Comments

```bash
# Add comment to a specific line
gh comment add 123 src/api.js 42 "this handles the rate limiting edge case"

# Add comment to a line range
gh comment add 123 src/auth.js 15:20 "updated auth flow for better security"

# Auto-detect PR from current branch
gh comment add src/api.js 42 "this fixes the jest scoping issue"
```

### List Comments

```bash
# List all comments on a PR with diff context
gh comment list 123

# List comments from specific author
gh comment list 123 --author octocat

# Auto-detect PR from current branch
gh comment list

# Clean output for human reading (hides URLs/IDs)
gh comment list 123 --quiet
```

### Reply to Comments

```bash
# Reply to a specific comment (use comment ID from verbose output)
gh comment reply 2246362251 "Good catch, thanks for the feedback!"

# Add a reaction to show appreciation
gh comment reply 2246362251 --reaction +1

# Reply with both message and reaction
gh comment reply 2246362251 "Fixed in latest commit!" --reaction heart

# Quick acknowledgment
gh comment reply 2246362251 --reaction eyes
```

### Key Feature: Diff Context

Unlike other tools, `gh comment list` shows the **exact code context** that comments refer to:

```
📍 Line-Specific Comments (1)
──────────────────────────────────────────────────
[1] 👤 reviewer • 5 minutes ago
📁 src/api.js:L42
📝 Code Context:
   🔹 @@ -40,6 +40,8 @@ function handleRequest(req) {
      if (!req.user) {
        throw new Error('Unauthorized');
      }
   ➕ +  // Add rate limiting check
   ➕ +  checkRateLimit(req.user.id);
      return processRequest(req);

   This needs error handling for rate limit failures
```

This makes it **perfect for AI-assisted code reviews** - no guessing what code the comment refers to!

> **🤖 AI-First Design**: By default, `gh comment list` shows URLs and comment IDs that AI needs to reply to comments. Use `--quiet` for human-only reading.

### Options

```bash
# Dry run - preview without posting
gh comment add --dry-run 123 src/api.js 42 "test comment"

# Verbose mode - show detailed API interactions
gh comment add --verbose 123 src/api.js 42 "test comment"

# Specify repository explicitly
gh comment add --repo owner/repo 123 src/api.js 42 "test comment"

# Skip diff validation
gh comment add --validate=false 123 src/api.js 42 "test comment"
```

## Examples

### Code Migration Comments
```bash
# Explain migration decisions
gh comment add 998 packages/knit/package.json 74 "switched from bugsnag to datadog rum as peerdependency"

# Document architectural changes
gh comment add 998 src/api/client.js 120:135 "centralized error handling - replaces scattered try-catch blocks"
```

### Complete Review Workflow
```bash
# 1. Someone reviews your PR and adds comments
# 2. List all feedback with diff context and comment IDs
gh comment list 998 --verbose

# 3. Reply to specific feedback in threaded conversations
gh comment reply 2246362251 "Good point, I'll fix this!"
gh comment reply 2246362252 --reaction +1

# 4. Add your own explanatory comments
gh comment add src/auth.js 42 "this handles the oauth callback edge case we discussed"

# 5. Verify all feedback is addressed
gh comment list 998 --author reviewer-username
```

## Roadmap

- [x] **Line-specific comments**: Add targeted comments to specific lines ✅
- [x] **List comments with context**: View all PR feedback with diff context ✅
- [x] **Reply to comments**: Threaded replies and emoji reactions ✅
- [ ] **Review-based comments**: Create professional grouped reviews
- [ ] **Batch operations**: Process multiple comments from config files
- [ ] **Interactive mode**: GUI-like comment selection
- [ ] **Comment templates**: Predefined comment patterns
- [ ] **Resolve comments**: Mark conversations as resolved

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Development

```bash
# Clone and build
git clone https://github.com/silouanwright/gh-comment
cd gh-comment
go build

# Install locally
gh extension install .

# Test
./gh-comment add --dry-run --repo owner/repo 123 file.js 42 "test comment"
```

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**gh-comment** - Making PR reviews more strategic, one line at a time.
