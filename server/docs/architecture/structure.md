# 🏗️ Project Structure & Architecture

## 1. Architectural Philosophy

The project is organized by **DDD** (Domain-Driven Design). We prioritize **interfaces as types** to ensure decoupling.

## 2. Standard Folder Structure

We adhere to the [Standard Go Project Layout](https://github.com/golang-standards/project-layout), adapted for DDD.

```text
.
├── cmd/
│   └── api/
│       └── main.go            # Entry Point: Dependency injection & server startup
├── internal/
│   ├── applications/          # Use Cases: Orchestrates domain logic and ties domains together
│   ├── config/                # Configuration loading (env vars)
│   ├── domain/                # Pure Domain Logic & Interfaces
│   │   ├── common/            # Shared value objects/types
│   │   └── iam/               # Identity & Access Management Bounded Context
│   │       ├── organization/
│   │       ├── organizationuser/
│   │       └── user/
│   ├── infrastructure/        # Infrastructure Implementation
│   │   └── persistence/
│   │       └── postgres/
│   │           ├── migrations/ # Goose migration files (.sql)
│   │           ├── queries/    # SQLC queries (.sql)
│   │           └── sqlc/       # Generated Go code from SQLC
│   ├── ports/                 # Driving Adapters (Entry Points)
│   │   └── http/              # Handlers, Router, Middleware
│   └── utils/                 # Utility functions
├── sqlc.yaml                  # SQLC configuration
├── Makefile                   # Build & Dev commands
└── go.mod                     # Go module definition
```

### Key Rules

1.  **Entry Point:** The application entry point MUST be `cmd/api/main.go`.
2.  **Concept Isolation:** Each domain concept has its own folder under `internal/domain`.
3.  **Application Layer:** Use cases reside in `internal/applications`. They are responsible for orchestrating domain services, managing transactions, and converting between DTOs and entities if necessary.
4.  **Persistence Layer:** All database-related code resides in `internal/infrastructure/persistence`.
5.  **Interface-First:** Handlers depend on Application Interfaces; Applications depend on Domain Service Interfaces; Domain Services depend on Repository Interfaces.
6.  **Separation of Concerns:** `ports/http` handles the web layer, `applications` handles use case orchestration, `domain` handles business logic, and `infrastructure` handles data access.