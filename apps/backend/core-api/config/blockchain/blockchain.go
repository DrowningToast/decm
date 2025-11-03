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
