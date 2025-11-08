-- name: CreateEventIssuer :one
INSERT INTO event_issuers (
    event_id,
    issuer_credential_id,
    is_signed,
    signature,
    sign_message_digest
) VALUES (
    sqlc.arg('event_id'),
    sqlc.arg('issuer_credential_id'),
    sqlc.arg('is_signed'),
    sqlc.arg('signature'),
    sqlc.arg('sign_message_digest')
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
    sign_message_digest = sqlc.arg('sign_message_digest'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEventIssuer :exec
DELETE FROM event_issuers WHERE id = sqlc.arg('id');

-- name: GetEventIssuerByEventIDAndIssuerCredentialID :one
SELECT * FROM event_issuers 
WHERE event_id = sqlc.arg('event_id') AND issuer_credential_id = sqlc.arg('issuer_credential_id');

-- name: GetEventIssuersByCredentialID :many
SELECT 
    ei.id,
    ei.event_id,
    ei.issuer_credential_id,
    ei.is_signed,
    ei.signature,
    ei.sign_message_digest,
    ei.created_at,
    ei.updated_at,
    e.id as event_id,
    e.title as event_title,
    e.short_description as event_short_description,
    e.start_date as event_start_date,
    e.end_date as event_end_date,
    e.location as event_location,
    e.owner_credential_id as event_owner_credential_id
FROM event_issuers ei
INNER JOIN events e ON ei.event_id = e.id
WHERE ei.issuer_credential_id = sqlc.arg('issuer_credential_id')
ORDER BY e.created_at DESC
LIMIT sqlc.arg('limit_count') OFFSET sqlc.arg('offset_count');
