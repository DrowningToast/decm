# 🚀 Quick Database Setup Guide

## For New Developers

### **1-Minute Setup**
```bash
# Clone and install
git clone <repo>
cd decm
bun install

# Database setup
cp .env.example .env
# Edit .env with your database preferences (optional - defaults work!)
bun db:setup

# Done! Your database is ready 🎉
```

### **Verify Setup**
```bash
# Check configuration
bun db:config

# View database tables  
bun db:status

# Check migration status
bun db:migrate:version
```

## **What Just Happened?**

The `bun db:setup` command automatically:
1. ✅ **Started PostgreSQL** in Docker container
2. ✅ **Waited for database** to be ready (health checks)  
3. ✅ **Read credentials** from your .env file
4. ✅ **Ran migrations** (extensions + core tables)
5. ✅ **Verified tables** were created successfully

## **Available Tables**

After setup, you'll have these tables:
- `authentication_credentials` - User authentication with public/private keys
- `profiles` - User profile data with privacy controls
- `events` - Event management with blockchain integration  
- `event_attendees` - Event participation tracking
- `event_certificates` - Certificate issuance system

## **Next Steps**

```bash
# Generate Go code for database queries
bun db:generate

# Start the backend server
bun backend:dev

# Start the frontend
bun dev
```

## **Common Commands**

```bash
# Database operations
bun db:start              # Start database
bun db:stop               # Stop database
bun db:reset              # Reset all data
bun db:console            # Open database console

# Development
bun backend:dev           # Start API server  
bun dev                   # Start frontend
bun gen-api              # Update API client
```

## **Configuration**

Default `.env` settings (work out of the box):
```env
DB_HOST=localhost
DB_PORT=5432  
DB_USER=decm_user
DB_PASSWORD=decm_password
DB_NAME=decm
DB_SSL_MODE=disable
```

**You can customize these values in your .env file!**

---
**That's it! You're ready to develop with DECM! 🎯**
