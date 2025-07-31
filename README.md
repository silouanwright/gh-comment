# gh-comment

Strategic line-specific PR commenting for GitHub CLI (optimized for AI)

## Overview

`gh-comment` is the first GitHub CLI extension designed for strategic, line-specific PR commenting workflows. It fills a genuine gap in the GitHub CLI ecosystem by enabling targeted comments on specific lines in pull requests, with smart tone transformation and batch operations. Built specifically for AI assistants and automated workflows.

## Features

- **AI-optimized design**: Specifically built for usage with AI assistants and automated workflows
- **Line-specific comments**: Add comments to individual lines or line ranges
- **Multi-line comment support**: Shell native and --message flags for complex comments
- **List all comments**: View all PR comments with diff context and author info
- **Reply to comments**: Respond to specific comments with threaded replies
- **Edit existing comments**: Update comment text to fix mistakes or add context
- **Emoji reactions**: Quick acknowledgments with GitHub reactions
- **Dry-run mode**: Preview comments before posting
- **Auto-detection**: Automatically detect current repository and PR
- **Verbose mode**: Detailed API interaction logging

## AI Workflows (Optimized for AI)

`gh-comment` was specifically designed for AI assistants like Claude, ChatGPT, and other automated workflows. Here's why AI tools work exceptionally well with this extension:

### **Why AI Loves gh-comment**

- **Immediate feedback**: Every command provides instant results, allowing AI to adjust strategy in real-time
- **Rich documentation**: Comprehensive `--help` text enables rapid AI learning of command patterns
- **Structured output**: Verbose mode provides detailed context for AI decision-making
- **Explicit operations**: Clear, predictable commands reduce AI interpretation overhead
- **Safety features**: Dry-run mode allows AI to preview actions before execution

### **Built for Speed**

Unlike complex MCP servers or API wrappers, `gh-comment` gives AI direct access to GitHub's PR commenting system with zero interpretation layers. This means faster execution and more reliable results.

## Prerequisites

Requires [GitHub CLI](https://cli.github.com/) installed and authenticated:

```bash
# Install GitHub CLI
brew install gh  # macOS
# or see: https://github.com/cli/cli#installation

# Authenticate
gh auth login
```

## Installation

```bash
gh extension install silouanwright/gh-comment
```

## Usage

### Basic Line Comments

```bash
# Add single-line comment
gh comment add 123 src/api.js 42 "this handles the rate limiting edge case"

# Add multi-line comment (shell native)
gh comment add 123 src/api.js 42 "Line 1
Line 2
Line 3"

# Add multi-line comment using --message flags (AI-friendly)
gh comment add 123 src/api.js 42 --message "First paragraph" --message "Second paragraph"

# Auto-detect PR from current branch
gh comment add src/api.js 42 "quick fix needed"
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

### Edit Comments

```bash
# Edit an existing comment
gh comment edit 2246362251 "Updated comment with better explanation"

# Edit with multi-line content (AI-friendly)
gh comment edit 2246362251 --message "First paragraph" --message "Second paragraph"

# Edit with multi-line content (shell native)
gh comment edit 2246362251 "Line 1
Line 2
Line 3"

# Test edit before applying
gh comment edit --dry-run 2246362251 "Test message"
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

# 4. Edit comments to fix mistakes or add context
gh comment edit 2246362251 "Updated: Good point, I'll fix this in the next commit!"

# 5. Add your own explanatory comments
gh comment add src/auth.js 42 "this handles the oauth callback edge case we discussed"

# 6. Verify all feedback is addressed
gh comment list 998 --author reviewer-username
```

## Limitations

### GitHub API Constraints

**Incremental Review Comments Not Supported**

While `gh-comment` supports both individual comments and batch review creation, it **does not support incremental review comments** (adding comments one-by-one to a pending review). This is a **GitHub API limitation**, not a design choice.

**What Works:**
- ✅ **Individual comments**: `gh comment add` posts comments immediately
- ✅ **Batch reviews**: `gh comment add-review` creates complete reviews with multiple comments in one API call
- ✅ **Submit pending reviews**: `gh comment submit-review` submits existing pending reviews

**What Doesn't Work:**
- ❌ **Incremental review building**: Cannot add comments one-by-one to a pending review
- ❌ **Mixed workflow**: Cannot combine individual comments with pending review comments

**Technical Explanation:**

GitHub's API has a fundamental constraint: **only one pending review per user per PR** is allowed. When a pending review exists, the API rejects attempts to:
- Create new reviews (returns 422: "Pull request review cannot be created")
- Add individual comments to existing pending reviews via the standard comments endpoint
- Mix review comments with standalone comments

This affects many GitHub integrations and has been a long-standing developer pain point:
- [GitHub Community Discussion #24854](https://github.com/orgs/community/discussions/24854): "Cannot add comments to existing pending review"
- [Stack Overflow #71421045](https://stackoverflow.com/questions/71421045/): "How to add comments to pending GitHub review via API"
- [GitHub API Docs Issue](https://docs.github.com/en/rest/pulls/comments): API documentation acknowledges this limitation

**Workarounds:**

1. **Use batch reviews**: Plan all comments and create them together with `add-review`
2. **Use individual comments**: Post comments immediately with `add` (not part of reviews)
3. **Submit and restart**: If you (or anyone) already started a review in the GitHub web interface, you must submit that pending review first using `gh comment submit-review`, then you can create new reviews with `add-review`

**Why This Matters for AI Workflows:**

AI tools often want to add comments incrementally as they analyze code. This API limitation forces a choice between:
- **Immediate feedback**: Comments appear right away but aren't grouped in reviews
- **Batched feedback**: Comments are professionally grouped but require planning all feedback upfront

`gh-comment` provides both options clearly, letting you choose the right approach for your workflow.

## Roadmap

- [ ] **Advanced filtering**: Filter comments by status, author, date, resolved state
- [ ] **Add tests**: Unit tests for core functionality and edge cases
- [ ] **What do you want to see?** [Let me know!](https://github.com/silouanwright/gh-comment/issues)

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
