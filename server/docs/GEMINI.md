# Template Fullstack: Go Fullstack Template

**Tech Stack:** Go (Golang), PostgreSQL, pgx, Goose, NATS.

> **⚠️ CRITICAL RULES:**
> Before writing any code, you MUST consult **[docs/RULES.md](RULES.md)**.
> This contains the "Living Laws" of the project (Design Patterns, Strict Style).

## 📚 Documentation Index

| Category         | Description                                     | Link                                            |
| :--------------- | :---------------------------------------------- | :---------------------------------------------- |
| **Architecture** | Project structure, EDA, simplified DDD.         | [View Structure](architecture/structure.md)     |
| **Persistence**  | Database workflow, pgx, Migrations.             | [View Persistence](architecture/persistence.md) |
| **API**          | REST conventions, JSON envelopes, File uploads. | [View API](api/conventions.md)                  |
| **Coding**       | Interfaces, Config, Boilerplate reduction.      | [View Standards](coding/standards.md)           |
| **Schema**       | ER Diagram and Entities.                        | [View Database](DATABASE.md)                    |

## 🚀 Quick Start

1.  **Entry Point:** `cmd/api/main.go`
2.  **Config:** Set environment variables (DB, NATS, etc).
3.  **Run:** `go run cmd/api/main.go`
4.  **Watch (Hot Reload):** `make watch` (Requires [Air](https://github.com/air-verse/air))

## 🛠️ Features

- **Auth Module:** Users authentication and registration.
- **CMS Module:** Content Management with hierarchical layout (Pages, Rows, Columns, Blocks).
- **Event-Driven:** Inter-service communication via NATS.
- **Architecture:** Simplified DDD for clean and scalable codebase.
