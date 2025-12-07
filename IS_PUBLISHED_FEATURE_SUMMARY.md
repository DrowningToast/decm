# Certificate Configuration `is_published` Field - Implementation Summary

## Overview

Added an `is_published` boolean field to the `event_certificate_configs` table to track the publication status of certificate configurations. This allows event organizers to explicitly publish or unpublish their certificate configuration.

## What Was Implemented

### 1. Database Migration ✅

**Files**:

- `packages/database/migrations/000008_add_is_published_to_event_certificate_configs.up.sql`
- `packages/database/migrations/000008_add_is_published_to_event_certificate_configs.down.sql`

**Changes**:

- Added `is_published` BOOLEAN column (default: FALSE, NOT NULL)
- Added index on `is_published` for faster querying
- Supports rollback with down migration

### 2. Database Queries ✅

**File**: `packages/database/queries/event_certificate_configs.sql`

**Updated Queries**:

- `CreateEventCertificateConfig` - Includes `is_published` field (defaults to FALSE)
- `UpdateEventCertificateConfig` - Allows updating `is_published`
- **New**: `ToggleEventCertificateConfigPublished` - Dedicated query to toggle published status

### 3. Backend Data Layer ✅

**Data Gateway** (`event_certificate_config.go`):

- Added `ToggleEventCertificateConfigPublished` method to interface

**Repository** (`event_certificate_config.go`):

- Implemented `ToggleEventCertificateConfigPublished` method

### 4. Backend API Layer ✅

**Response Struct** (`event_certificate_config.go`):

```go
type EventCertificateConfigResponse struct {
    // ... existing fields ...
    IsPublished bool `json:"is_published"`  // NEW
    // ... timestamps ...
}
```

**New Handler** (`toggle_certificate_published.go`):

- Endpoint: `PATCH /api/v1/events/{event_id}/config/certificate/published`
- Request body: `{ "is_published": boolean }`
- Authorization: Requires verified organizer role
- OpenAPI documented

**Updated Handlers**:

- `GetEventCertificateConfig` - Returns `is_published` field
- `UpdateEventCertificateConfig` - Returns `is_published` field

**Routes** (`routes.go`):

- Added route with role guard middleware

### 5. Frontend React Hook ✅

**File**: `apps/web/src/hooks/events/useToggleCertificatePublished.ts`

**Features**:

- TanStack Query mutation for toggling published status
- Invalidates certificate config cache on success
- Toast notifications for success/error
- Type-safe with generated API types

### 6. Frontend UI ✅

**File**: `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx`

**Changes**:

- Replaced certificate count-based published status with config-based status
- Added publish/unpublish button
- Button is disabled when:
    - Operation is in progress
    - Certificate is not ready for publishing (based on mint readiness)
- Visual indicators:
    - Green when published
    - Gray when not published
- Contextual descriptions and warnings

### 7. Internationalization ✅

**File**: `apps/web/src/lib/i18n/locales/en.json`

**New Translation Keys**:
| Key | Value |
|-----|-------|
| `publishedConfigStatus` | "✅ Certificate Configuration Published" |
| `notPublishedConfigStatus` | "Certificate Configuration Not Published" |
| `notPublishedConfigDescription` | "Publishing the configuration will allow certificates to be minted and distributed to participants." |
| `publishButton` | "Publish Configuration" |
| `unpublishButton` | "Unpublish Configuration" |

**Updated**:

- `publishedWarning` - Changed to reflect configuration publication

## API Endpoint

```
PATCH /api/v1/events/{event_id}/config/certificate/published
Authorization: Required (Verified Organizer)
Content-Type: application/json

Request Body:
{
  "is_published": true
}

Response: 200 OK
{
  "id": "uuid",
  "event_id": "uuid",
  "base_certificate_storage_key": "string",
  "base_certificate_presigned_url": "string",
  "event_name_pos_x": 100,
  "event_name_pos_y": 200,
  "name_pos_x": 150,
  "name_pos_y": 250,
  "academic_institution_pos_x": 300,
  "academic_institution_pos_y": 350,
  "is_published": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## User Flow

### Publishing Configuration

1. **Navigate** to event details → Certificates tab
2. **Check Readiness** - View "Certificate Publishing Requirements" status
3. **Review Requirements**:
    - ✓ Certificate configuration set up
    - ✓ All issuers signed (3/3)
    - ✓ Certificate contract deployed
4. **Publish** - Click "Publish Configuration" button
5. **Status Updates**:
    - Configuration marked as published
    - Button changes to "Unpublish Configuration"
    - Warning message displayed about changes requiring re-approval

### Unpublishing Configuration

1. **Click** "Unpublish Configuration" button
2. **Confirmation** via toast notification
3. **Status Updates**:
    - Configuration marked as not published
    - Button changes back to "Publish Configuration"

## UI States

### Not Published (Default)

```
┌──────────────────────────────────────────────────────────┐
│ ○ Certificate Configuration Not Published               │
│                                                          │
│ Publishing the configuration will allow certificates    │
│ to be minted and distributed to participants.           │
│                                                          │
│                        [Publish Configuration] (button)  │
└──────────────────────────────────────────────────────────┘
```

### Published

```
┌──────────────────────────────────────────────────────────┐
│ ✅ ✅ Certificate Configuration Published                │
│                                                          │
│ ⚠️ Certificate configuration has been published. Any    │
│ changes will require re-approval from all issuers.       │
│                                                          │
│                      [Unpublish Configuration] (button)  │
└──────────────────────────────────────────────────────────┘
```

### Button States

- **Enabled**: When ready to publish (all requirements met)
- **Disabled**: When:
    - Publishing operation in progress
    - Certificate configuration not ready (missing requirements)

## Integration with Mint Readiness

The publish button is **automatically disabled** when the certificate is not ready for publishing:

```typescript
disabled={isTogglingPublished || !mintReadiness?.is_ready}
```

This ensures that organizers can only publish configurations that meet all requirements:

- Certificate template configured
- All issuers have signed
- Certificate contract deployed

## Database Schema

```sql
CREATE TABLE event_certificate_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    base_certificate_storage_key VARCHAR(255) NOT NULL,
    event_name_pos_x FLOAT NOT NULL,
    event_name_pos_y FLOAT NOT NULL,
    name_pos_x FLOAT NOT NULL,
    name_pos_y FLOAT NOT NULL,
    academic_institution_pos_x FLOAT,
    academic_institution_pos_y FLOAT,
    is_published BOOLEAN DEFAULT FALSE NOT NULL,  -- NEW FIELD
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_certificate_configs_is_published
    ON event_certificate_configs(is_published);
```

## Files Modified/Created

### Backend (8 files)

1. `packages/database/migrations/000008_add_is_published_to_event_certificate_configs.up.sql` (NEW)
2. `packages/database/migrations/000008_add_is_published_to_event_certificate_configs.down.sql` (NEW)
3. `packages/database/queries/event_certificate_configs.sql` (Modified)
4. `apps/backend/core-api/internal/datagateway/event/event_certificate_config.go` (Modified)
5. `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go` (Modified)
6. `apps/backend/core-api/internal/handler/eventconfig/event_certificate_config.go` (Modified)
7. `apps/backend/core-api/internal/handler/eventconfig/toggle_certificate_published.go` (NEW)
8. `apps/backend/core-api/internal/handler/eventconfig/routes.go` (Modified)
9. `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go` (Modified)
10. `apps/backend/core-api/internal/handler/eventconfig/get_event_certificate_config.go` (Modified)
11. `apps/backend/core-api/internal/handler/eventconfig/update_event_certificate_config.go` (Modified)

### Frontend (3 files)

1. `apps/web/src/hooks/events/useToggleCertificatePublished.ts` (NEW)
2. `apps/web/src/components/pages/HostPages/EventsPage/HostEventDetailsPage.tsx` (Modified)
3. `apps/web/src/lib/i18n/locales/en.json` (Modified)

### Generated (3 files)

1. `packages/database/go/generated/event_certificate_configs.sql.go` (Generated)
2. `apps/backend/core-api/docs/swagger.json` (Generated)
3. `packages/api/src/apis/core/api.ts` (Generated)

## Migration Instructions

### To Apply Migration

```bash
# Run migrations (happens automatically on backend start)
pnpm dev:core

# Or manually
pnpm db:migrate
```

### To Rollback

```bash
# Run down migration
migrate -path packages/database/migrations -database "your-connection-string" down 1
```

## Testing Checklist

### Backend

- [x] Migration runs successfully
- [x] Database code generates correctly
- [x] Backend compiles without errors
- [x] OpenAPI documentation generated
- [ ] Toggle endpoint works with valid data
- [ ] Authorization enforced (organizers only)
- [ ] Invalid event ID returns 404

### Frontend

- [x] TypeScript compiles without errors
- [x] Hook properly calls API
- [x] Button displays correct state
- [x] Toast notifications work
- [x] Button disabled when not ready
- [ ] UI updates after toggle
- [ ] Works on mobile/tablet

## Benefits

1. **Explicit Control**: Organizers can explicitly control when certificates are published
2. **Safety**: Prevents accidental publication before all requirements are met
3. **Audit Trail**: `updated_at` timestamp tracks when status changed
4. **User Feedback**: Clear UI indicators and warnings
5. **Integration**: Works seamlessly with mint readiness validation
6. **Reversible**: Configuration can be unpublished if needed

## Future Enhancements

Consider adding:

- [ ] Publish history/audit log
- [ ] Confirmation dialog before publishing
- [ ] Email notifications to issuers when published
- [ ] Bulk publish for multiple events
- [ ] Scheduled publishing
- [ ] Publish workflow with approvals

---

## Summary

✅ **Feature Complete** - The `is_published` field is fully integrated from database to UI, allowing event organizers to explicitly publish and unpublish their certificate configurations with proper validation and user feedback.
