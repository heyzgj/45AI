---
id: T018
title: Implement image generation service
status: Done
epic: E03
effort: L
risk: H
dependencies: [T017]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Implement the core logic for the image generation service. This service will orchestrate the entire image generation process, including validating the request, checking user credits, calling the content safety service, and interacting with the ComfyUI API.

### Acceptance Criteria

- [ ] The `GenerateImage` method in `generation_service_impl.go` is fully implemented.
- [ ] It validates the user's uploaded image using the `ValidateImage` method.
- [ ] It checks if the user has enough credits to perform the generation.
- [ ] It calls the `CheckContentSafety` method to moderate the image.
- [ ] It interacts with a (mocked for now) ComfyUI client to generate the image.
- [ ] It deducts the credits from the user's account and creates a transaction record.
- [ ] It returns the generated image URLs and the credits used.

### Context Binding

- **Docs:** `@/docs/PRD.md#Feature-Template-Based-Image-Generation`, `@/docs/TECH_SPEC.md#Components`
- **Code:** `@/backend/internal/service/generation_service.go`, `@/backend/internal/repository/comfyui_repository.go`

### Agent Notes

*This is a complex service that interacts with multiple other services and repositories. A mock implementation of the ComfyUI repository should be created for now.* 