# 💾 Persistence & Database

## 1. Workflow (Goose + Repository Pattern)
1.  **Migrations:** Define schema changes in `/migrations` using **Goose**.
2.  **Implementation:** Create repository interfaces in the module's `domain` layer and implement them in the `repositories` layer using **pgx**.
3.  **Usage:** Repositories are injected into services to handle data persistence.

## 2. Database Conventions
- **Primary Keys:** MUST use `uuid.UUID` (`github.com/google/uuid`).
- **Timestamps:**
  - `created_at` (Default `NOW()`)
  - `updated_at` (On Update `NOW()`)
- **JSONB:** Used for flexible content structures in the CMS (block content, background configs).

## 3. Schema Documentation
Refer to [docs/DATABASE.md](../DATABASE.md) for the ER Diagram and specific table definitions.
