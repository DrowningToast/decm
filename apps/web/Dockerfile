# syntax=docker/dockerfile:1.7

# Base stage with common dependencies
FROM node:20-alpine AS base
RUN apk update && apk add --no-cache \
    libc6-compat \
    bash \
    git \
    curl \
    python3 \
    make \
    g++
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

# Install dependencies with cache mount for pnpm store
# This caches downloaded packages across builds for faster installations
# Skip prepare scripts (lefthook) as git hooks aren't needed in Docker
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile --ignore-scripts

# Then copy the full pruned source code
COPY --from=prepare /app/out/full/ .

# Build arguments for environment variables
ARG VITE_CORE_BACKEND_API
ARG VITE_APP_URL
ARG VITE_ENVIRONMENT=production
ARG VITE_GOOGLE_MAPS_API_KEY
ARG VITE_WALLETCONNECT_PROJECT_ID
ARG VITE_GOOGLE_OAUTH_CLIENT_ID

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
    VITE_GOOGLE_MAPS_API_KEY=${VITE_GOOGLE_MAPS_API_KEY} \
    VITE_WALLETCONNECT_PROJECT_ID=${VITE_WALLETCONNECT_PROJECT_ID} \
    VITE_GOOGLE_OAUTH_CLIENT_ID=${VITE_GOOGLE_OAUTH_CLIENT_ID}

# Create .env file for build process
RUN printf "VITE_CORE_BACKEND_API=%s\nVITE_APP_URL=%s\nVITE_ENVIRONMENT=%s\nVITE_GOOGLE_MAPS_API_KEY=%s\nVITE_WALLETCONNECT_PROJECT_ID=%s\nVITE_GOOGLE_OAUTH_CLIENT_ID=%s\n" \
        "${VITE_CORE_BACKEND_API:-}" \
        "${VITE_APP_URL:-}" \
        "${VITE_ENVIRONMENT:-production}" \
        "${VITE_GOOGLE_MAPS_API_KEY:-}" \
        "${VITE_WALLETCONNECT_PROJECT_ID:-}" \
        "${VITE_GOOGLE_OAUTH_CLIENT_ID:-}" \
        > .env

# Build using turbo with increased Node.js memory limit
RUN NODE_OPTIONS="--max-old-space-size=4096" pnpm turbo build --filter=@decm/web

# ---
# Runtime stage - Nginx to serve static files
FROM nginx:1-alpine AS runtime

# Install curl for healthcheck
RUN apk add --no-cache curl

WORKDIR /usr/share/nginx/html

# Copy nginx configuration
COPY apps/web/nginx/default.conf /etc/nginx/conf.d/default.conf

# Copy built static files from builder stage
COPY --from=builder /app/apps/web/dist ./

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD curl -f http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
