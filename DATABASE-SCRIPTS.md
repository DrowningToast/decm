# Database Migration Scripts Documentation

## 📋 Overview

This document describes the comprehensive database migration system that automates all database operations while reading configuration securely from environment files.

## 🎯 Key Features

- ✅ **Dynamic Configuration**: Reads database credentials from `.env` file
- ✅ **Error Handling**: Automatic recovery from dirty migration states  
- ✅ **Wait Mechanism**: Waits for database to be ready before operations
- ✅ **Security**: No hardcoded credentials in package.json
- ✅ **Flexibility**: Works with both .env file and environment variables
- ✅ **Comprehensive**: Complete database lifecycle management

## 🔧 Available Scripts

### **Core Database Operations**
```bash
# Complete database setup (start + wait + migrate)
bun db:setup

# Start/stop database container
bun db:start    # Start PostgreSQL with Docker Compose  
bun db:stop     # Stop database container

# Check database configuration
bun db:config   # Display current database configuration
```

### **Migration Management**
```bash
# Run migrations (with error recovery)
bun db:migrate              # Smart migration with dirty state handling
bun db:migrate:down         # Rollback migrations  
bun db:migrate:version      # Check current migration version
bun db:migrate:force        # Force migration to specific version
bun db:migrate:create name  # Create new migration files

# Reset database (down + up)
bun db:reset               # Complete reset of database
```

### **Database Tools**
```bash
# Code generation
bun db:generate            # Generate Go code from SQL queries (sqlc)

# Database inspection
bun db:status              # Show all database tables
bun db:console             # Open PostgreSQL interactive console
bun db:wait                # Wait for database to be ready
```

## 🔐 Environment Configuration

### **Configuration Priority**
1. **Environment variables** (highest priority)
2. **`.env` file** (if environment variables not set)
3. **`.env.example`** (fallback with warnings)

### **Required Environment Variables**
```env
# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password  
DB_NAME=decm
DB_SSL_MODE=disable
```

### **Setup Steps**
```bash
# 1. Copy environment template
cp .env.example .env

# 2. Edit with your database credentials
# Edit .env file with your actual database settings

# 3. Verify configuration
bun db:config

# 4. Start database and run migrations
bun db:setup
```

## 🚀 Complete Setup Process

### **From Scratch Setup**
```bash
# 1. Install dependencies
bun install

# 2. Configure environment
cp .env.example .env
# Edit .env with your database credentials

# 3. Complete database setup
bun db:setup
# This will:
#   - Start PostgreSQL container
#   - Wait for database to be ready  
#   - Run all pending migrations
#   - Verify tables were created

# 4. Generate Go code (if using sqlc)
bun db:generate
```

### **Development Workflow**
```bash
# Start database
bun db:start

# Create new migration
bun db:migrate:create add_new_table

# Edit migration files in packages/database/migrations/
# Then run migrations
bun db:migrate

# Generate updated Go code
bun db:generate

# Check database status
bun db:status
```

## 🛠️ Script Architecture

### **Core Components**

#### **1. `scripts/db-env.js`**
- **Purpose**: Environment configuration management
- **Features**: 
  - Reads `.env` file with fallback to `.env.example`
  - Supports environment variable overrides
  - Constructs PostgreSQL connection URL
  - Validates configuration completeness

#### **2. `scripts/db-migrate.js`**  
- **Purpose**: Intelligent migration runner
- **Features**:
  - Automatic dirty state recovery
  - Error handling and retry logic
  - Progress reporting and verification
  - Final status confirmation

#### **3. `scripts/wait-for-db.js`**
- **Purpose**: Database readiness checker
- **Features**:
  - Docker container health checking
  - PostgreSQL connection polling
  - Configurable timeout and retry logic
  - Graceful failure handling

#### **4. `scripts/db-command.js`**
- **Purpose**: Universal database command runner
- **Features**:
  - Dynamic credential injection
  - Support for migrate and psql commands
  - Consistent error handling
  - Command validation and help

## 🔍 Example Usage

### **Check Configuration**
```bash
$ bun db:config
📄 Using .env file for database configuration
🔧 Database Configuration:
   Host: localhost
   Port: 5432
   User: decm_user
   Database: decm
   SSL Mode: disable
   Source: .env file

🔗 Database URL:
postgres://decm_user:decm_password@localhost:5432/decm?sslmode=disable
```

### **Complete Setup**
```bash
$ bun db:setup
[+] Running 1/1
 ✔ Container decm-postgres  Started

🔍 Waiting for PostgreSQL database to be ready...
📄 Using .env file for database configuration
⏳ Polling database every 2s (max 30 attempts)...
   Attempt 1/30...
   🐳 Container not ready yet...
   Attempt 2/30...
   🔄 Database not ready yet...
   Attempt 3/30...
✅ Database is ready!

🗄️  Running database migrations...
📄 Using .env file for database configuration
⬆️  Attempting to run migrations...
📋 Running: migrate -path ./migrations -database [URL] up
1/u enable_extensions (12.345678ms)
2/u create_core_tables (45.678901ms)
✅ Migrations completed successfully!
📊 Current migration version: 2
🔍 Verifying database tables...
```

### **Migration with Error Recovery**
```bash
$ bun db:migrate
🗄️  Running database migrations...
⬆️  Attempting to run migrations...
⚠️  Migration failed, checking for dirty state...
🧹 Detected dirty database state, attempting to clean...
🔄 Forcing to version 1...
⬆️  Retrying migrations after cleanup...
✅ Migrations completed successfully after cleanup!
```

## ⚠️ Troubleshooting

### **Common Issues**

#### **Configuration Problems**
```bash
# Check configuration
bun db:config

# Common fixes:
cp .env.example .env  # Create .env file
# Edit .env with correct credentials
```

#### **Database Connection Issues**
```bash
# Ensure database is running
bun db:start

# Wait for database to be ready
bun db:wait

# Check Docker container status
docker ps | grep decm-postgres
```

#### **Migration Problems**
```bash
# Check current version
bun db:migrate:version

# Force clean state (if needed)
bun db:migrate:force 0
bun db:migrate

# Complete reset
bun db:reset
```

#### **Permission Issues**
```bash
# Make scripts executable
chmod +x scripts/*.js

# Check database permissions
bun db:console
# Then: \du to see users and permissions
```

### **Debug Commands**
```bash
# Show all database tables
bun db:status

# Check migration version
bun db:migrate:version

# Interactive database console
bun db:console

# Container logs
docker logs decm-postgres
```

## 📊 Migration Process Flow

```
1. Read .env configuration
2. Start database container (if needed)
3. Wait for database readiness
4. Check current migration version
5. Run pending migrations
6. Handle dirty state if needed
7. Verify tables were created
8. Generate Go code (if requested)
```

## 🔒 Security Considerations

- ✅ **No hardcoded credentials** in package.json
- ✅ **Environment-based configuration** from .env file
- ✅ **Secure credential handling** in scripts
- ✅ **SSL mode configuration** support
- ✅ **Connection string isolation** in helper functions

## 🎯 Best Practices

1. **Always use `.env` file** for local development
2. **Use environment variables** in production
3. **Run `bun db:config`** to verify settings before operations
4. **Use `bun db:setup`** for complete initial setup
5. **Test migrations** with `bun db:reset` before committing
6. **Keep `.env.example` updated** with latest required variables

This system provides a robust, secure, and maintainable approach to database operations while ensuring credentials are never hardcoded in the codebase.
