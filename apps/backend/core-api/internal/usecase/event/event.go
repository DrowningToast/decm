package event

import (
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/services/auth"
	"apps/backend/services/s3"
	"log/slog"
)

type EventUsecase struct {
	EventDataGateway         datagateway.EventDataGateway
	EventContractDataGateway datagateway.EventContractDataGateway
	EventIssuerDataGateway   datagateway.EventIssuerDataGateway

	S3Service   *s3.S3Service
	logger      *slog.Logger
	authService *auth.AuthService
}

func NewEventUsecase(eventDataGateway datagateway.EventDataGateway, eventContractDataGateway datagateway.EventContractDataGateway, eventIssuerDataGateway datagateway.EventIssuerDataGateway, s3Service *s3.S3Service, logger *slog.Logger, authService *auth.AuthService) *EventUsecase {
	return &EventUsecase{
		EventDataGateway:         eventDataGateway,
		EventContractDataGateway: eventContractDataGateway,
		EventIssuerDataGateway:   eventIssuerDataGateway,
		S3Service:                s3Service,
		logger:                   logger,
		authService:              authService,
	}
}
