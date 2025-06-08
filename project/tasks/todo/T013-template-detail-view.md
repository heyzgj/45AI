---
id: T013
title: Create template detail view/modal
status: To Do
epic: E02
effort: M
risk: Low
dependencies: [T010, T011]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-22
---

### Description

When a user taps on a `TemplateCard` from the gallery, a detail view should be presented. This could be a new page or a modal overlay. It will show a larger preview, the full description, and the "Generate" button.

### Acceptance Criteria

- [ ] Tapping a card in the gallery navigates to a detail page or opens a modal.
- [ ] The detail view fetches data from `GET /api/v1/templates/:id`.
- [ ] It displays a larger `preview_image_url`, the `name`, `description`, and `credit_cost`.
- [ ] A primary "Generate with this Template" button is prominent.
- [ ] The entry and exit animations for this view are fluid and follow the `STYLE_GUIDE.md`.

### Context Binding

- **Docs:** `@/docs/UX_FLOW.md#Flow-1-Core-Image-Generation`
- **Code:** `@/frontend/src/pages/template-detail/index.vue` (or as a component)

### Agent Notes

*A modal might be better to keep the user in the context of the gallery. Ensure the transition is smooth.* 