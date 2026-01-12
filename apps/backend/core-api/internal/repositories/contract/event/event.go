package contract

import (
	"apps/backend/common/customerror"
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

type EventContractRepository struct {
	client             *ethclient.Client
	contract           *eventContract.Event
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
}

func NewEventContractRepository(client *ethclient.Client, contract *eventContract.Event, blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway) *EventContractRepository {
	return &EventContractRepository{
		client:             client,
		contract:           contract,
		blockchainClientDg: blockchainClientDg,
	}
}

var _ eventcontract_datagateway.EventContractDataGateway = (*EventContractRepository)(nil)

func (r *EventContractRepository) GetCurrentSeatsCount(ctx context.Context) (*int64, error) {
	currentSeatsCount, err := r.contract.CurrentSeatsCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get current seats count"))
	}
	currentSeatsCountPtr := currentSeatsCount.Int64()
	return &currentSeatsCountPtr, nil
}

func (r *EventContractRepository) GetMaxSeatsCount(ctx context.Context) (*int64, error) {
	maxSeatsCount, err := r.contract.SeatsCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get max seats count"))
	}
	maxSeatsCountPtr := maxSeatsCount.Int64()
	return &maxSeatsCountPtr, nil
}

func (r *EventContractRepository) GetParticipants(ctx context.Context) ([]common.Address, error) {
	participants, err := r.contract.GetParticipants(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get participants"))
	}
	return participants, nil
}

func (r *EventContractRepository) AddParticipant(ctx context.Context, participantAddress common.Address, signMessage string, signature []byte) error {
	transactor, err := r.blockchainClientDg.GetTransactOpts(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get transaction options"))
	}

	tx, err := r.contract.AddParticipant(transactor, participantAddress, signMessage, signature)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to add participant"))
	}

	receipt, err := bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to wait for transaction to be mined"))
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return customerror.Parse(&customerror.ErrInternalServer, errors.New("transaction reverted"))
	}

	return nil
}

func (r *EventContractRepository) UpdateEvent(ctx context.Context, eventName string, eventDescription string, seatsCount int64, signMessage string, signature []byte) error {
	transactor, err := r.blockchainClientDg.GetTransactOpts(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get transaction options"))
	}

	tx, err := r.contract.UpdateEvent(transactor, eventName, eventDescription, big.NewInt(seatsCount), 0, signMessage, signature)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to update event"))
	}

	receipt, err := bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to wait for transaction to be mined"))
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return customerror.Parse(&customerror.ErrInternalServer, errors.New("transaction reverted"))
	}

	return nil
}
