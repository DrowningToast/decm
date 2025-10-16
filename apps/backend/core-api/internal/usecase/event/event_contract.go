package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"apps/backend/core-api/internal/repositories/postgres"
	"decm-database/go/generated"
)

type EventContractUsecase interface {
	CreateEventContract(ctx context.Context, eventID uuid.UUID, params CreateEventContractParams) (*generated.EventContract, error)
	GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventContract, error)
	UpdateEventContract(ctx context.Context, eventID uuid.UUID, params UpdateEventContractParams) (*generated.EventContract, error)
	DeleteEventContract(ctx context.Context, eventID uuid.UUID) error
}

type EventContractUsecaseImpl struct {
	repo *postgres.Repository
}

func NewEventContractUsecase(repo *postgres.Repository) EventContractUsecase {
	return &EventContractUsecaseImpl{
		repo: repo,
	}
}

type CreateEventContractParams struct {
	AccessManagerContractAddress string
	EventContractAddress         string
	TicketContractAddress        pgtype.Text
	CertificateContractAddress   pgtype.Text
}

func (u *EventContractUsecaseImpl) CreateEventContract(ctx context.Context, eventID uuid.UUID, params CreateEventContractParams) (*generated.EventContract, error) {
	// Check if contract already exists for this event
	existingContract, err := u.repo.GetEventContractByEventID(ctx, u.repo.GetDB(), eventID)
	if err == nil && existingContract != nil {
		return nil, fmt.Errorf("event contract already exists for event ID: %s", eventID.String())
	}

	// Create new contract
	createParams := generated.CreateEventContractParams{
		EventID:                      eventID,
		AccessManagerContractAddress: params.AccessManagerContractAddress,
		EventContractAddress:         params.EventContractAddress,
		TicketContractAddress:        params.TicketContractAddress,
		CertificateContractAddress:   params.CertificateContractAddress,
	}

	return u.repo.CreateEventContract(ctx, u.repo.GetDB(), createParams)
}

func (u *EventContractUsecaseImpl) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*generated.EventContract, error) {
	return u.repo.GetEventContractByEventID(ctx, u.repo.GetDB(), eventID)
}

type UpdateEventContractParams struct {
	AccessManagerContractAddress string
	EventContractAddress         string
	TicketContractAddress        pgtype.Text
	CertificateContractAddress   pgtype.Text
}

func (u *EventContractUsecaseImpl) UpdateEventContract(ctx context.Context, eventID uuid.UUID, params UpdateEventContractParams) (*generated.EventContract, error) {
	updateParams := generated.UpdateEventContractParams{
		EventID:                      eventID,
		AccessManagerContractAddress: params.AccessManagerContractAddress,
		EventContractAddress:         params.EventContractAddress,
		TicketContractAddress:        params.TicketContractAddress,
		CertificateContractAddress:   params.CertificateContractAddress,
	}

	return u.repo.UpdateEventContract(ctx, u.repo.GetDB(), updateParams)
}

func (u *EventContractUsecaseImpl) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return u.repo.DeleteEventContract(ctx, u.repo.GetDB(), eventID)
}
