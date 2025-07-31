# gh-comment

Strategic line-specific PR commenting for GitHub CLI

## Overview

`gh-comment` is the first GitHub CLI extension designed for strategic, line-specific PR commenting workflows. It fills a genuine gap in the GitHub CLI ecosystem by enabling targeted comments on specific lines in pull requests, with smart tone transformation and batch operations.

## Features

- **Line-specific comments**: Add comments to individual lines or line ranges
- **Smart tone transformation**: Automatic casual/formal/technical tone adjustment
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

### Tone Transformation

```bash
# Casual tone (default) - transforms formal language to casual
gh comment add 123 src/api.js 42 "This implements error fingerprinting" --tone casual
# → "this implements error fingerprinting - much cleaner approach"

# Formal tone - keeps professional language
gh comment add 123 src/api.js 42 "This implements error fingerprinting" --tone formal
# → "This implements error fingerprinting"

# Technical tone - adds technical precision
gh comment add 123 src/api.js 42 "This implements error fingerprinting" --tone technical
# → "This implements error fingerprinting."
```

### Options

```bash
# Dry run - preview without posting
gh comment add --dry-run 123 src/api.js 42 "test comment"

# Verbose mode - show API details
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

### Review Workflow
```bash
# Quick explanatory comments during review
gh comment add src/auth.js 42 "this handles the oauth callback edge case we discussed"
gh comment add src/utils/validation.js 15 "regex pattern updated for new email format requirements"
```

## Roadmap

- [ ] **Review-based comments**: Create professional grouped reviews
- [ ] **Batch operations**: Process multiple comments from config files
- [ ] **Interactive mode**: GUI-like comment selection
- [ ] **Comment templates**: Predefined comment patterns
- [ ] **Analytics**: Comment usage reporting

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
