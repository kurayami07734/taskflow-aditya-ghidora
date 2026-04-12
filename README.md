# Taskflow

Taskflow is a full-stack project and task management application. The backend is a RESTful API, and the frontend is a React-based SPA.

## Table of Contents

1. [Tech Stack](#1-tech-stack)
2. [Architecture](#2-architecture)
3. [Running with Docker Compose](#3-running-with-docker-compose)
4. [Local Development Setup](#4-local-development-setup)
5. [Test Credentials](#5-test-credentials)
6. [API Reference](#6-api-reference)
7. [What You'd Do With More Time](#7-what-youd-do-with-more-time)

---

## 1. Tech Stack

### Backend (Go)

- **Language**: Go 1.25
- **Framework**: chi (HTTP router)
- **Database**: PostgreSQL 15
- **ORM**: sqlx
- **Authentication**: JWT (HS256)
- **Documentation**: Swagger/OpenAPI (swaggo)
- **Dev**: Air (live reload)
- **Testing**: testcontainers-go

### Frontend (React)

- **Language**: TypeScript
- **Framework**: React 19 with Vite
- **UI Library**: Material UI (MUI) v9
- **State Management**: Zustand
- **Routing**: React Router v7
- **Data Fetching**: TanStack Query (React Query)
- **Date Handling**: Day.js with MUI Date Pickers

---

## 2. Architecture

### Backend Architecture

The backend follows a standard Go layout with clear separation of concerns:

- `handlers/` - HTTP request handlers (controllers)
- `models/` - Database models and queries
- `routers/` - Route definitions
- `middleware/` - HTTP middleware (auth, logging)
- `utils/` - Utilities (config, logger, JWT, pagination, responses)
- `migrations/` - Database migrations
- `seeds/` - Seed data for testing

#### Decisions Made

1. **No separate service layer** - Handlers directly call model functions. This works for a simple API but would need a service layer for complex business logic.

2. **Manual JWT claims** - Used custom JWT functions instead of jwt-go library for transparency and to avoid dependencies.

3. **Typed error responses** - Created typed error structs (BadRequest, Unauthorized, Forbidden, NotFound, Internal) with consistent JSON responses.

4. **Pagination** - Added default pagination (page=1, limit=10, max=30) to list endpoints.

5. **Slog logging** - Replaced default chi logger with custom slog-based middleware for better structured logging.

#### What Was Intentionally Left Out

- **Rate limiting** - Not needed for MVP
- **WebSocket** - Not needed for MVP
- **Caching** - Not needed for MVP
- **Separate service/repository layers** - Overhead for current scope
- **Graceful shutdown** - Implemented in main.go

### Frontend Architecture

The frontend follows a feature-based folder structure:

```
frontend/src/
├── api/              # API client and type definitions
├── components/       # Reusable components
│   ├── common/       # Shared components (Snackbar, Dialogs)
│   └── layout/       # Layout components (Navbar, AppLayout)
├── features/         # Feature-specific components
│   ├── projects/     # Project-related components
│   └── tasks/        # Task-related components
├── pages/            # Route page components
├── stores/           # Zustand state stores
├── theme/            # MUI theme configuration
└── types/            # Additional TypeScript types
```

#### State Management

- **Zustand** for global state (auth, theme)
- **TanStack Query** for server state (API data caching, mutations)

#### UI Pattern

- Material UI components with custom theme
- Dark/Light mode toggle with localStorage persistence
- Consistent spacing and typography using MUI's theme system

---

## 3. Running with Docker Compose

### Prerequisites

- Docker and Docker Compose

### Full Stack (API + Database + Frontend)

```bash
# Clone the repository
git clone https://github.com/anomalyco/taskflow-aditya-ghidora
cd taskflow-aditya-ghidora

# Copy environment file
cp .env.example .env

# Start all services
docker compose up
```

This starts:
- **PostgreSQL** database on port 5432
- **Backend API** on port 8080 (configurable via PORT)
- **Frontend** served via nginx on port 80

- **API**: http://localhost:8080
- **API Docs**: http://localhost:8080/docs/
- **Frontend**: http://localhost

The first run will execute migrations and seed data automatically.

---

## 4. Local Development Setup

### Prerequisites

- Docker and Docker Compose (for database)
- Go 1.25+
- Node.js 18+

### Backend (Go) with Hot Reload

```bash
# Start database only (runs migrations and seeds automatically)
docker compose -f compose.local.yml up

# In a new terminal, run the API with hot reload
cd backend
make dev
```

The API will be available at http://localhost:8080

### Frontend (React)

```bash
# Install dependencies
cd frontend
npm install

# Start development server
npm run dev
```

The frontend will be available at http://localhost:5173

### Development Commands

#### Backend

```bash
# Install dependencies
cd backend && go mod tidy

# Run tests
make test

# Generate swagger docs
make docs

# Run with hot reload
make dev
```

#### Frontend

```bash
# Install dependencies
cd frontend && npm install

# Run development server
npm run dev

# Build for production
npm run build
```

---

## 5. Test Credentials

The database is seeded with a test user:

- **Email**: test@example.com
- **Password**: password123

---

## 6. API Reference

Full API documentation (Swagger) is available at http://localhost:8080/docs/

---

## 7. What You'd Do With More Time

### Backend Shortcuts Taken

1. **No service layer** - Business logic mixed with handlers
2. **Limited validation** - Basic validation, no complex rules
3. **No transactions** - Some operations could benefit from atomic transactions
4. **Manual JWT** - Used custom implementation instead of standard library
5. **Limited error handling** - Could have more granular error codes

### Backend Improvements

1. **Add service layer** - Better separation of concerns
2. **Unit tests** - Currently only integration tests
3. **More middleware** - Request ID, correlation ID, tracing
4. **Database transactions** - For operations that modify multiple tables
5. **Caching** - Redis for frequently accessed data
6. **Rate limiting** - Prevent abuse
7. **Webhooks** - For external integrations
8. **CI/CD** - Setup pipelines to automate tests and linting checks on PRs

### Frontend Shortcuts Taken

1. **No code splitting** - Large bundle size, could benefit from lazy loading
2. **Limited error boundaries** - No granular error handling for components
3. **No form validation library** - Using basic controlled inputs
4. **Hardcoded API base URL** - Could use environment variables more
5. **Limited loading states** - Basic loading indicators only

### Frontend Improvements

1. **Code splitting** - Use React.lazy() for route-based splitting
2. **Add React Hook Form** - Better form validation and handling
3. **Add error boundaries** - Graceful error handling per feature
4. **Skeleton loaders** - Better UX during data loading
5. **Add unit tests** - Jest + React Testing Library
6. **Accessibility audit** - Ensure WCAG compliance
7. **Add internationalization** - i18next for multi-language support
8. **Optimistic updates** - Better UX for mutations