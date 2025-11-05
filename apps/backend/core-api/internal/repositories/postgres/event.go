package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

var _ datagateway.EventDataGateway = (*Repository)(nil)

func (r *Repository) CreateEvent(ctx context.Context, params datagateway.CreateEventParameters) (*entity.Event, error) {
	// Convert boolean values to pgtype.Int4 for database
	isPublic := pgmapper.Int32ToPgInt4(0)
	isBookingRequestRequired := pgmapper.Int32ToPgInt4(0)
	isVerified := pgmapper.Int32ToPgInt4(0)
	isTicketTransferable := pgmapper.Int32ToPgInt4(0)

	// Convert description to pgtype.Text
	longDescription := pgmapper.StringPtrToPgText(&params.Description)

	result, err := r.queries.CreateEvent(ctx, generated.CreateEventParams{
		ContactNumber:            params.ContactNumber,
		ContactAddress:           params.ContactAddress,
		OwnerCredentialID:        params.OwnerCredentialID,
		BannerStorageKey:         params.BannerStorageKey,
		IconStorageKey:           params.IconStorageKey,
		Title:                    params.Name,
		ShortDescription:         params.ShortDescription,
		LongDescription:          longDescription,
		StartDate:                params.StartDate,
		EndDate:                  params.EndDate,
		Location:                 params.Location,
		GoogleMapQuery:           params.GoogleMapQuery,
		MaxAttendees:             int32(params.SeatsCount),
		IsPublic:                 isPublic,
		IsBookingRequestRequired: isBookingRequestRequired,
		IsVerified:               isVerified,
		IsTicketTransferable:     isTicketTransferable,
		EventStatus:              generated.EventStatus(params.EventStatus),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Event{
		ID:                       result.ID,
		ChainID:                  int(result.ChainID),
		ContactNumber:            result.ContactNumber,
		ContactAddress:           result.ContactAddress,
		OwnerCredentialID:        result.OwnerCredentialID,
		BannerStorageKey:         result.BannerStorageKey,
		IconStorageKey:           result.IconStorageKey,
		Title:                    result.Title,
		ShortDescription:         result.ShortDescription,
		LongDescription:          result.LongDescription.String,
		StartDate:                result.StartDate,
		EndDate:                  result.EndDate,
		Location:                 result.Location,
		GoogleMapQuery:           result.GoogleMapQuery,
		MaxAttendees:             int(result.MaxAttendees),
		IsPublic:                 result.IsPublic.Int32 == 1,
		IsBookingRequestRequired: result.IsBookingRequestRequired.Int32 == 1,
		IsVerified:               result.IsVerified.Int32 == 1,
		IsTicketTransferable:     result.IsTicketTransferable.Int32 == 1,
		CreatedAt:                result.CreatedAt.Time,
		UpdatedAt:                result.UpdatedAt.Time,
		EventStatus:              entity.EventStatus(result.EventStatus),
	}, nil
}

func (r *Repository) GetEventById(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	result, err := r.queries.GetEventById(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Event{
		ID:                       result.ID,
		ChainID:                  int(result.ChainID),
		ContactNumber:            result.ContactNumber,
		ContactAddress:           result.ContactAddress,
		OwnerCredentialID:        result.OwnerCredentialID,
		BannerStorageKey:         result.BannerStorageKey,
		IconStorageKey:           result.IconStorageKey,
		Title:                    result.Title,
		ShortDescription:         result.ShortDescription,
		LongDescription:          result.LongDescription.String,
		StartDate:                result.StartDate,
		EndDate:                  result.EndDate,
		Location:                 result.Location,
		GoogleMapQuery:           result.GoogleMapQuery,
		MaxAttendees:             int(result.MaxAttendees),
		IsPublic:                 result.IsPublic.Int32 == 1,
		IsBookingRequestRequired: result.IsBookingRequestRequired.Int32 == 1,
		IsVerified:               result.IsVerified.Int32 == 1,
		IsTicketTransferable:     result.IsTicketTransferable.Int32 == 1,
		CreatedAt:                result.CreatedAt.Time,
		UpdatedAt:                result.UpdatedAt.Time,
		EventStatus:              entity.EventStatus(result.EventStatus),
	}, nil
}

func (r *Repository) ListEventsByOwnerCredentialID(ctx context.Context, ownerCredentialID uuid.UUID, limitCount int32, offsetCount int32) ([]*entity.Event, error) {
	events, err := r.queries.ListEventsByOwnerCredentialID(ctx, generated.ListEventsByOwnerCredentialIDParams{
		OwnerCredentialID: ownerCredentialID,
		LimitCount:        limitCount,
		OffsetCount:       offsetCount,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	eventsEntity := make([]*entity.Event, len(events))
	for i, event := range events {
		eventsEntity[i] = &entity.Event{
			ID:                       event.ID,
			ChainID:                  int(event.ChainID),
			ContactNumber:            event.ContactNumber,
			ContactAddress:           event.ContactAddress,
			OwnerCredentialID:        event.OwnerCredentialID,
			BannerStorageKey:         event.BannerStorageKey,
			IconStorageKey:           event.IconStorageKey,
			Title:                    event.Title,
			ShortDescription:         event.ShortDescription,
			LongDescription:          event.LongDescription.String,
			StartDate:                event.StartDate,
			EndDate:                  event.EndDate,
			Location:                 event.Location,
			GoogleMapQuery:           event.GoogleMapQuery,
			MaxAttendees:             int(event.MaxAttendees),
			IsPublic:                 event.IsPublic.Int32 == 1,
			IsBookingRequestRequired: event.IsBookingRequestRequired.Int32 == 1,
			IsVerified:               event.IsVerified.Int32 == 1,
			IsTicketTransferable:     event.IsTicketTransferable.Int32 == 1,
			CreatedAt:                event.CreatedAt.Time,
			UpdatedAt:                event.UpdatedAt.Time,
			EventStatus:              entity.EventStatus(event.EventStatus),
		}
	}

	return eventsEntity, nil
}

func (r *Repository) DeleteEvent(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	result, err := r.queries.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Event{
		ID:                       result.ID,
		EventType:                entity.EventType(result.EventType),
		ChainID:                  int(result.ChainID),
		ContactNumber:            result.ContactNumber,
		ContactAddress:           result.ContactAddress,
		OwnerCredentialID:        result.OwnerCredentialID,
		BannerStorageKey:         result.BannerStorageKey,
		IconStorageKey:           result.IconStorageKey,
		Title:                    result.Title,
		ShortDescription:         result.ShortDescription,
		LongDescription:          result.LongDescription.String,
		StartDate:                result.StartDate,
		EndDate:                  result.EndDate,
		Location:                 result.Location,
		GoogleMapQuery:           result.GoogleMapQuery,
		MaxAttendees:             int(result.MaxAttendees),
		IsPublic:                 result.IsPublic.Int32 == 1,
		IsBookingRequestRequired: result.IsBookingRequestRequired.Int32 == 1,
		IsVerified:               result.IsVerified.Int32 == 1,
		IsTicketTransferable:     result.IsTicketTransferable.Int32 == 1,
		CreatedAt:                result.CreatedAt.Time,
		UpdatedAt:                result.UpdatedAt.Time,
		EventStatus:              entity.EventStatus(result.EventStatus),
	}, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventParameters) (*entity.Event, error) {
	// Get the current event to preserve values that are not being updated
	currentEvent, err := r.queries.GetEventById(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Convert boolean values to pgtype.Int4 for database
	isPublic := pgmapper.Int32ToPgInt4(0)
	isVerified := pgmapper.Int32ToPgInt4(0)

	// Use values from params if provided, otherwise use current values
	isBookingRequestRequired := pgmapper.Int32ToPgInt4(0)
	if params.IsBookingRequestRequired != nil {
		if *params.IsBookingRequestRequired {
			isBookingRequestRequired = pgmapper.Int32ToPgInt4(1)
		}
	} else {
		isBookingRequestRequired = currentEvent.IsBookingRequestRequired
	}

	isTicketTransferable := pgmapper.Int32ToPgInt4(0)
	if params.IsTicketTransferable != nil {
		if *params.IsTicketTransferable {
			isTicketTransferable = pgmapper.Int32ToPgInt4(1)
		}
	} else {
		isTicketTransferable = currentEvent.IsTicketTransferable
	}

	// Use event type from params if provided, otherwise use current value
	eventType := currentEvent.EventType
	if params.EventType != nil {
		eventType = generated.EventType(*params.EventType)
	}

	// Use event status from current event (not updatable through this endpoint)
	eventStatus := currentEvent.EventStatus

	// Only update fields that are provided in params
	// For fields that are not provided, use current values
	name := currentEvent.Title
	if params.Name != nil {
		name = *params.Name
	}

	shortDescription := currentEvent.ShortDescription
	if params.ShortDescription != nil {
		shortDescription = *params.ShortDescription
	}

	description := currentEvent.LongDescription.String
	if params.Description != nil {
		description = *params.Description
	}

	startDate := currentEvent.StartDate
	if params.StartDate != nil {
		startDate = *params.StartDate
	}

	endDate := currentEvent.EndDate
	if params.EndDate != nil {
		endDate = *params.EndDate
	}

	location := currentEvent.Location
	if params.Location != nil {
		location = *params.Location
	}

	googleMapQuery := currentEvent.GoogleMapQuery
	if params.GoogleMapQuery != nil {
		googleMapQuery = *params.GoogleMapQuery
	}

	maxAttendees := currentEvent.MaxAttendees
	if params.SeatsCount != nil {
		maxAttendees = int32(*params.SeatsCount)
	}

	contactNumber := currentEvent.ContactNumber
	if params.ContactNumber != nil {
		contactNumber = *params.ContactNumber
	}

	contactAddress := currentEvent.ContactAddress
	if params.ContactAddress != nil {
		contactAddress = *params.ContactAddress
	}

	bannerStorageKey := currentEvent.BannerStorageKey
	if params.BannerStorageKey != nil {
		bannerStorageKey = *params.BannerStorageKey
	}

	iconStorageKey := currentEvent.IconStorageKey
	if params.IconStorageKey != nil {
		iconStorageKey = *params.IconStorageKey
	}

	// Convert description to pgtype.Text
	longDescription := pgmapper.StringPtrToPgText(&description)

	result, err := r.queries.UpdateEvent(ctx, generated.UpdateEventParams{
		ID:                       id,
		ContactNumber:            contactNumber,
		ContactAddress:           contactAddress,
		BannerStorageKey:         bannerStorageKey,
		IconStorageKey:           iconStorageKey,
		Title:                    name,
		ShortDescription:         shortDescription,
		LongDescription:          longDescription,
		StartDate:                startDate,
		EndDate:                  endDate,
		Location:                 location,
		GoogleMapQuery:           googleMapQuery,
		MaxAttendees:             maxAttendees,
		IsPublic:                 isPublic,
		IsBookingRequestRequired: isBookingRequestRequired,
		IsVerified:               isVerified,
		IsTicketTransferable:     isTicketTransferable,
		EventType:                eventType,
		EventStatus:              eventStatus,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.Event{
		ID:                       result.ID,
		EventType:                entity.EventType(result.EventType),
		ChainID:                  int(result.ChainID),
		ContactNumber:            result.ContactNumber,
		ContactAddress:           result.ContactAddress,
		OwnerCredentialID:        result.OwnerCredentialID,
		BannerStorageKey:         result.BannerStorageKey,
		IconStorageKey:           result.IconStorageKey,
		Title:                    result.Title,
		ShortDescription:         result.ShortDescription,
		LongDescription:          result.LongDescription.String,
		StartDate:                result.StartDate,
		EndDate:                  result.EndDate,
		Location:                 result.Location,
		GoogleMapQuery:           result.GoogleMapQuery,
		MaxAttendees:             int(result.MaxAttendees),
		IsPublic:                 result.IsPublic.Int32 == 1,
		IsBookingRequestRequired: result.IsBookingRequestRequired.Int32 == 1,
		IsVerified:               result.IsVerified.Int32 == 1,
		IsTicketTransferable:     result.IsTicketTransferable.Int32 == 1,
		CreatedAt:                result.CreatedAt.Time,
		UpdatedAt:                result.UpdatedAt.Time,
		EventStatus:              entity.EventStatus(result.EventStatus),
	}, nil
}
