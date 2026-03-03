package certificatecontract_datagateway

import "github.com/ethereum/go-ethereum/common"

type CertificateContractFactoryDataGateway interface {
	GetContract(address common.Address) (CertificateContractDataGateway, error)
}
