---
id: T011
title: Build template gallery UI with grid layout
status: Done
epic: E02
effort: M
risk: Low
dependencies: [T009, T043]
assignee: CursorAgent
created_at: 2024-05-22
updated_at: 2024-05-23
---

### Description

Create the main home page which functions as the template gallery. It should fetch the list of templates from the API and display them in a responsive grid.

### Acceptance Criteria

- [ ] The home page is located at `/frontend/src/pages/index/index.vue`.
- [ ] On page load, it calls the `GET /api/v1/templates` endpoint.
- [ ] It uses the `TemplateCard` component to display each template in a grid.
- [ ] The grid layout is responsive and adheres to the `STYLE_GUIDE.md` for spacing and aesthetics.
- [ ] A loading animation is shown while the templates are being fetched.

### Context Binding

- **Docs:** `@/docs/UX_FLOW.md#Flow-1-Core-Image-Generation`, `@/docs/STYLE_GUIDE.md`
- **Code:** `@/frontend/src/pages/index/index.vue`, `@/frontend/src/components/TemplateCard/TemplateCard.vue`, `@/frontend/src/api/template.js`

### Agent Notes

*Focus on the visual "airiness" and "softness" described in the style guide. Use appropriate padding and subtle entry animations for the cards.* 