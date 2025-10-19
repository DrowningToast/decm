package eventconfig

import (
	"apps/backend/core-api/internal/datagateway"
	eventDg "apps/backend/core-api/internal/datagateway/event"
	"log/slog"
)

type EventConfigUsecase struct {
	AuthenticationCredentialDg datagateway.AuthenticationCredentialDataGateway
	EventDg                    eventDg.EventDataGateway
	EventCertificateDg         eventDg.EventCertificateConfigDataGateway
	EventRegistrationDg        eventDg.EventRegistrationConfigDataGateway

	logger *slog.Logger
}

func NewEventConfigUsecase(
	authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway,
	eventDg eventDg.EventDataGateway,
	eventCertificateDg eventDg.EventCertificateConfigDataGateway,
	eventRegistrationDg eventDg.EventRegistrationConfigDataGateway,
	logger *slog.Logger,
) *EventConfigUsecase {
	return &EventConfigUsecase{
		AuthenticationCredentialDg: authenticationCredentialDg,
		EventDg:                    eventDg,
		EventCertificateDg:         eventCertificateDg,
		EventRegistrationDg:        eventRegistrationDg,
		logger:                     logger,
	}
}
