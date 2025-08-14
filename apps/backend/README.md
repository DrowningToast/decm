# DECM Backend API

Go Fiber-based Backend for Frontend (BFF) API for the DECM platform with integrated database support.

## Features

- 🚀 **Go Fiber** - Fast HTTP framework
- 🗄️ **Database Integration** - PostgreSQL with type-safe queries (sqlc)
- 🔄 **Auto-migrations** - Database schema management
- 📚 **OpenAPI/Swagger** - Comprehensive API documentation
- 🔄 **Auto-generation** - TypeScript client generation
- 🛡️ **CORS Support** - Configurable cross-origin resource sharing
- 🔍 **Request ID Tracking** - Request tracing and logging
- ⚡ **Hot Reload** - Development mode with auto-restart
- 🔐 **Environment Security** - No hardcoded credentials

## Architecture

This backend serves as a **Backend for Frontend (BFF)** with full database integration:
- Direct database access with type-safe queries
- Authentication credential management with BYOK support
- User profile system with privacy controls
- Event management with blockchain integration
- Automatic migrations and error recovery
- Generates TypeScript interfaces automatically

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL (via Docker or local)
- Environment file: Copy from root `.env.example` to `.env`

### Development
```bash
# Complete database setup (from project root)
bun db:setup

# Start development server with database integration
bun backend:dev

# Server will be available at: http://localhost:8080
# API documentation: http://localhost:8080/swagger/
```

### Build
```bash
# Build binary (from project root)
bun backend:build

# Or build directly  
go build -o bin/decm-backend cmd/main.go

# Run built binary
./bin/decm-backend
```

## Project Structure

```
apps/backend/
├── cmd/
│   └── main.go          # Application entry point with DB integration
├── api/
│   ├── routes/          # Route definitions
│   └── v1/              # API v1 handlers
│       └── users.go     # User management endpoints (example)
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Custom middleware
│   ├── services/        # Business logic services
│   │   └── database.go  # Database service integration
│   └── utils/           # Utility functions
├── docs/                # Generated OpenAPI documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod               # Go module definition with database deps
└── README.md
```

## Database Integration

### Current Schema
- **authentication_credentials** - User auth with BYOK and OAuth support
- **profiles** - User profiles with privacy controls
- **events** - Event management with blockchain integration
- **event_attendees** - Participation tracking
- **event_certificates** - Certificate issuance system

### Type-Safe Queries
All database operations use generated type-safe queries:
```go
// Generated from SQL queries via sqlc
user, err := h.DB.DB.Queries.GetAuthenticationCredentialByID(ctx, id)
```

### Migration System
- Automatic migrations in development mode
- Manual migration control: `bun db:migrate`
- Error recovery for dirty states
- Version tracking and rollback support

## Configuration

Configuration is loaded from environment variables (via `.env` file in project root):

```env
# Server
PORT=8080
ENVIRONMENT=development

# Database (PostgreSQL) - Required!
DB_HOST=localhost
DB_PORT=5432
DB_USER=decm_user
DB_PASSWORD=decm_password
DB_NAME=decm
DB_SSL_MODE=disable

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Security
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Blockchain Integration
BLOCKCHAIN_NETWORK=localhost
BLOCKCHAIN_RPC_URL=http://localhost:8545
BLOCKCHAIN_CHAIN_ID=1337
```

## API Documentation

When the server is running, interactive documentation is available at:
- **Swagger UI**: http://localhost:8080/swagger/
- **OpenAPI JSON**: http://localhost:8080/docs/swagger.json
- **OpenAPI YAML**: http://localhost:8080/docs/swagger.yaml

## Available Endpoints

### Health & Status
- `GET /` - Root health check with database status
- `GET /api/v1/health` - API health check

### Authentication & User Management (Example Implementation)
- `POST /api/v1/users` - Create user (example from old schema)
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /api/v1/users` - List users (paginated)
- `GET /api/v1/users/search` - Search users

### Future Endpoints (Prepared Schema)
The database schema is ready for:
- **Authentication Credentials**: CRUD with BYOK support
- **User Profiles**: Privacy-controlled profile management
- **Event Management**: Blockchain-integrated events
- **Event Participation**: Attendee tracking and requirements
- **Certificate System**: Issuance and publication management

## Development

### Database Workflow
1. **Create migrations**: `bun db:migrate:create description`
2. **Edit migration files**: `packages/database/migrations/`
3. **Run migrations**: `bun db:migrate`
4. **Create queries**: `packages/database/queries/table_name.sql`
5. **Generate Go code**: `bun db:generate`
6. **Use in handlers**: Import generated types and methods

### Adding New Endpoints
1. **Create SQL queries** in `packages/database/queries/`
2. **Generate Go code**: `bun db:generate`
3. **Create handler function** in `internal/handlers/` or `api/v1/`
4. **Add Swagger annotations** for documentation
5. **Register route** in `api/routes/routes.go`
6. **Generate documentation**: `bun backend:docs`
7. **Update TypeScript client**: `bun gen-api`

Example handler with database integration:
```go
// @Summary Create authentication credential
// @Description Create new user authentication credential
// @Tags authentication
// @Accept json
// @Produce json
// @Param credential body CreateCredentialRequest true "Credential data"
// @Success 201 {object} AuthenticationCredential
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/credentials [post]
func createAuthenticationCredential(h *handlers.Handlers) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
        defer cancel()
        
        var req CreateCredentialRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }
        
        credential, err := h.DB.DB.Queries.CreateAuthenticationCredential(ctx, generated.CreateAuthenticationCredentialParams{
            SolutionStatus:      req.SolutionStatus,
            PublicKey:           req.PublicKey,
            IsVerifiedOrganizer: req.IsVerifiedOrganizer,
            IsVerifiedStudent:   req.IsVerifiedStudent,
        })
        
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to create credential"})
        }
        
        return c.Status(201).JSON(credential)
    }
}
```

### Database Service Pattern
All database operations go through the DatabaseService:
```go
// In handlers
func NewHandlers(cfg *config.Config, dbService *services.DatabaseService) *Handlers {
    return &Handlers{
        config: cfg,
        DB:     dbService,  // Database service with connection and queries
    }
}
```

### Error Handling
Database-integrated error responses:
```json
{
  "error": "Failed to create authentication credential", 
  "code": 500,
  "request_id": "req_123456"
}
```

## Integration

### With Database
Full integration with PostgreSQL:
- Type-safe queries generated from SQL
- Automatic connection pooling
- Migration management
- Health check integration

### With Frontend
The backend provides seamless frontend integration:
- CORS configured for development
- TypeScript types auto-generated from database schema
- Request/response formats optimized for React

### With Blockchain
Architecture ready for Web3 integration:
- Contract addresses in events table
- Chain ID tracking
- Public key infrastructure
- BYOK (Bring Your Own Key) support

## Security Features

- **No hardcoded credentials** - All config from environment
- **Type-safe database operations** - Prevents SQL injection
- **Connection pooling** - Efficient resource management
- **Context timeouts** - Prevents hanging operations
- **Environment validation** - Ensures proper configuration

## Testing

### Database Testing
```bash
# Test database connection
bun db:config

# Verify tables
bun db:status

# Test migration system
bun db:reset
```

### API Testing
```bash
# Health check with database status
curl http://localhost:8080/

# Test database-backed endpoints
curl http://localhost:8080/api/v1/users
```

## Deployment

### Production Checklist
- [ ] Set `ENVIRONMENT=production` 
- [ ] Configure production database credentials
- [ ] Set strong `JWT_SECRET`
- [ ] Configure production CORS origins
- [ ] Run database migrations
- [ ] Build optimized binary

### Production Build
```bash
# Build optimized binary
go build -ldflags="-w -s" -o bin/decm-backend cmd/main.go

# Set environment
export ENVIRONMENT=production

# Run with production database
./bin/decm-backend
```

This backend provides a complete foundation for the DECM platform with full database integration, type-safe operations, and blockchain-ready architecture.