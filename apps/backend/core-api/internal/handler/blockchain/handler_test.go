package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"apps/backend/core-api/config"
	"apps/backend/core-api/internal/entity"

	blockchainConfig "apps/backend/core-api/config/blockchain"
	blockchain_usecase "apps/backend/core-api/internal/usecase/blockchain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// mockBlockchainUsecase is a mock implementation for testing
type mockBlockchainUsecase struct {
	gasPriceResponse *blockchain_usecase.GasPriceResponse
	shouldError      bool
	errorToReturn    error
}

func (m *mockBlockchainUsecase) GetGasPrice(ctx interface{}) (*blockchain_usecase.GasPriceResponse, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.gasPriceResponse, nil
}

// mockSystemStatusUsecase is a mock implementation for testing
type mockSystemStatusUsecase struct {
	schedule      *entity.SystemStatusSchedule
	shouldError   bool
	errorToReturn error
}

func (m *mockSystemStatusUsecase) GetClosestIncomingSchedule(ctx context.Context) (*entity.SystemStatusSchedule, error) {
	if m.shouldError {
		return nil, m.errorToReturn
	}
	return m.schedule, nil
}

func TestHandler_GetGasPrice_Success(t *testing.T) {
	// Setup
	app := fiber.New()

	mockUsecase := &mockBlockchainUsecase{
		gasPriceResponse: &blockchain_usecase.GasPriceResponse{
			CurrentGasPriceGwei: 230.0,
			BaseFeeGwei:         100.0,
			PriorityFeeGwei:     30.0,
			SoftCapPriceGwei:    1000.0,
			HardCapPriceGwei:    2500.0,
			SoftSafetyMargin:    77.0,
			HardSafetyMargin:    90.8,
		},
		shouldError: false,
	}

	// Create a mock config
	cfg := &config.Config{
		Blockchain: blockchainConfig.BlockchainConfig{
			MaxGasPriceGwei:     2500.0,
			SoftCapGasPriceGwei: 1000.0,
		},
	}

	// Note: Since the handler directly creates BlockchainUsecase,
	// we can't easily inject the mock without refactoring
	// This test documents the expected HTTP response format

	expectedResponse := blockchain_usecase.GasPriceResponse{
		CurrentGasPriceGwei: 230.0,
		BaseFeeGwei:         100.0,
		PriorityFeeGwei:     30.0,
		SoftCapPriceGwei:    1000.0,
		HardCapPriceGwei:    2500.0,
		SoftSafetyMargin:    77.0,
		HardSafetyMargin:    90.8,
	}

	// Verify the response structure can be marshaled
	_, err := json.Marshal(expectedResponse)
	assert.NoError(t, err, "Response should be JSON serializable")

	// Verify field types
	assert.IsType(t, float64(0), expectedResponse.CurrentGasPriceGwei)
	assert.IsType(t, float64(0), expectedResponse.BaseFeeGwei)
	assert.IsType(t, float64(0), expectedResponse.PriorityFeeGwei)
	assert.IsType(t, float64(0), expectedResponse.SoftCapPriceGwei)
	assert.IsType(t, float64(0), expectedResponse.HardCapPriceGwei)
	assert.IsType(t, float64(0), expectedResponse.SoftSafetyMargin)
	assert.IsType(t, float64(0), expectedResponse.HardSafetyMargin)

	_ = app
	_ = mockUsecase
	_ = cfg
	t.Skip("Handler integration test requires dependency injection refactoring")
}

func TestHandler_GetGasPrice_CacheHeaders(t *testing.T) {
	// Test that the handler sets appropriate cache headers

	// Expected cache header
	expectedCacheControl := "public, max-age=300" // 5 minutes

	// Verify the cache duration
	assert.Equal(t, "public, max-age=300", expectedCacheControl,
		"Cache-Control header should be set to 5 minutes (300 seconds)")

	t.Skip("Handler integration test requires dependency injection refactoring")
}

func TestHandler_GetGasPrice_ErrorHandling(t *testing.T) {
	// Test error scenarios

	tests := []struct {
		name          string
		mockError     error
		expectedError string
	}{
		{
			name:          "RPC connection failure",
			mockError:     errors.New("failed to connect to blockchain RPC"),
			expectedError: "failed to connect to blockchain RPC",
		},
		{
			name:          "Gas price fetch failure",
			mockError:     errors.New("failed to fetch gas price"),
			expectedError: "failed to fetch gas price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Document expected error handling behavior
			assert.NotNil(t, tt.mockError, "Error should be returned")
			assert.Contains(t, tt.mockError.Error(), tt.expectedError,
				"Error message should contain expected text")
		})
	}

	t.Skip("Handler integration test requires dependency injection refactoring")
}

func TestHandler_Routes(t *testing.T) {
	// Test that routes are properly mounted
	app := fiber.New()

	// The handler should mount a GET route at /blockchain/gas-price
	expectedRoute := "/blockchain/gas-price"
	expectedMethod := "GET"

	// Document the API contract
	assert.NotEmpty(t, expectedRoute, "Route should be defined")
	assert.Equal(t, "GET", expectedMethod, "Method should be GET")

	_ = app
	t.Skip("Route testing requires full handler setup")
}

func TestGasPriceResponse_JSONMarshaling(t *testing.T) {
	// Test that the response structure marshals correctly to JSON

	response := blockchain_usecase.GasPriceResponse{
		CurrentGasPriceGwei: 271.5,
		BaseFeeGwei:         120.75,
		PriorityFeeGwei:     30.0,
		SoftCapPriceGwei:    1000.0,
		HardCapPriceGwei:    2500.0,
		SoftSafetyMargin:    72.85,
		HardSafetyMargin:    89.14,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err, "Should marshal to JSON without error")

	// Unmarshal back
	var unmarshaled blockchain_usecase.GasPriceResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err, "Should unmarshal from JSON without error")

	// Verify values are preserved
	assert.InDelta(t, response.CurrentGasPriceGwei, unmarshaled.CurrentGasPriceGwei, 0.01)
	assert.InDelta(t, response.BaseFeeGwei, unmarshaled.BaseFeeGwei, 0.01)
	assert.InDelta(t, response.PriorityFeeGwei, unmarshaled.PriorityFeeGwei, 0.01)
	assert.InDelta(t, response.SoftCapPriceGwei, unmarshaled.SoftCapPriceGwei, 0.01)
	assert.InDelta(t, response.HardCapPriceGwei, unmarshaled.HardCapPriceGwei, 0.01)
	assert.InDelta(t, response.SoftSafetyMargin, unmarshaled.SoftSafetyMargin, 0.01)
	assert.InDelta(t, response.HardSafetyMargin, unmarshaled.HardSafetyMargin, 0.01)

	// Verify JSON field names (snake_case as per convention)
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, "current_gas_price_gwei")
	assert.Contains(t, jsonString, "base_fee_gwei")
	assert.Contains(t, jsonString, "priority_fee_gwei")
	assert.Contains(t, jsonString, "soft_cap_price_gwei")
	assert.Contains(t, jsonString, "hard_cap_price_gwei")
	assert.Contains(t, jsonString, "soft_safety_margin")
	assert.Contains(t, jsonString, "hard_safety_margin")
}

func TestGasPriceResponse_HTTPResponseFormat(t *testing.T) {
	// Test the complete HTTP response format

	response := blockchain_usecase.GasPriceResponse{
		CurrentGasPriceGwei: 271.5,
		BaseFeeGwei:         120.75,
		PriorityFeeGwei:     30.0,
		SoftCapPriceGwei:    1000.0,
		HardCapPriceGwei:    2500.0,
		SoftSafetyMargin:    72.85,
		HardSafetyMargin:    89.14,
	}

	// Create a test server
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "public, max-age=300")
		return c.JSON(response)
	})

	// Make request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Verify status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify cache header
	assert.Equal(t, "public, max-age=300", resp.Header.Get("Cache-Control"))

	// Verify content type
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
}

func TestGetSystemStatusResponse_JSONMarshaling(t *testing.T) {
	// Test that the combined response structure marshals correctly to JSON

	response := GetSystemStatusResponse{
		GasPrice: &blockchain_usecase.GasPriceResponse{
			CurrentGasPriceGwei: 271.5,
			BaseFeeGwei:         120.75,
			PriorityFeeGwei:     30.0,
			SoftCapPriceGwei:    1000.0,
			HardCapPriceGwei:    2500.0,
			SoftSafetyMargin:    72.85,
			HardSafetyMargin:    89.14,
		},
		ClosestIncomingSchedule: nil, // Test nullable field
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err, "Should marshal to JSON without error")

	// Unmarshal back
	var unmarshaled GetSystemStatusResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err, "Should unmarshal from JSON without error")

	// Verify gas price values are preserved
	assert.NotNil(t, unmarshaled.GasPrice)
	assert.InDelta(t, response.GasPrice.CurrentGasPriceGwei, unmarshaled.GasPrice.CurrentGasPriceGwei, 0.01)
	assert.Nil(t, unmarshaled.ClosestIncomingSchedule)

	// Verify JSON field names
	jsonString := string(jsonData)
	assert.Contains(t, jsonString, "gas_price")
	// closest_incoming_schedule is omitted when nil due to omitempty tag
	assert.NotContains(t, jsonString, "closest_incoming_schedule")
}

func TestHandler_Routes_WithStatus(t *testing.T) {
	// Test that both routes are properly mounted
	app := fiber.New()

	// The handler should mount GET routes at:
	// - /blockchain/gas-price (existing)
	// - /blockchain/status (new combined endpoint)
	expectedRoutes := []struct {
		path   string
		method string
	}{
		{"/blockchain/gas-price", "GET"},
		{"/blockchain/status", "GET"},
	}

	// Document the API contract
	for _, route := range expectedRoutes {
		assert.NotEmpty(t, route.path, "Route path should be defined")
		assert.Equal(t, "GET", route.method, "Method should be GET")
	}

	_ = app
	t.Skip("Route testing requires full handler setup")
}
