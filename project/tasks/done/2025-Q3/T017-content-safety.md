---
id: T017
title: Implement content safety moderation
status: Done
epic: E03
effort: M
risk: M
dependencies: []
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Integrate a third-party content safety API (e.g., Tencent Cloud, Volcano Engine) to scan all user-uploaded images for unsafe content before they are sent to the image generation model.

### Acceptance Criteria

- [ ] A new service is created to interact with the content safety API.
- [ ] The image generation service calls this new service to validate the user's uploaded image.
- [ ] If the image is flagged as unsafe, the generation process is aborted, and an appropriate error is returned to the user.
- [ ] The API key and URL for the content safety service are managed through the application configuration.

### Context Binding

- **Docs:** `@/docs/PRD.md#Feature-Content-Safety--Moderation`, `@/docs/TECH_SPEC.md#External-Services`
- **Code:** `@/backend/internal/service/content_safety_service.go`, `@/backend/internal/service/generation_service.go`

### Agent Notes

*The specific implementation will depend on the chosen content safety provider. For now, a mock implementation can be created that simulates the API call and response.* 