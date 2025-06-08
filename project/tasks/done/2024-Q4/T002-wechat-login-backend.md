---
id: T002
title: Implement WeChat OAuth login backend endpoint
status: Done
epic: E01
effort: M
risk: Medium
dependencies: [T001, T049, T050]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-23
---

### Description

Create a backend API endpoint (`POST /api/v1/auth/login`) that accepts a WeChat `code`, exchanges it for a `session_key` and `openid` with the WeChat API, finds or creates a user in the database, and returns a JWT.

### Acceptance Criteria

- [ ] API endpoint `POST /api/v1/auth/login` is created.
- [ ] Endpoint correctly calls WeChat's API to exchange the code.
- [ ] A new user is created in the `users` table if the `wechat_openid` is new.
- [ ] An existing user is retrieved if the `wechat_openid` already exists.
- [ ] A valid JWT is returned upon successful login/registration.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#API-Design`, `@/docs/PRD.md#Feature-Credit-Monetization-System`
- **Code:** `@/backend/internal/handler/auth.go`, `@/backend/internal/service/auth.go`

### Agent Notes

*Remember to store `wechat_openid` in the `users` table. Securely handle the app secret when calling WeChat's API, likely using environment variables managed by the config service.* 