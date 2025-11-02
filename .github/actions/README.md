# GitHub Composite Actions

This directory contains reusable composite actions for the DECM CI/CD workflows.

## Available Actions

### 1. setup-frontend

**File**: `setup-frontend/action.yml`

Complete frontend environment setup including Node.js, pnpm, environment files, and dependencies.

**Usage**:

```yaml
- name: Setup Frontend Environment
  uses: ./.github/actions/setup-frontend
  with:
      node-version: "20"
      pnpm-version: "9.15.0"
```

**Includes**:

- Setup pnpm
- Setup Node.js with cache
- Copy `.env.test` → `.env`
- Copy `.env.client.test` → `.env.client`
- Run `pnpm install --frozen-lockfile`

**Used in**: `frontend-eslint`, `frontend-build`, `frontend-tests`

---

### 2. setup-node-pnpm

**File**: `setup-node-pnpm/action.yml`

Basic Node.js and pnpm setup for backend jobs (single .env file).

**Usage**:

```yaml
- name: Setup Node.js and pnpm
  uses: ./.github/actions/setup-node-pnpm
  with:
      node-version: "20"
      pnpm-version: "9.15.0"
```

**Includes**:

- Setup pnpm
- Setup Node.js with cache
- Copy `.env.test` → `.env`
- Run `pnpm install --frozen-lockfile`

**Used in**: `docker-compose`, `database-migrate`, `database-generate`, `backend-build`, `backend-tests`, `backend-openapi`

---

### 3. setup-go-backend

**File**: `setup-go-backend/action.yml`

Go environment setup for backend compilation and testing.

**Usage**:

```yaml
- name: Setup Go Backend
  uses: ./.github/actions/setup-go-backend
  with:
      go-version: "1.23"
```

**Includes**:

- Setup Go with cache
- Install Go dependencies (`go mod download`)
- Setup Go workspace (use packages, sync)

**Used in**: `database-generate`, `backend-build`, `backend-tests`, `backend-openapi`

---

### 4. install-db-tools

**File**: `install-db-tools/action.yml`

Install database migration and code generation tools.

**Usage**:

```yaml
# Install both migrate and sqlc
- name: Install Database Tools
  uses: ./.github/actions/install-db-tools

# Install only migrate
- name: Install Database Tools
  uses: ./.github/actions/install-db-tools
  with:
      install-sqlc: "false"

# Custom versions
- name: Install Database Tools
  uses: ./.github/actions/install-db-tools
  with:
      migrate-version: "v4.17.0"
      sqlc-version: "v1.29.0"
```

**Includes**:

- Install golang-migrate
- Install sqlc (optional)

**Used in**: `database-migrate`, `database-generate`, `backend-build`, `backend-tests`

---

## Creating New Actions

### Structure

```
.github/actions/
└── your-action-name/
    └── action.yml
```

### Template

```yaml
name: "Action Name"
description: "Action description"

inputs:
    input-name:
        description: "Input description"
        required: true
        default: "default-value"

runs:
    using: "composite"
    steps:
        - name: Step Name
          shell: bash
          run: |
              echo "Your commands here"
```

### Requirements

1. **Shell directive**: All `run` steps must specify `shell: bash`
2. **Inputs**: Define all configurable parameters
3. **Defaults**: Provide sensible defaults when possible
4. **Documentation**: Add clear descriptions
5. **Single responsibility**: Keep actions focused

### Best Practices

- ✅ Accept versions as inputs (don't hardcode)
- ✅ Use environment variables for configuration
- ✅ Keep actions small and reusable
- ✅ Test with different configurations
- ✅ Document all inputs and outputs
- ❌ Don't checkout code (let workflows handle it)
- ❌ Don't install unnecessary tools
- ❌ Don't hardcode values that might change

## Benefits

### Code Reuse

- **Before**: 20+ lines of setup code per job
- **After**: 2-5 lines per job
- **Result**: ~70% reduction in workflow file size

### Maintainability

- Update setup logic in one place
- Changes automatically apply to all jobs
- Consistent setup across all workflows

### Performance

- Actions are cached by GitHub
- No overhead compared to inline steps
- Parallel jobs still benefit from caching

## Version History

### v1.0 (Current)

- `setup-frontend`: Frontend environment with dual .env files
- `setup-node-pnpm`: Basic Node.js setup for backend
- `setup-go-backend`: Go workspace and dependencies
- `install-db-tools`: Database tooling (migrate + sqlc)

## Related Documentation

- [CI/CD Workflow](../../.cursor/rules/ci-cd-workflow.mdc)
- [Development Workflow](../../.cursor/rules/development-workflow.mdc)
- [PR Checks Workflow](../workflows/pr-checks.yml)
