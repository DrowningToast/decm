package event

import (
	"context"
	"log/slog"

	authDg "apps/backend/core-api/internal/datagateway"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"apps/backend/services/s3"

	"github.com/google/uuid"
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

func (u *EventUsecase) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	events, err := u.EventDataGateway.ListEventsByOwnerCredentialID(ctx, ownerCredentialID, limitCount, offsetCount)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (u *EventUsecase) GetEventById(ctx context.Context, eventId uuid.UUID) (*entity.Event, error) {
	event, err := u.EventDataGateway.GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return event, nil
}
