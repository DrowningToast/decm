package eventconfig

import (
	"apps/backend/services/s3"
	"context"
	"decm-database/go/generated"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// EventCertificateConfigResponse represents the response structure for event certificate config
type EventCertificateConfigResponse struct {
	ID                          uuid.UUID `json:"id"`
	EventID                     uuid.UUID `json:"event_id"`
	BaseCertificateStorageKey   string    `json:"base_certificate_storage_key"`
	BaseCertificatePresignedURL string    `json:"base_certificate_presigned_url"`
	EventNamePosX               float64   `json:"event_name_pos_x"`
	EventNamePosY               float64   `json:"event_name_pos_y"`
	NamePosX                    float64   `json:"name_pos_x"`
	NamePosY                    float64   `json:"name_pos_y"`
	AcademicInstitutionPosX     *float64  `json:"academic_institution_pos_x"`
	AcademicInstitutionPosY     *float64  `json:"academic_institution_pos_y"`
	CreatedAt                   string    `json:"created_at"`
	UpdatedAt                   string    `json:"updated_at"`
}

type CreateEventCertificateConfigParams struct {
	BaseCertificateImage    multipart.FileHeader
	EventNamePosX           float64
	EventNamePosY           float64
	NamePosX                float64
	NamePosY                float64
	AcademicInstitutionPosX *float64
	AcademicInstitutionPosY *float64
}

func (uc *EventConfigUsecase) CreateEventCertificateConfig(ctx context.Context, eventID uuid.UUID, params CreateEventCertificateConfigParams) (*generated.EventCertificateConfig, error) {
	// Check if config already exists for this event
	existingConfig, err := uc.EventCertificateDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err == nil && existingConfig != nil {
		return nil, fmt.Errorf("event certificate config already exists for event ID: %s", eventID.String())
	}

	storageKey, err := uc.UploadBaseCertificateImage(ctx, eventID, &params.BaseCertificateImage)
	if err != nil {
		return nil, err
	}

	// Create new config
	createParams := generated.CreateEventCertificateConfigParams{
		EventID:                   eventID,
		BaseCertificateStorageKey: storageKey,
		EventNamePosX:             params.EventNamePosX,
		EventNamePosY:             params.EventNamePosY,
		NamePosX:                  params.NamePosX,
		NamePosY:                  params.NamePosY,
	}

	if params.AcademicInstitutionPosX != nil {
		createParams.AcademicInstitutionPosX = pgtype.Float8{Float64: *params.AcademicInstitutionPosX, Valid: true}
	}
	if params.AcademicInstitutionPosY != nil {
		createParams.AcademicInstitutionPosY = pgtype.Float8{Float64: *params.AcademicInstitutionPosY, Valid: true}
	}

	return uc.EventCertificateDg.CreateEventCertificateConfig(ctx, createParams)
}

type UpdateEventCertificateConfigParams struct {
	BaseCertificateImage    *multipart.FileHeader
	EventNamePosX           *float64
	EventNamePosY           *float64
	NamePosX                *float64
	NamePosY                *float64
	AcademicInstitutionPosX *float64
	AcademicInstitutionPosY *float64
}

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
	}

	if params.EventNamePosY != nil {
		updateParams.EventNamePosY = *params.EventNamePosY
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

func (uc *EventConfigUsecase) GetEventCertificateConfigByEventID(ctx context.Context, eventID uuid.UUID) (*EventCertificateConfigResponse, error) {
	eventCertConfig, err := uc.EventCertificateDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	baseConfigPresignedURL, err := uc.S3Service.GetPresignedURL(ctx, eventCertConfig.BaseCertificateStorageKey)
	if err != nil {
		return nil, err
	}

	return &EventCertificateConfigResponse{
		ID:                          eventCertConfig.ID,
		EventID:                     eventCertConfig.EventID,
		BaseCertificateStorageKey:   eventCertConfig.BaseCertificateStorageKey,
		BaseCertificatePresignedURL: baseConfigPresignedURL,
		EventNamePosX:               eventCertConfig.EventNamePosX,
		EventNamePosY:               eventCertConfig.EventNamePosY,
		NamePosX:                    eventCertConfig.NamePosX,
		NamePosY:                    eventCertConfig.NamePosY,
		AcademicInstitutionPosX:     &eventCertConfig.AcademicInstitutionPosX.Float64,
		AcademicInstitutionPosY:     &eventCertConfig.AcademicInstitutionPosY.Float64,
		CreatedAt:                   eventCertConfig.CreatedAt.Time.String(),
		UpdatedAt:                   eventCertConfig.UpdatedAt.Time.String(),
	}, nil
}

func (uc *EventConfigUsecase) DeleteEventCertificateConfig(ctx context.Context, eventID uuid.UUID) error {
	return uc.EventCertificateDg.DeleteEventCertificateConfig(ctx, eventID)
}

func (uc *EventConfigUsecase) UploadBaseCertificateImage(ctx context.Context, eventID uuid.UUID, file *multipart.FileHeader) (string, error) {
	requestObject, err := uc.S3Service.GetS3UploadRequestObject(s3.StorageKeyTypeEventCertificate, eventID, file)
	if err != nil {
		return "", err
	}

	storageKey, err := uc.S3Service.PutFile(ctx, requestObject)
	if err != nil {
		return "", err
	}

	return storageKey, nil
}
