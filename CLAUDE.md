# CLAUDE.md - DECM Platform Comprehensive Documentation

This file provides complete guidance to Claude Code when working with the DECM (Decentralized Event Management) platform codebase.

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
13. [Cursor Rules & Documentation](#cursor-rules--documentation)
14. [Common Patterns](#common-patterns)
15. [Troubleshooting](#troubleshooting)

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

**READ Operation:**

```go
func (r *Repository) GetProfile(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
    query, err := r.queries.GetProfileByID(ctx, id)

    // Decrypt all PII fields in Go
    email, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    return &entity.Profile{Email: email}, nil
}
```

**SEARCH Operation:**

```go
func (r *Repository) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
    // Encrypt the search term in Go (deterministic encryption)
    encryptedEmail, err := pgmapper.EncryptPII(email, r.piiEncryptionKey)

    query, err := r.queries.GetProfileByEmail(ctx, pgtype.Text{String: encryptedEmail, Valid: true})

    // Decrypt result
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    return &entity.Profile{Email: emailDec}, nil
}
```

**SQL Query Pattern (NO encryption in SQL):**

```sql
-- Encryption is done in Go application code, NOT in SQL
-- name: CreateProfile :one
INSERT INTO profiles (email, first_name)
VALUES (sqlc.narg(email), sqlc.narg(first_name))
RETURNING *;

-- name: GetProfileByEmail :one
SELECT * FROM profiles WHERE email = sqlc.arg(email);
```

#### Available Functions

- `EncryptStringPtrToPgText(field, key)` - Encrypt string pointer → pgtype.Text
- `DecryptPgTextToStringPtr(field, key)` - Decrypt pgtype.Text → string pointer
- `EncryptPII(plaintext, key)` - Encrypt raw strings (for search)
- `DecryptPII(ciphertext, key)` - Decrypt raw strings

#### PII Encryption Checklist

- [ ] PII fields stored as `TEXT` in database schema
- [ ] Encrypt using `pgmapper.EncryptStringPtrToPgText()` before INSERT/UPDATE in Go
- [ ] Decrypt using `pgmapper.DecryptPgTextToStringPtr()` after SELECT in Go
- [ ] Search by encrypting search term with `pgmapper.EncryptPII()`
- [ ] NO database-level encryption
- [ ] `PII_ENCRYPTION_KEY` from environment only
- [ ] No hardcoded keys
- [ ] Test encryption/decryption roundtrip

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

### Database Operations

```bash
pnpm db:generate           # Generate Go code from SQL queries (sqlc)
pnpm db:migrate            # Run pending migrations
pnpm db:migrate:create     # Create new migration file
pnpm db:migrate:down       # Rollback last migration
pnpm db:reset              # Rollback all and re-run migrations
pnpm db:console            # PostgreSQL CLI access
pnpm compose:up            # Start PostgreSQL container
pnpm compose:down          # Stop PostgreSQL container
```

### Frontend Development

```bash
pnpm dev                   # Start Vite dev server (Turbo)
pnpm build                 # Production build
pnpm lint                  # ESLint + TypeScript type checking (all workspaces)
pnpm lint:web              # Frontend linting only
pnpm lint:core             # Backend linting only
pnpm test                  # Run frontend unit tests
pnpm test:ui               # Interactive test dashboard
pnpm test:coverage         # Test coverage report
```

### Development URLs

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api/v1
- **Swagger Docs**: http://localhost:8080/swagger/
- **Database**: localhost:5432
- **Health Check**: http://localhost:8080/
- **Readiness**: http://localhost:8080/ready

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

#### Dependency Injection Flow

```
Config → PG Pool → Repositories → UseCases → Handlers → Routes
```

#### Architectural Rules

- Handlers NEVER access repositories directly
- Business logic belongs in UseCases, NOT handlers
- All database errors must be parsed with `pgerrutils.ParsePgError()`
- All user-facing errors use `customerror.New()` or `customerror.NewValidationErr()`
- All PII encryption happens in Repository layer only

### Frontend Architecture (React 19)

**Routing**: File-based routing with `@generouted/react-router`

- Routes defined in `apps/web/src/pages/`
- Auto-generated route configuration

**State Management**:

- React Query (`@tanstack/react-query`) for server state
- Zustand for client state
- React Context for auth/theme

**Component Structure**:

- `components/ui/` - Radix UI primitives
- `components/pages/` - Page-specific components
- `components/typography/` - Typography system (use for ALL text)
- `components/layouts/` - Layout wrappers

**API Integration**:

- Type-safe client generated from OpenAPI specs (`@decm/api` package)
- All API calls use generated client, never raw fetch/axios

---

## Backend Development

### Error Handling

```go
// User-facing error
return customerror.New(customerror.StatusBadRequest, "Invalid email", err)

// Validation error (auto-formats struct validation errors)
return customerror.NewValidationErr(validationErr)

// Database error parsing (auto-maps PG errors to HTTP codes)
return pgerrutils.ParsePgError(pgErr)
```

### Validation Pattern

```go
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=18,max=100"`
}

// In handler
if err := validatorutils.Validate(&req); err != nil {
    return customerror.NewValidationErr(err)
}
```

### Type Conversion

```go
// pgtype.Text ↔ *string
stringPtr := pgmapper.PgTextToStringPtr(pgText)
pgText := pgmapper.StringPtrToPgText(stringPtr)

// pgtype.Timestamptz ↔ *time.Time
timePtr := pgmapper.PgTimestampzToTimePtr(timestampz)
timestampz := pgmapper.TimePtrToPgTimestampz(timePtr)
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
// @Failure 401 {object} customerror.ErrResponse
// @Router /api/v1/profiles [post]
func (h *Handler) CreateProfile(ctx *fiber.Ctx) error { ... }
```

---

## Frontend Development

### Typography Component

ALWAYS use Typography component for text - never use plain text elements:

```typescript
import { Typography } from '@/components/typography/typography';

<Typography variant="h1" tag="h1">{t('title')}</Typography>
<Typography variant="text" tag="p">{content}</Typography>
```

### Internationalization (i18n)

All user-facing text must use `t()` translation function:

```typescript
import { useTranslation } from 'react-i18next';

const { t } = useTranslation();
return <h1>{t('common.welcome')}</h1>;
```

### Form Handling (React Hook Form + Zod)

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
    email: z.string().email("Invalid email"),
});

const {
    register,
    formState: { errors },
} = useForm({
    resolver: zodResolver(schema),
});
```

### React Query - Centralized Query Keys

Use centralized query keys from `apps/web/src/lib/queryKeys.ts`:

```typescript
import { queryKeys } from "@/lib/queryKeys";

useQuery({
    queryKey: queryKeys.event.byId(eventId),
    queryFn: () => fetchEvent(eventId),
});

useMutation({
    mutationFn: updateEvent,
    onSuccess: () => {
        queryClient.invalidateQueries({
            queryKey: queryKeys.event.all,
        });
    },
});
```

### API Client Usage

```typescript
import { DefaultApi } from "@decm/api";
import { config } from "@/config/config";

const api = new DefaultApi({
    basePath: config.apiUrl,
    withCredentials: true,
});

const profile = await api.createProfile({ email: "...", firstName: "..." });
```

### Protected Routes

Use `ProtectedRoute` component for authenticated pages:

```typescript
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";

export default function HostPage() {
    return (
        <ProtectedRoute requireHost={true}>
            <div>Host Content</div>
        </ProtectedRoute>
    );
}
```

**Props:**

- `requireHost?: boolean` - Require verified host/organizer role
- `requireIssuer?: boolean` - Require verified issuer role
- `requireAuthenticated?: boolean` - Explicitly control authentication (default: true)
- `unauthorizedRedirectTo?: Path` - Custom redirect path for failed checks

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

**Query File Structure:**

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, first_name)
VALUES ($1, $2)
RETURNING *;
```

### PII Field Schema

PII fields must be stored as `TEXT` (they'll be encrypted by application):

```sql
ALTER TABLE profiles ADD COLUMN email TEXT;
ALTER TABLE profiles ADD COLUMN first_name TEXT;
ALTER TABLE profiles ADD COLUMN last_name TEXT;
```

---

## API-First Development

### Development Workflow

1. **Add Swagger annotations** to Go handler
2. **Generate TypeScript client**: `pnpm gen-api:core`
3. **Use in frontend** with type-safe client
4. **Type safety maintained** end-to-end: Go → Swagger → TypeScript → React

### Example Handler with Swagger

```go
// @Summary Get event by ID
// @Description Retrieves event details with password/invitation status
// @ID get-event
// @Tags events
// @Param eventId path string true "Event ID"
// @Success 200 {object} EventResponse
// @Failure 404 {object} customerror.ErrResponse
// @Router /api/v1/events/{eventId} [get]
func (h *Handler) GetEvent(ctx *fiber.Ctx) error {
    eventId := ctx.Params("eventId")
    // Implementation
}
```

---

## Security Requirements

### Authentication

- **Wallet-based**: Ethereum message signing for registration
- **OAuth**: Google OAuth integration
- **Sessions**: JWT stored in HTTP-only cookies
- **Middleware**: `authentication_guard` and `verify_jwt` in `apps/backend/core-api/internal/middleware/`

### Key Environment Variables

```bash
# PII Encryption (CRITICAL)
PII_ENCRYPTION_KEY=your-256-bit-key

# JWT
JWT_SECRET_KEY=your-jwt-secret-key
JWT_EXPIRATION=24h
JWT_ISSUER=decm-core

# OAuth
GOOGLE_OAUTH_CLIENT_ID=your-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=decm
POSTGRES_PASSWORD=decm_password
POSTGRES_DATABASE=decm
```

**Frontend** (same `.env` file - MUST use `VITE_` prefix):

```bash
VITE_CORE_BACKEND_API=http://localhost:8080
VITE_ENVIRONMENT=development
VITE_GOOGLE_OAUTH_CLIENT_ID=your-client-id
VITE_GOOGLE_MAPS_API_KEY=your-maps-api-key
```

**Note**: Both backend and frontend variables are stored in the same root `.env` file. Frontend variables must be prefixed with `VITE_`.

## Development URLs

- React Hook Form + Zod for form validation
- Always use Typography component for text
- NEVER use `dangerouslySetInnerHTML`
- i18n: all user text must use `t()` translation function
- Validate on BOTH frontend (Zod) and backend (struct tags)

### Directory Structure

```
decm/
├── apps/
│   ├── web/                              # React 19 Frontend
│   │   ├── src/
│   │   │   ├── pages/                    # File-based routes (@generouted)
│   │   │   ├── components/
│   │   │   │   ├── ui/                   # Radix UI primitives
│   │   │   │   ├── pages/                # Page-specific components
│   │   │   │   ├── typography/           # Typography system
│   │   │   │   ├── auth/                 # Authentication components
│   │   │   │   └── layouts/              # Layout wrappers
│   │   │   ├── lib/
│   │   │   │   ├── api/                  # API client setup
│   │   │   │   ├── i18n/                 # i18next configuration
│   │   │   │   ├── queryKeys.ts          # Centralized React Query keys
│   │   │   │   └── hooks/                # Utility hooks
│   │   │   ├── hooks/                    # Custom React hooks
│   │   │   ├── services/                 # External services (api.ts)
│   │   │   └── test/                     # Test setup and utilities
│   │   └── package.json
│   └── backend/
│       ├── common/                       # Shared packages
│       │   ├── customerror/              # Error handling
│       │   ├── pgmapper/                 # Type conversion + PII encryption
│       │   ├── encryptutils/             # AES-GCM encryption
│       │   ├── pgerrutils/               # PostgreSQL error parsing
│       │   ├── validatorutils/           # Struct validation
│       │   └── pgclient/                 # PostgreSQL client
│       ├── services/                     # External services
│       │   ├── auth/                     # JWT service
│       │   └── oauth/                    # OAuth providers
│       └── core-api/                     # Main API
│           ├── cmd/main.go               # Application entry point
│           ├── config/                   # Configuration loading
│           ├── internal/
│           │   ├── handler/              # HTTP handlers (by feature)
│           │   ├── usecase/              # Business logic (by feature)
│           │   ├── repositories/         # Data access (postgres)
│           │   ├── entity/               # Domain entities
│           │   ├── datagateway/          # Repository interfaces
│           │   └── middleware/           # HTTP middleware
│           └── docs/                     # Generated Swagger docs
├── packages/
│   ├── database/                         # Database package
│   │   ├── migrations/                   # SQL migrations
│   │   ├── queries/                      # SQL queries for sqlc
│   │   └── go/generated/                 # sqlc-generated Go code
│   └── api/                              # Generated TypeScript client
│       └── src/                          # OpenAPI-generated code
└── scripts/                              # Node.js utility scripts
```

---

## Testing Strategy

### Frontend Unit Tests (Vitest + React Testing Library)

```bash
pnpm test                          # Run all tests
pnpm test --watch                  # Watch mode
pnpm test --coverage               # Coverage report
pnpm test --ui                     # Interactive UI dashboard
```

#### Test Categories

- **Hooks Tests** (26 tests)
    - `use-local-storage.test.ts` - Local storage persistence
    - `use-media-query.test.ts` - Media query matching
    - `usePaginationState.test.ts` - Pagination logic

- **Utility Functions Tests** (26 tests)
    - `utils.test.ts` - Class merging, delays
    - `event.utils.test.ts` - Event converters

- **Services Tests** (18 tests)
    - `OnboardService.test.ts` - Onboarding flows
    - `AuthService.test.ts` - Authentication

- **Component Tests** (44 tests)
    - `button.test.tsx` - Button variants and states
    - `input.test.tsx` - Input variations
    - `typography.test.tsx` - Typography styles

#### Test Configuration

- **Framework**: Vitest
- **Environment**: happy-dom
- **React Testing**: @testing-library/react
- **Mocking**: MSW for API mocking
- **Coverage Target**: 70%+ (100% for critical paths)

### Backend Testing (Go)

```bash
cd apps/backend
go test ./...        # Run all tests
go test -v ./...     # Verbose
go test -cover ./... # Coverage
go test -race ./...  # Race detector
```

---

## CI/CD Pipeline

### GitHub Actions Workflow

The CI/CD pipeline runs on every PR with comprehensive checks:

#### Pipeline Stages

1. **Label Gate** - `ready-to-review` label required
2. **Database** - Migrations + SQLC generation
3. **Backend** - Build + Tests + OpenAPI verification
4. **Frontend** - ESLint + Build + Tests
5. **Smart Contracts** - Solidity compilation
6. **Summary** - Aggregate results + failure handling

#### Key Features

- **Label-Based Gating**: PR must have `ready-to-review` label to run expensive checks
- **Composite Actions**: Reusable setup logic (frontend, backend, database tools)
- **Parallel Execution**: Independent checks run concurrently
- **Automatic Failure Handling**: Label removed on failure + comment posted
- **Time Savings**: 33-47% reduction in CI time through optimization

#### Job Dependencies

```
check-label (gate)
├── frontend-eslint ──┐
├── frontend-build ───┼──> frontend-tests
└── database-migrate ─> database-generate ─> backend-build ─> backend-tests
```

#### Composite Actions

Located in `.github/actions/`:

- `setup-frontend` - Node.js, pnpm, dependencies
- `setup-node-pnpm` - Basic Node.js + pnpm
- `setup-go-backend` - Go environment
- `install-db-tools` - migrate + sqlc tools

#### Developer Workflow

```bash
1. Create PR (WIP)
2. Add ready-to-review label when ready
3. Workflow runs (8-10 minutes)
4. If failure:
   - Label auto-removed
   - Comment posted with details
5. Fix issues, reapply label
6. Workflow reruns automatically
```

---

## Cursor Rules & Documentation

### Cursor Rules Files

Main Cursor rules are in `.cursorules` and `.cursor/rules/` directory:

#### Core Rules (`.cursorules`)

- Project Identity
- PII Encryption (⭐ Critical)
- Backend Architecture
- Swagger/OpenAPI
- Error Handling
- Validation Pattern
- SQL & sqlc Conventions
- Frontend Conventions
- Environment Configuration
- Development Endpoints
- Common Commands
- Code Quality
- Type Safety
- Testing
- Prohibited Patterns

#### Additional Rules

- **coderabbit.mdc** - CodeRabbit PR review automation
- **eslint-configuration.mdc** - ESLint and TypeScript linting
- **testing-setup.mdc** - Frontend testing with Vitest
- **ci-cd-workflow.mdc** - CI/CD workflow details
- **code-standards.mdc** - Code quality standards
- **project-conventions.mdc** - Project conventions
- **testing-conventions.mdc** - Testing conventions

### Documentation Files

- `README_CURSOR_RULES.md` - Comprehensive cursor rules guide
- `CI_CD_IMPROVEMENTS.md` - CI/CD optimization details
- `TEST_SUMMARY.md` - Testing and CI pipeline
- `QUERY_KEYS_MIGRATION_SUMMARY.md` - React Query keys pattern
- `CURSOR_RULES_GENERATED.md` - Generated rules summary
- `IMPLEMENTATION_SUMMARY.md` - Event detail implementation
- `PROTECTED_ROUTE_UPGRADE_SUMMARY.md` - Protected route upgrade

---

## Common Patterns

### Query Keys Pattern (React Query)

Centralized query keys in `apps/web/src/lib/queryKeys.ts`:

```typescript
export const queryKeys = {
    user: {
        all: ["user"],
        me: () => ["user", "me"],
    },
    event: {
        all: ["event"],
        byId: (id: string) => ["event", id],
        issuers: {
            byEventId: (id: string) => ["event", id, "issuers"],
        },
    },
    onboard: {
        status: {
            google: (accessToken: string, expiresIn: number) => [
                "onboard",
                "status",
                "google",
                accessToken,
                expiresIn,
            ],
            wallet: (signature: string) => ["onboard", "status", "wallet", signature],
        },
    },
};
```

### Event Detail Page Pattern

**Files:**

- `useEventDetailUsecase.ts` - Core logic (fetch, mutations, state)
- `EventDetailPage.tsx` - Component using usecase hook

**Features:**

- Password-protected events
- Invitation-only events
- Closed/read-only events
- Toast notifications
- Dynamic bottom nav variants

### Protected Route Pattern

```typescript
<ProtectedRoute requireHost={true} requireIssuer={false}>
    <HostDashboard />
</ProtectedRoute>
```

### Zustand Store Pattern

```typescript
import { create } from "zustand";

interface EventPasswordStore {
    password: string;
    setPassword: (password: string) => void;
    setOnSubmitCallback: (callback: (password: string) => Promise<void>) => void;
    resetPassword: () => void;
}

export const useEventPasswordStore = create<EventPasswordStore>((set) => ({
    password: "",
    setPassword: (password) => set({ password }),
    setOnSubmitCallback: (callback) => set({ onSubmitCallback: callback }),
    resetPassword: () => set({ password: "" }),
}));
```

### API Mutation Pattern

```typescript
const mutation = useMutation({
    mutationFn: async (password: string) => {
        return await api.submitEventPassword(eventId, password);
    },
    onSuccess: () => {
        toast.success("Password correct!");
        queryClient.invalidateQueries({
            queryKey: queryKeys.event.byId(eventId),
        });
    },
    onError: (error) => {
        toast.error("Invalid password");
    },
});
```

### Translation Pattern

**Files:**

- `apps/web/src/lib/i18n/locales/en.json`
- `apps/web/src/lib/i18n/locales/th.json`

```json
{
    "participant": {
        "events": {
            "detail": {
                "joined": "Joined",
                "accepted": "Accepted"
            }
        }
    }
}
```

Usage:

```typescript
const { t } = useTranslation();
return <Typography>{t("participant.events.detail.joined")}</Typography>;
```

---

## Certificate Features

### Overview

The DECM platform includes comprehensive certificate management features including:

- Certificate configuration with custom SVG templates
- Mint readiness validation (config, issuer signatures, contract deployment)
- Certificate publishing system with `is_published` flag
- Certificate image generation (SVG → PNG on-demand)
- Certificate claiming mechanism for participants
- Backend authorization guards for certificate operations

### Certificate Configuration (`is_published`)

**Database Field**: `event_certificate_configs.is_published` (BOOLEAN, default: FALSE)

**Purpose**: Track explicit publication status of certificate configurations

**Key Files**:

- Migration: `packages/database/migrations/000008_add_is_published_to_event_certificate_configs.up.sql`
- Backend Handler: `apps/backend/core-api/internal/handler/eventconfig/toggle_certificate_published.go`
- Frontend Hook: `apps/web/src/hooks/events/useToggleCertificatePublished.ts`
- UI: `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx`

**API Endpoint**:

```
PATCH /api/v1/events/{event_id}/config/certificate/published
Body: { "is_published": true }
```

**Usage Pattern**:

```typescript
// Frontend - Toggle published status
const { mutate: togglePublished } = useToggleCertificatePublished(eventId);
togglePublished({ is_published: true });

// Backend - Check published status
if config.IsPublished {
    // Certificate configuration is live
}
```

### Certificate Mint Readiness

**Purpose**: Validate all prerequisites before allowing certificate minting

**Validation Checks**:

1. ✅ Certificate configuration exists (SVG template uploaded)
2. ✅ ALL assigned issuers have signed (not just one)
3. ✅ Certificate contract deployed and address set

**Key Files**:

- Backend UseCase: `apps/backend/core-api/internal/usecase/eventconfig/check_certificate_mint_readiness.go`
- Backend Handler: `apps/backend/core-api/internal/handler/eventconfig/check_certificate_mint_readiness_response.go`
- Frontend Hook: `apps/web/src/hooks/events/useCertificateMintReadiness.ts`
- SQL Queries: `packages/database/queries/event_issuers.sql` (AllIssuersHaveSigned, GetTotalIssuersCount)

**API Endpoint**:

```
GET /api/v1/events/{event_id}/config/certificate/mint-readiness

Response:
{
  "is_ready": boolean,
  "has_certificate_config": boolean,
  "all_issuers_have_signed": boolean,
  "signed_issuers_count": number,
  "total_issuers_count": number,
  "has_certificate_contract": boolean,
  "certificate_contract_address": string | null,
  "missing_requirements": string[]
}
```

**Critical SQL Query**:

```sql
-- Returns true ONLY if ALL issuers have signed
SELECT
    CASE
        WHEN COUNT(*) = 0 THEN false
        WHEN COUNT(*) = COUNT(*) FILTER (WHERE is_signed = 1) THEN true
        ELSE false
    END AS all_issuers_signed
FROM event_issuers
WHERE event_id = $1 AND deleted_at IS NULL;
```

**Frontend Display**:

```typescript
const { data: mintReadiness } = useCertificateMintReadiness(eventId);

// Show status card
{mintReadiness?.is_ready ? (
    <div className="bg-green-50">✅ Ready to Mint Certificates</div>
) : (
    <div className="bg-blue-50">
        ○ Certificate Minting Requirements
        <ul>
            {mintReadiness?.missing_requirements.map(req => (
                <li key={req}>• {req}</li>
            ))}
        </ul>
    </div>
)}

// Disable publish button when not ready
<Button disabled={!mintReadiness?.is_ready}>
    Publish Configuration
</Button>
```

### Certificate Image Generation

**Purpose**: Generate PNG certificate images from SVG templates on-demand with participant data

**Architecture**: Server-side rendering (SVG → PNG) with template variable substitution

**Key Files**:

- Backend UseCase: `apps/backend/core-api/internal/usecase/event/generate_certificate_image.go`
- Backend UseCase (Participant): `apps/backend/core-api/internal/usecase/event/generate_certificate_image_for_participant.go`
- Backend Handler: `apps/backend/core-api/internal/handler/event/generate_certificate_image.go`
- Frontend Hook: `apps/web/src/hooks/useCertificateImage.ts`
- Frontend Component: `apps/web/src/components/pages/Participant/Certificates/CertificateDetail.tsx`

**Supported Template Variables**:

- `name` - Participant name (REQUIRED)
- `event_name` - Event name (REQUIRED)
- `academic_institution` - Academic institution (OPTIONAL)
- `certificate_title` - Certificate title (OPTIONAL)
- `certificate_subtitle` - Certificate subtitle (OPTIONAL)

**API Endpoint**:

```
GET /api/v1/certificates/{certificate_id}/image

Response Headers:
Content-Type: image/png
Cache-Control: public, max-age=86400
Content-Disposition: inline; filename=certificate.png

Security:
- Authentication Required (JWT cookie)
- Ownership Verification (user must own certificate)
- Returns 403 if unauthorized
```

**Backend Pattern**:

```go
// 1. Verify authorization
func (u *EventUsecase) GenerateCertificateImageForParticipant(
    ctx context.Context,
    certificateID uuid.UUID,
    authCredID uuid.UUID,
) ([]byte, error) {
    // Get certificate
    cert, err := u.eventRepo.GetCertificateByID(ctx, certificateID)

    // Check ownership
    if cert.ReceiverCredentialID != authCredID {
        return nil, customerror.New(customerror.StatusForbidden, "Unauthorized")
    }

    // Get template from S3
    template, err := u.s3Client.Get(cert.TemplateStorageKey)

    // Map variables
    variables := CertificateTemplateVariables{
        Name: cert.ReceiverName,
        EventName: cert.EventName,
        // ... other fields
    }

    // Render SVG to PNG
    return RenderSVGToPNG(template, variables)
}
```

**Frontend Pattern**:

```typescript
// Custom hook for image fetching
const useCertificateImage = ({ certificateId, enabled }) => {
    const [imageUrl, setImageUrl] = useState<string | null>(null);

    useEffect(() => {
        if (!enabled || !certificateId) return;

        fetch(`${API_URL}/certificates/${certificateId}/image`, {
            credentials: 'include', // Include auth cookies
        })
        .then(res => res.blob())
        .then(blob => {
            const url = URL.createObjectURL(blob);
            setImageUrl(url);
        });

        // Cleanup
        return () => {
            if (imageUrl) URL.revokeObjectURL(imageUrl);
        };
    }, [certificateId, enabled]);

    return { imageUrl, isLoading, error };
};

// Usage in component
const { imageUrl, isLoading } = useCertificateImage({
    certificateId: certificate.id,
    enabled: true,
});

{isLoading ? (
    <Spinner />
) : imageUrl ? (
    <img src={imageUrl} alt="Certificate" loading="lazy" />
) : (
    <div>Certificate not available</div>
)}
```

**SVG Template Requirements**:

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600">
    <!-- IDs must match database field names -->
    <text id="certificate_title">TITLE</text>
    <text id="name">NAME</text>
    <text id="event_name">EVENT</text>
    <text id="academic_institution">INSTITUTION</text>
    <text id="certificate_subtitle">SUBTITLE</text>
</svg>
```

### Certificate Claiming

**Purpose**: Allow participants to claim their certificates with password verification

**Key Files**:

- Frontend Hook: `apps/web/src/hooks/useClaimCertificate.ts`
- Frontend Component: `apps/web/src/components/pages/Participant/Certificates/CertificateDetail.tsx`
- Backend Endpoint: To be implemented at `POST /api/v1/events/{event_id}/certificates/{certificate_id}/claim`

**User Flow**:

1. Participant views certificate detail page
2. If unclaimed (status !== "completed"), sees "Claim Certificate" button
3. Clicks button → Password/PIN modal opens
4. Enters PIN/password for verification
5. Backend verifies password and mints NFT
6. Certificate status updated to "completed"
7. Button changes to "Certificate Claimed" badge

**Frontend Pattern**:

```typescript
// Custom hook
const useClaimCertificate = (certificateId: string, eventId: string) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (accountPassword: string) => {
            // TODO: Replace with real API call when backend ready
            return await coreApi.claimCertificate(
                { certificateId, eventId },
                { account_password: accountPassword }
            );
        },
        onSuccess: () => {
            toast.success(t("participant.certificates.claimSuccess"));
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.certificate.all });
            queryClient.invalidateQueries({ queryKey: QUERY_KEY.inbox.all });
        },
        onError: () => {
            toast.error(t("participant.certificates.claimError"));
        },
    });
};

// Usage in component
const { mutate: claimCertificate, isPending } = useClaimCertificate(certId, eventId);

const handleClaim = async () => {
    const password = await openPasswordPrompt({
        title: t("participant.certificates.claimTitle"),
        description: t("participant.certificates.claimDescription"),
        transactionType: "Certificate Claim",
    });

    claimCertificate(password);
};

// UI
{certificate.status === "completed" ? (
    <Badge variant="success">
        <Check className="w-4 h-4" />
        {t("participant.certificates.claimed")}
    </Badge>
) : (
    <Button onClick={handleClaim} disabled={isPending}>
        {isPending ? "Claiming..." : "Claim Certificate"}
    </Button>
)}
```

**Backend Requirements** (To be implemented):

```go
// Handler
// @Summary Claim certificate as participant
// @Description Mint certificate NFT and mark as claimed
// @Tags Certificates, Participant
// @Param event_id path string true "Event ID"
// @Param certificate_id path string true "Certificate ID"
// @Param body body ClaimCertificateRequest true "Account password"
// @Success 200 {object} ClaimCertificateResponse
// @Failure 401 {object} customerror.ErrResponse "Invalid password"
// @Failure 403 {object} customerror.ErrResponse "Not certificate owner"
// @Failure 409 {object} customerror.ErrResponse "Already claimed"
// @Router /api/v1/events/{event_id}/certificates/{certificate_id}/claim [post]
func (h *CertificateHandler) ClaimCertificate(ctx *fiber.Ctx) error

// Request/Response
type ClaimCertificateRequest struct {
    AccountPassword string `json:"account_password" binding:"required"`
}

type ClaimCertificateResponse struct {
    CertificateID      string `json:"certificate_id"`
    CertificateTokenID string `json:"certificate_token_id"`
    TransactionHash    string `json:"transaction_hash"`
    ClaimedAt          string `json:"claimed_at"`
}

// Business Logic Steps:
// 1. Verify JWT and get user credential ID
// 2. Verify account_password matches user's stored password
// 3. Get certificate and verify:
//    - Certificate exists
//    - Belongs to authenticated user (receiver_credential_id match)
//    - Not already claimed (certificate_token_id is NULL)
// 4. Mint NFT using system-managed wallet
// 5. Update database:
//    - Set certificate_token_id
//    - Set claimed timestamp
//    - Update inbox message to "claimed" status
// 6. Return certificate details with token ID and tx hash
```

### Certificate Settings Update Impact

**Key Consideration**: Changes to certificate configuration after publishing require re-approval from all issuers

**Warning Pattern**:

```typescript
{certificateConfig?.is_published && (
    <Alert variant="warning">
        ⚠️ Certificate configuration has been published.
        Any changes will require re-approval from all issuers.
    </Alert>
)}
```

**Implementation**:

- When config is published, show warning on settings page
- After updates, set `is_published = false`
- Require all issuers to sign again
- Re-check mint readiness before allowing publication

### Testing Certificate Features

**Backend Tests**:

```bash
cd apps/backend
go test -v ./core-api/internal/usecase/event -run "TestCertificate"
go test -v ./core-api/internal/usecase/eventconfig -run "TestCheckCertificateMintReadiness"
```

**Frontend Tests**:

```bash
pnpm test useCertificateImage
pnpm test useClaimCertificate
pnpm test useCertificateMintReadiness
pnpm test useToggleCertificatePublished
```

**Manual Testing Checklist**:

- [ ] Upload certificate template (SVG)
- [ ] Assign issuers to event
- [ ] Verify mint readiness shows correct status
- [ ] All issuers sign certificate config
- [ ] Mint readiness shows "ready" when all requirements met
- [ ] Publish certificate configuration
- [ ] Generate certificate image (loads correctly)
- [ ] Claim certificate (password verification works)
- [ ] Certificate status updates to "completed"
- [ ] Make changes to published config (warning appears)
- [ ] Verify `is_published` resets to false after update

---

## Deployment

### Production Deployment with Docker Compose

**Prerequisites**:

- Docker Engine 20.10+ and Docker Compose v2.0+
- At least 4GB RAM and 20GB disk space
- Domain name configured (optional, for SSL)

**Quick Start**:

```bash
# 1. Configure environment
cp .env.example .env
# Edit .env with production values

# 2. Start services
docker compose -f docker-compose.prod.yml up -d

# 3. Verify
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs -f
```

**Service Architecture**:

1. **PostgreSQL Database** - Persistent storage with auto-initialization
2. **Backend API** - Go REST API with auto-migrations on startup
3. **Frontend** - React 19 SPA served by Nginx
4. **Nginx Proxy Manager** (optional) - SSL termination and reverse proxy

**Critical Environment Variables**:

```bash
# Security (MUST CHANGE)
PII_ENCRYPTION_KEY=your-secure-32-character-key
JWT_SECRET=your-secure-jwt-secret

# Database
DB_NAME=decm
DB_USER=postgres
DB_PASSWORD=your-secure-password

# Blockchain
BLOCKCHAIN_NETWORK=mainnet
BLOCKCHAIN_RPC_URL=https://your-ethereum-node
BLOCKCHAIN_PRIVATE_KEY=your-private-key
BLOCKCHAIN_CHAIN_ID=1

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com

# Frontend
VITE_CORE_BACKEND_API=https://api.yourdomain.com
VITE_APP_URL=https://yourdomain.com
VITE_ENVIRONMENT=production
```

**Database Management**:

```bash
# Access PostgreSQL CLI
docker compose -f docker-compose.prod.yml exec postgres psql -U postgres -d decm

# Backup
docker compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres decm > backup.sql

# Restore
cat backup.sql | docker compose -f docker-compose.prod.yml exec -T postgres psql -U postgres -d decm
```

**SSL/HTTPS Setup with Nginx Proxy Manager**:

1. Access NPM admin UI at `http://your-server-ip:81`
2. Login with credentials from `.env` (NPM_ADMIN_EMAIL/NPM_ADMIN_PASSWORD)
3. Add Proxy Host for frontend:
    - Domain: yourdomain.com
    - Forward: frontend:80
    - SSL: Request Let's Encrypt certificate
4. Add Proxy Host for API:
    - Domain: api.yourdomain.com
    - Forward: backend:8080
    - SSL: Request Let's Encrypt certificate

**Monitoring**:

```bash
# Logs
docker compose -f docker-compose.prod.yml logs -f
docker compose -f docker-compose.prod.yml logs -f backend

# Health checks
curl http://localhost:8080/ready  # Backend
curl http://localhost:3000/        # Frontend

# Resource usage
docker stats
```

**Updates**:

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# Zero-downtime update (one service at a time)
docker compose -f docker-compose.prod.yml up -d --no-deps backend
docker compose -f docker-compose.prod.yml up -d --no-deps frontend
```

**Security Checklist**:

- [ ] Changed all default passwords
- [ ] Generated secure PII_ENCRYPTION_KEY (32 characters)
- [ ] Generated secure JWT_SECRET
- [ ] Configured SSL certificates
- [ ] Restricted PostgreSQL port (5432) to internal network
- [ ] Set up firewall rules (only 80, 443 exposed)
- [ ] Configured automated backups
- [ ] Updated CORS origins for production domain
- [ ] Verified S3 storage credentials
- [ ] Tested blockchain connection

**Troubleshooting**:

```bash
# Backend fails to start - check database connection
docker compose -f docker-compose.prod.yml logs postgres
docker compose -f docker-compose.prod.yml logs backend

# Frontend API errors - verify CORS
docker compose -f docker-compose.prod.yml exec backend env | grep CORS

# Out of disk space - clean up
docker system prune -a --volumes
docker system df -v
```

---

## Prohibited Patterns

❌ **NEVER do these:**

1. Use npm/yarn/bun - ONLY pnpm
2. Hardcode encryption keys - use environment variables
3. Skip PII encryption for any field
4. Use database-level encryption (pgcrypto) - encrypt in application layer
5. Write raw SQL in Go - use sqlc generated code
6. Access repositories directly from handlers
7. Return raw errors to clients - use `customerror` package
8. Skip Swagger annotations on endpoints
9. Use `dangerouslySetInnerHTML` in React
10. Skip database migrations - always use migration files
11. Store secrets in code
12. Skip type checking - ensure full TypeScript safety
13. Use hardcoded API URLs - use configuration/environment variables

---

## Troubleshooting

### Backend Won't Start

**Problem**: Backend fails to start with connection error

**Solutions**:

1. Check `.env` file has all required variables
2. Verify PostgreSQL is running: `pnpm compose:up`
3. Check database connection: `pnpm db:console`
4. Verify `PII_ENCRYPTION_KEY` is set and valid length (32 bytes for AES-256)

### PII Encryption Issues

- Regenerate TypeScript client: `pnpm gen-api:core`
- Check CORS settings in backend config
- Verify `VITE_CORE_BACKEND_API` in root `.env` file

**Solutions**:

1. Verify `PII_ENCRYPTION_KEY` is set in `.env`
2. Use `pgmapper` functions, not custom encryption
3. Check encryption key length: `echo $PII_ENCRYPTION_KEY | wc -c`
4. Ensure using application layer encryption, not database level
5. Test roundtrip: encrypt then decrypt to verify

### API Generation Fails

**Problem**: `pnpm gen-api:core` fails

**Solutions**:

1. Ensure backend compiles: `cd apps/backend && go build core-api/cmd/main.go`
2. Verify Swagger annotations are correct
3. Check for syntax errors in Go comments
4. Run `pnpm docs:core` first to see OpenAPI generation errors

### Frontend API Errors

**Problem**: Frontend cannot call backend API

**Solutions**:

1. Regenerate TypeScript client: `pnpm gen-api:core`
2. Check CORS settings in backend configuration
3. Verify `VITE_API_URL` in `.env.client`
4. Ensure backend is running on correct port (8080)
5. Check browser console for CORS errors

### Database Migration Issues

**Problem**: Migrations fail or don't run

**Solutions**:

1. Check migration files in `packages/database/migrations/`
2. View migration status: `pnpm db:migrate:version`
3. Rollback and retry: `pnpm db:reset`
4. Check PostgreSQL logs for errors
5. Verify database permissions

### Tests Failing in CI

**Problem**: Tests pass locally but fail in CI

**Solutions**:

1. Ensure `.env` and `.env.client` are set up
2. Check environment variables in GitHub Actions
3. Verify test setup in `apps/web/src/test/setup.ts`
4. Run tests with coverage: `pnpm test:coverage`
5. Check test output in GitHub Actions logs

### Linting Errors

**Problem**: ESLint or TypeScript errors

**Solutions**:

1. Run linting: `pnpm lint`
2. Auto-fix issues: `pnpm lint -- --fix`
3. Check specific errors: `pnpm lint:web` or `pnpm lint:core`
4. Ensure TypeScript configuration is correct
5. Check for missing imports or type definitions

---

## Environment Configuration

### Backend (.env file)

```bash
# Database (REQUIRED)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=decm
POSTGRES_PASSWORD=decm_password
POSTGRES_DATABASE=decm

# PII Encryption (CRITICAL - REQUIRED)
PII_ENCRYPTION_KEY=your-256-bit-encryption-key

# JWT (REQUIRED)
JWT_SECRET_KEY=your-jwt-secret-key
JWT_EXPIRATION=24h
JWT_ISSUER=decm-core

# OAuth (optional)
GOOGLE_OAUTH_CLIENT_ID=your-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env.client file - MUST use VITE\_ prefix)

```bash
VITE_API_URL=http://localhost:8080/api/v1
VITE_ENVIRONMENT=development
VITE_GOOGLE_OAUTH_CLIENT_ID=your-client-id
VITE_WALLETCONNECT_PROJECT_ID=your-project-id
```

### Generate PII Encryption Key

```bash
openssl rand -base64 32
# Output: REDACTED_PII_ENCRYPTION_KEY

# Add to .env
echo 'PII_ENCRYPTION_KEY=REDACTED_PII_ENCRYPTION_KEY' >> .env
```

---

## Quick Reference

### Most Used Commands

```bash
pnpm install           # Install dependencies
pnpm dev:core          # Start backend
pnpm dev               # Start frontend
pnpm compose:up        # Start database
pnpm gen-api:core      # Generate API client
pnpm db:generate       # Generate database code
pnpm lint              # Run linting and type checking
pnpm test              # Run tests
pnpm build             # Production build
pnpm db:reset          # Reset database
```

### Key Files to Know

- **`CLAUDE.md`** - This file
- **`.cursorules`** - Main cursor rules
- **`.env.example`** - Environment template
- **`apps/backend/core-api/cmd/main.go`** - Backend entry point
- **`apps/web/src/pages/`** - Frontend routes
- **`packages/database/migrations/`** - Database migrations
- **`.github/workflows/pr-checks.yml`** - CI/CD pipeline

### Important Directories

- `apps/backend/common/` - Shared utilities (encryption, error handling, validation)
- `apps/backend/core-api/internal/` - Backend application code
- `apps/web/src/` - Frontend application code
- `packages/database/` - Database and migrations
- `packages/api/` - Generated API client

---

## Additional Resources

### Documentation

- **README.md** - Quick start guide
- **README_CURSOR_RULES.md** - Cursor rules overview
- **CI_CD_IMPROVEMENTS.md** - CI/CD details
- **IMPLEMENTATION_SUMMARY.md** - Event detail implementation
- **PROTECTED_ROUTE_UPGRADE_SUMMARY.md** - Protected route features

### Files with Examples

- `apps/web/src/lib/queryKeys.ts` - React Query keys pattern
- `apps/web/src/components/auth/ProtectedRoute.tsx` - Protected route component
- `apps/backend/core-api/internal/handler/` - Handler examples
- `apps/backend/core-api/internal/repositories/postgres/` - Encryption patterns

---

## Final Checklist Before Committing

- [ ] Used `pnpm` for package management
- [ ] PII fields encrypted with `pgmapper` (if modified)
- [ ] No hardcoded secrets or keys
- [ ] Swagger annotations added (if backend endpoint)
- [ ] TypeScript types are complete (no `any`)
- [ ] Linting passes: `pnpm lint`
- [ ] Tests pass: `pnpm test`
- [ ] All translations added (EN + TH)
- [ ] Database migrations created (if schema changed)
- [ ] Error handling with `customerror` package
- [ ] No `dangerouslySetInnerHTML` in React components
- [ ] All PII stored as encrypted TEXT in database
- [ ] No raw SQL queries in Go code
- [ ] Repository layer handles encryption/decryption
- [ ] Handlers don't access repositories directly
- [ ] i18n keys used for all user-facing text

---

**Last Updated**: November 13, 2025
**Status**: ✅ Complete and Production Ready
**DECM Platform - Decentralized Event Management**
