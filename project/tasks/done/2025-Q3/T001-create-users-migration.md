---
id: T001
title: Create database migration for users table
status: Done
epic: E01
effort: S
risk: Low
dependencies: [T049]
assignee: CursorAgent
created_at: 2024-12-24
updated_at: 2024-12-24
---

### Description

Create the database migration script to define the `users` table schema as specified in the Data Map. The table will store user account information including WeChat OpenID and credit balance.

### Acceptance Criteria

- [ ] Migration file created with proper naming convention
- [ ] Users table includes all required fields with correct types
- [ ] Proper indexes added for wechat_openid (unique)
- [ ] Timestamps automatically managed by database
- [ ] Migration can be rolled back cleanly

### Context Binding

- **Docs:** `@/docs/DATA_MAP.md#users`
- **Code:** `@/backend/internal/model/`
- **Code:** `@/backend/migrations/`

### Agent Notes

Ensure the migration follows WeChat Cloud MySQL 8.0 compatibility. Use BIGINT for id field to handle large user volumes. 