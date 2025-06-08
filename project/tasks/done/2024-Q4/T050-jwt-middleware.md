---
id: T050
title: Implement JWT authentication middleware
status: Done
epic: E06
effort: M
risk: Medium
dependencies: []
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-23
---

### Description

Implement the logic for generating new JWTs for authenticated users and a Gin middleware for validating JWTs on protected API routes. The JWT should contain the user ID.

### Acceptance Criteria

- [ ] A function `GenerateToken(userID uint)` exists and returns a valid JWT string.
- [ ] A Gin middleware `AuthMiddleware()` is created in `/backend/internal/middleware/`.
- [ ] The middleware correctly extracts the token from the `Authorization` header.
- [ ] The middleware validates the token's signature and expiration using a secret key from the config.
- [ ] If valid, the user ID is extracted from the token and added to the request context (e.g., `c.Set("userID", userID)`).
- [ ] If invalid, the middleware aborts the request with a `401 Unauthorized` error and a clear JSON response.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#Auth-Method`
- **Code:** `@/backend/internal/middleware/auth.go`, `@/backend/internal/service/auth.go`

### Agent Notes

*Use the `golang-jwt/jwt/v5` library. The secret key for signing tokens must be stored securely and retrieved from the application configuration (`config.JWTSecret`).* 