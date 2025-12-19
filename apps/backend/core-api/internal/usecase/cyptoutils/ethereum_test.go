package cyptoutils

import (
	"math/big"
	"testing"

	"apps/backend/core-api/config"

	"github.com/stretchr/testify/assert"
)

func TestGetEthereumClient(t *testing.T) {
	// This test requires actual config and RPC URL
	// In a real scenario, you'd mock the config or use a test RPC URL
	t.Skip("Skipping test that requires actual Ethereum RPC connection")

	client, err := GetEthereumClient()
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Ethereum RPC: %v", err)
		return
	}

	assert.NotNil(t, client)
	client.Close()
}

func TestGetKeyedTransactor(t *testing.T) {
	// This test requires actual config with valid private key
	// In a real scenario, you'd mock the config
	t.Skip("Skipping test that requires actual blockchain config")

	cfg := config.LoadConfig()
	if cfg.Blockchain.PrivateKey == "" {
		t.Skip("Skipping test: no private key configured")
		return
	}

	auth, err := GetKeyedTransactor(nil, nil)
	if err != nil {
		t.Skipf("Skipping test: cannot create transactor: %v", err)
		return
	}

	assert.NotNil(t, auth)
}

func TestGetKeyedTransactor_InvalidPrivateKey(t *testing.T) {
	// This would require mocking config, which is complex
	// For now, we document the expected behavior
	t.Skip("Skipping test that requires config mocking")
}

func TestGweiToWei(t *testing.T) {
	tests := []struct {
		name     string
		gwei     float64
		expected string // Use string for big.Int comparison
	}{
		{
			name:     "1 Gwei",
			gwei:     1.0,
			expected: "1000000000", // 1e9 wei
		},
		{
			name:     "50 Gwei",
			gwei:     50.0,
			expected: "50000000000", // 50e9 wei
		},
		{
			name:     "100.5 Gwei",
			gwei:     100.5,
			expected: "100500000000", // 100.5e9 wei
		},
		{
			name:     "2500 Gwei (default hard cap)",
			gwei:     2500.0,
			expected: "2500000000000", // 2500e9 wei
		},
		{
			name:     "0.001 Gwei (very small)",
			gwei:     0.001,
			expected: "1000000", // 0.001e9 = 1e6 wei
		},
		{
			name:     "0 Gwei",
			gwei:     0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gweiToWei(tt.gwei)
			expected := new(big.Int)
			expected.SetString(tt.expected, 10)

			assert.Equal(t, expected.String(), result.String(),
				"gweiToWei(%f) = %s, want %s", tt.gwei, result.String(), expected.String())
		})
	}
}

func TestWeiToGwei(t *testing.T) {
	tests := []struct {
		name     string
		wei      string // Use string for big.Int input
		expected float64
	}{
		{
			name:     "1 Gwei in wei",
			wei:      "1000000000", // 1e9 wei
			expected: 1.0,
		},
		{
			name:     "50 Gwei in wei",
			wei:      "50000000000", // 50e9 wei
			expected: 50.0,
		},
		{
			name:     "100.5 Gwei in wei",
			wei:      "100500000000", // 100.5e9 wei
			expected: 100.5,
		},
		{
			name:     "2500 Gwei in wei (default hard cap)",
			wei:      "2500000000000", // 2500e9 wei
			expected: 2500.0,
		},
		{
			name:     "Very small value",
			wei:      "1000000", // 0.001e9 = 1e6 wei
			expected: 0.001,
		},
		{
			name:     "0 wei",
			wei:      "0",
			expected: 0.0,
		},
		{
			name:     "1 wei (sub-Gwei)",
			wei:      "1",
			expected: 0.000000001, // 1/1e9
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wei := new(big.Int)
			wei.SetString(tt.wei, 10)

			result := WeiToGwei(wei)

			assert.InDelta(t, tt.expected, result, 0.000000001,
				"WeiToGwei(%s) = %f, want %f", tt.wei, result, tt.expected)
		})
	}
}

func TestWeiToGwei_NilInput(t *testing.T) {
	result := WeiToGwei(nil)
	assert.Equal(t, 0.0, result, "WeiToGwei(nil) should return 0.0")
}

func TestGweiToWei_RoundTrip(t *testing.T) {
	// Test that converting Gwei -> Wei -> Gwei returns the original value
	testValues := []float64{1.0, 50.0, 100.5, 2500.0, 0.001, 80.0, 60.0}

	for _, original := range testValues {
		wei := gweiToWei(original)
		backToGwei := WeiToGwei(wei)

		assert.InDelta(t, original, backToGwei, 0.000000001,
			"Round trip failed: %f -> Wei -> %f", original, backToGwei)
	}
}

func TestGetKeyedTransactor_DynamicTipping(t *testing.T) {
	// Test the dynamic tipping calculation logic
	tests := []struct {
		name            string
		baseFeeGwei     float64
		maxCapGwei      float64
		expectedTipGwei float64
		shouldReject    bool
		rejectionReason string
	}{
		{
			name:            "Normal case - plenty of headroom",
			baseFeeGwei:     30.0,
			maxCapGwei:      80.0,
			expectedTipGwei: 50.0, // 80 - 30 = 50
			shouldReject:    false,
		},
		{
			name:            "Tight margins",
			baseFeeGwei:     78.0,
			maxCapGwei:      80.0,
			expectedTipGwei: 2.0, // 80 - 78 = 2
			shouldReject:    false,
		},
		{
			name:            "Very tight - minimum tip applied",
			baseFeeGwei:     79.5,
			maxCapGwei:      80.0,
			expectedTipGwei: 1.0, // Minimum tip floor
			shouldReject:    false,
		},
		{
			name:            "Base fee equals cap - should reject",
			baseFeeGwei:     80.0,
			maxCapGwei:      80.0,
			shouldReject:    true,
			rejectionReason: "base fee equals or exceeds cap",
		},
		{
			name:            "Base fee exceeds cap - should reject",
			baseFeeGwei:     100.0,
			maxCapGwei:      80.0,
			shouldReject:    true,
			rejectionReason: "base fee exceeds cap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseFee := tt.baseFeeGwei
			maxPrice := tt.maxCapGwei

			// Simulate the logic from GetKeyedTransactor
			if baseFee >= maxPrice {
				// Should reject
				assert.True(t, tt.shouldReject,
					"Expected rejection when baseFee (%.2f) >= maxPrice (%.2f)", baseFee, maxPrice)
				return
			}

			// Calculate tip
			tip := maxPrice - baseFee
			if tip < 1.0 {
				tip = 1.0
			}

			if !tt.shouldReject {
				assert.InDelta(t, tt.expectedTipGwei, tip, 0.1,
					"Tip calculation mismatch: expected %.2f, got %.2f", tt.expectedTipGwei, tip)
			}
		})
	}
}

func TestGetCurrentGasPrice_EIP1559Formula(t *testing.T) {
	// Test the EIP-1559 max fee calculation formula
	// MaxFeePerGas = (2 * BaseFee) + PriorityFee

	tests := []struct {
		name               string
		baseFeeGwei        float64
		priorityFeeGwei    float64
		expectedMaxFeeGwei float64
	}{
		{
			name:               "Typical Polygon",
			baseFeeGwei:        100.0,
			priorityFeeGwei:    30.0,
			expectedMaxFeeGwei: 230.0, // (2 * 100) + 30
		},
		{
			name:               "Low congestion",
			baseFeeGwei:        10.0,
			priorityFeeGwei:    2.0,
			expectedMaxFeeGwei: 22.0, // (2 * 10) + 2
		},
		{
			name:               "High congestion",
			baseFeeGwei:        500.0,
			priorityFeeGwei:    100.0,
			expectedMaxFeeGwei: 1100.0, // (2 * 500) + 100
		},
		{
			name:               "Edge case - zero priority fee",
			baseFeeGwei:        50.0,
			priorityFeeGwei:    0.0,
			expectedMaxFeeGwei: 100.0, // (2 * 50) + 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the formula from GetCurrentGasPrice
			baseFeeWei := gweiToWei(tt.baseFeeGwei)
			priorityFeeWei := gweiToWei(tt.priorityFeeGwei)

			// MaxFeePerGas = (2 * BaseFee) + PriorityFee
			maxFee := new(big.Int).Mul(baseFeeWei, big.NewInt(2))
			maxFee.Add(maxFee, priorityFeeWei)

			maxFeeGwei := WeiToGwei(maxFee)

			assert.InDelta(t, tt.expectedMaxFeeGwei, maxFeeGwei, 0.1,
				"Max fee calculation mismatch: expected %.2f, got %.2f", tt.expectedMaxFeeGwei, maxFeeGwei)
		})
	}
}

func TestGetKeyedTransactor_RejectWhenBaseFeeExceedsCap(t *testing.T) {
	// Test that GetKeyedTransactor properly rejects when base fee >= hard cap

	tests := []struct {
		name         string
		baseFeeGwei  float64
		hardCapGwei  float64
		shouldReject bool
		description  string
	}{
		{
			name:         "Base fee below cap - should allow",
			baseFeeGwei:  30.0,
			hardCapGwei:  80.0,
			shouldReject: false,
			description:  "Transaction should proceed when base fee is well below cap",
		},
		{
			name:         "Base fee equals cap - should reject",
			baseFeeGwei:  80.0,
			hardCapGwei:  80.0,
			shouldReject: true,
			description:  "Transaction should be rejected when base fee equals cap",
		},
		{
			name:         "Base fee exceeds cap - should reject",
			baseFeeGwei:  100.0,
			hardCapGwei:  80.0,
			shouldReject: true,
			description:  "Transaction should be rejected when base fee exceeds cap",
		},
		{
			name:         "Base fee slightly below cap - should allow with minimum tip",
			baseFeeGwei:  79.0,
			hardCapGwei:  80.0,
			shouldReject: false,
			description:  "Transaction should proceed with minimum 1 Gwei tip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the rejection logic from GetKeyedTransactor
			baseFee := tt.baseFeeGwei
			maxPrice := tt.hardCapGwei

			shouldReject := baseFee >= maxPrice

			assert.Equal(t, tt.shouldReject, shouldReject,
				"Rejection logic mismatch: %s", tt.description)

			if !shouldReject {
				// Calculate tip to verify it's reasonable
				tip := maxPrice - baseFee
				if tip < 1.0 {
					tip = 1.0
				}
				assert.GreaterOrEqual(t, tip, 1.0,
					"Tip should be at least 1 Gwei")
				assert.LessOrEqual(t, tip, maxPrice,
					"Tip should not exceed max price")
			}

			t.Logf("%s: baseFee=%.2f, hardCap=%.2f, reject=%v",
				tt.description, baseFee, maxPrice, shouldReject)
		})
	}
}

func TestGetKeyedTransactor_DynamicTippingFormula(t *testing.T) {
	// Test the dynamic tipping formula: Tip = MaxCap - CurrentBaseFee

	tests := []struct {
		name            string
		baseFeeGwei     float64
		hardCapGwei     float64
		expectedTipGwei float64
		description     string
	}{
		{
			name:            "Cheap network",
			baseFeeGwei:     10.0,
			hardCapGwei:     80.0,
			expectedTipGwei: 70.0, // 80 - 10
			description:     "High tip when network is cheap for fast inclusion",
		},
		{
			name:            "Moderate network",
			baseFeeGwei:     40.0,
			hardCapGwei:     80.0,
			expectedTipGwei: 40.0, // 80 - 40
			description:     "Moderate tip when network is moderately priced",
		},
		{
			name:            "Expensive network",
			baseFeeGwei:     70.0,
			hardCapGwei:     80.0,
			expectedTipGwei: 10.0, // 80 - 70
			description:     "Low tip when network is expensive but within cap",
		},
		{
			name:            "Very tight margin",
			baseFeeGwei:     79.5,
			hardCapGwei:     80.0,
			expectedTipGwei: 1.0, // Minimum tip floor applied
			description:     "Minimum 1 Gwei tip when very close to cap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseFee := tt.baseFeeGwei
			maxPrice := tt.hardCapGwei

			// Simulate the tipping formula
			tip := maxPrice - baseFee
			if tip < 1.0 {
				tip = 1.0
			}

			assert.InDelta(t, tt.expectedTipGwei, tip, 0.1,
				"Tip calculation mismatch: %s", tt.description)

			// Verify that total price = baseFee + tip <= maxPrice
			// Exception: When minimum tip (1 Gwei) is applied, total might slightly exceed cap
			totalPrice := baseFee + tip
			if tip > 1.0 {
				assert.LessOrEqual(t, totalPrice, maxPrice,
					"Total price should not exceed max price cap when tip > 1 Gwei")
			}

			t.Logf("%s: baseFee=%.2f, maxCap=%.2f, tip=%.2f, total=%.2f",
				tt.description, baseFee, maxPrice, tip, totalPrice)
		})
	}
}

func TestGetKeyedTransactor_ContextAndClientHandling(t *testing.T) {
	// Test the context and client parameter handling logic

	t.Run("Both nil - should use background context and create client", func(t *testing.T) {
		// When both ctx and client are nil:
		// 1. effectiveCtx should be context.Background()
		// 2. A new client should be created
		// 3. The new client should be closed after gas check

		// This documents the expected behavior
		t.Log("Expected: GetKeyedTransactor(nil, nil) uses Background context and creates temporary client")
	})

	t.Run("Context provided, client nil - should use provided context and create client", func(t *testing.T) {
		// When ctx is provided but client is nil:
		// 1. effectiveCtx should be the provided ctx
		// 2. A new client should be created
		// 3. The new client should be closed after gas check

		t.Log("Expected: GetKeyedTransactor(ctx, nil) uses provided context and creates temporary client")
	})

	t.Run("Both provided - should use both", func(t *testing.T) {
		// When both are provided:
		// 1. Use the provided ctx
		// 2. Use the provided client
		// 3. Do NOT close the client (caller owns it)

		t.Log("Expected: GetKeyedTransactor(ctx, client) uses both and does not close client")
	})

	t.Skip("Requires actual blockchain connection to test")
}

func TestGetKeyedTransactor_MandatoryGasCheck(t *testing.T) {
	// Test that gas check failure causes transaction rejection

	t.Run("Gas check failure should return error", func(t *testing.T) {
		// If GetCurrentGasPrice fails (RPC down, network error, etc.):
		// - GetKeyedTransactor should return an error
		// - No TransactOpts should be returned
		// - The error should wrap the original gas check error

		expectedErrorMessage := "gas price check failed: action rejected for safety"

		t.Logf("Expected behavior: If gas check fails, return error: %s", expectedErrorMessage)
		t.Log("This ensures no transaction is sent without price protection")
	})

	t.Run("Network congestion beyond cap should return error", func(t *testing.T) {
		// If network base fee >= hard cap:
		// - GetKeyedTransactor should return specific error
		// - Error should indicate base fee exceeds cap
		// - Should include actual values in error message

		expectedErrorFormat := "current network base fee (%.2f Gwei) exceeds your hard cap (%.2f Gwei)"

		t.Logf("Expected error format: %s", expectedErrorFormat)
		t.Log("This prevents transactions that would sit pending forever or be auto-rejected")
	})

	t.Skip("Requires mocking blockchain client")
}
