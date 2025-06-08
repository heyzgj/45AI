---
id: T009
title: Implement template listing API endpoint
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

Create a backend API endpoint (`GET /api/v1/templates`) that retrieves all active templates from the database and returns them as a JSON array.

### Acceptance Criteria

- [ ] API endpoint `GET /api/v1/templates` is created.
- [ ] The endpoint fetches all records from the `templates` table where `is_active` is true.
- [ ] The returned JSON array includes `id`, `name`, `preview_image_url`, and `credit_cost` for each template.
- [ ] The endpoint does not require authentication.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#API-Design`
- **Code:** `@/backend/internal/handler/template.go`, `@/backend/internal/repository/template.go`

### Agent Notes

*This is a public-facing endpoint and should be efficient. Consider future pagination.* 