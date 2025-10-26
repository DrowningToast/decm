-- name: CreateEventIssuer :one
INSERT INTO event_issuers (
    event_id,
    issuer_credential_id,
    is_signed,
    signature,
    sign_message
) VALUES (
    sqlc.arg('event_id'),
    sqlc.arg('issuer_credential_id'),
    sqlc.arg('is_signed'),
    sqlc.arg('signature'),
    sqlc.arg('sign_message')
) RETURNING *;

-- name: GetEventIssuersByEventID :many
SELECT * FROM event_issuers WHERE event_id = sqlc.arg('event_id');

-- name: GetEventIssuerByID :one
SELECT * FROM event_issuers WHERE id = sqlc.arg('id');

-- name: UpdateEventIssuer :one
UPDATE event_issuers
SET 
    is_signed = sqlc.arg('is_signed'),
    signature = sqlc.arg('signature'),
    sign_message = sqlc.arg('sign_message'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEventIssuer :exec
DELETE FROM event_issuers WHERE id = sqlc.arg('id');

-- name: GetEventIssuerByEventIDAndIssuerCredentialID :one
SELECT * FROM event_issuers 
WHERE event_id = sqlc.arg('event_id') AND issuer_credential_id = sqlc.arg('issuer_credential_id');
