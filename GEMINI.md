# GEMINI.md - DECM Platform Comprehensive Documentation

This file provides complete guidance to Gemini when working with the DECM (Decentralized Event Management) platform codebase.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Critical Rules](#critical-rules)
3. [Essential Commands](#essential-commands)
4. [Architecture Overview](#architecture-overview)
5. [Backend Development](#backend-development)
6. [Frontend Development](#frontend-development)
7. [Database & Migrations](#database--migrations)
8. [API-First Development](#api-first-development)
9. [Security Requirements](#security-requirements)
10. [Code Quality Standards](#code-quality-standards)
11. [Testing Strategy](#testing-strategy)
12. [CI/CD Pipeline](#cicd-pipeline)
13. [Common Patterns](#common-patterns)
14. [Troubleshooting](#troubleshooting)

---

## Project Overview

**DECM** (Decentralized Event Management) is a Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification. Built as a monorepo using pnpm workspaces and Turbo, with a Go backend and React 19 frontend.

### Tech Stack

**Frontend:**

- React 19 with TypeScript
- Vite for fast development
- Tailwind CSS + Radix UI for components
- React Query for API state management
- react-i18next for internationalization
- React Router with file-based routing (@generouted)

**Backend:**

- Go Fiber - Fast HTTP framework
- PostgreSQL with application-layer PII encryption (AES-GCM)
- sqlc for type-safe SQL queries
- JWT authentication with HTTP-only cookies
- Swagger/OpenAPI documentation

**Infrastructure:**

- Monorepo: Turbo + pnpm workspaces
- Database: PostgreSQL with automated migrations
- CI/CD: GitHub Actions with composite actions

---

## Critical Rules

### 1. Package Manager - ABSOLUTE ⚠️

**MUST use `pnpm`** for all package management operations:

- Never use npm, yarn, or bun
- Check `packageManager: "pnpm@9.15.0"` in package.json
- All scripts use `pnpm` exclusively

```bash
pnpm install              # Install dependencies
pnpm add package-name     # Add package
pnpm remove package-name  # Remove package
```

### 2. PII Encryption - NON-NEGOTIABLE 🔐

**ALL personally identifiable information MUST be encrypted at the application layer (repository) using AES-256-GCM encryption in Go code.**

**Note on conflicting documentation:** Some project files (`core-conventions.mdc`, `project-conventions.mdc`) mention `pgcrypto` for database-level encryption. However, the most detailed and recent documentation indicates that encryption is handled at the **application layer** in Go. **Follow the application-layer encryption pattern.**

#### PII Fields That Must Be Encrypted

**Authentication (2 fields):**

- `google_connector_ref`
- `github_connector_ref`

**Profile (9 fields):**

- `email`
- `first_name`, `last_name`
- `phone_number`
- `address`
- `bio`
- `profile_picture_url`
- `academic_institution`
- `academic_email`

#### Encryption Architecture

- **Location**: Application layer (Repository) ONLY
- **Algorithm**: AES-256-GCM in Go code (deterministic for searchability)
- **Key Management**: `PII_ENCRYPTION_KEY` from environment variables ONLY
- **Database Storage**: PII stored as `TEXT` columns (already encrypted)
- **Implementation**: Use `apps/backend/common/pgmapper` utilities
- **NO database-level encryption**: All encryption in Go application code

#### Encryption Patterns

**CREATE Operation:**

```go
import "apps/backend/common/pgmapper"

func (r *Repository) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error) {
    // 1. Encrypt PII fields in Go application layer
    emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, err
    }

    // 2. Insert encrypted data
    query, err := r.queries.CreateProfile(ctx, generated.CreateProfileParams{
        Email: emailEnc,
    })

    // 3. Decrypt for return in Go
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    return &entity.Profile{Email: emailDec}, nil
}
```

---

## Essential Commands

### Development Workflow (3 Terminals)

```bash
# Terminal 1: Database
pnpm compose:up

# Terminal 2: Backend (auto-runs migrations)
pnpm dev:core

# Terminal 3: Frontend
pnpm dev
```

### Backend Development

```bash
pnpm dev:core              # Start Go API server (with auto-migrations)
pnpm build:core            # Build Go binary to apps/backend/bin/
pnpm start:core            # Run built binary
pnpm docs:core             # Generate OpenAPI docs only (faster)
pnpm gen-api:core          # Full pipeline: OpenAPI → TypeScript client → build
```

---

## Architecture Overview

### Three-Layer Backend Architecture (Go)

The backend follows a strict **Handler → UseCase → Repository** pattern with dependency injection.

#### Handler Layer (`apps/backend/core-api/internal/handler/`)

- HTTP request/response handling
- Input parsing and validation with struct tags
- Swagger/OpenAPI documentation (REQUIRED for all endpoints)
- Maps errors to HTTP status codes
- Thin layer (~30 lines per handler)

#### UseCase Layer (`apps/backend/core-api/internal/usecase/`)

- Business logic orchestration
- Transaction management
- Domain-specific validation
- Coordinates multiple repositories
- No HTTP knowledge

#### Repository Layer (`apps/backend/core-api/internal/repositories/postgres/`)

- Database operations using sqlc-generated queries
- **PII encryption/decryption using AES-GCM (CRITICAL)**
- Error handling with `pgerrutils.ParsePgError()`
- Data mapping between database and domain entities

---

## Backend Development

### Go Naming Conventions

- Use `Id` instead of `ID` in variable names, struct fields, and function names. For example, `userId` is correct, `userID` is incorrect.

### Error Handling

```go
// User-facing error
return customerror.New(customerror.StatusBadRequest, "Invalid email", err)

// Validation error (auto-formats struct validation errors)
return customerror.NewValidationErr(validationErr)

// Database error parsing (auto-maps PG errors to HTTP codes)
return pgerrutils.ParsePgError(pgErr)
```

### Swagger/OpenAPI Documentation

ALL endpoints MUST have Swagger annotations:

```go
// @Summary Create user profile
// @Description Creates a new user profile with encrypted PII
// @ID create-profile
// @Tags profiles
// @Accept json
// @Produce json
// @Param profile body CreateProfileRequest true "Profile data"
// @Success 201 {object} ProfileResponse
// @Failure 400 {object} customerror.ErrResponse
// @Router /api/v1/profiles [post]
func (h *Handler) CreateProfile(ctx *fiber.Ctx) error { ... }
```

---

## Frontend Development

### API Key Transformation

Convert snake_case from the API to camelCase at the service level.

```typescript
// apps/web/src/services/EventService.ts
export async function getEventById(eventId: string) {
  const response = await api.getEvent({ id: eventId });
  // Transform snake_case to camelCase
  return {
    eventId: response.event_id,
    createdAt: response.created_at,
    updatedAt: response.updated_at,
  };
}

// Component receives camelCase props
<Component eventId={eventData.eventId} />
```

### Typography Component

ALWAYS use Typography component for text - never use plain text elements:

```typescript
import { Typography } from '@/components/typography/typography';

<Typography variant="h1" tag="h1">{t('title')}</Typography>
<Typography variant="text" tag="p">{content}</Typography>
```

### Internationalization (i18n)

All user-facing text must use the `t()` translation function:

```typescript
import { useTranslation } from 'react-i18next';

const { t } = useTranslation();
return <h1>{t('common.welcome')}</h1>;
```

---

## Database & Migrations

### Migration System

Migrations are in `packages/database/migrations/` and auto-run on backend start:

```bash
pnpm db:migrate:create     # Create new migration
pnpm db:migrate            # Run migrations
pnpm db:migrate:down       # Rollback last migration
pnpm db:reset              # Rollback all and re-run
```

### Query Development (sqlc)

1. Write SQL queries in `packages/database/queries/`
2. Run `pnpm db:generate` to generate Go code
3. Use generated types in repositories

---

## API-First Development

### Development Workflow

1. **Add Swagger annotations** to Go handler
2. **Generate TypeScript client**: `pnpm gen-api:core`
3. **Use in frontend** with type-safe client
4. **Type safety maintained** end-to-end: Go → Swagger → TypeScript → React

---

## Prohibited Patterns

❌ **NEVER do these:**

1. Use npm/yarn/bun - ONLY pnpm
2. Hardcode encryption keys - use environment variables
3. Skip PII encryption for any field. Follow the application-layer encryption guidance.
4. Write raw SQL in Go - use sqlc generated code
5. Access repositories directly from handlers
6. Return raw errors to clients - use `customerror` package
7. Skip Swagger annotations on endpoints
8. Use `dangerouslySetInnerHTML` in React
9. Skip database migrations - always use migration files
10. Store secrets in code
11. Skip type checking - ensure full TypeScript safety (no `any`)

---

This document is a guideline for interacting with the DECM platform. For more detailed information, please refer to the comprehensive `CLAUDE.md` file located in `.cursor/rules/`.
