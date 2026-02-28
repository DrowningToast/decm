package event

import (
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateEventContractParams struct {
	AccessManagerContractAddress string
	EventContractAddress         string
	TicketContractAddress        pgtype.Text
	CertificateContractAddress   pgtype.Text
}

func (u *EventUsecase) UpdateEventContract(ctx context.Context, eventID uuid.UUID, params UpdateEventContractParams) (*entity.EventContract, error) {
	updateParams := generated.UpdateEventContractParams{
		EventID:                      eventID,
		AccessManagerContractAddress: params.AccessManagerContractAddress,
		EventContractAddress:         params.EventContractAddress,
		TicketContractAddress:        params.TicketContractAddress,
		CertificateContractAddress:   params.CertificateContractAddress,
	}

	return u.EventContractDataGateway.UpdateEventContract(ctx, updateParams)
}
