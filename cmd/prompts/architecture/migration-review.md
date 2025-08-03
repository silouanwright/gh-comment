---
name: migration-review
title: Code Migration & Refactoring Review
category: architecture
estimated_time: 15-20 minutes
tags: [migration, refactoring, legacy, modernization]
examples:
  - "Use with: gh comment add 123 migration.sql 45 \"[migration safety feedback]\""
  - "Great for: Database migrations, framework upgrades, architecture changes"
---

Please review this migration/refactoring pull request with focus on safety, completeness, and maintaining functionality during the transition.

## Migration Review Framework

1. **Backwards Compatibility**
   🔧 Verify existing APIs remain functional
   🔧 Check for proper deprecation warnings, migration paths
   🔧 Look for database schema compatibility

2. **Data Integrity**
   🔧 Review data migration scripts, rollback procedures
   🔧 Check for proper data validation, transformation logic
   🔧 Verify no data loss during migration

3. **Functionality Preservation**
   📝 Ensure all existing features work as before
   📝 Check for proper test coverage of migrated code
   📝 Look for edge cases that might be missed

4. **Performance Impact**
   🤔 Compare performance before/after migration
   🤔 Look for potential bottlenecks introduced
   🤔 Check for proper caching, optimization strategies

5. **Rollback Strategy**
   🔧 Verify rollback procedures are documented, tested
   🔧 Check for feature flags, gradual rollout capabilities
   🔧 Look for monitoring, alerting during migration

6. **Team Knowledge Transfer**
   📝 Check for updated documentation, runbooks
   📝 Look for training materials, migration guides
   📝 Verify team readiness for new patterns

For migration feedback:
- Prioritize safety and risk mitigation
- Consider phased rollout strategies
- Focus on maintaining system stability
- Provide [SUGGEST: safer_approach] alternatives

Remember: Migrations are high-risk operations that require extra scrutiny.