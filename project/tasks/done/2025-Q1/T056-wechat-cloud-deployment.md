---
id: T056
title: Set up WeChat Cloud Hosting deployment
status: To Do
epic: E06
effort: L
risk: M
dependencies: [T047, T055]
assignee: CursorAgent
---

### Description

Configure deployment to WeChat Cloud Hosting platform for the Go backend. Create necessary configuration files, environment setup, and deployment scripts for WeChat Mini Program backend hosting.

### Acceptance Criteria

- [x] WeChat Cloud Hosting configuration file is created
- [x] Environment variables are properly configured for cloud deployment
- [x] Database connection is configured for cloud MySQL
- [x] Static asset serving is configured
- [x] Deployment script automates the release process

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#Infrastructure`, `@/docs/ARCHITECTURE_GUIDE.md#Deployment`
- **Code:** `@/backend/wechat-cloud-config.yaml`, `@/deploy/wechat-cloud.sh`

### Agent Notes

*WeChat Cloud Hosting has specific requirements for Go applications and database connections.* 