# 🏗️ Project Structure & Architecture

## 1. Architectural Philosophy

The project follows an **Event-Driven Architecture (EDA)**. Microservices themselves use a **simplified version of Domain-Driven Design (DDD)** to maintain clear boundaries while reducing boilerplate. We prioritize **interfaces as types** to ensure decoupling and use **NATS** for inter-service communication.

## 2. Standard Folder Structure

We adhere to the [Standard Go Project Layout](https://github.com/golang-standards/project-layout), adapted for DDD.

```text
.
├── cmd/
│   └── api/
│       └── main.go            # Entry Point: Dependency injection & server startup
├── internal/
│   ├── auth/                  # Authentication & Identity Module
│   │   ├── delivery/          # HTTP Handlers & Events
│   │   ├── domain/            # Domain Entities & Interfaces
│   │   ├── repositories/      # Persistence implementation
│   │   └── service/           # Business Logic
│   ├── cms/                   # Content Management System Module
│   │   ├── delivery/          # HTTP Handlers & Events
│   │   ├── domain/            # Domain Entities, DTOs & Interfaces
│   │   ├── repositories/      # Persistence implementation
│   │   └── services/          # Business Logic
│   ├── platform/              # Infrastructure (DB, NATS, Config)
│   └── platform/              # Shared infrastructure (DB, NATS, Router)
├── migrations/                # Database migrations (Goose)
├── pkg/                       # Shared libraries (jsonutil, httputil)
├── scripts/                   # Utility scripts
├── Makefile                   # Build & Dev commands
└── go.mod                     # Go module definition
```

### Key Rules

1.  **Entry Point:** The application entry point MUST be `cmd/api/main.go`.
2.  **Module Isolation:** Each business module (auth, cms) is self-contained under `internal/`.
3.  **Layered Architecture:** Each module follows a simplified layered structure:
    - **Domain:** Entities, DTOs, and repository/service interfaces.
    - **Service:** Business logic implementation and event publishing.
    - **Repositories:** Data access implementation.
    - **Delivery:** External interfaces (HTTP handlers and NATS event listeners).
4.  **Event-Driven Communication:** Modules communicate asynchronously using NATS. Services publish events (e.g., `cms.page.published`) that other modules can subscribe to.
5.  **Platform Layer:** Cross-cutting concerns like database connections, NATS, and configuration reside in `internal/platform`.
6.  **Interface-First:** High-level components depend on interfaces defined in the Domain layer, not on concrete implementations.
7.  **Separation of Concerns:** HTTP handlers manage request/response, services manage logic, and repositories manage data.