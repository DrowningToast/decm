package eventconfig

import (
	eventDg "apps/backend/core-api/internal/datagateway/event"
)

type EventConfigUsecase struct {
	EventCertificateDg  eventDg.EventCertificateConfigDataGateway
	EventRegistrationDg eventDg.EventRegistrationConfigDataGateway
}

func NewEventConfigUsecase(
	eventCertificateDg eventDg.EventCertificateConfigDataGateway,
	eventRegistrationDg eventDg.EventRegistrationConfigDataGateway,
) *EventConfigUsecase {
	return &EventConfigUsecase{
		EventCertificateDg:  eventCertificateDg,
		EventRegistrationDg: eventRegistrationDg,
	}
}
