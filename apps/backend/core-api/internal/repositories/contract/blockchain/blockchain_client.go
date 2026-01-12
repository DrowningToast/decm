package blockchain

import (
	"apps/backend/common/customerror"
	"apps/backend/common/metrics"
	"apps/backend/core-api/config/blockchain"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	"context"
	"math/big"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainClientRepository struct {
	client        *ethclient.Client
	config        *blockchain.BlockchainConfig
	getGasPriceFn func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error)
}

func NewBlockchainClientRepository(client *ethclient.Client, cfg *blockchain.BlockchainConfig) *BlockchainClientRepository {
	r := &BlockchainClientRepository{
		client: client,
		config: cfg,
	}
	r.getGasPriceFn = r.GetGasPrice
	return r
}

var _ blockchainclient_datagateway.BlockchainClientDataGateway = (*BlockchainClientRepository)(nil)

// blockchainRPCTimeout is the maximum time allowed for a single blockchain RPC
// call.  Without this guard the HTTP handler goroutine blocks indefinitely when
// the RPC node is slow or unreachable, causing the request to hang until the
// Fiber global timeout fires (which can be minutes).
const blockchainRPCTimeout = 5 * time.Second

// gasPriceRPCTimeout covers the full GetGasPrice operation which makes two
// sequential RPC calls (eth_maxPriorityFeePerGas + eth_getBlockByNumber).
const gasPriceRPCTimeout = 10 * time.Second

func (r *BlockchainClientRepository) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, blockchainRPCTimeout)
	defer cancel()
	blockNumber, err := r.client.BlockNumber(rpcCtx)
	if err != nil {
		return 0, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get current block number"))
	}
	return blockNumber, nil
}

func (r *BlockchainClientRepository) GetCalculatedDeadlineBlock(ctx context.Context) (uint64, error) {
	latestBlockNumber, err := r.GetCurrentBlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	// Calculate blocks for 1 week based on chain-specific block time
	// 1 week = 604,800 seconds
	secondsInWeek := uint64(604800)
	blockTimeSeconds := uint64(r.config.GetBlockTimeSeconds())
	blocksPerWeek := secondsInWeek / blockTimeSeconds

	return latestBlockNumber + blocksPerWeek, nil
}

func (r *BlockchainClientRepository) EstimateDeadlineTime(ctx context.Context, deadlineBlock uint64) (*time.Time, error) {
	currentBlock, err := r.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	// If deadline is in the past, return nil
	if deadlineBlock <= currentBlock {
		return nil, nil
	}

	blocksRemaining := deadlineBlock - currentBlock

	// Use chain-specific block time determined by chain ID
	secondsPerBlock := r.config.GetBlockTimeSeconds()
	estimatedTime := time.Now().Add(time.Duration(blocksRemaining*uint64(secondsPerBlock)) * time.Second)

	return &estimatedTime, nil
}

func (r *BlockchainClientRepository) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(r.config.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create private key")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(r.config.ChainID)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create keyed transactor")
	}

	// Apply Gas Caps if configured
	if r.config.MaxGasPriceGwei > 0 {
		// Perform mandatory gas price check
		gasInfo, err := r.getGasPriceFn(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "gas price check failed: action rejected for safety")
		}

		// REJECT: If the computed optimal fee exceeds our hard cap, we cannot proceed safely
		if gasInfo.MaxFeePerGasGwei >= r.config.MaxGasPriceGwei {
			return nil, errors.Errorf("transaction rejected: current network fee (%.2f Gwei) exceeds your hard cap (%.2f Gwei)", gasInfo.MaxFeePerGasGwei, r.config.MaxGasPriceGwei)
		}

		// Use the actual market rate, not the cap
		auth.GasTipCap = gweiToWei(gasInfo.MaxPriorityFeePerGasGwei)
		auth.GasFeeCap = gweiToWei(gasInfo.MaxFeePerGasGwei)
	}

	return auth, nil
}

func (r *BlockchainClientRepository) GetGasPrice(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
	// Guard the entire operation — it makes 2 sequential RPC calls.
	rpcCtx, cancel := context.WithTimeout(ctx, gasPriceRPCTimeout)
	defer cancel()

	start := time.Now()
	operation := "get_current_gas_price"

	// 1. Get SuggestPriorityFee
	tipCap, err := r.client.SuggestGasTipCap(rpcCtx)
	if err != nil {
		metrics.BlockchainOperationsTotal.WithLabelValues(operation, "error").Inc()
		metrics.BlockchainOperationDuration.WithLabelValues(operation, "error").Observe(time.Since(start).Seconds())
		return nil, errors.Wrap(err, "failed to suggest gas tip cap")
	}

	// 2. Get latest block header for Base Fee
	header, err := r.client.HeaderByNumber(rpcCtx, nil)
	if err != nil {
		metrics.BlockchainOperationsTotal.WithLabelValues(operation, "error").Inc()
		metrics.BlockchainOperationDuration.WithLabelValues(operation, "error").Observe(time.Since(start).Seconds())
		return nil, errors.Wrap(err, "failed to get latest block header")
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		// Fallback for non-EIP1559 networks
		suggestedPrice, err := r.client.SuggestGasPrice(rpcCtx)
		if err != nil {
			metrics.BlockchainOperationsTotal.WithLabelValues(operation, "error").Inc()
			metrics.BlockchainOperationDuration.WithLabelValues(operation, "error").Observe(time.Since(start).Seconds())
			return nil, errors.Wrap(err, "failed to suggest gas price")
		}

		info := &blockchainclient_datagateway.GasPriceInfo{
			MaxFeePerGasGwei:         r.WeiToGwei(suggestedPrice),
			MaxPriorityFeePerGasGwei: r.WeiToGwei(tipCap),
			BaseFeeGwei:              0,
		}

		metrics.BlockchainOperationsTotal.WithLabelValues(operation, "success").Inc()
		metrics.BlockchainOperationDuration.WithLabelValues(operation, "success").Observe(time.Since(start).Seconds())

		// Record Gas Prices
		metrics.BlockchainGasPrice.WithLabelValues("max_fee").Set(info.MaxFeePerGasGwei)
		metrics.BlockchainGasPrice.WithLabelValues("priority_fee").Set(info.MaxPriorityFeePerGasGwei)
		metrics.BlockchainGasPrice.WithLabelValues("base_fee").Set(info.BaseFeeGwei)

		return info, nil
	}

	// MaxFeePerGas = (2 * BaseFee) + PriorityFee (standard formula for safety)
	maxFee := new(big.Int).Mul(baseFee, big.NewInt(2))
	maxFee.Add(maxFee, tipCap)

	info := &blockchainclient_datagateway.GasPriceInfo{
		MaxFeePerGasGwei:         r.WeiToGwei(maxFee),
		MaxPriorityFeePerGasGwei: r.WeiToGwei(tipCap),
		BaseFeeGwei:              r.WeiToGwei(baseFee),
	}

	metrics.BlockchainOperationsTotal.WithLabelValues(operation, "success").Inc()
	metrics.BlockchainOperationDuration.WithLabelValues(operation, "success").Observe(time.Since(start).Seconds())

	// Record Gas Prices
	metrics.BlockchainGasPrice.WithLabelValues("max_fee").Set(info.MaxFeePerGasGwei)
	metrics.BlockchainGasPrice.WithLabelValues("priority_fee").Set(info.MaxPriorityFeePerGasGwei)
	metrics.BlockchainGasPrice.WithLabelValues("base_fee").Set(info.BaseFeeGwei)

	return info, nil
}

func (r *BlockchainClientRepository) WeiToGwei(wei *big.Int) float64 {
	if wei == nil {
		return 0
	}
	f := new(big.Float).SetInt(wei)
	f.Quo(f, big.NewFloat(1e9))
	gwei, _ := f.Float64()
	return gwei
}

// gweiToWei is a private helper to convert gwei to wei
func gweiToWei(gwei float64) *big.Int {
	weiFloat := new(big.Float).Mul(big.NewFloat(gwei), big.NewFloat(1e9))
	weiInt, _ := weiFloat.Int(nil)
	return weiInt
}
