---
id: E04
title: Credit System & Monetization
status: To Do
owner: ProductManager
---

### Summary

This epic covers the credit ("胶卷") system including balance tracking, purchase flows, transaction history, and payment integration with WeChat Pay and Apple IAP. Users must be able to purchase credits and track their spending.

### Related Documents

- `@/docs/PRD.md#Feature-Credit-Monetization-System`
- `@/docs/DATA_MAP.md#transactions`
- `@/docs/TECH_SPEC.md#API-Design`
- `@/docs/UX_FLOW.md#Flow-3-Credit-Recharge`

### Tasks in this Epic

- T026: Create database migration for transactions table
- T027: Implement credit balance management service
- T028: Build transaction logging system
- T029: Integrate WeChat Pay SDK
- T030: Integrate Apple In-App Purchase
- T031: Create purchase validation endpoint
- T032: Build credit purchase UI with pack selection
- T033: Implement transaction history API
- T034: Create transaction history UI
- T035: Add credit balance display with animations
- T036: Implement insufficient credits modal 