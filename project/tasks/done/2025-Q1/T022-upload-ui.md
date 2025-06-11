---
id: T022
title: Build upload UI with native file picker
status: To Do
epic: E03
effort: M
risk: L
dependencies: [T037, T039]
assignee: CursorAgent
---

### Description

Create a native image upload interface for both WeChat Mini Program and iOS app. The UI should allow users to select photos from their gallery or camera, with preview functionality and upload progress indicators.

### Acceptance Criteria

- [x] Native file picker integration for WeChat Mini Program (wx.chooseImage)
- [x] Native file picker integration for iOS app (Capacitor Camera plugin)
- [x] Image preview with crop/resize functionality
- [x] Upload progress indicator with percentage
- [x] Error handling for file size/format validation
- [x] Beautiful UI following Style Guide with animations

### Context Binding

- **Docs:** `@/docs/STYLE_GUIDE.md#Upload-Component`, `@/docs/UX_FLOW.md#Flow-1`
- **Code:** `@/frontend/src/components/ImageUpload/`, `@/frontend/src/pages/generate/`

### Agent Notes

*Platform-specific implementations required. WeChat uses wx.chooseImage, iOS uses Capacitor Camera API.* 