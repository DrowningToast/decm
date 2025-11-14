# syntax=docker/dockerfile:1.7

# Base stage with common dependencies
FROM node:20-alpine AS base
RUN apk update && apk add --no-cache libc6-compat bash git curl
WORKDIR /app

# Enable corepack and prepare pnpm
ENV PNPM_HOME="/root/.local/share/pnpm" \
    PATH="/root/.local/share/pnpm:$PATH"
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate

# ---
# Prepare stage - Use turbo prune to create minimal workspace
FROM base AS prepare
RUN pnpm install -g turbo@^2
COPY . .
# Generate a pruned monorepo for the "web" workspace with docker optimization
RUN turbo prune @decm/web --docker

# ---
# Builder stage - Install dependencies and build
FROM base AS builder

# First, copy only package.json files and lockfile (changes less often)
# This leverages Docker layer caching - if dependencies don't change, this layer is cached
COPY --from=prepare /app/out/json/ .
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile

# Then copy the full pruned source code
COPY --from=prepare /app/out/full/ .

# Build arguments for environment variables
ARG VITE_CORE_BACKEND_API
ARG VITE_APP_URL
ARG VITE_ENVIRONMENT=production
ARG VITE_GOOGLE_MAPS_API_KEY

# Uncomment to enable Turborepo Remote Caching
# ARG TURBO_TEAM
# ENV TURBO_TEAM=$TURBO_TEAM
# ARG TURBO_TOKEN
# ENV TURBO_TOKEN=$TURBO_TOKEN

# Set environment variables for build
ENV NODE_ENV=production \
    VITE_CORE_BACKEND_API=${VITE_CORE_BACKEND_API} \
    VITE_APP_URL=${VITE_APP_URL} \
    VITE_ENVIRONMENT=${VITE_ENVIRONMENT} \
    VITE_GOOGLE_MAPS_API_KEY=${VITE_GOOGLE_MAPS_API_KEY}

# Create .env file for build process
RUN printf "VITE_CORE_BACKEND_API=%s\nVITE_APP_URL=%s\nVITE_ENVIRONMENT=%s\nVITE_GOOGLE_MAPS_API_KEY=%s\n" \
        "${VITE_CORE_BACKEND_API:-}" \
        "${VITE_APP_URL:-}" \
        "${VITE_ENVIRONMENT:-production}" \
        "${VITE_GOOGLE_MAPS_API_KEY:-}" \
        > .env

# Build using turbo - explicitly build @decm/api first, then web app
RUN pnpm --filter=@decm/api build
RUN echo "Checking @decm/api build output:" && \
    ls -la packages/api/dist/ && \
    cat packages/api/package.json | grep -A 10 "exports"
RUN NODE_OPTIONS="--max-old-space-size=4096" pnpm turbo build --filter=@decm/web

# ---
# Runtime stage - Nginx to serve static files
FROM nginx:1.27-alpine AS runtime

WORKDIR /usr/share/nginx/html

# Copy nginx configuration
COPY apps/web/nginx/default.conf /etc/nginx/conf.d/default.conf

# Copy built static files from builder stage
COPY --from=builder /app/apps/web/dist ./

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl -f http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
