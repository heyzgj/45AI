
```markdown filename: docs/ARCHITECTURE_GUIDE.md
# Architecture Guide

## 1. High-Level Architecture
This application is built on a modern client-server model with a clear separation of concerns, designed for scalability and maintainability.

- **Frontend:** A single-page application built with UniBest (uni-app). It handles all rendering, user interaction, and client-side state. Its primary responsibility is crafting a premium user experience, implementing the fluid animations, soft visual styles, and micro-interactions defined in the Style Guide.
- **Backend:** A monolithic Go API server that acts as the "brain" of the application. It follows a layered architecture (Controller -> Service -> Repository) to manage business logic. It serves as a secure gateway, orchestrating calls to the database and external services, ensuring the client doesn't interact with them directly.
- **Infrastructure:** The primary infrastructure is WeChat Cloud Hosting, which provides a managed environment for the Go backend and MySQL database, simplifying deployment and integration with WeChat services. The computationally-intensive AI workload is intentionally decoupled and runs on a specialized GCP environment with GPU access, allowing each part of the system to be scaled independently based on its specific needs.

## 2. Notable Design Patterns
- **Repository Pattern:** The backend uses repositories to abstract data access. The services interact with repository interfaces, not concrete database implementations, making the code easier to test and maintain.
- **API Gateway (Backend as a Gateway):** The Go backend acts as a single entry point for the client. This pattern simplifies the client application and centralizes concerns like authentication, rate limiting, and logging.
- **Stateless Services:** Both the Go backend (using JWT for auth) and the ComfyUI API are stateless, which is crucial for horizontal scaling. Any number of instances can be run behind a load balancer.

## 3. Codebase Organization
/
├── frontend/ # UniBest (uni-app) source code
│ ├── src/
│ │ ├── pages/ # Application pages (Home, Profile, etc.)
│ │ ├── components/ # Reusable UI components (TemplateCard, AnimatedLoader)
│ │ ├── styles/ # Global styles, variables from Style Guide
│ │ ├── store/ # State management (Pinia or Vuex)
│ │ └── api/ # API client for communicating with the backend
│
├── backend/ # Go API server source code
│ ├── cmd/api/ # Main application entrypoint
│ ├── internal/
│ │ ├── handler/ # HTTP handlers (Controllers)
│ │ ├── service/ # Business logic
│ │ ├── repository/ # Database interaction logic
│ │ └── model/ # Data structures
│ └── go.mod
│
└── infra/ # Infrastructure as Code (optional)
└── comfyui/ # Dockerfile and configs for the GCP AI service


## 4. How to Extend
Adding a new feature (e.g., "AI Wardrobe") generally follows this workflow:

1.  **PRD & Design:** Update `PRD.md` and `UX_FLOW.md` with the new feature requirements. Consult the `STYLE_GUIDE.md` to design a visually consistent experience.
2.  **Data Model:** If needed, add a new table migration and define the model in `backend/internal/model/`. Update `DATA_MAP.md`.
3.  **Backend Logic:**
    - Create a new repository in `backend/internal/repository/`.
    - Implement the business logic in a new service in `backend/internal/service/`.
    - Expose the functionality via a new handler and route in `backend/internal/handler/`.
4.  **Frontend UI:**
    - Create new page(s) in `frontend/src/pages/`.
    - Build new components in `frontend/src/components/`, ensuring they follow all Style Guide rules for aesthetics and motion.
    - Add API calls in `frontend/src/api/` to connect to the new backend endpoint.
5.  **Testing:** Add unit, integration, and E2E tests, paying special attention to testing the new animations and UI states. Update `TEST_PLAN.md`.

## 5. Known Trade-offs and Future Improvements
- **Monolith vs. Microservices:** The backend is a monolith for simplicity. As features grow, services like billing or user management could be extracted.
- **No Caching Layer:** There is no caching layer (like Redis) for data like templates. This was a trade-off for faster initial development. If the template API comes under heavy load, implementing a Redis cache would be a key optimization.
- **Vendor Lock-in:** Using WeChat Cloud Hosting creates a degree of vendor lock-in but is a deliberate trade-off for speed-to-market and simplified integration with the primary target ecosystem.