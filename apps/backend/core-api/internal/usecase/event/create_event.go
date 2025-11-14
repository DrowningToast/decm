package event

import (
	"context"
	"errors"
	"math/big"
	"mime/multipart"
	"time"

	"apps/backend/common/customerror"
	datagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"

	eventAccessManagerContract "apps/backend/contracts/accessmanager"
	eventContract "apps/backend/contracts/event"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

type CreateEventParameters struct {
	Name             string
	ShortDescription string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	SeatsCount       int
	ContactNumber    string
	ContactAddress   string
	Location         string
	GoogleMapQuery   string
	EventBanner      *multipart.FileHeader
	EventIcon        *multipart.FileHeader
	HostPassword     string
}

func (uc *EventUsecase) CreateEvent(ctx context.Context, params CreateEventParameters, currentUser *auth.JwtClaims) (*entity.Event, common.Address, common.Address, *types.Transaction, error) {
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	isVerifiedOrganizer := credential.IsVerifiedOrganizer
	if !isVerifiedOrganizer {
		return nil, common.Address{}, common.Address{}, nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("user is not a verified organizer"))
	}

	bannerStorageKey, err := uc.UploadEventBanner(ctx, uuid.New(), params.EventBanner)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	iconStorageKey, err := uc.UploadEventIcon(ctx, uuid.New(), params.EventIcon)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	// Get chain ID from blockchain configuration
	chainId := uc.cfg.Blockchain.ChainID

	createEventParams := datagateway.CreateEventParameters{
		ChainId:                  chainId,
		Name:                     params.Name,
		ShortDescription:         params.ShortDescription,
		Description:              params.Description,
		StartDate:                params.StartDate,
		EndDate:                  params.EndDate,
		SeatsCount:               params.SeatsCount,
		ContactNumber:            params.ContactNumber,
		ContactAddress:           params.ContactAddress,
		Location:                 params.Location,
		GoogleMapQuery:           params.GoogleMapQuery,
		BannerStorageKey:         bannerStorageKey,
		IconStorageKey:           iconStorageKey,
		OwnerCredentialID:        currentUser.UserId,
		IsPublic:                 true,  // Default to public, should be parameterized
		IsBookingRequestRequired: false, // Default to false, should be parameterized
		IsVerified:               false, // Default to false, should be parameterized
		IsTicketTransferable:     false, // Default to false, should be parameterized
		EventStatus:              entity.EventStatus(entity.EventStatusActive),
	}

	event, err := uc.EventDataGateway.CreateEvent(ctx, createEventParams)
	if err != nil {
		uc.S3Service.DeleteFile(ctx, bannerStorageKey)
		uc.S3Service.DeleteFile(ctx, iconStorageKey)
		return nil, common.Address{}, common.Address{}, nil, err
	}

	decmAccessManagerAddress := common.HexToAddress(uc.cfg.Blockchain.DecmAccessManagerAddress)
	if decmAccessManagerAddress == (common.Address{}) {
		return nil, common.Address{}, common.Address{}, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("DecmAccessManagerAddress is required"))
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	auth, err := cyptoutils.GetKeyedTransactor()
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	_, hostAddress, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, params.HostPassword)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	eventAccessManagerAddress, tx, _, err := eventAccessManagerContract.DeployEventAccessManager(
		auth,
		client,
		decmAccessManagerAddress,
		*hostAddress,
	)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	// Wait for transaction to be mined
	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	eventContractAddress, tx, _, err := eventContract.DeployEvent(
		auth,
		client,
		eventAccessManagerAddress,             // DecmAccessManager address
		event.Title,                           // Event name
		event.ShortDescription,                // Event description
		big.NewInt(int64(event.MaxAttendees)), // Seats count
	)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	// Wait for transaction to be mined
	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	return event, eventAccessManagerAddress, eventContractAddress, tx, nil
}
