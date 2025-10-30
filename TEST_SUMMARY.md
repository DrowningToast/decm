# Frontend Unit Tests & CI Pipeline Summary

## ✅ Completed Tasks

### 1. GitHub Actions CI Workflow (`.github/workflows/pr-checks.yml`)

Created a comprehensive CI pipeline that runs on PRs to `main` with the following checks:

#### Database & SQLC Checks

- Sets up PostgreSQL service
- Installs sqlc and golang-migrate
- Runs database migrations
- Generates SQLC code and verifies it matches committed files

#### Backend Checks

- Builds Core API binary
- Runs backend unit tests with coverage
- Verifies Go dependencies

#### OpenAPI & Swagger Checks

- Generates OpenAPI documentation from Go code
- Verifies documentation matches committed files
- Ensures API contracts are up to date

#### Frontend Checks

- Runs TypeScript type checking
- Runs ESLint for code quality
- **Runs frontend unit tests** (NEW!)
- Builds production bundle
- Verifies build artifacts

#### Smart Contracts Checks

- Installs Foundry toolchain
- Compiles Solidity contracts
- Verifies contract artifacts

#### Summary Job

- Aggregates all check results
- Provides clear pass/fail status

### 2. Test Environment Configuration

Created test-specific environment files:

- `.env.test` - Backend test configuration
- `.env.client.test` - Frontend test configuration

### 3. Frontend Unit Tests (114 Tests Total ✅)

#### Hooks Tests (26 tests)

- **`use-local-storage.test.ts`** (11 tests)
    - Initialize with default/stored values
    - Update localStorage
    - Functional updates
    - Remove values
    - Custom serializer/deserializer
    - Handle invalid JSON
    - Sync across multiple hooks
    - InitializeWithValue option

- **`use-media-query.test.ts`** (7 tests)
    - Match media queries
    - Default values
    - Update on media query changes
    - Cleanup listeners
    - Support deprecated addListener API
    - Query change updates

- **`usePaginationState.test.ts`** (8 tests)
    - Initialize with default/custom values
    - Calculate correct offset
    - Handle page changes
    - Handle rows per page changes
    - Reset to page 1 on rows change
    - Stable function references
    - Complex pagination scenarios

#### Utility Functions Tests (26 tests)

- **`utils.test.ts`** (14 tests)
    - `cn()` utility for class merging
    - Conditional classes
    - Tailwind class conflicts
    - Array/object handling
    - Null/undefined values
    - `delay()` utility for async operations

- **`event.utils.test.ts`** (12 tests)
    - `toEventRegistrationConfigStatus()` converter
    - `toEventRegistrationConfigStatusNumber()` converter
    - Bidirectional conversion consistency
    - Edge cases and invalid inputs

#### Services Tests (18 tests)

- **`OnboardService.test.ts`** (9 tests)
    - Check onboard status via JWT cookie
    - Google OAuth method
    - Wallet method
    - Invalid parameter handling
    - Get sign message
    - Error handling

- **`AuthService.test.ts`** (9 tests)
    - Create account with Google OAuth
    - Create account with Wallet
    - Credential existence validation
    - Create profile (all fields & minimal)
    - Update profile
    - Sign out

#### Component Tests (44 tests)

- **`button.test.tsx`** (14 tests)
    - Render with children
    - Variant styles (primary, secondary-dark, secondary-light, ghost)
    - Size variants (sm, lg, xl, icon)
    - Click events
    - Disabled state
    - Loading state
    - Custom className
    - Props passthrough

- **`input.test.tsx`** (13 tests)
    - Text input
    - Different input types
    - Placeholder
    - onChange events
    - Disabled state
    - Custom className
    - Validation states
    - Controlled/uncontrolled

- **`typography.test.tsx`** (17 tests)
    - HTML tags (h1-h6, p, span, div)
    - Variants (header, text)
    - Sizes (small, base, subheader, header)
    - Colors (foreground, background, primary, secondary, muted)
    - Font families (inter, taviraj, anuphan, cormorant)
    - Font weights (normal, medium, semibold, bold)
    - Combined props
    - Custom className
    - Nested content

### 4. Test Infrastructure

#### Updated `apps/web/package.json`

Added test scripts:

```json
{
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage"
}
```

#### Enhanced `apps/web/src/test/setup.ts`

Added environment variable mocks:

- `VITE_CORE_BACKEND_API`
- `VITE_WALLETCONNECT_PROJECT_ID`
- `VITE_ENVIRONMENT`

## 📊 Test Results

```
✅ Test Files  10 passed (10)
✅ Tests       114 passed (114)
⏱️  Duration   1.96s
```

## 🚀 CI Pipeline Features

### Parallel Execution

Jobs run in parallel when possible for faster feedback:

- Database & SQLC checks
- Backend checks (depends on database)
- OpenAPI checks (depends on backend)
- Frontend checks (independent)
- Contracts checks (independent)

### Optimizations

- Caching for pnpm, Go modules
- PostgreSQL health checks
- Service containers for database
- Frozen lockfile for reproducibility

### Error Detection

- SQLC generation drift
- OpenAPI documentation drift
- TypeScript type errors
- ESLint violations
- Test failures
- Build failures

## 📝 Usage

### Running Tests Locally

```bash
# Run all frontend tests
pnpm --filter @decm/web test

# Run with UI
pnpm --filter @decm/web test:ui

# Run with coverage
pnpm --filter @decm/web test:coverage
```

### CI Pipeline

The workflow automatically runs on every PR targeting `main`. No manual intervention required.

## 🎯 Coverage

### Test Categories Covered

- ✅ Custom React hooks (localStorage, media queries, pagination)
- ✅ Utility functions (class merging, delays, converters)
- ✅ Services (authentication, onboarding)
- ✅ UI components (button, input, typography)

### Mocking Strategy

- Environment variables
- Window APIs (matchMedia, localStorage)
- External APIs (core API client)
- i18n translations

## 🔧 Technologies Used

- **Testing Framework**: Vitest
- **React Testing**: @testing-library/react
- **Test Environment**: happy-dom
- **Mocking**: vi (Vitest mocks)
- **CI/CD**: GitHub Actions

## 📖 Best Practices Implemented

1. **Comprehensive Coverage**: Tests cover happy paths, edge cases, and error scenarios
2. **Isolation**: Each test is independent with proper setup/teardown
3. **Clear Assertions**: Descriptive test names and clear expectations
4. **Mocking**: Proper mocking of external dependencies
5. **Fast Execution**: Tests run in under 2 seconds
6. **Type Safety**: Full TypeScript support in tests
7. **CI Integration**: Automated testing on every PR

## 🎉 Impact

- **Quality Assurance**: Catch bugs before they reach production
- **Confidence**: Refactor with confidence knowing tests will catch regressions
- **Documentation**: Tests serve as living documentation
- **Speed**: Fast feedback loop for developers
- **Standards**: Enforces code quality standards via CI
