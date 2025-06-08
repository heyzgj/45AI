---
id: T041
title: Set up API client and HTTP interceptors
status: Done
epic: E05
effort: M
risk: Low
dependencies: [T037]
assignee: CursorAgent
created_at: 2024-12-24
updated_at: 2024-12-24
---

### Description

Create the API client module using uni-app's network request APIs with proper interceptors for authentication, error handling, and loading states. Configure base URLs and request/response transformations.

### Acceptance Criteria

- [ ] API client module created with TypeScript support
- [ ] Request interceptor adds JWT token from store
- [ ] Response interceptor handles common errors
- [ ] Loading state management integrated
- [ ] Timeout and retry logic implemented
- [ ] Environment-based base URL configuration

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#API-Communication`
- **Code:** `@/frontend/src/api/`
- **Code:** `@/frontend/src/utils/request.ts`

### Agent Notes

Use uni.request as the base HTTP client. Ensure proper TypeScript typing for all API methods. Handle token refresh logic. 