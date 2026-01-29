# Beheer: Go Fullstack Template

**Tech Stack:** Go (Golang), PostgreSQL, SQLC, Goose.

> **⚠️ CRITICAL RULES:**
> Before writing any code, you MUST consult **[docs/rules.md](rules.md)**.
> This contains the "Living Laws" of the project (Design Patterns, Strict Style).

## 📚 Documentation Index

| Category         | Description                                     | Link                                            |
| :--------------- | :---------------------------------------------- | :---------------------------------------------- |
| **Architecture** | Project structure, DDD concepts, entry points.  | [View Structure](architecture/structure.md)     |
| **Persistence**  | Database workflow, SQLC, Migrations.            | [View Persistence](architecture/persistence.md) |
| **API**          | REST conventions, JSON envelopes, File uploads. | [View API](api/conventions.md)                  |
| **Coding**       | Interfaces, Config, Boilerplate reduction.      | [View Standards](coding/standards.md)           |
| **Schema**       | ER Diagram and Entities.                        | [View Database](DATABASE.md)                    |

## 🚀 Quick Start

1.  **Entry Point:** `cmd/api/main.go`
2.  **Config:** Set environment variables.
3.  **Run:** `go run cmd/api/main.go`
4.  **Watch (Hot Reload):** `make watch` (Requires [Air](github.com/air-verse/air))

## 🛠️ Features

- **IAM Module:** Implemented Users, Organizations, and OrganizationUsers (Memberships).
- **Authentication:** Middleware and handlers for secure access.
- **Database Schema:** Foundation migrations for UUIDs and Auth/IAM.
- **Architecture:** Clean DDD-inspired structure ready for expansion.
