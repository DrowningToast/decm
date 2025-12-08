# Frontend: Certificate Mint Readiness Implementation

## Overview

Added certificate mint readiness status display to the host's event certificate page, showing real-time validation of minting prerequisites.

## What Was Implemented

### 1. React Hook for API Integration

**File**: `apps/web/src/hooks/events/useCertificateMintReadiness.ts`

Created a custom React hook using TanStack Query:

- Fetches mint readiness status from the backend API
- Automatic caching with 30-second stale time
- Error handling and retry logic
- Conditionally enabled based on event ID availability

```typescript
export const useCertificateMintReadiness = (eventId: string | undefined) => {
    return useQuery<CoreApiInternalHandlerEventconfigCertificateMintReadinessResponse>({
        queryKey: ["certificate-mint-readiness", eventId],
        queryFn: async () => {
            if (!eventId) throw new Error("Event ID is required");
            return await coreApi.checkCertificateMintReadiness({ eventId });
        },
        enabled: !!eventId,
        staleTime: 30000,
        retry: 1,
    });
};
```

### 2. UI Component Integration

**File**: `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx`

Added mint readiness status card **before** the certificate published status section:

**Features**:

- ✅ **Visual Status Indicator**
    - Green background when ready to mint
    - Blue background when requirements pending
    - Check icon or empty circle for each requirement

- 📋 **Requirement Checklist**
    - Certificate configuration status
    - All issuers signed status (with X/Y count)
    - Certificate contract deployment status
    - Each item shows visual check/uncheck indicator

- ⚠️ **Missing Requirements List**
    - Displays specific missing requirements
    - Formatted as bullet points in a blue info box
    - Only shown when requirements are missing

**Layout**:

```
┌─────────────────────────────────────────────────┐
│ Certificate Mint Readiness Status (NEW!)        │
│ ✓ Certificate configuration set up              │
│ ✓ All issuers signed (3/3)                      │
│ ✓ Certificate contract deployed                 │
│ [Missing requirements if any]                   │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ Certificate Published Status (Existing)         │
│ ...                                              │
└─────────────────────────────────────────────────┘
```

### 3. Internationalization (i18n)

**File**: `apps/web/src/lib/i18n/locales/en.json`

Added translation keys under `events.hostDetails.certificates`:

| Key                   | Value                                       |
| --------------------- | ------------------------------------------- |
| `mintReady`           | "✅ Ready to Mint Certificates"             |
| `mintNotReady`        | "Certificate Minting Requirements"          |
| `configStatus`        | "Certificate configuration set up"          |
| `issuersSignedStatus` | "All issuers signed ({{signed}}/{{total}})" |
| `contractStatus`      | "Certificate contract deployed"             |
| `missingRequirements` | "Missing requirements"                      |

### 4. TypeScript API Client

**File**: `packages/api/src/apis/core/api.ts` (Generated)

The TypeScript API client was regenerated from the OpenAPI specification, providing:

- Type-safe method: `checkCertificateMintReadiness()`
- Complete type definitions for request/response
- Automatic request/response validation

## UI States

### State 1: Ready to Mint ✅

```
┌────────────────────────────────────────────────┐
│ ✅ ✅ Ready to Mint Certificates               │
│                                                 │
│ ✓ Certificate configuration set up             │
│ ✓ All issuers signed (3/3)                     │
│ ✓ Certificate contract deployed                │
└────────────────────────────────────────────────┘
```

### State 2: Requirements Pending ⏳

```
┌────────────────────────────────────────────────┐
│ ○ Certificate Minting Requirements             │
│                                                 │
│ ✓ Certificate configuration set up             │
│ ○ All issuers signed (1/3)                     │
│ ○ Certificate contract deployed                │
│                                                 │
│ Missing requirements:                          │
│ • Not all issuers have signed (1/3 signed)     │
│ • Certificate contract address is not set      │
└────────────────────────────────────────────────┘
```

### State 3: No Issuers Assigned ⚠️

```
┌────────────────────────────────────────────────┐
│ ○ Certificate Minting Requirements             │
│                                                 │
│ ✓ Certificate configuration set up             │
│ ○ All issuers signed (0/0)                     │
│ ✓ Certificate contract deployed                │
│                                                 │
│ Missing requirements:                          │
│ • No issuers have been assigned to this event  │
└────────────────────────────────────────────────┘
```

## Styling Details

### Color Scheme

- **Ready** (Green): `bg-green-50 border-green-200 text-green-900`
- **Not Ready** (Blue): `bg-blue-50 border-blue-200 text-blue-900`
- **Checkmarks**: `text-green-600` for completed items
- **Pending**: Gray border circles for incomplete items
- **Missing Info Box**: `bg-blue-100 text-blue-800`

### Responsive Design

- Full width container
- Flex layout with icon and content
- Proper spacing and padding
- Mobile-friendly text sizes

## Data Flow

```
┌──────────────┐
│  User Visits │
│ Certificate  │──┐
│    Page      │  │
└──────────────┘  │
                  │
                  ▼
┌─────────────────────────────────┐
│ useCertificateMintReadiness()   │
│ - Query Key: ["certificate-     │
│   mint-readiness", eventId]     │
│ - Stale Time: 30s               │
└─────────────┬───────────────────┘
              │
              ▼
┌────────────────────────────────┐
│ GET /api/v1/events/{eventId}/  │
│ config/certificate/mint-       │
│ readiness                      │
└────────────┬───────────────────┘
             │
             ▼
┌────────────────────────────────┐
│ Backend Validation:            │
│ 1. Check certificate config    │
│ 2. Validate ALL issuers signed │
│ 3. Verify contract deployed    │
└────────────┬───────────────────┘
             │
             ▼
┌────────────────────────────────┐
│ Response:                      │
│ {                              │
│   is_ready: boolean,           │
│   all_issuers_have_signed,     │
│   signed_issuers_count,        │
│   total_issuers_count,         │
│   has_certificate_config,      │
│   has_certificate_contract,    │
│   missing_requirements[]       │
│ }                              │
└────────────┬───────────────────┘
             │
             ▼
┌────────────────────────────────┐
│ React Component Renders:       │
│ - Status card with colors      │
│ - Checklist of requirements    │
│ - Missing requirements list    │
└────────────────────────────────┘
```

## Performance Considerations

1. **Caching**: 30-second stale time reduces API calls
2. **Conditional Fetching**: Only fetches when eventId is available
3. **Loading State**: Shows nothing while loading (no flicker)
4. **Error Handling**: Gracefully fails without breaking the page

## User Experience Benefits

1. **Visibility**: Users can see at a glance if certificates are ready to mint
2. **Clarity**: Specific checklist shows exactly what's needed
3. **Guidance**: Missing requirements list provides actionable feedback
4. **Consistency**: Matches existing UI patterns in the application
5. **Real-time**: Status updates automatically when page is refreshed

## Placement Rationale

The mint readiness status is placed **before** the published status because:

1. **Logical Flow**: Check readiness → Mint → Publish
2. **Priority**: Readiness must be checked before minting action
3. **Visual Hierarchy**: Most important status first
4. **User Journey**: Guides users through the correct sequence

## Testing Recommendations

### Unit Tests

- [ ] Hook returns correct data structure
- [ ] Hook handles loading state
- [ ] Hook handles error state
- [ ] Translation keys exist and render

### Integration Tests

- [ ] Component renders with ready state
- [ ] Component renders with not-ready state
- [ ] Missing requirements display correctly
- [ ] Issuer counts display correctly (X/Y)

### E2E Tests

- [ ] Navigate to certificate page
- [ ] Verify status card appears
- [ ] Check requirement checklist
- [ ] Verify missing requirements list

## Files Modified

1. **Created**:
    - `apps/web/src/hooks/events/useCertificateMintReadiness.ts`

2. **Modified**:
    - `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx`
    - `apps/web/src/lib/i18n/locales/en.json`

3. **Generated**:
    - `packages/api/src/apis/core/api.ts`

## Next Steps

Consider adding:

- [ ] Refresh button to manually check status
- [ ] Auto-refresh when issuers sign (WebSocket/polling)
- [ ] Estimated gas cost for minting
- [ ] "Mint Now" button that appears when ready
- [ ] Progress indicator showing % completion
- [ ] Tooltip explanations for each requirement
- [ ] Link to issuer management from issuer status

---

## Summary

✅ **Frontend integration complete** - Host event certificate page now displays real-time mint readiness status with detailed requirement checklist and actionable feedback.

The implementation seamlessly integrates with existing UI patterns and provides clear guidance to event organizers about certificate minting readiness.
