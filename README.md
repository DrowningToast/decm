# DECM Platform - Decentralized Event Management

Web 3.0 platform for NFT ticketing, digital credentials, and academic identity verification.

## 🚀 Quick Start

### Prerequisites
- **Node.js** 18+ and **pnpm** 9.15.0+
- **Go** 1.21+
- **PostgreSQL** 15+
- **Docker** (for containerized database)

### Installation

```bash
# Install dependencies
pnpm install

# Start PostgreSQL database
pnpm compose:up

# Start backend API (auto-runs migrations)
pnpm dev:core

# Start frontend (in another terminal)
pnpm dev
```

### Development URLs
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080/api/v1
- **API Docs**: http://localhost:8080/swagger/
- **Database**: localhost:5432

## 📦 Tech Stack

### Frontend
- **React 19** with TypeScript
- **Vite** for fast development
- **Tailwind CSS** + **Radix UI** for components
- **React Query** for API state management
- **react-i18next** for internationalization
- **React Router** with file-based routing (@generouted)

### Backend
- **Go Fiber** - Fast HTTP framework
- **PostgreSQL** with **pgcrypto** for PII encryption
- **sqlc** for type-safe SQL queries
- **JWT** authentication with HTTP-only cookies
- **Swagger/OpenAPI** documentation

### Infrastructure
- **Monorepo**: Turbo + pnpm workspaces
- **Database**: PostgreSQL with automated migrations
- **Package Manager**: pnpm (REQUIRED)

## 🏗️ Project Structure

```
decm/
├── apps/
│   ├── web/                    # React 19 Frontend
│   │   ├── src/
│   │   │   ├── components/     # UI components
│   │   │   ├── pages/          # File-based routes
│   │   │   ├── lib/            # Utilities, API client, i18n
│   │   │   └── hooks/          # React hooks
│   │   └── package.json
│   └── backend/
│       └── core-api/           # Go Fiber API
│           ├── cmd/            # Application entry point
│           ├── internal/       # Private application code
│           │   ├── handler/    # HTTP handlers
│           │   ├── usecase/    # Business logic
│           │   └── repositories/ # Data access
│           └── docs/           # Swagger documentation
├── packages/
│   ├── database/               # PostgreSQL + sqlc
│   │   ├── migrations/         # SQL migrations
│   │   ├── queries/            # SQL queries
│   │   └── go/generated/       # Generated Go code
│   └── api/                   # Generated TypeScript client
└── package.json               # Root scripts
```

## 🛠️ Available Commands

### Development
```bash
# Frontend
pnpm dev              # Start React dev server (Turbo)
pnpm build           # Production build
pnpm lint            # Lint code
pnpm check-types     # TypeScript type checking

# Backend
pnpm dev:core        # Start Go API server
pnpm build:core      # Build Go binary
pnpm start:core      # Run built binary
pnpm docs:core       # Generate OpenAPI docs

# Database
pnpm compose:up      # Start PostgreSQL container
pnpm compose:down    # Stop database
pnpm db:generate     # Generate Go code from SQL (sqlc)
pnpm db:migrate      # Run database migrations
pnpm db:setup        # Complete database setup
pnpm db:console      # PostgreSQL CLI access

# API Generation
pnpm gen-api:core    # Generate TypeScript client from OpenAPI
```

## 🌐 Internationalization (i18n)

The platform supports multiple languages using **react-i18next**.

### Supported Languages
- **English** (en) - Default
- **Thai** (th)

### Adding Translations

1. Add translations to JSON files:
   - `apps/web/src/lib/i18n/locales/en.json`
   - `apps/web/src/lib/i18n/locales/th.json`

2. Use translations in components:
```tsx
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();
  
  return <h1>{t('common.welcome')}</h1>;
}
```

3. Add the language switcher:
```tsx
import { LanguageSwitcher } from '@/components/LanguageSwitcher';

<LanguageSwitcher />
```

### Translation Keys Structure
```json
{
  "common": { ... },      // Common UI elements
  "nav": { ... },         // Navigation
  "auth": { ... },        // Authentication
  "signup": { ... },      // Sign up page
  "home": { ... },        // Home page
  "profile": { ... },     // Profile page
  "events": { ... },      // Events
  "credentials": { ... }, // Credentials
  "portfolio": { ... },   // Portfolio
  "validation": { ... },  // Form validation
  "errors": { ... }       // Error messages
}
```

### Language Detection
Language is detected and stored in the following order:
1. **localStorage** - User preference saved as `decm-language`
2. **Browser navigator** - Browser language settings
3. **HTML tag** - Document language

## 🔐 Security

### PII Encryption
All personally identifiable information (PII) is encrypted at the database level using PostgreSQL's `pgcrypto` extension.

```sql
-- Encryption pattern
INSERT INTO profiles (encrypted_email, encrypted_name)
VALUES (
    pgp_sym_encrypt(sqlc.arg(email), sqlc.arg(encryption_key)::varchar),
    pgp_sym_encrypt(sqlc.arg(name), sqlc.arg(encryption_key)::varchar)
);
```

### Authentication
- **Wallet-based authentication** with message signing
- **JWT sessions** via HTTP-only cookies
- **OAuth integration** (Google)
- Protected routes with authentication middleware

## 📚 Core Features

### 1. Event Management
- Multi-type events (academic, social, professional)
- NFT-based ticketing system
- Event creation and management

### 2. Digital Credentials
- Blockchain-based certificates
- QR code verification
- Certificate issuance and validation

### 3. Academic Identity
- LDAP integration for institutional verification
- "Verify Once, Use Everywhere" principle
- Academic credential management

### 4. e-Portfolio
- Personal achievement showcase
- Credential collection
- Public profile sharing

### 5. Evaluation System
- Reputation management
- Feedback system
- Participant evaluation

## 🔗 API Documentation

Access the interactive Swagger documentation at:
```
http://localhost:8080/swagger/
```

The API client is automatically generated from OpenAPI specs:
```typescript
import { DefaultApi } from '@decm/api';

const api = new DefaultApi({ 
  basePath: 'http://localhost:8080/api/v1' 
});

// Use generated methods
const response = await api.getProfile();
```

## 🗄️ Database

### Migrations
Database migrations are automatically run when starting the backend with `pnpm dev:core`.

Manual migration:
```bash
pnpm db:migrate
```

### Query Development
1. Write SQL queries in `packages/database/queries/`
2. Generate Go code: `pnpm db:generate`
3. Use generated types in repositories

## 🎨 UI Components

Built with **Radix UI** and **Tailwind CSS** for accessible, customizable components.

```tsx
import { Button } from '@/components/ui/button';
import { Dialog } from '@/components/ui/dialog';
import { Select } from '@/components/ui/select';
```

## 🧪 Development Workflow

### API-First Development
1. Define Go handler with OpenAPI annotations
2. Generate TypeScript client: `pnpm gen-api:core`
3. Use generated client in React frontend
4. Type safety maintained end-to-end

### Database-First Development
1. Write SQL queries in `packages/database/queries/`
2. Generate Go code: `pnpm db:generate`
3. Use generated types in repositories and handlers

## 📝 Environment Configuration

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
ENCRYPTION_KEY=your-encryption-key
```

## 🤝 Contributing

1. Follow the code standards in `.rules/code-standards`
2. Use conventional commits
3. Write OpenAPI documentation for all API endpoints
4. Ensure type safety with sqlc and TypeScript
5. Use pnpm for package management

## 📄 License

[Add your license here]

## 🔗 Links

- [API Documentation](http://localhost:8080/swagger/)
- [Frontend](http://localhost:3000)
- [Database Migrations](packages/database/migrations/)

---

**Built with ❤️ for Web 3.0 and decentralized identity management**
