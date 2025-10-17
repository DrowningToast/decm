package event

import (
	"apps/backend/common/customerror"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UpdateEventParameters struct {
	Name             *string
	ShortDescription *string
	Description      *string
	StartDate        *time.Time
	EndDate          *time.Time
	SeatsCount       *int
	ContactNumber    *string
	ContactAddress   *string
	Location         *string
	GoogleMapQuery   *string
	EventBanner      *multipart.FileHeader
	EventIcon        *multipart.FileHeader
}

func (uc *EventUsecase) UpdateEvent(ctx context.Context, id uuid.UUID, params UpdateEventParameters, currentUser *auth.JwtClaims) (*entity.Event, error) {
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	dbEvent, err := uc.EventDataGateway.GetEventById(ctx, id)
	if err != nil {
		return nil, err
	}

	if dbEvent == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
	}

	if credential.Id != dbEvent.OwnerCredentialID {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not the owner of the event"))
	}

	if *params.SeatsCount < dbEvent.MaxAttendees {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("seats count is less than the max attendees"))
	}

	newEventBannerStorageKey := dbEvent.BannerStorageKey // Keep existing by default
	newEventIconStorageKey := dbEvent.IconStorageKey     // Keep existing by default

	// Only upload new banner if a new file is provided
	if params.EventBanner != nil {
		previousBannerStorageKey := dbEvent.BannerStorageKey
		err := uc.S3Service.DeleteFile(ctx, previousBannerStorageKey)
		if err != nil {
			return nil, err
		}

		newBannerStorageKey, err := uc.UploadEventBanner(ctx, uuid.New(), params.EventBanner)
		if err != nil {
			return nil, err
		}
		newEventBannerStorageKey = newBannerStorageKey
	}

	// Only upload new icon if a new file is provided
	if params.EventIcon != nil {
		previousIconStorageKey := dbEvent.IconStorageKey

		newIconStorageKey, err := uc.UploadEventIcon(ctx, uuid.New(), params.EventIcon)
		if err != nil {
			return nil, err
		}
		newEventIconStorageKey = newIconStorageKey

		err = uc.S3Service.DeleteFile(ctx, previousIconStorageKey)
		if err != nil {
			return nil, err
		}

	}

	updateEventParams := datagateway.UpdateEventParameters{
		Name:             params.Name,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
		StartDate:        params.StartDate,
		EndDate:          params.EndDate,
		SeatsCount:       params.SeatsCount,
		ContactNumber:    params.ContactNumber,
		ContactAddress:   params.ContactAddress,
		Location:         params.Location,
		GoogleMapQuery:   params.GoogleMapQuery,
		BannerStorageKey: &newEventBannerStorageKey,
		IconStorageKey:   &newEventIconStorageKey,
	}

	event, err := uc.EventDataGateway.UpdateEvent(ctx, id, updateEventParams)
	if err != nil {
		return nil, err
	}

	return event, nil
}
