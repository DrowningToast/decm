# Font Families Table Migration Summary

## Overview

Created a new `event_certificate_font_families` table to store font family metadata with auto-incrementing IDs, and refactored the `event_certificate_configs` table to use foreign key references instead of storing font names as strings.

## Database Changes

### New Table: `event_certificate_font_families`

**Schema:**

```sql
CREATE TABLE event_certificate_font_families (
    id SERIAL PRIMARY KEY,
    font_family_name VARCHAR(100) NOT NULL UNIQUE,
    css_font_name VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    available_font_weights TEXT NOT NULL,
    is_support_italic BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
```

**Seeded Fonts (7 fonts):**

| ID  | Font Family Name | CSS Name       | Default | Weights  | Italic |
| --- | ---------------- | -------------- | ------- | -------- | ------ |
| 1   | Inter            | Inter          | ✅      | 100-900  | ✅     |
| 2   | Noto Sans Thai   | Noto Sans Thai | ❌      | 100-900  | ✅     |
| 3   | Prompt           | Prompt         | ❌      | 100-900  | ✅     |
| 4   | TH Sarabun New   | Sarabun        | ❌      | 100-800  | ✅     |
| 5   | Kanit            | Kanit          | ❌      | 100-900  | ✅     |
| 6   | Arial            | Arial          | ❌      | 400, 700 | ✅     |
| 7   | Tahoma           | Tahoma         | ❌      | 400, 700 | ✅     |

### Modified Table: `event_certificate_configs`

**Before:**

- `event_name_font_family` VARCHAR(100)
- `name_font_family` VARCHAR(100)
- `academic_institution_font_family` VARCHAR(100)
- `certificate_title_font_family` VARCHAR(100)
- `certificate_subtitle_font_family` VARCHAR(100)

**After:**

- `event_name_font_family_id` INTEGER (FK → event_certificate_font_families)
- `name_font_family_id` INTEGER (FK → event_certificate_font_families)
- `academic_institution_font_family_id` INTEGER (FK → event_certificate_font_families)
- `certificate_title_font_family_id` INTEGER (FK → event_certificate_font_families)
- `certificate_subtitle_font_family_id` INTEGER (FK → event_certificate_font_families)

All foreign keys are NULLABLE with `ON DELETE SET NULL` constraint.

## Migration Files

**Migration 000010:**

- **Up:** `packages/database/migrations/000010_create_event_certificate_font_families.up.sql`
- **Down:** `packages/database/migrations/000010_create_event_certificate_font_families.down.sql`

## Backend Code Changes

### 1. Generated Code (SQLC)

**New Files:**

- `packages/database/go/generated/event_certificate_font_families.sql.go`
    - Query functions for font families CRUD operations

**Updated Models:**

- `EventCertificateFontFamily` struct in `models.go`
- `EventCertificateConfig` struct updated with new `*FamilyID` fields

### 2. SQL Queries

**New Query File:** `packages/database/queries/event_certificate_font_families.sql`

Available queries:

- `GetAllEventCertificateFontFamilies` - Get all fonts (default first)
- `GetEventCertificateFontFamilyByID` - Get by ID
- `GetEventCertificateFontFamilyByName` - Get by name
- `GetDefaultEventCertificateFontFamily` - Get default font
- `CreateEventCertificateFontFamily` - Create new font
- `UpdateEventCertificateFontFamily` - Update font
- `SoftDeleteEventCertificateFontFamily` - Soft delete
- `DeleteEventCertificateFontFamily` - Hard delete

**Updated Query File:** `packages/database/queries/event_certificate_configs.sql`

- Changed from `sqlc.arg('event_name_font_family')` to `sqlc.narg('event_name_font_family_id')`
- All font family fields now use `narg` (nullable) with Int4 type

### 3. Entity Layer

**File:** `apps/backend/core-api/internal/entity/event_certificate_config.go`

**Changes:**

```go
// Before
EventNameFontFamily *string `json:"event_name_font_family"`

// After
EventNameFontFamilyID *int32 `json:"event_name_font_family_id"`
```

### 4. Repository Layer

**File:** `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go`

**Changes:**

- Mapper function updated to use `EventNameFontFamilyID` instead of `EventNameFontFamily`
- Changed from `PgTextToStringPtr()` to `pgInt4ToInt32Ptr()` for font family fields
- Removed unused `pgmapper` import

### 5. UseCase Layer

**File:** `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go`

**Changes:**

- Updated `EventCertificateConfigResponse` struct to use `*int32` IDs instead of `*string` names
- Updated `UpdateEventCertificateTextConfigParams` struct field types
- Updated param assignments from `pgtype.Text` to `pgtype.Int4`
- Updated response mapper to use ID fields

### 6. Handler Layer

**File:** `apps/backend/core-api/internal/handler/event/update_event_certificate_text_config.go`

**Changes:**

- Updated `UpdateEventCertificateTextConfigRequest` struct to accept `*int32` IDs
- Removed font family name validation (now validates IDs - TODO: add DB lookup validation)
- Updated request-to-usecase params mapping

### 7. Certificate Image Generation

**File:** `apps/backend/core-api/internal/usecase/event/generate_certificate_image.go`

**Temporary Fix:**

- Using default font "Inter" for all certificate text elements
- Added TODO comment to fetch font names from database using IDs
- This is a temporary solution - needs enhancement to lookup font CSS names from IDs

## API Breaking Changes

### Request Format Change

**Before:**

```json
{
    "event_name_font_family": "Prompt",
    "event_name_font_weight": 700
}
```

**After:**

```json
{
    "event_name_font_family_id": 3,
    "event_name_font_weight": 700
}
```

### Response Format Change

**Before:**

```json
{
    "event_name_font_family": "Prompt",
    "event_name_font_weight": 700
}
```

**After:**

```json
{
    "event_name_font_family_id": 3,
    "event_name_font_weight": 700
}
```

## Frontend Migration Required

⚠️ **Breaking Change**: Frontend code must be updated to:

1. Fetch font families list using new endpoint (to be created):

    ```
    GET /api/v1/certificate/font-families
    ```

2. Send font family IDs instead of names:

    ```typescript
    // Before
    {
        event_name_font_family: "Prompt";
    }

    // After
    {
        event_name_font_family_id: 3;
    }
    ```

3. Display font family names by looking up IDs in fetched fonts list

## Testing

```bash
# 1. Database Migration
pnpm db:migrate

# 2. Verify Font Families
psql "postgres://postgres:decm_password@localhost:5432/decm?sslmode=disable" \
  -c "SELECT * FROM event_certificate_font_families ORDER BY id;"

# 3. Generate Code
pnpm db:generate

# 4. Build Backend
go build core-api/cmd/main.go
```

## Font Information Research

### Google Fonts (Web Fonts)

- **Inter**: Variable font, 100-900 weights, true italics
- **Noto Sans Thai**: 100-900 weights, italic support
- **Prompt**: 100-900 weights (9 weights), italic styles
- **Kanit**: 100-900 weights (9 weights), 18 styles total (with italics)
- **Sarabun (TH Sarabun New)**: 100-800 weights (8 weights), italic styles

### System Fonts

- **Arial**: 400, 700 weights, italic support
- **Tahoma**: 400, 700 weights, italic support (2010+ versions)

## Future Enhancements

### TODO Items

1. **Certificate Image Generation:**
    - Implement font family lookup from database using IDs
    - Replace hardcoded "Inter" with actual font CSS names from `event_certificate_font_families` table

2. **API Endpoints:**
    - Create `GET /api/v1/certificate/font-families` endpoint
    - Create handler for listing all available font families
    - Add Swagger documentation

3. **Validation:**
    - Add database lookup validation in handler to verify font family IDs exist
    - Validate font weight matches available weights for selected font

4. **Frontend Updates:**
    - Update certificate configuration forms to use font family dropdown (fetch from API)
    - Store and send font family IDs instead of names
    - Update certificate preview to display correct font names

5. **Database:**
    - Consider adding more fonts (if needed)
    - Add font family management endpoints (create, update, delete)

## Migration Status

✅ **Completed:**

- Database migration (000010) created and applied
- Font families table created with 7 default fonts
- Foreign key constraints added to event_certificate_configs
- SQLC code generated successfully
- Backend code updated (entity, repository, usecase, handler)
- Backend compilation successful

⚠️ **Pending:**

- Frontend code migration (breaking change)
- Certificate image generation enhancement (font lookup)
- Font families list API endpoint
- Font family ID validation

## References

- Migration: `packages/database/migrations/000010_create_event_certificate_font_families.up.sql`
- Queries: `packages/database/queries/event_certificate_font_families.sql`
- Entity: `apps/backend/core-api/internal/entity/event_certificate_config.go`
- Repository: `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go`
- UseCase: `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go`
- Handler: `apps/backend/core-api/internal/handler/event/update_event_certificate_text_config.go`

---

**Migration Date:** December 9, 2024
**Migration Number:** 000010
**Status:** ✅ Backend Complete, ⚠️ Frontend Pending
