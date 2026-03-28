package event

import (
	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	eventContract "apps/backend/contracts/event"

	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

type DeleteEventParameters struct {
	HostPassword string
	Signature    string
	SignMessage  string
}

func (uc *EventUsecase) DeleteEvent(ctx context.Context, id uuid.UUID, currentUser *auth.JwtClaims, params DeleteEventParameters) (*entity.Event, error) {
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

	transactor, err := uc.BlockchainClientDg.GetTransactOpts(ctx)
	if err != nil {
		return nil, err
	}

	var signature []byte
	var signMessage string

	if params.Signature != "" && params.SignMessage != "" {
		// Wallet path: use provided signature and sign message
		sigHex := params.Signature
		if len(sigHex) >= 2 && sigHex[:2] == "0x" {
			sigHex = sigHex[2:]
		}
		sig, err := hex.DecodeString(sigHex)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid signature format"))
		}
		signature = sig
		signMessage = params.SignMessage

		// Verify signature matches credential wallet address
		recoveredAddress, err := cyptoutils.GetAddressFromSignature(signMessage, params.Signature)
		if err != nil {
			return nil, err
		}
		credentialAddress := common.HexToAddress(credential.WalletAddress)
		if recoveredAddress != credentialAddress {
			return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("wallet signature does not match credential wallet address"))
		}
	} else if params.HostPassword != "" {
		if credential.EncryptedPrivateKey == nil {
			return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("host_password is not supported for BYOK users; provide signature and sign_message instead"))
		}
		privateKey, hostAddress, err := cyptoutils.DecryptPrivateKey(
			*credential.EncryptedPrivateKey,
			params.HostPassword,
		)
		if err != nil {
			return nil, err
		}

		calculatedDeadlineBlock, err := uc.BlockchainClientDg.GetCalculatedDeadlineBlock(ctx)
		if err != nil {
			return nil, err
		}

		signMessage, err = cyptoutils.GetSignMessage(*hostAddress, eventContractAddress, calculatedDeadlineBlock)
		if err != nil {
			return nil, err
		}

		messageHash := cyptoutils.HashEthereumMessage(signMessage)
		signature, err = cyptoutils.Sign(messageHash.Bytes(), privateKey)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("either host_password or (signature + sign_message) is required"))
	}

	instance, err := eventContract.NewEvent(eventContractAddress, uc.ethClient)
	if err != nil {
		return nil, err
	}

	tx, err := instance.UpdateEvent(
		transactor,
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

	receipt, err := bind.WaitMined(ctx, uc.ethClient, tx)
	if err != nil {
		return nil, err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("transaction reverted (tx=%s)", tx.Hash().Hex()))
	}

	_, err = uc.EventDataGateway.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return dbEvent, nil
}
