package event

import (
	"context"
	"log/slog"
	"mime/multipart"

	"apps/backend/core-api/config"
	datagateway "apps/backend/core-api/internal/datagateway"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/auth"
	"apps/backend/services/s3"

	"github.com/google/uuid"
)

type EventUsecase struct {
	EventDataGateway                     eventdatagateway.EventDataGateway
	EventContractDataGateway             eventdatagateway.EventContractDataGateway
	EventIssuerDataGateway               eventdatagateway.EventIssuerDataGateway
	EventCertificateDataGateway          eventdatagateway.EventCertificateDataGateway
	EventCertificateSignatureDataGateway eventdatagateway.EventCertificateSignatureDataGateway
	EventCertificateConfigDg             eventdatagateway.EventCertificateConfigDataGateway
	EventCertificateFontFamilyDg         eventdatagateway.EventCertificateFontFamilyDataGateway
	AuthenticationCredentialDg           datagateway.AuthenticationCredentialDataGateway
	EventRegistrationInvitationDg        datagateway.EventRegistrationInvitationDataGateway
	EventAttendeeDg                      datagateway.EventAttendeeDataGateway
	InboxMessageDg                       datagateway.InboxMessageDataGateway
	S3Service                            *s3.S3Service
	logger                               *slog.Logger
	authService                          *auth.AuthService
	cfg                                  *config.Config
	UploadEventBanner                    func(ctx context.Context, entityID uuid.UUID, bannerFile *multipart.FileHeader) (string, error)
	UploadEventIcon                      func(ctx context.Context, entityID uuid.UUID, iconFile *multipart.FileHeader) (string, error)
}

func NewEventUsecase(eventDataGateway eventdatagateway.EventDataGateway, eventContractDataGateway eventdatagateway.EventContractDataGateway, eventIssuerDataGateway eventdatagateway.EventIssuerDataGateway, eventCertificateDataGateway eventdatagateway.EventCertificateDataGateway, eventCertificateSignatureDataGateway eventdatagateway.EventCertificateSignatureDataGateway, eventCertificateConfigDg eventdatagateway.EventCertificateConfigDataGateway, eventCertificateFontFamilyDg eventdatagateway.EventCertificateFontFamilyDataGateway, authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway, eventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway, eventAttendeeDg datagateway.EventAttendeeDataGateway, inboxMessageDg datagateway.InboxMessageDataGateway, s3Service *s3.S3Service, logger *slog.Logger, authService *auth.AuthService, cfg *config.Config) *EventUsecase {
	uc := &EventUsecase{
		EventDataGateway:                     eventDataGateway,
		EventContractDataGateway:             eventContractDataGateway,
		EventIssuerDataGateway:               eventIssuerDataGateway,
		EventAttendeeDg:                      eventAttendeeDg,
		EventCertificateDataGateway:          eventCertificateDataGateway,
		EventCertificateSignatureDataGateway: eventCertificateSignatureDataGateway,
		EventCertificateConfigDg:             eventCertificateConfigDg,
		EventCertificateFontFamilyDg:         eventCertificateFontFamilyDg,
		AuthenticationCredentialDg:           authenticationCredentialDg,
		EventRegistrationInvitationDg:        eventRegistrationInvitationDg,
		InboxMessageDg:                       inboxMessageDg,
		S3Service:                            s3Service,
		logger:                               logger,
		authService:                          authService,
		cfg:                                  cfg,
	}
	// Initialize upload function fields with their default implementations
	uc.UploadEventBanner = uc.uploadEventBannerImpl
	uc.UploadEventIcon = uc.uploadEventIconImpl
	return uc
}
