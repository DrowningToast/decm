# DECM Production Deployment Guide

This guide covers deploying DECM using Docker Compose in a production environment.

## Prerequisites

- Docker Engine 20.10+ and Docker Compose v2.0+
- At least 4GB RAM and 20GB disk space
- Domain name configured (optional, for SSL)
- Environment variables configured in `.env` file

## Quick Start

### 1. Configure Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

**CRITICAL**: Update these variables in `.env`:

```bash
# Security (MUST CHANGE IN PRODUCTION)
PII_ENCRYPTION_KEY=your-secure-32-character-key-here
JWT_SECRET=your-secure-jwt-secret-key-here

# Database
DB_NAME=decm
DB_USER=postgres
DB_PASSWORD=your-secure-database-password

# Blockchain
BLOCKCHAIN_NETWORK=mainnet
BLOCKCHAIN_RPC_URL=https://your-ethereum-node
BLOCKCHAIN_PRIVATE_KEY=your-private-key
BLOCKCHAIN_CHAIN_ID=1
BLOCKCHAIN_DECM_ACCESS_MANAGER_ADDRESS=0x...

# Google OAuth (if using)
GOOGLE_OAUTH_CLIENT_ID=your-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
GOOGLE_OAUTH_REDIRECT_URL=https://yourdomain.com/api/v1/auth/verify-google-oauth
GOOGLE_OAUTH_SUCCESS_URL=https://yourdomain.com/auth/success

# S3 Storage
S3_ACCESS_KEY_ID=your-access-key
S3_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET_NAME=your-bucket
S3_ENDPOINT=https://s3.amazonaws.com

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com

# Frontend URLs
VITE_CORE_BACKEND_API=https://api.yourdomain.com
VITE_APP_URL=https://yourdomain.com
VITE_ENVIRONMENT=production
VITE_GOOGLE_MAPS_API_KEY=your-google-maps-api-key
VITE_WALLETCONNECT_PROJECT_ID=your-walletconnect-project-id

# Nginx Proxy Manager (optional)
NPM_ADMIN_EMAIL=admin@yourdomain.com
NPM_ADMIN_PASSWORD=secure-admin-password
```

### 2. Build and Start Services

**Without Nginx Proxy Manager (recommended for simpler setups):**

```bash
docker compose -f docker-compose.prod.yml up -d
```

**With Nginx Proxy Manager (for SSL/reverse proxy):**

```bash
docker compose -f docker-compose.prod.yml --profile proxy up -d
```

### 3. Verify Services

Check all services are running:

```bash
docker compose -f docker-compose.prod.yml ps
```

Check logs:

```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f postgres
```

### 4. Access the Application

- **Frontend**: http://localhost:3000 (or your configured port)
- **Backend API**: http://localhost:8080 (or your configured port)
- **Backend Health Check**: http://localhost:8080/ready
- **Nginx Proxy Manager** (if enabled): http://localhost:81

## Service Architecture

The production setup includes:

1. **PostgreSQL Database** (`postgres`)
   - Persistent data storage with healthchecks
   - Automatic initialization
   - Volume: `postgres_data`

2. **Backend API** (`backend`)
   - Go-based REST API
   - Auto-runs database migrations on startup
   - Waits for database to be healthy
   - Healthcheck on `/ready` endpoint

3. **Frontend** (`frontend`)
   - React 19 SPA served by Nginx
   - Static files with optimized caching
   - Waits for backend to be healthy

4. **Nginx Proxy Manager** (`proxy`, optional)
   - SSL/TLS termination
   - Reverse proxy
   - Only starts with `--profile proxy`

## Port Configuration

Default ports (configurable via `.env`):

- Frontend: `3000` → `80` (container)
- Backend: `8080` → `8080` (container)
- Database: `5432` → `5432` (container)
- NPM HTTP: `80` → `80` (container)
- NPM HTTPS: `443` → `443` (container)
- NPM Admin: `81` → `81` (container)

## Database Management

### Access PostgreSQL CLI

```bash
docker compose -f docker-compose.prod.yml exec postgres psql -U postgres -d decm
```

### Backup Database

```bash
docker compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres decm > backup.sql
```

### Restore Database

```bash
cat backup.sql | docker compose -f docker-compose.prod.yml exec -T postgres psql -U postgres -d decm
```

## Scaling and Updates

### Update Services

1. Pull latest code:
   ```bash
   git pull origin main
   ```

2. Rebuild images:
   ```bash
   docker compose -f docker-compose.prod.yml build
   ```

3. Restart services:
   ```bash
   docker compose -f docker-compose.prod.yml up -d
   ```

### Zero-downtime Updates

```bash
# Build new images
docker compose -f docker-compose.prod.yml build

# Recreate services one by one
docker compose -f docker-compose.prod.yml up -d --no-deps backend
docker compose -f docker-compose.prod.yml up -d --no-deps frontend
```

## SSL/HTTPS Setup with Nginx Proxy Manager

If you started the proxy service:

1. Access NPM admin UI at http://your-server-ip:81
2. Login with configured credentials (`NPM_ADMIN_EMAIL` / `NPM_ADMIN_PASSWORD`)
3. Add a Proxy Host:
   - **Domain Names**: yourdomain.com, www.yourdomain.com
   - **Scheme**: http
   - **Forward Hostname/IP**: frontend
   - **Forward Port**: 80
   - **SSL**: Request a new SSL certificate (Let's Encrypt)

4. Add API proxy:
   - **Domain Names**: api.yourdomain.com
   - **Scheme**: http
   - **Forward Hostname/IP**: backend
   - **Forward Port**: 8080
   - **SSL**: Request a new SSL certificate

## Monitoring and Logs

### View Real-time Logs

```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service with timestamps
docker compose -f docker-compose.prod.yml logs -f --timestamps backend
```

### Check Resource Usage

```bash
docker stats
```

### Health Checks

```bash
# Backend health
curl http://localhost:8080/ready

# Frontend health
curl http://localhost:3000/

# Database health
docker compose -f docker-compose.prod.yml exec postgres pg_isready -U postgres
```

## Troubleshooting

### Backend fails to start

**Check database connection:**
```bash
docker compose -f docker-compose.prod.yml logs postgres
docker compose -f docker-compose.prod.yml logs backend
```

**Verify environment variables:**
```bash
docker compose -f docker-compose.prod.yml exec backend env | grep DB_
```

### Frontend shows API errors

**Check CORS configuration:**
- Ensure `CORS_ALLOWED_ORIGINS` in `.env` includes your frontend domain
- Verify `VITE_CORE_BACKEND_API` points to the correct backend URL

### Database migrations fail

**Manually run migrations:**
```bash
# Access backend container
docker compose -f docker-compose.prod.yml exec backend sh

# Inside container, run migrations manually
# (The backend auto-runs migrations on startup)
```

### Out of disk space

**Clean up Docker resources:**
```bash
docker system prune -a --volumes
```

**Check volume sizes:**
```bash
docker system df -v
```

## Security Best Practices

1. **Change default passwords**: Update all passwords in `.env`
2. **Use strong encryption keys**: Generate secure random keys for `PII_ENCRYPTION_KEY` and `JWT_SECRET`
3. **Enable SSL**: Use Nginx Proxy Manager or configure SSL certificates
4. **Firewall rules**: Only expose necessary ports (80, 443)
5. **Regular backups**: Automate database backups
6. **Update regularly**: Keep Docker images and code up to date
7. **Environment variables**: Never commit `.env` to version control
8. **Database access**: Restrict PostgreSQL port (5432) to internal network only

## Production Checklist

- [ ] `.env` file configured with production values
- [ ] `PII_ENCRYPTION_KEY` is a secure 32-character key
- [ ] `JWT_SECRET` is a strong random secret
- [ ] Database password changed from default
- [ ] CORS origins configured for production domain
- [ ] Frontend URLs updated for production domain
- [ ] S3 storage configured and tested
- [ ] Blockchain configuration verified
- [ ] SSL certificates configured (if using proxy)
- [ ] Firewall rules configured
- [ ] Backup strategy implemented
- [ ] Monitoring/alerting configured

## Shutting Down

### Stop all services

```bash
docker compose -f docker-compose.prod.yml down
```

### Stop and remove volumes (WARNING: deletes all data)

```bash
docker compose -f docker-compose.prod.yml down -v
```

## Support

For issues and questions:
- Check logs: `docker compose -f docker-compose.prod.yml logs`
- Review environment configuration in `.env`
- Consult main README.md and CLAUDE.md for development details
