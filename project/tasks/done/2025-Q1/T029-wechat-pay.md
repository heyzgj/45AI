---
id: T029
title: Integrate WeChat Pay SDK
status: Done
epic: E04
effort: L
risk: H
dependencies: [T027]
assignee: CursorAgent
---

### Description

Integrate the WeChat Pay SDK to handle credit purchases within the WeChat Mini Program. This will involve using the `uni-app` API for payments and a backend service to create and validate orders.

### Acceptance Criteria

- [x] The WeChat Pay SDK is added to the frontend project.
- [x] A backend service is created to generate a pre-pay order with WeChat.
- [x] The frontend calls this service to get a pre-pay ID and then initiates the payment with `uni.requestPayment`.
- [x] The payment result is sent back to the backend for validation.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#External-Services`
- **Code:** `@/frontend/src/api/payment.js`, `@/backend/internal/service/wechat_pay_service.go`

### Agent Notes

*This is a complex integration that requires careful handling of secrets and callbacks. A mock implementation should be created for initial development.* 