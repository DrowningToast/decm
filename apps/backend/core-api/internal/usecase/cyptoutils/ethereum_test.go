package cyptoutils

import (
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

	auth, err := GetKeyedTransactor()
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








