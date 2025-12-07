# Certificate Settings Update Impact Analysis

## Overview

This document describes the effects of updating certificate settings through the Certificate Settings Page (`/host/events/[eventId]/settings/certificate`).

## Update Flow

### 1. Frontend Update Process

**Location:** `apps/web/src/components/pages/HostPages/EventPages/CertificateSettingsPage.tsx`

When certificate settings are saved:

```183:256:apps/web/src/components/pages/HostPages/EventPages/CertificateSettingsPage.tsx
    // Handle form submission
    const handleSubmit = async () => {
        try {
            // Validate that at least one issuer is selected
            if (issuerManagement.selectedIssuers.length === 0) {
                toast.error(t("certificateSettings.noIssuersError"));
                return;
            }

            // TODO: Implement API call to save certificate settings
            console.log("Event ID:", eventId);
            console.log("Selected Issuers:", issuerManagement.selectedIssuers);
            console.log("SVG File:", certificateTemplate.svgFile);
            console.log("Detected Keywords:", allDetectedKeywords);

            const name = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ name }}",
            );

            const eventName = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ eventName }}",
            );

            const acedmicInstitutionName = allDetectedKeywords.find(
                (keyword) => keyword.keyword === "{{ academicInstitutionName }}",
            );

            if (certificateTemplate.svgFile && !name) {
                toast.error(t("certificateSettings.nameNotFound"));
                return;
            }

            const req: UpdateEventCertificateConfigPayload = {
                name_pos_x: name?.x ?? eventCertificateConfig?.name_pos_x ?? 0,
                name_pos_y: name?.y ?? eventCertificateConfig?.name_pos_y ?? 0,
                event_name_pos_x: eventName?.x ?? eventCertificateConfig?.event_name_pos_x ?? 0,
                event_name_pos_y: eventName?.y ?? eventCertificateConfig?.event_name_pos_y ?? 0,
                base_certificate_image: certificateTemplate.svgFile ?? undefined,
            };

            if (acedmicInstitutionName) {
                req.academic_institution_pos_x = acedmicInstitutionName.x;
                req.academic_institution_pos_y = acedmicInstitutionName.y;
            }

            await updateCertificateConfig(req);
            await updateEventIssuer(
                issuerManagement.selectedIssuers.map((issuer) => {
                    console.log(issuer);
                    return {
                        event_id: eventId,
                        issuer_credential_id: issuer.id,
                    };
                }),
            );

            // Refetch all queries to ensure the page is up to date
            await Promise.all([
                queryClient.refetchQueries({
                    queryKey: QUERY_KEY.event.certificate.config(eventId),
                }),
                queryClient.refetchQueries({
                    queryKey: QUERY_KEY.event.issuers.byEventId(eventId),
                }),
                queryClient.refetchQueries({
                    queryKey: QUERY_KEY.event.byId(eventId),
                }),
            ]);

            toast.success(t("certificateSettings.saveSuccess"));
        } catch (error) {
            console.error("Error saving certificate settings:", error);
            toast.error(t("certificateSettings.saveError"));
        }
    };
```

**Two separate API calls are made:**

1. `updateCertificateConfig` - Updates certificate template and position settings
2. `updateEventIssuer` - Updates issuer list

---

## 2. Certificate Config Update Impact

### Backend Handler

**Location:** `apps/backend/core-api/internal/handler/eventconfig/update_event_certificate_config.go`

```53:135:apps/backend/core-api/internal/handler/eventconfig/update_event_certificate_config.go
func (h *Handler) UpdateEventCertificateConfig(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	requestBody := UpdateEventCertificateConfigRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	if err := requestBody.IsValid(); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	params := eventconfig.UpdateEventCertificateConfigParams{}

	if requestBody.EventNamePosX != nil {
		params.EventNamePosX = requestBody.EventNamePosX
	}
	if requestBody.EventNamePosY != nil {
		params.EventNamePosY = requestBody.EventNamePosY
	}
	if requestBody.NamePosX != nil {
		params.NamePosX = requestBody.NamePosX
	}
	if requestBody.NamePosY != nil {
		params.NamePosY = requestBody.NamePosY
	}
	if requestBody.AcademicInstitutionPosX != nil {
		params.AcademicInstitutionPosX = requestBody.AcademicInstitutionPosX
	}
	if requestBody.AcademicInstitutionPosY != nil {
		params.AcademicInstitutionPosY = requestBody.AcademicInstitutionPosY
	}

	baseCertificateImage, _ := ctx.FormFile("base_certificate_image")
	if baseCertificateImage != nil {
		if err := validatorutils.ValidateImageFile(baseCertificateImage); err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		}
		params.BaseCertificateImage = baseCertificateImage
	}

	dbEventCertConfig, _ := h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if dbEventCertConfig == nil {
		_, err = h.EventConfigUc.CreateEventCertificateConfig(ctx.UserContext(), eventID, eventconfig.CreateEventCertificateConfigParams{
			BaseCertificateImage:    *params.BaseCertificateImage,
			EventNamePosX:           *params.EventNamePosX,
			EventNamePosY:           *params.EventNamePosY,
			NamePosX:                *params.NamePosX,
			NamePosY:                *params.NamePosX,
			AcademicInstitutionPosX: params.AcademicInstitutionPosX,
			AcademicInstitutionPosY: params.AcademicInstitutionPosY,
		})
		if err != nil {
			return customerror.Parse(&customerror.ErrInternalServer, err)
		}
	} else {
		_, err := h.EventConfigUc.UpdateEventCertificateConfig(ctx.UserContext(), eventID, params)
		if err != nil {
			return customerror.Parse(&customerror.ErrInternalServer, err)
		}
	}

	dbEventCertConfig, err = h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return ctx.Status(http.StatusOK).JSON(EventCertificateConfigResponse{
		ID:                        dbEventCertConfig.ID,
		EventID:                   dbEventCertConfig.EventID,
		BaseCertificateStorageKey: dbEventCertConfig.BaseCertificateStorageKey,
		EventNamePosX:             dbEventCertConfig.EventNamePosX,
		EventNamePosY:             dbEventCertConfig.EventNamePosY,
		NamePosX:                  dbEventCertConfig.NamePosX,
		NamePosY:                  dbEventCertConfig.NamePosY,
		AcademicInstitutionPosX:   dbEventCertConfig.AcademicInstitutionPosX,
		AcademicInstitutionPosY:   dbEventCertConfig.AcademicInstitutionPosY,
		CreatedAt:                 dbEventCertConfig.CreatedAt,
		UpdatedAt:                 dbEventCertConfig.UpdatedAt,
	})
}
```

### UseCase Layer

**Location:** `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go`

```84:143:apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go
func (uc *EventConfigUsecase) UpdateEventCertificateConfig(ctx context.Context, eventID uuid.UUID, params UpdateEventCertificateConfigParams) (*generated.EventCertificateConfig, error) {
	dbEventCertConfig, err := uc.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	updateParams := generated.UpdateEventCertificateConfigParams{
		EventID: eventID,
	}

	if params.EventNamePosX != nil {
		updateParams.EventNamePosX = *params.EventNamePosX
	} else {
		// Use existing value if not provided
		updateParams.EventNamePosX = dbEventCertConfig.EventNamePosX
	}

	if params.EventNamePosY != nil {
		updateParams.EventNamePosY = *params.EventNamePosY
	} else {
		// Use existing value if not provided
		updateParams.EventNamePosY = dbEventCertConfig.EventNamePosY
	}

	if params.NamePosX != nil {
		updateParams.NamePosX = *params.NamePosX
	}

	if params.NamePosY != nil {
		updateParams.NamePosY = *params.NamePosY
	}

	if params.AcademicInstitutionPosX != nil {
		updateParams.AcademicInstitutionPosX = pgtype.Float8{Float64: *params.AcademicInstitutionPosX, Valid: true}
	}

	if params.AcademicInstitutionPosY != nil {
		updateParams.AcademicInstitutionPosY = pgtype.Float8{Float64: *params.AcademicInstitutionPosY, Valid: true}
	}

	updateParams.BaseCertificateStorageKey = dbEventCertConfig.BaseCertificateStorageKey

	if params.BaseCertificateImage != nil {
		previousStorageKey := dbEventCertConfig.BaseCertificateStorageKey
		if previousStorageKey != "" {
			err := uc.S3Service.DeleteFile(ctx, previousStorageKey)
			if err != nil {
				return nil, err
			}
		}

		storageKey, err := uc.UploadBaseCertificateImage(ctx, eventID, params.BaseCertificateImage)
		if err != nil {
			return nil, err
		}
		updateParams.BaseCertificateStorageKey = storageKey
	}

	return uc.EventCertificateDg.UpdateEventCertificateConfig(ctx, updateParams)
}
```

### Database Update

**Location:** `packages/database/queries/event_certificate_configs.sql`

```25:37:packages/database/queries/event_certificate_configs.sql
-- name: UpdateEventCertificateConfig :one
UPDATE event_certificate_configs
SET
    base_certificate_storage_key = sqlc.arg('base_certificate_storage_key'),
    event_name_pos_x = COALESCE(sqlc.arg('event_name_pos_x'), event_name_pos_x),
    event_name_pos_y = COALESCE(sqlc.arg('event_name_pos_y'), event_name_pos_y),
    name_pos_x = sqlc.arg('name_pos_x'),
    name_pos_y = sqlc.arg('name_pos_y'),
    academic_institution_pos_x = sqlc.arg('academic_institution_pos_x'),
    academic_institution_pos_y = sqlc.arg('academic_institution_pos_y'),
    updated_at = NOW()
WHERE event_id = sqlc.arg('event_id')
RETURNING *;
```

---

## 3. Event Issuer Update Impact

### Backend UseCase

**Location:** `apps/backend/core-api/internal/usecase/event/update_event_issuer.go`

```21:55:apps/backend/core-api/internal/usecase/event/update_event_issuer.go
func (u *EventUsecase) UpdateEventIssuer(ctx context.Context, eventID uuid.UUID, params []UpdateEventIssuerParams, currentUser *auth.JwtClaims) ([]generated.EventIssuer, error) {
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	for _, param := range params {
		eventID := param.EventID
		issuerCredentialID := param.IssuerCredentialID

		_, err := u.EventIssuerDataGateway.GetEventIssuerByEventIDAndIssuerCredentialID(ctx, eventID, issuerCredentialID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err := u.EventIssuerDataGateway.CreateEventIssuer(ctx, generated.CreateEventIssuerParams{
					EventID:            eventID,
					IssuerCredentialID: issuerCredentialID,
					IsSigned:           0,
					Signature:          pgtype.Text{},
				})
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	return nil, nil
}
```

**⚠️ Important Behavior:** This function only **creates new issuers** if they don't exist. It does **NOT delete** issuers that are removed from the selection.

---

## Impact Summary

### ✅ What Gets Updated

1. **Database Record (`event_certificate_configs`)**
    - `base_certificate_storage_key` - Updated if new image uploaded
    - `event_name_pos_x`, `event_name_pos_x` - Position coordinates
    - `name_pos_x`, `name_pos_y` - Position coordinates
    - `academic_institution_pos_x`, `academic_institution_pos_y` - Position coordinates (nullable)
    - `updated_at` - Automatically set to current timestamp

2. **S3 Storage**
    - **Old certificate template image is DELETED** when a new one is uploaded
    - New template image is uploaded to S3 with a new storage key
    - Previous S3 file becomes inaccessible (deleted)

3. **Issuer List (`event_issuers`)**
    - New issuers are **created** if they don't exist
    - Existing issuers remain unchanged (no deletion)
    - Issuer signing status (`is_signed`) is **NOT reset** automatically

4. **Frontend Cache (React Query)**
    - Certificate config query is invalidated
    - Event issuers query is invalidated
    - All event queries are invalidated
    - Page data is refetched to show updated information

### ⚠️ What Does NOT Change

1. **Existing Certificates**
    - Already issued certificates are **NOT regenerated**
    - Certificate generation uses config at the time of creation
    - Changing template/positions only affects **future** certificate generation

2. **Removed Issuers**
    - Issuers removed from selection are **NOT deleted** from database
    - They remain in the `event_issuers` table
    - This is a potential data consistency issue

3. **Issuer Signing Status**
    - When certificate config changes, existing issuer signatures are **NOT invalidated**
    - Issuers who already signed certificates remain signed
    - New certificate template may require re-signing, but this is not enforced

### 🔄 Query Cache Invalidation

**Frontend hooks invalidate:**

```typescript
// useUpdateCertificateConfig
queryClient.invalidateQueries({
    queryKey: QUERY_KEY.event.certificate.config(eventId),
});
queryClient.invalidateQueries({
    queryKey: QUERY_KEY.event.all,
});

// useUpdateEventIssuer
queryClient.invalidateQueries({
    queryKey: QUERY_KEY.event.issuers.byEventId(eventID),
});
queryClient.invalidateQueries({
    queryKey: QUERY_KEY.event.all,
});

// Manual refetch in CertificateSettingsPage
queryClient.refetchQueries({
    queryKey: QUERY_KEY.event.certificate.config(eventId),
});
queryClient.refetchQueries({
    queryKey: QUERY_KEY.event.issuers.byEventId(eventId),
});
queryClient.refetchQueries({
    queryKey: QUERY_KEY.event.byId(eventId),
});
```

---

## Potential Issues & Considerations

### 1. **Missing Issuer Deletion**

- **Issue:** When issuers are removed from selection, they're not deleted from the database
- **Impact:** Orphaned issuer records, potential confusion
- **Recommendation:** Implement issuer deletion logic or use a sync approach (delete all, then recreate selected ones)

### 2. **No Signature Invalidation**

- **Issue:** When certificate template changes, existing signatures aren't invalidated
- **Impact:** Issuers may have signed certificates with the old template
- **Recommendation:** Consider resetting `is_signed = 0` for all issuers when config changes (similar to `ResetAllEventIssuersSigningStatus`)

### 3. **Old S3 File Deletion**

- **Issue:** Previous certificate template is deleted immediately
- **Impact:** If update fails partway, old template is lost
- **Recommendation:** Consider soft-delete or backup strategy for critical assets

### 4. **No Version History**

- **Issue:** Previous certificate configurations are not preserved
- **Impact:** Cannot rollback or audit configuration changes
- **Recommendation:** Consider adding versioning or audit logging

### 5. **Partial Update Failure**

- **Issue:** Two separate API calls (config + issuers) can fail independently
- **Impact:** System may be in inconsistent state if one succeeds and one fails
- **Recommendation:** Consider transaction/batch operation or better error handling

---

## Files Affected

### Frontend

- `apps/web/src/components/pages/HostPages/EventPages/CertificateSettingsPage.tsx`
- `apps/web/src/components/pages/HostPages/EventPages/useUpdateCertificateConfig.ts`
- `apps/web/src/components/pages/HostPages/EventPages/useUpdateEventIssuer.ts`

### Backend

- `apps/backend/core-api/internal/handler/eventconfig/update_event_certificate_config.go`
- `apps/backend/core-api/internal/usecase/eventconfig/event_certificate_config.go`
- `apps/backend/core-api/internal/usecase/event/update_event_issuer.go`
- `apps/backend/core-api/internal/repositories/postgres/event_certificate_config.go`
- `packages/database/queries/event_certificate_configs.sql`

---

## Testing Recommendations

1. **Test S3 file deletion** - Verify old template is deleted when new one uploaded
2. **Test partial updates** - Update only positions without uploading new image
3. **Test issuer sync** - Verify removed issuers are handled correctly
4. **Test signing status** - Verify issuer signatures when config changes
5. **Test cache invalidation** - Verify UI updates correctly after save
6. **Test error scenarios** - What happens if S3 upload fails? If issuer creation fails?
