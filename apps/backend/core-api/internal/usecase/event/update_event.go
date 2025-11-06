package event

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"mime/multipart"
	"time"

	"apps/backend/common/customerror"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	eventContract "apps/backend/contracts/event"
)

type UpdateEventParameters struct {
	Name             *string
	ShortDescription *string
	Description      *string
	StartDate        *time.Time
	EndDate          *time.Time
	SeatsCount       *int
	ContactNumber    *string
	ContactAddress   *string
	Location         *string
	GoogleMapQuery   *string
	EventBanner      *multipart.FileHeader
	EventIcon        *multipart.FileHeader
	HostPassword     string
}

func (uc *EventUsecase) UpdateEvent(ctx context.Context, id uuid.UUID, params UpdateEventParameters, currentUser *auth.JwtClaims) (*entity.Event, error) {
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
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

	if *params.SeatsCount < dbEvent.MaxAttendees {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("seats count is less than the max attendees"))
	}

	newEventBannerStorageKey := dbEvent.BannerStorageKey // Keep existing by default
	newEventIconStorageKey := dbEvent.IconStorageKey     // Keep existing by default

	// Only upload new banner if a new file is provided
	if params.EventBanner != nil {
		previousBannerStorageKey := dbEvent.BannerStorageKey
		err := uc.S3Service.DeleteFile(ctx, previousBannerStorageKey)
		if err != nil {
			return nil, err
		}

		newBannerStorageKey, err := uc.UploadEventBanner(ctx, uuid.New(), params.EventBanner)
		if err != nil {
			return nil, err
		}
		newEventBannerStorageKey = newBannerStorageKey
	}

	// Only upload new icon if a new file is provided
	if params.EventIcon != nil {
		previousIconStorageKey := dbEvent.IconStorageKey

		newIconStorageKey, err := uc.UploadEventIcon(ctx, uuid.New(), params.EventIcon)
		if err != nil {
			return nil, err
		}
		newEventIconStorageKey = newIconStorageKey

		err = uc.S3Service.DeleteFile(ctx, previousIconStorageKey)
		if err != nil {
			return nil, err
		}

	}

	updateEventParams := datagateway.UpdateEventParameters{
		Name:             params.Name,
		ShortDescription: params.ShortDescription,
		Description:      params.Description,
		StartDate:        params.StartDate,
		EndDate:          params.EndDate,
		SeatsCount:       params.SeatsCount,
		ContactNumber:    params.ContactNumber,
		ContactAddress:   params.ContactAddress,
		Location:         params.Location,
		GoogleMapQuery:   params.GoogleMapQuery,
		BannerStorageKey: &newEventBannerStorageKey,
		IconStorageKey:   &newEventIconStorageKey,
	}

	dbEventContracts, err := uc.EventContractDataGateway.GetEventContractByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	eventContractAddress := common.HexToAddress(dbEventContracts.EventContractAddress)
	if eventContractAddress == (common.Address{}) {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("event contract not found"))
	}

	event, err := uc.EventDataGateway.UpdateEvent(ctx, id, updateEventParams)
	if err != nil {
		return nil, err
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return nil, err
	}

	auth, err := cyptoutils.GetKeyedTransactor()
	if err != nil {
		return nil, err
	}

	privateKey, hostAddress, err := cyptoutils.DecryptPrivateKey(
		*credential.EncryptedPrivateKey,
		params.HostPassword,
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

	// Hash the message with Ethereum prefix for signing
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(signMessage), signMessage)
	messageHash := crypto.Keccak256Hash([]byte(prefixedMessage))

	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	// Adjust the recovery ID (v) to be compatible with Solidity's ECDSA.recover
	// Go-ethereum returns 0/1, but Solidity expects 27/28
	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	uc.logger.Info("Signature", "signature", hex.EncodeToString(signature[:]))
	uc.logger.Info("Message Hash", "messageHash", messageHash.String())
	uc.logger.Info("Host Address", "hostAddress", *hostAddress)

	tx, err := instance.UpdateEvent(
		auth,
		*params.Name,
		*params.Description,
		big.NewInt(int64(*params.SeatsCount)),
		signMessage, // Pass the original message, not the hash
		signature,
	)
	if err != nil {
		return nil, err
	}

	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, err
	}

	return event, nil
}
