# FluxMap
FluxMap is a task and project management application with a Go backend and a React frontend. It includes PERT/CPM analytics, project invitations/roles, and notification/update support.
Installation

Prerequisites
- Go 1.20+ (as specified in go.mod)
- Node 16+ and npm (for frontend)
- PostgreSQL (or use the provided Docker Compose which includes Postgres)
- Docker & docker-compose (recommended for quick local setup)
- psql CLI (for manual DB operations)

1) Quickstart with Docker Compose (recommended)
- From the repository root:
```sh
docker-compose up --build
```
This will build and start backend, frontend (if configured), and a PostgreSQL container according to docker-compose.yml. The default ports are:
- Backend API: :4000
- Frontend dev server: :3000
- PostgreSQL: :5432

2) Manual local setup

a. Database
- Create a PostgreSQL database and run the schema:
```sh
psql -d <your_db_name> -f schema.sql
```
- Ensure your environment variables (see next section) point to the created DB.

b. Backend
- Install dependencies and run:
```sh
go mod download
go run ./api/web/main.go
```
- Or build a binary:
```sh
go build -o fluxmap .
./fluxmap
```

c. Frontend
```sh
cd client
npm install
npm start
# Open http://localhost:8000
```

Environment variables
Provide required environment variables to the backend (example names from codebase):
- SESSION_SECRET — secret for cookie session store
- DATABASE_URL or individual DB connection params (host, port, user, pass, dbname)
- Any other env keys used in /internal/env or main entry points

Check api/web/main.go and internal/env for the exact environment key names expected by the application.

Running tests & simple load test
- Frontend: run tests from client:
```sh
cd client
npm test
```
- Simple load tester included: rate_test.sh (use at your own risk against local instances).

Project structure (high level)
- api/ — Go HTTP server, handlers, middleware, WebSocket handlers
  - api/web — route registration, handlers, helpers, websockets
- client/ — React app + components
- internal/ — models, repository interfaces and shared types
- docker/ — Dockerfiles and nginx configs
- schema.sql — PostgreSQL schema and migrations
- rate_test.sh — simple load script

API documentation
The definitive list of available HTTP routes and their attached handlers / middleware lives in:
- api/web/routes.go

Refer to that file for a machine‑readable mapping of endpoints, required HTTP methods, and which middleware (authentication, authorization, rate limiting) is applied. The file shows routes for:
- Auth (login, register, session, logout)
- Projects (list, create, read, update, remove user)
- Invitations (invite, list, confirm)
- Tasks (create, assign, approve, update, delete)
- PERT / CPM endpoints
- Graphs and analytic endpoints
- WebSocket endpoint

For handler details and business logic, see:
- api/web/handlers.go
- api/web/helpers.go
- api/web/graphs_handler.go
- api/web/ws_handler.go
- api/web/sockets.go

Contributing
- Follow existing patterns: add new handlers under api/web, repository implementations under internal/repository, and model logic under internal/models.
- Keep schema changes in schema.sql or use a migration strategy if you add a migration system.
- Add tests when feasible and document new endpoints by updating api/web/routes.go and related handler comments.

License
Add an appropriate LICENSE file to the repository.

Contact / support
Open issues or pull requests in this repository for bugs, feature requests, or deployment help.
