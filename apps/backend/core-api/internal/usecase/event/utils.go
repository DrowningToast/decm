package event

import (
	"apps/backend/core-api/internal/entity"
	"context"
)

func (u *EventUsecase) ToEventResponse(ctx context.Context, event *entity.Event) (*EventResponse, error) {
	bannerPresignedURL, err := u.S3DataGateway.GetPresignedURL(ctx, event.BannerStorageKey)
	if err != nil {
		return nil, err
	}
	iconPresignedURL, err := u.S3DataGateway.GetPresignedURL(ctx, event.IconStorageKey)
	if err != nil {
		return nil, err
	}

	return &EventResponse{
		Id:                       event.Id,
		ChainID:                  int32(event.ChainId),
		ContactNumber:            event.ContactNumber,
		ContactAddress:           event.ContactAddress,
		OwnerCredentialID:        event.OwnerCredentialId,
		BannerPresignedURL:       bannerPresignedURL,
		IconPresignedURL:         iconPresignedURL,
		Title:                    event.Title,
		ShortDescription:         event.ShortDescription,
		LongDescription:          event.LongDescription,
		StartDate:                event.StartDate,
		EndDate:                  event.EndDate,
		Location:                 event.Location,
		GoogleMapQuery:           event.GoogleMapQuery,
		MaxAttendees:             int32(event.MaxAttendees),
		AttendeesCount:           int32(event.AttendeesCount),
		IsPublic:                 event.IsPublic,
		IsBookingRequestRequired: event.IsBookingRequestRequired,
		IsVerified:               event.IsVerified,
		IsTicketTransferable:     event.IsTicketTransferable,
		CreatedAt:                event.CreatedAt,
		UpdatedAt:                event.UpdatedAt,
		EventStatus:              event.EventStatus,
		EventType:                event.EventType,
	}, nil
}
