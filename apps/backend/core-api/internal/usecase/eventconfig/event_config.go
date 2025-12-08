package eventconfig

import (
	"log/slog"

	"apps/backend/core-api/internal/datagateway"
	eventDg "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/s3"
)

type EventConfigUsecase struct {
	AuthenticationCredentialDg           datagateway.AuthenticationCredentialDataGateway
	EventDg                              eventDg.EventDataGateway
	EventDataGateway                     eventDg.EventDataGateway
	EventCertificateDg                   eventDg.EventCertificateConfigDataGateway
	EventCertificateDataGateway          eventDg.EventCertificateDataGateway
	EventCertificateSignatureDataGateway eventDg.EventCertificateSignatureDataGateway
	EventRegistrationDg                  eventDg.EventRegistrationConfigDataGateway
	EventIssuerDg                        eventDg.EventIssuerDataGateway
	EventContractDg                      eventDg.EventContractDataGateway
	InboxMessageDg                       datagateway.InboxMessageDataGateway
	S3Service                            s3.S3Service

	logger *slog.Logger
}

func NewEventConfigUsecase(
	authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway,
	eventDg eventDg.EventDataGateway,
	eventCertificateDg eventDg.EventCertificateConfigDataGateway,
	eventCertificateDataGateway eventDg.EventCertificateDataGateway,
	eventCertificateSignatureDataGateway eventDg.EventCertificateSignatureDataGateway,
	eventRegistrationDg eventDg.EventRegistrationConfigDataGateway,
	eventIssuerDg eventDg.EventIssuerDataGateway,
	eventContractDg eventDg.EventContractDataGateway,
	inboxMessageDg datagateway.InboxMessageDataGateway,
	s3Service s3.S3Service,
	logger *slog.Logger,
) *EventConfigUsecase {
	return &EventConfigUsecase{
		AuthenticationCredentialDg:           authenticationCredentialDg,
		EventDg:                              eventDg,
		EventDataGateway:                     eventDg,
		EventCertificateDg:                   eventCertificateDg,
		EventCertificateDataGateway:          eventCertificateDataGateway,
		EventCertificateSignatureDataGateway: eventCertificateSignatureDataGateway,
		EventRegistrationDg:                  eventRegistrationDg,
		EventIssuerDg:                        eventIssuerDg,
		EventContractDg:                      eventContractDg,
		InboxMessageDg:                       inboxMessageDg,
		S3Service:                            s3Service,
		logger:                               logger,
	}
}
