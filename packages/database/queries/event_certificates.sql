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
SELECT * FROM event_certificates WHERE event_id = sqlc.arg('event_id');

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
    revoked_at = sqlc.arg('revoked_at'),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEventCertificate :exec
DELETE FROM event_certificates WHERE id = sqlc.arg('id');
