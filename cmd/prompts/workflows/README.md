---
name: workflow-index
title: AI Workflow Prompts for gh-comment
category: workflows
estimated_time: 2-3 minutes
tags: [index, overview, workflow-guide]
examples:
  - "Use for: Overview of available AI workflow prompts"
  - "Great for: Getting started with gh-comment AI workflows"
---

# 🤖 AI Workflow Prompts for gh-comment

This directory contains ready-to-use AI prompts designed to work seamlessly with the `gh-comment` CLI extension. These prompts are optimized for AI assistants to perform common GitHub PR workflows efficiently.

## 📋 Available Workflows

### 1. **PR Code Review** (`pr-code-review.md`)
Comprehensive code review workflow that grades PRs, identifies issues, and submits batch review comments.

### 2. **Self-Review Documentation** (`self-review.md`)
Document your own PR changes by explaining the reasoning, decisions, and context behind your code.

### 3. **Address Review Comments** (`address-review.md`)
Systematically respond to reviewer feedback, fix issues, and resolve conversations with proper decorum.

## 🚀 How to Use

1. **Copy the prompt** from the relevant workflow file
2. **Paste into your AI assistant** (Claude, ChatGPT, etc.)
3. **Provide the PR context** (checkout the PR branch first)
4. **Let AI analyze and comment** using `gh-comment` commands

### 💡 Suggestion Syntax
Use `[SUGGEST: code]` in your comments for quick code suggestions that can be auto-expanded!

## 💡 Prerequisites

- GitHub CLI installed and authenticated
- `gh-comment` extension installed
- PR checked out locally: `gh pr checkout <PR_NUMBER>`

## 🔧 Customization

Feel free to modify these prompts to match your team's:
- Coding standards
- Review criteria
- Comment tone and style
- Grading scales

## 📝 Contributing

Have a workflow idea? Create a new prompt file following the existing format and submit a PR!
