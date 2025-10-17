package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (uc *EventUsecase) DeleteEvent(ctx context.Context, id uuid.UUID, currentUser *auth.JwtClaims) (*entity.Event, error) {
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialById(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	dbEvent, err := uc.EventDataGateway.GetEventById(ctx, id)
	if err != nil {
		return nil, err
	}

	if dbEvent == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event not found"))
	}

	if credential.Id != dbEvent.OwnerCredentialID {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not the owner of the event"))
	}

	err = uc.EventDataGateway.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
