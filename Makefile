.PHONY: proto wire migrate-up migrate-down run docker-up docker-down install-tools

DB_URL=postgresql://postgres:postgres@localhost:5432/tododb?sslmode=disable

## ── Tools ────────────────────────────────────────────────────────────────────
## Install all required code-generation tools.
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/google/wire/cmd/wire@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Seed the database with default groups and the initial admin user.
## Run once after migrate-up. Safe to run again (idempotent).
seed:
	go run ./cmd/seed

## ── Code generation ──────────────────────────────────────────────────────────
## Generate Go code from proto files using buf.
proto:
	buf generate

## Regenerate wire_gen.go from wire.go providers.
wire:
	cd wire && wire

## ── Database ─────────────────────────────────────────────────────────────────
## Apply all pending migrations.
migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

## Roll back the last applied migration.
migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

## ── Application ──────────────────────────────────────────────────────────────
## Run the gRPC server locally.
run:
	go run ./cmd/server

## Run the REST API gateway (translates REST to gRPC). Start the gRPC server first (make run).
run-gateway:
	go run ./cmd/gateway

## ── Docker ───────────────────────────────────────────────────────────────────
## Start Postgres container in the background.
docker-up:
	docker-compose up -d

## Stop and remove containers (data volume is preserved).
docker-down:
	docker-compose down
