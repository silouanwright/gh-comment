# Pre-Commit Hooks Setup Guide

This guide explains how to set up and use pre-commit hooks for the `gh-comment` project.

## 🎯 **Why Pre-Commit Hooks?**

Pre-commit hooks automatically run quality checks before each commit, ensuring:
- ✅ Code compiles and tests pass
- ✅ Code is properly formatted and linted
- ✅ Security vulnerabilities are caught early
- ✅ Consistent code quality across all contributors
- ✅ Faster CI/CD pipeline (issues caught locally)

## 📦 **Installation**

### 1. Install pre-commit

**macOS (using Homebrew):**
```bash
brew install pre-commit
```

**Linux/WSL:**
```bash
pip install pre-commit
```

**Alternative (using Go):**
```bash
go install github.com/pre-commit/pre-commit@latest
```

### 2. Install Required Go Tools

Some hooks require additional Go tools:

```bash
# Security scanner
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Import formatter
go install golang.org/x/tools/cmd/goimports@latest

# Advanced linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 3. Install Hooks in Repository

```bash
cd /path/to/gh-comment
pre-commit install
pre-commit install --hook-type commit-msg  # For commit message linting
```

## 🚀 **Usage**

### Automatic Execution
Once installed, hooks run automatically on every `git commit`. Example:
```bash
git add .
git commit -m "feat: add new feature"
# Hooks run automatically here
```

### Manual Execution
Run hooks on all files without committing:
```bash
pre-commit run --all-files
```

Run specific hook:
```bash
pre-commit run go-test-repo-mod
pre-commit run golangci-lint-repo-mod
```

### Skip Hooks (Emergency Use)
```bash
git commit --no-verify -m "emergency fix"
```

## 🔧 **Hook Configuration**

Our `.pre-commit-config.yaml` includes:

### **Basic Quality Checks:**
- Remove trailing whitespace
- Fix end-of-file issues
- Validate YAML files
- Prevent large files from being committed

### **Go-Specific Checks:**
- **`go-build-repo-mod`** - Ensures code compiles
- **`go-mod-tidy-repo`** - Keeps go.mod clean
- **`go-test-repo-mod`** - Runs unit tests (fast tests only)
- **`go-vet-repo-mod`** - Static analysis
- **`go-fmt-repo`** - Code formatting
- **`go-imports-repo`** - Import management
- **`go-sec-repo-mod`** - Security scanning
- **`golangci-lint-repo-mod`** - Comprehensive linting

### **Commit Message Linting:**
- Enforces conventional commit format
- Examples: `feat:`, `fix:`, `docs:`, `test:`

## ⚡ **Performance Optimization**

### Fast vs Comprehensive Checks

**Pre-commit (Fast - runs on every commit):**
- Unit tests with `-short` flag (excludes E2E tests)
- Basic linting and formatting
- Security scanning
- Build verification

**CI/CD (Comprehensive - runs on push):**
- Full test suite including E2E tests
- Cross-platform testing
- Coverage reporting
- Integration tests

### Test Exclusions

The pre-commit configuration excludes:
- E2E tests (too slow for commit-time)
- Integration tests that require external services
- Benchmark tests (run in CI)

## 🛠️ **Customization**

### Modify Hook Arguments

Edit `.pre-commit-config.yaml` to customize behavior:

```yaml
- id: go-test-repo-mod
  args: ['-short', '-race', '-timeout=60s', '-v']  # Add verbose output
```

### Add Custom Hooks

```yaml
- id: my-cmd-repo
  name: 'Custom Go Tool'
  args: [go, run, ./scripts/custom-check.go]
```

### Skip Specific Files

```yaml
exclude: |
  (?x)^(
    testdata/.*|
    vendor/.*|
    .*_generated\.go$
  )$
```

## 🔍 **Troubleshooting**

### Common Issues

**1. Hook fails with "command not found"**
```bash
# Install missing tool
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**2. Tests are too slow**
```bash
# Use -short flag to skip slow tests
args: ['-short', '-timeout=30s']
```

**3. Hook conflicts with IDE formatting**
```bash
# Run pre-commit to see what changed
pre-commit run --all-files
git diff
```

### Debug Mode

```bash
pre-commit run --verbose --all-files
```

## 🎯 **Integration with Existing Testing**

Our pre-commit setup complements the existing testing strategy:

### **Pre-Commit (Local):**
- Fast feedback (< 30 seconds)
- Catches obvious issues
- Ensures basic quality

### **CI/CD Pipeline (Remote):**
- Comprehensive testing
- Cross-platform verification
- Performance benchmarks
- Coverage reporting

### **Manual Testing:**
- E2E tests with real GitHub repos
- Integration testing
- Performance analysis

## 📋 **Best Practices**

1. **Keep hooks fast** - Pre-commit should complete in < 60 seconds
2. **Use `-short` flag** - Skip slow tests in pre-commit
3. **Run comprehensive tests in CI** - Full test suite on push
4. **Commit message standards** - Use conventional commits
5. **Regular updates** - Keep hook versions current

## 🔄 **Maintenance**

### Update Hooks
```bash
pre-commit autoupdate
```

### Clean Cache
```bash
pre-commit clean
```

### Reinstall Hooks
```bash
pre-commit uninstall
pre-commit install
```

## 🎉 **Benefits for gh-comment**

With pre-commit hooks, every commit to `gh-comment` will:
- ✅ Compile successfully
- ✅ Pass unit tests
- ✅ Follow Go formatting standards
- ✅ Pass security scans
- ✅ Have clean imports and dependencies
- ✅ Use conventional commit messages

This ensures high code quality and reduces CI/CD failures! 🚀
