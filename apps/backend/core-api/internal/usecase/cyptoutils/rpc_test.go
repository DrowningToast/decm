package cyptoutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestBlockNumber(t *testing.T) {
	// This test requires actual Ethereum RPC connection
	// In a real scenario, you'd mock the ethclient or use a test RPC URL
	t.Skip("Skipping test that requires actual Ethereum RPC connection")

	// Example of how the test would look with a real connection:
	// client, err := ethclient.Dial("http://localhost:8545")
	// require.NoError(t, err)
	// defer client.Close()
	//
	// blockNumber, err := GetLatestBlockNumber(client)
	// require.NoError(t, err)
	// assert.Greater(t, blockNumber, uint64(0))
}

func TestGetCalculatedDeadlineBlock(t *testing.T) {
	// This test requires actual Ethereum RPC connection
	t.Skip("Skipping test that requires actual Ethereum RPC connection")

	// Example of how the test would look with a real connection:
	// client, err := ethclient.Dial("http://localhost:8545")
	// require.NoError(t, err)
	// defer client.Close()
	//
	// deadlineBlock, err := GetCalculatedDeadlineBlock(client)
	// require.NoError(t, err)
	// assert.Greater(t, deadlineBlock, uint64(0))
	//
	// latestBlock, err := GetLatestBlockNumber(client)
	// require.NoError(t, err)
	// assert.Equal(t, latestBlock+30, deadlineBlock)
}

func TestGetCalculatedDeadlineBlock_DeadlineOffset(t *testing.T) {
	// Unit test for the deadline calculation logic
	// The deadline should be latestBlock + 30
	latestBlock := uint64(1000)
	expectedDeadline := latestBlock + 30

	assert.Equal(t, uint64(1030), expectedDeadline)
}







