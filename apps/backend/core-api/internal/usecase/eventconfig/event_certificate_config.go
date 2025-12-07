package eventconfig

import (
	"context"
	"decm-database/go/generated"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/s3"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// EventCertificateConfigResponse represents the response structure for event certificate config
type EventCertificateConfigResponse struct {
	ID                          uuid.UUID          `json:"id"`
	EventID                     uuid.UUID          `json:"event_id"`
	BaseCertificateStorageKey   string             `json:"base_certificate_storage_key"`
	BaseCertificatePresignedURL string             `json:"base_certificate_presigned_url"`
	EventNamePosX               float64            `json:"event_name_pos_x"`
	EventNamePosY               float64            `json:"event_name_pos_y"`
	NamePosX                    float64            `json:"name_pos_x"`
	NamePosY                    float64            `json:"name_pos_y"`
	AcademicInstitutionPosX     *float64           `json:"academic_institution_pos_x,omitempty"`
	AcademicInstitutionPosY     *float64           `json:"academic_institution_pos_y,omitempty"`
	IsPublished                 bool               `json:"is_published"`
	CreatedAt                   string             `json:"created_at"`
	UpdatedAt                   string             `json:"updated_at"`
	MintReadiness               *MintReadinessInfo `json:"mint_readiness,omitempty"`
}

// MintReadinessInfo contains mint readiness information embedded in certificate config
type MintReadinessInfo struct {
	IsReady                    bool     `json:"is_ready"`
	HasCertificateConfig       bool     `json:"has_certificate_config"`
	AllIssuersHaveSigned       bool     `json:"all_issuers_have_signed"`
	SignedIssuersCount         int64    `json:"signed_issuers_count"`
	TotalIssuersCount          int64    `json:"total_issuers_count"`
	HasCertificateContract     bool     `json:"has_certificate_contract"`
	CertificateContractAddress *string  `json:"certificate_contract_address,omitempty"`
	MissingRequirements        []string `json:"missing_requirements,omitempty"`
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

func (uc *EventConfigUsecase) GetEventCertificateConfigByEventID(ctx context.Context, eventID uuid.UUID) (*EventCertificateConfigResponse, error) {
	eventCertConfig, err := uc.EventCertificateDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event certificate config")
	}

	baseConfigPresignedURL, err := uc.S3Service.GetPresignedURL(ctx, eventCertConfig.BaseCertificateStorageKey)
	if err != nil {
		return nil, err
	}

	// Get mint readiness information
	mintReadiness, err := uc.CheckCertificateMintReadiness(ctx, eventID)
	if err != nil {
		// Log error but don't fail the entire request
		uc.logger.Error("failed to get mint readiness", "error", err, "event_id", eventID)
		mintReadiness = nil
	}

	var mintReadinessInfo *MintReadinessInfo
	if mintReadiness != nil {
		mintReadinessInfo = &MintReadinessInfo{
			IsReady:                    mintReadiness.IsReady,
			HasCertificateConfig:       mintReadiness.HasCertificateConfig,
			AllIssuersHaveSigned:       mintReadiness.AllIssuersHaveSigned,
			SignedIssuersCount:         mintReadiness.SignedIssuersCount,
			TotalIssuersCount:          mintReadiness.TotalIssuersCount,
			HasCertificateContract:     mintReadiness.HasCertificateContract,
			CertificateContractAddress: mintReadiness.CertificateContractAddress,
			MissingRequirements:        mintReadiness.MissingRequirements,
		}
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
		IsPublished:                 eventCertConfig.IsPublished,
		CreatedAt:                   eventCertConfig.CreatedAt.Time.String(),
		UpdatedAt:                   eventCertConfig.UpdatedAt.Time.String(),
		MintReadiness:               mintReadinessInfo,
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

// ToggleCertificatePublishedResult contains the result of toggling certificate published status
type ToggleCertificatePublishedResult struct {
	IsPublished          bool `json:"is_published"`
	InboxMessagesCreated int  `json:"inbox_messages_created,omitempty"`
}

// ToggleCertificatePublished toggles the certificate published status and creates inbox messages when publishing
func (uc *EventConfigUsecase) ToggleCertificatePublished(ctx context.Context, eventID uuid.UUID, isPublished bool, currentUserID uuid.UUID) (*ToggleCertificatePublishedResult, error) {
	response := &ToggleCertificatePublishedResult{
		IsPublished:          isPublished,
		InboxMessagesCreated: 0,
	}

	// If setting to published, validate and create inbox messages
	if isPublished {
		// 1. Get all certificates for the event
		certificates, err := uc.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get certificates")
		}

		if len(certificates) == 0 {
			return nil, errors.Wrap(fmt.Errorf("no certificates found for this event"), "cannot publish without certificates")
		}

		// 2. Get all event issuers
		eventIssuers, err := uc.EventIssuerDg.GetEventIssuersByEventID(ctx, eventID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get event issuers")
		}

		if len(eventIssuers) == 0 {
			return nil, errors.Wrap(fmt.Errorf("no issuers found for this event"), "cannot publish without issuers")
		}

		// 3. Validate: all issuers must have signed all certificates
		for _, certificate := range certificates {
			signatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateID(ctx, certificate.Id)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get certificate signatures")
			}

			// Check if we have signatures from all issuers
			if len(signatures) != len(eventIssuers) {
				return nil, fmt.Errorf("not all issuers have signature records for certificate %s", certificate.Id)
			}

			// Check if all issuers have actually signed (issuer_signature is not null)
			for _, signature := range signatures {
				if signature.IssuerSignature == nil || *signature.IssuerSignature == "" {
					return nil, fmt.Errorf("not all issuers have signed the certificates. Please ensure all issuers have completed signing before publishing")
				}
			}
		}

		// 4. Create inbox messages for all certificates
		inboxMessagesCreated, err := uc.createInboxMessagesForCertificates(ctx, eventID, currentUserID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create inbox messages")
		}
		response.InboxMessagesCreated = inboxMessagesCreated
	}

	// 5. Toggle the published status in database
	result, err := uc.EventCertificateDg.ToggleEventCertificateConfigPublished(
		ctx,
		generated.ToggleEventCertificateConfigPublishedParams{
			EventID:     eventID,
			IsPublished: isPublished,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to toggle certificate published status")
	}

	response.IsPublished = result.IsPublished
	return response, nil
}

// createInboxMessagesForCertificates creates inbox messages for all certificates of an event
func (uc *EventConfigUsecase) createInboxMessagesForCertificates(ctx context.Context, eventID uuid.UUID, currentUserID uuid.UUID) (int, error) {
	// Get event details for message content
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get event")
	}

	// Get all certificates for the event
	certificates, err := uc.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get certificates")
	}

	inboxMessagesCreated := 0
	for _, certificate := range certificates {
		// Skip if inbox message already exists for this certificate
		if certificate.InboxMessageId != nil {
			continue
		}

		// Skip if no receiver email (required for inbox message)
		if certificate.ReceiverEmail == nil || *certificate.ReceiverEmail == "" {
			uc.logger.Warn("skipping certificate without receiver email", "certificate_id", certificate.Id)
			continue
		}

		// Create message content with translations
		messageContent := map[string]string{
			"th": fmt.Sprintf("คุณได้รับเชิญให้รับใบประกาศนียบัตรสำหรับ: %s", event.Title),
			"en": fmt.Sprintf("You have been invited to claim your certificate for: %s", event.Title),
		}
		messageContentJSON, err := json.Marshal(messageContent)
		if err != nil {
			uc.logger.Error("failed to marshal message content", "error", err, "certificate_id", certificate.Id)
			continue
		}

		fallbackMessage := "You have been invited to claim your certificate"

		// Create inbox message
		inboxMessage, err := uc.InboxMessageDg.CreateInboxMessage(ctx, datagateway.CreateInboxMessageParameters{
			SenderCredentialID:     &currentUserID,
			ReceiverCredentialID:   certificate.ReceiverCredentialId,
			ReceiverEmail:          *certificate.ReceiverEmail,
			MessageType:            int(entity.InboxMessageTypeEventCertificateInvitation),
			MessageContent:         string(messageContentJSON),
			FallbackMessageContent: &fallbackMessage,
			IsRead:                 0,
		})
		if err != nil {
			uc.logger.Error("failed to create inbox message", "error", err, "certificate_id", certificate.Id)
			continue
		}

		// Update the certificate with the inbox_message_id
		_, err = uc.EventCertificateDataGateway.UpdateEventCertificateInboxMessageID(ctx, certificate.Id, inboxMessage.Id)
		if err != nil {
			uc.logger.Error("failed to update certificate with inbox message ID", "error", err, "certificate_id", certificate.Id)
			continue
		}

		inboxMessagesCreated++
	}

	return inboxMessagesCreated, nil
}
