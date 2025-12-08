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
WHERE event_id = sqlc.arg('event_id') 
  AND issuer_credential_id = sqlc.arg('issuer_credential_id')
  AND deleted_at IS NULL;

-- name: GetEventIssuersByCredentialID :many
SELECT 
    ei.id,
    ei.event_id as event_id,
    ei.issuer_credential_id,
    ei.is_signed,
    ei.signature,
    ei.sign_message_digest,
    ei.created_at,
    ei.updated_at,
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

-- name: ResetAllEventIssuersSigningStatus :exec
UPDATE event_issuers 
SET is_signed = 0, updated_at = NOW()
WHERE event_id = sqlc.arg('event_id');

-- name: GetIssuerEventsWithDetails :many
SELECT 
    ei.id,
    ei.event_id as event_id,
    ei.issuer_credential_id,
    ei.is_signed,
    ei.signature,
    ei.sign_message_digest,
    ei.created_at,
    ei.updated_at,
    e.title as event_title,
    e.short_description as event_short_description,
    e.start_date as event_start_date,
    e.end_date as event_end_date,
    e.location as event_location,
    e.owner_credential_id as event_owner_credential_id,
    p.first_name as owner_first_name,
    p.last_name as owner_last_name,
    p.email as owner_email,
    ac.wallet_address as owner_wallet_address,
    ac.google_connector_ref as owner_google_connector_ref,
    COALESCE(
        (SELECT COUNT(ec.id) 
         FROM event_certificates ec 
         WHERE ec.event_id = e.id 
           AND ec.revoked_at IS NULL
        ), 
    0)::INTEGER AS certificate_count
FROM event_issuers ei
INNER JOIN events e ON ei.event_id = e.id
LEFT JOIN profiles p ON e.owner_credential_id = p.authentication_credential_id
LEFT JOIN authentication_credentials ac ON e.owner_credential_id = ac.id
WHERE ei.issuer_credential_id = sqlc.arg('issuer_credential_id')
ORDER BY e.created_at DESC
LIMIT sqlc.arg('limit_count') OFFSET sqlc.arg('offset_count');

-- name: HasSignedIssuers :one
SELECT EXISTS(
    SELECT 1 
    FROM event_issuers 
    WHERE event_id = sqlc.arg('event_id') 
      AND is_signed = 1
      AND deleted_at IS NULL
) AS has_signed_issuers;

-- name: GetSignedIssuersCount :one
SELECT COUNT(*) AS count
FROM event_issuers 
WHERE event_id = sqlc.arg('event_id') 
  AND is_signed = 1
  AND deleted_at IS NULL;

-- name: GetTotalIssuersCount :one
SELECT COUNT(*) AS count
FROM event_issuers 
WHERE event_id = sqlc.arg('event_id') 
  AND deleted_at IS NULL;

-- name: AllIssuersHaveSigned :one
SELECT 
    CASE 
        WHEN COUNT(*) = 0 THEN false
        WHEN COUNT(*) = COUNT(*) FILTER (WHERE is_signed = 1) THEN true
        ELSE false
    END AS all_issuers_signed
FROM event_issuers
WHERE event_id = sqlc.arg('event_id')
  AND deleted_at IS NULL;
