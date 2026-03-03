package certificate_share

import (
	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"
)

type CertificateShareUsecase struct {
	EventCertificateDataGateway eventdatagateway.EventCertificateDataGateway
	CertificateShareDg          eventdatagateway.CertificateShareDataGateway
}

func NewCertificateShareUsecase(
	eventCertificateDataGateway eventdatagateway.EventCertificateDataGateway,
	certificateShareDg eventdatagateway.CertificateShareDataGateway,
) *CertificateShareUsecase {
	return &CertificateShareUsecase{
		EventCertificateDataGateway: eventCertificateDataGateway,
		CertificateShareDg:          certificateShareDg,
	}
}
