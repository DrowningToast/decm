# DECM Backend API Server

Backend for Frontend (BFF) API server for the DECM platform, built with Go Fiber.

## Quick Start

### Prerequisites
- Go 1.21+ installed
- Copy `env.template` from root to `.env` in root directory

### Development

From the **root directory**:

```bash
# Start backend in development mode
bun backend:dev

# Or directly with npm/yarn
npm run backend:dev
```

This will:
- Load environment variables from root `.env` file
- Start the Go Fiber server on port 8080 (default)
- Enable hot-reload during development

### Production

```bash
# Build the backend
bun backend:build

# Start the built binary
bun backend:start
```

## API Endpoints

- `GET /` - Health check and service info
- `GET /api/v1/health` - API health check

## Configuration

The backend loads configuration from the root `.env` file:

- **Server**: Port, environment settings
- **CORS**: Frontend origins configuration  
- **Database**: PostgreSQL connection settings
- **Blockchain**: Web3/blockchain integration settings
- **JWT**: Authentication token configuration
- **LDAP**: Academic identity verification settings

## Project Structure

```
apps/backend/
├── cmd/               # Application entry points
│   └── main.go       # Main server application
├── internal/         # Private application code
│   ├── config/       # Configuration management
│   ├── handlers/     # HTTP request handlers
│   ├── middleware/   # Custom middleware
│   ├── models/       # Data models and types
│   ├── services/     # Business logic services
│   └── utils/        # Utility functions
├── api/              # API definitions
│   ├── routes/       # Route definitions
│   └── v1/           # API version 1
└── pkg/              # Public packages
```

## Features

- **CORS Support**: Configured for frontend integration
- **Request Logging**: Structured logging with request IDs
- **Error Handling**: Centralized error handling
- **Environment Config**: Flexible configuration management
- **Monorepo Integration**: Integrated with Turbo build system
