package event

import (
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/services/s3"
	"log/slog"
)

type EventUsecase struct {
	EventDataGateway datagateway.EventDataGateway
	S3Service        *s3.S3Service
	logger           *slog.Logger
}

func NewEventUsecase(eventDataGateway datagateway.EventDataGateway, s3Service *s3.S3Service, logger *slog.Logger) *EventUsecase {
	return &EventUsecase{
		EventDataGateway: eventDataGateway,
		S3Service:        s3Service,
		logger:           logger,
	}
}
