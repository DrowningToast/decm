# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DECM (Decentralized Event Management) is a Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification. Built as a monorepo using pnpm workspaces and Turbo, with a Go backend and React 19 frontend.

**Critical**: This project uses **pnpm** exclusively - never use npm, yarn, or bun for package management.

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
pnpm lint                  # ESLint
pnpm check-types          # TypeScript type checking
```

## Architecture Overview

### Three-Layer Backend Architecture (Go)

The backend follows a strict **Handler → UseCase → Repository** pattern with dependency injection:

**Handler Layer** (`apps/backend/core-api/internal/handler/`):
- HTTP request/response handling
- Input parsing and validation with struct tags
- Swagger/OpenAPI documentation (REQUIRED for all endpoints)
- Maps errors to HTTP status codes
- Thin layer (~30 lines per handler)

**UseCase Layer** (`apps/backend/core-api/internal/usecase/`):
- Business logic orchestration
- Transaction management
- Domain-specific validation
- Coordinates multiple repositories

**Repository Layer** (`apps/backend/core-api/internal/repositories/postgres/`):
- Database operations using sqlc-generated queries
- **PII encryption/decryption (CRITICAL - see Security section)**
- Error handling with `pgerrutils.ParsePgError()`
- Data mapping between database and domain entities

**Dependency Injection Flow** (see `apps/backend/core-api/cmd/main.go`):
```
Config → PG Pool → Repositories → UseCases → Handlers → Routes
```

Key architectural rules:
- Handlers NEVER access repositories directly
- Business logic belongs in UseCases, NOT handlers
- All database errors must be parsed with `pgerrutils.ParsePgError()`
- All user-facing errors use `customerror.New()` or `customerror.NewValidationErr()`

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

### Database Architecture (PostgreSQL + sqlc)

**Migration System**:
- Migrations in `packages/database/migrations/`
- Auto-run on backend start in development
- Use `pnpm db:migrate:create` to create new migrations

**Query Development**:
1. Write SQL queries in `packages/database/queries/`
2. Run `pnpm db:generate` to generate Go code
3. sqlc generates type-safe Go code in `packages/database/go/generated/`
4. Use generated types in repositories

**Query File Structure**:
```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, first_name)
VALUES ($1, $2)
RETURNING *;
```

## Critical Security Requirements

### PII Encryption (NON-NEGOTIABLE)

**All PII MUST be encrypted at the repository layer** using AES-GCM encryption.

**PII Fields**:
- Authentication: `google_connector_ref`, `github_connector_ref`
- Profile: `first_name`, `last_name`, `email`, `phone_number`, `address`, `bio`, `profile_picture_url`, `academic_institution`, `academic_email`

**Encryption Architecture**:
- **Database Layer**: Stores encrypted data as `TEXT` columns
- **Repository Layer**: Handles all encryption/decryption using `apps/backend/common/pgmapper`
- **Algorithm**: AES-256-GCM (deterministic for searchability)
- **Key Management**: `PII_ENCRYPTION_KEY` from environment variables ONLY

**Encryption Pattern** (see `.cursor/rules/database-security.mdc`):

```go
import "apps/backend/common/pgmapper"

// CREATE - Encrypt before insert
func (r *Repository) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error) {
    // 1. Encrypt PII fields
    emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
    if err != nil {
        return nil, err
    }

    // 2. Insert encrypted data
    query, err := r.queries.CreateProfile(ctx, generated.CreateProfileParams{
        Email: emailEnc,
    })

    // 3. Decrypt for return
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    return &entity.Profile{Email: emailDec}, nil
}

// READ - Decrypt after query
func (r *Repository) GetProfile(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
    query, err := r.queries.GetProfileByID(ctx, id)

    // Decrypt all PII fields
    email, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    return &entity.Profile{Email: email}, nil
}

// SEARCH - Encrypt search term
func (r *Repository) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
    // Encrypt the search term (deterministic encryption allows direct comparison)
    encryptedEmail, err := pgmapper.EncryptPII(email, r.piiEncryptionKey)

    query, err := r.queries.GetProfileByEmail(ctx, pgtype.Text{String: encryptedEmail, Valid: true})

    // Decrypt result
    emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
    return &entity.Profile{Email: emailDec}, nil
}
```

**SQL Queries (NO encryption in SQL)**:
```sql
-- Encryption handled in Go - SQL is clean
-- name: CreateProfile :one
INSERT INTO profiles (email, first_name)
VALUES (sqlc.narg(email), sqlc.narg(first_name))
RETURNING *;

-- name: GetProfileByEmail :one
SELECT * FROM profiles WHERE email = sqlc.arg(email);
```

**Checklist for PII Changes**:
- [ ] PII fields stored as `TEXT` in database schema
- [ ] Encrypt using `pgmapper.EncryptStringPtrToPgText()` before INSERT/UPDATE
- [ ] Decrypt using `pgmapper.DecryptPgTextToStringPtr()` after SELECT
- [ ] Search by encrypting search term with `pgmapper.EncryptPII()`
- [ ] NO `pgp_sym_encrypt`/`pgp_sym_decrypt` in SQL queries
- [ ] `PII_ENCRYPTION_KEY` from environment only
- [ ] Test encryption/decryption in repository tests

### Authentication

- **Wallet-based**: Ethereum message signing for registration
- **OAuth**: Google OAuth integration
- **Sessions**: JWT stored in HTTP-only cookies
- **Middleware**: `authentication_guard` and `verify_jwt` in `apps/backend/core-api/internal/middleware/`

## API-First Development Workflow

When adding/modifying endpoints:

1. **Add Swagger annotations** to Go handler:
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

2. **Generate TypeScript client**:
```bash
pnpm gen-api:core  # Generates OpenAPI spec → TypeScript client → builds @decm/api package
```

3. **Use in frontend**:
```typescript
import { DefaultApi } from '@decm/api';

const api = new DefaultApi({
  basePath: config.apiUrl,
  withCredentials: true,
});

const profile = await api.createProfile({ email: '...', firstName: '...' });
```

**Type Safety Flow**: Go structs → Swagger/OpenAPI → TypeScript types → React components

## Common Packages & Utilities

### Backend (Go)

**Error Handling** (`apps/backend/common/customerror/`):
```go
// User-facing error
return customerror.New(customerror.StatusBadRequest, "Invalid email", err)

// Validation error (auto-formats struct validation errors)
return customerror.NewValidationErr(validationErr)

// Database error parsing (auto-maps PG errors to HTTP codes)
return pgerrutils.ParsePgError(pgErr)
```

**PII Encryption** (`apps/backend/common/pgmapper/`):
```go
// Encrypt string pointer → pgtype.Text (most common)
encrypted, err := pgmapper.EncryptStringPtrToPgText(field, encryptionKey)

// Decrypt pgtype.Text → string pointer
decrypted, err := pgmapper.DecryptPgTextToStringPtr(field, encryptionKey)

// Encrypt/decrypt raw strings (for search)
encrypted, err := pgmapper.EncryptPII(plaintext, encryptionKey)
decrypted, err := pgmapper.DecryptPII(ciphertext, encryptionKey)
```

**Validation** (`apps/backend/common/validatorutils/`):
```go
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=18,max=100"`
}

if err := validatorutils.Validate(&req); err != nil {
    return customerror.NewValidationErr(err)
}
```

**Type Conversion** (`apps/backend/common/pgmapper/`):
```go
// pgtype.Text ↔ *string
stringPtr := pgmapper.PgTextToStringPtr(pgText)
pgText := pgmapper.StringPtrToPgText(stringPtr)

// pgtype.Timestamptz ↔ *time.Time
timePtr := pgmapper.PgTimestampzToTimePtr(timestampz)
timestampz := pgmapper.TimePtrToPgTimestampz(timePtr)
```

### Frontend (TypeScript/React)

**Form Handling**:
- React Hook Form + Zod for validation
- See `.cursor/rules/form-patterns.mdc` for patterns

**i18n**:
```typescript
import { useTranslation } from 'react-i18next';

const { t } = useTranslation();
return <h1>{t('common.welcome')}</h1>;
```

**Typography** (ALWAYS use for text):
```typescript
import { Typography } from '@/components/typography/typography';

<Typography variant="h1" tag="h1">{t('title')}</Typography>
<Typography variant="text" tag="p">{content}</Typography>
```

**API Client**:
```typescript
import { DefaultApi } from '@decm/api';
import { config } from '@/config/config';

const api = new DefaultApi({
  basePath: config.apiUrl,
  withCredentials: true,
});
```

## Environment Configuration

**Backend** (`.env` in repository root):
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

**Frontend** (`.env.client` - MUST use `VITE_` prefix):
```bash
VITE_API_URL=http://localhost:8080/api/v1
VITE_ENVIRONMENT=development
VITE_GOOGLE_OAUTH_CLIENT_ID=your-client-id
```

## Development URLs

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1
- Swagger Docs: http://localhost:8080/swagger/
- Database: localhost:5432
- Health Check: http://localhost:8080/
- Readiness: http://localhost:8080/ready

## Project Structure Details

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
│   │   │   │   └── layouts/              # Layout wrappers
│   │   │   ├── lib/
│   │   │   │   ├── api/                  # API client setup
│   │   │   │   └── i18n/                 # i18next configuration
│   │   │   └── hooks/                    # React hooks
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
│   │   ├── migrations/                   # SQL migrations (golang-migrate)
│   │   ├── queries/                      # SQL queries for sqlc
│   │   └── go/generated/                 # sqlc-generated Go code
│   └── api/                              # Generated TypeScript client
│       └── src/                          # OpenAPI-generated code
└── scripts/                              # Node.js utility scripts
    ├── db-migrate.js                     # Migration runner
    └── db-env.js                         # Database config loader
```

## Key Cursor Rules

The `.cursor/rules/` directory contains detailed patterns (summarized here):

**database-security.mdc**: PII encryption patterns at repository layer (CRITICAL - read first)

**go-backend-architecture.mdc**: Three-layer architecture, dependency injection, error handling

**repository-patterns.mdc**: Database operations with encryption (CREATE/READ/UPDATE/DELETE patterns)

**api-generation.mdc**: OpenAPI → TypeScript workflow

**form-patterns.mdc**: React Hook Form + Zod validation patterns

**error-handling.mdc**: Unified error handling across Go backend and React frontend

**authentication-security.mdc**: Wallet-based auth, OAuth, JWT sessions

**environment-config.mdc**: Configuration management, .env patterns

**development-workflow.mdc**: Development commands, troubleshooting

**testing-conventions.mdc**: Testing patterns for Go and React

## Important Development Notes

1. **Package Manager**: Use `pnpm` exclusively - the project won't work with npm/yarn/bun
2. **PII Encryption**: ALWAYS encrypt PII at repository layer - this is non-negotiable for security
3. **Swagger Docs**: ALL endpoints MUST have Swagger annotations
4. **API Generation**: Run `pnpm gen-api:core` after backend changes to sync TypeScript client
5. **Database Queries**: Use sqlc - never write raw SQL in Go code
6. **Error Handling**: Always use `customerror` package, never return raw errors to clients
7. **Typography**: Use Typography component for ALL text in frontend
8. **i18n**: All user-facing text MUST use `t()` translation function
9. **Form Validation**: Validate on BOTH frontend (Zod) and backend (struct tags)
10. **Migrations**: Auto-run in dev mode, but test manually in production scenarios

## Testing

**Backend (Go)**:
```bash
cd apps/backend
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -cover ./...             # With coverage
go test -race ./...              # With race detector
```

**Frontend (TypeScript/React)**:
```bash
pnpm test                        # Run tests (when configured)
pnpm test:watch                  # Watch mode
pnpm test:coverage               # Coverage report
```

## Troubleshooting

**Backend won't start**:
- Check `.env` file exists with all required variables
- Verify PostgreSQL is running: `pnpm compose:up`
- Check database connection: `pnpm db:console`

**API generation fails**:
- Ensure backend compiles: `cd apps/backend && go build core-api/cmd/main.go`
- Verify Swagger annotations are correct
- Run steps separately: `pnpm docs:core` then check errors

**Frontend API errors**:
- Regenerate TypeScript client: `pnpm gen-api:core`
- Check CORS settings in backend config
- Verify `VITE_API_URL` in `.env.client`

**Database migration issues**:
- Check migration files in `packages/database/migrations/`
- View migration status: `pnpm db:migrate:version`
- Rollback and retry: `pnpm db:reset`

**PII encryption errors**:
- Verify `PII_ENCRYPTION_KEY` is set in `.env`
- Ensure using `pgmapper` functions, not raw encryption
- Check encryption key length (must be 32 bytes for AES-256)
