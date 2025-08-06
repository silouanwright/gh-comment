# Go Formatting Consistency Research & Implementation

## Problem Statement

We had inconsistent formatting between local pre-commit hooks and CI, causing endless format failures where local tools pass but CI fails with "File is not properly formatted (goimports)".

## Root Cause Analysis

### Current State Issues
1. **Local pre-commit**: Using `tekwizely/pre-commit-golang` with separate `go-fmt-repo` hook  
2. **CI**: Using `golangci-lint` with `goimports` linter enabled
3. **Result**: Different formatting rules and import grouping standards

### Tool Differences
- **gofmt**: Basic Go source code formatting
- **goimports**: Extends gofmt + manages imports (grouping, sorting, unused removal)  
- **golangci-lint**: Aggregated linter that includes goimports + many other checks

## Research Findings (December 2025)

### Best Practices from Industry Sources

#### 1. **Unified Tool Approach (Recommended)**
- Use `golangci-lint` for BOTH pre-commit and CI
- Eliminates tool inconsistencies
- Single source of truth for formatting rules
- Better performance through caching and parallel execution

#### 2. **Official golangci-lint Pre-commit Hooks**
Available at: `https://github.com/golangci/golangci-lint`

```yaml
- golangci-lint: Only checks changed files (`--new-from-rev HEAD --fix`)
- golangci-lint-full: Checks all files (recommended for CI)  
- golangci-lint-fmt: Formats all files (`golangci-lint fmt`)
- golangci-lint-config-verify: Validates configuration
```

#### 3. **Version Consistency**
- Pin golangci-lint version across environments
- Use same .golangci.yml configuration everywhere
- Ensure reproducible results

### Industry Recommendations

#### **Pre-commit Strategy**
```yaml
# Option A: golangci-lint for changed files only (fast feedback)
- repo: https://github.com/golangci/golangci-lint
  rev: v1.64.8  # Pin version
  hooks:
    - id: golangci-lint
      args: [--timeout=5m]

# Option B: Format-only hook for speed
- repo: https://github.com/golangci/golangci-lint  
  rev: v1.64.8
  hooks:
    - id: golangci-lint-fmt
```

#### **CI Strategy**
```yaml
# GitHub Actions - use official action for performance
- name: golangci-lint
  uses: golangci/golangci-lint-action@v6
  with:
    version: v1.64.8  # Same version as pre-commit
    args: --timeout=5m
```

## Implementation Decision

### **Hybrid Approach (Chosen)**
1. **Pre-commit**: `golangci-lint` for changed files (fast feedback)
2. **CI**: `golangci-lint` with same config (comprehensive check)
3. **Configuration**: Single `.golangci.yml` for both environments

### **Why This Approach**
- ✅ **Consistency**: Same tool, same rules everywhere
- ✅ **Performance**: Pre-commit only checks changed files
- ✅ **Reliability**: CI checks all files as safety net  
- ✅ **Developer Experience**: Fast local feedback, thorough CI validation
- ✅ **Maintainability**: Single configuration file

## Alternative Approaches Considered

### **Separate Tools Approach (Rejected)**
```yaml
# Not recommended - leads to inconsistencies
- gofmt → CI: gofmt  
- goimports → CI: goimports
- go vet → CI: go vet  
- golangci-lint → CI: golangci-lint
```
**Problems**: Different tools may have conflicting formatting rules

### **CI-Only Approach (Rejected)**  
- No pre-commit hooks, rely only on CI
- **Problems**: Late feedback, poor developer experience

## Configuration Details

### **Optimized .golangci.yml**
- Enable core linters: `gofmt`, `goimports`, `govet`, `staticcheck`
- Disable overly strict rules for development productivity
- Exclude test files from some checks
- Focus on critical issues vs. style preferences

### **Pre-commit Hook Benefits**
1. **Immediate Feedback**: Catch issues before push
2. **Selective Checking**: Only modified files (performance) 
3. **Auto-fixing**: `--fix` flag for automatic corrections
4. **Staging Integration**: Works with git staged files

### **CI Benefits**  
1. **Comprehensive**: Full codebase validation
2. **Safety Net**: Catches anything pre-commit missed
3. **Consistent Environment**: Same rules for all PRs
4. **Performance**: Official GitHub Action optimizations

## Implementation Steps

1. ✅ **Research completed**
2. 🔄 **Update pre-commit config** to use official golangci-lint hooks
3. 🔄 **Verify CI uses same version** as pre-commit  
4. 🔄 **Test consistency** between local and CI formatting
5. 🔄 **Document process** for team adoption

## Expected Outcomes

- ✅ **Zero formatting inconsistencies** between local and CI
- ✅ **Faster development cycle** with immediate local feedback  
- ✅ **Better code quality** through comprehensive linting
- ✅ **Reduced CI failures** due to formatting issues
- ✅ **Single source of truth** for code standards

## Sources & References

- [golangci-lint Pre-commit Hooks](https://github.com/golangci/golangci-lint/blob/main/.pre-commit-hooks.yaml)
- [Go Linting Best Practices for CI/CD](https://medium.com/@tedious/go-linting-best-practices-for-ci-cd-with-github-actions-aa6d96e0c509)
- [Complete Guide to Linting Go Programs](https://freshman.tech/linting-golang/)
- [golangci-lint Configuration Guide](https://golangci-lint.run/usage/configuration/)

---
*Research completed: December 2025*  
*Status: Ready for implementation*