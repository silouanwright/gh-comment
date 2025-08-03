---
name: performance-optimization
title: Performance Optimization Review
category: performance
estimated_time: 10-15 minutes
tags: [performance, optimization, scaling, efficiency]
examples:
  - "Use with: gh comment add 123 database.py 89 \"[performance analysis for specific function]\""
  - "Great for: Pre-deployment reviews, scalability planning, optimization sprints"
---

Please conduct a comprehensive performance review of this pull request, focusing on identifying optimization opportunities and potential bottlenecks.

## Performance Analysis Areas

1. **Database Performance**
   🤔 Look for N+1 query problems, missing indexes, inefficient queries
   🤔 Check for proper connection pooling and query optimization
   🤔 Identify opportunities for caching, batch operations

2. **Memory Management** 
   🤔 Check for memory leaks, excessive object creation
   🤔 Look for inefficient data structures, large object retention
   🤔 Verify proper cleanup of resources, event listeners

3. **Computational Efficiency**
   🤔 Identify expensive operations in hot paths, loops
   🤔 Look for redundant calculations, unnecessary processing
   🤔 Check for proper algorithm complexity (O(n) vs O(n²))

4. **Network & I/O Optimization**
   🤔 Review API call patterns, bundling opportunities  
   🤔 Check for proper async/await usage, parallel processing
   🤔 Look for image optimization, asset compression

5. **Frontend Performance**
   🤔 Check for unnecessary re-renders, large bundle sizes
   🤔 Look for lazy loading opportunities, code splitting
   🤔 Verify proper caching strategies

For each optimization opportunity:
- Use 🤔 emoji for performance suggestions (not blocking)
- Provide [SUGGEST: optimized_code] with more efficient alternatives
- Quantify potential impact where possible (e.g., "reduces database calls by 80%")
- Consider trade-offs between performance and code readability

Focus on changes that will have measurable impact on user experience.