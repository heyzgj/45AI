---
id: T020
title: Implement image generation API endpoint
status: Done
epic: E03
effort: S
risk: M
dependencies: [T019]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Register the image generation API endpoint (`POST /api/v1/generate`) and protect it with the authentication middleware.

### Acceptance Criteria

- [x] The `POST /api/v1/generate` endpoint is registered in `main.go`.
- [x] The endpoint is protected by the `AuthMiddleware`.
- [x] The endpoint is mapped to the `GenerationHandler.GenerateImage` method.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#API-Design`
- **Code:** `@/backend/cmd/api/main.go`

### Agent Notes

*This endpoint was implemented as part of task T018. Marking as done.* 