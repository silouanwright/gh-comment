# Context Continuation Guide
**Date**: August 2, 2025  
**Status**: Mid-implementation of intelligent API endpoint handling

## 🎯 Current State Summary

### ✅ **Completed Tasks**
1. **Fixed nil pointer crashes** - All commands now properly initialize GitHub clients
2. **Added intelligent error analysis** - Enhanced error messages with auto-correction suggestions
3. **Fixed some API endpoints** - Reactions now work for review comments
4. **Integration test success** - All core functionality working (add, list, suggestions, reactions)

### 🚧 **Current Issue Being Addressed**
**Problem**: Commands use invisible API fallbacks (trying multiple endpoints silently)
- AddReaction: tries pulls/comments then issues/comments 
- EditComment: tries pulls/comments then issues/comments
- RemoveReaction: tries pulls/comments then issues/comments

**User Feedback**: "The sort of invisible fallbacks of like trying other APIs rather than simply just saying in the error messaging, hey, try this way. Do you think that's the right approach?"

**Answer**: No, invisible fallbacks are problematic. Better approach is intelligent detection + clear messaging.

## 🎯 **Next Steps to Complete**

### **1. Replace Invisible Fallbacks with Smart Detection**

**Files to modify**:
- `/internal/github/real_client.go` - Remove fallback logic in:
  - `AddReaction()` (lines ~290-304)
  - `EditComment()` (lines ~366-382) 
  - `RemoveReaction()` (lines ~325-342)

**Approach**:
```go
// Instead of trying both endpoints, detect comment type first
func (c *RealClient) AddReaction(owner, repo string, commentID int, reaction string) error {
    // 1. Detect if comment is issue or review type
    commentInfo, err := c.DetectCommentType(owner, repo, commentID, prNumber)
    if err != nil {
        return CreateSmartError(c, "add_reaction", "reply", commentID, prNumber, err)
    }
    
    // 2. Use correct endpoint based on detection
    var endpoint string
    if commentInfo.Type == "review" {
        endpoint = fmt.Sprintf("repos/%s/%s/pulls/comments/%d/reactions", owner, repo, commentID)
    } else {
        endpoint = fmt.Sprintf("repos/%s/%s/issues/comments/%d/reactions", owner, repo, commentID)
    }
    
    // 3. Single API call with helpful error if it fails
    err = c.restClient.Post(endpoint, body, nil)
    if err != nil {
        return CreateSmartError(c, "add_reaction", "reply", commentID, prNumber, err)
    }
}
```

### **2. Fix PR Number Detection Issue**

**Problem**: `DetectCommentType()` needs PR number but commands don't always have it.

**Solutions**:
1. Add PR detection to commands that need it
2. Use the existing `getPRContext()` helper from commands
3. Update command signatures to pass PR number to client methods

**Files to check**:
- `cmd/reply.go` - Add PR context to reaction calls
- `cmd/edit.go` - Add PR context to edit calls

### **3. Enhance Error Messages Integration**

**Current**: `AnalyzeAndEnhanceError()` in `/internal/github/error_helper.go`
**Missing**: Integration with the new `comment_detector.go` functionality

**Tasks**:
1. Update `AnalyzeAndEnhanceError()` to use comment type detection
2. Provide specific command suggestions based on detected comment type
3. Include file paths and line numbers in error messages for review comments

### **4. Complete Review Command Fix**

**Current Issue**: Review command fails with GraphQL `commitId` error
**Root Cause**: Missing commit ID in review comment creation

**Fix needed** in `CreateReview()`:
```go
// Get latest commit SHA for the PR
prDetails, err := c.GetPRDetails(owner, repo, pr)
latestCommitSHA := prDetails["head"].(map[string]interface{})["sha"].(string)

// Add commit_id to each comment
for i := range review.Comments {
    review.Comments[i].CommitID = latestCommitSHA
}
```

### **5. Add Comprehensive Tests**

**File**: `/internal/github/endpoint_selection_test.go` (already created)
**Missing**:
1. Tests for smart detection logic
2. Tests for enhanced error messages
3. Integration tests that verify single API calls (no fallbacks)

**Test Cases Needed**:
- Comment type detection accuracy
- Error message helpfulness
- Command suggestions correctness
- Performance (single API call vs multiple)

## 📋 **Implementation Checklist**

### **Phase 1: Remove Fallbacks** 
- [ ] Remove fallback logic from `AddReaction()`
- [ ] Remove fallback logic from `EditComment()`  
- [ ] Remove fallback logic from `RemoveReaction()`
- [ ] Test that single endpoint calls work correctly

### **Phase 2: Integrate Smart Detection**
- [ ] Fix PR number passing from commands to client methods
- [ ] Integrate `DetectCommentType()` into reaction methods
- [ ] Integrate `DetectCommentType()` into edit methods
- [ ] Update error handling to use smart detection

### **Phase 3: Fix Review Command**
- [ ] Fix GraphQL commitId issue in `CreateReview()`
- [ ] Test review creation with multiple comments
- [ ] Verify review command works end-to-end

### **Phase 4: Testing & Validation**
- [ ] Add unit tests for smart detection
- [ ] Add integration tests for single API calls
- [ ] Test error message helpfulness
- [ ] Performance testing (fewer API calls)

## 🚨 **Key Files Modified Recently**

1. **`/internal/github/real_client.go`**
   - Fixed nil pointer issues (✅)
   - Added reaction endpoint fallbacks (🚧 needs removal)
   - Added edit endpoint fallbacks (🚧 needs removal)

2. **`/internal/github/error_helper.go`** 
   - Added intelligent error analysis (✅)
   - Provides context-aware suggestions (✅)

3. **`/internal/github/comment_detector.go`**
   - Created smart comment type detection (✅)
   - Needs integration with actual commands (🚧)

4. **All `cmd/*.go` files**
   - Fixed nil pointer client initialization (✅)
   - Use `createGitHubClient()` helper properly (✅)

## 🎯 **Success Criteria**

When complete, the system should:
1. **Make single API calls** - No more invisible fallbacks
2. **Provide helpful errors** - Clear guidance on comment types and alternatives
3. **Auto-detect comment types** - Use `gh comment list` data to determine endpoints
4. **Work reliably** - All commands (reply, edit, review) work consistently
5. **Be fast** - Fewer API calls = better performance and rate limit usage

## 🔗 **Context Links**

- **Integration test PR**: https://github.com/silouanwright/gh-comment/pull/4
- **Test comment ID**: 2249338891 (review comment, works with reactions now)
- **Current status**: Basic functionality works, but needs architectural improvement

## 💡 **Design Philosophy**

The goal is **transparent, helpful tooling** rather than **magic that sometimes works**. Users should understand:
- What type of comment they're working with
- Why operations fail
- How to fix issues
- What alternatives exist

This aligns with your existing comprehensive help system and professional CLI approach.