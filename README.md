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

# List all comments on a PR
gh comment list 123

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

```
add                  Add general PR comments
review               Create line-specific code reviews  
list                 List all PR comments
edit                 Edit existing comments
react                Add emoji reactions
batch                Process multiple comments from YAML
lines                Show which lines can be commented on
review-reply         Reply to review comments
prompts              Get code review prompt templates
export               Export comments to JSON
config               Manage configuration settings
close-pending-review Submit pending reviews from GitHub UI
```

Run `gh comment --help` for comprehensive documentation and examples.

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