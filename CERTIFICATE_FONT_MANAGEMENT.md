# Certificate Font Management Implementation

This document describes the implementation of font family and font weight management for certificate templates in the DECM platform.

## Overview

The system allows event hosts to customize the font family and weight for each text template in certificate designs, including:

- Event Name
- Participant Name
- Academic Institution
- Certificate Title
- Certificate Subtitle

## Backend Implementation

### 1. Database Schema

Font fields were added to the `event_certificate_configs` table via migration `000009`:

- `event_name_font_family` (VARCHAR)
- `event_name_font_weight` (INTEGER)
- `name_font_family` (VARCHAR)
- `name_font_weight` (INTEGER)
- `academic_institution_font_family` (VARCHAR)
- `academic_institution_font_weight` (INTEGER)
- `certificate_title_font_family` (VARCHAR)
- `certificate_title_font_weight` (INTEGER)
- `certificate_subtitle_font_family` (VARCHAR)
- `certificate_subtitle_font_weight` (INTEGER)

### 2. SQL Query

**File**: `packages/database/queries/event_certificate_configs.sql`

The query `UpdateEventCertificateTextConfig` updates all font configuration fields:

```sql
-- name: UpdateEventCertificateTextConfig :one
UPDATE event_certificate_configs
SET
    event_name_font_family = sqlc.arg('event_name_font_family'),
    event_name_font_weight = sqlc.arg('event_name_font_weight'),
    name_font_family = sqlc.arg('name_font_family'),
    name_font_weight = sqlc.arg('name_font_weight'),
    academic_institution_font_family = sqlc.arg('academic_institution_font_family'),
    academic_institution_font_weight = sqlc.arg('academic_institution_font_weight'),
    certificate_title_font_family = sqlc.arg('certificate_title_font_family'),
    certificate_title_font_weight = sqlc.arg('certificate_title_font_weight'),
    certificate_subtitle_font_family = sqlc.arg('certificate_subtitle_font_family'),
    certificate_subtitle_font_weight = sqlc.arg('certificate_subtitle_font_weight'),
    updated_at = NOW()
WHERE event_id = sqlc.arg('event_id')
RETURNING *;
```

### 3. Data Gateway

**File**: `apps/backend/core-api/internal/datagateway/event/event_certificate_config.go`

Interface method:

```go
UpdateEventCertificateTextConfig(ctx context.Context, params generated.UpdateEventCertificateTextConfigParams) (*generated.EventCertificateConfig, error)
```

### 4. Repository

**File**: `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go`

Implementation delegates to generated sqlc queries.

### 5. Usecase

**File**: `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go`

New method: `UpdateEventCertificateTextConfig`

- Validates certificate config exists
- Converts optional parameters to pgtype values
- Updates font configuration in database
- Returns updated configuration with all font fields

Updated `EventCertificateConfigResponse` to include font fields:

- All font family and weight fields as optional pointers
- Properly extracts values from pgtype fields

### 6. Handler

**File**: `apps/backend/core-api/internal/handler/event/update_event_certificate_text_config.go`

New endpoint: `PUT /events/{event_id}/certificates/text-config`

**Request Body**:

```json
{
    "event_name_font_family": "Prompt",
    "event_name_font_weight": 700,
    "name_font_family": "Sarabun",
    "name_font_weight": 600,
    "academic_institution_font_family": "Kanit",
    "academic_institution_font_weight": 500,
    "certificate_title_font_family": "Arial",
    "certificate_title_font_weight": 700,
    "certificate_subtitle_font_family": "Georgia",
    "certificate_subtitle_font_weight": 400
}
```

**Validation**:

- Font families: Prompt, Sarabun, Kanit, Arial, Helvetica, Times New Roman, Georgia, Verdana, Courier New
- Font weights: 100, 200, 300, 400, 500, 600, 700, 800, 900

**Response**: Full certificate configuration with updated font settings

### 7. Routes

**File**: `apps/backend/core-api/internal/handler/event/routes.go`

Added route:

```go
eventGroup.Put("/:event_id/certificates/text-config", h.UpdateEventCertificateTextConfig)
```

## Frontend Implementation

### 1. API Client

**Generated from OpenAPI specs**: `packages/api/src/apis/core/api.ts`

Method: `updateEventCertificateTextConfig(eventId, textConfig)`

### 2. Custom Hook

**File**: `apps/web/src/components/pages/HostPages/EventPages/useUpdateCertificateTextConfig.ts`

Hook: `useUpdateCertificateTextConfig(eventId)`

- Uses React Query mutation
- Invalidates certificate config query on success
- Returns `updateCertificateTextConfig` function and `isUpdatingCertificateTextConfig` status

### 3. Font Settings Component

**File**: `apps/web/src/components/CertificateFontSettings.tsx`

Features:

- Dropdown selects for font family (9 options)
- Dropdown selects for font weight (9 options: 100-900)
- Separate configuration for each text field:
    - Event Name
    - Participant Name
    - Academic Institution
    - Certificate Title
    - Certificate Subtitle
- Local state management with React hooks
- Save button to submit changes
- Loading state during updates

### 4. Page Integration

**File**: `apps/web/src/components/pages/HostPages/EventPages/CertificateSettingsPage.tsx`

Added "Step 3: Font Settings" section:

- Only shown when certificate config exists
- Integrates `CertificateFontSettings` component
- Handles font config updates with toast notifications
- Updates certificate config state after successful save

## Usage Flow

1. **Event Host** navigates to Certificate Settings page
2. **Upload Certificate Template** (Step 1-2) - creates initial config
3. **Configure Fonts** (Step 3):
    - Select font family from dropdown for each field
    - Select font weight from dropdown for each field
    - Click "Save Font Settings"
4. **Backend processes request**:
    - Validates font values
    - Updates database record
    - Returns updated configuration
5. **Frontend updates UI**:
    - Shows success toast
    - Refreshes certificate config data
    - Font settings persist for future certificate generation

## Certificate Generation

**File**: `apps/backend/core-api/internal/usecase/event/generate_certificate_image.go`

When generating certificate images, the system:

1. Retrieves certificate config including font settings
2. Applies font family and weight to each text element:
    ```go
    fontFamily := getValueOrDefault(pgmapper.PgTextToStringPtr(config.NameFontFamily), "Prompt")
    fontWeight := fontWeightToString(config.NameFontWeight, "bold")
    ```
3. Creates SVG text elements with configured fonts:
    ```go
    createTextElement(data.Name, config.NamePosX, config.NamePosY, fontFamily, fontWeight, 16)
    ```
4. Renders final certificate as PNG with customized fonts

## API Documentation

**OpenAPI Endpoint**: `PUT /api/v1/events/{event_id}/certificates/text-config`

**Tags**: Events, Certificates
**Authentication**: Required
**Authorization**: Event host

**Success Response** (200):

```json
{
    "id": "uuid",
    "event_id": "uuid",
    "event_name_font_family": "Prompt",
    "event_name_font_weight": 700,
    "name_font_family": "Sarabun",
    "name_font_weight": 600
    // ... other fields
}
```

**Error Responses**:

- 400: Invalid font family or weight
- 404: Certificate configuration not found
- 500: Internal server error

## Testing

To test the implementation:

1. Start backend: `pnpm dev:core`
2. Start frontend: `pnpm dev`
3. Navigate to: `http://localhost:3000/host/events/{event_id}/settings/certificate`
4. Upload a certificate template (if not done)
5. Scroll to "Step 3: Font Settings"
6. Select different fonts and weights
7. Click "Save Font Settings"
8. Verify success toast appears
9. Generate a test certificate to see font changes

## Future Enhancements

1. **Font Size Configuration**: Allow customizing font size per field
2. **Font Color**: Add color picker for text colors
3. **Font Preview**: Show live preview of font changes
4. **Custom Fonts**: Support uploading custom font files
5. **Font Templates**: Save and reuse font configuration presets
6. **Alignment Control**: Add text-anchor configuration (left, center, right)

## Files Modified/Created

### Backend

- ✅ `packages/database/queries/event_certificate_configs.sql` (query already existed)
- ✅ `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go` (updated)
- ✅ `apps/backend/core-api/internal/handler/event/update_event_certificate_text_config.go` (new)
- ✅ `apps/backend/core-api/internal/handler/event/routes.go` (updated)

### Frontend

- ✅ `apps/web/src/components/CertificateFontSettings.tsx` (new)
- ✅ `apps/web/src/components/pages/HostPages/EventPages/useUpdateCertificateTextConfig.ts` (new)
- ✅ `apps/web/src/components/pages/HostPages/EventPages/CertificateSettingsPage.tsx` (updated)

### Generated

- ✅ `packages/api/src/apis/core/api.ts` (regenerated)
- ✅ `apps/backend/core-api/docs/*` (regenerated)

## Conclusion

The font management system is fully implemented and integrated into the certificate upload page. Event hosts can now customize fonts for all text fields in their certificate templates, providing greater flexibility in certificate design while maintaining type safety and validation throughout the stack.
