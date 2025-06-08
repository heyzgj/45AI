---
id: T005
title: Build login UI with WeChat SDK integration
status: To Do
epic: E01
effort: M
risk: Medium
dependencies: [T002]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-22
---

### Description

Create the frontend login page. This page should feature a single "Login with WeChat" button that uses the UniBest/WeChat SDK to initiate the login flow, retrieve the `code`, and send it to the backend login endpoint.

### Acceptance Criteria

- [ ] A new page is created at `/frontend/src/pages/login/index.vue`.
- [ ] The page UI matches the minimalist design from the `STYLE_GUIDE.md` and `UX_FLOW.md`.
- [ ] The "Login with WeChat" button correctly calls `uni.login()` to get the user's `code`.
- [ ] On success, the `code` is sent to the backend's `POST /api/v1/auth/login` endpoint.
- [ ] The returned JWT is stored securely on the client (e.g., in Pinia store and uni-storage).
- [ ] User is redirected to the home page on successful login.

### Context Binding

- **Docs:** `@/docs/UX_FLOW.md#Flow-2-User-Onboarding-Login`, `@/docs/STYLE_GUIDE.md`
- **Code:** `@/frontend/src/pages/login/index.vue`, `@/frontend/src/stores/auth.js`

### Agent Notes

*The UI should be simple and elegant. Pay attention to the button's style and press-state animations as defined in the Style Guide.* 