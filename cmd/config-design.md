# Configuration File Design Specification

## Overview
This specification defines the configuration file format for `gh-comment` to support default flags, repository settings, and user preferences.

## Configuration File Locations

### Priority Order (highest to lowest)
1. `--config` flag specified file
2. `.gh-comment.yaml` in current directory
3. `.gh-comment.yaml` in repository root (if in git repo)
4. `~/.config/gh-comment/config.yaml` (user config)
5. `~/.gh-comment.yaml` (legacy user config)

## Configuration File Format (YAML)

```yaml
# Default repository and PR settings
defaults:
  repository: "owner/repo"  # Default repository
  author: "username"        # Default author filter
  pr: 0                     # Default PR (0 = auto-detect)
  
# Command behavior settings
behavior:
  dry_run: false           # Default --dry-run behavior
  verbose: false           # Default --verbose behavior
  validate: true           # Default --validate behavior
  no_expand_suggestions: false  # Default suggestion expansion

# Display preferences
display:
  format: "table"          # Output format: table, json, quiet
  color: "auto"            # Color output: auto, always, never
  quiet: false             # Default --quiet behavior
  
# Filtering defaults (for list command)
filters:
  status: "all"            # Default status filter: open, resolved, all
  type: "all"              # Default type filter: issue, review, all
  since: ""                # Default since filter
  until: ""                # Default until filter

# Review defaults
review:
  event: "COMMENT"         # Default review event: APPROVE, REQUEST_CHANGES, COMMENT

# API settings
api:
  timeout: 30              # API timeout in seconds
  retry_count: 3           # Number of retries for failed requests
  rate_limit_buffer: 10    # Seconds to wait before rate limit reset

# Suggestion syntax settings
suggestions:
  expand_by_default: true  # Whether to expand [SUGGEST:] syntax by default
  max_offset: 999          # Maximum allowed offset for [SUGGEST:+N:]
  
# Custom aliases for common commands
aliases:
  approve: "review --event APPROVE"
  reject: "review --event REQUEST_CHANGES"
  lgtm: "add 'LGTM! Looks good to merge 🚀'"

# Template settings
templates:
  default_review_body: "Code review complete"
  default_approval_message: "LGTM! Ready to merge"
```

## JSON Format Support

The same configuration can be provided in JSON format:

```json
{
  "defaults": {
    "repository": "owner/repo",
    "author": "username",
    "pr": 0
  },
  "behavior": {
    "dry_run": false,
    "verbose": false,
    "validate": true
  },
  "display": {
    "format": "table",
    "color": "auto",
    "quiet": false
  },
  "filters": {
    "status": "all",
    "type": "all"
  },
  "review": {
    "event": "COMMENT"
  }
}
```

## Environment Variable Override

Configuration values can be overridden by environment variables:

- `GH_COMMENT_REPO` -> defaults.repository
- `GH_COMMENT_AUTHOR` -> defaults.author
- `GH_COMMENT_DRY_RUN` -> behavior.dry_run
- `GH_COMMENT_VERBOSE` -> behavior.verbose
- `GH_COMMENT_FORMAT` -> display.format
- `GH_COMMENT_COLOR` -> display.color

## Flag Priority

Priority order (highest to lowest):
1. Command-line flags
2. Environment variables
3. Configuration file values
4. Built-in defaults

## Configuration Validation

### Required Validations
1. **Repository format**: Must match `owner/repo` pattern if specified
2. **Enum values**: status, type, event, format, color must be valid values
3. **Numeric ranges**: timeout, retry_count, max_offset must be positive
4. **Path validation**: Template paths must exist if specified

### Warning Validations
1. **Large timeouts**: Warn if timeout > 300 seconds
2. **High retry counts**: Warn if retry_count > 10
3. **Large offsets**: Warn if max_offset > 100

## Configuration Commands

### Generate Default Config
```bash
gh comment config init [--format yaml|json] [--global]
```

### Show Current Config
```bash
gh comment config show [--effective] [--source]
```

### Validate Config
```bash
gh comment config validate [config-file]
```

### Set Config Values
```bash
gh comment config set defaults.repository owner/repo
gh comment config set display.format json
```

## Implementation Notes

### Config Loading Order
1. Parse command-line flags for `--config`
2. Search for config files in priority order
3. Merge configs (file -> env vars -> flags)
4. Validate final configuration
5. Apply defaults for missing values

### Error Handling
- Invalid YAML/JSON: Show syntax error with line number
- Missing files: Silently skip (except when --config specified)
- Invalid values: Show validation error and exit
- Permission errors: Show clear error message

### Backwards Compatibility
- All existing flags continue to work
- Configuration is purely additive
- No breaking changes to existing behavior

## Example Usage

### Project-specific config
```yaml
# .gh-comment.yaml in project root
defaults:
  repository: "myorg/myproject"
  
filters:
  author: "team-*"
  
review:
  event: "REQUEST_CHANGES"  # Default to request changes for this project
```

### User-wide preferences
```yaml
# ~/.config/gh-comment/config.yaml
display:
  format: "json"
  color: "always"
  
behavior:
  verbose: true
  
suggestions:
  expand_by_default: true
```

This design provides comprehensive configuration support while maintaining backward compatibility and following standard conventions.