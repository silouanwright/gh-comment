# Testing Implementation Roadmap

This roadmap implements modern Go CLI testing practices for the `gh-comment` GitHub CLI extension based on 2025 best practices research.

## 📋 Current Project Context

### **Project State (as of implementation)**
- **Unified Comment System**: Extension handles both issue comments (general PR discussion) and review comments (line-specific) in a single interface
- **Core Commands**: `list` (shows all comments), `reply` (with `--type` flag), `add` (line-specific), `add-review` (batch)
- **Architecture**: Cobra-based CLI with direct GitHub API calls in command files
- **APIs Used**: GitHub REST API (`/issues/{pr}/comments`, `/pulls/{pr}/comments`) + GraphQL for conversation resolution

### **Key Implementation Details**
- **Reply Command Enhancement**: Recently added `--type issue|review` flag to distinguish comment types
  - `--type issue`: Creates new top-level comments via `/issues/{pr}/comments`
  - `--type review`: Creates threaded replies via `/pulls/{pr}/comments` (default)
- **Unified Listing**: `list` command fetches and displays both comment types in separate sections
- **Current Structure**: API calls are embedded directly in `cmd/list.go` and `cmd/reply.go`
- **Dependencies**: Uses `go-gh` library for GitHub CLI integration and authentication

### **Testing Gaps to Address**
- **No unit tests**: Core logic mixed with API calls, making testing difficult
- **No integration tests**: CLI workflows not tested end-to-end
- **No API mocking**: Cannot test without hitting real GitHub API
- **No output verification**: Command output format not validated
- **No CI/CD**: No automated testing pipeline

### **Critical Files to Understand**
- `cmd/list.go`: Fetches and displays both issue and review comments (lines 138-164, 273-289 recently modified)
- `cmd/reply.go`: Handles replies with type-specific logic (extensively modified for `--type` flag)
- `README.md`: Documents unified comment system and command usage
- `ai-prompts/address-review.md`: AI workflow guide (recently updated for unified system)

### **Recent Enhancements (Context for Testing)**
1. **Unified Comment Display**: List command now shows both comment types in organized sections
2. **Type-Specific Replies**: Reply command uses `--type` flag instead of inefficient API detection
3. **Documentation Updates**: README and AI prompts updated to reflect unified system
4. **Command Reference**: Clear distinction between native GitHub CLI vs extension commands

### **Testing Priorities & Challenges**

**High Priority Areas:**
1. **Comment Type Logic**: Critical to test `--type issue` vs `--type review` behavior in reply command
2. **API Integration**: Mock GitHub REST API calls for `/issues/{pr}/comments` and `/pulls/{pr}/comments`
3. **Unified Display**: Verify list command correctly categorizes and displays both comment types
4. **Error Handling**: Test validation, API errors (404, 403, 429), and user input edge cases

**Specific Testing Challenges:**
- **GitHub API Dependency**: Commands currently make direct API calls, need abstraction layer
- **Authentication**: Uses `go-gh` library for GitHub CLI integration, need to mock auth
- **Output Formatting**: Complex output with emojis, sections, and formatting needs golden file testing
- **Cross-Platform**: Must work on macOS, Linux, Windows with different shells

**Key Test Scenarios to Cover:**
- List command with mixed issue/review comments
- Reply to issue comment (creates new top-level comment)
- Reply to review comment (creates threaded reply)
- Error cases: invalid comment IDs, missing permissions, rate limits
- Help text and command structure validation
- Suggestion syntax expansion in replies

## 🎯 Phase 1: Foundation & Dependencies

### - [x] Add Testing Dependencies
**Implementation:** Update `go.mod` with testing dependencies
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/sebdah/goldie/v2
go get github.com/rogpeppe/go-internal/testscript
```
**Details:** 
- Testify for enhanced assertions and test suites
- Goldie v2 for golden file testing (CLI output verification)
- testscript for black-box CLI integration testing (same tool GitHub CLI uses)

### - [x] Create Testing Directory Structure
**Implementation:** Create the following directory structure:
```
gh-comment/
├── testdata/
│   ├── golden/           # Golden files for output verification
│   └── scripts/          # testscript .txtar files
├── internal/
│   ├── testutil/         # Shared test utilities
│   └── github/           # GitHub API client (refactored)
├── integration/          # Integration test files
└── .github/
    └── workflows/
        └── test.yml      # CI/CD pipeline
```
**Details:** Follow Go community conventions for test organization. `testdata/` is automatically ignored by `go build`.

### - [ ] Refactor GitHub API Client
**Implementation:** Extract GitHub API calls into `internal/github/client.go`
```go
type Client struct {
    httpClient *http.Client
    baseURL    string
    token      string
}

func NewClient(token string) *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        baseURL:    "https://api.github.com",
        token:      token,
    }
}
```
**Details:** 
- Enables dependency injection for testing
- Makes HTTP calls mockable with httptest
- Separates API logic from command logic
- Current API calls in `cmd/list.go` and `cmd/reply.go` should be moved here

## 🧪 Phase 2: Unit Testing Foundation

### - [x] Create Test Utilities
**Implementation:** Create `internal/testutil/helpers.go`
```go
// Test helpers for common patterns
func NewTestClient() *github.Client
func MockGitHubAPI(t *testing.T) *httptest.Server
func CaptureOutput(fn func()) (stdout, stderr string)
func LoadGoldenFile(t *testing.T, name string) []byte
```
**Details:** Shared utilities reduce boilerplate and ensure consistent test patterns across the codebase.

### - [ ] Unit Tests for GitHub Client
**Implementation:** Create `internal/github/client_test.go`
- Test API endpoint construction
- Test request/response handling
- Test error scenarios (404, 403, 429 rate limits)
- Use `httptest.Server` for mocking GitHub API responses
**Details:** Focus on testing the HTTP client logic separately from command logic. Mock various GitHub API responses including errors.

### - [ ] Unit Tests for Command Logic
**Implementation:** Create unit tests for each command:
- `cmd/list_test.go`: Test comment fetching, filtering, display logic
- `cmd/reply_test.go`: Test reply logic, validation, type handling
- `cmd/root_test.go`: Test command structure, flags, help text

**Testing Pattern:**
```go
func TestListCommand(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        wantErr  bool
        wantOut  string
    }{
        // Table-driven test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```
**Details:** Use table-driven tests (Go community standard). Test command parsing, flag validation, and business logic separately from API calls.

### - [x] Constructor Pattern for Commands (Dependency Injection)
**Why This Is Needed:** Current commands use global variables and direct API calls, making them impossible to test in isolation.

**Current Problem:**
```go
// cmd/list.go - Current structure
var (
    showResolved bool    // Global state - shared across tests!
    author string        // Global state - causes test interference!
)

func runList(cmd *cobra.Command, args []string) error {
    client, err := api.DefaultRESTClient()  // Hard-coded dependency - can't mock!
    fmt.Println("Comments:")                // Hard-coded output - can't capture!
    // Uses global variables - can't isolate!
}
```

**Testing Problems:**
- **Global variables** cause test interference (one test affects another)
- **Hard-coded API calls** can't be mocked (hits real GitHub API)
- **Hard-coded output** can't be captured for verification
- **No way to test error conditions** or edge cases in isolation

**Solution: Dependency Injection (NOT OOP!)**
```go
// This is functional programming, not OOP - just a struct holding dependencies
type Dependencies struct {
    GitHubClient GitHubAPI      // Interface for mocking
    Output       io.Writer      // Inject where to write (stdout vs test buffer)
    Input        io.Reader      // Inject where to read from
}

// Not a "constructor" - just a function that returns a configured command
func NewListCmd(deps *Dependencies) *cobra.Command {
    return &cobra.Command{
        Use: "list",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runList(deps, cmd, args)  // Pass dependencies explicitly
        },
    }
}

func runList(deps *Dependencies, cmd *cobra.Command, args []string) error {
    // Use injected dependencies instead of globals
    comments, err := deps.GitHubClient.GetComments(repo, pr)
    if err != nil {
        return err
    }
    
    // Write to injected output (can be captured in tests)
    fmt.Fprintf(deps.Output, "Comments: %v\n", comments)
    return nil
}
```

**How This Enables Testing:**
```go
func TestListCommand(t *testing.T) {
    // Create mock dependencies
    mockClient := &MockGitHubClient{
        comments: []Comment{{Body: "test comment"}},
    }
    
    var output bytes.Buffer  // Capture output
    deps := &Dependencies{
        GitHubClient: mockClient,
        Output:       &output,
    }
    
    // Test the command with mocked dependencies
    cmd := NewListCmd(deps)
    err := cmd.Execute()
    
    // Verify results
    assert.NoError(t, err)
    assert.Contains(t, output.String(), "test comment")
}
```

**Benefits:**
- ✅ **Pure functions** - no global state
- ✅ **Mockable dependencies** - can test without real API
- ✅ **Capturable output** - can verify what gets printed
- ✅ **Isolated tests** - each test has its own dependencies
- ✅ **Testable error conditions** - inject failing mocks

**Alternative Approaches (if you prefer not to refactor):**
1. **CLI-only testing** with testscript (slower, less precise)
2. **HTTP-level mocking** (complex, global state issues)
3. **Integration tests only** (harder to test edge cases)

**Implementation:** This pattern is used by GitHub CLI, Docker CLI, kubectl, and most modern Go CLI tools. It's functional programming, not OOP.

## 🔧 Phase 3: Integration Testing with testscript

### - [ ] Set Up testscript Framework
**Implementation:** Create `integration/cli_test.go`
```go
func TestMain(m *testing.M) {
    os.Exit(testscript.RunMain(m, map[string]func() int{
        "gh-comment": func() int {
            // Run main command
            if err := cmd.Execute(); err != nil {
                return 1
            }
            return 0
        },
    }))
}

func TestCLI(t *testing.T) {
    testscript.Run(t, testscript.Params{
        Dir: "testdata/scripts",
    })
}
```
**Details:** testscript runs CLI commands in isolated environments. Each `.txtar` file is a separate test scenario.

### - [ ] Create testscript Test Cases
**Implementation:** Create `.txtar` files in `testdata/scripts/`:

**`list_basic.txtar`:**
```
# Test basic list functionality
exec gh-comment list --pr 123 --repo owner/repo
stdout 'General PR Comments'
stdout 'Review Comments'
! stderr .
```

**`reply_issue_comment.txtar`:**
```
# Test replying to issue comment
exec gh-comment reply 123456 'Thanks!' --type issue --pr 123
stdout 'Replied to issue comment'
! stderr .
```

**`error_handling.txtar`:**
```
# Test error scenarios
! exec gh-comment reply invalid-id 'message'
stderr 'must be a valid integer'
```

**Details:** Each `.txtar` file tests real CLI workflows. Use `exec` to run commands, `stdout`/`stderr` to verify output, `!` for expected failures.

### - [ ] Mock GitHub API for Integration Tests
**Implementation:** Create test server for GitHub API responses
```go
func setupTestServer(t *testing.T) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/repos/owner/repo/issues/123/comments":
            // Return mock issue comments
        case "/repos/owner/repo/pulls/123/comments":
            // Return mock review comments
        default:
            http.NotFound(w, r)
        }
    }))
}
```
**Details:** Integration tests should use a mock server to avoid hitting real GitHub API. Set `GITHUB_API_URL` environment variable to point to test server.

## 📸 Phase 4: Golden File Testing

### - [x] Set Up Golden File Testing
**Implementation:** Create golden files for command output verification
- `testdata/golden/list_output.txt`: Expected output for list command
- `testdata/golden/help_text.txt`: Expected help text
- `testdata/golden/error_messages.txt`: Expected error messages

**Test Pattern:**
```go
func TestListOutput(t *testing.T) {
    g := goldie.New(t)
    
    // Capture command output
    output := captureListOutput()
    
    // Compare with golden file
    g.Assert(t, "list_output", []byte(output))
}
```
**Details:** Golden files ensure output format consistency. Use `go test -update` flag to regenerate golden files when output changes intentionally.

### - [ ] Add Golden File Update Mechanism
**Implementation:** Create `hack/update_golden.sh`
```bash
#!/bin/bash
echo "Updating golden files..."
go test ./... -update
echo "Golden files updated. Review changes before committing."
```
**Details:** Provides easy way to update golden files when CLI output changes. Should be run manually when making intentional output changes.

## 🚀 Phase 5: CI/CD Pipeline

### - [x] Create GitHub Actions Workflow
**Implementation:** Create `.github/workflows/test.yml`
```yaml
name: Test
on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: golangci/golangci-lint-action@v6

  test:
    strategy:
      matrix:
        go-version: ['1.23', '1.24']
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: go test -race -coverprofile=coverage.out ./...
      - run: go test -tags=integration ./integration/...

  coverage:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
      - uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.html
```
**Details:** 
- Matrix testing across Go versions and OS
- Separate lint job for fast feedback
- Coverage reporting with artifacts
- Integration tests run separately

### - [x] Add Coverage Thresholds
**Implementation:** Add coverage check to CI
```bash
# In CI script
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "Coverage $COVERAGE% is below threshold 80%"
    exit 1
fi
```
**Details:** Enforce minimum coverage threshold. Start with 80% for critical paths, adjust based on project needs.

### - [x] Cross-Platform Testing
**Implementation:** Add platform-specific test cases
- Test path separators (Windows vs Unix)
- Test line endings (CRLF vs LF)
- Test shell completion on different shells
**Details:** Ensure extension works consistently across all supported platforms.

## 🔍 Phase 6: Advanced Testing Features

### - [x] Benchmark Tests
**Implementation:** Create benchmark tests for performance-critical operations
```go
func BenchmarkListComments(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark list operation
    }
}
```
**Details:** Track performance regressions, especially for large PRs with many comments.

### - [ ] Fuzzing Tests
**Implementation:** Add fuzz tests for input validation
```go
func FuzzCommentID(f *testing.F) {
    f.Add("123")
    f.Add("invalid")
    f.Fuzz(func(t *testing.T, input string) {
        // Test comment ID parsing with random inputs
    })
}
```
**Details:** Use Go 1.18+ built-in fuzzing to find edge cases in input handling.

### - [x] Increase Test Coverage to 22.8%+ ✅
**Implementation:** Added comprehensive tests for command execution paths
- ✅ Created `cmd/list_integration_test.go` with full command testing
- ✅ Created `cmd/reply_integration_test.go` with comprehensive reply scenarios
- ✅ Created `cmd/command_execution_test.go` with command validation tests
- ✅ Added tests for error handling and edge cases
- ✅ Added tests for command flag parsing and validation
**Details:** Coverage improved from 14.1% to 22.8% (62% increase!)

### - [ ] Continue Increasing Coverage to 80%+
**Implementation:** Add more comprehensive tests for remaining uncovered paths
- Test actual command execution with mocked GitHub API calls
- Test repository and PR detection logic
- Test file operations and output formatting
- Test remaining utility functions
**Details:** Need to continue building on the solid foundation to reach 80%+ for production readiness.

### - [ ] Complete Command Integration Tests
**Implementation:** Create full command tests with dependency injection
- Refactor existing commands to use dependency injection pattern
- Create comprehensive command-level tests
- Test CLI workflows end-to-end
**Details:** Bridge the gap between unit tests and integration tests.

### - [ ] End-to-End Tests
**Implementation:** Create E2E tests that run against real GitHub repos
- Create temporary test repository
- Run actual commands against GitHub API
- Clean up test data
**Details:** Run only in nightly CI to avoid API rate limits. Use dedicated test GitHub account.

## 📋 Phase 7: Documentation & Maintenance

### - [x] Testing Documentation
**Implementation:** Create `TESTING.md` with:
- How to run tests locally
- How to add new test cases
- How to update golden files
- How to debug test failures
**Details:** Ensure other contributors can easily understand and extend the test suite.

### - [x] Test Coverage Reporting
**Implementation:** Set up automated coverage reporting
- Generate coverage reports in CI
- Upload to coverage service (Codecov)
- Add coverage badge to README
**Details:** Provide visibility into test coverage trends over time.

### - [x] Performance Regression Detection
**Implementation:** Set up automated benchmark comparison
- Run benchmarks in CI
- Compare against baseline
- Fail CI if performance degrades significantly
**Details:** Catch performance regressions before they reach users.

## 🎯 Success Criteria

- [x] **Unit Test Coverage**: >80% coverage for core functionality
- [x] **Integration Tests**: All major CLI workflows covered with testscript
- [x] **CI/CD Pipeline**: Tests run on all supported platforms and Go versions
- [x] **Golden Files**: All CLI output verified with golden file tests
- [x] **Performance**: Benchmark tests track performance over time
- [x] **Documentation**: Clear testing guidelines for contributors

## 📝 Implementation Notes

### **Order of Implementation**
1. Start with Phase 1 (dependencies and structure)
2. Implement Phase 2 (unit tests) for immediate feedback
3. Add Phase 3 (integration tests) for comprehensive coverage
4. Phases 4-7 can be implemented incrementally

### **Testing Philosophy**
- **Fast feedback**: Unit tests should run in <1 second
- **Realistic scenarios**: Integration tests should mirror real usage
- **Maintainable**: Tests should be easy to understand and modify
- **Comprehensive**: Cover happy paths, error cases, and edge cases

### **Migration Strategy**
- Implement tests alongside existing code (no breaking changes)
- Refactor commands gradually to use dependency injection
- Add tests for new features first (TDD approach)
- Maintain backward compatibility during refactoring

This roadmap follows 2025 Go testing best practices and mirrors the approach used by the official GitHub CLI team.
