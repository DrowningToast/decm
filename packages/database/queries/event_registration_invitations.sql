-- name: CreateEventRegistrationInvitation :one
INSERT INTO event_registration_invitations (
    event_id,
    inbox_message_id,
    valid_until,
    code
) VALUES (
    sqlc.arg(event_id),
    sqlc.arg(inbox_message_id),
    sqlc.narg(valid_until),
    sqlc.narg(code)
)
RETURNING *;

-- name: GetEventRegistrationInvitationByID :one
SELECT * FROM event_registration_invitations WHERE id = sqlc.arg(id);

-- name: GetEventRegistrationInvitationsByEventID :many
SELECT * FROM event_registration_invitations 
WHERE event_id = sqlc.arg(event_id)
ORDER BY created_at DESC;

-- name: UpdateEventRegistrationInvitation :one
UPDATE event_registration_invitations 
SET 
    valid_until = sqlc.narg(valid_until),
    code = sqlc.narg(code),
    cancelled_at = sqlc.narg(cancelled_at),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEventRegistrationInvitation :exec
DELETE FROM event_registration_invitations WHERE id = sqlc.arg(id);
