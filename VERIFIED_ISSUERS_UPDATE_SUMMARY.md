# Verified Issuers Update Summary

## Overview

Updated the verified issuers endpoint to follow DECM conventions with proper ViewModel pattern, request validation, and frontend camelCase transformation.

## Backend Changes

### 1. Created ViewModel in Usecase Layer

**File:** `apps/backend/core-api/internal/usecase/issuer/verified_issuer_viewmodel.go`

- Created `VerifiedIssuerViewModel` struct with proper JSON tags
- Follows the convention of defining ViewModels in the usecase layer (not handler or repository)

### 2. Updated Usecase Implementation

**File:** `apps/backend/core-api/internal/usecase/issuer/get_verified_issuers.go`

- Added `GetVerifiedIssuersRequest` struct for request parameters
- Updated `GetVerifiedIssuers` method to:
    - Accept request struct instead of individual parameters
    - Return `[]VerifiedIssuerViewModel` instead of `[]entity.Profile`
    - Convert entity.Profile to ViewModel before returning
- Maintained existing search logic for encrypted PII fields

### 3. Updated Handler with Proper Validation

**File:** `apps/backend/core-api/internal/handler/issuer/get_verified_issuers.go`

- Implemented `GetVerifiedIssuersRequest` struct with query parameters
- Added `Parse()` method following DECM pattern:
    - Parses query parameters with default values (limit: 25, offset: 0)
    - Returns `*customerror.Err` on parsing errors
- Added `IsValid()` method with validation:
    - Validates limit is non-negative and max 1000
    - Validates offset is non-negative
    - Returns proper custom errors
- Updated handler to:
    - Follow Parse → IsValid → UseCase → Response pattern
    - Return ViewModels in JSON response
    - Use proper error handling

### 4. Updated OpenAPI Documentation

- Changed from `{object}` to `{array}` for response type
- Updated to reference `issuer.VerifiedIssuerViewModel`
- Added detailed parameter descriptions with defaults
- Changed error responses to use `customerror.ErrResponse`

### 5. Generated Database Queries and API Client

- Ran `pnpm db:generate` to regenerate sqlc code
- Ran `pnpm docs:core` to regenerate OpenAPI/Swagger docs
- Ran `pnpm gen-api:core` to generate TypeScript API client

## Frontend Changes

### 1. Created IssuerService

**File:** `apps/web/src/services/IssuerService.ts`

- Created `VerifiedIssuer` interface with camelCase properties
- Implements snake_case → camelCase transformation at service layer
- Exported `getVerifiedIssuers()` function that:
    - Calls the backend API
    - Transforms all snake_case keys to camelCase
    - Returns type-safe `VerifiedIssuer[]`

**Transformation Examples:**

```typescript
// Backend response (snake_case)
{
  authentication_credential_id: "...",
  first_name: "John",
  is_first_name_public: true
}

// Frontend data (camelCase)
{
  authenticationCredentialId: "...",
  firstName: "John",
  isFirstNamePublic: true
}
```

### 2. Updated Query Keys

**File:** `apps/web/src/lib/queryKeys.ts`

- Changed `issuers.verified` from constant to function
- Now accepts parameters: `verified(search?, limit?, offset?)`
- Returns proper query key array for React Query caching

### 3. Updated useVerifiedIssuers Hook

**File:** `apps/web/src/hooks/events/useVerifiedIssuers.ts`

- Imports `VerifiedIssuer` type and `getVerifiedIssuers` service
- Uses proper query key function with parameters
- Returns typed `VerifiedIssuer[]` data
- Removed direct API client usage (now uses service layer)

### 4. Updated useIssuerManagement Hook

**File:** `apps/web/src/hooks/useIssuerManagement.ts`

- Updated to use `VerifiedIssuer` instead of `EntityProfile`
- Updated all property access to camelCase:
    - `authentication_credential_id` → `authenticationCredentialId`
    - `first_name` → `firstName`
    - `last_name` → `lastName`
    - `academic_institution` → `academicInstitution`
- Maintains compatibility with existing Issuer interface

## Architecture Compliance

### ✅ Backend Conventions

1. **ViewModel Pattern**: ViewModels defined in usecase layer
2. **Handler Pattern**: Parse → Validate → UseCase → Response
3. **Error Handling**: Custom error types with TryParseAsCustomErr
4. **OpenAPI Documentation**: Complete annotations with proper types
5. **Data Layer Separation**: Repository returns entities, Usecase returns ViewModels

### ✅ Frontend Conventions

1. **Service Layer Transformation**: snake_case → camelCase at service boundary
2. **Type Safety**: Proper TypeScript interfaces for all data
3. **Query Keys**: Parameterized query keys for cache management
4. **Hook Composition**: Separation of data fetching (hook) and business logic (component)

## Testing Recommendations

### Backend

```bash
# Test the endpoint
curl "http://localhost:8080/api/v1/issuers?limit=10&offset=0&search=john"
```

### Frontend

1. Test search functionality in certificate settings page
2. Verify camelCase properties are accessible
3. Check that query caching works properly
4. Verify no console errors for undefined properties

## Files Modified

### Backend

- `apps/backend/core-api/internal/usecase/issuer/verified_issuer_viewmodel.go` (new)
- `apps/backend/core-api/internal/usecase/issuer/get_verified_issuers.go`
- `apps/backend/core-api/internal/handler/issuer/get_verified_issuers.go`
- `apps/backend/core-api/docs/*` (generated)
- `packages/database/go/generated/*` (generated)

### Frontend

- `apps/web/src/services/IssuerService.ts` (new)
- `apps/web/src/hooks/events/useVerifiedIssuers.ts`
- `apps/web/src/hooks/useIssuerManagement.ts`
- `apps/web/src/lib/queryKeys.ts`
- `packages/api/src/apis/core/api.ts` (generated)

## Migration Notes

### Breaking Changes

None - The API endpoint signature remains the same, only internal structure improved.

### Backward Compatibility

- Frontend components using `EntityProfile` directly should be updated to use `VerifiedIssuer`
- Any direct API calls should be replaced with service layer calls

## Next Steps

1. **Testing**: Test the search functionality in the certificate settings page
2. **Monitoring**: Monitor for any type errors in production
3. **Documentation**: Update API documentation if needed
4. **Similar Endpoints**: Consider applying same pattern to other issuer-related endpoints

## Summary

This update brings the verified issuers endpoint into full compliance with DECM coding standards:

- ✅ ViewModels in usecase layer
- ✅ Proper request validation in handlers
- ✅ Service layer transformation for frontend
- ✅ Proper query key management
- ✅ Type-safe end-to-end implementation
- ✅ No linting errors

The implementation follows the principle: **Backend (snake_case) → Service (transform) → Frontend (camelCase)**.
