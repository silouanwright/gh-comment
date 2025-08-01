# Tasks & Development Progress

This file tracks ongoing development tasks, features, and improvements for `gh-comment`. Tasks are organized by priority and status.

## 🚧 In Progress

### High Priority
- [ ] **Binary Distribution Setup** - Add automated binary releases for better user experience
  - [x] Create GitHub Actions workflow for binary releases
  - [x] Set up `gh-extension-precompile` action
  - [x] Test binary distribution with tag releases
  - [ ] Update installation documentation
  - [ ] Announce binary distribution to users

- [ ] **AI-Optimized README Overhaul** - Complete command documentation for AI processing
  - [ ] Audit all existing commands vs README coverage gaps
  - [ ] Research and propose README optimization approaches (with A-F ratings)
  - [ ] Document every command with full flag specifications
  - [ ] Add comprehensive examples for each command combination
  - [ ] Structure content for AI parsing and synthesis
  - [ ] Test README with AI systems for completeness
  - [ ] Implement chosen optimization approach

### Medium Priority
- [ ] **Increase Test Coverage to 80%** - Refactor commands with dependency injection
  - [ ] Refactor `cmd/list.go` to use new GitHub API client
  - [ ] Refactor `cmd/reply.go` to use dependency injection
  - [ ] Refactor `cmd/add.go` to use dependency injection
  - [ ] Test actual command execution with mocked GitHub API calls
  - [ ] Test repository and PR detection logic
  - [ ] Test file operations and output formatting
  - [ ] Test remaining utility functions
  - [ ] Add comprehensive unit tests for refactored commands
  - [ ] Update integration tests to use new structure

- [ ] **Cross-Platform Testing** - Ensure compatibility across all platforms
  - [ ] Add Windows-specific test scenarios (path separators, line endings)
  - [ ] Test shell compatibility (bash, zsh, fish, PowerShell)
  - [ ] Verify testscript behavior on different operating systems
  - [ ] Add platform-specific golden files if needed
  - [ ] Test GitHub CLI integration across platforms

- [ ] **Automated Test Data Cleanup** - Implement cleanup routines for E2E tests
  - [ ] Add test repository cleanup after E2E test runs
  - [ ] Implement comment cleanup for failed test scenarios
  - [ ] Add test data isolation to prevent cross-test contamination
  - [ ] Create test data lifecycle management

- [x] **License Evaluation** - Research and potentially switch from MIT to ISC license
  - [x] Research MIT vs ISC license differences and benefits
  - [x] Analyze legal implications and compatibility
  - [x] Evaluate impact on contributors and users
  - [x] Check ecosystem adoption (Go CLI tools, GitHub extensions)
  - [x] Make recommendation based on research findings
  - [x] **Decision: Keep MIT License** - Better ecosystem alignment, legal clarity, and industry adoption

## 📋 Planned Features

### Core Features
- [ ] **GitLab-style line offset syntax** - Support `[SUGGEST:+2: code]` and `[SUGGEST:-1: code]`
  - [ ] Design syntax specification
  - [ ] Implement relative line positioning
  - [ ] Add tests for offset syntax
  - [ ] Update documentation

- [ ] **Advanced filtering** - Filter comments by status, author, date, resolved state
  - [ ] Add `--status` flag (open, resolved, all)
  - [ ] Add `--since` and `--until` date filtering
  - [ ] Add `--resolved` boolean filter
  - [ ] Extend `--author` filtering capabilities

- [ ] **Configuration file support** - Default flags and repository settings
  - [ ] Design configuration file format (YAML/JSON)
  - [ ] Implement config file parsing
  - [ ] Add `--config` flag support
  - [ ] Create default config generation command

- [ ] **Template system** - Reusable comment patterns and workflows
  - [ ] Design template file format
  - [ ] Implement template loading and substitution
  - [ ] Add built-in templates for common scenarios
  - [ ] Create template sharing mechanism

### Quality & Performance
- [ ] **Cross-Platform Testing** - Ensure compatibility across all platforms
  - [ ] Add Windows-specific test scenarios
  - [ ] Test path separators and line endings
  - [ ] Verify shell compatibility (bash, zsh, fish, PowerShell)
  - [ ] Add platform-specific golden files if needed

- [ ] **Performance Optimizations**
  - [ ] Optimize comment fetching with pagination
  - [ ] Add caching for frequently accessed data
  - [ ] Implement parallel API calls where possible
  - [ ] Monitor and optimize memory usage

### User Experience
- [ ] **Batch operations** - Apply operations to multiple comments at once
  - [ ] Design batch operation syntax
  - [ ] Implement batch comment creation
  - [ ] Add batch reaction management
  - [ ] Create batch editing capabilities

- [ ] **Export functionality** - Export comments to various formats
  - [ ] JSON export format
  - [ ] CSV export for spreadsheet analysis
  - [ ] Markdown export for documentation
  - [ ] HTML export for presentations

## ✅ Recently Completed

### August 2025
- [x] **Fixed Failing Tests** - Resolved test failures in `TestHelperFunctions` and `TestPRContext`
- [x] **GitHub API Client Refactoring** - Created clean abstraction layer in `internal/github/`
  - [x] Created `GitHubAPI` interface
  - [x] Implemented `RealClient` with actual GitHub API calls
  - [x] Enhanced `MockClient` for comprehensive testing
  - [x] Added support for all extension operations
- [x] **Performance Regression Testing** - Full CI/CD integration with benchstat comparison
  - [x] Enhanced GitHub Actions workflow with benchmark comparison
  - [x] Added scheduled benchmark tracking for historical data
  - [x] Created local benchmark script (`scripts/benchmark.sh`)
  - [x] Configured performance regression alerts
- [x] **Complete Testing Infrastructure** - Comprehensive testing suite implementation
  - [x] Unit tests with mocks and table-driven patterns
  - [x] Integration tests with testscript (3 .txtar files, 100% passing)
  - [x] Fuzz tests for edge case discovery (5 comprehensive fuzz functions)
  - [x] E2E tests with real GitHub API testing and safety measures
  - [x] Benchmark tests with performance monitoring
  - [x] Dependency injection tests with full workflow testing
  - [x] Pre-commit hooks for automated quality checks
- [x] **Documentation Updates** - Updated roadmaps and testing documentation
  - [x] Updated `TESTING_ROADMAP.md` with current status (now merged)
  - [x] Updated `README.md` roadmap section
  - [x] Enhanced testing documentation with new capabilities
  - [x] Created comprehensive `TESTING.md` and `E2E_TESTING.md` guides

### Earlier 2025
- [x] **CI/CD Pipeline** - Complete GitHub Actions workflow with multi-platform testing
- [x] **Coverage Reporting** - Automated coverage tracking and thresholds
- [x] **Golden File Testing** - CLI output verification system

## 🎯 Success Metrics

### Code Quality
- **Test Coverage**: Currently 16.1% (target: 80%+)
- **Test Success Rate**: 100% passing
- **Performance**: All benchmarks stable with regression detection

### User Experience
- **Installation**: Source-based (migrating to binary distribution)
- **Update Mechanism**: `gh extension upgrade` supported
- **Cross-Platform**: macOS, Linux, Windows supported
- **Documentation**: Comprehensive guides available

---

## 📝 Task Management Notes

### How to Use This File
1. **Add new tasks** under appropriate sections
2. **Move tasks** between sections as they progress
3. **Check off subtasks** using `- [x]` syntax
4. **Archive completed tasks** to the "Recently Completed" section
5. **Update dates** when moving tasks to completed

### Task Priorities
- **High**: Critical functionality, bug fixes, security issues
- **Medium**: Important features, performance improvements
- **Low**: Nice-to-have features, documentation improvements

### Status Indicators
- `🚧` In Progress
- `📋` Planned
- `✅` Completed
- `🎯` Success Metrics
- `⚠️` Blocked/Issues
- `🔄` Under Review

Last updated: August 2025 (merged from TESTING_ROADMAP.md)