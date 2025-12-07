-- name: CreateEventRegistrationInvitation :one
INSERT INTO event_registration_invitations (
    event_id,
    inbox_message_id,
    valid_until,
    code,
    first_name,
    last_name,
    email,
    phone_number,
    academic_institution
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(inbox_message_id),
    sqlc.narg(valid_until),
    sqlc.narg(code),
    sqlc.narg(first_name),
    sqlc.narg(last_name),
    sqlc.narg(email),
    sqlc.narg(phone_number),
    sqlc.narg(academic_institution)
)
RETURNING *;

-- name: GetEventRegistrationInvitationByID :one
SELECT * FROM event_registration_invitations WHERE id = sqlc.arg(id);

-- name: GetEventRegistrationInvitationsByEventID :many
SELECT * FROM event_registration_invitations 
WHERE event_id = sqlc.arg(event_id) 
AND cancelled_at IS NULL
ORDER BY created_at DESC;

-- name: GetEventRegistrationInvitationByInboxMessageID :one
SELECT * FROM event_registration_invitations 
WHERE inbox_message_id = sqlc.arg(inbox_message_id)
AND cancelled_at IS NULL
ORDER BY created_at DESC;

-- name: GetEventRegistrationInvitationByEventIdAndCredentialId :one
SELECT 
    event_registration_invitations.id as event_registration_invitation_id,
    event_registration_invitations.event_id as event_id,
    event_registration_invitations.valid_until as valid_until,
    event_registration_invitations.code as code,
    event_registration_invitations.first_name as first_name,
    event_registration_invitations.last_name as last_name,
    event_registration_invitations.email as email,
    event_registration_invitations.phone_number as phone_number,
    event_registration_invitations.academic_institution as academic_institution,
    event_registration_invitations.created_at as event_registration_invitation_created_at,
    event_registration_invitations.updated_at as event_registration_invitation_updated_at,
    event_registration_invitations.cancelled_at as event_registration_invitation_cancelled_at,
    event_registration_invitations.accepted_at as event_registration_invitation_accepted_at,
    inbox_messages.id as inbox_message_id,
    inbox_messages.id as sender_credential_id,
    inbox_messages.receiver_credential_id as receiver_credential_id,
    inbox_messages.receiver_email as receiver_email,
    inbox_messages.receiver_wallet_address as receiver_wallet_address,
    inbox_messages.message_type as message_type,
    inbox_messages.message_content as message_content,
    inbox_messages.fallback_message_content as fallback_message_content,
    inbox_messages.created_at as inbox_message_created_at,
    inbox_messages.updated_at as inbox_message_updated_at,
    inbox_messages.is_read as inbox_message_is_read,
    inbox_messages.hidden_at as inbox_message_hidden_at,
    inbox_messages.deleted_at as inbox_message_deleted_at
FROM event_registration_invitations 
INNER JOIN inbox_messages ON event_registration_invitations.inbox_message_id = inbox_messages.id
WHERE event_registration_invitations.event_id = sqlc.arg(event_id) 
AND event_registration_invitations.cancelled_at IS NULL
AND (
    inbox_messages.receiver_credential_id = sqlc.arg(credential_id)
    OR inbox_messages.receiver_email = sqlc.narg(email)
    OR inbox_messages.receiver_wallet_address = sqlc.narg(wallet_address)
)
ORDER BY event_registration_invitations.created_at DESC;

-- name: UpdateEventRegistrationInvitation :one
UPDATE event_registration_invitations 
SET 
    valid_until = sqlc.narg(valid_until),
    code = sqlc.narg(code),
    cancelled_at = sqlc.narg(cancelled_at),
    first_name = sqlc.narg(first_name),
    last_name = sqlc.narg(last_name),
    email = sqlc.narg(email),
    phone_number = sqlc.narg(phone_number),
    academic_institution = sqlc.narg(academic_institution),
    accepted_at = sqlc.narg(accepted_at),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateEventRegistrationInvitationAcceptedStatus :one
UPDATE event_registration_invitations 
SET 
    accepted_at = sqlc.narg(accepted_at),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEventRegistrationInvitation :exec
DELETE FROM event_registration_invitations WHERE id = sqlc.arg(id);
