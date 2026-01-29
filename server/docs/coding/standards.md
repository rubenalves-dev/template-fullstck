# 💻 Coding Standards

## 1. Interface-Driven Design
> **Rule:** Implement against interfaces, not structs.

- **Handlers (Ports)** depend on `Application Use Case Interfaces`.
- **Application Use Cases** depend on `Domain Service Interfaces` or `Repository Interfaces`.
- **Domain Services** depend on `Repository Interfaces`.
- **Goal:** Easy mocking and implementation swapping.

## 2. Configuration
- **Source:** Environment variables ONLY.
- **Loading:** Parse into a typesafe `Config` struct at startup (e.g., using `caarlos0/env`).

## 3. Boilerplate & Utils
- Use `restutil` for HTTP input/output (binding JSON, rendering responses).
- Use `errutil` for mapping domain errors to HTTP status codes.
- Do not repeat `json.NewDecoder` in every handler.

## 4. Logging
- Use `log/slog`.
- Always include `request_id` from the context in logs.
