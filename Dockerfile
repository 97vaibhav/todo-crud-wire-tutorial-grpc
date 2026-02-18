# ── Stage 1: Build ────────────────────────────────────────────────────────────
# Use the full Go image to compile. This stage is large but temporary.
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first so Docker caches this layer.
# Only re-downloads modules when go.mod / go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build a statically linked binary.
# CGO_ENABLED=0 disables C bindings — required for scratch/alpine base images.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ── Stage 2: Run ──────────────────────────────────────────────────────────────
# Use a minimal Alpine image. Contains only libc and a shell — nothing else.
FROM alpine:3.19

WORKDIR /app

# Copy only the compiled binary from the builder stage.
COPY --from=builder /app/server .

# Copy migrations so you can run `migrate` inside the container if needed.
COPY --from=builder /app/migrations ./migrations

EXPOSE 50051

CMD ["./server"]
