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
WHERE event_id = sqlc.arg(event_id) AND cancelled_at IS NULL
ORDER BY created_at DESC;

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
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEventRegistrationInvitation :exec
DELETE FROM event_registration_invitations WHERE id = sqlc.arg(id);
