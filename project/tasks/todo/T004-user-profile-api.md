---
id: T004
title: Create user profile API endpoints
status: To Do
epic: E01
effort: S
risk: Low
dependencies: [T050]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-22
---

### Description

Create API endpoints to get the current user's profile (`GET /api/v1/me`) and to update it (`PUT /api/v1/me`). The GET endpoint should return user details like nickname, avatar, and credit balance.

### Acceptance Criteria

- [ ] `GET /api/v1/me` endpoint is created and protected by the JWT middleware.
- [ ] It retrieves the user ID from the context and fetches the user from the database.
- [ ] It returns the user's `nickname`, `avatar_url`, and `credits`.
- [ ] `PUT /api/v1/me` endpoint is created for future use (updating profile).

### Context Binding

- **Docs:** `@/docs/DATA_MAP.md#users`
- **Code:** `@/backend/internal/handler/user.go`, `@/backend/internal/service/user.go`

### Agent Notes

*Ensure the API only returns non-sensitive user data.* 