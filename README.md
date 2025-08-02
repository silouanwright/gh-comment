# gh-comment 💬

Strategic line-specific PR commenting for GitHub CLI (optimized for AI)

## Overview

`gh-comment` is the first GitHub CLI extension designed for comprehensive PR comment management. It provides a unified system for both general PR discussion and line-specific code review comments, filling a genuine gap in the GitHub CLI ecosystem. Features smart suggestion expansion, complete comment visibility, and universal reply capabilities. Built specifically for AI assistants and automated workflows.

## Features

- 🤖 **AI-optimized design**: Specifically built for usage with AI assistants and automated workflows
- 📍 **Line-specific comments**: Add comments to individual lines or line ranges
- 🚀 **Suggestion expansion**: Simple `[SUGGEST: code]` syntax automatically expands to GitHub suggestion blocks
- 📝 **Multi-line comment support**: Shell native and --message flags for complex comments
- 📋 **Unified comment listing**: View ALL PR comments (both general discussion and line-specific) with diff context
- 💬 **Universal reply system**: Reply to any comment type (issue or review) with automatic API selection
- ✏️ **Edit existing comments**: Update comment text to fix mistakes or add context
- 😀 **Emoji reactions**: Quick acknowledgments with GitHub reactions
- 🧪 **Dry-run mode**: Preview comments before posting
- 🔍 **Auto-detection**: Automatically detect current repository and PR
- 🔊 **Verbose mode**: Detailed API interaction logging
- 🐠 **Cross-shell compatibility**: Works identically in Fish, Bash, Zsh, and PowerShell

## 🤖 AI Workflows (Optimized for AI)

`gh-comment` was specifically designed for AI assistants like Claude, ChatGPT, and other automated workflows. Here's why AI tools work exceptionally well with this extension:

### **Why AI Loves gh-comment**

- **Immediate feedback**: Every command provides instant results, allowing AI to adjust strategy in real-time
- **Rich documentation**: Comprehensive `--help` text enables rapid AI learning of command patterns
- **Structured output**: Verbose mode provides detailed context for AI decision-making
- **Explicit operations**: Clear, predictable commands reduce AI interpretation overhead
- **Safety features**: Dry-run mode allows AI to preview actions before execution

### **🚀 Built for Speed**

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

### 🚀 Installation

**One-command install** with automatic platform detection:
```bash
gh extension install silouanwright/gh-comment
```

GitHub CLI automatically:
- Detects your platform (macOS, Linux, Windows)
- Downloads the correct pre-compiled binary
- Falls back to source compilation if needed

**That's it!** No need to specify your platform or architecture.

### 🔄 Upgrading

```bash
# Upgrade to latest version
gh extension upgrade gh-comment

# Or reinstall
gh extension remove gh-comment
gh extension install silouanwright/gh-comment
```

### 🗑️ Uninstalling

```bash
gh extension remove gh-comment
```

### ✅ Verification

```bash
# Check installation
gh comment --version

# Test basic functionality
gh comment --help
```

## Comment System

`gh-comment` provides unified management for both types of PR comments:

- **General PR Comments**: Discussion like "LGTM", questions, or general feedback
- **Line-Specific Review Comments**: Code review comments tied to specific files and lines

`gh comment list` shows both types with full context, and `gh comment reply` works with both using the `--type` flag to automatically select the correct GitHub API endpoint.

## Usage

### 🚀 GitHub Suggestion Expansion

`gh-comment` includes innovative **suggestion expansion** that transforms simple syntax into GitHub's complex suggestion markdown:

```bash
# Simple inline suggestions - expands automatically
gh comment add 123 src/api.js 42 "Try [SUGGEST: const result = input?.value || 'default'] instead"

# Multiline suggestions with intuitive syntax
gh comment add 123 src/api.js 42 "Consider this approach:

<<<SUGGEST
function processData(input) {
  if (!input) return null;
  return input.map(item => item.value);
}
SUGGEST>>>

This handles the null case better."

# Disable expansion when needed
gh comment add 123 src/api.js 42 "Raw text: [SUGGEST: code]" --no-expand-suggestions
```

**What happens behind the scenes:**

Your input:
```
[SUGGEST: const result = data.filter(x => x.active);]
```

Becomes this GitHub suggestion block:
````
```suggestion
const result = data.filter(x => x.active);
```
````

- Works identically across Fish, Bash, Zsh, and PowerShell
- No complex shell escaping required!

### Creating Comments

#### Line-Specific Comments (gh-comment extension)
```bash
# Add single-line comment to specific code line
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

#### General PR Comments (native GitHub CLI)
```bash
# Add general discussion comment (not tied to specific code)
gh pr comment 123 --body "LGTM! Just a few minor suggestions below."

# Multi-line general comment
gh pr comment 123 --body "Thanks for this PR!

I've reviewed the changes and they look good overall.
Just a couple of questions about the implementation."
```

**Key Distinction:**
- **Line-specific**: Use `gh comment add <pr> <file> <line>` (gh-comment extension)
- **General discussion**: Use `gh pr comment <pr> --body` (native GitHub CLI)

### List All Comments (Advanced Filtering System)

`gh comment list` shows **ALL** comments on a PR with powerful filtering capabilities:

```bash
# List all comments on a PR with diff context
gh comment list 123

# Filter by author (supports wildcards and partial matching)
gh comment list 123 --author octocat           # Exact match
gh comment list 123 --author "octo*"           # Wildcard match
gh comment list 123 --author "*@company.com"   # Email domain match

# Filter by comment type
gh comment list 123 --type issue    # General PR comments only
gh comment list 123 --type review   # Line-specific review comments only

# Filter by date range
gh comment list 123 --since "2024-01-01"       # Absolute date
gh comment list 123 --since "3 days ago"       # Relative date
gh comment list 123 --until "1 week ago"       # Before date
gh comment list 123 --since "2 days ago" --until "1 day ago"  # Date range

# Filter by status (open/resolved/all)
gh comment list 123 --status open     # Active comments only
gh comment list 123 --status resolved # Resolved comments only

# Combine multiple filters
gh comment list 123 --author "dev*" --type review --since "1 week ago"

# Auto-detect PR from current branch
gh comment list

# Clean output for human reading (hides URLs/IDs)
gh comment list 123 --quiet
```

**Output includes two distinct sections:**
- **💬 General PR Comments**: Discussion comments like "LGTM" or general feedback
- **📋 Review Comments**: Line-specific code review comments with file/line context

### Reply to Comments (Universal System)

`gh comment reply` works with **both** comment types using the `--type` flag:

```bash
# Reply to line-specific review comment (default)
gh comment reply 2246362251 "Good point, I'll fix this!"

# Reply to general PR discussion comment
gh comment reply 3141344022 "Thanks for the feedback!" --type issue

# Add reaction to any comment type
gh comment reply 2246362251 --reaction +1

# Reply and resolve conversation (review comments only)
gh comment reply 2246362251 "Fixed in latest commit" --resolve

# Remove reaction from any comment
gh comment reply 2246362251 --remove-reaction +1
```

**Comment Types:**
- `--type review` (default): Line-specific code review comments - creates threaded replies
- `--type issue`: General PR discussion comments - creates new top-level comments

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

### Key Feature: Unified Comment View

Unlike other tools, `gh comment list` shows **ALL** PR comments with full context:

```
📝 Comments on PR #123 (5 total)

💬 General PR Comments (2)
──────────────────────────────────────────────────
[1] 👤 reviewer • 2 hours ago
   LGTM! Just a few minor suggestions below.
   🔗 https://github.com/owner/repo/pull/123#issuecomment-3141344022

[2] 👤 maintainer • 1 hour ago
   Thanks for the contribution! Please address the review feedback.
   🔗 https://github.com/owner/repo/pull/123#issuecomment-3141355432


📋 Review Comments (3)
──────────────────────────────────────────────────
[1] 👤 reviewer • 2 hours ago
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
   🔗 https://github.com/owner/repo/pull/123#discussion_r2246362251
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

### Complete Review Workflow (Unified System)
```bash
# 1. Someone reviews your PR and adds both general and line-specific comments
# 2. List ALL feedback with full context and comment IDs
gh comment list 998

# 3. Reply to general PR discussion
gh comment reply 3141344022 "Thanks for the overall feedback!" --type issue

# 4. Reply to specific line-specific review comments
gh comment reply 2246362251 "Good point, I'll fix this!" --type review
gh comment reply 2246362252 --reaction +1

# 5. Resolve conversations after addressing feedback
gh comment reply 2246362251 "Fixed in commit abc123" --resolve

# 6. Edit comments to fix mistakes or add context
gh comment edit 2246362251 "Updated: Fixed in commit abc123 with proper error handling"

# 7. Add your own explanatory comments
gh comment add src/auth.js 42 "this handles the oauth callback edge case we discussed"

# 8. Verify all feedback is addressed (both types)
gh comment list 998 --author reviewer-username
```

## Limitations

**[Incremental Review Comments](https://github.com/orgs/community/discussions/168380) Not Supported**

Queuing individual comments as part of a review would be helpful, but GitHub's API constraint allows only one pending review per user per PR. This affects many integrations and remains a long-standing developer pain point.

**Workarounds:**
- **Batch reviews**: Plan all feedback upfront with `add-review`
- **Individual comments**: Post immediately with `add` (not grouped in reviews)
- **Submit and restart**: Complete any pending review first, then create new ones

**Help us get this feature added!** Please visit https://github.com/orgs/community/discussions/168380 and upvote the discussion to gain more traction with GitHub's product team.

## Roadmap

### 🚧 In Progress (August 2025)
- [ ] **Test Coverage to 80%**: Refactoring commands with dependency injection for better testing
- [ ] **Cross-Platform Testing**: Ensuring consistent behavior across Windows, macOS, and Linux

### 📋 Planned Features
- [ ] **GitLab-style line offset syntax**: Support `[SUGGEST:+2: code]` and `[SUGGEST:-1: code]` for relative line positioning in suggestions
- [ ] **Advanced filtering**: Filter comments by status, author, date, resolved state
- [ ] **Configuration file support**: Default flags and repository settings
- [ ] **Template system**: Reusable comment patterns and workflows
- [ ] **Batch operations**: Apply operations to multiple comments at once
- [ ] **Export functionality**: Export comments to various formats (JSON, CSV, Markdown)
- [ ] **What do you want to see?** [Let me know!](https://github.com/silouanwright/gh-comment/issues)

### ✅ Recently Completed
- [x] **API Abstraction Layer**: Clean separation of GitHub API calls for better testing (August 2025)
- [x] **Performance Regression Testing**: Automated benchmark comparison on every PR (August 2025)
- [x] **Pre-commit Hooks**: Automated code quality checks on every commit (August 2025)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Development

### Setup

```bash
# Clone and build
git clone https://github.com/silouanwright/gh-comment
cd gh-comment
go build

# Install locally
gh extension install .

# Test basic functionality
./gh-comment add --dry-run --repo owner/repo 123 file.js 42 "test comment"
```

### Testing

The project includes comprehensive testing with multiple layers:

```bash
# Run all tests
go test ./...

# Run only fast unit tests
go test ./... -short

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests (uses testscript)
go test ./... -run TestIntegration

# Run E2E tests (requires GitHub token and test repo)
export GH_TOKEN="your-token"
export GH_E2E_REPO="owner/repo"
export GH_E2E_PR="123"
export RUN_E2E_TESTS="true"
go test ./cmd -run TestE2E

# Run benchmarks and compare performance
./scripts/benchmark.sh                    # Run benchmarks
./scripts/benchmark.sh compare main       # Compare with main branch
./scripts/benchmark.sh profile            # Generate CPU/memory profiles
```

**Test Architecture:**
- **Unit Tests**: Core function testing with mocks
- **Integration Tests**: CLI workflow testing with testscript
- **Fuzz Tests**: Edge case discovery with Go 1.18+ fuzzing
- **E2E Tests**: Real GitHub API testing with safety measures
- **Benchmark Tests**: Performance monitoring with regression detection

**Recent Testing Improvements (August 2025):**
- ✅ GitHub API abstraction layer for better testability
- ✅ Automated performance regression testing in CI/CD
- ✅ Local benchmark comparison script for developers

See `docs/testing/TESTING.md` and `docs/testing/E2E_TESTING.md` for detailed testing documentation.

### Code Quality

Pre-commit hooks automatically ensure code quality:

```bash
# Install pre-commit (one-time setup)
brew install pre-commit
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install hooks
pre-commit install
pre-commit install --hook-type commit-msg
```

**Automatic Checks on Every Commit:**
- Go build verification
- Unit test execution
- Code formatting (`go fmt`)
- Static analysis (`go vet`)
- Dependency management (`go mod tidy`)
- Conventional commit message format

See `docs/development/PRE_COMMIT_SETUP.md` for complete setup instructions.

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**gh-comment** - Making PR reviews more strategic, one line at a time.
