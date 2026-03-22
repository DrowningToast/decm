package blockchain

import (
	"context"
	"errors"
	"math/big"
	"testing"

	blockchainConfig "apps/backend/core-api/config/blockchain"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testPrivateKey is a well-known Hardhat/Anvil test account private key.
// It must never be used with real funds.
const testPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// newTestRepo builds a BlockchainClientRepository with no RPC client (nil),
// replacing getGasPriceFn with the provided stub so RPC calls are never made.
func newTestRepo(cfg *blockchainConfig.BlockchainConfig, gasPriceFn func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error)) *BlockchainClientRepository {
	r := &BlockchainClientRepository{
		client:        nil,
		config:        cfg,
		getGasPriceFn: gasPriceFn,
	}
	return r
}

// ─── Invalid private key ──────────────────────────────────────────────────────

func TestGetTransactOpts_InvalidPrivateKey(t *testing.T) {
	cfg := &blockchainConfig.BlockchainConfig{
		PrivateKey:      "not-valid-hex",
		ChainID:         1337,
		MaxGasPriceGwei: 80,
	}
	r := newTestRepo(cfg, nil) // getGasPriceFn never reached

	_, err := r.GetTransactOpts(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create private key")
}

// ─── Gas price cap disabled (MaxGasPriceGwei = 0) ────────────────────────────

func TestGetTransactOpts_NoCap(t *testing.T) {
	cfg := &blockchainConfig.BlockchainConfig{
		PrivateKey:      testPrivateKey,
		ChainID:         1337,
		MaxGasPriceGwei: 0, // disabled
	}
	called := false
	r := newTestRepo(cfg, func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
		called = true
		return nil, nil
	})

	auth, err := r.GetTransactOpts(context.Background())

	require.NoError(t, err)
	require.NotNil(t, auth)
	assert.False(t, called, "getGasPriceFn should not be called when cap is disabled")
	assert.Nil(t, auth.GasTipCap, "GasTipCap should not be set when cap is disabled")
	assert.Nil(t, auth.GasFeeCap, "GasFeeCap should not be set when cap is disabled")
	assert.Equal(t, uint64(0), auth.GasLimit, "GasLimit is always 0 (auto-estimate)")
}

// ─── Gas price fetch failure ──────────────────────────────────────────────────

func TestGetTransactOpts_GasPriceFetchError(t *testing.T) {
	cfg := &blockchainConfig.BlockchainConfig{
		PrivateKey:      testPrivateKey,
		ChainID:         1337,
		MaxGasPriceGwei: 80,
	}
	r := newTestRepo(cfg, func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
		return nil, errors.New("rpc connection refused")
	})

	_, err := r.GetTransactOpts(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "gas price check failed")
	assert.Contains(t, err.Error(), "rpc connection refused")
}

// ─── Market fee exceeds hard cap ──────────────────────────────────────────────

func TestGetTransactOpts_MarketFeeExceedsCap(t *testing.T) {
	cfg := &blockchainConfig.BlockchainConfig{
		PrivateKey:      testPrivateKey,
		ChainID:         1337,
		MaxGasPriceGwei: 0.5,
	}
	r := newTestRepo(cfg, func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
		return &blockchainclient_datagateway.GasPriceInfo{
			MaxFeePerGasGwei:         1.0, // above cap
			MaxPriorityFeePerGasGwei: 0.02,
			BaseFeeGwei:              0.49,
		}, nil
	})

	_, err := r.GetTransactOpts(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "transaction rejected")
	assert.Contains(t, err.Error(), "1.00 Gwei")
	assert.Contains(t, err.Error(), "0.50 Gwei")
}

// ─── Market fee exactly at cap (boundary: >= means reject) ───────────────────

func TestGetTransactOpts_MarketFeeAtExactCap(t *testing.T) {
	cfg := &blockchainConfig.BlockchainConfig{
		PrivateKey:      testPrivateKey,
		ChainID:         1337,
		MaxGasPriceGwei: 0.1,
	}
	r := newTestRepo(cfg, func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
		return &blockchainclient_datagateway.GasPriceInfo{
			MaxFeePerGasGwei:         0.1, // exactly at cap
			MaxPriorityFeePerGasGwei: 0.02,
			BaseFeeGwei:              0.04,
		}, nil
	})

	_, err := r.GetTransactOpts(context.Background())

	require.Error(t, err, "should reject when market fee equals cap (>= check)")
	assert.Contains(t, err.Error(), "transaction rejected")
}

// ─── Happy path: uses actual market rate, not the cap ────────────────────────

func TestGetTransactOpts_SetsActualMarketRate_NotCap(t *testing.T) {
	const capGwei = 80.0
	const marketFeeGwei = 0.1
	const marketTipGwei = 0.05

	cfg := &blockchainConfig.BlockchainConfig{
		PrivateKey:      testPrivateKey,
		ChainID:         1337,
		MaxGasPriceGwei: capGwei,
	}
	r := newTestRepo(cfg, func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
		return &blockchainclient_datagateway.GasPriceInfo{
			MaxFeePerGasGwei:         marketFeeGwei,
			MaxPriorityFeePerGasGwei: marketTipGwei,
			BaseFeeGwei:              0.025,
		}, nil
	})

	auth, err := r.GetTransactOpts(context.Background())

	require.NoError(t, err)
	require.NotNil(t, auth)

	expectedFeeCap := gweiToWei(marketFeeGwei)
	expectedTipCap := gweiToWei(marketTipGwei)
	wrongFeeCap := gweiToWei(capGwei)

	assert.Equal(t, expectedFeeCap, auth.GasFeeCap, "GasFeeCap must be the market rate, not the hard cap")
	assert.Equal(t, expectedTipCap, auth.GasTipCap, "GasTipCap must be the actual network-suggested tip")
	assert.NotEqual(t, wrongFeeCap, auth.GasFeeCap, "GasFeeCap must NOT be set to MaxGasPriceGwei")
}

// ─── Fee cap values round-trip correctly through Gwei↔Wei conversion ─────────

func TestGetTransactOpts_FeeCapConversionAccuracy(t *testing.T) {
	tests := []struct {
		name           string
		marketFeeGwei  float64
		marketTipGwei  float64
		capGwei        float64
		expectedFeeWei *big.Int
		expectedTipWei *big.Int
	}{
		{
			name:           "Sepolia low gas (0.1 Gwei fee, 0.05 Gwei tip)",
			marketFeeGwei:  0.1,
			marketTipGwei:  0.05,
			capGwei:        80,
			expectedFeeWei: big.NewInt(100_000_000),
			expectedTipWei: big.NewInt(50_000_000),
		},
		{
			name:           "Moderate gas (5 Gwei fee, 1 Gwei tip)",
			marketFeeGwei:  5.0,
			marketTipGwei:  1.0,
			capGwei:        80,
			expectedFeeWei: big.NewInt(5_000_000_000),
			expectedTipWei: big.NewInt(1_000_000_000),
		},
		{
			name:           "High gas (50 Gwei fee, 2 Gwei tip) — under 80 Gwei cap",
			marketFeeGwei:  50.0,
			marketTipGwei:  2.0,
			capGwei:        80,
			expectedFeeWei: big.NewInt(50_000_000_000),
			expectedTipWei: big.NewInt(2_000_000_000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &blockchainConfig.BlockchainConfig{
				PrivateKey:      testPrivateKey,
				ChainID:         1337,
				MaxGasPriceGwei: tt.capGwei,
			}
			r := newTestRepo(cfg, func(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
				return &blockchainclient_datagateway.GasPriceInfo{
					MaxFeePerGasGwei:         tt.marketFeeGwei,
					MaxPriorityFeePerGasGwei: tt.marketTipGwei,
				}, nil
			})

			auth, err := r.GetTransactOpts(context.Background())

			require.NoError(t, err)
			assert.Equal(t, 0, auth.GasFeeCap.Cmp(tt.expectedFeeWei),
				"GasFeeCap: got %s wei, want %s wei", auth.GasFeeCap, tt.expectedFeeWei)
			assert.Equal(t, 0, auth.GasTipCap.Cmp(tt.expectedTipWei),
				"GasTipCap: got %s wei, want %s wei", auth.GasTipCap, tt.expectedTipWei)
		})
	}
}
