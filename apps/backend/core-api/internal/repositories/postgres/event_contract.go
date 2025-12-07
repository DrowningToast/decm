package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common/pgmapper"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
)

var _ datagateway.EventContractDataGateway = (*Repository)(nil)

func mapGeneratedToEntityEventContract(generatedEventContract *generated.EventContract) *entity.EventContract {
	return &entity.EventContract{
		ID:                           generatedEventContract.ID,
		EventID:                      generatedEventContract.EventID,
		AccessManagerContractAddress: generatedEventContract.AccessManagerContractAddress,
		EventContractAddress:         generatedEventContract.EventContractAddress,
		TicketContractAddress:        pgmapper.PgTextToStringPtr(generatedEventContract.TicketContractAddress),
		CertificateContractAddress:   pgmapper.PgTextToStringPtr(generatedEventContract.CertificateContractAddress),
		CreatedAt:                    generatedEventContract.CreatedAt.Time,
		UpdatedAt:                    generatedEventContract.UpdatedAt.Time,
	}
}

func (r *Repository) CreateEventContract(ctx context.Context, params generated.CreateEventContractParams) (*entity.EventContract, error) {
	result, err := r.queries.CreateEventContract(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapGeneratedToEntityEventContract(&result), nil
}

func (r *Repository) GetEventContractByEventID(ctx context.Context, eventID uuid.UUID) (*entity.EventContract, error) {
	result, err := r.queries.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return mapGeneratedToEntityEventContract(&result), nil
}

func (r *Repository) UpdateEventContract(ctx context.Context, params generated.UpdateEventContractParams) (*entity.EventContract, error) {
	result, err := r.queries.UpdateEventContract(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapGeneratedToEntityEventContract(&result), nil
}

func (r *Repository) DeleteEventContract(ctx context.Context, eventID uuid.UUID) error {
	return r.queries.DeleteEventContract(ctx, eventID)
}
