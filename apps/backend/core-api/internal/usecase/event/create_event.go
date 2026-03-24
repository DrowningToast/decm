package event

import (
	"apps/backend/common/customerror"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	"apps/backend/core-api/internal/entity"
	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"errors"
	"mime/multipart"
	"time"

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
	// Wallet-based auth (alternative to HostPassword)
	Signature  string
	SignMessage string
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

	// Resolve host address early (before file uploads) so we fail fast on bad auth
	var hostAddress common.Address
	if params.Signature != "" {
		// Wallet path: recover address from signature and verify against stored wallet address
		recoveredAddress, err := cyptoutils.GetAddressFromSignature(params.SignMessage, params.Signature)
		if err != nil {
			return nil, common.Address{}, common.Address{}, nil, err
		}
		credentialAddress := common.HexToAddress(credential.WalletAddress)
		if recoveredAddress != credentialAddress {
			return nil, common.Address{}, common.Address{}, nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("wallet signature does not match credential wallet address"))
		}
		hostAddress = recoveredAddress
	} else {
		// Password path: decrypt private key to get host address
		_, addr, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, params.HostPassword)
		if err != nil {
			return nil, common.Address{}, common.Address{}, nil, err
		}
		hostAddress = *addr
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

	// Validate chain ID
	if chainId <= 0 {
		return nil, common.Address{}, common.Address{}, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("invalid blockchain chain ID"))
	}
	createEventParams := event_datagateway.CreateEventParameters{
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
		_ = uc.S3DataGateway.DeleteFile(ctx, bannerStorageKey)
		_ = uc.S3DataGateway.DeleteFile(ctx, iconStorageKey)
		return nil, common.Address{}, common.Address{}, nil, err
	}

	decmAccessManagerAddress := common.HexToAddress(uc.cfg.Blockchain.DecmAccessManagerAddress)
	if decmAccessManagerAddress == (common.Address{}) {
		return nil, common.Address{}, common.Address{}, nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("DecmAccessManagerAddress is required"))
	}

	// Create smart contracts using the factory
	contractResponse, err := uc.EventContractFactoryDg.CreateContract(ctx, eventcontract_datagateway.CreateContractParams{
		AccessManagerContractAddress: decmAccessManagerAddress,
		HostAddress:                  hostAddress,
		EventName:                    event.Title,
		EventDescription:             event.ShortDescription,
		SeatsCount:                   int64(event.MaxAttendees),
	})
	if err != nil {
		return nil, common.Address{}, common.Address{}, nil, err
	}

	return event, contractResponse.AccessManagerContractAddress, contractResponse.EventContractAddress, nil, nil
}
