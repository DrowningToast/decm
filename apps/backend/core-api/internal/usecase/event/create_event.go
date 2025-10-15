package event

import (
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"context"

	"github.com/google/uuid"
)

func (uc *EventUsecase) CreateEvent(ctx context.Context, params datagateway.CreateEventParameters) (*entity.Event, error) {
	bannerStorageKey, err := uc.UploadEventBanner(ctx, uuid.New(), params.EventBanner)
	if err != nil {
		return nil, err
	}

	iconStorageKey, err := uc.UploadEventIcon(ctx, uuid.New(), params.EventIcon)
	if err != nil {
		return nil, err
	}

	return &entity.Event{}, nil
}
