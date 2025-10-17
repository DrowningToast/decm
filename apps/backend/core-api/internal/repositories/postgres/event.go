package postgres

import (
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"context"

	"decm-database/go/generated"

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
		}
	}

	return eventsEntity, nil
}

func (r *Repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEvent(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateEvent(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventParameters) (*entity.Event, error) {
	// Convert boolean values to pgtype.Int4 for database
	isPublic := pgmapper.Int32ToPgInt4(0)
	isBookingRequestRequired := pgmapper.Int32ToPgInt4(0)
	isVerified := pgmapper.Int32ToPgInt4(0)
	isTicketTransferable := pgmapper.Int32ToPgInt4(0)

	// Convert description to pgtype.Text
	longDescription := pgmapper.StringPtrToPgText(params.Description)

	result, err := r.queries.UpdateEvent(ctx, generated.UpdateEventParams{
		ID:                       id,
		ContactNumber:            *params.ContactNumber,
		ContactAddress:           *params.ContactAddress,
		BannerStorageKey:         *params.BannerStorageKey,
		IconStorageKey:           *params.IconStorageKey,
		Title:                    *params.Name,
		ShortDescription:         *params.Description,
		LongDescription:          longDescription,
		StartDate:                *params.StartDate,
		EndDate:                  *params.EndDate,
		Location:                 *params.Location,
		GoogleMapQuery:           *params.GoogleMapQuery,
		MaxAttendees:             int32(*params.SeatsCount),
		IsPublic:                 isPublic,
		IsBookingRequestRequired: isBookingRequestRequired,
		IsVerified:               isVerified,
		IsTicketTransferable:     isTicketTransferable,
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
	}, nil
}
