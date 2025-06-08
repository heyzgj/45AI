---
id: E03
title: AI Image Generation Pipeline
status: To Do
owner: ProductManager
---

### Summary

This epic covers the core image generation functionality including photo upload, content safety validation, credit deduction, ComfyUI integration, and result display. The end-to-end generation time must be under 30 seconds as per NFRs.

### Related Documents

- `@/docs/PRD.md#Feature-Template-Based-Image-Generation`
- `@/docs/PRD.md#Feature-Content-Safety-Moderation`
- `@/docs/TECH_SPEC.md#Components`
- `@/docs/UX_FLOW.md#Flow-1-Core-Image-Generation`

### Tasks in this Epic

- T016: Set up ComfyUI API on GCP with GPU
- T017: Implement image upload endpoint with validation
- T018: Integrate Tencent Cloud content safety API
- T019: Build ComfyUI service integration layer
- T020: Implement generation request queueing system
- T021: Create image generation API endpoint
- T022: Build upload UI with native file picker
- T023: Create loading screen with pulsing animation
- T024: Implement result display with staggered fade-in
- T025: Add save to album functionality 