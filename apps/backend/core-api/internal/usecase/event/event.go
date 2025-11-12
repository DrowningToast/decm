package event

import (
	"log/slog"

	authDg "apps/backend/core-api/internal/datagateway"
	datagateway "apps/backend/core-api/internal/datagateway"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/auth"
	"apps/backend/services/s3"
)

type EventUsecase struct {
	EventDataGateway                     eventdatagateway.EventDataGateway
	EventContractDataGateway             eventdatagateway.EventContractDataGateway
	EventIssuerDataGateway               eventdatagateway.EventIssuerDataGateway
	EventCertificateDataGateway          eventdatagateway.EventCertificateDataGateway
	EventCertificateSignatureDataGateway eventdatagateway.EventCertificateSignatureDataGateway
	AuthenticationCredentialDg           authDg.AuthenticationCredentialDataGateway
	EventRegistrationInvitationDg        datagateway.EventRegistrationInvitationDataGateway
	EventAttendeeDg                      datagateway.EventAttendeeDataGateway
	S3Service                            *s3.S3Service
	logger                               *slog.Logger
	authService                          *auth.AuthService
}

func NewEventUsecase(eventDataGateway eventdatagateway.EventDataGateway, eventContractDataGateway eventdatagateway.EventContractDataGateway, eventIssuerDataGateway eventdatagateway.EventIssuerDataGateway, eventCertificateDataGateway eventdatagateway.EventCertificateDataGateway, eventCertificateSignatureDataGateway eventdatagateway.EventCertificateSignatureDataGateway, authenticationCredentialDg authDg.AuthenticationCredentialDataGateway, eventRegistrationInvitationDg datagateway.EventRegistrationInvitationDataGateway, s3Service *s3.S3Service, logger *slog.Logger, authService *auth.AuthService) *EventUsecase {
	return &EventUsecase{
		EventDataGateway:                     eventDataGateway,
		EventContractDataGateway:             eventContractDataGateway,
		EventIssuerDataGateway:               eventIssuerDataGateway,
		EventCertificateDataGateway:          eventCertificateDataGateway,
		EventCertificateSignatureDataGateway: eventCertificateSignatureDataGateway,
		AuthenticationCredentialDg:           authenticationCredentialDg,
		EventRegistrationInvitationDg:        eventRegistrationInvitationDg,
		S3Service:                            s3Service,
		logger:                               logger,
		authService:                          authService,
	}
}
