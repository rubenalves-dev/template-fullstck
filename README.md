# Go Fullstack Template

A modern, production-ready Go fullstack template using Clean Architecture and Domain-Driven Design (DDD) principles.

## 🚀 Overview

This repository serves as a foundation for building scalable SaaS applications. It comes pre-configured with essential IAM (Identity and Access Management) features and a robust architectural pattern.

## 🛠️ Tech Stack

- **Backend:** Go (Golang)
- **Database:** PostgreSQL
- **SQL Generation:** [SQLC](https://sqlc.dev/)
- **Migrations:** [Goose](https://github.com/pressly/goose)
- **Router:** [go-chi](https://github.com/go-chi/chi)
- **Validation:** [caarlos0/env](https://github.com/caarlos0/env) for config

## 🏗️ Architecture

The project follows a DDD-inspired structure with clearly defined layers:

1.  **Ports (Entry Points):** HTTP handlers and middleware.
2.  **Applications (Use Cases):** Orchestrates domain logic and ties different domains together.
3.  **Domain:** Core business logic, entities, and repository interfaces.
4.  **Infrastructure:** Implementation details like persistence (Postgres), external APIs, etc.

For more details, see [server/docs/architecture/structure.md](server/docs/architecture/structure.md).

## 🚦 Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- [SQLC](https://sqlc.dev/docs/install/)
- [Air](https://github.com/air-verse/air) (for hot-reload)

### Setup

1.  Clone the repository.
2.  Navigate to the `server` directory.
3.  Copy `.env.example` to `.env`.
4.  Run `docker-compose up -d` to start the database.
5.  Run `make watch` to start the development server.

## 📚 Documentation

Comprehensive documentation is available in the `server/docs` directory:

- [Project Overview](server/docs/GEMINI.md)
- [Database Schema](server/docs/DATABASE.md)
- [Architecture](server/docs/architecture/structure.md)
- [Coding Standards](server/docs/coding/standards.md)
