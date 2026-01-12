package eventconfig

import (
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/services/s3"
	"log/slog"

	eventDg "apps/backend/core-api/internal/datagateway/offchain/event"
)

type EventConfigUsecase struct {
	AuthenticationCredentialDg           offchain_datagateway.AuthenticationCredentialDataGateway
	EventDg                              eventDg.EventDataGateway
	EventDataGateway                     eventDg.EventDataGateway
	EventCertificateDg                   eventDg.EventCertificateConfigDataGateway
	EventCertificateDataGateway          eventDg.EventCertificateDataGateway
	EventCertificateSignatureDataGateway eventDg.EventCertificateSignatureDataGateway
	EventCertificateFontFamilyDg         eventDg.EventCertificateFontFamilyDataGateway
	EventRegistrationDg                  eventDg.EventRegistrationConfigDataGateway
	EventIssuerDg                        eventDg.EventIssuerDataGateway
	EventContractDg                      eventDg.EventContractDataGateway
	InboxMessageDg                       offchain_datagateway.InboxMessageDataGateway
	S3Service                            s3.S3Service

	logger *slog.Logger
}

func NewEventConfigUsecase(
	authenticationCredentialDg offchain_datagateway.AuthenticationCredentialDataGateway,
	eventDg eventDg.EventDataGateway,
	eventCertificateDg eventDg.EventCertificateConfigDataGateway,
	eventCertificateDataGateway eventDg.EventCertificateDataGateway,
	eventCertificateSignatureDataGateway eventDg.EventCertificateSignatureDataGateway,
	eventCertificateFontFamilyDg eventDg.EventCertificateFontFamilyDataGateway,
	eventRegistrationDg eventDg.EventRegistrationConfigDataGateway,
	eventIssuerDg eventDg.EventIssuerDataGateway,
	eventContractDg eventDg.EventContractDataGateway,
	inboxMessageDg offchain_datagateway.InboxMessageDataGateway,
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
		EventCertificateFontFamilyDg:         eventCertificateFontFamilyDg,
		EventRegistrationDg:                  eventRegistrationDg,
		EventIssuerDg:                        eventIssuerDg,
		EventContractDg:                      eventContractDg,
		InboxMessageDg:                       inboxMessageDg,
		S3Service:                            s3Service,
		logger:                               logger,
	}
}
