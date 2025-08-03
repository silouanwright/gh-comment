---
name: address-review-workflow
title: Address Review Comments Workflow
category: workflows
estimated_time: 15-25 minutes
tags: [response, feedback, resolution, systematic, unified]
examples:
  - "Use for: Systematically addressing all reviewer feedback"
  - "Great for: Complex PRs with multiple reviewers and extensive feedback"
---

# 🔧 Address Review Comments Workflow

You are helping the PR author systematically address ALL review feedback using `gh-comment`'s unified comment system. This workflow handles **both general PR discussion and line-specific code review comments**, focusing on **responding to feedback**, **fixing issues**, and **resolving conversations** professionally and efficiently.

## 🎯 Purpose

When reviewers leave comments on your PR (both types), you need to:
- **Acknowledge** their feedback appropriately (general discussion + code-specific)
- **Fix** the issues they've identified in the code
- **Respond** with context or questions when needed
- **Resolve** conversations once addressed (review comments only)
- **Maintain** positive team relationships across all communication

## 🛠️ Command Reference (Important!)

**Before you start**, understand which commands to use:

```bash
# 💬 GENERAL PR COMMENTS (use native GitHub CLI)
gh pr comment <PR_NUMBER> --body "Your message"     # Create new general comment
gh comment reply <COMMENT_ID> "Reply" --type issue   # Reply to existing general comment

# 📋 LINE-SPECIFIC COMMENTS (use gh-comment extension)
gh comment add <PR> <file> <line> "Comment"          # Create new line-specific comment
gh comment reply <COMMENT_ID> "Reply" --type review  # Reply to existing line-specific comment
```

**Key Distinction:**
- **General discussion** = `gh pr comment` (native CLI) or `gh comment reply --type issue`
- **Code-specific feedback** = `gh comment add` (extension) or `gh comment reply --type review`

## 📋 Address Review Process

### Step 1: List and Analyze All Comments (Unified System)

Start by getting a complete view of ALL feedback - both general discussion and line-specific:

```bash
gh comment list <PR_NUMBER>
```

**The output shows two distinct sections:**
- **💬 General PR Comments**: Overall feedback, questions, LGTM, general discussion
- **📋 Review Comments**: Line-specific code review comments with file/line context

**Categorize ALL comments by action needed:**
- 🐛 **Bugs/Issues** - Must fix before merge (usually line-specific)
- 💡 **Suggestions** - Consider implementing (both types)
- ❓ **Questions** - Need to respond with clarification (both types)
- ✨ **Praise** - Acknowledge with reactions (both types)
- 🤔 **Discussions** - Engage in conversation (usually general)

### Step 2: React Appropriately to Comments (Both Types)

**Use reactions to acknowledge feedback quickly (works with both comment types):**

```bash
# Great suggestions or catches (any comment type)
gh comment reply <COMMENT_ID> --reaction "heart"

# Good suggestions you'll implement (any comment type)
gh comment reply <COMMENT_ID> --reaction "+1"

# Funny or clever observations (any comment type)
gh comment reply <COMMENT_ID> --reaction "laugh"

# Mind-blowing insights (any comment type)
gh comment reply <COMMENT_ID> --reaction "rocket"

# Surprising discoveries (any comment type)
gh comment reply <COMMENT_ID> --reaction "eyes"
```

**Note**: Reactions work identically for both general PR comments and line-specific review comments.

### Step 3: Respond with Context and Fixes (Type-Specific)

**For line-specific review comments (default behavior):**

```bash
# Reply to code review comment (creates threaded reply)
gh comment reply <COMMENT_ID> "Thanks for catching this! Fixed in [commit hash]. The issue was [explanation]." --resolve

# Or explicitly specify review type
gh comment reply <COMMENT_ID> "Fixed the null check as suggested!" --type review --resolve
```

**For general PR discussion comments:**

```bash
# Reply to general discussion (creates new top-level comment)
gh comment reply <COMMENT_ID> "Thanks for the overall feedback! I've addressed all the points below." --type issue

# Note: --resolve doesn't work with issue comments (no conversation threading)
```

### Step 4: Batch Address Multiple Issues

**IMPORTANT:** After fixing multiple issues, summarize your changes with a **general PR comment**:

```bash
# Use native GitHub CLI for general PR comments (not gh-comment extension)
gh pr comment <PR_NUMBER> --body "## 🔧 Review Comments Addressed

Thanks everyone for the thorough review! I've addressed all the feedback:

### ✅ Issues Fixed:
- **@reviewer1's comment on line 45**: Fixed null pointer issue with proper validation
- **@reviewer2's suggestion on line 78**: Refactored to use Map for O(n) performance  
- **@reviewer3's question about error handling**: Added comprehensive error messages

### 📝 Changes Made:
- Commit abc123: Fix validation logic
- Commit def456: Performance optimization
- Commit ghi789: Improve error handling

### 🤔 Still Discussing:
- **Architecture question from @reviewer4**: Let's discuss the trade-offs in our next sync

Ready for re-review! 🚀"
```

**Note**: `gh comment add` is for line-specific comments only. For general PR comments, use `gh pr comment`.

## 💬 Response Best Practices & Decorum

### 🎯 Professional Response Patterns

#### Acknowledging Good Catches
```
nice catch! fixed in abc123. totally missed that edge case 🤦‍♂️

you're absolutely right about the null handling - updated the logic to check for that

smart suggestion! hadn't thought about the performance implications. implemented your approach and it's much cleaner
```

#### Asking for Clarification
```
could you help me understand what you mean by the "async handling" part? 

quick question - are you thinking we should debounce the input or throttle it?

want to make sure I get this right - should I move the validation to the middleware layer?
```

#### Explaining Your Reasoning
```
went with this approach because it matches the pattern we use in the user service. happy to change if you see issues with the consistency angle

tested both approaches and this one was 2x faster with large datasets. the memory trade-off seemed worth it but let me know if you disagree
```

#### Disagreeing Respectfully
```
i see your point, though i'm leaning toward keeping the current approach because it's more explicit about error handling. what do you think?

interesting idea! my concern is that it might make debugging harder when things go wrong. could we maybe try a hybrid approach?

let's chat about this - feels like there might be some trade-offs we should think through together
```

### 🚫 What NOT to Do

- ❌ "That's wrong" → ✅ "I see it differently because..."
- ❌ "This is fine" → ✅ "Good point! Here's my reasoning..."
- ❌ Ignoring comments → ✅ At least react or acknowledge
- ❌ Being defensive → ✅ Being curious and collaborative
- ❌ "Will fix later" → ✅ "Fixed in commit abc123" or "Created issue #123"

## 🔄 Resolve Conversations Systematically (Review Comments Only)

**IMPORTANT**: Conversation resolution only works with **line-specific review comments**. General PR discussion comments don't have conversation threading.

### When to Resolve (Review Comments):
- ✅ **Issue is fixed** and you've explained the fix
- ✅ **Question is answered** completely
- ✅ **Suggestion is implemented** or politely declined with reasoning
- ✅ **Discussion reached conclusion** or agreement to continue elsewhere

### How to Resolve (Review Comments Only):
```bash
# Fix and resolve in one action (review comments)
gh comment reply <COMMENT_ID> "Fixed in commit abc123! The validation now handles edge cases properly." --resolve

# Or explicitly specify type and resolve
gh comment reply <COMMENT_ID> "Addressed this concern" --type review --resolve
```

### For General PR Comments:
```bash
# General discussion - just reply (no resolve option)
gh comment reply <COMMENT_ID> "Thanks for the feedback! All points addressed." --type issue
```

## 📊 Track Your Progress

### Create a Response Checklist (Both Comment Types):
```
## 📋 Review Response Checklist

### 💬 General PR Discussion:
- [ ] Comment #301: Overall LGTM → React with heart
- [ ] Comment #302: Architecture question → Explain reasoning (--type issue)
- [ ] Comment #303: Thanks for contribution → Reply with gratitude (--type issue)

### 📋 Line-Specific Review Comments:
- [ ] Comment #123: Null pointer bug → Fix validation (--type review --resolve)
- [ ] Comment #456: Security issue → Add input sanitization (--type review --resolve)

### 💡 Suggestions (Both Types):
- [ ] Comment #789: Performance optimization → Benchmark and implement
- [ ] Comment #101: Code organization → Refactor if time permits

### ❓ Questions (Need to Respond):
- [ ] Comment #112: General question → Reply (--type issue)
- [ ] Comment #134: Line-specific question → Reply (--type review)

### ✨ Praise (Acknowledge - Both Types):
- [ ] Comment #145: Good implementation → React with heart
- [ ] Comment #167: Clever solution → React with rocket
```

## 🚀 Ready to Address Reviews? (Unified System)

1. **List ALL comments on your PR** (both general discussion and line-specific)
2. **Categorize by comment type and priority** (issue vs review)
3. **React to acknowledgeable comments** (works with both types)
4. **Reply with appropriate type flag** (--type issue for general, --type review for line-specific)
5. **Fix issues and respond with context** (use --resolve for review comments only)
6. **Resolve review conversations systematically** (only works with review comments)
7. **Summarize all changes made** (address both types of feedback)

**Remember**: 
- General PR comments (💬) use `--type issue` and create new top-level comments
- Line-specific comments (📋) use `--type review` (default) and support `--resolve`
- Use `[SUGGEST: code]` syntax in replies for actionable code changes!

What PR has review comments that need addressing?
