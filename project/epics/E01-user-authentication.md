---
id: E01
title: User Authentication & Profile Management
status: To Do
owner: ProductManager
---

### Summary

This epic covers all work required for user authentication via WeChat login, JWT token management, and user profile functionality. Users should be able to log in with their WeChat account and view their profile information including credit balance.

### Related Documents

- `@/docs/PRD.md#Feature-Credit-Monetization-System`
- `@/docs/TECH_SPEC.md#API-Design`
- `@/docs/DATA_MAP.md#users`
- `@/docs/UX_FLOW.md#Flow-2-User-Onboarding-Login`

### Tasks in this Epic

- T001: Create database migration for users table
- T002: Implement WeChat OAuth login backend endpoint
- T003: Implement JWT token generation and validation
- T004: Create user profile API endpoints
- T005: Build login UI with WeChat SDK integration
- T006: Create user profile page UI
- T007: Implement authentication state management 