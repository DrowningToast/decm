package contract

import (
	"apps/backend/common/customerror"
	certificateContract "apps/backend/contracts/certificate"
	blockchainclient_datagateway "apps/backend/core-api/internal/datagateway/onchain/blockchain_client"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CertificateContractFactoryRepository struct {
	client             *ethclient.Client
	blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway
}

func NewCertificateContractFactoryRepository(client *ethclient.Client, blockchainClientDg blockchainclient_datagateway.BlockchainClientDataGateway) *CertificateContractFactoryRepository {
	return &CertificateContractFactoryRepository{
		client:             client,
		blockchainClientDg: blockchainClientDg,
	}
}

var _ certificatecontract_datagateway.CertificateContractFactoryDataGateway = (*CertificateContractFactoryRepository)(nil)

func (r *CertificateContractFactoryRepository) GetContract(address common.Address) (certificatecontract_datagateway.CertificateContractDataGateway, error) {
	instance, err := certificateContract.NewEventCertificate(address, r.client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return NewCertificateContractRepository(r.client, instance, r.blockchainClientDg), nil
}
