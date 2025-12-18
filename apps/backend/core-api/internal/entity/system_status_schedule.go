package entity

import "time"

type SystemStatus int

const (
	SystemStatusMaintenance SystemStatus = 0
	SystemStatusOperating   SystemStatus = 1
)

type SystemStatusSchedule struct {
	ID             int32        `json:"id"`
	OrderId        int32        `json:"orderId"`
	StartTime      time.Time    `json:"startTime"`
	PlannedEndTime *time.Time   `json:"plannedEndTime"`
	Status         SystemStatus `json:"status"`
	IsPlanned      bool         `json:"isPlanned"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	DeletedAt      *time.Time   `json:"deletedAt"`
}
