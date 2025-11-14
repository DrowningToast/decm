package eventconfig

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *EventConfigUsecase) CheckEventPassword(ctx context.Context, eventID uuid.UUID, password string) (bool, error) {
	passwordHash, err := u.EventRegistrationDg.GetEventRegistrationConfigPasswordByEventId(ctx, eventID)
	if err != nil {
		return false, errors.Wrap(err, "failed to get event registration config")
	}
	if passwordHash == nil {
		return false, customerror.NewWithPreset(&customerror.ErrInvalidArgument, errors.New("registration password is not set"))
	}

	match, err := hashutils.CompareHash(password, *passwordHash)
	if err != nil {
		return false, err
	}
	if match {
		return true, nil
	}

	return false, nil
}
