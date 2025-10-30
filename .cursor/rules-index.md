# Cursor Rules Index

## 🎯 Core Rules (Always Applied)

### [core-conventions.mdc](.cursor/rules/core-conventions.mdc)

**Compact reference** for DECM conventions covering:

- pnpm package manager requirements
- Backend (Go) architecture patterns
- Frontend (React/TypeScript) structure
- Essential commands
- Security & validation basics

## 📋 Available Rules by Topic

### Backend Development

- `go-backend-architecture.mdc` - Handler, usecase, repository layers
- `database-migrations.mdc` - PostgreSQL patterns
- `database-security.mdc` - PII encryption & security
- `repository-patterns.mdc` - Data access layer
- `error-handling.mdc` - Error handling patterns

### Frontend Development

- `ui-components.mdc` - Radix UI + Tailwind patterns
- `form-patterns.mdc` - React Hook Form integration
- `page-components-pattern.mdc` - Routes vs components
- `i18n-translations.mdc` - Translation patterns
- `typography-usage.mdc` - Text component guidelines

### Code Quality & Standards

- `code-standards.mdc` - Unified Go/TypeScript standards
- `eslint-configuration.mdc` - TypeScript linting rules
- `testing-setup.mdc` - Vitest + React Testing Library
- `testing-conventions.mdc` - Testing patterns
- `coderabbit.mdc` - PR review standards

### API & Infrastructure

- `api-generation.mdc` - OpenAPI → TypeScript client
- `environment-config.mdc` - .env management
- `development-workflow.mdc` - Git workflow & processes

### Web3 & Authentication

- `wallet-integration.mdc` - Web3 wallet patterns
- `authentication-security.mdc` - JWT & auth flows
- `onboarding-authentication-flows.mdc` - Complete flows
- `smart-contracts.mdc` - NFT contracts

## 🚀 Quick Access

Run in Cursor:

```
fetch_rules(["core-conventions"])  # Start here
```

For specific topics, use:

```
fetch_rules(["<topic>-<subtopic>"])  # e.g., "database-security"
```

---

**Start with `core-conventions.mdc` for a quick overview!**
