---
id: T032
title: Build credit purchase UI with pack selection
status: Done
epic: E04
effort: M
risk: L
dependencies: [T029, T030]
assignee: CursorAgent
---

### Description

Create the frontend UI for purchasing credits. This page will display a list of available credit packs and allow the user to select one to purchase.

### Acceptance Criteria

- [x] A new page is created at `/frontend/src/pages/purchase/index.vue`.
- [x] The page displays a list of credit packs with clear pricing.
- [x] When a user selects a pack, it initiates the appropriate payment flow (WeChat Pay or Apple IAP).
- [x] The UI adheres to the `STYLE_GUIDE.md`, focusing on a clean and trustworthy design.

### Context Binding

- **Docs:** `@/docs/UX_FLOW.md#Flow-3-Credit-Recharge`
- **Code:** `@/frontend/src/pages/purchase/index.vue`

### Agent Notes

*The design of this page is critical for user trust. It should be clean, simple, and transparent.* 