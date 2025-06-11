---
id: T012
title: Implement template card component with hover animations
status: Done
epic: E02
effort: M
risk: Low
dependencies: []
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Create a reusable `TemplateCard` component that displays a template's preview image, name, and credit cost. The card should have hover and tap animations as defined in the `STYLE_GUIDE.md`.

### Acceptance Criteria

- [x] A `TemplateCard.vue` component is created in `/frontend/src/components/TemplateCard/`.
- [x] The component accepts a `template` object as a prop.
- [x] On hover, the card's preview image scales up and an overlay appears.
- [x] On tap, the card lifts up and forwards with a stronger shadow.
- [x] The component is responsive and adheres to the `STYLE_GUIDE.md`.

### Context Binding

- **Docs:** `@/docs/STYLE_GUIDE.md#Components`
- **Code:** `@/frontend/src/components/TemplateCard/TemplateCard.vue`

### Agent Notes

*This component was implemented as part of other tasks. Marking as done.* 