package blockchain

import (
	"errors"
	"time"
)

type BlockchainConfig struct {
	Network    string `env:"NETWORK" envDefault:"localhost"`
	RPCURL     string `env:"RPC_URL" envDefault:"http://localhost:8545"`
	PrivateKey string `env:"PRIVATE_KEY" envDefault:""`
	ChainID    int    `env:"CHAIN_ID" envDefault:"1337"`
	// DecmAccessManagerAddress is the address of the DecmAccessManager contract
	DecmAccessManagerAddress string `env:"DECM_ACCESS_MANAGER_ADDRESS" envDefault:""`
	// Etherscan API Key
	EtherscanAPIKey string `env:"ETHERSCAN_API_KEY" envDefault:""`

	// Gas Management (Optional)
	// MaxGasPriceGwei limits the total fee per gas (Base Fee + Priority Fee)
	// On prod max out at 80 gwei
	MaxGasPriceGwei float64 `env:"MAX_GAS_PRICE_GWEI" envDefault:"80"`
	// SoftCapGasPriceGwei is a softer limit for UI warnings
	// On prod set to 60 gwei
	SoftCapGasPriceGwei float64 `env:"SOFT_CAP_GAS_PRICE_GWEI" envDefault:"60"`

	// PollInterval is the duration between worker runs
	// Default: 1h
	PollInterval time.Duration `env:"POLL_INTERVAL" envDefault:"1h"`
}

func (c *BlockchainConfig) Validate() error {
	if c.Network == "" {
		return errors.New("network is required")
	}
	if c.RPCURL == "" {
		return errors.New("RPC URL is required")
	}
	if c.PrivateKey == "" {
		return errors.New("private key is required")
	}
	if c.ChainID == 0 {
		return errors.New("chain ID is required")
	}
	if c.EtherscanAPIKey == "" {
		return errors.New("etherscan API key is required")
	}
	return nil
}

// GetBlockTimeSeconds returns the average block time in seconds based on the chain ID
// This is used for deadline calculations and time estimations
func (c *BlockchainConfig) GetBlockTimeSeconds() int {
	switch c.ChainID {
	case 1: // Ethereum Mainnet
		return 12
	case 137: // Polygon PoS Mainnet
		return 2
	case 80001, 80002: // Polygon Mumbai Testnet / Amoy Testnet
		return 2
	case 56: // BSC Mainnet
		return 3
	case 97: // BSC Testnet
		return 3
	case 11155111: // Sepolia Testnet
		return 12
	case 5: // Goerli Testnet (deprecated)
		return 12
	case 1337, 31337: // Localhost (Hardhat/Ganache)
		return 2
	default:
		// Default to 2 seconds for unknown chains (conservative estimate)
		return 2
	}
}
