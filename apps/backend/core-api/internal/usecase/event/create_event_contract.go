package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"decm-database/go/generated"
)

type CreateEventContractParams struct {
	AccessManagerContractAddress string
	EventContractAddress         string
	TicketContractAddress        pgtype.Text
	CertificateContractAddress   pgtype.Text
}

func (u *EventUsecase) CreateEventContract(ctx context.Context, eventID uuid.UUID, params CreateEventContractParams) (*generated.EventContract, error) {
	// Check if contract already exists for this event
	existingContract, err := u.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
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

	return u.EventContractDataGateway.CreateEventContract(ctx, createParams)
}
