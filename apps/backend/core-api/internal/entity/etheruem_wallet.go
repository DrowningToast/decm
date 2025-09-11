package entity

type EthereumWallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Mnemonic   string `json:"mnemonic,omitempty"`
}
