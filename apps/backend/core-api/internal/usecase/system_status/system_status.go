package system_status

import (
	"apps/backend/core-api/internal/datagateway"
)

type SystemStatusUsecase struct {
	systemStatusRepo datagateway.SystemStatusScheduleDataGateway
}

func NewSystemStatusUsecase(systemStatusRepo datagateway.SystemStatusScheduleDataGateway) *SystemStatusUsecase {
	return &SystemStatusUsecase{
		systemStatusRepo: systemStatusRepo,
	}
}
