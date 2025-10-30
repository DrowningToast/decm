package event

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"apps/backend/common/customerror"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	"github.com/google/uuid"
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

func (uc *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParameters, currentUser *auth.JwtClaims) (*entity.Event, error) {
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	bannerStorageKey, err := uc.UploadEventBanner(ctx, uuid.New(), params.EventBanner)
	if err != nil {
		return nil, err
	}

	iconStorageKey, err := uc.UploadEventIcon(ctx, uuid.New(), params.EventIcon)
	if err != nil {
		return nil, err
	}

	createEventParams := datagateway.CreateEventParameters{
		Name:                     params.Name,
		ShortDescription:         params.ShortDescription,
		Description:              params.Description,
		StartDate:                params.StartDate,
		EndDate:                  params.EndDate,
		SeatsCount:               params.SeatsCount,
		ContactNumber:            params.ContactNumber,
		ContactAddress:           params.ContactAddress,
		Location:                 params.Location,
		GoogleMapQuery:           params.GoogleMapQuery,
		BannerStorageKey:         bannerStorageKey,
		IconStorageKey:           iconStorageKey,
		OwnerCredentialID:        currentUser.UserId,
		IsPublic:                 true,  // Default to public, should be parameterized
		IsBookingRequestRequired: false, // Default to false, should be parameterized
		IsVerified:               false, // Default to false, should be parameterized
		IsTicketTransferable:     false, // Default to false, should be parameterized
	}

	event, err := uc.EventDataGateway.CreateEvent(ctx, createEventParams)
	if err != nil {
		uc.S3Service.DeleteFile(ctx, bannerStorageKey)
		uc.S3Service.DeleteFile(ctx, iconStorageKey)
		return nil, err
	}

	// TODO: Deploy event contracts

	return event, nil
}
