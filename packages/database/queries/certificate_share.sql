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
