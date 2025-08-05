# Human-Centered Code Review Guide
*Making code reviews feel like collaborative problem-solving, not performance evaluations*

## 🧠 The Core Problem

Code reviews trigger the same fight-or-flight response as being criticized in public. **Your goal**: Make authors feel like you're solving problems together, not pointing out their failures.

**The Golden Rule**: Separate the code from the coder
- ❌ "You always forget error handling"
- ✅ "What do you think about adding error handling for the database call?"

### The Three Psychological Triggers (Avoid These!)
1. **Truth Triggers**: "They got it wrong!" → Ask questions instead of making statements
2. **Relationship Triggers**: "Why are they saying this?" → Build trust through explanation
3. **Identity Triggers**: "I'm a bad developer" → Focus on code improvement, not person judgment

*[Learn more about psychological safety in code reviews →](https://agilesparks.com/build-psychological-safety-in-teams-through-code-reviews)*

## 💬 Communication That Works

### Frame Everything as Questions
**Instead of Commands → Use Curious Questions**
- ❌ "Use a Map here"
- ✅ "What do you think about using a Map here? It might help with lookup performance."

**Instead of "You" → Use "We" Language**
- ❌ "You need to refactor this function"
- ✅ "Could we consider breaking this function into smaller methods?"

### The OIR Framework for Tough Feedback
**Observation + Impact + Request** = Less defensiveness
- **Observation**: "I notice this function handles multiple responsibilities"
- **Impact**: "which makes it harder to test and maintain"
- **Request**: "Could we consider breaking it into separate functions?"

*[Deep dive into constructive feedback techniques →](https://www.michaelagreiler.com/respectful-constructive-code-review-feedback/)*

## 🎯 Make Suggestions Specific, Not Vague

**The #1 problem in code reviews**: Vague feedback that leaves authors guessing what you actually want.

### Use GitHub's Suggestion Feature for Concrete Changes
Instead of describing what to change, **show exactly what you want** using GitHub's built-in suggestion feature.

**How to use it:**
1. In the "Files changed" tab, click the + icon next to the line you want to change
2. Click the +/- suggestion button in the comment toolbar
3. Edit the code to show exactly what you want
4. Author can apply your change with one click

**Before** (Vague and frustrating):
> "🤔 What about using a Map here? Better performance."

**After** (Specific and actionable):
> "🤔 What about using a Map here? Better lookup performance and clearer intent:
>
> ```suggestion
> const users = new Map();
> users.set(id, userData);
> ```

**Why this transforms reviews:**
- **No guessing**: Author sees exactly what you want
- **One-click apply**: Reduces friction to implement feedback
- **Shared credit**: Both reviewer and author get commit credit
- **Less back-and-forth**: Fewer "is this what you meant?" rounds

*[Learn more about GitHub suggestions →](https://docs.github.com/articles/incorporating-feedback-in-your-pull-request)*

### [CREG: Code Review Emoji Guide](https://devblogs.microsoft.com/appcenter/how-the-visual-studio-mobile-center-team-does-code-review/)

CREG, pioneered by Microsoft, is a system for adding emojis before your comment to clarify intent. Microsoft found that emojis "help separate well-meant suggestions, simple questions, and must-have requests"—crucial because authors need to instantly know what's blocking vs. what's optional. The emojis also add "an additional human component to the conversation, so we don't forget there's a human on the other side of the screen."

**The Essential CREG Emojis:**
- 🔧 **Must fix** - Required changes
- ⛏️ **Nitpick** - Minor style issues, not blocking
- 😃 **Praise** - Highlight good work
- 🤔 **Question** - Need clarification or thinking out loud
- 📝 **Note** - Educational info, no action needed
- ♻️ **Refactor** - Structural improvements
- 📌 **Future** - Out of scope, note for later

**Examples:**
```
🔧 What do you think about adding null checking before the database call?
This might prevent runtime errors.

⛏️ I'm wondering if we could use camelCase here for consistency
with the rest of the codebase?

😃 Brilliant use of the factory pattern here!

🤔 What was the reasoning behind choosing async over sync here?

📝 I noticed this pattern is called dependency injection -
thought it might be helpful context!

♻️ Could we consider extracting this into a helper function?
It might improve readability.

📌 I'm wondering if we could add rate limiting here in a future iteration?
```

*[Full CREG emoji system →](https://github.com/erikthedeveloper/code-review-emoji-guide)*

## 📝 Handling Tricky Scenarios

**Security Issue Using OIR Framework**:
> "I notice this endpoint doesn't have authentication (Observation), which could allow unauthorized access to user data (Impact). What do you think about adding auth middleware before the handler? (Request)"

**Performance Concern with Collaborative Language**:
> "I'm wondering if this approach might struggle with large datasets since we're loading everything into memory. Could we consider pagination or streaming? Happy to pair on this if it gets complex!"

**Architecture Question Without Commands**:
> "I'm curious about the reasoning behind putting business logic in the controller. What do you think about moving it to a service layer? This might make testing easier, but I could be missing something about the architecture."

**Handling Disagreement Constructively**:
> "I see a different approach here than what we discussed earlier. Help me understand your thinking—what made this solution feel like the better path? I want to make sure I'm not missing something important."

## 🚫 Toxic Patterns That Kill Teams

### The "Death by a Thousand Round Trips"
Don't provide feedback incrementally across multiple rounds. Give comprehensive feedback upfront.

### The "Hostage Situation"
Never block PRs to force unrelated work:
- ❌ "This PR looks good, but first refactor that legacy module"
- ✅ "This looks great! I'm wondering if that legacy module could use similar improvements in a future iteration?"

### Communication Red Flags
- **Sarcasm**: "Did you even test this?"
- **Hostility**: "This is wrong"
- **Gatekeeping**: "That's not how we do things here"
- **Perfectionism**: Blocking adequate solutions for theoretical perfect ones

*[More code review anti-patterns to avoid →](https://blog.submain.com/toxic-code-review-culture/)*

## ⚡ Process Essentials

### Size Matters: The 400-Line Rule
**Analogy**: Code reviews are like proofreading essays. Beyond 400 lines, your brain starts missing important details.

- **Sweet spot**: Under 200 lines
- **Maximum**: 400 lines (defect detection drops significantly after this)

### When Large PRs Are Actually Better
**Exception**: Library migrations, breaking changes, and automated refactoring.

Sometimes splitting a large change creates **more** mental overhead than reviewing it all at once. This happens when:

**Library Migrations & Breaking Changes**
- Upgrading from React 16 → 18, or migrating from Bugsnag → DataDog
- The same pattern gets applied across dozens of files
- **Why keep it together**: Reviewers need to see the complete migration pattern to give meaningful feedback. If they suggest a different approach on the last small PR, you have to rework all the previous ones.

**Automated Refactoring**
- IDE-generated changes like "rename method across codebase"
- **Why it's different**: You're really reviewing [the tool](https://softwareengineering.stackexchange.com/questions/381343) and the decision to use it, not 1000+ individual line changes.

**Real example**: A [major library migration](https://bssw.io/blog_posts/pull-request-size-matters) touched 450K lines across 2,000+ files—splitting it would have been counterproductive.

**The trade-off**: Higher review overhead vs. avoiding fragmented context and rework cycles.

*[Research on optimal code review size →](https://smartbear.com/learn/code-review/best-practices-for-peer-code-review/)*

### Response Time Expectations
- **Initial response**: Within 4-6 hours during overlapping work hours
- **Full review**: Within 24 hours for business-critical changes

### Match Your Communication to Experience Level
**For Junior Developers**: More context, suggest learning resources, offer pairing
**For Senior Developers**: More concise, engage in technical debates, question decisions

## 🎯 The Conciseness Principle

**Think Twitter, Not Novel**: Every word adds value. Respect cognitive load.

**Before** (Cognitive overload):
> "I think you might want to consider refactoring this approach to use a Map data structure because it would provide better performance characteristics for lookups, and also it would be more semantically appropriate for this use case where we're doing key-value operations, and additionally it would make the code more maintainable in the long run since Maps have built-in iteration methods that are more efficient than what we're currently doing with arrays, plus the syntax is cleaner and more readable, and I believe most modern JavaScript engines optimize Map operations better than object property access in this scenario."

**After** (Concise and clear):
> "🤔 What about using a Map here? Better lookup performance and clearer intent for key-value operations."

*[Why cognitive load matters in code reviews →](https://link.springer.com/article/10.1007/s10664-022-10123-8)*

## 🚀 Quick Reference Cheat Sheet

### Before You Comment, Ask:
1. **Is it true?** (Facts vs. opinions)
2. **Is it necessary?** (Does it meaningfully improve the code?)
3. **Is it kind?** (Will this build up or tear down?)

### Essential Phrases:
- "What do you think about..."
- "Could we consider..."
- "I'm wondering if..."
- "What was the reasoning behind..."

### Red Flag Words to Avoid:
- "Obviously" • "Just" • "Simply" • "Clearly" • "You should" • "Wrong"

## 💭 Be Your Own First Reviewer

**Before hitting "Create Pull Request"** - review your own changes and add proactive comments.

This isn't about approving your own work; it's about **guiding reviewers through your thinking** and catching issues early. [Industry](https://medium.com/@sahilseth/pr-guidelines-how-to-author-and-review-pull-requests-d4f3450acec4) [experts](https://medium.com/google-developer-experts/how-to-pull-request-d75ac81449a5) [consistently](https://www.gustavwengel.dk/2025/02/19/pr-reviewer-practices.html) recommend this practice.

### Self-Review Your Changes
Look through your diff as if you're seeing it for the first time:
- Did you leave any debug code or TODO comments?
- Are there parts that might confuse reviewers?
- Do variable names make sense out of context?
- Are there non-obvious decisions that need explanation?

### Add Proactive Comments
Help reviewers focus on what matters by explaining:

**Complex Logic**:
> "This algorithm handles the edge case where users can have duplicate emails across different tenants. Had to use a compound key here instead of the simpler approach."

**Uncertain Decisions**:
> "I'm not sure about this variable name - does `processedUserData` make sense, or would you suggest something clearer?"

**Non-Obvious Approaches**:
> "Using Promise.allSettled instead of Promise.all here because we want to continue processing even if some user validations fail."

**Why this works**: Self-review ["makes you a better developer"](https://medium.com/google-developer-experts/how-to-pull-request-d75ac81449a5) and reviewers consistently report that ["author comments guiding the review"](https://www.gustavwengel.dk/2025/02/19/pr-reviewer-practices.html) are ["always helpful."](https://www.gustavwengel.dk/2025/02/19/pr-reviewer-practices.html)

## 🎯 The Ultimate Goal

**Remember**: Code reviews are collaboration, not evaluation. You're not a gatekeeper—you're a thinking partner helping create better code together.

Every review is an opportunity to share knowledge, build relationships, improve code quality, and create psychological safety.

---

## Want to Dive Deeper?

### 🧠 Psychology & Communication
**Understanding the human side of code reviews**
- [Building Psychological Safety in Code Reviews](https://agilesparks.com/build-psychological-safety-in-teams-through-code-reviews) - Why safety is the foundation of effective reviews
- [Respectful and Constructive Code Review Feedback](https://www.michaelagreiler.com/respectful-constructive-code-review-feedback/) - Practical communication techniques
- [Human Code Reviews (Part One)](https://mtlynch.io/human-code-reviews-1/) - Excellent deep dive into the human aspects
- [Compassionate Code Reviews](https://www.youtube.com/watch?v=Ea8EiIPZvh0) - April Wensel's talk on empathetic reviewing
- [The Role of Psychological Safety in Software Quality](https://link.springer.com/article/10.1007/s10664-024-10512-1) - Academic research on team dynamics

### 🎯 Industry Best Practices
**How successful companies do code reviews**
- [Google's Code Review Best Practices](https://google.github.io/eng-practices/review/) - The gold standard guide from Google
- [Microsoft's Code Review Culture](https://devblogs.microsoft.com/appcenter/how-the-visual-studio-mobile-center-team-does-code-review/) - Real-world emoji usage and team practices
- [GitHub Staff Engineer's Review Philosophy](https://github.blog/developer-skills/github/how-to-review-code-effectively-a-github-staff-engineers-philosophy/) - Insights from 7,000+ reviews
- [Code Review Best Practices](https://www.michaelagreiler.com/code-review-best-practices/) - 30+ proven practices from Microsoft research

### 🛠️ Tools & Systems
**Practical systems for better reviews**
- [Code Review Emoji Guide (CREG)](https://github.com/erikthedeveloper/code-review-emoji-guide) - The original emoji system
- [Conventional Comments](https://conventionalcomments.org/) - Structured labeling system for feedback
- [GitHub Suggestions Feature](https://docs.github.com/articles/incorporating-feedback-in-your-pull-request) - How to give specific, actionable feedback
- [CREG Browser Extension](https://www.raycast.com/russellyeo/code-review-emojis) - Tools to make emoji reviews easier

### 📊 Research & Studies
**Academic and industry research on code review effectiveness**
- [Modern Code Review: A Case Study at Google](https://research.google/pubs/modern-code-review-a-case-study-at-google/) - Comprehensive research on review practices
- [Cognitive Load in Code Reviews](https://link.springer.com/article/10.1007/s10664-022-10123-8) - Why brevity and clarity matter
- [Pull Request Size Matters](https://bssw.io/blog_posts/pull-request-size-matters) - Research on optimal PR sizes and exceptions
- [Code Review Best Practices Research](https://smartbear.com/learn/code-review/best-practices-for-peer-code-review/) - Industry studies on effective review processes

### ⚠️ Anti-Patterns & Pitfalls
**What NOT to do in code reviews**
- [Toxic Code Review Culture](https://blog.submain.com/toxic-code-review-culture/) - Warning signs and how to avoid them
- [Code Review Anti-Patterns](https://www.chiark.greenend.org.uk/~sgtatham/quasiblog/code-review-antipatterns/) - Common mistakes that damage teams
- [Unlearning Toxic Behaviors](https://medium.com/@sandya.sankarram/unlearning-toxic-behaviors-in-a-code-review-culture-b7c295452a3c) - How to recover from bad review cultures

### 🚀 Advanced Topics
**For teams ready to level up their review game**
- [Inclusion in Code Review](https://microsoft.github.io/code-with-engineering-playbook/code-reviews/inclusion-in-code-review/) - Making reviews work for diverse teams
- [AI-Assisted Code Reviews](https://newsletter.getdx.com/p/ai-assisted-code-reviews-at-google) - How AI can enhance (not replace) human reviews
- [Remote-First Code Review](https://medium.com/bbc-product-technology/looks-good-to-me-making-code-reviews-better-for-remote-first-teams-95bd92ee4e27) - Adapting reviews for distributed teams
- [Code Review as Mentorship](https://smartbear.com/blog/developing-a-culture-of-mentorship-with-code-revie/) - Using reviews for knowledge transfer
