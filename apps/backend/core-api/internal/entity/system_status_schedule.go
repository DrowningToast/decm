package entity

import "time"

type SystemStatus int

const (
	SystemStatusMaintenance SystemStatus = 0
	SystemStatusOperating   SystemStatus = 1
)

type SystemStatusSchedule struct {
	ID             int32        `json:"id"`
	OrderId        int32        `json:"order_id"`
	StartTime      time.Time    `json:"start_time"`
	PlannedEndTime *time.Time   `json:"planned_end_time"`
	Status         SystemStatus `json:"status"`
	IsPlanned      bool         `json:"is_planned"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at"`
}
