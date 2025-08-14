# DECM - Decentralized Event Management Platform

> A Web 3.0 and Blockchain-powered platform for event management, NFT ticketing, digital credentials, and academic identity verification.

[![Built with Bun](https://img.shields.io/badge/Built%20with-Bun-ff69b4?logo=bun)](https://bun.sh/)
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
- **Package Manager**: Bun for all operations
- **Monorepo**: Turbo + Bun workspaces

### 🗄️ Database Design
- **Type-safe Operations**: sqlc generates Go code from SQL
- **Migration System**: Version-controlled schema changes
- **Full-text Search**: PostgreSQL trigram matching
- **UUID Primary Keys**: Blockchain-compatible identifiers
- **JSON Support**: Flexible metadata storage

### 🔒 Security Features
- **SQL Injection Prevention**: Prepared statements for all queries
- **Context Timeouts**: Prevents hanging database operations
- **Connection Pooling**: Efficient database resource management
- **Environment-based Config**: Secure credential management

## 🚀 Quick Start

### Prerequisites
- **Bun**: `curl -fsSL https://bun.sh/install | bash`
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
bun install

# Setup environment
cp .env.example .env
# Edit .env with your database credentials
```

### 2. Start Database
```bash
# Start PostgreSQL with Docker Compose
bun compose:up

# Database will be ready at localhost:5432
```

### 3. Start Development Servers

#### Backend (with auto-migrations)
```bash
# Start Go API server - automatically runs migrations
bun backend:dev

# Available at: http://localhost:8080
# API docs at: http://localhost:8080/swagger/
```

#### Frontend
```bash
# Start React development server
bun dev

# Available at: http://localhost:3000
```

### 4. Generate TypeScript API Client
```bash
# Generate type-safe API client from OpenAPI spec
bun gen-api

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
bun dev              # Start React dev server
bun build            # Build for production
bun lint             # Run ESLint
bun check-types      # TypeScript checking

# Backend development
bun backend:dev      # Start API with auto-migrations
bun backend:build    # Build Go binary
bun backend:docs     # Generate OpenAPI documentation
```

### 🗄️ Database Operations
```bash
# Code generation
bun db:generate      # Generate Go code from SQL queries

# Migration management
bun db:migrate:up    # Run pending migrations
bun db:migrate:down  # Rollback last migration
bun db:migrate:create name  # Create new migration

# Docker database
bun compose:up       # Start PostgreSQL container
bun compose:down     # Stop database
```

### 🔄 API Generation
```bash
# Full API client generation pipeline
bun gen-api          # OpenAPI docs → TypeScript client

# Individual steps
bun backend:docs     # Generate OpenAPI specification
cd packages/api && bun generate  # Generate TypeScript client
```

## 🗄️ Database Schema

### Core Tables

#### 👤 Users
- Academic identity with LDAP verification
- Institution email validation  
- Web3 wallet integration
- Full-text search capabilities

#### 🎯 Events
- Multi-type event support
- NFT contract integration
- Capacity and location management
- Organizer relationships

#### 🎫 NFT Tickets  
- Blockchain token tracking
- QR code check-in system
- Transaction hash recording
- Seat assignments

#### 🏆 Credentials
- Digital certificates and badges
- Blockchain verification
- QR code validation  
- Skills and criteria tracking

#### 📋 Portfolios
- Personal achievement collections
- Public/private sharing
- Token-based access control

#### ⭐ Evaluations
- Multi-criteria rating system
- Anonymous feedback support
- Reputation scoring

## 🌐 API Documentation

### Health Check
- `GET /` - System health with database status
- `GET /api/v1/health` - API health check

### User Management
- `POST /api/v1/users` - Create user account
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user profile  
- `GET /api/v1/users/search?q=query` - Search users
- `DELETE /api/v1/users/{id}` - Delete user account

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

### Test User API
```bash
# Create a new user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@university.edu",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "academic_institution": "University of Technology"
  }'

# Search users
curl "http://localhost:8080/api/v1/users/search?q=john&limit=10"
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
DB_USER=postgres
DB_PASSWORD=your_password
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
- [x] Database schema and migrations
- [x] User management API
- [x] Authentication system
- [x] Basic frontend structure

### Phase 2: Event Management 🚧
- [ ] Event CRUD operations
- [ ] Event search and filtering
- [ ] Organizer dashboard
- [ ] Event categories and tags

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
3. **Follow** the coding standards and use Bun for all operations
4. **Test** your changes: `bun test` (when tests are implemented)
5. **Commit** your changes: `git commit -m 'Add amazing feature'`
6. **Push** to the branch: `git push origin feature/amazing-feature`
7. **Open** a Pull Request

### Development Setup
```bash
# Install Bun (required)
curl -fsSL https://bun.sh/install | bash

# Install dependencies
bun install

# Start development environment
bun compose:up      # Database
bun backend:dev     # API server  
bun dev             # Frontend
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