# GitHub API Feature Request

**Submit to:** https://github.com/orgs/community/discussions/categories/api-and-webhooks

---

## Title
[API Feature Request] Allow adding comments to existing pending PR reviews (enable incremental review workflows)

## Body

### Summary
GitHub's web interface allows users to queue up review comments before submitting, but the API forces all comments to be created in a single call. This prevents building interactive review tools that match the natural workflow users expect, creating a significant capability gap between web and programmatic access.

### Problem Statement
When building GitHub CLI tools that facilitate code review workflows, developers cannot:
- Add additional comments to a pending review after it's been created
- Build interactive review tools that let users add comments iteratively to the same review
- Create review workflows where comments are discovered and added over time to a single review

#### Current API Behavior
```bash
# This works:
POST /repos/owner/repo/pulls/123/reviews
{
  "body": "Overall looks good",
  "event": "COMMENT",
  "comments": [
    {"path": "file.js", "line": 42, "body": "Comment 1"},
    {"path": "file.js", "line": 50, "body": "Comment 2"}
  ]
}

# This fails with 422 error:
POST /repos/owner/repo/pulls/123/reviews/{review_id}/comments
{"path": "file.js", "line": 60, "body": "Comment 3"}
```

#### Exact Error Response
```json
{
  "message": "Pull request review cannot be created",
  "errors": [
    {
      "resource": "PullRequestReview",
      "code": "custom",
      "message": "A review cannot be created because a pending review already exists"
    }
  ],
  "documentation_url": "https://docs.github.com/rest/pulls/reviews#create-a-review-for-a-pull-request"
}
```

### Technical Details
- **Core Issue**: No API endpoint exists to add comments to an existing pending review
- **Current Limitation**: Once a pending review is created, additional comments cannot be added to it via API
- **Affects Both APIs**: Neither REST nor GraphQL provide this capability
- **Workaround Required**: All review comments must be included in the initial review creation call

### API-Web Interface Parity Issue

**GitHub Web Interface Workflow:**
```
1. User writes a comment on a line
2. User clicks "Add review comment" (instead of "Add single comment")
3. ✅ Pending review is created
4. User continues browsing, finds more issues
5. User adds more comments, each one goes to the pending review
6. User clicks "Submit review" when ready
7. All comments are submitted together as a cohesive review
```

**GitHub API Workflow:**
```
1. Create review with ALL comments at once:
   POST /repos/owner/repo/pulls/123/reviews
   {
     "body": "Review summary",
     "event": "COMMENT",
     "comments": [
       {"path": "file.js", "line": 42, "body": "Comment 1"},
       {"path": "file.js", "line": 50, "body": "Comment 2"},
       {"path": "file.js", "line": 60, "body": "Comment 3"}
     ]
   }
2. ❌ Review is immediately submitted - no pending state
3. ❌ Cannot add more comments to this review later
```

**The Gap:** The API has no equivalent to the web interface's "pending review" workflow where you can queue up comments before submitting.

### Real-World Impact: gh-comment CLI Tool

I've built [`gh-comment`](https://github.com/silouanwright/gh-comment), a GitHub CLI extension specifically for AI-assisted code review workflows. This limitation directly impacts:

**Current User Experience:**
- ❌ Users must plan entire review upfront 
- ❌ Cannot add comments as they discover issues
- ❌ Forces choice between immediate individual comments OR batched reviews
- ❌ Breaks natural review flow where insights emerge incrementally

**Desired User Experience:**
- ✅ Start review, add initial comments
- ✅ Continue analysis, add more comments to same review
- ✅ Submit comprehensive review when complete
- ✅ Matches GitHub web interface behavior via API

### Compelling Use Case: AI Code Review Workflow

**The Problem in Action:**
An AI code reviewer analyzes a PR and immediately finds 3 critical issues, so it starts a review. During deeper analysis, it discovers 2 more subtle problems and wants to add them to the same review for a comprehensive evaluation. 

**Current Reality:** The AI must either:
- ❌ Submit an incomplete review with only the first 3 issues
- ❌ Start over and re-analyze everything to create one big review
- ❌ Post the new issues as separate individual comments (breaking review cohesion)

**With This Feature:** The AI could naturally add the 2 additional comments to the existing pending review, then submit a complete, professional review - exactly like a human would do in the web interface.

### Community Evidence

**Tools Affected by This Limitation:**
- [`gh-comment`](https://github.com/silouanwright/gh-comment) - CLI for strategic PR commenting (workaround implemented)
- [PyGithub Issue #3038](https://github.com/PyGithub/PyGithub/issues/3038) - Python library users affected
- Multiple Stack Overflow questions about batching review comments

### Developer Ecosystem Impact

**Affected Use Cases:**
- 🤖 AI-assisted code review tools
- 🔄 Interactive CLI review workflows  
- 📱 Mobile apps for code review
- 🔧 IDE integrations for PR review
- 📊 Review analytics tools that need incremental data collection

**Business Impact:**
- Limits innovation in developer tooling
- Forces suboptimal UX in review tools
- Creates barrier to building GitHub-integrated products

### Proposed Solution
Add API support for the "pending review" workflow that already exists in the GitHub web interface.

#### Requested API Enhancement
Enable this workflow via API (matching the web interface):
```
1. Create pending review (without submitting):
   POST /repos/owner/repo/pulls/123/reviews
   {
     "body": "Starting my review",
     "event": "PENDING",  # New event type
     "comments": [
       {"path": "file.js", "line": 42, "body": "First comment"}
     ]
   }
   
2. Add more comments to the pending review:
   POST /repos/owner/repo/pulls/123/reviews/{review_id}/comments
   {
     "path": "file.js", 
     "line": 50, 
     "body": "Second comment found later"
   }
   
3. Submit the review when ready:
   PATCH /repos/owner/repo/pulls/123/reviews/{review_id}
   {
     "event": "COMMENT"  # Or APPROVE/REQUEST_CHANGES
   }
```

This would replicate the exact workflow available in the GitHub web interface.

### Context
- **Tool Type**: GitHub CLI application (`gh-comment`)
- **APIs Used**: Both REST API and GraphQL API
- **Language**: Go
- **Current Workaround**: Create all comments in single API call (limiting interactivity)
- **Repository**: https://github.com/silouanwright/gh-comment

### References
- [PyGithub Issue #3038](https://github.com/PyGithub/PyGithub/issues/3038): "issues with pull request review, adding many comments always fails"
- [Stack Overflow #71421045](https://stackoverflow.com/questions/71421045/): "How to add comments to pending GitHub review via API"
- [GitHub Community Discussion #24854](https://github.com/orgs/community/discussions/24854): "Cannot add comments to existing pending review"
- [GitHub API Documentation](https://docs.github.com/en/rest/pulls/comments): Current API limitations

---

**Search Confirmation**: ✅ I have searched for existing discussions and found no duplicate requests for this specific API enhancement.

**Impact**: This change would enable a new generation of interactive GitHub review tools and bring API functionality to parity with the web interface experience.

---

**🚀 If you've experienced this API limitation while building GitHub tools, please upvote this discussion and share your specific use case in the comments. The more evidence we can provide of developer impact, the stronger the case for GitHub to prioritize this enhancement.**
