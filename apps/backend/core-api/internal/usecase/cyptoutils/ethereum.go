package cyptoutils

import (
	"context"
	"math/big"
	"time"

	"apps/backend/common/metrics"
	"apps/backend/core-api/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func GetEthereumClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.LoadConfig().Blockchain.RPCURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ethereum client")
	}

	return client, nil
}

func GetKeyedTransactor(ctx context.Context, client *ethclient.Client) (*bind.TransactOpts, error) {
	cfg := config.LoadConfig().Blockchain
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create private key")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(cfg.ChainID)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create keyed transactor")
	}

	// 1. Determine effective context
	effectiveCtx := ctx
	if effectiveCtx == nil {
		effectiveCtx = context.Background()
	}

	// 2. Apply Gas Caps if configured
	if cfg.MaxGasPriceGwei > 0 {
		effectiveClient := client
		var isLocalClient bool

		if effectiveClient == nil {
			var err error
			effectiveClient, err = GetEthereumClient()
			if err != nil {
				return nil, errors.Wrap(err, "failed to create ethereum client for mandatory gas price check")
			}
			isLocalClient = true
		}

		// Ensure local client is closed after the check
		if isLocalClient {
			defer effectiveClient.Close()
		}

		// Perform mandatory gas price check
		gasInfo, err := GetCurrentGasPrice(effectiveCtx, effectiveClient)
		if err != nil {
			return nil, errors.Wrap(err, "gas price check failed: action rejected for safety")
		}

		baseFee := gasInfo.BaseFeeGwei
		maxPrice := cfg.MaxGasPriceGwei

		// REJECT: If base fee alone is higher than our hard cap, we cannot proceed safely
		if baseFee >= maxPrice {
			return nil, errors.Errorf("transaction rejected: current network base fee (%.2f Gwei) exceeds your hard cap (%.2f Gwei)", baseFee, maxPrice)
		}

		// Calculate tip: Fill the gap between base fee and our cap
		tip := maxPrice - baseFee
		if tip < 1.0 {
			tip = 1.0 // Minimum 1 Gwei tip
		}

		auth.GasTipCap = gweiToWei(tip)
		auth.GasFeeCap = gweiToWei(maxPrice)
	}

	// Set explicit gas limit if configured
	if cfg.GasLimit > 0 {
		auth.GasLimit = cfg.GasLimit
	}

	return auth, nil
}

func gweiToWei(gwei float64) *big.Int {
	weiFloat := new(big.Float).Mul(big.NewFloat(gwei), big.NewFloat(1e9))
	weiInt, _ := weiFloat.Int(nil)
	return weiInt
}

func WeiToGwei(wei *big.Int) float64 {
	if wei == nil {
		return 0
	}
	f := new(big.Float).SetInt(wei)
	f.Quo(f, big.NewFloat(1e9))
	gwei, _ := f.Float64()
	return gwei
}

type GasPriceInfo struct {
	MaxFeePerGasGwei         float64 `json:"max_fee_per_gas_gwei"`
	MaxPriorityFeePerGasGwei float64 `json:"max_priority_fee_per_gas_gwei"`
	BaseFeeGwei              float64 `json:"base_fee_gwei"`
}

func GetCurrentGasPrice(ctx context.Context, client *ethclient.Client) (*GasPriceInfo, error) {
	start := time.Now()
	operation := "get_current_gas_price"

	// 1. Get SuggestPriorityFee
	tipCap, err := client.SuggestGasTipCap(ctx)

	if err != nil {
		metrics.BlockchainOperationsTotal.WithLabelValues(operation, "error").Inc()
		metrics.BlockchainOperationDuration.WithLabelValues(operation, "error").Observe(time.Since(start).Seconds())
		return nil, errors.Wrap(err, "failed to suggest gas tip cap")
	}

	// 2. Get latest block header for Base Fee
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		metrics.BlockchainOperationsTotal.WithLabelValues(operation, "error").Inc()
		metrics.BlockchainOperationDuration.WithLabelValues(operation, "error").Observe(time.Since(start).Seconds())
		return nil, errors.Wrap(err, "failed to get latest block header")
	}

	baseFee := header.BaseFee
	if baseFee == nil {
		// Fallback for non-EIP1559 networks
		suggestedPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			metrics.BlockchainOperationsTotal.WithLabelValues(operation, "error").Inc()
			metrics.BlockchainOperationDuration.WithLabelValues(operation, "error").Observe(time.Since(start).Seconds())
			return nil, errors.Wrap(err, "failed to suggest gas price")
		}
		
		info := &GasPriceInfo{
			MaxFeePerGasGwei:         WeiToGwei(suggestedPrice),
			MaxPriorityFeePerGasGwei: WeiToGwei(tipCap),
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

	info := &GasPriceInfo{
		MaxFeePerGasGwei:         WeiToGwei(maxFee),
		MaxPriorityFeePerGasGwei: WeiToGwei(tipCap),
		BaseFeeGwei:              WeiToGwei(baseFee),
	}

	metrics.BlockchainOperationsTotal.WithLabelValues(operation, "success").Inc()
	metrics.BlockchainOperationDuration.WithLabelValues(operation, "success").Observe(time.Since(start).Seconds())

	// Record Gas Prices
	metrics.BlockchainGasPrice.WithLabelValues("max_fee").Set(info.MaxFeePerGasGwei)
	metrics.BlockchainGasPrice.WithLabelValues("priority_fee").Set(info.MaxPriorityFeePerGasGwei)
	metrics.BlockchainGasPrice.WithLabelValues("base_fee").Set(info.BaseFeeGwei)

	return info, nil
}
