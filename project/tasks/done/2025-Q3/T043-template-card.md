---
id: T043
title: Create TemplateCard component
status: Done
epic: E05
effort: S
risk: Low
dependencies: [T037, T042]
assignee: CursorAgent
created_at: 2024-12-24
updated_at: 2024-12-24
---

### Description

Create a reusable TemplateCard component that displays template preview, name, description, and credit cost. Should follow the Style Guide for animations and visual design.

### Acceptance Criteria

- [ ] Component displays template preview image
- [ ] Shows template name and description
- [ ] Displays credit cost with 🎞️ icon
- [ ] Hover/press effects as per Style Guide
- [ ] Loading and error states
- [ ] Emits click event for selection

### Context Binding

- **Docs:** `@/docs/STYLE_GUIDE.md#Component-Design`
- **Code:** `@/frontend/src/components/`
- **Type:** `@/frontend/src/types/api.ts#Template`

### Agent Notes

Follow the card styles already defined in components.scss. Ensure smooth animations on interaction. 