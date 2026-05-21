# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies first — separate layer improves cache reuse
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bin/server ./cmd/api

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.20

# Non-root user for least-privilege execution
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/bin/server ./server

USER app

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./server"]
