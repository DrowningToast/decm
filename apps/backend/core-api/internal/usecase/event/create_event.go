package event

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"mime/multipart"
	"time"
)

type CreateEventParameters struct {
	Name             string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	SeatsCount       int
	ContactNumber    string
	ContactAddress   string
	Location         string
	GoogleMapQuery   string
	EventBanner      *multipart.FileHeader
	EventIcon        *multipart.FileHeader
}

func (uc *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParameters) (*entity.Event, error) {
	// bannerStorageKey, err := uc.UploadEventBanner(ctx, uuid.New(), params.EventBanner)
	// if err != nil {
	// 	return nil, err
	// }

	// iconStorageKey, err := uc.UploadEventIcon(ctx, uuid.New(), params.EventIcon)
	// if err != nil {
	// 	return nil, err
	// }

	return &entity.Event{}, nil
}
