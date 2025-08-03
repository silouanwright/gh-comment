---
name: architecture-review
title: Architecture & Design Pattern Review
category: architecture
estimated_time: 20-25 minutes
tags: [architecture, design patterns, scalability, maintainability]
examples:
  - "Use with: gh comment add 123 service.go 134:150 \"[architectural feedback]\""
  - "Great for: Major feature reviews, refactoring planning, technical debt reduction"
---

Please provide an architectural review of this pull request, focusing on design patterns, code organization, and long-term maintainability.

## Architecture Analysis Framework

1. **Design Patterns & Principles**
   ♻️ Evaluate adherence to SOLID principles
   ♻️ Check for appropriate use of design patterns (Factory, Observer, etc.)
   ♻️ Look for violations of DRY, KISS, YAGNI principles

2. **Code Organization & Structure**
   ♻️ Review module boundaries, separation of concerns
   ♻️ Check for proper layering (presentation, business, data)
   ♻️ Evaluate package/folder structure, naming conventions

3. **Dependency Management**
   ♻️ Look for tight coupling, circular dependencies
   ♻️ Check for proper dependency injection patterns
   ♻️ Review interface usage, abstraction levels

4. **Scalability & Future Growth**
   ♻️ Identify potential bottlenecks for scaling
   ♻️ Check for extensibility points, plugin architectures
   ♻️ Look for configuration management, feature toggles

5. **Error Handling & Resilience**
   ♻️ Review error propagation patterns, retry logic
   ♻️ Check for circuit breakers, graceful degradation
   ♻️ Evaluate logging, monitoring integration points

6. **API Design**
   ♻️ Check for consistent API patterns, versioning strategy
   ♻️ Review request/response schemas, data contracts
   ♻️ Look for proper HTTP status codes, error responses

For each architectural concern:
- Use ♻️ emoji for refactoring suggestions
- Explain the long-term benefits of proposed changes
- Consider team knowledge and migration complexity
- Provide [SUGGEST: refactored_code] for structural improvements

Focus on changes that improve maintainability and team productivity.