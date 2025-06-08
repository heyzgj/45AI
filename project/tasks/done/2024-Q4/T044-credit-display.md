---
id: T044
title: Create CreditDisplay component
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

Create a reusable CreditDisplay component that shows the user's current credit balance with the 🎞️ icon. Should support different sizes and display modes (compact, full).

### Acceptance Criteria

- [ ] Component shows credit balance with 🎞️ icon
- [ ] Supports different sizes (small, medium, large)
- [ ] Animated number changes
- [ ] Loading state while fetching
- [ ] Low balance warning state
- [ ] Click to navigate to purchase page

### Context Binding

- **Docs:** `@/docs/STYLE_GUIDE.md`
- **Code:** `@/frontend/src/components/`
- **Code:** `@/frontend/src/stores/user.ts`

### Agent Notes

Use the dusty rose color theme. Consider adding a subtle pulse animation when credits are low (<20). 