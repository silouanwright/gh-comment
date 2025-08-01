# End-to-End Testing Guide

This document describes how to run and maintain end-to-end (E2E) tests for the `gh-comment` CLI extension.

## Overview

E2E tests verify the complete functionality of `gh-comment` by running against real GitHub repositories. These tests ensure that:

- Commands work correctly with the GitHub API
- Authentication flows function properly
- Real-world scenarios are handled correctly
- Integration with the `gh` CLI works as expected

## Prerequisites

### Required Environment Variables

```bash
# GitHub personal access token with appropriate permissions
export GH_TOKEN="ghp_your_token_here"

# Test repository (should be a repository you have access to)
export GH_E2E_REPO="owner/repo"

# PR number for testing (should exist in the test repository)
export GH_E2E_PR="123"

# Enable E2E tests (they are skipped by default)
export RUN_E2E_TESTS="true"
```

### GitHub Token Permissions

Your `GH_TOKEN` needs the following permissions:
- `repo` - Full repository access
- `read:org` - Read organization membership (if testing on org repos)

### Test Repository Setup

The test repository should have:
- At least one open PR with comments
- Mix of issue comments and review comments (if possible)
- Permissions for your token to create/modify comments

## Running E2E Tests

### Basic E2E Test Run

```bash
# Set up environment
export GH_TOKEN="your_token"
export GH_E2E_REPO="your-org/test-repo"
export GH_E2E_PR="1"
export RUN_E2E_TESTS="true"

# Run E2E tests
go test -run=TestE2E -v ./cmd/
```

### Running Specific E2E Test Cases

```bash
# Test only the list functionality
go test -run=TestE2E/list_comments_e2e -v ./cmd/

# Test only the comment workflow
go test -run=TestE2E/comment_workflow_e2e -v ./cmd/
```

### E2E Test Setup Verification

```bash
# Test that E2E environment is configured correctly
go test -run=TestE2ESetup -v ./cmd/
```

### E2E Benchmarks

```bash
# Run E2E performance benchmarks
export RUN_E2E_TESTS="true"
go test -bench=BenchmarkE2E -v ./cmd/
```

## Test Structure

### Test Categories

1. **Environment Setup Tests** (`TestE2ESetup`)
   - Verify environment variables are set correctly
   - Test command capture mechanism
   - Validate test infrastructure

2. **List Command Tests** (`testListCommentsE2E`)
   - Test basic comment listing
   - Test command flags (--quiet, --author, etc.)
   - Test error handling for non-existent PRs

3. **Comment Workflow Tests** (`testCommentWorkflowE2E`)
   - Test dry-run functionality
   - Test comment creation workflow
   - Test reaction and resolution workflows

### Test Safety

E2E tests are designed to be **non-destructive**:
- Use `--dry-run` mode for operations that would modify data
- Generate unique test identifiers to avoid conflicts
- Clean up any test data created during runs
- Fail gracefully if permissions are insufficient

## CI/CD Integration

### GitHub Actions

Add E2E tests to your CI pipeline:

```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Run E2E Tests
        env:
          GH_TOKEN: ${{ secrets.GH_TOKEN }}
          GH_E2E_REPO: ${{ secrets.GH_E2E_REPO }}
          GH_E2E_PR: ${{ secrets.GH_E2E_PR }}
          RUN_E2E_TESTS: "true"
        run: go test -run=TestE2E -v ./cmd/
```

### Required Secrets

Set these secrets in your GitHub repository:
- `GH_TOKEN` - GitHub personal access token
- `GH_E2E_REPO` - Test repository name
- `GH_E2E_PR` - Test PR number

## Troubleshooting

### Common Issues

#### Tests are Skipped

```
E2E tests skipped. Set RUN_E2E_TESTS=true to run.
```

**Solution**: Set the `RUN_E2E_TESTS` environment variable:
```bash
export RUN_E2E_TESTS="true"
```

#### Missing GitHub Token

```
GH_TOKEN not set, skipping E2E tests
```

**Solution**: Set your GitHub personal access token:
```bash
export GH_TOKEN="ghp_your_token_here"
```

#### Repository Access Issues

```
failed to fetch comments: 404 Not Found
```

**Solutions**:
- Verify the repository exists and is accessible
- Check that your token has the required permissions
- Ensure the PR number exists in the specified repository

#### Rate Limiting

```
failed to fetch comments: 403 rate limit exceeded
```

**Solutions**:
- Wait for the rate limit to reset
- Use a token with higher rate limits
- Reduce the frequency of E2E test runs

### Debug Mode

Enable verbose output for debugging:

```bash
export GH_DEBUG=1
go test -run=TestE2E -v ./cmd/
```

## Best Practices

### Test Repository Management

1. **Use a Dedicated Test Repository**
   - Create a repository specifically for testing
   - Keep it separate from production repositories
   - Regularly clean up test data

2. **Test Data Lifecycle**
   - Generate unique identifiers for test comments
   - Clean up test data after each run
   - Use descriptive test messages (e.g., "🤖 E2E Test Comment")

3. **Permission Management**
   - Use tokens with minimal required permissions
   - Rotate tokens regularly
   - Store tokens securely in CI/CD systems

### Test Reliability

1. **Handle Flaky Tests**
   - Implement retry logic for network operations
   - Add appropriate timeouts
   - Handle rate limiting gracefully

2. **Environment Isolation**
   - Don't rely on specific comment IDs
   - Test against multiple PR states
   - Handle missing or changed test data

3. **Monitoring**
   - Track E2E test success rates
   - Monitor API usage and rate limits
   - Alert on consistent failures

## Extending E2E Tests

### Adding New Test Cases

1. **Create Test Function**
   ```go
   func testNewFeatureE2E(t *testing.T, repo string, prNumber int) {
       // Test implementation
   }
   ```

2. **Add to Main Test**
   ```go
   t.Run("new_feature_e2e", func(t *testing.T) {
       testNewFeatureE2E(t, repo, prNumber)
   })
   ```

3. **Update Documentation**
   - Document new environment variables
   - Update troubleshooting guide
   - Add examples

### Testing New Commands

When adding new commands to `gh-comment`:

1. **Add Command Handler**
   ```go
   case "new-command":
       return runNewCommandE2E(cmdArgs)
   ```

2. **Implement E2E Function**
   ```go
   func runNewCommandE2E(args []string) (string, error) {
       // Command simulation for E2E testing
   }
   ```

3. **Add Test Cases**
   ```go
   func testNewCommandE2E(t *testing.T, repo string, prNumber int) {
       // Test the new command
   }
   ```

## Security Considerations

### Token Security

- **Never commit tokens to version control**
- **Use environment variables or secure secret management**
- **Rotate tokens regularly**
- **Use tokens with minimal required permissions**

### Test Data Privacy

- **Don't include sensitive information in test comments**
- **Use public repositories for testing when possible**
- **Clean up test data promptly**
- **Be mindful of repository visibility**

### Rate Limiting

- **Respect GitHub's rate limits**
- **Implement backoff strategies**
- **Monitor API usage**
- **Use authenticated requests for higher limits**

---

For more information about testing strategies, see [TESTING.md](TESTING.md).
