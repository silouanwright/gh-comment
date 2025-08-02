# Makefile for gh-comment
# Follows GitHub CLI and Kubernetes patterns for test organization

.PHONY: help test test-unit test-integration test-all coverage coverage-html build clean

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Build targets
build: ## Build the binary
	@echo "🔨 Building gh-comment..."
	@go build -o gh-comment .
	@echo "✅ Build completed: ./gh-comment"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -f gh-comment test-build coverage.out coverage.filtered.out coverage.html
	@echo "✅ Clean completed"

# Test targets
test: test-unit ## Run unit tests (default)

test-unit: ## Run unit tests with coverage (excludes integration tests)
	@echo "🧪 Running unit tests..."
	@go test -cover -coverprofile=coverage.out -coverpkg=./cmd,./internal/github ./cmd ./internal/...
	@echo "📊 Filtering coverage (excluding test utilities)..."
	@grep -v -E "(integration-scenarios|test-integration|testutil|_mock\.go|test/)" coverage.out > coverage.filtered.out || true
	@echo "📈 Coverage Report:"
	@go tool cover -func coverage.filtered.out | tail -1

test-integration: ## Run integration tests (requires GitHub token)
	@echo "🌐 Running integration tests..."
	@if [ -z "$$GITHUB_TOKEN" ]; then \
		echo "⚠️  GITHUB_TOKEN not set - integration tests will be skipped"; \
		echo "   Set GITHUB_TOKEN to run real GitHub API tests"; \
	fi
	@go test -tags=integration -v ./cmd ./test/integration/

test-all: test-unit test-integration ## Run all tests (unit + integration)
	@echo "✅ All tests completed"

# Coverage targets
coverage: test-unit ## Generate coverage report (unit tests only)
	@echo "📊 Unit test coverage:"
	@go tool cover -func coverage.filtered.out

coverage-html: test-unit ## Generate HTML coverage report
	@echo "🌐 Generating HTML coverage report..."
	@go tool cover -html=coverage.filtered.out -o coverage.html
	@echo "📄 Coverage report generated: coverage.html"
	@echo "💡 Open with: open coverage.html"

# Development helpers
test-watch: ## Watch for changes and run unit tests
	@echo "👀 Watching for changes..."
	@which fswatch > /dev/null || (echo "❌ fswatch not found. Install with: brew install fswatch" && exit 1)
	@fswatch -o . | xargs -n1 -I{} make test-unit

lint: ## Run linter
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "❌ golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run

format: ## Format code
	@echo "💄 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

# CI/CD helpers
ci-test: ## Run tests in CI environment
	@echo "🤖 Running CI tests..."
	@make test-unit
	@echo "📊 Final coverage:"
	@go tool cover -func coverage.filtered.out | tail -1

# Coverage thresholds (matching industry standards)
coverage-check: test-unit ## Check coverage meets thresholds
	@echo "🎯 Checking coverage thresholds..."
	@COVERAGE=$$(go tool cover -func coverage.filtered.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE < 80" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below 80% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets 80% threshold"; \
	fi

# Integration test helpers
integration-dry-run: ## Test integration framework without creating PRs
	@echo "🏃 Testing integration framework (dry run)..."
	@go run -tags=integration . test-integration --dry-run --scenario=comments

integration-inspect: ## Run integration tests and leave PR open for inspection
	@echo "🔍 Running integration tests with inspection..."
	@go run -tags=integration . test-integration --inspect --scenario=comments

# Development status
status: ## Show project status
	@echo "📊 Project Status:"
	@echo "  Repository: $$(git remote get-url origin 2>/dev/null || echo 'Not a git repository')"
	@echo "  Branch: $$(git branch --show-current 2>/dev/null || echo 'Unknown')"
	@echo "  Go version: $$(go version)"
	@echo "  Build status: $$(if [ -f gh-comment ]; then echo '✅ Built'; else echo '❌ Not built'; fi)"
	@if [ -f coverage.filtered.out ]; then \
		echo "  Coverage: $$(go tool cover -func coverage.filtered.out | tail -1 | awk '{print $$3}')"; \
	else \
		echo "  Coverage: ❓ Run 'make coverage' to generate"; \
	fi

# Default target
.DEFAULT_GOAL := help
