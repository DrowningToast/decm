package system_status

import (
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
)

type SystemStatusUsecase struct {
	systemStatusRepo offchain_datagateway.SystemStatusScheduleDataGateway
}

func NewSystemStatusUsecase(systemStatusRepo offchain_datagateway.SystemStatusScheduleDataGateway) *SystemStatusUsecase {
	return &SystemStatusUsecase{
		systemStatusRepo: systemStatusRepo,
	}
}
