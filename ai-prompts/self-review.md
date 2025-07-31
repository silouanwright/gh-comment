# 📝 Self-Review Documentation Workflow

You are helping the PR author document and explain their own changes using `gh-comment`. This workflow focuses on **explaining the "why"** behind code decisions rather than critiquing the implementation.

## 🎯 Purpose

When you create a PR, reviewers often need context to understand:
- **Why** you chose this approach
- **What** alternatives you considered  
- **How** this fits into the bigger picture
- **Where** the tricky parts are
- **When** certain decisions were made

This workflow helps you proactively document your thought process, making reviews faster and more productive.

## 📋 Self-Documentation Process

### Step 1: Identify Key Areas to Explain

Look for code that reviewers might question:
- **🤔 Complex Logic** - Algorithms, business rules, edge cases
- **🏗️ Architecture Decisions** - Why you structured code this way
- **🔄 Refactoring Choices** - Why you changed existing patterns
- **⚡ Performance Optimizations** - Trade-offs you made
- **🔧 Technical Debt** - Temporary solutions and future plans
- **📚 Domain Knowledge** - Business context reviewers might not have

### Step 2: Explain Your Reasoning

For each area, document:
- **Context**: What problem were you solving?
- **Decision**: What approach did you choose?
- **Alternatives**: What other options did you consider?
- **Trade-offs**: Why this choice was best
- **Future**: Any follow-up work planned

### Step 3: Add Explanatory Comments

**IMPORTANT:** Prepare ALL your explanatory comments before submission:

```bash
gh comment add-review <PR_NUMBER> \
  --event "COMMENT" \
  --body "## 📝 Self-Review: Context & Decisions

👋 I've added some context comments to help with the review process.

### 🎯 Key Changes Summary:
- [CHANGE 1]: [Brief explanation]
- [CHANGE 2]: [Brief explanation]

### 🤔 Areas for Extra Review:
- [AREA 1]: [Why it needs attention]
- [AREA 2]: [Why it needs attention]

Feel free to ask questions about anything that's unclear!" \
  --comment "src/utils/parser.js:45:💡 **Context: Complex parsing logic** - This handles edge cases where user input contains nested brackets. I considered using a regex but chose manual parsing for better error messages and debugging." \
  --comment "src/api/client.js:78:🏗️ **Architecture: New caching layer** - Added Redis caching here because this endpoint gets hit 1000+ times/day. The 5-minute TTL balances freshness vs performance."
```

## 💬 Self-Review Comment Templates

### 🗣️ Casual Self-Review Tone

**Keep it conversational and explanatory:**

**✅ Good: Casual & Informative**
```
had to inline this mock because jest was throwing scoping errors when we tried to reference setupDatadogMock() from the external function. basically jest mock factories can't reference variables from outer scope unless they're prefixed with 'mock' - learned this the hard way when all the tests started failing with reference errors.
```

**❌ Bad: Formal & Bot-like**
```
🔧 **Technical Implementation**: This inline mock definition resolves Jest variable scoping violations as documented in our migration guide section 5.
```

### Comment Templates

#### Complex Logic Explanation
```
this parsing logic handles the tricky case where users can nest brackets in their input. tried regex first but debugging was a nightmare when edge cases failed, so went with manual parsing for clearer error messages.
```

#### Architecture Decision
```
went with redis caching here because this endpoint gets hammered (1000+ requests/day) and the data only changes once per hour. tried in-memory caching first but that doesn't work across multiple server instances. 5-minute TTL keeps things fresh enough while cutting response times from 800ms to 50ms.
```

#### Performance Optimization
```
replaced the nested loops with a Map lookup - went from O(n²) to O(n) which makes a huge difference when processing large datasets. tested with 10k records and saw 15x speedup. the trade-off is slightly more memory usage but totally worth it for the performance gain.
```

#### Technical Debt Note
```
this is a bit of a hack for now - ideally we'd have proper validation middleware but that's a bigger refactor. added a TODO to clean this up in the next sprint. it works fine for current use cases but might need attention if we add more complex validation rules.
```

#### Domain Knowledge Context
```
the 72-hour rule here comes from our compliance team - users have 3 business days to dispute charges. the weird timezone handling is because we need to use the user's local timezone, not server time. sarah from legal can explain more if needed - this tripped up the last dev who worked on billing.
```

## 🎯 Self-Review Benefits

### For Reviewers:
- ✅ Faster review process
- ✅ Better understanding of decisions
- ✅ More focused feedback
- ✅ Reduced back-and-forth questions

### For You:
- ✅ Forces you to think through decisions
- ✅ Documents your reasoning for future reference
- ✅ Builds trust with your team
- ✅ Reduces review cycle time

## 🚀 Ready for Self-Review?

1. **Review your own PR changes**
2. **Identify areas that need explanation**
3. **Prepare context comments for each area**
4. **Submit all explanatory comments at once**

What PR would you like to self-document?
