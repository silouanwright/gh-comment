# Batch Configuration Examples

This directory contains example YAML configuration files for the `gh comment batch` command.

## Files

- **review-config.yaml**: Comprehensive code review with multiple line-specific comments and review submission
- **security-checklist.yaml**: Security audit configuration with vulnerability-focused comments
- **bulk-comments.yaml**: Mixed comment types including both general discussion and line-specific feedback

## Usage

```bash
# Use an example configuration
gh comment batch 123 examples/review-config.yaml

# Validate configuration without executing
gh comment batch 123 examples/security-checklist.yaml --dry-run

# Run with verbose output
gh comment batch 123 examples/bulk-comments.yaml --verbose
```

## Configuration Format

```yaml
# Optional: Override PR number (can be set via command line)
pr: 123
repo: owner/repository  # Optional: Override repository

# Optional: Create a review with summary
review:
  body: "Review summary message"
  event: APPROVE | REQUEST_CHANGES | COMMENT

# Comments to add
comments:
  - file: path/to/file.js
    line: 42                    # Single line comment
    message: "Comment text"
    type: review               # review or issue

  - file: path/to/file.js
    range: "10-15"             # Range comment
    message: "Comment text" 
    type: review

  - file: ""                   # Empty file for general PR comments
    message: "General discussion"
    type: issue
```