---
id: T031
title: Create purchase validation endpoint
status: Done
epic: E04
effort: M
risk: H
dependencies: [T029, T030]
assignee: CursorAgent
---

### Description

Create a backend API endpoint (`POST /api/v1/billing/purchase`) to validate purchase receipts from both WeChat Pay and Apple IAP. This endpoint will be the single point of contact for the frontend to confirm purchases.

### Acceptance Criteria

- [x] A `POST /api/v1/billing/purchase` endpoint is created.
- [x] The endpoint accepts a `productId` and a platform-specific `receipt`.
- [x] It calls the appropriate service (`WeChatPayService` or `AppleIAPService`) to validate the receipt.
- [x] On successful validation, it updates the user's credit balance using the `CreditService`.
- [x] It returns a confirmation to the frontend.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#API-Design`
- **Code:** `@/backend/internal/handler/billing_handler.go`

### Agent Notes

*This endpoint needs to be robust and secure, as it directly affects the monetization of the application.* 