---
id: T038
title: Configure build scripts and environment variables
status: To Do
epic: E05
effort: S
risk: L
dependencies: [T037]
assignee: CursorAgent
---

### Description

Set up build scripts and environment variable configuration for the UniBest project. Configure different build targets for WeChat Mini Program and iOS app, with proper environment management for development, staging, and production.

### Acceptance Criteria

- [x] Build scripts for WeChat Mini Program compilation
- [x] Build scripts for iOS app compilation
- [x] Environment variable configuration (.env files)
- [x] Development, staging, and production configurations
- [x] Hot reload and development server setup

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#Build-Process`
- **Code:** `@/frontend/package.json`, `@/frontend/.env.*`, `@/frontend/vite.config.ts`

### Agent Notes

*UniBest uses Vite for building. Configure proper build targets and environment variables for different platforms.* 