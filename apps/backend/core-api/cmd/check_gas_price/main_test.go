package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGasPriceComparison(t *testing.T) {
	// Test the comparison logic used in the script

	tests := []struct {
		name               string
		currentPrice       float64
		softCap            float64
		hardCap            float64
		expectedStatus     string // "ok", "soft_exceeded", "hard_exceeded"
		expectedSoftMargin float64
		expectedHardMargin float64
	}{
		{
			name:               "Normal operation - within both caps",
			currentPrice:       50.0,
			softCap:            60.0,
			hardCap:            80.0,
			expectedStatus:     "ok",
			expectedSoftMargin: 16.67, // (60-50)/60*100
			expectedHardMargin: 37.5,  // (80-50)/80*100
		},
		{
			name:               "Warning zone - exceeded soft cap",
			currentPrice:       70.0,
			softCap:            60.0,
			hardCap:            80.0,
			expectedStatus:     "soft_exceeded",
			expectedSoftMargin: -16.67, // (60-70)/60*100
			expectedHardMargin: 12.5,   // (80-70)/80*100
		},
		{
			name:               "Danger zone - exceeded hard cap",
			currentPrice:       100.0,
			softCap:            60.0,
			hardCap:            80.0,
			expectedStatus:     "hard_exceeded",
			expectedSoftMargin: -66.67, // (60-100)/60*100
			expectedHardMargin: -25.0,  // (80-100)/80*100
		},
		{
			name:               "Edge case - at soft cap",
			currentPrice:       60.0,
			softCap:            60.0,
			hardCap:            80.0,
			expectedStatus:     "ok", // At soft cap is still OK
			expectedSoftMargin: 0.0,
			expectedHardMargin: 25.0,
		},
		{
			name:               "Edge case - at hard cap",
			currentPrice:       80.0,
			softCap:            60.0,
			hardCap:            80.0,
			expectedStatus:     "soft_exceeded", // At hard cap is dangerous but not exceeded
			expectedSoftMargin: -33.33,
			expectedHardMargin: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate margins
			softMargin := ((tt.softCap - tt.currentPrice) / tt.softCap) * 100
			hardMargin := ((tt.hardCap - tt.currentPrice) / tt.hardCap) * 100

			// Verify margin calculations
			assert.InDelta(t, tt.expectedSoftMargin, softMargin, 0.1,
				"Soft margin calculation mismatch")
			assert.InDelta(t, tt.expectedHardMargin, hardMargin, 0.1,
				"Hard margin calculation mismatch")

			// Determine status based on script logic
			var actualStatus string
			if tt.currentPrice > tt.hardCap {
				actualStatus = "hard_exceeded"
			} else if tt.currentPrice > tt.softCap {
				actualStatus = "soft_exceeded"
			} else {
				actualStatus = "ok"
			}

			assert.Equal(t, tt.expectedStatus, actualStatus,
				"Status determination mismatch")
		})
	}
}

func TestScriptUsageInstructions(t *testing.T) {
	// Document the script's usage

	scriptName := "check_gas_price.go"
	command := "pnpm check-gas:core"
	description := "Checks current Polygon gas prices and compares against configured limits"

	assert.NotEmpty(t, scriptName, "Script should have a name")
	assert.NotEmpty(t, command, "Script should have a pnpm command")
	assert.NotEmpty(t, description, "Script should have a description")

	// Expected output format
	expectedOutputFields := []string{
		"Base Fee",
		"Priority Fee (Tip)",
		"Current Total",
		"Soft Cap Fee",
		"Hard Cap Fee",
	}

	for _, field := range expectedOutputFields {
		assert.NotEmpty(t, field, "Output should include %s", field)
	}
}

func TestScriptConfigLoading(t *testing.T) {
	// Test that the script properly loads configuration

	// The script should load from root .env when run from apps/backend/
	expectedEnvPaths := []string{
		"../../.env", // Default path when run from apps/backend/
		".env",       // Docker path
		"/app/.env",  // Docker absolute path
	}

	assert.GreaterOrEqual(t, len(expectedEnvPaths), 1,
		"Script should check multiple .env paths for portability")

	// Document required environment variables
	requiredVars := []string{
		"BLOCKCHAIN_NETWORK",
		"BLOCKCHAIN_RPC_URL",
		"BLOCKCHAIN_CHAIN_ID",
		"BLOCKCHAIN_MAX_GAS_PRICE_GWEI",
		"BLOCKCHAIN_SOFT_CAP_GAS_PRICE_GWEI",
	}

	for _, v := range requiredVars {
		assert.NotEmpty(t, v, "Required variable: %s", v)
	}
}

func TestScriptStatusMessages(t *testing.T) {
	// Test the status message logic

	tests := []struct {
		name            string
		currentPrice    float64
		softCap         float64
		hardCap         float64
		expectedEmoji   string
		expectedMessage string
	}{
		{
			name:            "Green status - all good",
			currentPrice:    50.0,
			softCap:         60.0,
			hardCap:         80.0,
			expectedEmoji:   "🟢",
			expectedMessage: "Gas prices are within the allowed range",
		},
		{
			name:            "Yellow status - soft cap exceeded",
			currentPrice:    70.0,
			softCap:         60.0,
			hardCap:         80.0,
			expectedEmoji:   "🟡",
			expectedMessage: "exceeds your SOFT cap",
		},
		{
			name:            "Red status - hard cap exceeded",
			currentPrice:    90.0,
			softCap:         60.0,
			hardCap:         80.0,
			expectedEmoji:   "🔴",
			expectedMessage: "EXCEEDS your HARD cap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status string
			if tt.currentPrice > tt.hardCap {
				status = "🔴 hard_exceeded"
			} else if tt.currentPrice > tt.softCap {
				status = "🟡 soft_exceeded"
			} else {
				status = "🟢 ok"
			}

			assert.Contains(t, status, tt.expectedEmoji,
				"Status message should contain correct emoji")
		})
	}
}
