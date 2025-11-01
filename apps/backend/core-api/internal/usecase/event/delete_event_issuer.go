package event

import (
	"apps/backend/common/customerror"
	"apps/backend/services/auth"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (u *EventUsecase) DeleteEventIssuer(ctx context.Context, id uuid.UUID, currentUser *auth.JwtClaims) error {
	credential, err := u.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	err = u.EventIssuerDataGateway.DeleteEventIssuer(ctx, id)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return nil
}
