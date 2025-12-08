# Certificate Image Generation

## Overview

The Certificate Image Generation feature allows the backend to dynamically generate certificate images (PNG) from SVG templates stored in S3. This provides a powerful way to create personalized certificates with participant-specific information rendered on-the-fly.

## Architecture

```
┌─────────────────┐
│   Client/Web    │
└────────┬────────┘
         │ HTTP Request
         ▼
┌─────────────────┐
│   Handler       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Use Case      │
│  - Fetch SVG    │
│  - Replace Vars │
│  - Render PNG   │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌──────┐  ┌──────┐
│  S3  │  │  DB  │
└──────┘  └──────┘
```

## Components

### 1. Use Case: `GenerateCertificateImage`

**Location**: `apps/backend/core-api/internal/usecase/event/generate_certificate_image.go`

**Purpose**:

- Fetches SVG templates from S3
- Replaces template variables by matching element `id` attributes
- Renders the modified SVG to PNG binary

**Key Methods**:

#### `GenerateCertificateImage(ctx, certificateConfigID, params)`

Generates a PNG from a certificate config ID.

```go
params := event.GenerateCertificateImageParams{
    TemplateVariables: event.CertificateTemplateVariables{
        Name:      "John Doe",
        EventName: "Tech Conference 2024",
    },
}

pngBytes, err := eventUsecase.GenerateCertificateImage(ctx, configID, params)
```

#### `GenerateCertificateImageByEventID(ctx, eventID, params)`

Convenience method that fetches the certificate config by event ID first.

```go
pngBytes, err := eventUsecase.GenerateCertificateImageByEventID(ctx, eventID, params)
```

### 2. Template Variables

The system uses **strictly-typed template variables** matching the `event_certificates` database schema to ensure type safety:

```go
type CertificateTemplateVariables struct {
    // Required fields
    Name      string   // Participant name (event_certificates.name) -> id="name"
    EventName string   // Event name (events.title) -> id="event_name"

    // Optional fields (from event_certificates table)
    AcademicInstitution *string  // event_certificates.academic_institution -> id="academic_institution"
    CertificateTitle    *string  // event_certificates.certificate_title -> id="certificate_title"
    CertificateSubtitle *string  // event_certificates.certificate_subtitle -> id="certificate_subtitle"
}
```

**Supported Fields** (matching database schema):

- ✅ `name` - Participant name (REQUIRED)
- ✅ `event_name` - Event name (REQUIRED)
- ✅ `academic_institution` - Academic institution (OPTIONAL)
- ✅ `certificate_title` - Certificate title (OPTIONAL)
- ✅ `certificate_subtitle` - Certificate subtitle (OPTIONAL)

#### SVG Template ID Syntax

The SVG template should include elements with `id` attributes matching the fields above. Multiple syntax patterns are supported:

- `id="name"` ✅ Recommended
- `id="{{ name }}"`
- `id="{{name}}"`
- `id="{{ .name }}"`
- `id="{{.name}}"`

#### Example SVG Template

```xml
<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
    <!-- Background -->
    <rect width="800" height="600" fill="#f8f9fa"/>

    <!-- Border -->
    <rect x="20" y="20" width="760" height="560"
          fill="none" stroke="#333" stroke-width="2"/>

    <!-- Certificate Title (optional) -->
    <text id="certificate_title" x="400" y="100"
          text-anchor="middle" font-size="36" font-weight="bold" fill="#333">
        TITLE_PLACEHOLDER
    </text>

    <!-- Participant Name (required) -->
    <text id="name" x="400" y="250"
          text-anchor="middle" font-size="32" fill="#000">
        NAME_PLACEHOLDER
    </text>

    <!-- Event Name (required) -->
    <text id="event_name" x="400" y="320"
          text-anchor="middle" font-size="24" fill="#666">
        EVENT_PLACEHOLDER
    </text>

    <!-- Academic Institution (optional) -->
    <text id="academic_institution" x="400" y="380"
          text-anchor="middle" font-size="18" fill="#999">
        INSTITUTION_PLACEHOLDER
    </text>

    <!-- Certificate Subtitle (optional) -->
    <text id="certificate_subtitle" x="400" y="440"
          text-anchor="middle" font-size="16" fill="#999">
        SUBTITLE_PLACEHOLDER
    </text>
</svg>
```

### 3. Database Schema

The SVG template is stored in S3, referenced by the `event_certificate_configs` table:

```sql
CREATE TABLE event_certificate_configs (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL,
    base_certificate_storage_key VARCHAR(255) NOT NULL, -- S3 key for SVG template
    -- ... other fields
);
```

**New Query Added**:

```sql
-- name: GetEventCertificateConfigByID :one
SELECT * FROM event_certificate_configs WHERE id = $1;
```

### 4. Dependencies

Added Go libraries for SVG rendering:

```go
github.com/srwiley/oksvg   // SVG parsing
github.com/srwiley/rasterx  // Rasterization to PNG
```

Installed via:

```bash
go get github.com/srwiley/oksvg github.com/srwiley/rasterx
```

## Usage Examples

### Example 1: Basic Certificate Generation

```go
// In your handler
func (h *Handler) GenerateCertificate(ctx *fiber.Ctx) error {
    eventID := uuid.MustParse(ctx.Params("eventId"))

    params := event.GenerateCertificateImageParams{
        TemplateVariables: event.CertificateTemplateVariables{
            Name:      "Jane Smith",
            EventName: "Annual Tech Summit",
        },
    }

    pngBytes, err := h.EventUsecase.GenerateCertificateImageByEventID(
        ctx.Context(),
        eventID,
        params,
    )
    if err != nil {
        return err
    }

    ctx.Set("Content-Type", "image/png")
    return ctx.Send(pngBytes)
}
```

### Example 2: Certificate for Existing Participant

```go
func (h *Handler) GetParticipantCertificate(ctx *fiber.Ctx) error {
    eventID := uuid.MustParse(ctx.Params("eventId"))
    participantID := uuid.MustParse(ctx.Params("participantId"))

    // Get participant certificate data
    cert, err := h.EventUsecase.EventCertificateDataGateway.GetEventCertificateByID(
        ctx.Context(),
        participantID,
    )
    if err != nil {
        return err
    }

    // Get event data
    event, err := h.EventUsecase.EventDataGateway.GetEventById(
        ctx.Context(),
        eventID,
    )
    if err != nil {
        return err
    }

    // Build strictly-typed template variables from certificate data
    templateVars := event.CertificateTemplateVariables{
        Name:                *cert.Name,
        EventName:           event.Title,
        AcademicInstitution: cert.AcademicInstitution,
        CertificateTitle:    cert.CertificateTitle,
        CertificateSubtitle: cert.CertificateSubtitle,
    }

    params := event.GenerateCertificateImageParams{
        TemplateVariables: templateVars,
    }

    pngBytes, err := h.EventUsecase.GenerateCertificateImageByEventID(
        ctx.Context(),
        eventID,
        params,
    )
    if err != nil {
        return err
    }

    ctx.Set("Content-Type", "image/png")
    ctx.Set("Content-Disposition", "inline; filename=certificate.png")
    ctx.Set("Cache-Control", "public, max-age=86400") // Cache 24h

    return ctx.Send(pngBytes)
}
```

### Example 3: Batch Certificate Generation

```go
func (h *Handler) GenerateBatchCertificates(ctx *fiber.Ctx) error {
    eventID := uuid.MustParse(ctx.Params("eventId"))

    // Get all certificates for event
    certificates, err := h.EventUsecase.EventCertificateDataGateway.GetEventCertificatesByEventID(
        ctx.Context(),
        eventID,
    )
    if err != nil {
        return err
    }

    // Generate images for all participants
    results := make([]struct {
        ParticipantID uuid.UUID
        ImageData     []byte
    }, 0, len(certificates))

    for _, cert := range certificates {
        templateVars := event.CertificateTemplateVariables{
            Name:                *cert.Name,
            EventName:           event.Title,
            AcademicInstitution: cert.AcademicInstitution,
            CertificateTitle:    cert.CertificateTitle,
            CertificateSubtitle: cert.CertificateSubtitle,
        }

        params := event.GenerateCertificateImageParams{
            TemplateVariables: templateVars,
        }

        pngBytes, err := h.EventUsecase.GenerateCertificateImageByEventID(
            ctx.Context(),
            eventID,
            params,
        )
        if err != nil {
            h.logger.Error("Failed to generate certificate", "error", err)
            continue
        }

        results = append(results, struct {
            ParticipantID uuid.UUID
            ImageData     []byte
        }{
            ParticipantID: cert.Id,
            ImageData:     pngBytes,
        })

        // Optionally save to S3 for caching
    }

    return ctx.JSON(fiber.Map{
        "count": len(results),
        "message": "Certificates generated successfully",
    })
}
```

## API Endpoints (Suggested)

### POST `/api/v1/events/{eventId}/certificate-image`

**Description**: Generate a certificate image with custom template variables.

**Request Body**:

```json
{
    "name": "John Doe",
    "event_name": "Tech Conference 2024",
    "academic_institution": "University of Technology",
    "certificate_title": "Certificate of Achievement",
    "certificate_subtitle": "For Outstanding Performance"
}
```

**Note**: Only `name` and `event_name` are required. All other fields are optional and match the `event_certificates` table schema.

**Response**: PNG image (binary)

**Headers**:

- `Content-Type: image/png`
- `Content-Disposition: inline; filename=certificate.png`

### GET `/api/v1/events/{eventId}/participants/{participantId}/certificate-image`

**Description**: Generate a certificate for a specific participant using their stored data.

**Response**: PNG image (binary)

## Performance Considerations

### Caching Strategy

Since rendering SVG to PNG can be resource-intensive, consider:

1. **Client-Side Caching**: Set appropriate `Cache-Control` headers

    ```go
    ctx.Set("Cache-Control", "public, max-age=86400") // 24 hours
    ```

2. **Server-Side Caching**: Cache generated PNGs in S3 or Redis

    ```go
    cacheKey := fmt.Sprintf("cert-image:%s:%s", eventID, participantID)
    // Check cache first, generate only if not cached
    ```

3. **CDN**: Serve cached images through CloudFront or similar CDN

### Optimization Tips

- **SVG Size**: Keep SVG templates reasonably sized (<500KB)
- **Font Embedding**: Consider embedding fonts in the SVG for consistency
- **Background Jobs**: For batch generation, use background workers
- **Lazy Loading**: Generate on-demand rather than pre-generating for all participants

## Testing

Test file: `generate_certificate_image_test.go`

```bash
# Run all certificate image generation tests
go test -v ./core-api/internal/usecase/event -run TestReplaceSVGTemplateVariables
go test -v ./core-api/internal/usecase/event -run TestRenderSVGToPNG

# Run integration tests
go test -v ./core-api/internal/usecase/event -run TestIntegrationReplacementAndRender

# Benchmark
go test -bench=BenchmarkRenderSVGToPNG ./core-api/internal/usecase/event
```

## Troubleshooting

### Common Issues

**1. SVG Parsing Errors**

- **Cause**: Invalid or malformed SVG
- **Solution**: Validate SVG using online validators before uploading

**2. Template Variables Not Replacing**

- **Cause**: `id` attribute doesn't match the variable name
- **Solution**: Ensure exact match (case-sensitive)

**3. Poor Image Quality**

- **Cause**: SVG viewBox dimensions too small
- **Solution**: Use higher resolution viewBox (e.g., 1600x1200 instead of 800x600)

**4. Performance Issues**

- **Cause**: Complex SVG or large file size
- **Solution**: Simplify SVG, remove unnecessary elements, implement caching

## Future Enhancements

1. **Dynamic Font Loading**: Support custom fonts from S3
2. **QR Code Integration**: Embed verification QR codes in certificates
3. **Watermarks**: Add dynamic watermarks for security
4. **Multi-Format Support**: Generate PDF in addition to PNG
5. **Background Images**: Composite with custom background images
6. **Template Versioning**: Support multiple template versions per event

## Related Files

- Use Case: `apps/backend/core-api/internal/usecase/event/generate_certificate_image.go`
- Tests: `apps/backend/core-api/internal/usecase/event/generate_certificate_image_test.go`
- Handler Example: `apps/backend/core-api/internal/handler/event/generate_certificate_image_example.go.example`
- Data Gateway: `apps/backend/core-api/internal/datagateway/event/event_certificate_config.go`
- Repository: `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go`
- SQL Queries: `packages/database/queries/event_certificate_configs.sql`

## Summary

The Certificate Image Generation feature provides a complete server-side solution for creating personalized certificate images. By leveraging SVG templates and Go's powerful image processing libraries, the backend can generate high-quality PNG certificates on-the-fly, eliminating the need for client-side rendering and ensuring consistent output across all platforms.
