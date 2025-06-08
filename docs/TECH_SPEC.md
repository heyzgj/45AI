# Technical Specification

## 1. Overview
The system employs a client-server architecture with a decoupled AI processing service. The frontend is a cross-platform UniBest application communicating with a Go backend API. The backend manages users, templates, and transactions, while orchestrating calls to external services for content moderation and AI image generation. This design separates business logic from intensive computation.

## 2. Architecture Diagram
`UniBest App (iOS/WeChat)` <--> `Go Backend API (WeChat Cloud)`
          |                                  |
          |                                  +--> `MySQL 8.0` (for user/transaction data)
          |                                  |
          +----------------------------------+--> `Tencent Cloud Content Safety API` (on upload)
                                             |
                                             +--> `ComfyUI API (GCP)` (for image generation)

## 3. Components
- **Frontend:** A UniBest (uni-app) project, responsible for all UI, state management, and client-side interactions. It communicates with the backend via a RESTful API and is responsible for implementing the sophisticated animations and transitions defined in the Style Guide. It uses the `wot-design-uni` component library to ensure a consistent and high-quality user interface.
- **Backend:** A Go API server using the Gin framework, deployed on WeChat Cloud Hosting. Its responsibilities include:
  - User authentication (JWT-based).
  - Serving template metadata.
  - Processing credit transactions and IAP validation.
  - Acting as a proxy/orchestrator for external API calls.
- **Database:** A MySQL 8.0 instance on WeChat Cloud Hosting. It stores user profiles, AI templates metadata, and transaction logs.
- **AI Service:** A self-hosted ComfyUI instance running on Google Cloud Platform (GCP) with GPU support. It exposes a private API endpoint that the Go backend calls to perform the image generation task.
- **External Services:**
  - **Tencent Cloud Content Safety:** Used for synchronous moderation of all user image uploads.
  - **WeChat/Apple IAP:** Platform-native services for processing payments.

## 4. Data Models
- **User:** `id, wechat_openid, nickname, avatar_url, credits, created_at, updated_at`
- **Template:** `id, name, description, preview_image_url, credit_cost, is_active, created_at`
- **Transaction:** `id, user_id (FK), type ('purchase' or 'generation'), amount (credits), external_payment_id, related_template_id, created_at`

## 5. API Design
- **Auth**
  - `POST /api/v1/auth/login`: Takes WeChat `code`, returns a JWT token.
- **Templates**
  - `GET /api/v1/templates`: Returns a list of all active templates.
- **Generation**
  - `POST /api/v1/generate`:
    - Requires JWT auth.
    - `multipart/form-data` request with `image` file and `template_id`.
    - Backend validates image with Tencent API, checks user credits, deducts credits, calls ComfyUI API, and returns an array of generated image URLs.
- **Billing**
  - `GET /api/v1/me/transactions`: Returns the user's transaction history.
  - `POST /api/v1/billing/purchase`: Takes a `productId` and platform `receipt`, validates it, and updates user credits.

## 6. Key Technical Decisions
- **Language & Framework (Backend):** Go with Gin. Chosen for its high performance, concurrency, and low resource footprint, which is ideal for an API gateway that orchestrates multiple services.
- **Framework (Frontend):** UniBest. Chosen to meet the requirement of a single codebase for both WeChat Mini Program and iOS. It utilizes the `wot-design-uni` component library.
- **Auth Method:** JWT (JSON Web Tokens). Stateless authentication is a good fit for the microservice-oriented architecture and simplifies scaling the backend.
- **AI Deployment:** Decoupled service on GCP. Separating the GPU-intensive ComfyUI workload from the main business logic API allows for independent scaling and resource management. GCP is chosen for its robust GPU offerings.
- **Deployment:** WeChat Cloud Hosting for the backend/DB. This simplifies integration with WeChat login and payment services.

## 7. Constraints & Non-Functional Requirements
- **Performance:** P95 latency for API responses (excluding generation) must be under 200ms. End-to-end generation time should be < 30 seconds.
- **Security:** All user-uploaded content must be moderated before processing. Passwords or sensitive keys are not stored. All communication is over HTTPS.
- **Scalability:** The Go backend must be horizontally scalable. The GCP AI service can be scaled by adding more GPU instances.