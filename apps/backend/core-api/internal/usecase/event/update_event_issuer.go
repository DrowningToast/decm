package event

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"decm-database/go/generated"
)

type UpdateEventIssuerParams struct {
	EventID            uuid.UUID
	IssuerCredentialID uuid.UUID
}

func (u *EventUsecase) UpdateEventIssuer(ctx context.Context, eventID uuid.UUID, params []UpdateEventIssuerParams, currentUser *auth.JwtClaims) ([]generated.EventIssuer, error) {
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	for _, param := range params {
		eventID := param.EventID
		issuerCredentialID := param.IssuerCredentialID

		_, err := u.EventIssuerDataGateway.GetEventIssuerByEventIDAndIssuerCredentialID(ctx, eventID, issuerCredentialID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				_, err := u.EventIssuerDataGateway.CreateEventIssuer(ctx, generated.CreateEventIssuerParams{
					EventID:            eventID,
					IssuerCredentialID: issuerCredentialID,
					IsSigned:           0,
					Signature:          pgtype.Text{},
					SignMessage:        pgtype.Text{},
				})
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	return nil, nil
}
