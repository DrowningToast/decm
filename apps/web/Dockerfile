# syntax=docker/dockerfile:1.7

FROM node:20-alpine AS build

WORKDIR /workspace

ENV PNPM_HOME="/root/.local/share/pnpm" \
    PATH="/root/.local/share/pnpm:$PATH" \
    NODE_ENV=production

RUN apk add --no-cache bash git \
    && corepack enable \
    && corepack prepare pnpm@9.15.0 --activate

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json ./
COPY apps ./apps
COPY packages ./packages

RUN --mount=type=cache,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile

ARG VITE_API_URL
ARG VITE_APP_URL
ARG VITE_ENVIRONMENT=production
ARG VITE_GOOGLE_MAPS_API_KEY

ENV VITE_API_URL=${VITE_API_URL} \
    VITE_APP_URL=${VITE_APP_URL} \
    VITE_ENVIRONMENT=${VITE_ENVIRONMENT} \
    VITE_GOOGLE_MAPS_API_KEY=${VITE_GOOGLE_MAPS_API_KEY}

RUN printf "VITE_API_URL=%s\nVITE_APP_URL=%s\nVITE_ENVIRONMENT=%s\nVITE_GOOGLE_MAPS_API_KEY=%s\n" \
        "${VITE_API_URL:-}" \
        "${VITE_APP_URL:-}" \
        "${VITE_ENVIRONMENT:-production}" \
        "${VITE_GOOGLE_MAPS_API_KEY:-}" \
        > .env.client

RUN pnpm --filter @decm/api build
RUN pnpm --filter @decm/web build

FROM nginx:1.27-alpine AS runtime

WORKDIR /usr/share/nginx/html

COPY apps/web/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /workspace/apps/web/dist ./

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --retries=3 CMD wget -qO- http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
