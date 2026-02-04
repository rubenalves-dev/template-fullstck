# 📏 Project Rules & Design Patterns

This document serves as the **Single Source of Truth** for:
1.  **Code Style & Conventions**
2.  **Design Patterns** (e.g., State Pattern, Factory Pattern)
3.  **Architectural Rules**

**Note to Agent:** Every time a new design pattern or strict code style is established during the conversation, it MUST be registered here.

## 1. Design Patterns

### 1.1 State Management
- **Pattern:** State Pattern via Dates (not Enums).
- **Rule:** Do not use `ENUM` for entity states (e.g., `Active`, `Inactive`). Instead, use nullable timestamps (e.g., `activated_at`, `archived_at`, `canceled_at`).
- **Reasoning:** Provides historical context and allows for soft-deletes/restoration without losing data.

### 1.2 Line Item Association (Invoices)
- **Pattern:** Link `InvoiceItem` to the specific entity (Enrollment, Subscription, Product) via nullable Foreign Keys.
- **Rule:** Do not link the `Invoice` header directly to the consumable. Link the *Line Item*.

## 2. Architectural Rules

### 2.1 Interface-First Development
- **Rule:** Prioritize interfaces over specific implementations.
- **Mandate:** Handlers MUST depend on Service Interfaces. Services MUST depend on Repository Interfaces.
- **Reasoning:** Ensures high decoupling, allows for easy mocking in tests, and enables swapping implementations (e.g., SQL repository for a Mock or API client) without touching business logic.

### 2.2 Folder Structure & Layering
- **Module Root:** `internal/<module>/` (e.g., `internal/auth/`, `internal/cms/`).
- **Domain:** `internal/<module>/domain/`. Contains entities, DTOs, and repository/service interfaces.
- **Service:** `internal/<module>/service/` or `services/`. Contains business logic and event publishing.
- **Repositories:** `internal/<module>/repositories/`. Contains concrete persistence implementations using `pgx`.
- **Delivery:**
    - **HTTP Handlers:** `internal/<module>/delivery/http/`.
    - **Event Listeners:** `internal/<module>/delivery/events/`.
- **Platform:** `internal/platform/`. Shared infrastructure (DB, NATS, Router, Config).

