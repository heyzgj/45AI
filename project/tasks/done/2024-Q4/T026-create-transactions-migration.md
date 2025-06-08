---
id: T026
title: Create database migration for transactions table
status: Done
epic: E04
effort: S
risk: Low
dependencies: [T049]
assignee: CursorAgent
created_at: 2024-12-24
updated_at: 2024-12-24
---

### Description

Create the database migration script to define the `transactions` table schema as specified in the Data Map. This table will log all credit changes for auditing and user history.

### Acceptance Criteria

- [ ] Migration file created with proper naming convention
- [ ] Transactions table includes all required fields with correct types
- [ ] Foreign key constraints properly defined
- [ ] ENUM type for transaction type (purchase/generation)
- [ ] Proper indexes for performance
- [ ] Migration can be rolled back cleanly

### Context Binding

- **Docs:** `@/docs/DATA_MAP.md#transactions`
- **Code:** `@/backend/migrations/`

### Agent Notes

Ensure foreign key constraints reference users and templates tables. The amount field should support both positive (purchases) and negative (generation) values. 