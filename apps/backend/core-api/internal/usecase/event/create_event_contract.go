package event

import (
	"apps/backend/common/customerror"
	"context"
	"errors"

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
	if err != nil {
		// Check if error is a not-found sentinel
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			if customErr.Code != nil && *customErr.Code == customerror.ErrNotFound.Code {
				// Not found is expected, continue with creation
				existingContract = nil
			} else {
				// Return other custom errors as-is
				return nil, err
			}
		} else {
			// Wrap non-custom errors using customerror.New
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
	}

	// Check if contract already exists
	if err == nil && existingContract != nil {
		return nil, customerror.Parse(&customerror.ErrDuplicateEntry, errors.New("event contract already exists for event ID: "+eventID.String()))
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
