package event

import (
	authDg "apps/backend/core-api/internal/datagateway"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/auth"
	"apps/backend/services/s3"
	"log/slog"
)

type EventUsecase struct {
	EventDataGateway           datagateway.EventDataGateway
	EventContractDataGateway   datagateway.EventContractDataGateway
	EventIssuerDataGateway     datagateway.EventIssuerDataGateway
	AuthenticationCredentialDg authDg.AuthenticationCredentialDataGateway

	S3Service   *s3.S3Service
	logger      *slog.Logger
	authService *auth.AuthService
}

func NewEventUsecase(eventDataGateway datagateway.EventDataGateway, eventContractDataGateway datagateway.EventContractDataGateway, eventIssuerDataGateway datagateway.EventIssuerDataGateway, authenticationCredentialDg authDg.AuthenticationCredentialDataGateway, s3Service *s3.S3Service, logger *slog.Logger, authService *auth.AuthService) *EventUsecase {
	return &EventUsecase{
		EventDataGateway:           eventDataGateway,
		EventContractDataGateway:   eventContractDataGateway,
		EventIssuerDataGateway:     eventIssuerDataGateway,
		AuthenticationCredentialDg: authenticationCredentialDg,
		S3Service:                  s3Service,
		logger:                     logger,
		authService:                authService,
	}
}
