# Integration Testing Scripts

Automated integration testing suite for `gh-comment` that validates all help text examples work correctly with real GitHub APIs.

## Overview

This modular script suite replaces manual AI-driven integration testing with fast, reliable automation. It extracts examples from help text, executes them against real GitHub PRs, and provides comprehensive reporting.

## Quick Start

```bash
# Full integration test cycle (setup → test → cleanup)
./scripts/integration/run_integration.sh

# Setup test environment only (for development)
./scripts/integration/run_integration.sh --setup-only

# Resume testing from existing environment
./scripts/integration/run_integration.sh --resume
```

## Script Architecture

### Core Scripts

| Script | Purpose | Dependencies |
|--------|---------|--------------|
| `00_preflight.sh` | Environment validation | None |
| `01_setup.sh` | Create PR, branch, test files | `00_preflight.sh` |
| `02_parse_help.sh` | Extract examples from help text | `01_setup.sh` |
| `03_run_tests.sh` | Execute all extracted examples | `02_parse_help.sh` |
| `04_cleanup.sh` | Close PR, cleanup resources | Any previous |

### Orchestrator

| Script | Purpose |
|--------|---------|
| `run_integration.sh` | Master orchestrator with phase control |

## Workflow Phases

### 1. 🔍 Preflight Checks (`00_preflight.sh`)
- Validates working tree is clean
- Checks required tools (go, git, gh)
- Verifies GitHub CLI authentication
- Tests binary compilation
- Ensures test coverage >80%

### 2. 🚀 Setup (`01_setup.sh`)
- Creates unique test branch (`integration-test-YYYYMMDD-HHMMSS`)
- Generates realistic test files:
  - `src/api.js` - Express.js middleware
  - `src/main.go` - Go CLI application
  - `tests/auth_test.js` - Jest test suite
- Creates GitHub PR for testing
- Saves environment variables for other scripts

### 3. 📖 Parse Help Text (`02_parse_help.sh`)
- Discovers all available commands
- Extracts examples from each command's help text
- Generates executable test cases with variable substitution
- Creates master test runner script

### 4. 🧪 Run Tests (`03_run_tests.sh`)
- Executes all extracted help text examples
- Captures detailed results and logs
- Generates comprehensive summary reports
- Provides success/failure statistics

### 5. 🧹 Cleanup (`04_cleanup.sh`)
- Archives all test results and logs
- Closes GitHub PR with summary comment
- Deletes test branch (local and remote)
- Cleans up temporary files

## Usage Examples

### Full Integration Test
```bash
# Complete cycle with cleanup
./scripts/integration/run_integration.sh
```

### Development Workflow
```bash
# 1. Set up test environment
./scripts/integration/run_integration.sh --setup-only

# 2. Iterate on help text fixes, then test
./scripts/integration/run_integration.sh --resume --test-only

# 3. Clean up when done
./scripts/integration/run_integration.sh --cleanup-only
```

### Debugging Failed Tests
```bash
# Run tests without cleanup to inspect PR
./scripts/integration/run_integration.sh --no-cleanup

# Manual cleanup later
./scripts/integration/run_integration.sh --cleanup-only --force-cleanup
```

## Configuration

### Environment Variables
The setup script creates environment files in `logs/integration/test-env-TIMESTAMP.sh`:

```bash
export TEST_BRANCH="integration-test-20250805-143022"
export PR_NUM="42"
export PR_URL="https://github.com/user/repo/pull/42"
export TIMESTAMP="20250805-143022"
# ... and more
```

### Generated Test Files
The setup creates realistic test files that match help text examples:

- **`src/api.js`** - Express.js middleware with auth, rate limiting
- **`src/main.go`** - Go CLI with command processing
- **`tests/auth_test.js`** - Jest tests with various scenarios
- **`README.md`** - Documentation for the test repository

## Output and Logging

### Directory Structure
```
logs/integration/
├── test-env-TIMESTAMP.sh           # Environment variables
├── test-results-TIMESTAMP.txt      # Detailed test results
├── summary-TIMESTAMP.md            # Markdown summary report
├── test-cases-TIMESTAMP/           # Generated test cases
│   ├── add_examples.sh
│   ├── review_examples.sh
│   └── ...
├── cmd-*-TIMESTAMP.log             # Individual command logs
└── archived-TIMESTAMP/             # Complete archive after cleanup
    ├── README.md                   # Archive summary
    └── [all test artifacts]
```

### Test Results Format
```
# test-results-TIMESTAMP.txt
command:example_id:status:details
add:example_1:SUCCESS:Comment added successfully
review:example_2:FAILED:File not found in PR diff
list:example_1:SUCCESS:Listed 5 comments
```

### Summary Report
The summary report (`summary-TIMESTAMP.md`) includes:
- Test metrics and success rate
- Failed tests with details
- Commands tested with individual results
- Links to detailed logs

## Advanced Usage

### Phase Control
```bash
# Run specific phases
./scripts/integration/run_integration.sh --parse-only
./scripts/integration/run_integration.sh --test-only
./scripts/integration/run_integration.sh --cleanup-only

# Skip phases
./scripts/integration/run_integration.sh --no-cleanup
./scripts/integration/run_integration.sh --resume
```

### Force Operations
```bash
# Force cleanup without confirmation
./scripts/integration/run_integration.sh --cleanup-only --force-cleanup
```

## Error Handling

### Common Issues and Solutions

| Issue | Solution |
|-------|----------|
| Working tree dirty | Commit or stash changes |
| GitHub CLI not authenticated | Run `gh auth login` |
| Binary won't build | Fix compilation errors |
| Test environment not found | Run setup phase first |
| PR creation fails | Check repository permissions |

### Recovery Procedures

**Lost test environment:**
```bash
# Check for existing environments
ls logs/integration/test-env-*.sh

# Resume from specific environment
source logs/integration/test-env-TIMESTAMP.sh
./scripts/integration/run_integration.sh --resume --test-only
```

**Stuck test branch:**
```bash
# Force cleanup
./scripts/integration/run_integration.sh --cleanup-only --force-cleanup
```

## Integration with Development

### Pre-commit Testing
Add to your development workflow:
```bash
# Before major changes
./scripts/integration/run_integration.sh --setup-only

# After fixing help text
./scripts/integration/run_integration.sh --resume --test-only
```

### CI Considerations
⚠️ **Note:** These scripts are designed for local development only. GitHub doesn't allow creating PRs in CI for security reasons. Use for manual validation of help text accuracy.

## Benefits Over Manual Testing

- **Speed:** 10x faster than manual AI-driven testing
- **Reliability:** Consistent execution, no human error
- **Completeness:** Tests every help text example automatically
- **Traceability:** Detailed logs and reports for debugging
- **Repeatability:** Identical test conditions every run
- **Scalability:** Easy to add new commands and examples

## Contributing

When adding new commands or changing help text:

1. Update help text in command files
2. Run integration tests to validate examples
3. Fix any broken examples found
4. Archive test results for documentation

The integration tests ensure help text stays accurate and trustworthy for users.