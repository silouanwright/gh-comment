# Integration Tests

This directory contains the integration testing framework for `gh-comment` using a "dogfooding" approach.

## Overview

The integration tests work by:
1. Creating real pull requests on this repository
2. Exercising all `gh-comment` functionality against the live PRs
3. Validating results through multiple methods
4. Cleaning up test artifacts automatically

## Usage

```bash
# Run all integration tests with auto-cleanup
go run . test-integration

# Run specific scenario and inspect results
go run . test-integration --scenario=comments --inspect

# Run tests without cleanup for debugging
go run . test-integration --no-cleanup
```

## Available Scenarios

- **comments**: Basic line and range commenting
- **reviews**: Review comment creation and submission
- **reactions**: Reactions and replies to comments
- **batch**: YAML-based batch operations
- **suggestions**: Suggestion syntax testing

## Directory Structure

```
integration-tests/
├── README.md              # This file
├── scenarios/              # (Reserved for future scenario-specific files)
├── scripts/                # (Reserved for helper scripts)
├── templates/              # Test file templates
│   └── dummy-code.js      # JavaScript template with intentional issues
└── results/               # Test execution logs and temporary files
    └── integration-*.log  # Timestamped log files
```

## How It Works

### 1. PR Creation
- Creates a unique branch with timestamp
- Copies template file with intentional code issues
- Commits and pushes the branch
- Opens a pull request via `gh pr create`

### 2. Test Execution
- Runs all `gh-comment` commands against the live PR
- Tests line comments, range comments, reviews, reactions, and replies
- Validates results using both command output and GitHub API

### 3. Validation Methods
- **Command Output**: Verifies commands succeed and produce expected output
- **GitHub API**: Cross-validates results using `gh` CLI API calls
- **Cross-Command**: Adds comments then verifies with `list` command

### 4. Cleanup
- Closes the test PR
- Deletes local and remote test branches
- Removes temporary files
- Logs all operations for debugging

## Safety Features

- **Confirmation prompts** for destructive operations
- **Rate limiting** to respect GitHub API limits
- **Comprehensive logging** of all operations
- **Automatic cleanup** to prevent repository pollution
- **Error handling** with detailed failure information

## Flags

- `--cleanup`: Auto-close PR after tests (default: true)
- `--inspect`: Leave PR open for manual inspection
- `--scenario=<name>`: Run only specific scenario
- `--no-cleanup`: Disable automatic cleanup

## Logs

All test executions create detailed logs in `results/integration-YYYYMMDD-HHMMSS.log` containing:
- All executed commands and their output
- API responses and validation results
- Timing information and error details
- PR URLs and cleanup status

## Integration with CI/CD

This framework can be integrated into CI/CD pipelines with conditional execution:

```bash
# Run integration tests every 10th execution
EXECUTION_COUNT=$(cat .execution-count 2>/dev/null || echo 0)
if (( (EXECUTION_COUNT + 1) % 10 == 0 )); then
    go run . test-integration --cleanup
fi
echo $((EXECUTION_COUNT + 1)) > .execution-count
```