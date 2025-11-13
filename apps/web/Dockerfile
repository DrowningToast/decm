# ==========================================
# Stage 1: Build Stage
# ==========================================
FROM node:20-alpine AS builder

# Install pnpm
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

# Set working directory
WORKDIR /build

# Copy package files
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY turbo.json ./
COPY apps/web/package.json ./apps/web/
COPY packages/api/package.json ./packages/api/
COPY packages/eslint-config/package.json ./packages/eslint-config/
COPY packages/typescript-config/package.json ./packages/typescript-config/

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY apps/web/ ./apps/web/
COPY packages/api/ ./packages/api/
COPY packages/eslint-config/ ./packages/eslint-config/
COPY packages/typescript-config/ ./packages/typescript-config/

# Build the application
RUN pnpm --filter @decm/web build

# ==========================================
# Stage 2: Runtime Stage with Nginx
# ==========================================
FROM nginx:1.27-alpine

# Install runtime dependencies
RUN apk add --no-cache \
    curl \
    tzdata

# Copy custom nginx configuration
COPY apps/web/nginx.conf /etc/nginx/nginx.conf

# Copy built assets from builder
COPY --from=builder /build/apps/web/dist /usr/share/nginx/html

# Create non-root user for nginx
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid

# Switch to non-root user
USER nginx

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:80/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
