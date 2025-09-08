# DECM - Decentralized Event Management Platform

> A Web 3.0 and Blockchain-powered platform for event management, NFT ticketing, digital credentials, and academic identity verification.

[![Built with pnpm](https://img.shields.io/badge/Built%20with-pnpm-yellow?logo=pnpm)](https://pnpm.io/)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue?logo=postgresql)](https://postgresql.org/)
[![Go](https://img.shields.io/badge/Backend-Go%20Fiber-00ADD8?logo=go)](https://gofiber.io/)
[![React](https://img.shields.io/badge/Frontend-React%2019-61DAFB?logo=react)](https://react.dev/)

## 🌟 Platform Features

### 🎫 Event Management & NFT Ticketing
- **Multi-type Events**: Student activities, faculty events, training seminars, public events
- **NFT-based Tickets**: Unique digital tickets with blockchain verification
- **QR Code Check-in**: Seamless event entry with QR code scanning
- **Smart Contracts**: Immutable ticket ownership and transfer

### 🎓 Academic Identity Verification (DAI)
- **LDAP Integration**: University authentication and verification
- **"Verify Once, Use Everywhere"**: Cross-platform identity validation  
- **Institution Email Verification**: Academic status confirmation
- **Decentralized Identity**: User-owned academic credentials

### 🏆 Digital Credentials & Badges
- **Blockchain Certificates**: Anti-forgery digital credentials
- **Automated Issuance**: Event completion certificates
- **QR Verification**: Instant credential authenticity checking
- **Skills Tracking**: Competency and achievement recording

### 📋 e-Portfolio System
- **Achievement Portfolios**: Personal credential and event collections
- **Public/Private Sharing**: Controlled access to achievements
- **QR/Link Sharing**: Easy sharing with employers and institutions
- **Portfolio Analytics**: Track engagement and views

### ⭐ Evaluation & Reputation
- **Multi-criteria Ratings**: Content, organization, venue assessments
- **Anonymous Feedback**: Optional anonymous evaluations
- **Reputation Scoring**: Lecturer and organization reputation tracking
- **Real-time Analytics**: Live feedback processing

## 🏗️ Technical Architecture

### 🎯 Full-Stack Monorepo
- **Frontend**: React 19 + Tailwind CSS + Radix UI
- **Backend**: Go Fiber BFF API with OpenAPI documentation
- **Database**: PostgreSQL with sqlc type-safe queries
- **Package Manager**: pnpm for all operations
- **Monorepo**: Turbo + pnpm workspaces

### 🗄️ Database Design
- **Type-safe Operations**: sqlc generates Go code from SQL
- **Migration System**: Version-controlled schema changes with automatic recovery
- **Environment Configuration**: No hardcoded credentials, reads from .env
- **Blockchain Integration**: Contract addresses, chain IDs, and key management
- **Privacy Controls**: Granular privacy settings for user data

### 🔒 Security Features
- **SQL Injection Prevention**: Prepared statements for all queries
- **Context Timeouts**: Prevents hanging database operations
- **Connection Pooling**: Efficient database resource management
- **Environment-based Config**: Secure credential management

## 🚀 Quick Start

### Prerequisites
- **pnpm**: `npm install -g pnpm` or `curl -fsSL https://get.pnpm.io/install.sh | sh -`
- **Go**: 1.24.1 or higher
- **Docker**: For PostgreSQL database
- **sqlc**: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- **migrate**: `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### 1. Setup Environment
```bash
# Clone repository
git clone <repository-url>
cd decm

# Install dependencies  
pnpm install

# Setup environment
cp .env.example .env
# Edit .env with your database credentials
```

### 2. Start Database
```bash
# Start PostgreSQL with Docker Compose
pnpm compose:up

# Database will be ready at localhost:5432
```

### 3. Start Development Servers

#### Backend (with auto-migrations)
```bash
# Start Go API server - automatically runs migrations
pnpm backend:dev

# Available at: http://localhost:8080
# API docs at: http://localhost:8080/swagger/
```

#### Frontend
```bash
# Start React development server
pnpm dev

# Available at: http://localhost:3000
```

### 4. Generate TypeScript API Client
```bash
# Generate type-safe API client from OpenAPI spec
pnpm gen-api

# Creates packages/api/ with TypeScript types
```

## 📁 Project Structure

```
decm/
├── 📱 apps/
│   ├── 🌐 web/                    # React 19 Frontend
│   │   ├── src/
│   │   │   ├── index.tsx          # App entry point
│   │   │   ├── components/ui/     # Radix UI components
│   │   │   └── lib/               # Utility functions
│   │   └── package.json
│   └── ⚙️ backend/                # Go Fiber API
│       ├── cmd/main.go            # Server entry point
│       ├── api/                   # Route handlers
│       ├── internal/              # Private application code
│       └── go.mod
├── 📦 packages/
│   ├── 🗄️ database/              # Database package
│   │   ├── migrations/            # SQL migration files
│   │   ├── queries/               # SQL query definitions  
│   │   ├── go/generated/          # Generated Go code
│   │   └── connection.go          # Database utilities
│   ├── 🌐 api/                    # Generated TypeScript client
│   └── ⚙️ typescript-config/      # Shared TypeScript configs
├── 🐳 docker-compose.isolated.yml # PostgreSQL container
├── 📋 package.json                # Root scripts and dependencies
└── ⚡ turbo.json                  # Monorepo build config
```

## 🛠️ Available Commands

### 🎯 Development
```bash
# Frontend development
pnpm dev              # Start React dev server
pnpm build            # Build for production
pnpm lint             # Run ESLint
pnpm check-types      # TypeScript checking

# Backend development
pnpm backend:dev      # Start API with auto-migrations
pnpm backend:build    # Build Go binary
pnpm backend:docs     # Generate OpenAPI documentation
```

### 🗄️ Database Operations
```bash
# Code generation
pnpm db:generate      # Generate Go code from SQL queries

# Migration management
pnpm db:migrate:up    # Run pending migrations
pnpm db:migrate:down  # Rollback last migration
pnpm db:migrate:create name  # Create new migration

# Docker database
pnpm compose:up       # Start PostgreSQL container
pnpm compose:down     # Stop database
```

### 🔄 API Generation
```bash
# Full API client generation pipeline
pnpm gen-api          # OpenAPI docs → TypeScript client

# Individual steps
pnpm backend:docs     # Generate OpenAPI specification
cd packages/api && pnpm generate  # Generate TypeScript client
```

## 🗄️ Database Schema

### Core Tables

#### 🔐 Authentication Credentials
- **BYOK Support**: Bring Your Own Key or system-managed keys
- **OAuth Integration**: Google and GitHub connector support
- **Verification System**: Organizer and student verification flags
- **Public Key Infrastructure**: Blockchain-compatible key management

#### 👤 User Profiles
- **Privacy Controls**: Granular privacy settings for each field
- **Academic Integration**: Institution and academic email verification
- **Personal Data**: Contact information with privacy controls
- **Profile Customization**: Bio, profile pictures, and display preferences

#### 🎯 Events
- **Blockchain Integration**: Contract addresses and chain ID support
- **Requirement System**: Configurable attendee requirements
- **Verification**: Event verification and validation system
- **Location Management**: Google Maps integration and location tracking

#### 🎫 Event Participation
- **Attendee Tracking**: Registration and acceptance status
- **Data Provision**: Track what information users provide
- **Contact Management**: Blockchain contact addresses for participants

#### 🏆 Event Certificates
- **Publication Control**: Manage certificate publication status
- **Event Linking**: Certificates linked to specific events and users
- **Issuance Tracking**: Track certificate creation and distribution

## 🌐 API Documentation

### Health Check
- `GET /` - System health with database status
- `GET /api/v1/health` - API health check

### Authentication & User Management
- `POST /api/v1/authentication-credentials` - Create authentication credential
- `GET /api/v1/authentication-credentials/{id}` - Get credential by ID
- `PUT /api/v1/authentication-credentials/{id}` - Update credential
- `GET /api/v1/profiles` - Get user profiles
- `POST /api/v1/profiles` - Create user profile

### Interactive Documentation
When the backend is running, visit:
- **Swagger UI**: http://localhost:8080/swagger/
- **OpenAPI JSON**: http://localhost:8080/docs/swagger.json

## 🧪 Testing the Integration

### Test Database Connection
```bash
# Check system health with database status
curl http://localhost:8080/

# Should return: {"database": true, "status": "healthy", ...}
```

### Test Authentication API
```bash
# Create authentication credential
curl -X POST http://localhost:8080/api/v1/authentication-credentials \
  -H "Content-Type: application/json" \
  -d '{
    "solution_status": 1,
    "public_key": "your_public_key",
    "is_verified_organizer": 0,
    "is_verified_student": 1
  }'

# Check database tables
pnpm db:status
```

## 🔧 Environment Configuration

Copy `.env.example` to `.env` and configure:

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=decm_user
DB_PASSWORD=decm_password
DB_NAME=decm
DB_SSL_MODE=disable

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000

# JWT Configuration  
JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION=24h

# Blockchain Configuration
BLOCKCHAIN_NETWORK=localhost
BLOCKCHAIN_RPC_URL=http://localhost:8545
BLOCKCHAIN_CHAIN_ID=1337

# LDAP Configuration (Academic Identity)
LDAP_HOST=ldap.university.edu
LDAP_PORT=389
LDAP_BASE_DN=dc=university,dc=edu
```

## 📚 Documentation

- **[Database Integration Guide](DATABASE-INTEGRATION.md)** - Complete database setup and usage
- **[Migration Integration Summary](MIGRATION-INTEGRATION.md)** - Migration system details  
- **[API Generation Guide](API-GENERATION.md)** - TypeScript client generation
- **[Project Proposal](project-proposal.md)** - Original platform concept and requirements

## 🤝 Development Guidelines

### Code Standards
- **TypeScript**: Strict mode enabled with proper typing
- **Go**: Standard formatting with `gofmt` and `golint`  
- **SQL**: Well-documented migrations with rollback support
- **React**: Functional components with hooks, no unsafe HTML

### Database Rules
- **Migrations**: Always create up/down pairs with safety checks
- **Queries**: Use sqlc annotations for type-safe code generation
- **Context**: All operations must use context with timeouts
- **Transactions**: Wrap related operations in database transactions

### API Standards  
- **REST**: Follow RESTful conventions for endpoints
- **Validation**: Validate all inputs with proper error messages
- **Documentation**: All endpoints must have OpenAPI annotations
- **Error Handling**: Consistent JSON error response format

## 🎯 Target Users

- **🎓 Students**: Event participation, credential collection, identity verification
- **👨‍🏫 Faculty/Staff**: Event creation, credential issuance, academic verification  
- **🏢 Organizations**: Event hosting, participant verification, credential validation
- **🌐 External Users**: Public event access, professional credentialing

## 🚧 Future Roadmap

### Phase 1: Core Platform ✅
- [x] Database schema and migrations with blockchain support
- [x] Authentication credentials system with BYOK support
- [x] Environment-based configuration system
- [x] Smart migration system with error recovery
- [x] Basic frontend structure

### Phase 2: User Profiles & Event Management 🚧
- [ ] User profile CRUD operations with privacy controls
- [ ] Event CRUD operations with blockchain integration
- [ ] Event search and filtering
- [ ] Organizer dashboard with verification system
- [ ] Event categories and requirement management

### Phase 3: NFT Ticketing 📋
- [ ] Smart contract integration
- [ ] NFT minting for tickets
- [ ] QR code generation
- [ ] Check-in system

### Phase 4: Digital Credentials 📋
- [ ] Certificate templates
- [ ] Automated issuance
- [ ] Blockchain verification
- [ ] QR code validation

### Phase 5: Portfolio System 📋
- [ ] Portfolio creation
- [ ] Public sharing
- [ ] Analytics dashboard
- [ ] Export functionality

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Follow** the coding standards and use pnpm for all operations
4. **Test** your changes: `pnpm test` (when tests are implemented)
5. **Commit** your changes: `git commit -m 'Add amazing feature'`
6. **Push** to the branch: `git push origin feature/amazing-feature`
7. **Open** a Pull Request

### Development Setup
```bash
# Install pnpm (required)
npm install -g pnpm

# Install dependencies
pnpm install

# Start development environment
pnpm compose:up      # Database
pnpm backend:dev     # API server  
pnpm dev             # Frontend
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

- **Documentation**: Check the `/docs` directory for detailed guides
- **API Reference**: http://localhost:8080/swagger/ when backend is running
- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and community support

---

<div align="center">

**Built with ❤️ for the future of decentralized event management**

[🌟 Star this repo](https://github.com/your-org/decm) • [🐛 Report Bug](https://github.com/your-org/decm/issues) • [💡 Request Feature](https://github.com/your-org/decm/issues)

</div>