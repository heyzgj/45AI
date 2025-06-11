---
id: T006
title: Create user profile page UI
status: Done
epic: E01
effort: M
risk: Low
dependencies: [T004, T044]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-23
---

### Description

Create the user profile page which displays the user's avatar, nickname, and their current credit balance. It should also include a list of their recent transactions.

### Acceptance Criteria

- [ ] A new page is created at `/frontend/src/pages/profile/index.vue`.
- [ ] The page fetches user data from the `/api/v1/me` endpoint and displays the `nickname` and `avatar_url`.
- [ ] The page uses the `CreditDisplay` component to show the user's credit balance.
- [ ] The page fetches data from `/api/v1/me/transactions` and displays a list of transactions.
- [ ] The UI adheres to the `STYLE_GUIDE.md`, focusing on clean layout and typography.

### Context Binding

- **Docs:** `@/docs/DATA_MAP.md#users`, `@/docs/STYLE_GUIDE.md`
- **Code:** `@/frontend/src/pages/profile/index.vue`, `@/frontend/src/components/CreditDisplay/CreditDisplay.vue`, `@/frontend/src/api/user.js`

### Agent Notes

*This page is primarily for display. Ensure the layout is spacious and easy to read.* 