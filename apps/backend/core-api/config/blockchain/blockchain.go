package blockchain

import "errors"

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
	// GasLimit is the maximum gas to use for transactions
	// 0 means automatic estimation (unreliable on some networks)
	// Recommended: 3000000 for contract deployments, 500000 for regular txs
	GasLimit uint64 `env:"GAS_LIMIT" envDefault:"3000000"`
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
