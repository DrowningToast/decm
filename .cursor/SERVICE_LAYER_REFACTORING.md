# Service Layer Refactoring Summary

## Issue Fixed

The certificate preview was broken because:

1. Backend Swagger annotation was missing `/api/v1` prefix
2. API calls were not going through service layer with proper mappers
3. No camelCase/snake_case conversion between frontend and backend

## Changes Made

### 1. Backend Route Fix

**File**: `apps/backend/core-api/internal/handler/event/generate_certificate_image.go`

- Updated `@Router` annotation from `/certificates/{certificate_id}/image` to `/api/v1/certificates/{certificate_id}/image`
- Ensures generated API client has correct path structure

### 2. Certificate Service Layer (`apps/web/src/services/CertificateService/`)

#### CertificateService.ts

Added comprehensive certificate operations:

- `getCertificateImage(certificateId)` - Fetch certificate image as blob
- `getMyCertificatesList()` - Get user's certificates
- `getEventCertificates(eventId)` - Get event certificates
- `getFontFamilies()` - Get available fonts
- `signEventCertificates(eventId, issuerPin)` - Sign certificates
- `importCertificates(params)` - Import certificate receivers
- `revokeCertificates(params)` - Revoke specific certificates
- `revokeAllCertificates(eventId)` - Revoke all certificates
- `publishCertificates(eventId)` - Publish certificates
- `toggleCertificatePublished(eventId, isPublished)` - Toggle publish status
- `checkCertificateMintReadiness(eventId)` - Check mint readiness
- `getClaimCertificateSignMessage(certificateId)` - Get sign message for claiming
- `claimCertificateWithPin(params)` - Claim certificate with PIN/password
- `claimCertificateWithSignature(params)` - Claim certificate with wallet signature
- `claimCertificate(params)` - Unified claim method (auto-detects PIN or signature)

#### mapper.ts

Created proper snake_case → camelCase mappers:

- `Certificate` interface (frontend camelCase)
- `mapCertificate()` - Maps single certificate
- `mapToMyCertificatesViewModel()` - Maps certificate list
- `mapToSignedCertificatesResult()` - Maps signed certificates
- `mapImportCertificatesResponse()` - Maps import result
- `mapRevokeCertificatesResponse()` - Maps revoke result
- `mapRevokeAllCertificatesResponse()` - Maps revoke all result
- `mapPublishCertificatesResponse()` - Maps publish result
- `mapClaimCertificateSignMessage()` - Maps claim sign message
- `mapClaimCertificateResponse()` - Maps claim certificate result
- All mappers convert snake_case API responses to camelCase frontend interfaces

### 3. Updated Hooks to Use Service Layer

**All 13 certificate hooks** now use `certificateService` instead of direct `coreApiClient` calls:

#### Certificate Hooks

- `useMyCertificatesListViewModel` - Uses `certificateService.getMyCertificatesList()`
- `useEventCertificates` - Uses `certificateService.getEventCertificates()`
- `useCertificateImage` - Uses `certificateService.getCertificateImage()`
- `useCertificateFontFamilies` - Uses `certificateService.getFontFamilies()`
- `useCertificateMintReadiness` - Uses `certificateService.checkCertificateMintReadiness()`

#### Certificate Management Hooks

- `useSignEventCertificates` - Uses `certificateService.signEventCertificates()`
- `useImportCertificates` - Uses `certificateService.importCertificates()`
- `useRevokeEventCertificate` - Uses `certificateService.revokeCertificates()`
- `useRevokeAllEventCertificates` - Uses `certificateService.revokeAllCertificates()`
- `useToggleCertificatePublished` - Uses `certificateService.toggleCertificatePublished()`

#### Certificate Claiming Hooks

- `useGetClaimCertificateSignMessage` - Uses `certificateService.getClaimCertificateSignMessage()`
- `useClaimCertificate` - Uses `certificateService.claimCertificate()` (supports both PIN and signature methods)

### 4. Type Safety Improvements

- All API responses properly typed with camelCase interfaces
- Consistent data transformation through mapper layer
- No more mixing of snake_case and camelCase in frontend code

## Architecture Pattern

```
┌─────────────────┐
│  React Component│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Custom Hook   │ (useEventCertificates, useCertificateImage, etc.)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ CertificateService│ (Service layer with business logic)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Mapper       │ (snake_case ↔ camelCase conversion)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  API Client     │ (Generated from OpenAPI - coreApiClient.v1.*)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Backend API    │ (/api/v1/certificates/...)
└─────────────────┘
```

## Benefits

1. **Separation of Concerns**: Business logic in service layer, not in hooks
2. **Consistent Data Format**: All frontend code uses camelCase
3. **Type Safety**: Proper TypeScript interfaces at each layer
4. **Maintainability**: Single source of truth for API calls
5. **Reusability**: Service methods can be used across different components
6. **Testing**: Easier to mock service layer for testing

## Convention

### Frontend → Backend (Request)

```typescript
// Frontend (camelCase)
const params = {
    eventId: "123",
    hostPin: "1234",
};

// Service layer converts to snake_case
await coreApiClient.v1.importCertificates(
    { eventId: params.eventId },
    { host_pin: params.hostPin },
);
```

### Backend → Frontend (Response)

```typescript
// Backend response (snake_case)
{
  certificate_id: "123",
  event_id: "456",
  receiver_email: "user@example.com"
}

// Mapper converts to camelCase
{
  certificateId: "123",
  eventId: "456",
  receiverEmail: "user@example.com"
}
```

## Next Steps

Continue refactoring remaining hooks and components to use service layer pattern:

- Event hooks (useEvent, useEventsList, etc.)
- Profile hooks
- Auth hooks
- Registration hooks

All direct `coreApiClient` calls should go through appropriate service layer with proper mappers.
