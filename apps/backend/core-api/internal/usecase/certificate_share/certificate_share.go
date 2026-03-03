package certificate_share

import (
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	certificatecontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/certificate_contract"
)

type CertificateShareUsecase struct {
	EventCertificateDataGateway eventdatagateway.EventCertificateDataGateway
	CertificateShareDg          eventdatagateway.CertificateShareDataGateway
	CertificateContractFactoryDg certificatecontract_datagateway.CertificateContractFactoryDataGateway
}

func NewCertificateShareUsecase(
	eventCertificateDataGateway eventdatagateway.EventCertificateDataGateway,
	certificateShareDg eventdatagateway.CertificateShareDataGateway,
	certificateContractFactoryDg certificatecontract_datagateway.CertificateContractFactoryDataGateway,
) *CertificateShareUsecase {
	return &CertificateShareUsecase{
		EventCertificateDataGateway:  eventCertificateDataGateway,
		CertificateShareDg:           certificateShareDg,
		CertificateContractFactoryDg: certificateContractFactoryDg,
	}
}
