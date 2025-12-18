-- name: CreateSystemStatusSchedule :one
INSERT INTO system_status_schedules (
    start_time,
    planned_end_time,
    status,
    is_planned
) VALUES (
    sqlc.arg(start_time),
    sqlc.arg(planned_end_time),
    sqlc.arg(status),
    sqlc.arg(is_planned)
) RETURNING *;

-- name: GetSystemStatusScheduleById :one
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE id = sqlc.arg(id)
AND deleted_at IS NULL;

-- name: GetSystemStatusScheduleByOrderId :one
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE order_id = sqlc.arg(order_id)
AND deleted_at IS NULL;

-- name: GetCurrentSystemStatus :one
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE deleted_at IS NULL
AND start_time <= NOW()
ORDER BY start_time DESC
LIMIT 1;

-- name: GetUpcomingSystemStatusSchedules :many
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE deleted_at IS NULL
AND start_time > NOW()
ORDER BY start_time ASC
LIMIT sqlc.arg(limit_count)::INTEGER;

-- name: GetSystemStatusScheduleHistory :many
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE deleted_at IS NULL
ORDER BY start_time DESC
LIMIT sqlc.arg(limit_count)::INTEGER
OFFSET sqlc.arg(offset_count)::INTEGER;

-- name: UpdateSystemStatusSchedule :one
UPDATE system_status_schedules
SET
    start_time = COALESCE(sqlc.narg(start_time), start_time),
    planned_end_time = COALESCE(sqlc.narg(planned_end_time), planned_end_time),
    status = COALESCE(sqlc.narg(status), status),
    is_planned = COALESCE(sqlc.narg(is_planned), is_planned),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
AND deleted_at IS NULL
RETURNING *;

-- name: DeleteSystemStatusSchedule :exec
UPDATE system_status_schedules
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
AND deleted_at IS NULL;

-- name: GetPlannedMaintenanceSchedules :many
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE deleted_at IS NULL
AND is_planned = TRUE
AND status = 0
AND start_time > NOW()
ORDER BY start_time ASC;

-- name: CountSystemStatusSchedules :one
SELECT COUNT(*) FROM system_status_schedules
WHERE deleted_at IS NULL;

-- name: GetSystemStatusSchedulesUpdatedBetween :many
SELECT 
    id,
    order_id,
    start_time,
    planned_end_time,
    status,
    is_planned,
    created_at,
    updated_at,
    deleted_at
FROM system_status_schedules
WHERE deleted_at IS NULL
AND updated_at >= sqlc.arg(start_date)
AND updated_at <= sqlc.arg(end_date)
ORDER BY updated_at DESC;
