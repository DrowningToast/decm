# ==========================================
# Stage 1: Build Stage
# ==========================================
FROM golang:1.24.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files from both backend and database package
COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/
COPY apps/backend/go.work apps/backend/go.work.sum ./apps/backend/
COPY packages/database/go.mod packages/database/go.sum ./packages/database/

# Download dependencies
WORKDIR /build/apps/backend
RUN go mod download

# Copy source code
WORKDIR /build
COPY apps/backend/ ./apps/backend/
COPY packages/database/ ./packages/database/

# Build the application
WORKDIR /build/apps/backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/bin/decm-backend \
    core-api/cmd/main.go

# ==========================================
# Stage 2: Runtime Stage
# ==========================================
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    postgresql-client

# Create non-root user
RUN addgroup -g 1000 decm && \
    adduser -D -u 1000 -G decm decm

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bin/decm-backend /app/decm-backend

# Copy migrations (if needed for auto-migration)
COPY packages/database/migrations /app/migrations

# Change ownership to non-root user
RUN chown -R decm:decm /app

# Switch to non-root user
USER decm

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/ready || exit 1

# Run the application
CMD ["/app/decm-backend"]
