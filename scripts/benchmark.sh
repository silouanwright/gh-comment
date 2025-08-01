#!/bin/bash

# Benchmark script for local performance testing
set -e

echo "🏃 Running benchmarks..."
echo

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if benchstat is installed
if ! command -v benchstat &> /dev/null; then
    echo -e "${YELLOW}benchstat not found. Installing...${NC}"
    go install golang.org/x/perf/cmd/benchstat@latest
fi

# Function to run benchmarks
run_benchmarks() {
    local output_file=$1
    local label=$2

    echo -e "${GREEN}Running benchmarks for: ${label}${NC}"
    go test -bench=. -benchmem -run=^$ -count=10 ./... > "$output_file" 2>&1

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Benchmarks completed${NC}"
    else
        echo -e "${RED}✗ Benchmark failed${NC}"
        cat "$output_file"
        exit 1
    fi
}

# Check for command line arguments
if [ "$1" == "compare" ]; then
    if [ -z "$2" ]; then
        echo "Usage: $0 compare <branch-or-commit>"
        echo "Example: $0 compare main"
        exit 1
    fi

    COMPARE_REF=$2

    # Save current branch
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

    # Run benchmarks on current branch
    run_benchmarks "current-bench.txt" "$CURRENT_BRANCH"

    # Stash any changes
    git stash push -m "benchmark-comparison" --include-untracked

    # Checkout comparison branch/commit
    echo -e "${YELLOW}Switching to $COMPARE_REF...${NC}"
    git checkout "$COMPARE_REF" --quiet

    # Run benchmarks on comparison branch
    run_benchmarks "compare-bench.txt" "$COMPARE_REF"

    # Return to original branch
    git checkout "$CURRENT_BRANCH" --quiet

    # Restore stashed changes if any
    if git stash list | grep -q "benchmark-comparison"; then
        git stash pop --quiet
    fi

    # Compare results
    echo
    echo -e "${GREEN}Benchmark Comparison: $COMPARE_REF -> $CURRENT_BRANCH${NC}"
    echo "================================================================"
    benchstat compare-bench.txt current-bench.txt

    # Clean up
    rm -f compare-bench.txt current-bench.txt

elif [ "$1" == "profile" ]; then
    # CPU profiling
    echo -e "${GREEN}Running CPU profile...${NC}"
    go test -bench=. -benchmem -cpuprofile=cpu.prof -run=^$ ./cmd

    echo -e "${GREEN}Running memory profile...${NC}"
    go test -bench=. -benchmem -memprofile=mem.prof -run=^$ ./cmd

    echo
    echo "View profiles with:"
    echo "  go tool pprof -http=:8080 cpu.prof"
    echo "  go tool pprof -http=:8081 mem.prof"

else
    # Default: just run benchmarks
    run_benchmarks "benchmark.txt" "current code"

    echo
    echo -e "${GREEN}Benchmark Results:${NC}"
    echo "=================="
    grep -E "Benchmark.*ns/op" benchmark.txt | column -t

    echo
    echo "Full results saved to: benchmark.txt"
    echo
    echo "Other options:"
    echo "  $0 compare <branch>  - Compare with another branch"
    echo "  $0 profile          - Generate CPU and memory profiles"
fi
