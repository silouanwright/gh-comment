# Testing Best Practices for gh-comment

This guide documents the testing patterns, strategies, and best practices established in the gh-comment project. It serves as a reference for maintaining high test coverage and code quality.

## Overview

The gh-comment project has achieved **77.0% test coverage** through comprehensive testing strategies that include:

- **Dependency Injection** for testability
- **Mock clients** for GitHub API interactions  
- **Table-driven tests** for comprehensive coverage
- **Edge case testing** with validation
- **Integration testing** with real command execution
- **Fuzz testing** for robustness
- **Output capture testing** for CLI functions

## Architecture Patterns

### Dependency Injection Pattern

All commands use dependency injection to enable testability:

```go
var (
    // Client for dependency injection (tests can override)
    commandClient github.GitHubAPI
)

func runCommand(cmd *cobra.Command, args []string) error {
    // Initialize client if not set (production use)
    if commandClient == nil {
        commandClient = &github.RealClient{}
    }
    
    // Use injected client for all API operations
    result, err := commandClient.SomeOperation(owner, repo, data)
    // ... rest of function
}
```

**Benefits:**
- Commands are fully unit testable
- No external API calls during tests
- Predictable test behavior
- Easy error simulation

### Mock Client Implementation

The `MockClient` implements the `GitHubAPI` interface for testing:

```go
type MockClient struct {
    // Data to return
    IssueComments  []Comment
    ReviewComments []Comment
    CreatedComment *Comment
    
    // Error simulation
    ListIssueCommentsError error
    CreateCommentError     error
    // ... other error fields
}
```

**Usage in tests:**
```go
func TestCommand(t *testing.T) {
    mockClient := github.NewMockClient()
    commandClient = mockClient
    
    // Configure mock behavior
    mockClient.CreateCommentError = errors.New("API error")
    
    // Test the command
    err := runCommand(nil, []string{"arg1", "arg2"})
    assert.Error(t, err)
}
```

## Testing Patterns

### Table-Driven Tests

Use table-driven tests for comprehensive scenario coverage:

```go
func TestCommandScenarios(t *testing.T) {
    tests := []struct {
        name           string
        args           []string
        setupMock      func(*github.MockClient)
        wantErr        bool
        expectedErrMsg string
    }{
        {
            name: "successful operation",
            args: []string{"123", "message"},
            setupMock: func(m *github.MockClient) {
                // Configure successful mock behavior
            },
            wantErr: false,
        },
        {
            name: "API error",
            args: []string{"123", "message"},
            setupMock: func(m *github.MockClient) {
                m.CreateCommentError = errors.New("API failed")
            },
            wantErr: true,
            expectedErrMsg: "API failed",
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := github.NewMockClient()
            if tt.setupMock != nil {
                tt.setupMock(mockClient)
            }
            commandClient = mockClient

            err := runCommand(nil, tt.args)
            if tt.wantErr {
                assert.Error(t, err)
                if tt.expectedErrMsg != "" {
                    assert.Contains(t, err.Error(), tt.expectedErrMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Test Isolation and Cleanup

Always save and restore original state:

```go
func TestCommand(t *testing.T) {
    // Save original state
    originalClient := commandClient
    originalRepo := repo
    originalPR := prNumber
    defer func() {
        commandClient = originalClient
        repo = originalRepo
        prNumber = originalPR
    }()

    // Set up test state
    mockClient := github.NewMockClient()
    commandClient = mockClient
    repo = "owner/repo"
    prNumber = 123

    // Run test
    // ...
}
```

### Error Testing Patterns

Test all error paths with specific error checking:

```go
func TestErrorScenarios(t *testing.T) {
    tests := []struct {
        name              string
        setupMockError    func(*github.MockClient)
        expectedErrMsg    string
    }{
        {
            name: "find comment error",
            setupMockError: func(m *github.MockClient) {
                m.FindCommentError = assert.AnError
            },
            expectedErrMsg: "failed to find comment",
        },
        {
            name: "create comment error", 
            setupMockError: func(m *github.MockClient) {
                m.CreateCommentError = assert.AnError
            },
            expectedErrMsg: "failed to create comment",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockClient := github.NewMockClient()
            tt.setupMockError(mockClient)
            commandClient = mockClient

            err := runCommand(nil, []string{"123"})
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedErrMsg)
        })
    }
}
```

### Validation Testing

Test input validation comprehensively:

```go
func TestInputValidation(t *testing.T) {
    tests := []struct {
        name           string
        input          string
        wantErr        bool
        expectedErrMsg string
    }{
        {"valid input", "123", false, ""},
        {"invalid number", "abc", true, "must be a valid integer"},
        {"empty input", "", true, "must be a valid integer"},
        {"negative number", "-1", false, ""}, // if negative is valid
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateInput(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                if tt.expectedErrMsg != "" {
                    assert.Contains(t, err.Error(), tt.expectedErrMsg)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Advanced Testing Techniques

### Output Capture Testing

For functions that print to stdout:

```go
func captureOutput(fn func()) string {
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    outputChan := make(chan string)
    go func() {
        var buf bytes.Buffer
        buf.ReadFrom(r)
        outputChan <- buf.String()
    }()

    fn()

    w.Close()
    os.Stdout = oldStdout
    return <-outputChan
}

func TestDisplayFunction(t *testing.T) {
    output := captureOutput(func() {
        displaySomething("test data")
    })
    
    assert.Contains(t, output, "expected content")
    assert.True(t, strings.HasSuffix(output, "\n"))
}
```

### Fuzz Testing

Add fuzz tests for robust input handling:

```go
func FuzzCommentID(f *testing.F) {
    // Seed with known inputs
    f.Add("123")
    f.Add("0")
    f.Add("-1")
    f.Add("abc")
    f.Add("")

    f.Fuzz(func(t *testing.T, input string) {
        // Function should not panic
        defer func() {
            if r := recover(); r != nil {
                t.Errorf("Function panicked with input %q: %v", input, r)
            }
        }()

        validateCommentID(input)
    })
}
```

### Integration Testing

Test real command execution:

```go
func TestCommandIntegration(t *testing.T) {
    // Set up test environment
    mockClient := github.NewMockClient()
    commandClient = mockClient
    
    // Configure expected behavior
    mockClient.IssueComments = []github.Comment{
        {ID: 123, Body: "test comment"},
    }

    // Execute command
    err := runCommand(nil, []string{"list", "123"})
    assert.NoError(t, err)
    
    // Verify mock was called correctly
    // (could add call tracking to mock if needed)
}
```

## Coverage Strategies

### Achieving High Coverage

1. **Test all public functions** - Every exported function should have tests
2. **Test all code paths** - Use conditional tests for different branches
3. **Test error paths** - Ensure error handling works correctly
4. **Test edge cases** - Empty inputs, boundary values, special characters
5. **Test user workflows** - Integration tests for common use cases

### Coverage Analysis

Monitor coverage with:

```bash
# Generate coverage profile
go test ./cmd -coverprofile=coverage.out

# View coverage percentage
go test ./cmd -cover

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# View detailed coverage by function
go tool cover -func=coverage.out
```

### Identifying Coverage Gaps

Look for:
- **Uncovered error paths** - Add error simulation tests
- **Unused helper functions** - Either test or remove
- **Complex conditionals** - Add tests for all branches
- **Init functions** - Consider testable alternatives

## Test Organization

### File Structure

```
cmd/
├── command.go              # Command implementation
├── command_test.go         # Main unit tests
├── command_integration_test.go  # Integration tests (if needed)
├── helpers.go              # Helper functions
├── helpers_test.go         # Helper function tests
└── fuzz_test.go           # Fuzz tests (shared file)
```

### Test Naming

- **Test functions**: `TestFunctionName`
- **Subtests**: Descriptive names with underscores
- **Test files**: `*_test.go` 
- **Integration tests**: `*_integration_test.go`
- **Fuzz tests**: `FuzzFunctionName`

### Test Categories

1. **Unit Tests** - Test individual functions in isolation
2. **Integration Tests** - Test command execution end-to-end
3. **Error Tests** - Test all error conditions
4. **Edge Case Tests** - Test boundary conditions
5. **Fuzz Tests** - Test with random inputs
6. **Benchmark Tests** - Performance testing (if needed)

## Common Pitfalls

### Avoid These Patterns

```go
// DON'T: Test implementation details
func TestInternal(t *testing.T) {
    // Testing private variables or internal state
}

// DON'T: Tests that depend on external state
func TestWithRealAPI(t *testing.T) {
    // Making real API calls
}

// DON'T: Tests without cleanup
func TestWithoutCleanup(t *testing.T) {
    // Modifying global state without restoration
}

// DON'T: Weak assertions
func TestWeak(t *testing.T) {
    result := someFunction()
    assert.NotNil(t, result) // Too generic
}
```

### Best Practices

```go
// DO: Test public behavior
func TestPublicBehavior(t *testing.T) {
    input := "test"
    result, err := PublicFunction(input)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
}

// DO: Use dependency injection
func TestWithMocks(t *testing.T) {
    mockClient := github.NewMockClient()
    commandClient = mockClient
    // ... test with predictable mock behavior
}

// DO: Clean up state
func TestWithCleanup(t *testing.T) {
    original := globalVar
    defer func() { globalVar = original }()
    // ... modify state for test
}

// DO: Strong assertions  
func TestStrong(t *testing.T) {
    result, err := someFunction()
    assert.NoError(t, err)
    assert.Equal(t, "expected value", result.SpecificField)
    assert.Len(t, result.Items, 3)
}
```

## Continuous Improvement

### Metrics to Track

- **Coverage percentage** - Aim for >75%
- **Test execution time** - Keep tests fast
- **Test reliability** - No flaky tests
- **Code quality** - High signal-to-noise ratio

### Regular Reviews

1. **Add tests for new features** - Before implementation
2. **Update tests for changes** - Keep tests current  
3. **Refactor test code** - Maintain readability
4. **Remove obsolete tests** - Clean up unused tests

## Tools and Libraries

### Essential Testing Tools

- **testify/assert** - Rich assertion library
- **testify/require** - Failing assertions  
- **testify/mock** - Mock generation (if needed)
- **go test** - Built-in test runner
- **go tool cover** - Coverage analysis

### Useful Patterns

```go
// Setup helper for common test state
func setupTest(t *testing.T) (*github.MockClient, func()) {
    mockClient := github.NewMockClient()
    originalClient := commandClient
    commandClient = mockClient
    
    cleanup := func() {
        commandClient = originalClient
    }
    
    return mockClient, cleanup
}

// Usage in tests
func TestSomething(t *testing.T) {
    mockClient, cleanup := setupTest(t)
    defer cleanup()
    
    // Configure mock and run test
    mockClient.SomeField = "test value"
    // ... test code
}
```

## Conclusion

The gh-comment project demonstrates how systematic testing practices can achieve high coverage and maintainable code. The key principles are:

1. **Dependency injection** enables isolated testing
2. **Comprehensive test coverage** catches bugs early
3. **Consistent patterns** make tests maintainable
4. **Edge case testing** improves robustness
5. **Integration testing** validates user workflows

By following these patterns, the codebase maintains professional quality with confidence in refactoring and extending functionality.

---

*This guide reflects the testing practices established while improving gh-comment test coverage from 30.6% to 77.0%.*