# GitHub Actions Workflows

## 📁 Workflow Structure

The CI/CD pipeline is split into modular, reusable workflows for better maintainability:

```
.github/workflows/
├── pr-checks-new.yml    # Main orchestrator (lightweight)
├── database.yml         # Database migration & SQLC generation
├── backend.yml          # Backend build, tests, OpenAPI
├── frontend.yml         # Frontend lint, build, tests
├── docker.yml           # Docker builds + startup tests ⭐ NEW
└── contracts.yml        # Smart contract builds
```

## 🚀 What's New

### Service Startup Test (docker.yml)

**Problem it solves:** Catches startup crashes before production deployment

The new startup test:

- ✅ Builds Docker images for backend and frontend
- ✅ Starts services using `docker-compose.ci.yml`
- ✅ Waits for PostgreSQL to be healthy
- ✅ Waits for Backend `/ready` endpoint to respond
- ✅ Waits for Frontend to be accessible
- ✅ Verifies all services can communicate
- ✅ **Catches missing environment variables** before prod
- ✅ **Detects startup crashes** early

**Real-world example:** Would have caught the `panic: failed to open log file` error we encountered.

## 📋 Workflow Details

### 1. pr-checks-new.yml (Main Orchestrator)

**Purpose:** Coordinates all checks and enforces label-based gating

**Flow:**

```
check-label
    ├─> database-checks (parallel)
    ├─> backend-checks (parallel)
    ├─> frontend-checks (parallel)
    └─> contracts-checks (parallel)
              ↓
         docker-checks (waits for backend + frontend)
              ↓
         pr-checks-summary
```

**Key Features:**

- Requires `ready-to-review` label to run
- Removes label automatically on failure
- Posts detailed failure comment on PR

### 2. database.yml

**Jobs:**

- `migrate` - Tests database migrations
- `generate` - Validates SQLC generated code

**Dependencies:** PostgreSQL service

### 3. backend.yml

**Jobs:**

- `build` - Compiles Go binary
- `tests` - Runs unit tests with DB
- `openapi` - Validates OpenAPI docs

**Dependencies:** PostgreSQL for tests

### 4. frontend.yml

**Jobs:**

- `eslint` - Linting checks
- `build` - Builds production bundle
- `tests` - Unit tests with coverage

**Dependencies:** None

### 5. docker.yml ⭐ NEW

**Jobs:**

- `build-backend` - Builds backend Docker image
- `build-frontend` - Builds frontend Docker image
- `startup-test` - **Tests service startup** (NEW!)

**Startup Test Steps:**

1. Load Docker images from artifacts
2. Start services with docker-compose
3. Wait for PostgreSQL health check
4. Wait for Backend health check (`/ready`)
5. Wait for Frontend health check
6. Test Backend endpoints (200 OK)
7. Test Frontend accessibility (200 OK)
8. Verify all services running
9. Show logs on failure
10. Cleanup resources

**Timeout:** 10 minutes

### 6. contracts.yml

**Jobs:**

- `build` - Compiles Solidity contracts with Foundry

**Dependencies:** None

## 🔧 Configuration Files

### docker-compose.ci.yml

Minimal docker-compose for CI testing:

- PostgreSQL with test database
- Backend with all required env vars
- Frontend with test build args
- Health checks for all services
- Uses non-conflicting ports (5433, 8081, 3002)

**Key Features:**

- All required environment variables set (catches missing vars)
- Mock values for external services (OAuth, S3, Blockchain)
- Fast startup with minimal dependencies
- Comprehensive health checks

## 📊 Environment Variables

### Required for Backend Startup

```bash
# Core
ENVIRONMENT=production  # Using production for maximum compatibility
PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=decm_ci
DB_PASSWORD=***
DB_NAME=decm_ci

# Security (CI test values)
PII_ENCRYPTION_KEY=test-encryption-key-for-ci-only-32bytes-minimum
JWT_SECRET=test-jwt-secret-key-for-ci-only
JWT_ISSUER=decm-service-ci
JWT_EXPIRATION=24h
COOKIE_DOMAIN=localhost

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Blockchain (mock)
BLOCKCHAIN_NETWORK=sepolia
BLOCKCHAIN_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/test
BLOCKCHAIN_PRIVATE_KEY=0x0000...0001
BLOCKCHAIN_CHAIN_ID=11155111
BLOCKCHAIN_ETHERSCAN_API_KEY=test-api-key
BLOCKCHAIN_DECM_ACCESS_MANAGER_ADDRESS=0x0000...0000

# OAuth (mock)
GOOGLE_OAUTH_CLIENT_ID=test-client-id
GOOGLE_OAUTH_CLIENT_SECRET=test-client-secret
GOOGLE_OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/authentication/oauth/google/callback
GOOGLE_OAUTH_SUCCESS_URL=http://localhost:3000/auth/success

# S3 Storage (mock)
S3_ACCESS_KEY_ID=test-access-key
S3_SECRET_ACCESS_KEY=test-secret-key
S3_BUCKET_NAME=test-bucket
S3_ENDPOINT=http://localhost:9000
```

## 🎯 Benefits of New Structure

### Maintainability

- ✅ Each workflow is focused on one area
- ✅ Easy to find and update specific checks
- ✅ Reduced duplication with reusable workflows
- ✅ Clear separation of concerns

### Performance

- ✅ Parallel execution where possible
- ✅ Artifacts shared between jobs
- ✅ GitHub Actions cache for Docker layers

### Reliability

- ✅ **Startup test catches crashes before production**
- ✅ Health checks verify services are actually working
- ✅ Detailed logs on failure for debugging
- ✅ Graceful cleanup even on failure

### Developer Experience

- ✅ Clear status for each area (Database, Backend, Frontend, Docker, Contracts)
- ✅ Fast feedback on what failed
- ✅ Easy to re-run individual workflows
- ✅ Detailed PR comments on failure

## 🚦 Migration Guide

### Before (Old Structure)

```yaml
# Single 652-line pr-checks.yml with:
- 11 jobs
- Lots of duplication
- Hard to maintain
- No startup testing
```

### After (New Structure)

```yaml
# Modular structure with:
- 6 workflow files
- 1 main orchestrator (140 lines)
- 5 reusable workflows
- NEW: Service startup test
- Easy to maintain
```

### To Adopt

1. **Backup old workflow:**

    ```bash
    mv .github/workflows/pr-checks.yml .github/workflows/pr-checks.old.yml
    ```

2. **Activate new workflow:**

    ```bash
    mv .github/workflows/pr-checks-new.yml .github/workflows/pr-checks.yml
    ```

3. **Test on a PR:**
    - Create a test PR
    - Add `ready-to-review` label
    - Verify all checks pass

4. **Remove old workflow once confirmed:**
    ```bash
    rm .github/workflows/pr-checks.old.yml
    ```

## 🐛 Debugging Failed Checks

### Database Checks Failed

- Check migration syntax
- Verify SQLC generated code is committed
- Check `packages/database/go/generated/`

### Backend Checks Failed

- Check Go compilation errors
- Verify tests pass locally
- Check OpenAPI docs are up to date

### Frontend Checks Failed

- Check ESLint errors
- Verify build succeeds locally
- Check test coverage

### Docker Checks Failed (Startup Test) ⭐

- **Check backend logs** - shows startup errors
- **Verify environment variables** - missing vars cause crashes
- **Check health endpoints** - `/ready` should return 200
- **Database connection** - verify DB is accessible
- **Port conflicts** - CI uses 5433, 8081, 3001

Common issues:

```bash
# Missing env var
Error: required environment variable "PII_ENCRYPTION_KEY" is not set

# Health check timeout
Backend is starting... (waiting for /ready endpoint)
❌ Backend container exited unexpectedly!

# Database connection failed
Error: failed to connect to database: connection refused
```

### Contracts Checks Failed

- Check Foundry compilation
- Verify submodules are initialized
- Check contract syntax

## 📈 Monitoring

### Workflow Execution Time

**Before:** ~15-20 minutes (sequential in many places)
**After:** ~12-15 minutes (more parallelization)

**Startup Test:** Adds ~3-5 minutes but catches critical issues

### Success Rate

The startup test helps maintain high PR quality by catching:

- Environment variable misconfigurations
- Service startup crashes
- Health endpoint failures
- Database connectivity issues
- Docker image build problems

## 🔄 Future Improvements

Potential additions:

- [ ] E2E tests after startup test
- [ ] Performance benchmarks
- [ ] Security scanning (Snyk, Trivy)
- [ ] Load testing
- [ ] Smoke tests against staging

## 📝 Notes

- All workflows use GitHub Actions cache for speed
- Docker layer caching reduces build times
- Artifacts are shared between jobs when needed
- Startup test has 10-minute timeout
- All services must pass health checks
- Logs are saved on failure for debugging

---

**Last Updated:** 2026-01-09
**Maintainer:** DevOps Team
