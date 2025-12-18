# syntax=docker/dockerfile:1.7

FROM golang:1.24-alpine AS builder

WORKDIR /workspace

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache build-base git ca-certificates

COPY apps/backend/go.mod apps/backend/go.sum apps/backend/go.work apps/backend/go.work.sum ./apps/backend/
COPY packages/database/go.mod packages/database/go.sum ./packages/database/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd apps/backend && go mod download

COPY apps/backend ./apps/backend
COPY packages/database ./packages/database

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd apps/backend && go build -o bin/decm-backend core-api/cmd/main.go

FROM alpine:3.20 AS runtime

WORKDIR /app

# Install runtime dependencies including migrate CLI and Chromium for certificate image generation
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    curl \
    chromium \
    ttf-freefont && \
    wget -O /tmp/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz && \
    tar -xzf /tmp/migrate.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/migrate && \
    rm /tmp/migrate.tar.gz && \
    # Create symlink for chromedp to find chromium
    ln -s /usr/bin/chromium /usr/bin/google-chrome

# Copy application binary and migration files
COPY --from=builder /workspace/apps/backend/bin/decm-backend /app/decm-backend
COPY --from=builder /workspace/packages/database/migrations /app/migrations

# Create simple entrypoint script
RUN printf '#!/bin/sh\nset -e\nmigrate -path /app/migrations -database "$DATABASE_URL" up || { echo "Migration failed" >&2; exit 1; }\nexec /app/decm-backend\n' > /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

EXPOSE 8080

ENV PORT=8080 \
    ENVIRONMENT=production

HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=40s \
    CMD wget --quiet --tries=1 -O /dev/null http://localhost:8080/ready || exit 1

ENTRYPOINT ["/app/entrypoint.sh"]
