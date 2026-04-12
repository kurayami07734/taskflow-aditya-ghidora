# Taskflow API

## 1. Overview

Taskflow is a RESTful API for managing projects and tasks. Users can create projects, manage tasks within projects, and organize work with priorities and status tracking.

### Tech Stack

- **Language**: Go 1.25
- **Framework**: chi (HTTP router)
- **Database**: PostgreSQL 15
- **ORM**: sqlx
- **Authentication**: JWT (HS256)
- **Documentation**: Swagger/OpenAPI (swaggo)
- **Dev**: Air (live reload)
- **Testing**: testcontainers-go

## 2. Architecture Decisions

### Structure

The project follows a standard Go layout with clear separation of concerns:

- `handlers/` - HTTP request handlers (controllers)
- `models/` - Database models and queries
- `routers/` - Route definitions
- `middleware/` - HTTP middleware (auth, logging)
- `utils/` - Utilities (config, logger, JWT, pagination, responses)
- `migrations/` - Database migrations
- `seeds/` - Seed data for testing

### Decisions Made

1. **No separate service layer** - Handlers directly call model functions. This works for a simple API but would need a service layer for complex business logic.

2. **Manual JWT claims** - Used custom JWT functions instead of jwt-go library for transparency and to avoid dependencies.

3. **Typed error responses** - Created typed error structs (BadRequest, Unauthorized, Forbidden, NotFound, Internal) with consistent JSON responses.

4. **Pagination** - Added default pagination (page=1, limit=10, max=30) to list endpoints.

5. **Slog logging** - Replaced default chi logger with custom slog-based middleware for better structured logging.

### What Was Intentionally Left Out

- **Rate limiting** - Not needed for MVP
- **WebSocket** - Not needed for MVP
- **Caching** - Not needed for MVP
- **Separate service/repository layers** - Overhead for current scope
- **Graceful shutdown** - Implemented in main.go

## 3. Running Locally

```bash
# Clone the repository
git clone https://github.com/anomalyco/taskflow-aditya-ghidora
cd taskflow-aditya-ghidora

# Copy environment file
cp .env.example .env

# Start the application
docker compose up
# API available at http://localhost:8080 (default port)
```

The first run will execute migrations and seed data automatically.

## 4. Running Migrations

Migrations run automatically on startup via Docker Compose. If you need to run them manually:

```bash
docker compose run migrate
```

## 5. Test Credentials

The database is seeded with a test user:

- **Email**: test@example.com
- **Password**: password123

## 6. API Reference

Full API documentation (swagger) is available at http://localhost:8080/docs/

## 7. What You'd Do With More Time

### Shortcuts Taken

1. **No service layer** - Business logic mixed with handlers
2. **Limited validation** - Basic validation, no complex rules
3. **No transactions** - Some operations could benefit from atomic transactions
4. **Manual JWT** - Used custom implementation instead of standard library
5. **Limited error handling** - Could have more granular error codes

### Improvements

1. **Add service layer** - Better separation of concerns
2. **Unit tests** - Currently only integration tests
3. **More middleware** - Request ID, correlation ID, tracing
4. **Database transactions** - For operations that modify multiple tables
5. **Caching** - Redis for frequently accessed data
6. **Rate limiting** - Prevent abuse
7. **Webhooks** - For external integrations
8. **CI/CD** - Setup pipelines to automate tests and linting checks on PRs

## Development Commands

```bash
# start database for local dev
docker compose -f compose.local.yml up

# Install dependencies
cd backend && go mod tidy

# Run tests
make test

# Generate swagger docs
make docs

# Run with hot reload
make dev
```
