# Contributing to gh-comment

Thank you for your interest in contributing to gh-comment! This guide will help you get started with development, testing, and submitting contributions.

## Development Setup

### Prerequisites

1. **Go 1.21+** - Install from [golang.org](https://golang.org/doc/install)
2. **GitHub CLI** - Install from [cli.github.com](https://cli.github.com/)
3. **Git** - For version control

### Getting Started

```bash
# Clone the repository
git clone https://github.com/silouanwright/gh-comment.git
cd gh-comment

# Install dependencies
go mod download

# Build the project
go build

# Run tests
go test ./cmd -v

# Install locally for testing
gh extension install .
```

## Project Structure

```
├── cmd/                 # Command implementations
│   ├── add.go          # Add command
│   ├── review.go       # Review command
│   └── ...             # Other commands
├── internal/github/     # GitHub API client
│   ├── client.go       # Interface definitions
│   ├── real_client.go  # Production client
│   └── mock_client.go  # Test client
├── docs/               # Documentation
└── README.md           # Main documentation
```

## Development Workflow

### Making Changes

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Follow the established patterns**:
   - All commands use dependency injection with `GitHubAPI` interface
   - Use the `MockClient` for testing
   - Follow the table-driven test pattern
   - Add comprehensive error handling

3. **Write tests** for new functionality:
   ```go
   func TestYourFeature(t *testing.T) {
       // Save original client
       originalClient := yourClient
       defer func() { yourClient = originalClient }()
       
       // Set up mock
       mockClient := github.NewMockClient()
       yourClient = mockClient
       
       // Your test logic
   }
   ```

### Code Standards

- **Go formatting**: Run `go fmt ./...` before committing
- **Error handling**: All errors should provide actionable guidance
- **Help text**: Keep examples working and realistic
- **Commit messages**: Use [Conventional Commits](https://www.conventionalcommits.org/)

### Testing

We have comprehensive testing at multiple levels:

#### Unit Tests
```bash
# Run all tests
go test ./cmd -v

# Run specific test
go test ./cmd -run TestAddCommand -v

# Run with coverage
go test ./cmd -cover
```

#### Integration Tests
```bash
# Build and test with real GitHub PRs
go build
./gh-comment list <your-test-pr>
```

See [docs/testing/TESTING_GUIDE.md](testing/TESTING_GUIDE.md) for comprehensive testing documentation.

## Architecture Principles

### Dependency Injection Pattern

All commands follow this pattern for testability:

```go
var commandClient github.GitHubAPI

func runCommand(cmd *cobra.Command, args []string) error {
    if commandClient == nil {
        commandClient = &github.RealClient{}
    }
    // Use commandClient for all operations
}
```

### Error Handling

Provide actionable error messages:

```go
func formatActionableError(operation string, err error) error {
    // Analyze error and provide helpful suggestions
    return fmt.Errorf("error during %s: %w\n\n💡 Suggestions:\n  • Check X\n  • Try Y", operation, err)
}
```

### Help Text Guidelines

- **Examples must work**: All help text examples should be copy/pasteable
- **Use realistic values**: Avoid placeholders like `<file>` in examples
- **Show common use cases**: Focus on what users actually do
- **Progressive complexity**: Simple examples first, advanced later

## Submitting Changes

### Pull Request Process

1. **Ensure tests pass**:
   ```bash
   go test ./cmd
   go build
   ```

2. **Update documentation** if needed:
   - Update help text for command changes
   - Update README.md for new features
   - Add/update tests

3. **Create pull request**:
   - Use descriptive title
   - Explain the problem and solution
   - Reference any related issues
   - Include testing evidence

### Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): brief description

feat(review): add optional review body support
fix(batch): resolve GitHub API structure issue  
docs(readme): simplify quick start section
test(add): add comprehensive edge case coverage
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`

## Common Development Tasks

### Adding a New Command

1. Create `cmd/newcommand.go` with the command structure
2. Add dependency injection pattern
3. Create `cmd/newcommand_test.go` with comprehensive tests
4. Update help text and examples
5. Add to command list in README.md

### Modifying GitHub API Interactions

1. Update the interface in `internal/github/client.go`
2. Implement in `internal/github/real_client.go`
3. Add mock implementation in the `MockClient`
4. Update all affected commands
5. Add tests for the new functionality

### Improving Error Messages

1. Identify common error scenarios
2. Add pattern matching in `formatActionableError`
3. Provide specific, actionable guidance
4. Test with real error conditions

## Getting Help

- **GitHub Issues**: Report bugs or request features
- **Discussions**: Ask questions or propose ideas
- **Code Review**: All PRs get thorough review and feedback

## Code of Conduct

- **Be respectful**: Treat all contributors with respect
- **Be constructive**: Provide helpful feedback
- **Be collaborative**: Work together toward better solutions
- **Follow standards**: Maintain code quality and consistency

## Release Process

Releases are handled by maintainers:

1. Version bump in appropriate files
2. Update CHANGELOG.md
3. Create GitHub release with binaries
4. Update documentation if needed

Thank you for contributing to gh-comment!