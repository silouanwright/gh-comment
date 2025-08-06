window.BENCHMARK_DATA = {
  "lastUpdate": 1754448356606,
  "repoUrl": "https://github.com/silouanwright/gh-comment",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "hello@silouan.io",
            "name": "Silouan Wright",
            "username": "silouanwright"
          },
          "committer": {
            "email": "hello@silouan.io",
            "name": "Silouan Wright",
            "username": "silouanwright"
          },
          "distinct": true,
          "id": "cfd245057c02f1806912c47deea4f2324e002e5e",
          "message": "feat: enhance review-reply command with intelligent error handling\n\nMajor improvements to address GitHub API 404 limitations:\n\n**Enhanced Error Handling:**\n- Detect GitHub API limitation where only 'standalone' review comments support replies\n- Comments from 'Start Review' → 'Submit Review' cannot be replied to (API limitation)\n- Distinguish between issue comment vs review comment type mismatches\n- Provide specific, actionable alternatives when threading fails\n\n**Intelligent Fallback Suggestions:**\n- Alternative 1: Add new comment at same location with file/line lookup\n- Alternative 2: Use emoji reactions for quick feedback\n- Alternative 3: Try resolving conversation (often works when replies don't)\n- Alternative 4: Use general PR discussion for issue comments\n\n**Improved Documentation:**\n- Updated help text to clearly explain GitHub API limitations\n- Added realistic examples showing fallback approaches\n- Clear explanation of 'standalone' vs 'review thread' comment types\n\n**Comprehensive Testing:**\n- Added TestHandleReviewReplyError for error handling scenarios\n- Added TestReviewReplyIntelligentErrorHandling for specific cases\n- Fixed existing tests to work with enhanced error handling\n\n**Root Cause Analysis:**\nBased on research, the 404 errors occur because:\n- GitHub's REST API only supports /replies endpoint for standalone comments\n- Review comments created via review flow are part of 'review threads'\n- Review threads have limited REST API support (better GraphQL support exists)\n- This is a known GitHub architectural limitation, not a tool issue\n\nUsers now get helpful guidance instead of cryptic 404 errors.",
          "timestamp": "2025-08-05T15:24:21-05:00",
          "tree_id": "c4359a8ff3609dfb8f726240574fb8902bb371fa",
          "url": "https://github.com/silouanwright/gh-comment/commit/cfd245057c02f1806912c47deea4f2324e002e5e"
        },
        "date": 1754425649778,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkListComments",
            "value": 14113,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "71370 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14113,
            "unit": "ns/op",
            "extra": "71370 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "71370 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "71370 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14295,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "82244 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14295,
            "unit": "ns/op",
            "extra": "82244 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "82244 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "82244 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14182,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "79420 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14182,
            "unit": "ns/op",
            "extra": "79420 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "79420 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "79420 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 15003,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "79762 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 15003,
            "unit": "ns/op",
            "extra": "79762 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "79762 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "79762 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14237,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "84508 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14237,
            "unit": "ns/op",
            "extra": "84508 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "84508 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "84508 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14614,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "81514 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14614,
            "unit": "ns/op",
            "extra": "81514 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "81514 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "81514 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14307,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "80973 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14307,
            "unit": "ns/op",
            "extra": "80973 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "80973 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "80973 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14201,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "80458 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14201,
            "unit": "ns/op",
            "extra": "80458 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "80458 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "80458 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14329,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "81091 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14329,
            "unit": "ns/op",
            "extra": "81091 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "81091 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "81091 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14253,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "77917 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14253,
            "unit": "ns/op",
            "extra": "77917 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "77917 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "77917 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 581.6,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2033559 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 581.6,
            "unit": "ns/op",
            "extra": "2033559 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2033559 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2033559 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 581.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2038110 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 581.8,
            "unit": "ns/op",
            "extra": "2038110 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2038110 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2038110 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 581.1,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2043529 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 581.1,
            "unit": "ns/op",
            "extra": "2043529 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2043529 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2043529 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 581.2,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2043487 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 581.2,
            "unit": "ns/op",
            "extra": "2043487 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2043487 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2043487 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 580.6,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2040500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 580.6,
            "unit": "ns/op",
            "extra": "2040500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2040500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2040500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 580.9,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2015983 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 580.9,
            "unit": "ns/op",
            "extra": "2015983 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2015983 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2015983 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 582.1,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2038980 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 582.1,
            "unit": "ns/op",
            "extra": "2038980 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2038980 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2038980 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 581.3,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2021157 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 581.3,
            "unit": "ns/op",
            "extra": "2021157 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2021157 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2021157 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 580.2,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2038075 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 580.2,
            "unit": "ns/op",
            "extra": "2038075 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2038075 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2038075 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 580.9,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2039420 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 580.9,
            "unit": "ns/op",
            "extra": "2039420 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2039420 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2039420 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3121,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3121,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3126,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3126,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3127,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3127,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3203,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3203,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3287,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3287,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3115,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3115,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3114,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3114,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3111,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3111,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3111,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3111,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.3118,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.3118,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 71096,
            "unit": "ns/op\t   67449 B/op\t     675 allocs/op",
            "extra": "16797 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 71096,
            "unit": "ns/op",
            "extra": "16797 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67449,
            "unit": "B/op",
            "extra": "16797 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "16797 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 70711,
            "unit": "ns/op\t   67454 B/op\t     675 allocs/op",
            "extra": "16938 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 70711,
            "unit": "ns/op",
            "extra": "16938 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67454,
            "unit": "B/op",
            "extra": "16938 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "16938 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 69105,
            "unit": "ns/op\t   67458 B/op\t     675 allocs/op",
            "extra": "17086 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 69105,
            "unit": "ns/op",
            "extra": "17086 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67458,
            "unit": "B/op",
            "extra": "17086 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17086 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 70105,
            "unit": "ns/op\t   67426 B/op\t     675 allocs/op",
            "extra": "17212 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 70105,
            "unit": "ns/op",
            "extra": "17212 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67426,
            "unit": "B/op",
            "extra": "17212 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17212 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 69825,
            "unit": "ns/op\t   67436 B/op\t     675 allocs/op",
            "extra": "17040 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 69825,
            "unit": "ns/op",
            "extra": "17040 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67436,
            "unit": "B/op",
            "extra": "17040 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17040 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 71133,
            "unit": "ns/op\t   67449 B/op\t     675 allocs/op",
            "extra": "17025 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 71133,
            "unit": "ns/op",
            "extra": "17025 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67449,
            "unit": "B/op",
            "extra": "17025 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17025 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 70012,
            "unit": "ns/op\t   67436 B/op\t     675 allocs/op",
            "extra": "16766 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 70012,
            "unit": "ns/op",
            "extra": "16766 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67436,
            "unit": "B/op",
            "extra": "16766 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "16766 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 69954,
            "unit": "ns/op\t   67445 B/op\t     675 allocs/op",
            "extra": "16929 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 69954,
            "unit": "ns/op",
            "extra": "16929 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67445,
            "unit": "B/op",
            "extra": "16929 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "16929 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 69888,
            "unit": "ns/op\t   67436 B/op\t     675 allocs/op",
            "extra": "17061 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 69888,
            "unit": "ns/op",
            "extra": "17061 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67436,
            "unit": "B/op",
            "extra": "17061 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17061 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 69951,
            "unit": "ns/op\t   67427 B/op\t     675 allocs/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 69951,
            "unit": "ns/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67427,
            "unit": "B/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3506860,
            "unit": "ns/op\t     516 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3506860,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 516,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3498329,
            "unit": "ns/op\t     405 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3498329,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 405,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3536529,
            "unit": "ns/op\t     627 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3536529,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 627,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3508001,
            "unit": "ns/op\t     516 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3508001,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 516,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3483906,
            "unit": "ns/op\t     406 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3483906,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 406,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3494128,
            "unit": "ns/op\t     295 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3494128,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 295,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3505410,
            "unit": "ns/op\t     515 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3505410,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 515,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3509544,
            "unit": "ns/op\t     515 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3509544,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 515,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3499603,
            "unit": "ns/op\t     276 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3499603,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 276,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3493851,
            "unit": "ns/op\t     406 B/op\t       5 allocs/op",
            "extra": "338 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3493851,
            "unit": "ns/op",
            "extra": "338 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 406,
            "unit": "B/op",
            "extra": "338 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "338 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 403258,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2724 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 403258,
            "unit": "ns/op",
            "extra": "2724 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2724 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2724 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 412660,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3116 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 412660,
            "unit": "ns/op",
            "extra": "3116 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3116 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3116 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 407685,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 407685,
            "unit": "ns/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3121 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 409267,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2946 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 409267,
            "unit": "ns/op",
            "extra": "2946 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2946 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2946 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 407087,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3074 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 407087,
            "unit": "ns/op",
            "extra": "3074 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3074 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3074 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 404712,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 404712,
            "unit": "ns/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 406786,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2650 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 406786,
            "unit": "ns/op",
            "extra": "2650 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2650 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2650 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 404830,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3168 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 404830,
            "unit": "ns/op",
            "extra": "3168 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3168 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3168 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 410939,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3120 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 410939,
            "unit": "ns/op",
            "extra": "3120 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3120 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3120 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 405884,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "3166 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 405884,
            "unit": "ns/op",
            "extra": "3166 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "3166 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "3166 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15490,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "73305 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15490,
            "unit": "ns/op",
            "extra": "73305 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "73305 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "73305 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15582,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "72584 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15582,
            "unit": "ns/op",
            "extra": "72584 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "72584 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "72584 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15555,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "74280 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15555,
            "unit": "ns/op",
            "extra": "74280 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "74280 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "74280 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15449,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "73557 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15449,
            "unit": "ns/op",
            "extra": "73557 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "73557 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "73557 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15467,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "75501 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15467,
            "unit": "ns/op",
            "extra": "75501 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "75501 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "75501 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15325,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "75522 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15325,
            "unit": "ns/op",
            "extra": "75522 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "75522 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "75522 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15401,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "74281 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15401,
            "unit": "ns/op",
            "extra": "74281 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "74281 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "74281 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15376,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "73923 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15376,
            "unit": "ns/op",
            "extra": "73923 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "73923 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "73923 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15401,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "74594 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15401,
            "unit": "ns/op",
            "extra": "74594 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "74594 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "74594 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15369,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "75105 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15369,
            "unit": "ns/op",
            "extra": "75105 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "75105 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "75105 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1348,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "936794 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1348,
            "unit": "ns/op",
            "extra": "936794 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "936794 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "936794 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1351,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "935479 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1351,
            "unit": "ns/op",
            "extra": "935479 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "935479 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "935479 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1345,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "958736 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1345,
            "unit": "ns/op",
            "extra": "958736 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "958736 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "958736 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1353,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "922958 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1353,
            "unit": "ns/op",
            "extra": "922958 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "922958 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "922958 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1348,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "962340 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1348,
            "unit": "ns/op",
            "extra": "962340 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "962340 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "962340 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1347,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "964342 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1347,
            "unit": "ns/op",
            "extra": "964342 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "964342 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "964342 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1347,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "960210 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1347,
            "unit": "ns/op",
            "extra": "960210 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "960210 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "960210 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1363,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "964152 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1363,
            "unit": "ns/op",
            "extra": "964152 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "964152 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "964152 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1359,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "962775 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1359,
            "unit": "ns/op",
            "extra": "962775 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "962775 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "962775 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1353,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "960458 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1353,
            "unit": "ns/op",
            "extra": "960458 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "960458 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "960458 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 450.4,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2583230 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 450.4,
            "unit": "ns/op",
            "extra": "2583230 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2583230 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2583230 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 453,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2572965 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 453,
            "unit": "ns/op",
            "extra": "2572965 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2572965 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2572965 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 450.4,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2609536 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 450.4,
            "unit": "ns/op",
            "extra": "2609536 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2609536 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2609536 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 452.6,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2569737 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 452.6,
            "unit": "ns/op",
            "extra": "2569737 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2569737 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2569737 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 453.6,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2614562 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 453.6,
            "unit": "ns/op",
            "extra": "2614562 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2614562 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2614562 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 458.8,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2599766 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 458.8,
            "unit": "ns/op",
            "extra": "2599766 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2599766 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2599766 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 452.3,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2597563 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 452.3,
            "unit": "ns/op",
            "extra": "2597563 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2597563 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2597563 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 454,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2591505 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 454,
            "unit": "ns/op",
            "extra": "2591505 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2591505 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2591505 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 451.2,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2558782 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 451.2,
            "unit": "ns/op",
            "extra": "2558782 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2558782 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2558782 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 451.8,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2595267 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 451.8,
            "unit": "ns/op",
            "extra": "2595267 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2595267 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2595267 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "hello@silouan.io",
            "name": "Silouan Wright",
            "username": "silouanwright"
          },
          "committer": {
            "email": "hello@silouan.io",
            "name": "Silouan Wright",
            "username": "silouanwright"
          },
          "distinct": true,
          "id": "45dd95ff2f1cc5f921a41a63d169c2c73e952dd5",
          "message": "fix: use proper Go formatting (tabs) in integration script\n\n- Fix src/main.go template to use tabs instead of spaces\n- Prevents formatting failures in CI when integration tests run\n- Ensures consistency between script-generated files and Go standards",
          "timestamp": "2025-08-05T21:41:04-05:00",
          "tree_id": "1eb08fd849b0b219b5266d11539ad99f6f0e4440",
          "url": "https://github.com/silouanwright/gh-comment/commit/45dd95ff2f1cc5f921a41a63d169c2c73e952dd5"
        },
        "date": 1754448207754,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkListComments",
            "value": 13549,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "75147 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13549,
            "unit": "ns/op",
            "extra": "75147 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "75147 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "75147 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13516,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "85375 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13516,
            "unit": "ns/op",
            "extra": "85375 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "85375 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "85375 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13951,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "85566 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13951,
            "unit": "ns/op",
            "extra": "85566 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "85566 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "85566 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13469,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "84577 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13469,
            "unit": "ns/op",
            "extra": "84577 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "84577 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "84577 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13523,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "82976 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13523,
            "unit": "ns/op",
            "extra": "82976 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "82976 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "82976 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13504,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "83925 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13504,
            "unit": "ns/op",
            "extra": "83925 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "83925 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "83925 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 15578,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "85447 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 15578,
            "unit": "ns/op",
            "extra": "85447 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "85447 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "85447 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13484,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "85653 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13484,
            "unit": "ns/op",
            "extra": "85653 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "85653 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "85653 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13535,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "87175 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13535,
            "unit": "ns/op",
            "extra": "87175 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "87175 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "87175 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 13544,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "84502 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 13544,
            "unit": "ns/op",
            "extra": "84502 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "84502 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "84502 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 602,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1981143 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 602,
            "unit": "ns/op",
            "extra": "1981143 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1981143 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1981143 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 589.9,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2001799 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 589.9,
            "unit": "ns/op",
            "extra": "2001799 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2001799 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2001799 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591.9,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1982797 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591.9,
            "unit": "ns/op",
            "extra": "1982797 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1982797 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1982797 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 590.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1987336 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 590.8,
            "unit": "ns/op",
            "extra": "1987336 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1987336 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1987336 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 596.4,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1999755 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 596.4,
            "unit": "ns/op",
            "extra": "1999755 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1999755 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1999755 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591.1,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2005940 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591.1,
            "unit": "ns/op",
            "extra": "2005940 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2005940 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2005940 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 594.5,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2003593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 594.5,
            "unit": "ns/op",
            "extra": "2003593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2003593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2003593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1981035 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591.8,
            "unit": "ns/op",
            "extra": "1981035 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1981035 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1981035 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 590,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1987729 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 590,
            "unit": "ns/op",
            "extra": "1987729 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1987729 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1987729 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591.6,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1967712 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591.6,
            "unit": "ns/op",
            "extra": "1967712 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1967712 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1967712 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6227,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6227,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.623,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.623,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6344,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6344,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6308,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6308,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6271,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6271,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6237,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6237,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6262,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6262,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.623,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.623,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.622,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.622,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.623,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.623,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 64862,
            "unit": "ns/op\t   67430 B/op\t     675 allocs/op",
            "extra": "18388 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 64862,
            "unit": "ns/op",
            "extra": "18388 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67430,
            "unit": "B/op",
            "extra": "18388 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18388 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 64625,
            "unit": "ns/op\t   67432 B/op\t     675 allocs/op",
            "extra": "18350 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 64625,
            "unit": "ns/op",
            "extra": "18350 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67432,
            "unit": "B/op",
            "extra": "18350 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18350 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 64609,
            "unit": "ns/op\t   67431 B/op\t     675 allocs/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 64609,
            "unit": "ns/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67431,
            "unit": "B/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 64562,
            "unit": "ns/op\t   67431 B/op\t     675 allocs/op",
            "extra": "18200 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 64562,
            "unit": "ns/op",
            "extra": "18200 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67431,
            "unit": "B/op",
            "extra": "18200 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18200 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 65071,
            "unit": "ns/op\t   67422 B/op\t     675 allocs/op",
            "extra": "18429 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 65071,
            "unit": "ns/op",
            "extra": "18429 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67422,
            "unit": "B/op",
            "extra": "18429 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18429 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 65061,
            "unit": "ns/op\t   67441 B/op\t     675 allocs/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 65061,
            "unit": "ns/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67441,
            "unit": "B/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18277 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 65204,
            "unit": "ns/op\t   67420 B/op\t     675 allocs/op",
            "extra": "18397 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 65204,
            "unit": "ns/op",
            "extra": "18397 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67420,
            "unit": "B/op",
            "extra": "18397 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18397 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 64583,
            "unit": "ns/op\t   67419 B/op\t     675 allocs/op",
            "extra": "18187 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 64583,
            "unit": "ns/op",
            "extra": "18187 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67419,
            "unit": "B/op",
            "extra": "18187 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18187 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 65076,
            "unit": "ns/op\t   67442 B/op\t     675 allocs/op",
            "extra": "18392 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 65076,
            "unit": "ns/op",
            "extra": "18392 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67442,
            "unit": "B/op",
            "extra": "18392 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18392 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 64672,
            "unit": "ns/op\t   67426 B/op\t     675 allocs/op",
            "extra": "18462 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 64672,
            "unit": "ns/op",
            "extra": "18462 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67426,
            "unit": "B/op",
            "extra": "18462 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18462 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3496483,
            "unit": "ns/op\t     276 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3496483,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 276,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3486749,
            "unit": "ns/op\t     514 B/op\t       5 allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3486749,
            "unit": "ns/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 514,
            "unit": "B/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3484655,
            "unit": "ns/op\t     513 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3484655,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 513,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3480868,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3480868,
            "unit": "ns/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3484988,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3484988,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3487068,
            "unit": "ns/op\t     276 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3487068,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 276,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3484903,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3484903,
            "unit": "ns/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3483829,
            "unit": "ns/op\t     276 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3483829,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 276,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3482862,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3482862,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3487754,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3487754,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395291,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2852 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395291,
            "unit": "ns/op",
            "extra": "2852 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2852 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2852 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 396152,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2811 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 396152,
            "unit": "ns/op",
            "extra": "2811 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2811 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2811 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 396564,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2750 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 396564,
            "unit": "ns/op",
            "extra": "2750 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2750 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2750 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395089,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2671 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395089,
            "unit": "ns/op",
            "extra": "2671 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2671 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2671 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395017,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2686 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395017,
            "unit": "ns/op",
            "extra": "2686 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2686 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2686 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 406762,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2721 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 406762,
            "unit": "ns/op",
            "extra": "2721 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2721 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2721 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 398464,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 398464,
            "unit": "ns/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 398081,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2752 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 398081,
            "unit": "ns/op",
            "extra": "2752 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2752 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2752 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 396254,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2892 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 396254,
            "unit": "ns/op",
            "extra": "2892 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2892 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2892 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395405,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2910 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395405,
            "unit": "ns/op",
            "extra": "2910 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2910 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15098,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76036 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15098,
            "unit": "ns/op",
            "extra": "76036 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76036 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76036 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15005,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76318 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15005,
            "unit": "ns/op",
            "extra": "76318 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76318 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76318 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15023,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76668 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15023,
            "unit": "ns/op",
            "extra": "76668 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76668 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76668 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15014,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15014,
            "unit": "ns/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15186,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76210 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15186,
            "unit": "ns/op",
            "extra": "76210 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76210 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76210 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 14995,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "77156 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 14995,
            "unit": "ns/op",
            "extra": "77156 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "77156 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "77156 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 14935,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "75910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 14935,
            "unit": "ns/op",
            "extra": "75910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "75910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "75910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15014,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "74593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15014,
            "unit": "ns/op",
            "extra": "74593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "74593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "74593 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 14980,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76366 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 14980,
            "unit": "ns/op",
            "extra": "76366 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76366 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76366 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 14976,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "77388 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 14976,
            "unit": "ns/op",
            "extra": "77388 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "77388 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "77388 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1376,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "870374 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1376,
            "unit": "ns/op",
            "extra": "870374 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "870374 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "870374 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1377,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "912322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1377,
            "unit": "ns/op",
            "extra": "912322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "912322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "912322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1380,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "925100 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1380,
            "unit": "ns/op",
            "extra": "925100 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "925100 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "925100 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1376,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "888307 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1376,
            "unit": "ns/op",
            "extra": "888307 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "888307 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "888307 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1371,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "941971 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1371,
            "unit": "ns/op",
            "extra": "941971 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "941971 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "941971 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1377,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "923505 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1377,
            "unit": "ns/op",
            "extra": "923505 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "923505 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "923505 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1384,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "897586 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1384,
            "unit": "ns/op",
            "extra": "897586 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "897586 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "897586 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1374,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "860532 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1374,
            "unit": "ns/op",
            "extra": "860532 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "860532 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "860532 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1373,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "890650 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1373,
            "unit": "ns/op",
            "extra": "890650 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "890650 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "890650 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1373,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "905972 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1373,
            "unit": "ns/op",
            "extra": "905972 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "905972 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "905972 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 460.3,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2596347 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 460.3,
            "unit": "ns/op",
            "extra": "2596347 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2596347 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2596347 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 461.5,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2581302 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 461.5,
            "unit": "ns/op",
            "extra": "2581302 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2581302 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2581302 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 460.1,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2570349 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 460.1,
            "unit": "ns/op",
            "extra": "2570349 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2570349 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2570349 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 461.7,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2575671 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 461.7,
            "unit": "ns/op",
            "extra": "2575671 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2575671 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2575671 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 460.1,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2556424 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 460.1,
            "unit": "ns/op",
            "extra": "2556424 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2556424 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2556424 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 462.5,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2569292 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 462.5,
            "unit": "ns/op",
            "extra": "2569292 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2569292 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2569292 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 460,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2566495 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 460,
            "unit": "ns/op",
            "extra": "2566495 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2566495 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2566495 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 464.9,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2544687 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 464.9,
            "unit": "ns/op",
            "extra": "2544687 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2544687 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2544687 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 462.1,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2565244 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 462.1,
            "unit": "ns/op",
            "extra": "2565244 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2565244 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2565244 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 460.2,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2554566 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 460.2,
            "unit": "ns/op",
            "extra": "2554566 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2554566 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2554566 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "hello@silouan.io",
            "name": "Silouan Wright",
            "username": "silouanwright"
          },
          "committer": {
            "email": "hello@silouan.io",
            "name": "Silouan Wright",
            "username": "silouanwright"
          },
          "distinct": true,
          "id": "49189b63d79db0b2f846ee0619a3f50841e0351c",
          "message": "fix: remove --no-verify from integration script commits\n\n- Integration script was bypassing pre-commit hooks entirely\n- This caused formatting issues to go undetected locally but fail in CI\n- Now integration script will run through proper formatting validation\n- Ensures consistency between local pre-commit and CI checks",
          "timestamp": "2025-08-05T21:43:34-05:00",
          "tree_id": "a6c514999d38a553969a9a020cd203727c5e1b52",
          "url": "https://github.com/silouanwright/gh-comment/commit/49189b63d79db0b2f846ee0619a3f50841e0351c"
        },
        "date": 1754448356302,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkListComments",
            "value": 14097,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "72927 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14097,
            "unit": "ns/op",
            "extra": "72927 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "72927 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "72927 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14584,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "79070 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14584,
            "unit": "ns/op",
            "extra": "79070 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "79070 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "79070 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14017,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "81236 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14017,
            "unit": "ns/op",
            "extra": "81236 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "81236 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "81236 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14245,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "82010 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14245,
            "unit": "ns/op",
            "extra": "82010 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "82010 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "82010 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14152,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "81598 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14152,
            "unit": "ns/op",
            "extra": "81598 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "81598 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "81598 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14196,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "80168 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14196,
            "unit": "ns/op",
            "extra": "80168 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "80168 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "80168 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14085,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "80907 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14085,
            "unit": "ns/op",
            "extra": "80907 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "80907 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "80907 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14122,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "81735 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14122,
            "unit": "ns/op",
            "extra": "81735 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "81735 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "81735 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14170,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "77894 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14170,
            "unit": "ns/op",
            "extra": "77894 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "77894 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "77894 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments",
            "value": 14234,
            "unit": "ns/op\t   56688 B/op\t       8 allocs/op",
            "extra": "78346 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - ns/op",
            "value": 14234,
            "unit": "ns/op",
            "extra": "78346 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - B/op",
            "value": 56688,
            "unit": "B/op",
            "extra": "78346 times\n4 procs"
          },
          {
            "name": "BenchmarkListComments - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "78346 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 594,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1995927 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 594,
            "unit": "ns/op",
            "extra": "1995927 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1995927 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1995927 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2005096 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591,
            "unit": "ns/op",
            "extra": "2005096 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2005096 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2005096 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 590.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2005791 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 590.8,
            "unit": "ns/op",
            "extra": "2005791 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2005791 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2005791 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591.2,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2002910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591.2,
            "unit": "ns/op",
            "extra": "2002910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2002910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2002910 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 591,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2005404 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 591,
            "unit": "ns/op",
            "extra": "2005404 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2005404 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2005404 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 600.1,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1996192 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 600.1,
            "unit": "ns/op",
            "extra": "1996192 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1996192 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1996192 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 598.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2004616 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 598.8,
            "unit": "ns/op",
            "extra": "2004616 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2004616 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2004616 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 596.2,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2000500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 596.2,
            "unit": "ns/op",
            "extra": "2000500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2000500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2000500 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 595.7,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "1999117 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 595.7,
            "unit": "ns/op",
            "extra": "1999117 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "1999117 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "1999117 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo",
            "value": 589.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "2005116 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - ns/op",
            "value": 589.8,
            "unit": "ns/op",
            "extra": "2005116 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "2005116 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatTimeAgo - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "2005116 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.623,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.623,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6248,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6248,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6235,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6235,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6244,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6244,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6226,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6226,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6229,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6229,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6223,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6223,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.623,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.623,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6234,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6234,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations",
            "value": 0.6228,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - ns/op",
            "value": 0.6228,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkMockClientOperations - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 66129,
            "unit": "ns/op\t   67438 B/op\t     675 allocs/op",
            "extra": "17810 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 66129,
            "unit": "ns/op",
            "extra": "17810 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67438,
            "unit": "B/op",
            "extra": "17810 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17810 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 66641,
            "unit": "ns/op\t   67431 B/op\t     675 allocs/op",
            "extra": "17910 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 66641,
            "unit": "ns/op",
            "extra": "17910 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67431,
            "unit": "B/op",
            "extra": "17910 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17910 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 68023,
            "unit": "ns/op\t   67439 B/op\t     675 allocs/op",
            "extra": "18339 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 68023,
            "unit": "ns/op",
            "extra": "18339 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67439,
            "unit": "B/op",
            "extra": "18339 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18339 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 67697,
            "unit": "ns/op\t   67428 B/op\t     675 allocs/op",
            "extra": "17708 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 67697,
            "unit": "ns/op",
            "extra": "17708 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67428,
            "unit": "B/op",
            "extra": "17708 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17708 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 66735,
            "unit": "ns/op\t   67433 B/op\t     675 allocs/op",
            "extra": "17949 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 66735,
            "unit": "ns/op",
            "extra": "17949 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67433,
            "unit": "B/op",
            "extra": "17949 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17949 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 67400,
            "unit": "ns/op\t   67457 B/op\t     675 allocs/op",
            "extra": "17745 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 67400,
            "unit": "ns/op",
            "extra": "17745 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67457,
            "unit": "B/op",
            "extra": "17745 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17745 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 66849,
            "unit": "ns/op\t   67426 B/op\t     675 allocs/op",
            "extra": "17756 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 66849,
            "unit": "ns/op",
            "extra": "17756 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67426,
            "unit": "B/op",
            "extra": "17756 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "17756 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 66518,
            "unit": "ns/op\t   67426 B/op\t     675 allocs/op",
            "extra": "18055 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 66518,
            "unit": "ns/op",
            "extra": "18055 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67426,
            "unit": "B/op",
            "extra": "18055 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18055 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 66498,
            "unit": "ns/op\t   67424 B/op\t     675 allocs/op",
            "extra": "18044 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 66498,
            "unit": "ns/op",
            "extra": "18044 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67424,
            "unit": "B/op",
            "extra": "18044 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18044 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions",
            "value": 65685,
            "unit": "ns/op\t   67426 B/op\t     675 allocs/op",
            "extra": "18079 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - ns/op",
            "value": 65685,
            "unit": "ns/op",
            "extra": "18079 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - B/op",
            "value": 67426,
            "unit": "B/op",
            "extra": "18079 times\n4 procs"
          },
          {
            "name": "BenchmarkExpandSuggestions - allocs/op",
            "value": 675,
            "unit": "allocs/op",
            "extra": "18079 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3493843,
            "unit": "ns/op\t     624 B/op\t       5 allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3493843,
            "unit": "ns/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 624,
            "unit": "B/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3495909,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3495909,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3480181,
            "unit": "ns/op\t     513 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3480181,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 513,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3480097,
            "unit": "ns/op\t     515 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3480097,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 515,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3484340,
            "unit": "ns/op\t     386 B/op\t       5 allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3484340,
            "unit": "ns/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 386,
            "unit": "B/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3481881,
            "unit": "ns/op\t     518 B/op\t       5 allocs/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3481881,
            "unit": "ns/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 518,
            "unit": "B/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3489623,
            "unit": "ns/op\t     534 B/op\t       5 allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3489623,
            "unit": "ns/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 534,
            "unit": "B/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3480677,
            "unit": "ns/op\t     407 B/op\t       5 allocs/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3480677,
            "unit": "ns/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 407,
            "unit": "B/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3478371,
            "unit": "ns/op\t     426 B/op\t       5 allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3478371,
            "unit": "ns/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 426,
            "unit": "B/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "340 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody",
            "value": 3484590,
            "unit": "ns/op\t     404 B/op\t       5 allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - ns/op",
            "value": 3484590,
            "unit": "ns/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - B/op",
            "value": 404,
            "unit": "B/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateCommentBody - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "344 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 403054,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2814 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 403054,
            "unit": "ns/op",
            "extra": "2814 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2814 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2814 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 399762,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 399762,
            "unit": "ns/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2636 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 396207,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2733 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 396207,
            "unit": "ns/op",
            "extra": "2733 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2733 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2733 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395529,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2884 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395529,
            "unit": "ns/op",
            "extra": "2884 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2884 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2884 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395683,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2788 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395683,
            "unit": "ns/op",
            "extra": "2788 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2788 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2788 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 396211,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2659 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 396211,
            "unit": "ns/op",
            "extra": "2659 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2659 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2659 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395396,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395396,
            "unit": "ns/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2804 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 397083,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2798 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 397083,
            "unit": "ns/op",
            "extra": "2798 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2798 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2798 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395788,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2779 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395788,
            "unit": "ns/op",
            "extra": "2779 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2779 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2779 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig",
            "value": 395575,
            "unit": "ns/op\t  145488 B/op\t    2559 allocs/op",
            "extra": "2917 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - ns/op",
            "value": 395575,
            "unit": "ns/op",
            "extra": "2917 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - B/op",
            "value": 145488,
            "unit": "B/op",
            "extra": "2917 times\n4 procs"
          },
          {
            "name": "BenchmarkParseBatchConfig - allocs/op",
            "value": 2559,
            "unit": "allocs/op",
            "extra": "2917 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15112,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "75961 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15112,
            "unit": "ns/op",
            "extra": "75961 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "75961 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "75961 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15051,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "79255 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15051,
            "unit": "ns/op",
            "extra": "79255 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "79255 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "79255 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15102,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "77430 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15102,
            "unit": "ns/op",
            "extra": "77430 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "77430 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "77430 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15049,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "77649 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15049,
            "unit": "ns/op",
            "extra": "77649 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "77649 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "77649 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15006,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "75951 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15006,
            "unit": "ns/op",
            "extra": "75951 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "75951 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "75951 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15327,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76491 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15327,
            "unit": "ns/op",
            "extra": "76491 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76491 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76491 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15256,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76816 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15256,
            "unit": "ns/op",
            "extra": "76816 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76816 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76816 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 14992,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76347 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 14992,
            "unit": "ns/op",
            "extra": "76347 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76347 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76347 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15203,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "76458 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15203,
            "unit": "ns/op",
            "extra": "76458 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "76458 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "76458 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError",
            "value": 15192,
            "unit": "ns/op\t    5736 B/op\t     107 allocs/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - ns/op",
            "value": 15192,
            "unit": "ns/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - B/op",
            "value": 5736,
            "unit": "B/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkFormatActionableError - allocs/op",
            "value": 107,
            "unit": "allocs/op",
            "extra": "77170 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1382,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "902001 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1382,
            "unit": "ns/op",
            "extra": "902001 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "902001 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "902001 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1378,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "907322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1378,
            "unit": "ns/op",
            "extra": "907322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "907322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "907322 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1370,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "912559 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1370,
            "unit": "ns/op",
            "extra": "912559 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "912559 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "912559 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1373,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "937708 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1373,
            "unit": "ns/op",
            "extra": "937708 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "937708 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "937708 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1370,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "933367 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1370,
            "unit": "ns/op",
            "extra": "933367 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "933367 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "933367 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1366,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "939121 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1366,
            "unit": "ns/op",
            "extra": "939121 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "939121 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "939121 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1368,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "897993 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1368,
            "unit": "ns/op",
            "extra": "897993 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "897993 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "897993 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1366,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "952942 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1366,
            "unit": "ns/op",
            "extra": "952942 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "952942 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "952942 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1368,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "939178 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1368,
            "unit": "ns/op",
            "extra": "939178 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "939178 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "939178 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt",
            "value": 1369,
            "unit": "ns/op\t     704 B/op\t      29 allocs/op",
            "extra": "953481 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - ns/op",
            "value": 1369,
            "unit": "ns/op",
            "extra": "953481 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - B/op",
            "value": 704,
            "unit": "B/op",
            "extra": "953481 times\n4 procs"
          },
          {
            "name": "BenchmarkParsePositiveInt - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "953481 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 465.6,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2571867 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 465.6,
            "unit": "ns/op",
            "extra": "2571867 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2571867 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2571867 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 469.5,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2468774 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 469.5,
            "unit": "ns/op",
            "extra": "2468774 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2468774 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2468774 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 464.6,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2546236 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 464.6,
            "unit": "ns/op",
            "extra": "2546236 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2546236 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2546236 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 464.2,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2534811 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 464.2,
            "unit": "ns/op",
            "extra": "2534811 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2534811 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2534811 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 466.5,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2526234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 466.5,
            "unit": "ns/op",
            "extra": "2526234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2526234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2526234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 461.5,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2542234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 461.5,
            "unit": "ns/op",
            "extra": "2542234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2542234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2542234 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 464.8,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2555898 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 464.8,
            "unit": "ns/op",
            "extra": "2555898 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2555898 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2555898 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 462.6,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2539532 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 462.6,
            "unit": "ns/op",
            "extra": "2539532 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2539532 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2539532 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 463.5,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2536276 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 463.5,
            "unit": "ns/op",
            "extra": "2536276 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2536276 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2536276 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess",
            "value": 460.9,
            "unit": "ns/op\t     384 B/op\t       8 allocs/op",
            "extra": "2575356 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - ns/op",
            "value": 460.9,
            "unit": "ns/op",
            "extra": "2575356 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - B/op",
            "value": 384,
            "unit": "B/op",
            "extra": "2575356 times\n4 procs"
          },
          {
            "name": "BenchmarkColorizeSuccess - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2575356 times\n4 procs"
          }
        ]
      }
    ]
  }
}