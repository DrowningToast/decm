package eventconfig

import (
	"log/slog"

	"apps/backend/core-api/internal/datagateway"
	eventDg "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/s3"
)

type EventConfigUsecase struct {
	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
	EventDg                    eventDg.EventDataGateway
	EventCertificateDg         eventDg.EventCertificateConfigDataGateway
	EventRegistrationDg        eventDg.EventRegistrationConfigDataGateway
	S3Service                  s3.S3Service

	logger *slog.Logger
}

func NewEventConfigUsecase(
	authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway,
	eventDg eventDg.EventDataGateway,
	eventCertificateDg eventDg.EventCertificateConfigDataGateway,
	eventRegistrationDg eventDg.EventRegistrationConfigDataGateway,
	s3Service s3.S3Service,
	logger *slog.Logger,
) *EventConfigUsecase {
	return &EventConfigUsecase{
		AuthenticationCredentialDg: authenticationCredentialDg,
		EventDg:                    eventDg,
		EventCertificateDg:         eventCertificateDg,
		EventRegistrationDg:        eventRegistrationDg,
		S3Service:                  s3Service,
		logger:                     logger,
	}
}
