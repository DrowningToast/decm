# Generated Cursor Rules Summary

Three new Cursor Rules have been created for the DECM platform to improve code quality, testing practices, and PR review processes.

## 1. CodeRabbit Configuration & PR Review Guide

**File**: `.cursor/rules/coderabbit.mdc`

### Purpose

Provides guidelines for CodeRabbit PR review automation and best practices for pull request submissions.

### Key Sections

- **Automated Review Coverage**: What CodeRabbit analyzes for Go backend and TypeScript/React frontend
- **Key Review Areas**: Code standards, security, database patterns, API documentation, testing requirements, linting, and type safety
- **Common Review Comments**: Approved patterns and typical findings
- **PR Description Best Practices**: Template for PR submissions
- **Pre-PR Checks**: Commands to run before creating a PR
- **Performance Considerations**: N+1 queries, re-renders, bundle size
- **Security Focus**: PII encryption, authentication, and data protection

### When to Use

- Creating pull requests
- Understanding CodeRabbit feedback
- Learning DECM code review standards
- Setting up CI/CD pipelines

---

## 2. ESLint Configuration & TypeScript Linting Guide

**File**: `.cursor/rules/eslint-configuration.mdc`
**Applies to**: `*.ts, *.tsx, *.js, *.jsx`

### Purpose

Comprehensive guide to ESLint rules, configurations, and best practices for the DECM frontend TypeScript/React codebase.

### Key Sections

- **ESLint Configuration**: File locations and extended configurations
- **Running ESLint**: Commands for linting, fixing, and checking
- **ESLint Rules by Category**: TypeScript rules, React Hooks rules, React Refresh rules, general JavaScript rules
- **Common ESLint Errors & Fixes**: 8 detailed error scenarios with solutions
- **Component Linting Standards**: Props interfaces, hook usage, return types
- **Pre-commit Linting**: Husky integration with lint-staged
- **IDE Integration**: VS Code setup for automatic linting and formatting
- **Best Practices**: What to do and what not to do
- **Troubleshooting**: Cache issues, module resolution, performance

### When to Use

- Writing or reviewing TypeScript/React code
- Fixing linting errors
- Understanding type safety requirements
- Setting up IDE for development
- Configuring pre-commit hooks

### Commands Reference

```bash
pnpm lint                          # Run linting + TypeScript checking
pnpm lint -- --fix                 # Auto-fix issues
pnpm format                        # Format with Prettier
```

---

## 3. Frontend Testing Setup - Vitest & React Testing Library

**File**: `.cursor/rules/testing-setup.mdc`
**Applies to**: `*.test.ts, *.test.tsx, *.spec.ts, *.spec.tsx, *_test.go, vitest.config.ts`

### Purpose

Complete guide to frontend testing using Vitest and React Testing Library, including setup, patterns, and best practices.

### Key Sections

- **Test Configuration**: Vitest setup with happy-dom environment
- **Test File Structure**: Naming conventions and project organization
- **Component Testing Patterns**: Basic components, user interactions, form submissions
- **API Mocking with MSW**: Setup, handlers, and testing with mocked API
- **Hook Testing**: Custom hooks and context testing
- **Async Testing**: Waiting for elements, using waitFor
- **Accessibility Testing**: Testing roles and query priorities
- **Snapshot Testing**: When and how to use (sparingly)
- **Best Practices**: Do's and don'ts for effective testing
- **Debugging Tests**: Debug output, screen queries, UI dashboard
- **Coverage Reporting**: Generating and analyzing coverage
- **CI/CD Integration**: GitHub Actions example
- **Troubleshooting**: Common issues and solutions

### When to Use

- Writing unit and integration tests for React components
- Mocking API calls in tests
- Testing hooks and context providers
- Setting up test environment
- Understanding test coverage requirements

### Commands Reference

```bash
pnpm test                          # Run all tests
pnpm test --watch                  # Watch mode
pnpm test --coverage               # With coverage report
pnpm test --ui                     # Interactive UI dashboard
pnpm test Button.test.tsx          # Specific file
pnpm test --grep "pattern"         # Match pattern
```

### Coverage Targets

- Frontend: Minimum 70% coverage
- Critical paths (auth, payments): 100% coverage
- Utilities: 90%+ coverage

---

## Rule Metadata

### CodeRabbit Rule

```
Description: "CodeRabbit PR review configuration and automation for DECM platform"
Applies: All files (global)
```

### ESLint Rule

```
Description: "ESLint configuration and TypeScript linting rules for DECM frontend"
Applies: *.ts, *.tsx, *.js, *.jsx (frontend files)
```

### Testing Setup Rule

```
Description: "Frontend testing setup with Vitest and React Testing Library"
Applies: *.test.ts, *.test.tsx, *.spec.ts, *.spec.tsx, *_test.go, vitest.config.ts
```

---

## Integration with Existing Rules

These new rules complement existing DECM Cursor Rules:

1. **CodeRabbit** relates to:
    - [code-standards.mdc](.cursor/rules/code-standards.mdc)
    - [testing-conventions.mdc](.cursor/rules/testing-conventions.mdc)
    - [project-conventions.mdc](.cursor/rules/project-conventions.mdc)

2. **ESLint** relates to:
    - [code-standards.mdc](.cursor/rules/code-standards.mdc)
    - [project-conventions.mdc](.cursor/rules/project-conventions.mdc)

3. **Testing Setup** relates to:
    - [testing-conventions.mdc](.cursor/rules/testing-conventions.mdc)
    - [code-standards.mdc](.cursor/rules/code-standards.mdc)

---

## How to Use These Rules

### In Cursor IDE

1. Rules are automatically available in `.cursor/rules/` directory
2. Cursor AI will reference relevant rules based on file context
3. Use `fetch_rules` command to explicitly request a specific rule

### Fetching Rules Manually

```
fetch_rules(rule_names: ["coderabbit", "eslint-configuration", "testing-setup"])
```

### Development Workflow

1. **Write Code** → ESLint rule provides linting guidance
2. **Write Tests** → Testing Setup rule provides testing patterns
3. **Create PR** → CodeRabbit rule provides PR best practices
4. **Submit PR** → CodeRabbit automatically reviews based on all three rules

---

## Key Takeaways

### For Developers

- ✅ Always run `pnpm lint` and `pnpm test` before committing
- ✅ Follow explicit return types and props interfaces
- ✅ Write tests for components and hooks
- ✅ Mock APIs with MSW in tests
- ✅ TypeScript validation is included in `pnpm lint`

### For PR Reviews

- ✅ CodeRabbit will verify code standards automatically
- ✅ Address CodeRabbit suggestions before merging
- ✅ Ensure tests pass and coverage is adequate
- ✅ Follow PR description template
- ✅ Run pre-PR checks locally

### For Maintainers

- ✅ Use these rules to onboard new developers
- ✅ Enforce quality standards through CodeRabbit
- ✅ Monitor test coverage trends
- ✅ Update rules as standards evolve

---

## File Locations

All rules are located in: `.cursor/rules/`

```
.cursor/rules/
├── coderabbit.mdc                 (NEW)
├── eslint-configuration.mdc       (NEW)
├── testing-setup.mdc              (NEW)
├── code-standards.mdc             (existing)
├── project-conventions.mdc        (existing)
├── testing-conventions.mdc        (existing)
└── ... (other DECM rules)
```

---

## Next Steps

1. **Review the new rules** in `.cursor/rules/` directory
2. **Share with team** for alignment on standards
3. **Run `pnpm lint`** to verify ESLint works
4. **Run `pnpm test`** to verify testing setup
5. **Create test PR** to validate CodeRabbit integration

---

Generated: October 30, 2025
DECM Platform - Decentralized Event Management
