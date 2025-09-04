# gh-comment

Add comments to GitHub PRs from the command line.

**GitHub CLI Extension** - Extends the `gh` command with PR commenting capabilities. Learn more about [GitHub CLI extensions](https://cli.github.com/manual/gh_extension) (including how to upgrade, remove, and manage extensions).

## Installation

```bash
gh extension install silouanwright/gh-comment
```

Requires [GitHub CLI](https://cli.github.com/) (`gh auth login` to authenticate).

## Quick Start

```bash
# Add a general PR comment
gh comment add 123 "LGTM! Ready to merge"

# Add a line-specific code review comment  
gh comment review 123 --comment src/api.js:42:"Add error handling here"

# List unresolved comments on a PR (default behavior)
gh comment list 123

# List ALL comments including resolved ones
gh comment list 123 --show-all

# React to a comment
gh comment react 2254752948 +1
```

## Common Workflows

### Basic Code Review
```bash
# Review with multiple line comments
gh comment review 123 "Found some issues" \
  --comment src/api.js:42:"Missing error handling" \
  --comment src/auth.js:15:"Use bcrypt instead" \
  --event REQUEST_CHANGES
```

### Working with Comments
```bash
# Edit a comment
gh comment edit 2254752948 "Updated: Fixed the issue"

# Add a code suggestion
gh comment add 123 "Try this: [SUGGEST: return data || default]"

# Reply to a review comment
gh comment review-reply 2254752948 "Good catch, fixed!"
```

### Batch Operations
```yaml
# review.yaml
pr: 123
comments:
  - file: src/api.js
    line: 42
    message: "Add rate limiting"
  - file: src/auth.js
    line: 15
    message: "Update to OAuth2"
```

```bash
gh comment batch 123 review.yaml
```

## Commands

### Core Commands
```bash
# General PR discussion comments
gh comment add <pr> <message>                    # Add general discussion comment
gh comment add <pr> <file> <line> <message>      # Add line-specific issue comment

# Line-specific code reviews  
gh comment review <pr> [body] --comment <file:line:message> --event <APPROVE|REQUEST_CHANGES|COMMENT>
gh comment review-reply <comment-id> <message>   # Reply to review comments

# Comment management (shows unresolved by default, use --show-all for all)
gh comment list <pr> [--author] [--since] [--type] [--show-all] [--quiet]
gh comment edit <comment-id> <new-message>       # Modify existing comments
gh comment react <comment-id> <emoji>            # Add/remove emoji reactions
```

### Advanced Features
```bash
# Batch operations from YAML
gh comment batch <pr> <config.yaml> [--dry-run] [--verbose]

# Workflow helpers
gh comment lines <pr> <file>                     # Show commentable lines
gh comment close-pending-review <pr> <message>   # Submit pending reviews from GitHub UI
gh comment prompts [list|<template>]             # AI code review templates
gh comment export <pr>                           # Export comments to JSON
```

### Global Flags
```bash
-p, --pr <number>        # PR number (auto-detects from branch)
-R, --repo <owner/repo>  # Repository (auto-detects from current directory)
    --dry-run            # Preview changes without executing
-v, --verbose            # Show detailed API interactions  
    --validate           # Validate lines exist in diff (default: false)
```

Run `gh comment <command> --help` for detailed examples and options.

## Key Features

- **Line-specific comments**: Comment on exact lines or ranges
- **Auto-detects PR**: Works without specifying PR number  
- **Suggestion syntax**: `[SUGGEST: code]` expands to GitHub suggestions
- **Batch processing**: Review multiple files systematically

## Documentation

- **[Advanced Usage](docs/ADVANCED_USAGE.md)** - Power user features, automation workflows, and complex scenarios
- **[Contributing](docs/CONTRIBUTING.md)** - Development setup, testing, and contribution guidelines

## Troubleshooting

**"line does not exist in diff"**: Use `gh comment lines <pr> <file>` to see commentable lines. Only changed lines in a PR can have review comments.

**"Resource not accessible"**: Run `gh auth status` to check authentication.

See `gh comment --help` for detailed usage, examples, and all available options.