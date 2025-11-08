-- name: CreateEventCertificateSignature :one
INSERT INTO event_certificate_signatures (
    event_certificate_id,
    issuer_credential_id,
    issuer_signature,
    host_signature,
    sign_message,
    sign_message_digest
) VALUES (
    sqlc.arg('event_certificate_id'),
    sqlc.arg('issuer_credential_id'),
    sqlc.arg('issuer_signature'),
    sqlc.arg('host_signature'),
    sqlc.arg('sign_message'),
    sqlc.arg('sign_message_digest')
) RETURNING *;

-- name: GetEventCertificateSignatureByID :one
SELECT * FROM event_certificate_signatures WHERE id = sqlc.arg('id');

-- name: GetEventCertificateSignaturesByEventCertificateID :many
SELECT * FROM event_certificate_signatures WHERE event_certificate_id = sqlc.arg('event_certificate_id');

-- name: UpdateEventCertificateSignature :one
UPDATE event_certificate_signatures
SET 
    issuer_signature = sqlc.arg('issuer_signature'),
    host_signature = sqlc.arg('host_signature'),
    sign_message = sqlc.arg('sign_message'),
    sign_message_digest = sqlc.arg('sign_message_digest'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEventCertificateSignature :exec
DELETE FROM event_certificate_signatures WHERE id = sqlc.arg('id');
