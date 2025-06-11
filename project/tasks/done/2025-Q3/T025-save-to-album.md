---
id: T025
title: Add save to album functionality
status: Done
epic: E03
effort: S
risk: L
dependencies: [T021]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Implement the functionality to save a generated image to the user's device photo album. This will involve using the UniBest API for saving images.

### Acceptance Criteria

- [ ] A "Save to Album" button is added to the image result page.
- [ ] Tapping the button saves the selected image to the user's photo album.
- [ ] A confirmation message is displayed to the user after the image is saved successfully.

### Context Binding

- **Docs:** `@/docs/UX_FLOW.md#Flow-1-Core-Image-Generation`
- **Code:** `@/frontend/src/pages/generate/index.vue`

### Agent Notes

*Use the `uni.saveImageToPhotosAlbum` API for this functionality. Ensure that the app has the necessary permissions to save images.* 