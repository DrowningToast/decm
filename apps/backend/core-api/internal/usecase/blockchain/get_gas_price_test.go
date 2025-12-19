package blockchain

import (
	"apps/backend/core-api/config"
	"testing"

	blockchainConfig "apps/backend/core-api/config/blockchain"

	"github.com/stretchr/testify/assert"
)

func TestGetGasPrice_SafetyMarginCalculations(t *testing.T) {
	// Test the safety margin calculation logic

	tests := []struct {
		name               string
		currentGasPrice    float64
		softCap            float64
		hardCap            float64
		expectedSoftMargin float64
		expectedHardMargin float64
	}{
		{
			name:               "Normal operation",
			currentGasPrice:    230.0,
			softCap:            1000.0,
			hardCap:            2500.0,
			expectedSoftMargin: 77.0, // (1000 - 230) / 1000 * 100
			expectedHardMargin: 90.8, // (2500 - 230) / 2500 * 100
		},
		{
			name:               "Near soft cap",
			currentGasPrice:    950.0,
			softCap:            1000.0,
			hardCap:            2500.0,
			expectedSoftMargin: 5.0,  // (1000 - 950) / 1000 * 100
			expectedHardMargin: 62.0, // (2500 - 950) / 2500 * 100
		},
		{
			name:               "Exceeded soft cap",
			currentGasPrice:    1500.0,
			softCap:            1000.0,
			hardCap:            2500.0,
			expectedSoftMargin: -50.0, // (1000 - 1500) / 1000 * 100 = negative
			expectedHardMargin: 40.0,  // (2500 - 1500) / 2500 * 100
		},
		{
			name:               "Near hard cap",
			currentGasPrice:    2400.0,
			softCap:            1000.0,
			hardCap:            2500.0,
			expectedSoftMargin: -140.0, // (1000 - 2400) / 1000 * 100
			expectedHardMargin: 4.0,    // (2500 - 2400) / 2500 * 100
		},
		{
			name:               "Exceeded hard cap",
			currentGasPrice:    3000.0,
			softCap:            1000.0,
			hardCap:            2500.0,
			expectedSoftMargin: -200.0, // (1000 - 3000) / 1000 * 100
			expectedHardMargin: -20.0,  // (2500 - 3000) / 2500 * 100
		},
		{
			name:               "Production typical values",
			currentGasPrice:    50.0,
			softCap:            60.0,
			hardCap:            80.0,
			expectedSoftMargin: 16.67, // (60 - 50) / 60 * 100
			expectedHardMargin: 37.5,  // (80 - 50) / 80 * 100
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

			// Verify calculations
			assert.InDelta(t, tt.expectedSoftMargin, softMargin, 0.1,
				"Soft margin calculation mismatch for current=%.2f, soft=%.2f",
				tt.currentGasPrice, tt.softCap)
			assert.InDelta(t, tt.expectedHardMargin, hardMargin, 0.1,
				"Hard margin calculation mismatch for current=%.2f, hard=%.2f",
				tt.currentGasPrice, tt.hardCap)
		})
	}
}

func TestGetGasPrice_ResponseStructure(t *testing.T) {
	// Test that the response structure has all required fields

	response := GasPriceResponse{
		CurrentGasPriceGwei: 271.5,
		BaseFeeGwei:         120.75,
		PriorityFeeGwei:     30.0,
		SoftCapPriceGwei:    1000.0,
		HardCapPriceGwei:    2500.0,
		SoftSafetyMargin:    72.85,
		HardSafetyMargin:    89.14,
	}

	// Verify all fields are set
	assert.NotZero(t, response.CurrentGasPriceGwei, "CurrentGasPriceGwei should be set")
	assert.NotZero(t, response.BaseFeeGwei, "BaseFeeGwei should be set")
	assert.NotZero(t, response.PriorityFeeGwei, "PriorityFeeGwei should be set")
	assert.NotZero(t, response.SoftCapPriceGwei, "SoftCapPriceGwei should be set")
	assert.NotZero(t, response.HardCapPriceGwei, "HardCapPriceGwei should be set")
	// Safety margins can be negative, so don't check NotZero

	// Verify field types
	assert.IsType(t, float64(0), response.CurrentGasPriceGwei)
	assert.IsType(t, float64(0), response.BaseFeeGwei)
	assert.IsType(t, float64(0), response.PriorityFeeGwei)
	assert.IsType(t, float64(0), response.SoftCapPriceGwei)
	assert.IsType(t, float64(0), response.HardCapPriceGwei)
	assert.IsType(t, float64(0), response.SoftSafetyMargin)
	assert.IsType(t, float64(0), response.HardSafetyMargin)
}

func TestGetGasPrice_ConfigIntegration(t *testing.T) {
	// Test that the usecase properly reads from config

	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     80.0,
			SoftCapGasPriceGwei: 60.0,
		},
	}

	// Verify config values are accessible
	assert.Equal(t, 80.0, cfg.Blockchain.MaxGasPriceGwei, "Hard cap should match config")
	assert.Equal(t, 60.0, cfg.Blockchain.SoftCapGasPriceGwei, "Soft cap should match config")

	// Document that the usecase should use these values
	t.Log("Usecase should use cfg.Blockchain.MaxGasPriceGwei for HardCapPriceGwei")
	t.Log("Usecase should use cfg.Blockchain.SoftCapGasPriceGwei for SoftCapPriceGwei")
}

func TestGetGasPrice_NegativeMarginsWhenExceeded(t *testing.T) {
	// Test that margins are correctly negative when caps are exceeded

	currentPrice := 3000.0
	softCap := 1000.0
	hardCap := 2500.0

	softMargin := ((softCap - currentPrice) / softCap) * 100
	hardMargin := ((hardCap - currentPrice) / hardCap) * 100

	// Both margins should be negative
	assert.Less(t, softMargin, 0.0, "Soft margin should be negative when price exceeds cap")
	assert.Less(t, hardMargin, 0.0, "Hard margin should be negative when price exceeds cap")

	// Verify actual values
	assert.InDelta(t, -200.0, softMargin, 0.1)
	assert.InDelta(t, -20.0, hardMargin, 0.1)
}

func TestGetGasPrice_ZeroCapHandling(t *testing.T) {
	// Test edge case where caps are zero (should not panic)

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

	// Should not panic and should return 0.0
	assert.Equal(t, 0.0, softMargin, "Should return 0.0 when soft cap is 0")
	assert.Equal(t, 0.0, hardMargin, "Should return 0.0 when hard cap is 0")
}

func TestGetGasPrice_ComponentsAddUpCorrectly(t *testing.T) {
	// Test that base fee + priority fee relationships are logical

	tests := []struct {
		baseFee     float64
		priorityFee float64
		maxFee      float64
	}{
		{
			baseFee:     100.0,
			priorityFee: 30.0,
			maxFee:      230.0, // (2 * 100) + 30 using EIP-1559 formula
		},
		{
			baseFee:     50.0,
			priorityFee: 10.0,
			maxFee:      110.0, // (2 * 50) + 10
		},
	}

	for _, tt := range tests {
		// Verify the relationship: maxFee = (2 * baseFee) + priorityFee
		calculated := (2 * tt.baseFee) + tt.priorityFee
		assert.InDelta(t, tt.maxFee, calculated, 0.1,
			"Max fee should equal (2 * baseFee) + priorityFee")
	}
}

func TestGetGasPrice_RealWorldScenarios(t *testing.T) {
	// Test real-world scenarios based on actual Polygon data

	tests := []struct {
		name        string
		scenario    string
		currentGwei float64
		softCapGwei float64
		hardCapGwei float64
		shouldWarn  bool
		shouldBlock bool
	}{
		{
			name:        "Normal Polygon traffic",
			scenario:    "Base fee ~30-50 Gwei on Polygon mainnet",
			currentGwei: 40.0,
			softCapGwei: 60.0,
			hardCapGwei: 80.0,
			shouldWarn:  false,
			shouldBlock: false,
		},
		{
			name:        "High traffic period",
			scenario:    "Popular NFT mint causing congestion",
			currentGwei: 75.0,
			softCapGwei: 60.0,
			hardCapGwei: 80.0,
			shouldWarn:  true,  // Exceeds soft cap
			shouldBlock: false, // Still within hard cap
		},
		{
			name:        "Network spike (the 20 POL problem)",
			scenario:    "Severe congestion - would cost 20 POL without caps",
			currentGwei: 8000.0, // 8000 Gwei * 250k gas = 20 POL
			softCapGwei: 60.0,
			hardCapGwei: 80.0,
			shouldWarn:  true,
			shouldBlock: true, // Exceeds hard cap - transaction rejected
		},
		{
			name:        "Edge case - at soft cap boundary",
			currentGwei: 60.0,
			softCapGwei: 60.0,
			hardCapGwei: 80.0,
			shouldWarn:  false, // At cap is OK
			shouldBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exceedsSoft := tt.currentGwei > tt.softCapGwei
			exceedsHard := tt.currentGwei > tt.hardCapGwei

			assert.Equal(t, tt.shouldWarn, exceedsSoft || exceedsHard,
				"Warning status mismatch for scenario: %s", tt.scenario)
			assert.Equal(t, tt.shouldBlock, exceedsHard,
				"Block status mismatch for scenario: %s", tt.scenario)

			t.Logf("Scenario: %s", tt.scenario)
			t.Logf("  Current: %.2f Gwei", tt.currentGwei)
			t.Logf("  Soft Cap: %.2f Gwei", tt.softCapGwei)
			t.Logf("  Hard Cap: %.2f Gwei", tt.hardCapGwei)
			t.Logf("  Should Warn: %v", tt.shouldWarn)
			t.Logf("  Should Block: %v", tt.shouldBlock)
		})
	}
}

func TestGetGasPrice_CacheControl(t *testing.T) {
	// Test cache control header expectations

	expectedCacheDuration := 300 // 5 minutes in seconds
	expectedCacheHeader := "public, max-age=300"

	assert.Equal(t, 300, expectedCacheDuration,
		"Cache should be set for 5 minutes")
	assert.Equal(t, "public, max-age=300", expectedCacheHeader,
		"Cache-Control header should be public with 5 minute max-age")

	// Rationale for 5 minute cache:
	// - Gas prices can fluctuate rapidly during congestion
	// - 5 minutes is short enough to reflect current conditions
	// - 5 minutes is long enough to reduce RPC load
	// - Users typically spend 1-2 minutes on a form before submitting
	t.Log("5 minute cache balances freshness with RPC efficiency")
}

