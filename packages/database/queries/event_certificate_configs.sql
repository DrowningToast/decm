-- name: CreateEventCertificateConfig :one
INSERT INTO event_certificate_configs (
    event_id,
    base_certificate_storage_key,
    event_name_pos_x,
    event_name_pos_y,
    name_pos_x,
    name_pos_y,
    academic_institution_pos_x,
    academic_institution_pos_y
) VALUES (
    sqlc.arg('event_id'),
    sqlc.arg('base_certificate_storage_key'),
    sqlc.arg('event_name_pos_x'),
    sqlc.arg('event_name_pos_y'),
    sqlc.arg('name_pos_x'),
    sqlc.arg('name_pos_y'),
    sqlc.arg('academic_institution_pos_x'),
    sqlc.arg('academic_institution_pos_y')
) RETURNING *;

-- name: GetEventCertificateConfigByEventID :one
SELECT * FROM event_certificate_configs WHERE event_id = sqlc.arg('event_id');

-- name: UpdateEventCertificateConfig :one
UPDATE event_certificate_configs
SET 
    base_certificate_storage_key = sqlc.arg('base_certificate_storage_key'),
    event_name_pos_x = sqlc.arg('event_name_pos_x'),
    event_name_pos_y = sqlc.arg('event_name_pos_y'),
    name_pos_x = sqlc.arg('name_pos_x'),
    name_pos_y = sqlc.arg('name_pos_y'),
    academic_institution_pos_x = sqlc.arg('academic_institution_pos_x'),
    academic_institution_pos_y = sqlc.arg('academic_institution_pos_y'),
    updated_at = NOW()
WHERE event_id = sqlc.arg('event_id')
RETURNING *;

-- name: DeleteEventCertificateConfig :exec
DELETE FROM event_certificate_configs WHERE event_id = sqlc.arg('event_id');
