# 💾 Persistence & Database

## 1. Workflow (sqlc + goose)
1.  **Migrations:** Define schema changes in `internal/infrastructure/persistence/postgres/migrations` using **Goose**.
2.  **Queries:** Define SQL queries in `internal/infrastructure/persistence/postgres/queries/*.sql`.
3.  **Generation:** Run `sqlc generate`. This outputs Go code to `internal/infrastructure/persistence/postgres/sqlc`.
4.  **Usage:** Domain repositories (e.g., `internal/domain/iam/user/repo.go`) wrap the generated `sqlc` code to satisfy domain interfaces. These repositories are injected into Application Use Cases or Domain Services.

## 2. Database Conventions
- **Primary Keys:** MUST use `uuid.UUID` (`github.com/google/uuid`).
- **Timestamps:**
  - `created_at` (Default `NOW()`)
  - `updated_at` (On Update `NOW()`)
  - `deleted_at` (Optional, for soft deletes only)
- **Organization Isolation:** core tables MUST include `organization_id` for multi-tenancy.

## 3. Schema Documentation
Refer to [docs/DATABASE.md](../DATABASE.md) for the ER Diagram and specific table definitions.
