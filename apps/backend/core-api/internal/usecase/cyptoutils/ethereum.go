package cyptoutils

import (
	"apps/backend/core-api/config"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func GetEthereumClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.LoadConfig().Blockchain.RPCURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ethereum client")
	}

	return client, nil
}

func GetKeyedTransactor() (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(config.LoadConfig().Blockchain.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create private key")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(config.LoadConfig().Blockchain.ChainID)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create keyed transactor")
	}

	return auth, nil
}
