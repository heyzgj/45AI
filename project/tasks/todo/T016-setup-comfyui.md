---
id: T016
title: Set up ComfyUI API on GCP with GPU
status: To Do
epic: E03
effort: L
risk: H
dependencies: []
assignee: CursorAgent
created_at: 2024-05-23
updated_at: 2024-05-23
---

### Description

Set up the ComfyUI API on a Google Cloud Platform (GCP) virtual machine with GPU support. This involves creating a new GCP project, configuring a VM instance with the appropriate machine type and GPU, installing the necessary drivers and dependencies, and deploying the ComfyUI application.

### Acceptance Criteria

- [ ] A new GCP project is created.
- [ ] A new VM instance is created with a GPU attached.
- [ ] The NVIDIA drivers are installed and configured correctly.
- [ ] The ComfyUI application is installed and running.
- [ ] The ComfyUI API is accessible from the backend application.

### Context Binding

- **Docs:** `@/docs/TECH_SPEC.md#AI-Service`, `@/infra/comfyui/README.md`
- **Code:** `@/infra/comfyui/`

### Agent Notes

*This is a complex infrastructure task that may require manual steps. The user's GCP credentials will be required.* 