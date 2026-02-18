# Todo gRPC Service

A production-style Todo API built in Go to demonstrate clean architecture, gRPC with buf, dependency injection with Wire, Postgres with GORM, and database migrations — all wired together from scratch.

---

## Tech Stack

| Concern | Tool |
|---|---|
| API protocol | [gRPC](https://grpc.io/) |
| Proto management | [buf CLI](https://buf.build/docs/installation) |
| Database | PostgreSQL 16 |
| ORM | [GORM](https://gorm.io/) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Dependency injection | [Wire](https://github.com/google/wire) |
| Config | Environment variables via [godotenv](https://github.com/joho/godotenv) |
| Containerisation | Docker + Docker Compose |

---

## Architecture

```
Request
  └── gRPC Handler          (internal/handler/grpc)   — translate proto ↔ domain
        └── Usecase         (internal/usecase)         — business logic
              └── Repository interface (internal/domain) — abstraction boundary
                    └── Repository impl (internal/infrastructure/repository) — GORM
                              └── Postgres DB
```

### Layer rules

- **`domain/`** — pure Go structs and interfaces. Zero external imports. Every other layer depends on this; this depends on nothing.
- **`usecase/`** — business logic only. Talks to `domain.TodoRepository` (the interface), never to GORM directly.
- **`infrastructure/`** — GORM lives here and only here. Implements `domain.TodoRepository`. Converts between GORM models and domain entities at the boundary.
- **`handler/grpc/`** — translates gRPC request/response types to/from domain types. No business logic.
- **`wire/`** — wires all layers together at compile time. No runtime reflection.

---

## Prerequisites

| Tool | Install |
|---|---|
| Go 1.23+ | https://go.dev/dl |
| Docker + Docker Compose | https://docs.docker.com/get-docker |
| buf CLI | `brew install bufbuild/buf/buf` or https://buf.build/docs/installation |
| protoc-gen-go | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` |
| protoc-gen-go-grpc | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| golang-migrate CLI | `brew install golang-migrate` or see [docs](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) |
| wire CLI | `go install github.com/google/wire/cmd/wire@latest` |
| grpcurl (optional, for testing) | `brew install grpcurl` |

Install all Go tools at once:

```bash
make install-tools
```

---

## Getting started

### 1. Clone and enter the project

```bash
git clone https://github.com/97vaibhav/todo-crud-wire-tutorial-grpc.git
cd todo-crud-wire-tutorial-grpc
```

### 2. Copy environment file

```bash
cp .env.example .env
```

The defaults work out of the box with the Docker Compose setup. Edit `.env` only if you use a different Postgres instance.

### 3. Install Go dependencies

```bash
go mod download
```

### 4. Generate proto code

This reads `proto/todo/v1/todo.proto` and writes Go files into `gen/`:

```bash
make proto
```

Generated files:
- `gen/todo/v1/todo.pb.go` — message types
- `gen/todo/v1/todo_grpc.pb.go` — service client + server interfaces

### 5. Start Postgres

```bash
make docker-up
```

This starts a Postgres 16 container on port `5432` with the credentials from `docker-compose.yml`.

### 6. Run migrations

```bash
make migrate-up
```

Creates the `todos` table. Migration files live in `migrations/`.

### 7. Start the server

```bash
make run
```

The gRPC server starts on port `50051`.

---

## Testing the API

### With grpcurl

Because the server registers gRPC reflection, grpcurl can discover the schema at runtime — no proto file needed on the client side.

**List available services:**
```bash
grpcurl -plaintext localhost:50051 list
```

**List methods on TodoService:**
```bash
grpcurl -plaintext localhost:50051 list todo.v1.TodoService
```

**Create a todo:**
```bash
grpcurl -plaintext \
  -d '{"title": "Buy milk", "description": "From the corner store"}' \
  localhost:50051 todo.v1.TodoService/CreateTodo
```

Expected response:
```json
{
  "todo": {
    "id": "a1b2c3d4-...",
    "title": "Buy milk",
    "description": "From the corner store",
    "status": "TODO_STATUS_PENDING",
    "createdAt": "2026-02-18T10:00:00Z",
    "updatedAt": "2026-02-18T10:00:00Z"
  }
}
```

---

## All Makefile commands

```bash
make install-tools   # install protoc-gen-go, protoc-gen-go-grpc, wire, migrate
make proto           # run buf generate → regenerate gen/ from proto files
make wire            # run wire → regenerate wire/wire_gen.go
make docker-up       # start Postgres container
make docker-down     # stop Postgres container (data volume preserved)
make migrate-up      # apply all pending migrations
make migrate-down    # roll back the last migration
make run             # start the gRPC server
```

---

## Project structure

```
.
├── buf.yaml                                   # buf module config (proto root)
├── buf.gen.yaml                               # buf codegen config (plugins + output)
├── Makefile
├── Dockerfile                                 # multi-stage build → small Alpine image
├── docker-compose.yml                         # Postgres for local development
├── .env.example                               # copy to .env
│
├── proto/
│   └── todo/v1/todo.proto                     # API contract — edit this to change the API
│
├── gen/                                       # generated by `make proto` — do not edit
│   └── todo/v1/
│       ├── todo.pb.go
│       └── todo_grpc.pb.go
│
├── internal/
│   ├── config/
│   │   └── config.go                          # loads env vars into a Config struct
│   │
│   ├── domain/
│   │   └── todo.go                            # Todo entity + TodoRepository interface
│   │
│   ├── usecase/
│   │   ├── todo_usecase.go                    # TodoUsecase interface
│   │   └── todo_usecase_impl.go               # business logic implementation
│   │
│   ├── infrastructure/
│   │   ├── db/
│   │   │   └── postgres.go                    # GORM DB connection provider
│   │   └── repository/
│   │       └── todo_repository.go             # GORM implementation of TodoRepository
│   │
│   └── handler/
│       └── grpc/
│           └── todo_handler.go                # gRPC handler (proto ↔ usecase)
│
├── wire/
│   ├── wire.go                                # providers + injector (read by wire CLI)
│   └── wire_gen.go                            # generated by `make wire` — do not edit
│
├── migrations/
│   ├── 000001_create_todos_table.up.sql
│   └── 000001_create_todos_table.down.sql
│
└── cmd/
    └── server/
        └── main.go                            # entry point — intentionally thin
```

---

## How Wire works

Wire is a **compile-time** dependency injector. You declare providers (functions that build a struct from its dependencies), and Wire generates plain Go code that calls them in the right order.

```
config.Load()                 → *Config
infradb.NewPostgresDB(cfg)    → *gorm.DB
repository.NewTodoRepository(db) → domain.TodoRepository
usecase.NewTodoUsecase(repo)  → usecase.TodoUsecase
grpchandler.NewTodoHandler(uc) → *TodoHandler
```

Wire reads `wire/wire.go` (build-tagged `wireinject` so the Go compiler skips it) and writes `wire/wire_gen.go` (build-tagged `!wireinject` so only this is compiled). To regenerate after adding a new provider:

```bash
make wire
```

---

## Adding a new endpoint (e.g. GetTodo)

1. Add the RPC + messages to `proto/todo/v1/todo.proto`
2. Run `make proto` to regenerate the Go interface
3. Add `GetTodo` to `domain.TodoRepository` interface
4. Implement `GetTodo` in `infrastructure/repository/todo_repository.go`
5. Add `GetTodo` to `usecase.TodoUsecase` interface
6. Implement it in `usecase/todo_usecase_impl.go`
7. Implement the handler method in `internal/handler/grpc/todo_handler.go`

No changes needed to Wire or `main.go`.

---

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `postgres` | Postgres user |
| `DB_PASSWORD` | `postgres` | Postgres password |
| `DB_NAME` | `tododb` | Database name |
| `GRPC_PORT` | `50051` | Port the gRPC server listens on |
