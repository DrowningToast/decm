package contract

import (
	"apps/backend/common/customerror"
	eventAccessManagerContract "apps/backend/contracts/accessmanager"
	eventContract "apps/backend/contracts/event"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"
	"context"
	"math/big"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventContractFactoryRepository struct {
	client             *ethclient.Client
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
}

func NewEventContractFactoryRepository(client *ethclient.Client, blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway) *EventContractFactoryRepository {
	return &EventContractFactoryRepository{
		client:             client,
		blockchainClientDg: blockchainClientDg,
	}
}

var _ eventcontract_datagateway.EventContractFactoryDataGateway = (*EventContractFactoryRepository)(nil)

func (r *EventContractFactoryRepository) GetContract(address common.Address) (eventcontract_datagateway.EventContractDataGateway, error) {
	eventContractInstance, err := eventContract.NewEvent(address, r.client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return NewEventContractRepository(r.client, eventContractInstance, r.blockchainClientDg), nil
}

func (r *EventContractFactoryRepository) CreateContract(ctx context.Context, params eventcontract_datagateway.CreateContractParams) (*eventcontract_datagateway.CreateContractResponse, error) {
	transactor, err := r.blockchainClientDg.GetTransactOpts(ctx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// create access manager
	eventAccessManagerAddress, tx, _, err := eventAccessManagerContract.DeployEventAccessManager(
		transactor,
		r.client,
		params.AccessManagerContractAddress,
		params.HostAddress,
	)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	receipt, err := bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("transaction reverted"))
	}

	// create event
	eventContractAddress, tx, _, err := eventContract.DeployEvent(
		transactor,
		r.client,
		eventAccessManagerAddress,            // DecmAccessManager address
		params.EventName,                     // Event name
		params.EventDescription,              // Event description
		big.NewInt(int64(params.SeatsCount)), // Seats count
	)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to deploy event contract"))
	}
	receipt, err = bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to wait for transaction to be mined"))
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("transaction reverted"))
	}

	return &eventcontract_datagateway.CreateContractResponse{
		EventContractAddress:         eventContractAddress,
		AccessManagerContractAddress: eventAccessManagerAddress,
	}, nil
}
