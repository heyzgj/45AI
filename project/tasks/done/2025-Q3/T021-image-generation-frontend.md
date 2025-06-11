---
id: T021
title: Implement image generation frontend
status: Done
epic: E03
effort: L
risk: M
dependencies: [T020]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Implement the frontend logic for the image generation flow. This includes handling the image upload, calling the backend API, and displaying the results.

### Acceptance Criteria

- [ ] The "Generate with this Template" button on the template detail page is functional.
- [ ] Tapping the button opens the native file picker to select an image.
- [ ] The selected image is uploaded to the backend's `/api/v1/generate` endpoint along with the template ID.
- [ ] A loading animation is displayed while the image is being generated.
- [ ] The generated images are displayed on a result page.

### Context Binding

- **Docs:** `@/docs/UX_FLOW.md#Flow-1-Core-Image-Generation`
- **Code:** `@/frontend/src/pages/template-detail/index.vue`, `@/frontend/src/pages/generate/index.vue`

### Agent Notes

*The result page can be a new page or a modal. For now, focus on the flow of uploading the image and receiving the response.* 