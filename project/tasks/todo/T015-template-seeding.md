---
id: T015
title: Implement template data seeding script
status: To Do
epic: E02
effort: S
risk: Low
dependencies: [T008]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-22
---

### Description

Create a script to populate the `templates` table with initial data. This is crucial for development and testing environments.

### Acceptance Criteria

- [ ] A Go script is created, likely in `/backend/cmd/seeder/`.
- [ ] When run, the script connects to the database.
- [ ] It inserts a predefined list of 5-10 templates into the `templates` table.
- [ ] The script is idempotent; running it multiple times does not create duplicate templates.

### Context Binding

- **Docs:** `@/docs/DATA_MAP.md#templates`
- **Code:** `@/backend/cmd/seeder/main.go`

### Agent Notes

*This can be a simple standalone Go program that uses the existing repository layer to insert data.* 