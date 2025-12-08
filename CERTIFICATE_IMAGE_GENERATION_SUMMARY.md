# Certificate Image Generation - Implementation Summary

## Overview

This document describes the implementation of the certificate image generation feature that converts SVG certificate templates into PNG images with personalized data.

## Implementation

### Location

`apps/backend/core-api/internal/usecase/event/generate_certificate_image.go`

### Main Function

```go
func (u *EventUsecase) GenerateCertificateImage(ctx context.Context, certificateID uuid.UUID) ([]byte, error)
```

### Process Flow

1. **Retrieve Certificate Data**
    - Fetches certificate by ID from database
    - Retrieves associated certificate configuration
    - Gets event details for event name

2. **Download SVG Template**
    - Downloads SVG file from S3 using `base_certificate_storage_key`
    - Reads SVG content into memory

3. **Replace Template Variables**
    - Processes SVG content with certificate data
    - Supports two replacement strategies:
        - **ID-based**: Replaces text within elements with specific IDs
        - **Placeholder-based**: Replaces inline placeholders like `{{name}}` or `{name}`

4. **Render to PNG**
    - Parses processed SVG using `oksvg`
    - Renders to high-quality PNG (2x scale) using `rasterx`
    - Encodes as PNG bytes

5. **Return PNG Bytes**
    - Returns byte array ready for HTTP response

## Template Variables

The following variables can be used in SVG templates:

| Variable               | Source                          | Required |
| ---------------------- | ------------------------------- | -------- |
| `name`                 | Certificate.Name                | Yes      |
| `event_name`           | Event.Title                     | Yes      |
| `academic_institution` | Certificate.AcademicInstitution | No       |
| `certificate_title`    | Certificate.CertificateTitle    | No       |
| `certificate_subtitle` | Certificate.CertificateSubtitle | No       |

## SVG Template Format

### Method 1: ID-Based Replacement

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600">
  <!-- Element IDs match variable names -->
  <text id="certificate_title" x="400" y="100">TITLE</text>
  <text id="name" x="400" y="300">NAME HERE</text>
  <text id="event_name" x="400" y="350">EVENT NAME</text>
  <text id="academic_institution" x="400" y="400">INSTITUTION</text>
  <text id="certificate_subtitle" x="400" y="500">SUBTITLE</text>
</svg>
```

**How it works:**

- Matches `<text id="name">...</text>` tags
- Replaces entire text content between opening and closing tags
- Preserves all attributes and styling

### Method 2: Placeholder-Based Replacement

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600">
  <!-- Use {{variable}} or {variable} syntax -->
  <text x="400" y="100">Certificate of Completion</text>
  <text x="400" y="300">This certifies that {{name}}</text>
  <text x="400" y="350">has completed {event_name}</text>
  <text x="400" y="400">at {{academic_institution}}</text>
</svg>
```

**How it works:**

- Searches for `{{variable}}` or `{variable}` patterns anywhere in SVG
- Replaces with actual values
- Works with inline text and mixed content

## Rendering Details

### Image Quality

- **Base dimensions**: Read from SVG viewBox
- **Default dimensions**: 800x600px if not specified
- **Rendering scale**: 2x (for high-quality output)
- **Output format**: PNG with alpha channel

### Dependencies

```go
import (
    "github.com/srwiley/oksvg"   // SVG parsing
    "github.com/srwiley/rasterx"  // Rasterization to PNG
)
```

These are already included in `go.mod`:

```
github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c
github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef
```

## Usage in Handler

### Example Handler Implementation

```go
// @Summary Generate certificate image
// @Description Generates a PNG image from certificate template with participant data
// @ID generate-certificate-image
// @Tags certificates
// @Produce image/png
// @Param certificate_id path string true "Certificate ID"
// @Success 200 {file} binary "PNG image"
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/certificates/{certificate_id}/image [get]
func (h *EventHandler) GenerateCertificateImage(ctx *fiber.Ctx) error {
    // Parse certificate ID
    certificateID, err := uuid.Parse(ctx.Params("certificate_id"))
    if err != nil {
        return customerror.New(customerror.StatusBadRequest, "Invalid certificate ID", err)
    }

    // Generate PNG
    pngBytes, err := h.EventUsecase.GenerateCertificateImage(ctx.Context(), certificateID)
    if err != nil {
        return err
    }

    // Set response headers
    ctx.Set("Content-Type", "image/png")
    ctx.Set("Content-Disposition", fmt.Sprintf("inline; filename=certificate-%s.png", certificateID))
    ctx.Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours

    // Return PNG bytes
    return ctx.Send(pngBytes)
}
```

### Authorization Considerations

For participant-facing endpoints, add authorization check:

```go
func (h *EventHandler) GenerateCertificateImageForParticipant(ctx *fiber.Ctx) error {
    // Get authenticated user ID from JWT
    authCredID := ctx.Locals("auth_credential_id").(uuid.UUID)

    certificateID, err := uuid.Parse(ctx.Params("certificate_id"))
    if err != nil {
        return customerror.New(customerror.StatusBadRequest, "Invalid certificate ID", err)
    }

    // Verify ownership
    certificate, err := h.EventUsecase.EventCertificateDataGateway.GetEventCertificateByID(
        ctx.Context(),
        certificateID,
    )
    if err != nil {
        return customerror.New(customerror.StatusNotFound, "Certificate not found", err)
    }

    if certificate.ReceiverCredentialId == nil ||
       *certificate.ReceiverCredentialId != authCredID {
        return customerror.New(
            customerror.StatusForbidden,
            "You do not have permission to access this certificate",
            nil,
        )
    }

    // Generate PNG
    pngBytes, err := h.EventUsecase.GenerateCertificateImage(ctx.Context(), certificateID)
    if err != nil {
        return err
    }

    ctx.Set("Content-Type", "image/png")
    ctx.Set("Content-Disposition", fmt.Sprintf("inline; filename=certificate-%s.png", certificateID))
    ctx.Set("Cache-Control", "public, max-age=86400")

    return ctx.Send(pngBytes)
}
```

## Error Handling

The function returns customer-facing errors for:

| Error                    | Status | Reason                            |
| ------------------------ | ------ | --------------------------------- |
| Certificate not found    | 404    | Invalid certificate ID            |
| Config not found         | 404    | No certificate config for event   |
| Template download failed | 500    | S3 service error                  |
| Template read failed     | 500    | IO error                          |
| Event not found          | 404    | Invalid event reference           |
| Rendering failed         | 500    | SVG parsing or PNG encoding error |

## Frontend Usage

### React Hook Example

```typescript
import { useState, useEffect } from 'react';

const useCertificateImage = ({ certificateId, enabled }) => {
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!enabled || !certificateId) return;

    setIsLoading(true);

    fetch(`${API_URL}/certificates/${certificateId}/image`, {
      credentials: 'include', // Include auth cookies
    })
      .then(res => {
        if (!res.ok) throw new Error('Failed to load certificate');
        return res.blob();
      })
      .then(blob => {
        const url = URL.createObjectURL(blob);
        setImageUrl(url);
        setIsLoading(false);
      })
      .catch(err => {
        setError(err);
        setIsLoading(false);
      });

    // Cleanup: revoke object URL when component unmounts
    return () => {
      if (imageUrl) URL.revokeObjectURL(imageUrl);
    };
  }, [certificateId, enabled]);

  return { imageUrl, isLoading, error };
};

// Usage in component
const CertificateViewer = ({ certificateId }) => {
  const { imageUrl, isLoading, error } = useCertificateImage({
    certificateId,
    enabled: true,
  });

  if (isLoading) return <Spinner />;
  if (error) return <div>Error loading certificate</div>;
  if (!imageUrl) return null;

  return (
    <img
      src={imageUrl}
      alt="Certificate"
      className="w-full h-auto"
      loading="lazy"
    />
  );
};
```

## Testing

### Unit Test Example

```go
func TestGenerateCertificateImage(t *testing.T) {
    // Setup mocks
    mockEventDg := &mockEventDataGateway{}
    mockCertDg := &mockEventCertificateDataGateway{}
    mockConfigDg := &mockEventCertificateConfigDataGateway{}
    mockS3 := &mockS3Service{}

    usecase := &EventUsecase{
        EventDataGateway:             mockEventDg,
        EventCertificateDataGateway:  mockCertDg,
        EventCertificateConfigDg:     mockConfigDg,
        S3Service:                    mockS3,
    }

    ctx := context.Background()
    certificateID := uuid.New()

    // Mock certificate data
    certificate := &entity.EventCertificate{
        Id:                  certificateID,
        EventId:             uuid.New(),
        Name:                stringPtr("John Doe"),
        AcademicInstitution: stringPtr("University"),
    }
    mockCertDg.EXPECT().GetEventCertificateByID(ctx, certificateID).Return(certificate, nil)

    // Mock config
    config := &generated.EventCertificateConfig{
        BaseCertificateStorageKey: "event/test/certificate/template.svg",
    }
    mockConfigDg.EXPECT().GetEventCertificateConfigByEventID(ctx, certificate.EventId).Return(config, nil)

    // Mock S3 download
    svgContent := `<svg><text id="name">PLACEHOLDER</text></svg>`
    mockS3.EXPECT().GetFile(ctx, config.BaseCertificateStorageKey).Return(
        io.NopCloser(strings.NewReader(svgContent)),
        nil,
    )

    // Mock event
    event := &entity.Event{Title: "Test Event"}
    mockEventDg.EXPECT().GetEventById(ctx, certificate.EventId).Return(event, nil)

    // Execute
    pngBytes, err := usecase.GenerateCertificateImage(ctx, certificateID)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, pngBytes)
    assert.Greater(t, len(pngBytes), 0)
}
```

## Performance Considerations

1. **Caching**: Consider caching generated PNGs in S3 or CDN
2. **Async Generation**: For bulk operations, use background jobs
3. **Image Size**: Current 2x scale (1600x1200) provides good quality
4. **Memory**: Each rendering uses ~10-20MB of memory temporarily

### Optimization: Cache Generated Images

```go
func (u *EventUsecase) GenerateCertificateImageCached(ctx context.Context, certificateID uuid.UUID) ([]byte, error) {
    // Check if already generated and cached
    cacheKey := fmt.Sprintf("certificate-images/%s.png", certificateID)

    cached, err := u.S3Service.GetFile(ctx, cacheKey)
    if err == nil {
        return io.ReadAll(cached)
    }

    // Generate new
    pngBytes, err := u.GenerateCertificateImage(ctx, certificateID)
    if err != nil {
        return nil, err
    }

    // Cache in S3 (fire and forget)
    go func() {
        _ = u.S3Service.PutFile(context.Background(), &s3.S3UploadRequestObject{
            storageKey:  cacheKey,
            file:        bytes.NewReader(pngBytes),
            contentType: "image/png",
        })
    }()

    return pngBytes, nil
}
```

## Troubleshooting

### SVG Not Rendering

- **Check SVG syntax**: Ensure SVG is well-formed XML
- **ViewBox required**: SVG must have `viewBox` attribute
- **Supported features**: oksvg has limited CSS support

### Text Not Replaced

- **Check IDs**: Element IDs must match variable names exactly
- **Case sensitive**: Variable names are case-sensitive
- **Placeholder format**: Use `{{name}}` or `{name}`, not `$name` or `%name%`

### Poor Image Quality

- **Increase scale**: Change `scale := 2.0` to `scale := 3.0` or higher
- **Vector fonts**: Ensure fonts are embedded or use system fonts
- **SVG dimensions**: Provide explicit width/height in SVG

---

## Summary

The certificate image generation feature provides a robust, scalable solution for generating personalized certificate images:

✅ **Secure**: Downloads templates from private S3 storage  
✅ **Flexible**: Supports multiple variable replacement strategies  
✅ **High Quality**: 2x rendering scale for crisp images  
✅ **Type-Safe**: Full Go type safety throughout  
✅ **Testable**: Clean architecture enables easy unit testing  
✅ **Performant**: Efficient SVG parsing and PNG encoding

**Next Steps:**

1. Create handler endpoint for the usecase
2. Add authorization guards for participant access
3. Generate OpenAPI documentation
4. Create TypeScript client with `pnpm gen-api:core`
5. Implement frontend certificate viewer component
6. Add caching layer for generated images (optional)
7. Monitor performance and optimize if needed
