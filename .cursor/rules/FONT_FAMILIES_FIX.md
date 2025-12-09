# Certificate Font Families API Fix

## Problem

The `getEventCertificateFontFamilies` endpoint was returning 404 errors because the OpenAPI annotation was missing the `/api/v1` prefix.

## Root Cause

In `apps/backend/core-api/internal/handler/eventconfig/get_event_certificate_font_families.go`, the `@Router` annotation was:

```go
// @Router /eventconfig/certificate-font-families [get]
```

This caused the OpenAPI spec generator to create the endpoint path without `/api/v1`, which meant:

- The generated TypeScript API client placed it under `certificateFontFamilies` namespace instead of `v1`
- The actual route was mounted at `/api/v1/eventconfig/certificate-font-families` (because all handlers are mounted under `apiV1` group)
- But the generated client was calling `/eventconfig/certificate-font-families` (missing `/api/v1`)
- Result: **404 Not Found**

## Solution

### 1. Fixed the OpenAPI Annotation

Updated the `@Router` annotation to include the full path:

```go
// @Router /api/v1/eventconfig/certificate-font-families [get]
```

### 2. Regenerated the API

```bash
pnpm docs:core      # Regenerate Swagger/OpenAPI docs
pnpm gen-api:core   # Regenerate TypeScript API client
```

### 3. Updated the Hook

Changed from:

```typescript
const response = await coreApiClient.certificateFontFamilies.getEventCertificateFontFamilies();
```

To:

```typescript
const response = await coreApiClient.v1.getEventCertificateFontFamilies();
```

## Files Modified

1. **apps/backend/core-api/internal/handler/eventconfig/get_event_certificate_font_families.go**
    - Fixed `@Router` annotation to include `/api/v1` prefix

2. **packages/api/src/apis/core/api.ts** (GENERATED)
    - Regenerated with correct endpoint path
    - Method now under `v1` namespace instead of `certificateFontFamilies`

3. **apps/web/src/hooks/events/useCertificateFontFamilies.ts**
    - Updated to use `coreApiClient.v1.getEventCertificateFontFamilies()`

4. **apps/web/src/components/CertificateFontSettings.tsx**
    - Added placeholders to all Select components for better UX

## Verification

The endpoint now works correctly:

- **Request**: `GET http://localhost:8080/api/v1/eventconfig/certificate-font-families`
- **Response**: 200 OK with font families data

## Key Takeaway

**Always include the full route path (including `/api/v1`) in OpenAPI `@Router` annotations**, even though the route is mounted under an `apiV1` group in the code. The annotation is for documentation generation, not runtime routing.
