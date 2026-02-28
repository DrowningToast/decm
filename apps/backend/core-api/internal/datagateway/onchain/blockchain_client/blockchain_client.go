package blockchainclient_datagateway

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// GasPriceInfo contains gas price information from the blockchain
type GasPriceInfo struct {
	MaxFeePerGasGwei         float64 `json:"max_fee_per_gas_gwei"`
	MaxPriorityFeePerGasGwei float64 `json:"max_priority_fee_per_gas_gwei"`
	BaseFeeGwei              float64 `json:"base_fee_gwei"`
}

// BlockchainClientDataGateway provides blockchain client operations
type BlockchainClientDataGateway interface {
	// GetCurrentBlockNumber returns the latest block number from the blockchain
	GetCurrentBlockNumber(ctx context.Context) (uint64, error)

	// GetCalculatedDeadlineBlock returns the deadline block number (current block + 1 week worth of blocks)
	GetCalculatedDeadlineBlock(ctx context.Context) (uint64, error)

	// EstimateDeadlineTime estimates the time when a specific deadline block will be reached
	// Returns nil if the deadline block is in the past or if estimation fails
	EstimateDeadlineTime(ctx context.Context, deadlineBlock uint64) (*time.Time, error)

	// GetTransactOpts creates a transaction options object for signing blockchain transactions
	// Applies gas price caps and limits based on configuration
	GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error)

	// GetGasPrice retrieves current gas price information from the blockchain
	GetGasPrice(ctx context.Context) (*GasPriceInfo, error)

	// WeiToGwei converts wei (smallest unit) to gwei (gigawei) for gas price readability
	WeiToGwei(wei *big.Int) float64
}
