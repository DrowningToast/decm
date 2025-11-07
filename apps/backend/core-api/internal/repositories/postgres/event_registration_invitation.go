package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

var _ datagateway.EventRegistrationInvitationDataGateway = (*Repository)(nil)

func (r *Repository) CreateEventRegistrationInvitation(ctx context.Context, params datagateway.CreateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	result, err := r.queries.CreateEventRegistrationInvitation(ctx, generated.CreateEventRegistrationInvitationParams{
		EventID:        params.EventID,
		InboxMessageID: params.InboxMessageID,
		ValidUntil:     pgmapper.TimePtrToPgTimestampz(params.ValidUntil),
		Code:           pgmapper.StringPtrToPgText(params.Code),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventRegistrationInvitation{
		ID:             result.ID,
		EventID:        result.EventID,
		InboxMessageID: result.InboxMessageID,
		ValidUntil:     pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:           pgmapper.PgTextToStringPtr(result.Code),
		CreatedAt:      *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:      *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:    pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) GetEventRegistrationInvitationByID(ctx context.Context, id uuid.UUID) (*entity.EventRegistrationInvitation, error) {
	result, err := r.queries.GetEventRegistrationInvitationByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventRegistrationInvitation{
		ID:             result.ID,
		EventID:        result.EventID,
		InboxMessageID: result.InboxMessageID,
		ValidUntil:     pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:           pgmapper.PgTextToStringPtr(result.Code),
		CreatedAt:      *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:      *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:    pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) GetEventRegistrationInvitationsByEventID(ctx context.Context, eventID uuid.UUID) ([]*entity.EventRegistrationInvitation, error) {
	results, err := r.queries.GetEventRegistrationInvitationsByEventID(ctx, eventID)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	invitations := make([]*entity.EventRegistrationInvitation, len(results))
	for i, result := range results {
		invitations[i] = &entity.EventRegistrationInvitation{
			ID:             result.ID,
			EventID:        result.EventID,
			InboxMessageID: result.InboxMessageID,
			ValidUntil:     pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
			Code:           pgmapper.PgTextToStringPtr(result.Code),
			CreatedAt:      *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
			UpdatedAt:      *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
			CancelledAt:    pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
		}
	}

	return invitations, nil
}

func (r *Repository) UpdateEventRegistrationInvitation(ctx context.Context, id uuid.UUID, params datagateway.UpdateEventRegistrationInvitationParameters) (*entity.EventRegistrationInvitation, error) {
	result, err := r.queries.UpdateEventRegistrationInvitation(ctx, generated.UpdateEventRegistrationInvitationParams{
		ID:          id,
		ValidUntil:  pgmapper.TimePtrToPgTimestampz(params.ValidUntil),
		Code:        pgmapper.StringPtrToPgText(params.Code),
		CancelledAt: pgmapper.TimePtrToPgTimestampz(params.CancelledAt),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return &entity.EventRegistrationInvitation{
		ID:             result.ID,
		EventID:        result.EventID,
		InboxMessageID: result.InboxMessageID,
		ValidUntil:     pgmapper.PgTimestampzToTimePtr(result.ValidUntil),
		Code:           pgmapper.PgTextToStringPtr(result.Code),
		CreatedAt:      *pgmapper.PgTimestampzToTimePtr(result.CreatedAt),
		UpdatedAt:      *pgmapper.PgTimestampzToTimePtr(result.UpdatedAt),
		CancelledAt:    pgmapper.PgTimestampzToTimePtr(result.CancelledAt),
	}, nil
}

func (r *Repository) DeleteEventRegistrationInvitation(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEventRegistrationInvitation(ctx, id)
	if err != nil {
		return pgerrutils.ParsePgError(err)
	}
	return nil
}
