# @decm/api

Generated TypeScript API client for the DECM (Decentralized Event Management) backend.

## Overview

This package contains automatically generated TypeScript interfaces and API client code based on the OpenAPI/Swagger specification from the DECM Go Fiber backend.

## Usage

```typescript
import { /* API classes and types */ } from '@decm/api';

// Use the generated API client
const api = new DefaultApi();
```

## Generation

The TypeScript code is automatically generated from the backend's OpenAPI specification:

```bash
# Generate API client from backend OpenAPI spec
bun gen-api
```

This will:
1. Start the backend server to ensure latest OpenAPI spec
2. Generate TypeScript client code using OpenAPI Generator
3. Post-process and organize the generated files
4. Build the TypeScript to JavaScript

## Development

```bash
# Install dependencies
bun install

# Generate API client
bun generate

# Build TypeScript
bun build

# Watch mode for development
bun dev
```

## Package Structure

```
packages/api/
├── src/           # Source TypeScript files
├── dist/          # Built JavaScript files
├── generated/     # Raw generated files (temporary)
├── scripts/       # Post-processing scripts
└── package.json
```

## Features

- **Type Safety**: Full TypeScript support with generated interfaces
- **Fetch-based**: Uses modern fetch API for HTTP requests
- **Tree Shakeable**: ES modules with selective imports
- **Auto-generated**: Always in sync with backend API changes
- **Monorepo Ready**: Designed for use across the DECM monorepo

## API Documentation

The API documentation is available at `http://localhost:8080/swagger/` when the backend is running.
