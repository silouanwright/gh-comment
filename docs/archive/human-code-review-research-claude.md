# Human Code Review Communication: Research for Natural AI-Generated Comments

This comprehensive research examines how experienced developers communicate during code reviews, synthesizing findings from major tech companies, academic studies, and industry experts to inform AI prompt design for more natural and effective review comments.

## The Psychology Behind Effective Code Review Feedback

Code review psychology fundamentally shapes how developers receive and respond to feedback. Research from Harvard's Amy Edmondson¹ and Google's Project Aristotle reveals that **psychological safety is the single most important factor** determining review effectiveness. When developers perceive criticism as threatening their professional identity, their amygdala triggers the same fight-or-flight response as physical danger, making constructive dialogue nearly impossible.

The most effective reviewers understand three psychological triggers that activate defensive responses.² **Truth triggers** fire when feedback feels inaccurate ("They got it wrong!"). **Relationship triggers** activate based on reviewer trust ("Why are they saying this?"). **Identity triggers** - the most powerful - occur when criticism attacks professional self-worth ("I'm a bad developer"). Successful review cultures actively design around these triggers by separating code quality from developer identity and focusing on collaborative improvement rather than error detection.

## Communication Patterns That Build Rather Than Break

The difference between mechanical and constructive feedback lies in subtle but crucial communication patterns. Google's engineering practices³ demonstrate that **framing feedback as questions rather than commands** dramatically improves reception. Instead of "Rename to sendMessage," effective reviewers ask "What do you think about calling this 'sendMessage'? It might make the intent clearer." This approach invites collaboration rather than demanding compliance.

Microsoft's research⁴ confirms that using **"we" language** creates psychological partnership. "Can we break this function into smaller methods?" feels fundamentally different from "You should refactor this function." The most respected reviewers across all studied companies consistently frame feedback as collaborative exploration, provide context and reasoning for suggestions, offer alternatives rather than just pointing out problems, and explicitly acknowledge good practices when they find them.

**The OIR framework**⁵ (Observation, Impact, Request) provides a structured approach that reduces defensiveness. Rather than "This code is confusing," reviewers state: "I notice this function handles multiple responsibilities (Observation), which makes it harder to test and maintain (Impact). Could we consider breaking it into separate functions? (Request)."

## Strategic Use of Human Elements in Technical Communication

Research from 2023-2024 reveals that **strategic emoji use significantly improves code review communication**¹². Studies show emojis reduce misinterpretation of written feedback, increase developer satisfaction, and clarify intent when tone might be ambiguous. The Code Review Emoji Guide (CREG) framework¹³ has emerged as an industry standard, with specific emojis conveying clear intent: 🔧 for necessary changes, ❓ for questions, 🤔 for thinking out loud, and ⛏️ for non-blocking nitpicks.

Microsoft's Mobile Center team¹⁴ reports that emoji use adds a "human component to conversation" and helps separate well-meant suggestions from mandatory changes, reducing hours of miscommunication. However, emojis work best when used strategically - to soften potentially harsh feedback ("🤔 This logic seems complex - what if we broke it down?"), clarify non-blocking comments ("⛏️ Style nitpick: consistent spacing would be nice"), or show appreciation ("😃 Great use of the factory pattern here!").

Casual language elements like contractions ("Can't we simplify this?" vs "Cannot we simplify this?") and conversational phrases ("This looks great, but I'm wondering...") make feedback feel more approachable while maintaining professionalism. The key is balancing humanity with technical accuracy - never using casual language for serious security or performance issues.

## Balancing Review Thoroughness with Team Velocity

Google achieves remarkable efficiency with median review times under 4 hours for most changes¹⁵, significantly outperforming industry averages. Their success stems from clear principles about when to write longer versus shorter comments. **Longer explanatory comments** serve new team members who need educational context, complex architectural decisions requiring "why" explanations, and security or performance issues needing detailed solutions. **Brief comments** work for experienced developers, obvious bugs, and repeated patterns where you can reference earlier feedback.

The most effective teams follow a **high-level first strategy**: initial reviews focus on architecture and major logic flaws, subsequent rounds address implementation details, and final passes handle minor optimizations. Research shows that reviews exceeding 20-50 comments in a single round enter a "danger zone" where developers become overwhelmed¹⁶. Successful teams group related issues and focus on the **80/20 rule** - 80% of value comes from addressing major functionality, security, and maintainability issues, while only 20% comes from style preferences and minor optimizations.

## Workflow Patterns from High-Performing Teams

Different code changes require fundamentally different review approaches. For **feature development**, reviewers focus on architecture and integration points while ensuring adequate documentation. **Bug fixes** demand root cause analysis and verification that fixes don't introduce new issues. **Refactoring** requires justifying business value and ensuring behavior preservation through comprehensive testing. **Documentation changes** need accuracy verification and clarity assessment for the target audience.

The research reveals that teams with sub-24-hour review response times show **40% better cycle time performance**. Optimal practices include initial responses within 4-6 hours during overlapping hours, review sessions limited to 60-90 minutes to maintain focus, and inspection rates of 300-500 lines per hour for thorough analysis. Changes under 225 lines receive the most effective reviews, with quality dropping significantly beyond 400 lines.

Successful teams strategically balance **asynchronous and synchronous reviews**. Async reviews work best for standard development with clear scope, while sync reviews excel for complex architectural decisions, when async threads exceed 3 back-and-forth exchanges, or during new developer onboarding. The key to making async reviews feel more personal involves positive reinforcement of good solutions, asking "why" questions rather than making demands, acknowledging effort while assuming good intent, and explaining reasoning behind all suggestions.

## Toxic Patterns That Destroy Review Culture

The research identifies devastating anti-patterns that teams must actively prevent¹⁷. **"Death by a Thousand Round Trips"** occurs when reviewers provide feedback incrementally, forcing multiple revision cycles for issues that could have been caught initially. This can extend review cycles from days to weeks and signals profound disrespect for developer time.

**"The Hostage Situation"** involves blocking pull requests until unrelated work is completed - either demanding personal style preferences not mandated by the organization or requiring substantial refactors of untouched legacy code. As one developer noted¹⁸, this "strong-arms the submitter into resolving historic issues by using their need to complete assigned work as leverage, which is bullying."

**Toxic communication patterns**¹⁹ include sarcastic or hostile language ("Did you even test this?"), judgmental questions implying incompetence ("Why didn't you just use a constants file?"), stating opinions as absolute facts, and using negative emojis without explanation. The research shows these patterns directly cause developer burnout, with 68-80% of developers reporting burnout symptoms, much of it attributed to poor review processes.

## Building Positive Review Cultures: Lessons from Industry Leaders

The most compelling transformation stories demonstrate that **successful review cultures evolve organically rather than through mandates**. At Airbnb, the shift began when "a few motivated engineers" started highlighting great code reviews in weekly meetings. This grassroots approach reached a cultural tipping point where not requesting reviews became unusual, transforming from "juggling chainsaws blindfolded" to a quality-focused culture where "we literally can't afford even small mistakes."

Google's approach treats code review primarily as an educational tool rather than a bug-catching mechanism. Their culture emphasizes that reviews ensure code meets standards, follows best practices, and that the team understands and can maintain the code going forward. This learning-focused approach, combined with tools that clearly indicate whose "turn" it is to act, enables their industry-leading review speeds.

Microsoft's inclusion-first design addresses the realities of 60,000 engineers collaborating globally. They explicitly distinguish between opinions, preferences, best practices, and facts while building awareness of how communication styles vary across cultures. Their approach specifically addresses impostor syndrome and creates pathways for underrepresented developers to participate fully in the review process.

## Remote Team Success Patterns

Successful remote teams **design for asynchronous collaboration** from the ground up. They over-document context with detailed PR descriptions, clear acceptance criteria, relevant documentation links, and even video explanations for complex changes. Time zone optimization involves identifying "golden hours" when multiple zones overlap and using these strategically for synchronous discussions when needed.

The BBC's rapid transformation to remote-first development revealed that code reviews are "a source of difficulties" when scaling with distributed teams. Their solution involved creating anonymous collaborative documents to source honest feedback about review processes, leading to human-centered approaches that treat every review as "an opportunity to ask questions, share knowledge and consider alternative ways of doing things."

## Making Code Reviews Natural Mentorship Opportunities

The most successful teams explicitly design reviews as mentorship platforms. Natural mentoring during reviews accelerates learning because lessons are applied daily in real work contexts. Effective mentorship involves providing specific, actionable feedback with clear explanations of the "why" behind suggestions. Successful teams balance criticism with recognition and use collaborative language that builds confidence.

For junior developers, creating safe environments is crucial. As one developer noted, "There is no such thing as a stupid question. If someone doesn't understand a bit of code, they need to feel free to ask." Teams that excel at junior participation set explicit expectations that everyone's perspective is valuable and learning flows in both directions. Some use pair review sessions where junior and senior developers review together, building both confidence and skills.

## The Critical Importance of Concise Communication

Recent research reveals that **cognitive load fundamentally limits review effectiveness**⁶, making conciseness a critical factor rather than a nice-to-have preference. Studies show that humans can only hold approximately **4 "chunks" of information in working memory**⁷ simultaneously, and reviewer attention degrades linearly after just 10 minutes⁸. This means lengthy, verbose comments actively harm comprehension by overwhelming cognitive capacity.

Google's engineering practices⁹ demonstrate the power of brevity through strategic use of **"Nit:" prefixes** to separate essential feedback from polish suggestions, allowing reviewers to quickly categorize and prioritize comments. Their approach focuses on **signal-to-noise ratio** - every word should contribute to understanding, with extraneous information actively removed to reduce cognitive burden.

The research on cognitive load in software engineering¹⁰ shows that **ambiguous or verbose feedback creates extraneous cognitive load** as reviewers expend mental effort parsing unnecessarily complex language instead of focusing on the code issues. Effective comments follow the principle of **"expressing an idea with the fewest characters possible, but no fewer"**¹¹ - achieving precision through economy of language rather than elaborate explanation.

**Practical conciseness patterns** from successful teams include using questions instead of lengthy explanations ("What about using a Map here?" vs "I think you should consider refactoring this approach to use a Map data structure because it would provide better performance characteristics"), leveraging code suggestions for specific changes rather than verbose descriptions, and providing context links instead of embedding full explanations ("See the auth service pattern for similar validation").

The key insight is that **conciseness respects developer time and cognitive capacity**. When reviewers can quickly understand feedback and take action, review cycles accelerate and quality improves. Verbose comments, even well-intentioned ones, can actually slow down the review process by forcing developers to parse unnecessary information to find actionable insights.

## Key Principles for Natural AI-Generated Reviews

Based on this research, AI-generated code review comments should embody several core principles to feel natural and effective. **Frame suggestions as questions** rather than commands, using phrases like "What do you think about..." or "Have you considered..." to invite collaboration. **Provide context and reasoning** for every suggestion, explaining the why behind recommendations. **Acknowledge positives** explicitly when identifying good practices or clever solutions.

**Use strategic emojis** to clarify intent and add human warmth - 🤔 for thoughtful suggestions, ❓ for genuine questions, 😃 for praise, and ⛏️ for minor nitpicks. **Employ collaborative language** with "we" statements and inclusive phrasing. **Match communication depth to context** - provide more explanation for junior developers or complex issues, but keep it brief for experienced developers on straightforward changes.

**Separate observations from judgments** by stating what the code does rather than labeling it as good or bad. **Offer specific alternatives** with code examples rather than vague suggestions. **Respect developer time** by batching related feedback and clearly distinguishing between blocking issues and nice-to-have improvements. Most importantly, **maintain a learning orientation** that treats every review as an opportunity for mutual growth rather than one-sided evaluation.

## Conclusion: The Human Heart of Technical Excellence

This research demonstrates that creating effective code review cultures requires intentional focus on human psychology, structured communication frameworks, and sustained cultural evolution. The technical aspects of code review pale in importance compared to the human elements that determine whether feedback leads to growth or defensiveness. As teams increasingly adopt AI assistance in code reviews, these human-centered principles become even more critical to ensure that automation enhances rather than replaces the collaborative learning that makes code reviews valuable.

The future of code review will likely involve continued evolution toward more intelligent tooling and better distributed collaboration support. However, the fundamental human elements - psychological safety, constructive communication, mentorship, and team building - will remain central to success. AI-generated review comments must embody these principles to feel natural and drive the positive outcomes that the best human reviewers achieve through years of experience and emotional intelligence.

---

## References

1. Edmondson, A. (1999). Psychological Safety and Learning Behavior in Work Teams. *Administrative Science Quarterly*, 44(2), 350-383. https://journals.sagepub.com/doi/abs/10.2307/2666999

2. Artstain, R. (2024). Once more, with feeling: A radical approach to code review. https://rinaarts.com/once-more-with-feeling-a-radical-approach-to-code-review/

3. Google Engineering Practices. (2024). How to write code review comments. https://google.github.io/eng-practices/review/reviewer/comments.html

4. Greiler, M. (2024). How Code Reviews work at Microsoft. https://www.michaelagreiler.com/code-reviews-at-microsoft-how-to-code-review-at-a-large-software-company/

5. Yoodli. (2024). How to Use the Situation-Behavior-Impact (SBI)™ Feedback Model. https://yoodli.ai/blog/sbi-feedback-model

6. Andaloussi, A. A., et al. (2022). Do explicit review strategies improve code review performance? Towards understanding the role of cognitive load. *Empirical Software Engineering*, 27, 1-29. https://link.springer.com/article/10.1007/s10664-022-10123-8

7. Baddeley, A. (1992). Working memory. *Science*, 255(5044), 556-559.

8. Giese, A. (2020). Cracking the code review (Part 2): Make them seem small. https://gieseanw.wordpress.com/2020/06/25/cracking-the-code-review-part-2-make-them-seem-small/

9. Google Engineering Practices. (2024). The Standard of Code Review. https://google.github.io/eng-practices/review/reviewer/standard.html

10. Costa, M., et al. (2021). Measuring the cognitive load of software developers: An extended Systematic Mapping Study. *Information and Software Technology*, 131, 106491. https://www.sciencedirect.com/science/article/abs/pii/S095058492100046X

11. Fry, A. (2022). Code Review How To: Brevity and Repetition. https://andyfry.co/code-review-how-to-brevity-repetition

12. Code Review Emoji Guide. (2024). GitHub repository. https://github.com/erikthedeveloper/code-review-emoji-guide

13. Conventional Comments. (2024). https://conventionalcomments.org/

14. Microsoft Developer Blogs. (2024). How We Do Code Review - App Center Blog. https://devblogs.microsoft.com/appcenter/how-the-visual-studio-mobile-center-team-does-code-review/

15. Engineer's Codex. (2024). How Google takes the pain out of code reviews, with 97% dev satisfaction. https://read.engineerscodex.com/p/how-google-takes-the-pain-out-of

16. Research.Google. (2024). Modern Code Review: A Case Study at Google. https://research.google/pubs/modern-code-review-a-case-study-at-google/

17. Tatham, S. (2024). Code review antipatterns. https://www.chiark.greenend.org.uk/~sgtatham/quasiblog/code-review-antipatterns/

18. AWS Well-Architected. (2024). Anti-patterns for code review. https://docs.aws.amazon.com/wellarchitected/latest/devops-guidance/anti-patterns-for-code-review.html

19. Sankarram, S. (2024). Unlearning toxic behaviors in a code review culture. *Medium*. https://medium.com/@sandya.sankarram/unlearning-toxic-behaviors-in-a-code-review-culture-b7c295452a3c

## Additional Resources

20. Lynch, M. (2024). How to Do Code Reviews Like a Human (Part One). https://mtlynch.io/human-code-reviews-1/

21. Feather, I. (2024). Radical Candor in Code Review. https://www.ianfeather.co.uk/radical-candor-in-code-review/

22. Greiler, M. (2024). How to Give Respectful and Constructive Code Review Feedback. https://www.michaelagreiler.com/respectful-constructive-code-review-feedback/

23. Bos, A. (2024). My Case for Conventional Comments. https://aaronbos.dev/posts/case-for-conventional-comments

24. Microsoft Engineering Playbook. (2024). Inclusion in Code Review. https://microsoft.github.io/code-with-engineering-playbook/code-reviews/inclusion-in-code-review/

25. BBC Product & Technology. (2024). Looks Good To Me: Making code reviews better for remote-first teams. https://medium.com/bbc-product-technology/looks-good-to-me-making-code-reviews-better-for-remote-first-teams-95bd92ee4e27

26. SmartBear. (2024). Developing a Culture of Mentorship with Code Review. https://smartbear.com/blog/developing-a-culture-of-mentorship-with-code-revie/

27. Elbre, E. (2018). Psychology of Code Readability. *Medium*. https://medium.com/@egonelbre/psychology-of-code-readability-d23b1ff1258a

28. GitHub. (2024). Cognitive load is what matters. https://github.com/zakirullin/cognitive-load

29. Stitcher.io. (2024). A programmer's cognitive load. https://stitcher.io/blog/a-programmers-cognitive-load

30. Demircioğlu, A. (2024). Cognitive load in software engineering. *Medium*. https://atakde.medium.com/cognitive-load-in-software-engineering-6e9059266b79
