# ---- Build stage ----
FROM golang:1.24.0 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./api/web

# ---- Final runtime image ----
FROM debian:bookworm-slim

# Use 'apt-get' because this is Debian, not Alpine
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 4000
CMD ["./server"]