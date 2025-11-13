# Docker Production Setup

This document describes the production Docker setup for the DECM platform.

## Architecture

The production environment consists of 4 Docker containers:

1. **PostgreSQL** - Database (postgres:16-alpine)
2. **Backend API** - Go Fiber application (custom image)
3. **Frontend Web** - React 19 + Nginx (custom image)
4. **Nginx Proxy Manager** - Reverse proxy with SSL/TLS management

## Quick Start

### 1. Configure Environment

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

**Critical variables to set:**

```bash
# Security (REQUIRED)
PII_ENCRYPTION_KEY=your-32-byte-encryption-key-here
JWT_SECRET=your-jwt-secret-min-32-chars-here

# Database
DB_PASSWORD=strong-database-password-here

# Production URLs
CORS_ALLOWED_ORIGINS=https://yourdomain.com
VITE_API_URL=https://api.yourdomain.com/api/v1
GOOGLE_OAUTH_REDIRECT_URL=https://api.yourdomain.com/api/v1/auth/verify-google-oauth
GOOGLE_OAUTH_SUCCESS_URL=https://yourdomain.com/auth/success

# OAuth (Get from Google Cloud Console)
GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret

# Storage (AWS S3 or compatible)
S3_ACCESS_KEY_ID=your-s3-access-key
S3_SECRET_ACCESS_KEY=your-s3-secret-key
S3_BUCKET_NAME=your-bucket-name
S3_ENDPOINT=https://s3.amazonaws.com
S3_REGION=us-east-1

# Blockchain (Optional)
BLOCKCHAIN_RPC_URL=your-ethereum-rpc-url
BLOCKCHAIN_PRIVATE_KEY=your-private-key
```

### 2. Build Images

Build both frontend and backend Docker images:

```bash
# Build backend
docker build -f apps/backend/Dockerfile -t decm-backend:latest .

# Build frontend
docker build -f apps/web/Dockerfile -t decm-frontend:latest .
```

### 3. Start Production Stack

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 4. Configure Nginx Proxy Manager

1. Access Nginx Proxy Manager admin UI:
   - URL: http://your-server-ip:81
   - Default login: `admin@example.com` / `changeme`
   - **IMPORTANT**: Change default password immediately!

2. Add Proxy Hosts:
   - **Frontend**: `yourdomain.com` → `http://frontend:80`
   - **Backend API**: `api.yourdomain.com` → `http://backend:8080`

3. Enable SSL/TLS:
   - Request Let's Encrypt certificates
   - Force SSL
   - Enable HTTP/2

## Container Details

### Backend Container

- **Image**: `decm-backend:latest`
- **Port**: 8080 (internal)
- **Health Check**: `http://localhost:8080/ready`
- **Features**:
  - Multi-stage build (builder + runtime)
  - Non-root user execution
  - Alpine-based (minimal size)
  - Auto-migrations on startup

### Frontend Container

- **Image**: `decm-frontend:latest`
- **Port**: 80 (internal)
- **Web Server**: Nginx 1.27
- **Health Check**: `http://localhost:80/health`
- **Features**:
  - Multi-stage build (Node builder + Nginx runtime)
  - Non-root user execution
  - Gzip compression enabled
  - React Router support (SPA)
  - Static asset caching

### PostgreSQL Container

- **Image**: `postgres:16-alpine`
- **Port**: 5432
- **Volumes**: 
  - `postgres_data` - Database files
  - `./database/backups` - Backup directory

### Nginx Proxy Manager

- **Image**: `jc21/nginx-proxy-manager:latest`
- **Ports**:
  - 80 - HTTP
  - 443 - HTTPS
  - 81 - Admin UI
- **Features**:
  - Web-based UI for proxy configuration
  - Let's Encrypt SSL certificate automation
  - Access lists and security

## Docker Compose Commands

```bash
# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Stop all services
docker-compose -f docker-compose.prod.yml down

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# View specific service logs
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f frontend

# Restart a service
docker-compose -f docker-compose.prod.yml restart backend

# Pull latest images
docker-compose -f docker-compose.prod.yml pull

# Rebuild images
docker-compose -f docker-compose.prod.yml build --no-cache

# Remove volumes (WARNING: deletes data!)
docker-compose -f docker-compose.prod.yml down -v
```

## CI/CD Integration

The GitHub Actions workflow `.github/workflows/docker-build.yml` automatically:

1. Builds Docker images on push to `main`, `develop`, or `production/*` branches
2. Pushes images to GitHub Container Registry (ghcr.io)
3. Tags images with:
   - Branch name
   - Git commit SHA
   - `latest` for main branch
   - Semantic version for tagged releases (v1.0.0)

### Using CI-Built Images

Pull images from GitHub Container Registry:

```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull images
docker pull ghcr.io/your-org/decm/backend:latest
docker pull ghcr.io/your-org/decm/frontend:latest

# Tag for local use
docker tag ghcr.io/your-org/decm/backend:latest decm-backend:latest
docker tag ghcr.io/your-org/decm/frontend:latest decm-frontend:latest
```

## Production Deployment Checklist

- [ ] Set strong `PII_ENCRYPTION_KEY` (32 bytes)
- [ ] Set strong `JWT_SECRET` (min 32 characters)
- [ ] Set strong `DB_PASSWORD`
- [ ] Configure Google OAuth credentials
- [ ] Configure S3 storage
- [ ] Set production domain in `CORS_ALLOWED_ORIGINS`
- [ ] Set production API URL in `VITE_API_URL`
- [ ] Configure Nginx Proxy Manager SSL certificates
- [ ] Change Nginx Proxy Manager default password
- [ ] Set up database backups
- [ ] Configure firewall rules
- [ ] Set up monitoring and logging
- [ ] Test health checks
- [ ] Test frontend access
- [ ] Test backend API
- [ ] Test OAuth login flow

## Monitoring

### Health Checks

All containers have health checks configured:

```bash
# Check container health status
docker ps --format "table {{.Names}}\t{{.Status}}"

# Manual health check
curl http://localhost:8080/ready  # Backend
curl http://localhost:3000/health # Frontend
```

### Logs

View logs in real-time:

```bash
# All services
docker-compose -f docker-compose.prod.yml logs -f

# Specific service with tail
docker-compose -f docker-compose.prod.yml logs -f --tail=100 backend
```

## Backup and Restore

### Database Backup

```bash
# Create backup
docker exec decm-postgres-prod pg_dump -U decm_user decm > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore from backup
docker exec -i decm-postgres-prod psql -U decm_user decm < backup_20250113_120000.sql
```

### Volume Backup

```bash
# Backup volumes
docker run --rm -v decm_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_data_backup.tar.gz /data
```

## Security Best Practices

1. **Secrets Management**:
   - Never commit `.env` files
   - Use environment variable injection or secrets manager
   - Rotate keys regularly

2. **Network Security**:
   - Use internal Docker network for service communication
   - Only expose necessary ports
   - Configure firewall rules

3. **SSL/TLS**:
   - Always use HTTPS in production
   - Enable HSTS headers
   - Use modern TLS versions only

4. **Container Security**:
   - Run as non-root user
   - Keep base images updated
   - Scan images for vulnerabilities

5. **Database Security**:
   - Use strong passwords
   - Enable SSL connections
   - Regular backups
   - PII encryption at application layer

## Troubleshooting

### Backend won't start

```bash
# Check logs
docker logs decm-backend-prod

# Common issues:
# - Database connection failure: Check DB_HOST, DB_PASSWORD
# - Missing environment variables: Check .env file
# - Migration errors: Check database/migrations/
```

### Frontend returns 404 for routes

```bash
# Verify nginx.conf is properly configured for React Router
docker exec decm-frontend-prod cat /etc/nginx/nginx.conf

# Check for try_files configuration
```

### Database connection issues

```bash
# Test database connectivity
docker exec decm-backend-prod pg_isready -h postgres -U decm_user

# Check network
docker network inspect decm_decm-network
```

### SSL certificate issues

1. Check Nginx Proxy Manager logs
2. Verify domain DNS points to server
3. Ensure ports 80/443 are accessible
4. Check Let's Encrypt rate limits

## Performance Tuning

### PostgreSQL

```bash
# Adjust shared_buffers, work_mem in postgres container
# Add to docker-compose.prod.yml:
command:
  - "postgres"
  - "-c"
  - "shared_buffers=256MB"
  - "-c"
  - "max_connections=200"
```

### Nginx

Nginx configuration in `apps/web/nginx.conf`:
- Gzip compression enabled
- Static asset caching (1 year)
- Client max body size: 20MB

### Docker Resources

Limit container resources:

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

## Scaling

For horizontal scaling:

1. Use external PostgreSQL service (RDS, Cloud SQL)
2. Deploy multiple backend/frontend containers
3. Use external load balancer or Nginx Proxy Manager
4. Share session store (Redis)
5. Use external file storage (S3)

## Support

For issues or questions:
- Check logs: `docker-compose logs`
- Review health checks
- Verify environment configuration
- Consult main README.md for development setup
