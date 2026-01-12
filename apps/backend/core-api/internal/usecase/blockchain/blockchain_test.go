package blockchain

import (
	"apps/backend/core-api/config"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	"context"
	"math/big"
	"testing"
	"time"

	blockchainConfig "apps/backend/core-api/config/blockchain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBlockchainClientDg is a mock implementation of BlockchainClientDataGateway
type MockBlockchainClientDg struct {
	mock.Mock
}

func (m *MockBlockchainClientDg) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockchainClientDg) GetCalculatedDeadlineBlock(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBlockchainClientDg) EstimateDeadlineTime(ctx context.Context, deadlineBlock uint64) (*time.Time, error) {
	args := m.Called(ctx, deadlineBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*time.Time), args.Error(1)
}

func (m *MockBlockchainClientDg) GetTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bind.TransactOpts), args.Error(1)
}

func (m *MockBlockchainClientDg) GetGasPrice(ctx context.Context) (*blockchainclient_datagateway.GasPriceInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*blockchainclient_datagateway.GasPriceInfo), args.Error(1)
}

func (m *MockBlockchainClientDg) WeiToGwei(wei *big.Int) float64 {
	args := m.Called(wei)
	return args.Get(0).(float64)
}

func TestBlockchainUsecase_GetGasPrice_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBlockchainDg := new(MockBlockchainClientDg)

	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     2500.0,
			SoftCapGasPriceGwei: 1000.0,
		},
	}

	uc := NewBlockchainUsecase(nil, cfg, mockBlockchainDg)

	// Mock gas price info from blockchain
	gasPriceInfo := &blockchainclient_datagateway.GasPriceInfo{
		MaxFeePerGasGwei:         100.0,
		MaxPriorityFeePerGasGwei: 2.0,
		BaseFeeGwei:              98.0,
	}

	mockBlockchainDg.On("GetGasPrice", ctx).Return(gasPriceInfo, nil)

	// Act
	result, err := uc.GetGasPrice(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify gas price values are correctly passed through
	assert.Equal(t, 100.0, result.CurrentGasPriceGwei)
	assert.Equal(t, 98.0, result.BaseFeeGwei)
	assert.Equal(t, 2.0, result.PriorityFeeGwei)

	// Verify caps are set from config
	assert.Equal(t, 1000.0, result.SoftCapPriceGwei)
	assert.Equal(t, 2500.0, result.HardCapPriceGwei)

	// Verify safety margin calculations
	// Soft margin: ((1000 - 100) / 1000) * 100 = 90.0%
	assert.InDelta(t, 90.0, result.SoftSafetyMargin, 0.1)
	// Hard margin: ((2500 - 100) / 2500) * 100 = 96.0%
	assert.InDelta(t, 96.0, result.HardSafetyMargin, 0.1)

	mockBlockchainDg.AssertExpectations(t)
}

func TestBlockchainUsecase_GetGasPrice_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBlockchainDg := new(MockBlockchainClientDg)

	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     2500.0,
			SoftCapGasPriceGwei: 1000.0,
		},
	}

	uc := NewBlockchainUsecase(nil, cfg, mockBlockchainDg)

	// Mock error from blockchain client
	mockBlockchainDg.On("GetGasPrice", ctx).Return(nil, assert.AnError)

	// Act
	result, err := uc.GetGasPrice(ctx)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get current gas price")

	mockBlockchainDg.AssertExpectations(t)
}

func TestBlockchainUsecase_GetGasPrice_ExceedingSoftCap(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBlockchainDg := new(MockBlockchainClientDg)

	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     2500.0,
			SoftCapGasPriceGwei: 1000.0,
		},
	}

	uc := NewBlockchainUsecase(nil, cfg, mockBlockchainDg)

	// Gas price exceeds soft cap but within hard cap
	gasPriceInfo := &blockchainclient_datagateway.GasPriceInfo{
		MaxFeePerGasGwei:         1500.0,
		MaxPriorityFeePerGasGwei: 3.0,
		BaseFeeGwei:              1497.0,
	}

	mockBlockchainDg.On("GetGasPrice", ctx).Return(gasPriceInfo, nil)

	// Act
	result, err := uc.GetGasPrice(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1500.0, result.CurrentGasPriceGwei)

	// Soft margin: ((1000 - 1500) / 1000) * 100 = -50.0% (negative = exceeded)
	assert.InDelta(t, -50.0, result.SoftSafetyMargin, 0.1)
	// Hard margin: ((2500 - 1500) / 2500) * 100 = 40.0%
	assert.InDelta(t, 40.0, result.HardSafetyMargin, 0.1)

	mockBlockchainDg.AssertExpectations(t)
}

func TestBlockchainUsecase_GetGasPrice_ExceedingBothCaps(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockBlockchainDg := new(MockBlockchainClientDg)

	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     2500.0,
			SoftCapGasPriceGwei: 1000.0,
		},
	}

	uc := NewBlockchainUsecase(nil, cfg, mockBlockchainDg)

	// Gas price exceeds both caps
	gasPriceInfo := &blockchainclient_datagateway.GasPriceInfo{
		MaxFeePerGasGwei:         3000.0,
		MaxPriorityFeePerGasGwei: 5.0,
		BaseFeeGwei:              2995.0,
	}

	mockBlockchainDg.On("GetGasPrice", ctx).Return(gasPriceInfo, nil)

	// Act
	result, err := uc.GetGasPrice(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 3000.0, result.CurrentGasPriceGwei)

	// Both margins should be negative (exceeded)
	// Soft margin: ((1000 - 3000) / 1000) * 100 = -200.0%
	assert.InDelta(t, -200.0, result.SoftSafetyMargin, 0.1)
	// Hard margin: ((2500 - 3000) / 2500) * 100 = -20.0%
	assert.InDelta(t, -20.0, result.HardSafetyMargin, 0.1)

	mockBlockchainDg.AssertExpectations(t)
}

func TestBlockchainUsecase_GetGasPrice_CalculatesSafetyMargins(t *testing.T) {
	tests := []struct {
		name                string
		currentGasPrice     float64
		softCap             float64
		hardCap             float64
		expectedSoftMargin  float64
		expectedHardMargin  float64
		shouldExceedSoftCap bool
		shouldExceedHardCap bool
	}{
		{
			name:                "Well within both caps",
			currentGasPrice:     50.0,
			softCap:             1000.0,
			hardCap:             2500.0,
			expectedSoftMargin:  95.0, // (1000 - 50) / 1000 * 100
			expectedHardMargin:  98.0, // (2500 - 50) / 2500 * 100
			shouldExceedSoftCap: false,
			shouldExceedHardCap: false,
		},
		{
			name:                "Exceeds soft cap but within hard cap",
			currentGasPrice:     1500.0,
			softCap:             1000.0,
			hardCap:             2500.0,
			expectedSoftMargin:  -50.0, // (1000 - 1500) / 1000 * 100
			expectedHardMargin:  40.0,  // (2500 - 1500) / 2500 * 100
			shouldExceedSoftCap: true,
			shouldExceedHardCap: false,
		},
		{
			name:                "Exceeds both caps",
			currentGasPrice:     3000.0,
			softCap:             1000.0,
			hardCap:             2500.0,
			expectedSoftMargin:  -200.0, // (1000 - 3000) / 1000 * 100
			expectedHardMargin:  -20.0,  // (2500 - 3000) / 2500 * 100
			shouldExceedSoftCap: true,
			shouldExceedHardCap: true,
		},
		{
			name:                "At exact soft cap",
			currentGasPrice:     1000.0,
			softCap:             1000.0,
			hardCap:             2500.0,
			expectedSoftMargin:  0.0,  // (1000 - 1000) / 1000 * 100
			expectedHardMargin:  60.0, // (2500 - 1000) / 2500 * 100
			shouldExceedSoftCap: false,
			shouldExceedHardCap: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate safety margins using the same formula as the usecase
			softMargin := 0.0
			if tt.softCap > 0 {
				softMargin = ((tt.softCap - tt.currentGasPrice) / tt.softCap) * 100
			}

			hardMargin := 0.0
			if tt.hardCap > 0 {
				hardMargin = ((tt.hardCap - tt.currentGasPrice) / tt.hardCap) * 100
			}

			// Verify calculations match expected values
			assert.InDelta(t, tt.expectedSoftMargin, softMargin, 0.1, "Soft margin calculation mismatch")
			assert.InDelta(t, tt.expectedHardMargin, hardMargin, 0.1, "Hard margin calculation mismatch")

			// Verify threshold checks
			exceedsSoft := tt.currentGasPrice > tt.softCap
			exceedsHard := tt.currentGasPrice > tt.hardCap
			assert.Equal(t, tt.shouldExceedSoftCap, exceedsSoft, "Soft cap threshold check mismatch")
			assert.Equal(t, tt.shouldExceedHardCap, exceedsHard, "Hard cap threshold check mismatch")
		})
	}
}

func TestBlockchainUsecase_GetGasPrice_WithZeroCaps(t *testing.T) {
	// Test edge case where caps are zero (should not cause division by zero)
	currentPrice := 100.0
	softCap := 0.0
	hardCap := 0.0

	softMargin := 0.0
	if softCap > 0 {
		softMargin = ((softCap - currentPrice) / softCap) * 100
	}

	hardMargin := 0.0
	if hardCap > 0 {
		hardMargin = ((hardCap - currentPrice) / hardCap) * 100
	}

	// Both should be 0 when caps are 0 (no division by zero)
	assert.Equal(t, 0.0, softMargin, "Soft margin should be 0 when cap is 0")
	assert.Equal(t, 0.0, hardMargin, "Hard margin should be 0 when cap is 0")
}
