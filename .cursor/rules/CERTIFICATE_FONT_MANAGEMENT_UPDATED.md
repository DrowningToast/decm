# Certificate Font Management - Updated Implementation with Font Family Table

## Overview

This document describes the updated implementation that uses a database table to manage font families instead of hardcoded values. The system now uses **font family IDs** as foreign keys, providing better data integrity and easier font management.

## Database Schema

### New Table: `event_certificate_font_families`

```sql
CREATE TABLE event_certificate_font_families (
    id SERIAL PRIMARY KEY,
    font_family_name VARCHAR(100) NOT NULL UNIQUE,
    css_font_name VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    available_font_weights TEXT NOT NULL, -- Comma-separated: "100,200,300,400,500,600,700,800,900"
    is_support_italic BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ -- For soft delete
);
```

**Default Font Families** (from migration):

- Inter (default) - weights: 100-900, italic: yes
- Noto Sans Thai - weights: 100-900, italic: yes
- Prompt - weights: 100-900, italic: yes
- TH Sarabun New (Sarabun) - weights: 100-800, italic: yes
- Kanit - weights: 100-900, italic: yes
- Arial - weights: 400,700, italic: yes
- Tahoma - weights: 400,700, italic: yes

### Updated: `event_certificate_configs` Table

Changed from VARCHAR font_family fields to INTEGER foreign keys:

**Old Fields** (removed):

- `event_name_font_family` VARCHAR
- `name_font_family` VARCHAR
- `academic_institution_font_family` VARCHAR
- `certificate_title_font_family` VARCHAR
- `certificate_subtitle_font_family` VARCHAR

**New Fields** (added):

- `event_name_font_family_id` INTEGER → FK to `event_certificate_font_families(id)`
- `name_font_family_id` INTEGER → FK to `event_certificate_font_families(id)`
- `academic_institution_font_family_id` INTEGER → FK to `event_certificate_font_families(id)`
- `certificate_title_font_family_id` INTEGER → FK to `event_certificate_font_families(id)`
- `certificate_subtitle_font_family_id` INTEGER → FK to `event_certificate_font_families(id)`

## Backend Implementation

### 1. SQL Queries

**File**: `packages/database/queries/event_certificate_font_families.sql`

```sql
-- Get all active font families (ordered by default first, then alphabetically)
-- name: GetAllEventCertificateFontFamilies :many
SELECT * FROM event_certificate_font_families
WHERE deleted_at IS NULL
ORDER BY
    CASE WHEN is_default THEN 0 ELSE 1 END,
    font_family_name;

-- Get font family by ID
-- name: GetEventCertificateFontFamilyByID :one
SELECT * FROM event_certificate_font_families
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- Get default font family
-- name: GetDefaultEventCertificateFontFamily :one
SELECT * FROM event_certificate_font_families
WHERE is_default = TRUE AND deleted_at IS NULL
LIMIT 1;
```

**Updated**: `packages/database/queries/event_certificate_configs.sql`

```sql
-- name: UpdateEventCertificateTextConfig :one
UPDATE event_certificate_configs
SET
    event_name_font_family_id = sqlc.narg('event_name_font_family_id'),
    event_name_font_weight = sqlc.narg('event_name_font_weight'),
    name_font_family_id = sqlc.narg('name_font_family_id'),
    name_font_weight = sqlc.narg('name_font_weight'),
    academic_institution_font_family_id = sqlc.narg('academic_institution_font_family_id'),
    academic_institution_font_weight = sqlc.narg('academic_institution_font_weight'),
    certificate_title_font_family_id = sqlc.narg('certificate_title_font_family_id'),
    certificate_title_font_weight = sqlc.narg('certificate_title_font_weight'),
    certificate_subtitle_font_family_id = sqlc.narg('certificate_subtitle_font_family_id'),
    certificate_subtitle_font_weight = sqlc.narg('certificate_subtitle_font_weight'),
    updated_at = NOW()
WHERE event_id = sqlc.arg('event_id')
RETURNING *;
```

### 2. Data Gateway

**File**: `apps/backend/core-api/internal/datagateway/event/event_certificate_font_family.go`

```go
type EventCertificateFontFamilyDataGateway interface {
    GetAllEventCertificateFontFamilies(ctx context.Context) ([]generated.EventCertificateFontFamily, error)
    GetEventCertificateFontFamilyByID(ctx context.Context, id int32) (*generated.EventCertificateFontFamily, error)
    GetDefaultEventCertificateFontFamily(ctx context.Context) (*generated.EventCertificateFontFamily, error)
}
```

### 3. Repository

**File**: `apps/backend/core-api/internal/repositories/postgres/event_certificate_font_family.go`

Implements the data gateway interface with simple pass-through to sqlc queries.

**Updated**: `event_certificate_config.go` mapper to use font family IDs:

```go
EventNameFontFamilyID:           pgInt4ToInt32Ptr(gen.EventNameFontFamilyID),
NameFontFamilyID:                pgInt4ToInt32Ptr(gen.NameFontFamilyID),
// ... etc
```

### 4. Usecase

**File**: `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_font_family.go`

```go
type EventCertificateFontFamilyResponse struct {
    ID                   int32    `json:"id"`
    FontFamilyName       string   `json:"font_family_name"`
    CssFontName          string   `json:"css_font_name"`
    IsDefault            bool     `json:"is_default"`
    AvailableFontWeights []int32  `json:"available_font_weights"` // Parsed from CSV string
    IsSupportItalic      bool     `json:"is_support_italic"`
}

func GetAllEventCertificateFontFamilies(ctx) ([]EventCertificateFontFamilyResponse, error)
```

**Features**:

- Parses comma-separated font weights into int32 array
- Returns ordered list (default font first)

**Updated**: `event_certificate_config.go`

- Added `EventCertificateFontFamilyDg` field to usecase struct
- Updated `UpdateEventCertificateTextConfigParams` to use font family IDs

### 5. Handler

**File**: `apps/backend/core-api/internal/handler/eventconfig/get_event_certificate_font_families.go`

**New Endpoint**: `GET /eventconfig/certificate-font-families`

**Response**:

```json
{
    "font_families": [
        {
            "id": 1,
            "font_family_name": "Inter",
            "css_font_name": "Inter",
            "is_default": true,
            "available_font_weights": [100, 200, 300, 400, 500, 600, 700, 800, 900],
            "is_support_italic": true
        },
        {
            "id": 3,
            "font_family_name": "Prompt",
            "css_font_name": "Prompt",
            "is_default": false,
            "available_font_weights": [100, 200, 300, 400, 500, 600, 700, 800, 900],
            "is_support_italic": true
        }
        // ... more fonts
    ]
}
```

**Updated**: `update_event_certificate_text_config.go`

- Request body uses `event_name_font_family_id` (int32) instead of `event_name_font_family` (string)
- Validation removed for font family names (will be enforced by foreign key constraint)
- TODO: Add validation to check if font family IDs exist in database

### 6. Routes

**File**: `apps/backend/core-api/internal/handler/eventconfig/routes.go`

```go
// Public route (no authentication required)
publicGroup := r.Group("/eventconfig")
publicGroup.Get("/certificate-font-families", h.GetEventCertificateFontFamilies)
```

**Note**: Font families endpoint is public so users can see available fonts before creating an event.

### 7. Dependency Injection

**File**: `apps/backend/core-api/cmd/main.go`

```go
// Updated EventConfigUsecase constructor to include font family data gateway
eventConfigUc := eventconfig_usecase.NewEventConfigUsecase(
    pgRepo, // AuthenticationCredentialDg
    pgRepo, // EventDg
    pgRepo, // EventCertificateDg
    pgRepo, // EventCertificateDataGateway
    pgRepo, // EventCertificateSignatureDataGateway
    pgRepo, // EventCertificateFontFamilyDg (NEW)
    pgRepo, // EventRegistrationDg
    pgRepo, // EventIssuerDg
    pgRepo, // EventContractDg
    pgRepo, // InboxMessageDg
    *s3Service,
    logger,
)

// Updated Handler constructor to include logger
eventConfigHandler := eventconfig_handler.NewHandler(
    eventConfigUc,
    eventUc,
    authService,
    authenticationGuardMiddleware,
    roleGuardMiddleware,
    logger, // NEW
)
```

## Frontend Implementation

### 1. API Hook

**File**: `apps/web/src/hooks/useCertificateFontFamilies.ts`

```typescript
export function useCertificateFontFamilies() {
    const { data, isLoading, error } = useQuery({
        queryKey: QUERY_KEY.event.certificate.fontFamilies,
        queryFn: async () => {
            const response = await coreApiClient.getEventCertificateFontFamilies();
            return response.data;
        },
    });

    return {
        fontFamilies: data?.font_families || [],
        isLoading,
        error,
    };
}
```

### 2. Query Keys

**File**: `apps/web/src/lib/queryKeys.ts`

```typescript
event: {
    certificate: {
        config: (eventId: string) => ["event", eventId, "certificate", "config"] as const,
        fontFamilies: ["event", "certificate", "font-families"] as const, // NEW
    },
}
```

### 3. Updated Component

**File**: `apps/web/src/components/CertificateFontSettings.tsx`

**Key Changes**:

1. **Uses Font Family IDs**:

```typescript
export interface CertificateFontConfig {
    event_name_font_family_id?: number; // Changed from string to number
    event_name_font_weight?: number;
    // ... similar for all fields
}
```

2. **Fetches Font Families from API**:

```typescript
const { fontFamilies, isLoading: isLoadingFonts } = useCertificateFontFamilies();
```

3. **Dynamic Font Weight Options**:

```typescript
// Get available weights for a given font family ID
const getAvailableWeights = (fontFamilyId?: number) => {
    if (!fontFamilyId) return [];
    const fontFamily = fontFamilies.find((f) => f.id === fontFamilyId);
    return fontFamily?.available_font_weights || [];
};
```

4. **Dropdown Selects**:

```typescript
<Select value={fontConfig.event_name_font_family_id?.toString()}>
    <SelectContent>
        {fontFamilies.map((font) => (
            <SelectItem key={font.id} value={font.id.toString()}>
                {font.font_family_name}
            </SelectItem>
        ))}
    </SelectContent>
</Select>
```

5. **Smart Weight Filtering**:
    - Font weight options dynamically change based on selected font family
    - Only shows weights supported by the selected font
    - Example: Arial only shows 400 and 700, while Inter shows 100-900

### 4. Updated Hook

**File**: `apps/web/src/components/pages/HostPages/EventPages/useUpdateCertificateTextConfig.ts`

Uses `coreApiClient` to call the updated endpoint with font family IDs.

## API Endpoints

### GET /eventconfig/certificate-font-families

**Public endpoint** - No authentication required

**Response**:

```json
{
    "font_families": [
        {
            "id": 1,
            "font_family_name": "Inter",
            "css_font_name": "Inter",
            "is_default": true,
            "available_font_weights": [100, 200, 300, 400, 500, 600, 700, 800, 900],
            "is_support_italic": true
        }
    ]
}
```

### PUT /events/{event_id}/certificates/text-config

**Updated Request**:

```json
{
    "event_name_font_family_id": 1,
    "event_name_font_weight": 700,
    "name_font_family_id": 3,
    "name_font_weight": 600,
    "academic_institution_font_family_id": 5,
    "academic_institution_font_weight": 500,
    "certificate_title_font_family_id": 1,
    "certificate_title_font_weight": 700,
    "certificate_subtitle_font_family_id": 4,
    "certificate_subtitle_font_weight": 400
}
```

## Benefits of Font Family Table Approach

### 1. **Data Integrity**

- Foreign key constraints ensure only valid font families are used
- No typos or invalid font names
- Database enforces referential integrity

### 2. **Centralized Font Management**

- Add/remove fonts from one table
- No code changes needed to add new fonts
- Update font properties (weights, italic support) in one place

### 3. **Better UX**

- Frontend shows only supported weights for each font
- Default font automatically selected for new configurations
- Font families ordered with default first

### 4. **Scalability**

- Easy to add custom fonts uploaded by users
- Can add font metadata (preview images, categories, etc.)
- Support for soft delete (retain history)

### 5. **Type Safety**

- Integer IDs prevent string comparison issues
- Generated TypeScript types ensure compile-time safety
- sqlc generates type-safe Go code

## Usage Flow

1. **User opens Certificate Settings page**
2. **Frontend fetches available fonts** via `GET /eventconfig/certificate-font-families`
3. **User selects font family** from dropdown (shows font names)
4. **Font weight dropdown updates** to show only supported weights for selected font
5. **User selects font weight** from filtered options
6. **User clicks "Save Font Settings"**
7. **Frontend sends font family IDs** (not names) to backend
8. **Backend validates** via foreign key constraints
9. **Database updates** configuration
10. **Certificate generation** fetches font family names by ID when rendering

## Certificate Generation Updates Required

The `generate_certificate_image.go` file currently has a TODO to fetch font family names from the database. This needs to be implemented:

```go
// TODO in addTextOverlaysToSVG function:
// Instead of using defaultFontFamily = "Inter" for everything,
// Fetch font family names from event_certificate_font_families table
// using the font_family_id values from config

// Example implementation needed:
func (u *EventUsecase) getFontFamilyName(ctx context.Context, fontFamilyID *int32) string {
    if fontFamilyID == nil {
        return "Inter" // default
    }
    fontFamily, err := u.EventCertificateFontFamilyDg.GetEventCertificateFontFamilyByID(ctx, *fontFamilyID)
    if err != nil {
        u.logger.Warn("failed to get font family, using default", "id", *fontFamilyID)
        return "Inter"
    }
    return fontFamily.CssFontName
}
```

## Files Created/Updated

### Backend

**Created**:

- ✅ `packages/database/migrations/000010_create_event_certificate_font_families.up.sql`
- ✅ `packages/database/migrations/000010_create_event_certificate_font_families.down.sql`
- ✅ `packages/database/queries/event_certificate_font_families.sql`
- ✅ `apps/backend/core-api/internal/datagateway/event/event_certificate_font_family.go`
- ✅ `apps/backend/core-api/internal/repositories/postgres/event_certificate_font_family.go`
- ✅ `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_font_family.go`
- ✅ `apps/backend/core-api/internal/handler/eventconfig/get_event_certificate_font_families.go`

**Updated**:

- ✅ `apps/backend/core-api/internal/entity/event_certificate_config.go` - Font family IDs
- ✅ `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go` - Mapper updated
- ✅ `apps/backend/core-api/internal/usecase/eventconfig/event_config.go` - Added font family DG
- ✅ `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go` - Font family IDs in response
- ✅ `apps/backend/core-api/internal/handler/event/update_event_certificate_text_config.go` - Font family IDs in request/response
- ✅ `apps/backend/core-api/internal/handler/eventconfig/handler.go` - Added Logger field
- ✅ `apps/backend/core-api/internal/handler/eventconfig/routes.go` - Added font families route
- ✅ `apps/backend/core-api/cmd/main.go` - Updated DI wiring

### Frontend

**Created**:

- ✅ `apps/web/src/hooks/useCertificateFontFamilies.ts`

**Updated**:

- ✅ `apps/web/src/lib/queryKeys.ts` - Added fontFamilies query key
- ✅ `apps/web/src/components/CertificateFontSettings.tsx` - Uses font family IDs and dynamic weights
- ✅ `apps/web/src/components/pages/HostPages/EventPages/useUpdateCertificateTextConfig.ts` - Uses coreApiClient

### Generated

- ✅ `packages/database/go/generated/*` - Updated models and queries
- ✅ `packages/api/src/apis/core/api.ts` - New endpoint and updated types
- ✅ `apps/backend/core-api/docs/*` - Updated OpenAPI specs

## Next Steps

### 1. Update Certificate Generation (REQUIRED)

The `generate_certificate_image.go` currently uses a hardcoded default font. Update it to:

- Accept `EventCertificateFontFamilyDataGateway` in EventUsecase
- Fetch font family names by ID when rendering certificates
- Use `css_font_name` field for SVG rendering

### 2. Font Family Management (Optional)

Add CRUD endpoints for admins to manage font families:

- POST `/admin/certificate-font-families` - Add custom font
- PUT `/admin/certificate-font-families/{id}` - Update font
- DELETE `/admin/certificate-font-families/{id}` - Soft delete font

### 3. Font Preview (Enhancement)

Add font preview in the settings UI:

- Show sample text in each font
- Real-time preview of selected font combination
- Preview certificate with selected fonts before saving

### 4. Custom Font Upload (Future)

Allow organizations to upload custom web fonts:

- Upload .woff2 font files
- Store in S3
- Generate CSS @font-face rules
- Add to font families table

## Testing

```bash
# 1. Run migrations
pnpm db:migrate

# 2. Start backend
pnpm dev:core

# 3. Test font families endpoint
curl http://localhost:8080/api/v1/eventconfig/certificate-font-families

# 4. Start frontend
pnpm dev

# 5. Navigate to certificate settings
http://localhost:3000/host/events/{event_id}/settings/certificate

# 6. Select different fonts and weights
# 7. Save and verify in database:
# SELECT * FROM event_certificate_configs WHERE event_id = 'your-event-id';
```

## Summary

The updated implementation provides a robust, scalable solution for font management:

✅ **Database-driven** font library  
✅ **Foreign key constraints** ensure data integrity  
✅ **Dynamic font weight filtering** based on font capabilities  
✅ **Public API endpoint** for font families  
✅ **Type-safe** throughout (Go + TypeScript)  
✅ **Extensible** - easy to add more fonts  
✅ **Future-proof** - supports custom font uploads

All code builds and lints successfully!
