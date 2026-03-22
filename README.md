# DECM Platform - Decentralized Event Credential Management

Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification.

## Quick Start

### Prerequisites

- **Node.js** 18+ and **pnpm** 9.15.0+ (required — do not use npm or yarn)
- **Go** 1.24+
- **PostgreSQL** 15+
- **Docker** (for containerized services)

### Installation

```bash
# Install dependencies
pnpm install

# Start PostgreSQL and monitoring services
pnpm compose:up

# Start backend API (auto-runs migrations)
pnpm dev:core

# Start frontend (in another terminal)
pnpm dev
```

### Development URLs

| Service            | URL                            |
| ------------------ | ------------------------------ |
| Frontend           | http://localhost:3000          |
| Backend API        | http://localhost:8080/api/v1   |
| API Docs (Swagger) | http://localhost:8080/swagger/ |
| pgAdmin            | http://localhost:5050          |
| Grafana            | http://localhost:3001          |
| Prometheus         | http://localhost:9090          |
| Database           | localhost:5432                 |

## Tech Stack

### Frontend

- **React 19** with TypeScript
- **Vite** for fast development and builds
- **Tailwind CSS** + **Radix UI** for accessible components
- **React Query** for server state management
- **react-i18next** for internationalization (EN, TH)
- **React Router** with file-based routing (@generouted)
- **Wagmi** + **Viem** + **ReOwn AppKit** for Web3/wallet integration

### Backend

- **Go Fiber** — fast HTTP framework
- **PostgreSQL** with application-layer PII encryption (AES-256-GCM)
- **sqlc** for type-safe SQL queries
- **JWT** authentication with HTTP-only cookies
- **Swagger/OpenAPI** documentation

### Smart Contracts

- **Solidity** contracts in `apps/contracts/`
- Ethereum blockchain integration with gas price management

### Infrastructure

- **Monorepo**: Turbo + pnpm workspaces
- **Database**: PostgreSQL 16 with automated migrations
- **Monitoring**: Prometheus + Loki + Grafana + Alloy

## Project Structure

```
decm/
├── apps/
│   ├── web/                        # React 19 Frontend
│   │   └── src/
│   │       ├── components/         # Radix UI-based components
│   │       ├── pages/              # File-based routes
│   │       ├── lib/                # API client, i18n, utilities
│   │       ├── hooks/              # Custom React hooks
│   │       ├── services/           # External service integrations
│   │       └── context/            # React contexts
│   ├── backend/
│   │   └── core-api/               # Go Fiber API
│   │       ├── cmd/                # Application entry point
│   │       ├── internal/
│   │       │   ├── handler/        # HTTP handlers (13 feature domains)
│   │       │   ├── usecase/        # Business logic
│   │       │   ├── repositories/   # Data access + PII encryption
│   │       │   ├── middleware/     # Auth, logging, validation
│   │       │   ├── entity/         # Data models
│   │       │   ├── datagateway/    # External service contracts
│   │       │   └── worker/         # Background jobs
│   │       └── docs/               # Swagger/OpenAPI specs
│   └── contracts/                  # Solidity smart contracts
├── packages/
│   ├── database/                   # PostgreSQL + sqlc
│   │   ├── migrations/             # SQL migration files
│   │   ├── queries/                # SQL query files
│   │   └── go/generated/           # Generated Go code
│   ├── api/                        # Generated TypeScript client
│   ├── eslint-config/
│   └── typescript-config/
├── documentations/                 # Architecture and deployment docs
├── docker-compose.isolated.yml     # Development services
├── docker-compose.prod.yml         # Production
└── package.json                    # Root scripts
```

## Available Commands

### Development

```bash
pnpm dev              # Start all apps (Turbo)
pnpm dev:web          # Start React dev server only
pnpm dev:core         # Start Go API server only
pnpm build            # Production build (all apps)
pnpm build:core       # Build Go binary
pnpm start:core       # Run built Go binary
```

### Testing & Linting

```bash
pnpm lint             # Lint all (ESLint + TypeScript)
pnpm lint:go          # GolangCI-lint for backend
pnpm test:web         # Frontend tests (Vitest)
pnpm test:core        # Backend Go tests
```

### Database

```bash
pnpm compose:up       # Start PostgreSQL + monitoring containers
pnpm compose:down     # Stop all containers
pnpm db:setup         # Start DB + wait + run migrations
pnpm db:migrate       # Run database migrations
pnpm db:generate      # Generate Go code from SQL (sqlc)
pnpm db:console       # PostgreSQL CLI access
```

### API & Contracts

```bash
pnpm docs:core        # Generate OpenAPI docs (Swagger)
pnpm gen-api:core     # Generate TypeScript client from OpenAPI
pnpm contract:build   # Build smart contracts
pnpm contract:all     # Build + generate TypeScript contract types
```

## Core Features

### 1. Authentication & Onboarding

- **Wallet-based** login and registration with message signing
- **Google OAuth** login and registration
- Onboarding status tracking
- Role detection (Host, Issuer, Verified)

### 2. Profile Management

- Personal information with per-field privacy controls (name, email, bio, phone, address)
- Academic institution and academic email fields
- Profile picture management

### 3. Event Management

- Create, update, and delete events (academic, social, professional)
- Assign and manage issuers per event
- Manage event participants and attendee lists
- Deploy and manage per-event certificate smart contracts

### 4. Event Registration & Invitations

- Bulk import participant invitation lists
- Wallet-signed event join flow (on-chain registration)
- Participant invitation management (view, cancel)
- Password-protected event access

### 5. Certificate Lifecycle

- Bulk import certificate recipients
- Certificate template design (fonts, colors, text, layout)
- Publish, sign, revoke, and re-issue certificates
- Certificate image generation (PNG rendering)
- Wallet-signed certificate claiming on blockchain
- Mint readiness verification before on-chain operations

### 6. Certificate Sharing

- Generate shareable public links per certificate
- Toggle share active/inactive and set password protection
- Rate-limited public endpoints for share data and image retrieval

### 7. Issuer Management

- List all verified issuers in the system
- View events where the current user is an assigned issuer

### 8. Event Configuration

- Certificate design configuration (fonts, colors, text, templates)
- List available certificate font families
- Event registration configuration (requirements, parameters)
- Toggle certificate template published state

### 9. Inbox & Messaging

- Receive event invitations, certificate notifications, and system messages
- Mark individual or all messages as read
- Messages routable by credential ID, email, or wallet address

### 10. Blockchain & System

- Query current blockchain gas price
- System health status dashboard
- Planned maintenance schedules and closest upcoming schedule

## Internationalization (i18n)

Supported languages: **English** (en, default), **Thai** (th)

Translation files:

- `apps/web/src/lib/i18n/locales/en.json`
- `apps/web/src/lib/i18n/locales/th.json`

```tsx
import { useTranslation } from "react-i18next";

function MyComponent() {
    const { t } = useTranslation();
    return <h1>{t("common.welcome")}</h1>;
}
```

Language detection order: `localStorage` (`decm-language`) → browser navigator → HTML tag

## Security

### PII Encryption

All personally identifiable information is encrypted at the application layer using **AES-256-GCM** in the repository layer. See `documentations/pii-encryption.md` for details.

```go
emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
query, err := r.queries.CreateProfile(ctx, generated.CreateProfileParams{
    Email: emailEnc,
})
emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
```

### Authentication

- **Wallet-based authentication** with message signing
- **JWT sessions** via HTTP-only cookies
- **Google OAuth** integration
- Protected routes with authentication middleware

## API Documentation

Interactive Swagger docs at http://localhost:8080/swagger/

The TypeScript API client is generated from OpenAPI specs:

```typescript
import { DefaultApi } from "@decm/api";

const api = new DefaultApi({
    basePath: "http://localhost:8080/api/v1",
});

const response = await api.getProfile();
```

Regenerate after backend changes: `pnpm gen-api:core`

## Database

Migrations run automatically on `pnpm dev:core`. To run manually:

```bash
pnpm db:migrate
```

### Query Development

1. Write SQL queries in `packages/database/queries/`
2. Generate Go code: `pnpm db:generate`
3. Use generated types in repositories

## Development Workflow

### API-First

1. Define Go handler with OpenAPI annotations
2. Generate TypeScript client: `pnpm gen-api:core`
3. Use generated client in React frontend

### Database-First

1. Write SQL in `packages/database/queries/`
2. Generate Go code: `pnpm db:generate`
3. Use generated types in repositories

## Environment Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=decm
POSTGRES_PASSWORD=decm_password
POSTGRES_DB=decm

# JWT
JWT_SECRET=your-secret-key

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret

# Encryption
ENCRYPTION_KEY=your-32-byte-hex-key
```

See `.env.example` for the full list of configuration options.

## Documentation

Additional documentation in `documentations/`:

| File                                | Description                           |
| ----------------------------------- | ------------------------------------- |
| `backend-architecture.md`           | Backend layered architecture          |
| `frontend-architecture.md`          | Frontend structure and patterns       |
| `pii-encryption.md`                 | AES-256-GCM encryption implementation |
| `pii-encryption-quick-reference.md` | Quick encryption guide                |
| `qa-environment.md`                 | QA deployment guide                   |
| `production-environment.md`         | Production deployment guide           |

## Contributing

1. Follow code standards in `.rules/code-standards`
2. Use conventional commits
3. Write OpenAPI documentation for all API endpoints
4. Ensure type safety with sqlc and TypeScript
5. Use `pnpm` for package management (not npm/yarn)
