# Local GitHub Actions Testing with Act

This document explains how to test GitHub Actions workflows locally using [act](https://github.com/nektos/act).

## Prerequisites

1. **Install act**
    - macOS: `brew install act`
    - Linux/Windows: See [act installation guide](https://github.com/nektos/act#installation)

2. **Docker** must be running on your machine

## Quick Start

### List Available Jobs

```bash
pnpm act:list
# or
act pull_request --list
```

### Run All Jobs

```bash
pnpm act
# or
act pull_request
```

### Run Specific Job

```bash
# Using pnpm scripts
pnpm act:docker-compose
pnpm act:backend-build
pnpm act:frontend-build

# Using the helper script
node scripts/act-run.js docker-compose
node scripts/act-run.js backend-build -v
pnpm act:job docker-compose

# Direct act command (with workflow directory specified)
act pull_request -W .github/workflows/ --job backend-build --secret-file .env.test
```

## Available Scripts

All scripts use `pnpm`:

- `pnpm act` - Run all jobs (simulates pull_request event)
- `pnpm act:list` - List all available jobs
- `pnpm act:job [job-name]` - Run specific job using helper script
- `pnpm act:docker-compose` - Run docker-compose job
- `pnpm act:database-migrate` - Run database migration job
- `pnpm act:database-generate` - Run database generation job
- `pnpm act:backend-build` - Run backend build job
- `pnpm act:backend-tests` - Run backend tests job
- `pnpm act:backend-openapi` - Run OpenAPI generation job
- `pnpm act:frontend-build` - Run frontend build job
- `pnpm act:frontend-tests` - Run frontend tests job
- `pnpm act:contracts-build` - Run contracts build job

## Configuration

### `.actrc` File

The `.actrc` file in the project root configures act defaults:

- Uses `catthehacker/ubuntu:act-latest` runner image
- Sets workflow directory to `.github/workflows/` (only scans workflow files)
- Enables verbose output

**Note**: The workflow directory is explicitly set to avoid act scanning other YAML files (like `.coderabbit.yaml`).

### Secrets File

Act uses the `.env.test` file in the project root for environment variables and secrets. This file is committed to the repository and contains test-specific configurations.

The workflow automatically loads secrets from `.env.test` when running with act.

## Job Dependencies

The workflow has the following dependency chain:

```
docker-compose (no dependencies)
  └── database-migrate
        └── database-generate
              └── backend-build
                    ├── backend-tests (also needs database-migrate)
                    └── backend-openapi

frontend-build (no dependencies)
  └── frontend-tests

contracts-build (no dependencies)
```

To test a specific job, act will automatically handle dependencies based on the `needs:` declarations in the workflow.

## Common Options

### Verbose Output

```bash
act pull_request --job backend-build -v
```

### List Steps Without Running

```bash
act pull_request --job backend-build --list
```

### Use Specific Event

```bash
act push --secret-file .env.test
act pull_request --secret-file .env.test
```

### Dry Run (Show What Would Run)

```bash
act pull_request --dryrun --secret-file .env.test
```

## Limitations

1. **GitHub Services**: Some GitHub Actions features don't work locally (e.g., `services:` in jobs)
    - For PostgreSQL, act uses a local service container
    - Docker Compose jobs require Docker to be available

2. **File System**: Act runs in containers, so paths are relative to the container

3. **Secrets**: Use `.env.act` for local secrets (never commit this file)

4. **Performance**: Local runs may be slower than GitHub Actions runners

## Troubleshooting

### Act Not Found

```bash
# macOS
brew install act

# Or download from releases
# https://github.com/nektos/act/releases
```

### Docker Not Running

```bash
# Make sure Docker is running
docker ps
```

### Permission Issues

```bash
# Make sure scripts are executable
chmod +x scripts/act-run.js
```

### Job Fails Locally But Works on GitHub

- Check that all dependencies are installed
- Verify environment variables are set correctly
- Ensure Docker is accessible from within act containers
- Some actions may behave differently locally

### Error: "Unknown Property" or "workflow is not valid"

If you see errors like:

```
Error: workflow is not valid. '.coderabbit.yaml': Line: 1 Column 1: Unknown Property reviews
```

This means act is trying to parse non-workflow YAML files. The configuration is already set to only scan `.github/workflows/`, but if you're running act directly without the flags, make sure to include:

```bash
act pull_request -W .github/workflows/ [other-options]
```

Or use the provided scripts which handle this automatically:

```bash
pnpm act:job [job-name]
```

## Best Practices

1. **Test Individual Jobs**: Start by testing individual jobs before running the full workflow
2. **Use Verbose Mode**: Use `-v` flag to see detailed output
3. **Check Logs**: Review act output carefully for errors
4. **Clean Up**: Act may leave containers running - use `docker ps` and clean up if needed
5. **Commit Often**: Test workflows locally before pushing to catch issues early

## References

- [act GitHub Repository](https://github.com/nektos/act)
- [act Documentation](https://github.com/nektos/act#overview)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
