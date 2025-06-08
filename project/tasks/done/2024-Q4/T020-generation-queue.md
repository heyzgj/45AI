---
id: T020
title: Implement generation request queueing system
status: Done
epic: E03
effort: L
risk: H
dependencies: [T018]
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Implement a queueing system to manage image generation requests. This is crucial for handling concurrent requests and ensuring that the ComfyUI API is not overloaded. The system should use a message queue (e.g., RabbitMQ, Redis) to process requests asynchronously.

### Acceptance Criteria

- [ ] A message queue is integrated into the backend.
- [ ] When a generation request is received, it is added to the queue.
- [ ] A separate worker process consumes requests from the queue and calls the `GenerationService`.
- [ ] The API endpoint returns an initial response with a request ID.
- [ ] The frontend can use the request ID to poll for the generation status.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#Scalability`
- **Code:** `@/backend/internal/service/queue_service.go`, `@/backend/internal/worker/`

### Agent Notes

*This is a significant architectural change that will require careful planning and implementation. A simple in-memory queue can be used for initial development, but a more robust solution like RabbitMQ should be used for production.* 