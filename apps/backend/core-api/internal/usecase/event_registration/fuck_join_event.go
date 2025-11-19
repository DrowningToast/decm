package event_registration

import (
	"apps/backend/common/customerror"
	eventContract "apps/backend/contracts/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type FuckJoinEventPayload struct {
	FirstName           *string `json:"first_name,omitempty"`
	LastName            *string `json:"last_name,omitempty"`
	Email               *string `json:"email,omitempty"`
	PhoneNumber         *string `json:"phone_number,omitempty"`
	AcademicInstitution *string `json:"academic_institution,omitempty"`
	AcademicEmail       *string `json:"academic_email,omitempty"`
	Address             *string `json:"address,omitempty"`
	Bio                 *string `json:"bio,omitempty"`
	PinCode             *string `json:"pin_code,omitempty"`
}

func (uc *EventRegistrationUsecase) FuckJoinEvent(ctx context.Context, currentUser *auth.JwtClaims, eventId uuid.UUID, payload FuckJoinEventPayload) (*entity.EventAttendee, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	dbEventContracts, err := uc.EventContractDg.GetEventContractByEventID(ctx, eventId)
	if err != nil {
		return nil, err
	}

	eventContractAddress := common.HexToAddress(dbEventContracts.EventContractAddress)
	if eventContractAddress == (common.Address{}) {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event contract not found"))
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return nil, err
	}

	auth, err := cyptoutils.GetKeyedTransactor()
	if err != nil {
		return nil, err
	}

	privateKey, participantAddress, err := cyptoutils.DecryptPrivateKey(
		*credential.EncryptedPrivateKey,
		*payload.PinCode,
	)
	if err != nil {
		return nil, err
	}
	instance, err := eventContract.NewEvent(eventContractAddress, client)
	if err != nil {
		return nil, err
	}

	calculatedDeadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
	if err != nil {
		return nil, err
	}

	signMessage, err := cyptoutils.GetSignMessage(*participantAddress, eventContractAddress, calculatedDeadlineBlock)
	if err != nil {
		return nil, err
	}

	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	tx, err := instance.AddParticipant(
		auth,
		*participantAddress,
		signMessage,
		signature,
	)
	if err != nil {
		return nil, err
	}

	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
