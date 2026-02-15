# ---- Build stage ----
FROM golang:1.24.0-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY server/ ./server/

WORKDIR /app/server
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./api/web

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /

COPY --from=builder /app/server/server /server
COPY --from=builder /app/server/config.json /config.json

EXPOSE 4000
ENTRYPOINT ["/server"]