-- name: CreateCertificateShare :one
INSERT INTO certificate_share (event_certificate_id, active, handle, password)
VALUES (
    sqlc.arg('event_certificate_id'),
    sqlc.arg('active'),
    sqlc.arg('handle'),
    sqlc.arg('password')
)
RETURNING *;

-- name: GetCertificateShareByHandle :one
SELECT * FROM certificate_share WHERE handle = sqlc.arg('handle') LIMIT 1;

-- name: GetCertificateShareByID :one
SELECT * FROM certificate_share WHERE id = sqlc.arg('id') LIMIT 1;

-- name: GetCertificateShareByEventCertificateID :one
SELECT * FROM certificate_share WHERE event_certificate_id = sqlc.arg('event_certificate_id') LIMIT 1;

-- name: UpdateCertificateShare :one
UPDATE certificate_share
SET password = sqlc.arg('password'),
    active = sqlc.arg('active'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;
