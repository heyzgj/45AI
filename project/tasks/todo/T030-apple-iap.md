---
id: T030
title: Integrate Apple In-App Purchase
status: To Do
epic: E04
effort: L
risk: H
dependencies: [T027]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Integrate Apple's In-App Purchase (IAP) framework to handle credit purchases on the iOS platform. This will involve using StoreKit on the frontend and a backend service to validate receipts.

### Acceptance Criteria

- [ ] The StoreKit framework is used on the iOS frontend to initiate purchases.
- [ ] A backend service is created to validate the App Store receipt.
- [ ] On successful validation, the user's credit balance is updated.
- [ ] A transaction is logged for the purchase.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#External-Services`
- **Code:** `@/frontend/src/api/payment.js`, `@/backend/internal/service/apple_iap_service.go`

### Agent Notes

*This is a complex integration that requires careful handling of secrets and callbacks. A mock implementation should be created for initial development.* 