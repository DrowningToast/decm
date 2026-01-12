package cyptoutils

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetLatestBlockNumber(client *ethclient.Client) (uint64, error) {
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func GetCalculatedDeadlineBlock(client *ethclient.Client) (uint64, error) {
	latestBlockNumber, err := GetLatestBlockNumber(client)
	if err != nil {
		return 0, err
	}

	// Polygon PoS (Proof of Stake) mainnet is approximately 2 seconds
	// https://polygonscan.com/chart/blocktime
	// 302,400 blocks = 1 week (604,800 seconds)
	deadlineBlock := 302400

	return latestBlockNumber + uint64(deadlineBlock), nil
}
