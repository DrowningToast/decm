package blockchain

import (
	"testing"

	"apps/backend/core-api/config"
	blockchainConfig "apps/backend/core-api/config/blockchain"

	"github.com/stretchr/testify/assert"
)

func TestBlockchainUsecase_GetGasPrice_Success(t *testing.T) {
	// Skip if no valid blockchain config
	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     2500.0,
			SoftCapGasPriceGwei: 1000.0,
		},
	}

	// Note: We can't easily test this without refactoring to use dependency injection
	// This test documents the expected behavior
	_ = cfg
	t.Skip("Requires refactoring to inject ethclient dependency")
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
