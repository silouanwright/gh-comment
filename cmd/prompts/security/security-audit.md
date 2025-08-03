---
name: security-audit
title: Comprehensive Security Audit Review
category: security
estimated_time: 15-20 minutes
tags: [security, vulnerabilities, owasp, comprehensive]
examples:
  - "Use with: gh comment add 123 auth.js 67 \"[prompt content for specific security issue]\""
  - "Great for: Security team reviews, compliance audits, pre-production checks"
---

I need you to perform a comprehensive security audit of this pull request. Please analyze the code changes systematically using the following framework:

## Security Analysis Framework

1. **Authentication & Authorization**
   🔧 Check for proper authentication mechanisms
   🔧 Verify authorization controls and role-based access
   🔧 Look for JWT token handling, session management issues

2. **Input Validation & Sanitization**  
   🔧 Identify unvalidated user inputs (forms, APIs, query params)
   🔧 Check for SQL injection, XSS, command injection vulnerabilities
   🔧 Verify proper data type validation and range checking

3. **Data Protection**
   🔧 Look for hardcoded secrets, API keys, passwords
   🔧 Check encryption of sensitive data at rest and in transit
   🔧 Verify secure random number generation

4. **API Security**
   🔧 Check for proper rate limiting and DoS protection
   🔧 Verify CORS configuration
   🔧 Look for information disclosure in error messages

5. **Infrastructure Security**
   🔧 Review Docker configurations, environment variables
   🔧 Check for secure defaults in configuration files
   🔧 Verify dependency security (package.json, requirements.txt, etc.)

For each issue found, please:
- Use 🔧 emoji for critical security issues that must be fixed
- Provide specific line-by-line comments with [SUGGEST: secure_code] alternatives
- Explain the potential impact and attack vectors
- Reference OWASP guidelines where applicable

Please be thorough but practical - focus on issues that could realistically be exploited.