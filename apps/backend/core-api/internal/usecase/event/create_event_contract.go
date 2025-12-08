package event

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEventContractParams struct {
	AccessManagerContractAddress string
	EventContractAddress         string
	TicketContractAddress        pgtype.Text
	CertificateContractAddress   pgtype.Text
}

func (u *EventUsecase) CreateEventContract(ctx context.Context, eventID uuid.UUID, params CreateEventContractParams) (*entity.EventContract, error) {
	_, err := u.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Create new contract
	createParams := generated.CreateEventContractParams{
		EventID:                      eventID,
		AccessManagerContractAddress: params.AccessManagerContractAddress,
		EventContractAddress:         params.EventContractAddress,
		TicketContractAddress:        params.TicketContractAddress,
		CertificateContractAddress:   params.CertificateContractAddress,
	}

	return u.EventContractDataGateway.CreateEventContract(ctx, createParams)
}
