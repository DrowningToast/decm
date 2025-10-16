package event

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/core-api/internal/repositories/postgres"
	"decm-database/go/generated"
)

type EventIssuerUsecase interface {
	CreateEventIssuer(ctx context.Context, params CreateEventIssuerParams) (*generated.EventIssuer, error)
	GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error)
	GetEventIssuerByID(ctx context.Context, id uuid.UUID) (*generated.EventIssuer, error)
	UpdateEventIssuer(ctx context.Context, id uuid.UUID, params UpdateEventIssuerParams) (*generated.EventIssuer, error)
	DeleteEventIssuer(ctx context.Context, id uuid.UUID) error
}

type EventIssuerUsecaseImpl struct {
	repo *postgres.Repository
}

func NewEventIssuerUsecase(repo *postgres.Repository) EventIssuerUsecase {
	return &EventIssuerUsecaseImpl{
		repo: repo,
	}
}

type CreateEventIssuerParams struct {
	EventID            uuid.UUID
	IssuerCredentialID uuid.UUID
	IsSigned           int32
	Signature          pgtype.Text
}

func (u *EventIssuerUsecaseImpl) CreateEventIssuer(ctx context.Context, params CreateEventIssuerParams) (*generated.EventIssuer, error) {
	createParams := generated.CreateEventIssuerParams{
		EventID:            params.EventID,
		IssuerCredentialID: params.IssuerCredentialID,
		IsSigned:           params.IsSigned,
		Signature:          params.Signature,
	}

	return u.repo.CreateEventIssuer(ctx, u.repo.GetDB(), createParams)
}

func (u *EventIssuerUsecaseImpl) GetEventIssuersByEventID(ctx context.Context, eventID uuid.UUID) ([]generated.EventIssuer, error) {
	return u.repo.GetEventIssuersByEventID(ctx, u.repo.GetDB(), eventID)
}

func (u *EventIssuerUsecaseImpl) GetEventIssuerByID(ctx context.Context, id uuid.UUID) (*generated.EventIssuer, error) {
	return u.repo.GetEventIssuerByID(ctx, u.repo.GetDB(), id)
}

type UpdateEventIssuerParams struct {
	IsSigned  int32
	Signature pgtype.Text
}

func (u *EventIssuerUsecaseImpl) UpdateEventIssuer(ctx context.Context, id uuid.UUID, params UpdateEventIssuerParams) (*generated.EventIssuer, error) {
	updateParams := generated.UpdateEventIssuerParams{
		ID:        id,
		IsSigned:  params.IsSigned,
		Signature: params.Signature,
	}

	return u.repo.UpdateEventIssuer(ctx, u.repo.GetDB(), updateParams)
}

func (u *EventIssuerUsecaseImpl) DeleteEventIssuer(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteEventIssuer(ctx, u.repo.GetDB(), id)
}
