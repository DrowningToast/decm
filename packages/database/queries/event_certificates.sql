-- name: CreateEventCertificate :one
INSERT INTO event_certificates (
    event_id,
    receiver_credential_id,
    receiver_email,
    name,
    academic_institution,
    certificate_title,
    certificate_subtitle,
    event_contract_address,
    event_certificate_address,
    certificate_token_id,
    certificate_digest
) VALUES (
    sqlc.arg('event_id'),
    sqlc.arg('receiver_credential_id'),
    sqlc.arg('receiver_email'),
    sqlc.arg('name'),
    sqlc.arg('academic_institution'),
    sqlc.arg('certificate_title'),
    sqlc.arg('certificate_subtitle'),
    sqlc.arg('event_contract_address'),
    sqlc.arg('event_certificate_address'),
    sqlc.arg('certificate_token_id'),
    sqlc.arg('certificate_digest')
) RETURNING *;

-- name: GetEventCertificateByID :one
SELECT * FROM event_certificates WHERE id = sqlc.arg('id');

-- name: GetEventCertificatesByEventID :many
SELECT * FROM event_certificates WHERE event_id = sqlc.arg('event_id') AND revoked_at IS NULL;

-- name: UpdateEventCertificate :one
UPDATE event_certificates
SET 
    receiver_credential_id = sqlc.arg('receiver_credential_id'),
    receiver_email = sqlc.arg('receiver_email'),
    name = sqlc.arg('name'),
    academic_institution = sqlc.arg('academic_institution'),
    certificate_title = sqlc.arg('certificate_title'),
    certificate_subtitle = sqlc.arg('certificate_subtitle'),
    event_contract_address = sqlc.arg('event_contract_address'),
    event_certificate_address = sqlc.arg('event_certificate_address'),
    certificate_token_id = sqlc.arg('certificate_token_id'),
    revoked_at = sqlc.arg('revoked_at')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEventCertificate :exec
DELETE FROM event_certificates WHERE id = sqlc.arg('id');

-- name: GetAllEventCertificateIDsByEventID :many
SELECT id FROM event_certificates WHERE event_id = sqlc.arg('event_id');

-- name: GetEventCertificateByInboxMessageID :one
SELECT * FROM event_certificates WHERE inbox_message_id = sqlc.arg('inbox_message_id');

-- name: UpdateEventCertificateInboxMessageID :one
UPDATE event_certificates
SET inbox_message_id = sqlc.arg('inbox_message_id')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: GetClaimedCertificatesByEventID :many
SELECT ec.* 
FROM event_certificates ec
WHERE ec.event_id = sqlc.arg('event_id') 
  AND ec.certificate_token_id IS NOT NULL
  AND ec.revoked_at IS NULL
ORDER BY ec.created_at DESC;

-- name: GetUnclaimedReadyCertificatesByEventID :many
SELECT ec.* 
FROM event_certificates ec
INNER JOIN event_certificate_configs ecc ON ec.event_id = ecc.event_id
WHERE ec.event_id = sqlc.arg('event_id') 
  AND ec.certificate_token_id IS NULL
  AND ecc.is_published = TRUE
  AND ec.revoked_at IS NULL
ORDER BY ec.created_at DESC;

-- name: GetClaimedCertificatesByCredentialID :many
SELECT 
    ec.id,
    ec.event_id,
    ec.receiver_credential_id,
    ec.receiver_email,
    ec.name,
    ec.academic_institution,
    ec.certificate_title,
    ec.certificate_subtitle,
    ec.event_contract_address,
    ec.event_certificate_address,
    ec.certificate_token_id,
    ec.certificate_digest,
    ec.inbox_message_id,
    ec.created_at,
    ec.revoked_at,
    e.title as event_name
FROM event_certificates ec
INNER JOIN events e ON ec.event_id = e.id
WHERE (
    ec.receiver_credential_id = sqlc.arg('receiver_credential_id')
    OR ec.receiver_email = sqlc.arg('receiver_email')
  )
  AND ec.certificate_token_id IS NOT NULL
  AND ec.revoked_at IS NULL
ORDER BY ec.created_at DESC;

-- name: GetUnclaimedReadyCertificatesByCredentialID :many
SELECT 
    ec.id,
    ec.event_id,
    ec.receiver_credential_id,
    ec.receiver_email,
    ec.name,
    ec.academic_institution,
    ec.certificate_title,
    ec.certificate_subtitle,
    ec.event_contract_address,
    ec.event_certificate_address,
    ec.certificate_token_id,
    ec.certificate_digest,
    ec.inbox_message_id,
    ec.created_at,
    ec.revoked_at,
    e.title as event_name
FROM event_certificates ec
INNER JOIN event_certificate_configs ecc ON ec.event_id = ecc.event_id
INNER JOIN events e ON ec.event_id = e.id
INNER JOIN event_attendees ea ON ec.event_id = ea.event_id AND ea.attendee_credential_id = sqlc.arg('receiver_credential_id')
WHERE (
    ec.receiver_credential_id = sqlc.arg('receiver_credential_id')
    OR ec.receiver_email = sqlc.arg('receiver_email')
  )
  AND ec.certificate_token_id IS NULL
  AND ecc.is_published = TRUE
  AND ec.revoked_at IS NULL
ORDER BY ec.created_at DESC;
