# ---- Build stage ----
FROM golang:1.23.4 AS builder

WORKDIR /app

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy full repo
COPY . .

# Build the Go API (main is in api/web)
RUN go build -o server ./api/web

# ---- Final runtime image ----
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 4000
CMD ["./server"]
