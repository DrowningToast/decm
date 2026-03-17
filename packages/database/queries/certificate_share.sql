-- name: CreateCertificateShare :one
INSERT INTO event_certificate_shares (event_certificate_id, active, handle, password)
VALUES (
    sqlc.arg('event_certificate_id'),
    sqlc.arg('active'),
    sqlc.arg('handle'),
    sqlc.arg('password')
)
RETURNING *;

-- name: GetCertificateShareByHandle :one
SELECT * FROM event_certificate_shares WHERE handle = sqlc.arg('handle') LIMIT 1;

-- name: GetCertificateShareByID :one
SELECT * FROM event_certificate_shares WHERE id = sqlc.arg('id') LIMIT 1;

-- name: GetCertificateShareByEventCertificateID :one
SELECT * FROM event_certificate_shares WHERE event_certificate_id = sqlc.arg('event_certificate_id') LIMIT 1;

-- name: UpdateCertificateShare :one
UPDATE event_certificate_shares
SET password = sqlc.arg('password'),
    active = sqlc.arg('active'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;
