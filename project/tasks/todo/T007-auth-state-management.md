---
id: T007
title: Implement authentication state management
status: To Do
epic: E01
effort: M
risk: Low
dependencies: [T040, T005]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-22
---

### Description

Implement a Pinia store to manage authentication state, including the JWT and the user's profile information.

### Acceptance Criteria

- [ ] A new Pinia store module is created at `/frontend/src/stores/auth.js`.
- [ ] The store has state for `token`, `user`, and `isAuthenticated`.
- [ ] It has actions for `login`, `logout`, and `fetchUser`.
- [ ] The `login` action saves the token to device storage (`uni.setStorageSync`).
- [ ] The store can be initialized from storage, so the user stays logged in across app sessions.
- [ ] The API client is configured to send the `token` in the `Authorization` header for all authenticated requests.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#Auth-Method`
- **Code:** `@/frontend/src/stores/auth.js`, `@/frontend/src/api/index.js`

### Agent Notes

*Ensure sensitive data is not stored in a way that is easily accessible. The token should be the primary piece of stored auth data.* 