---
id: T031
title: Create purchase validation endpoint
status: To Do
epic: E04
effort: M
risk: H
dependencies: [T029, T030]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Create a backend API endpoint (`POST /api/v1/billing/purchase`) to validate purchase receipts from both WeChat Pay and Apple IAP. This endpoint will be the single point of contact for the frontend to confirm purchases.

### Acceptance Criteria

- [ ] A `POST /api/v1/billing/purchase` endpoint is created.
- [ ] The endpoint accepts a `productId` and a platform-specific `receipt`.
- [ ] It calls the appropriate service (`WeChatPayService` or `AppleIAPService`) to validate the receipt.
- [ ] On successful validation, it updates the user's credit balance using the `CreditService`.
- [ ] It returns a confirmation to the frontend.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#API-Design`
- **Code:** `@/backend/internal/handler/billing_handler.go`

### Agent Notes

*This endpoint needs to be robust and secure, as it directly affects the monetization of the application.* 