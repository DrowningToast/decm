package event

import (
	"context"
	"errors"
	"math/big"

	"apps/backend/common/customerror"
	eventContract "apps/backend/contracts/event"
	"apps/backend/core-api/internal/entity"
	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

func (uc *EventUsecase) DeleteEvent(ctx context.Context, id uuid.UUID, currentUser *auth.JwtClaims, hostPassword string) (*entity.Event, error) {
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	dbEvent, err := uc.EventDataGateway.GetEventById(ctx, id)
	if err != nil {
		return nil, err
	}

	if credential.Id != dbEvent.OwnerCredentialId {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not owner of the event"))
	}

	dbEventContracts, err := uc.EventContractDataGateway.GetEventContractByEventID(ctx, id)
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
	defer client.Close()

	auth, err := cyptoutils.GetKeyedTransactor(ctx, client)
	if err != nil {
		return nil, err
	}

	privateKey, hostAddress, err := cyptoutils.DecryptPrivateKey(
		*credential.EncryptedPrivateKey,
		hostPassword,
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

	signMessage, err := cyptoutils.GetSignMessage(*hostAddress, eventContractAddress, calculatedDeadlineBlock)
	if err != nil {
		return nil, err
	}

	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	tx, err := instance.UpdateEvent(
		auth,
		dbEvent.Title,
		dbEvent.LongDescription,
		big.NewInt(int64(dbEvent.MaxAttendees)),
		2,
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

	_, err = uc.EventDataGateway.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return dbEvent, nil
}
