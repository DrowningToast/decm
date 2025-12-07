-- name: GetEventAttendeeByEventIDAndCredentialID :one
SELECT *
FROM event_attendees
WHERE event_id = sqlc.arg(event_id)
AND attendee_credential_id = sqlc.arg(attendee_credential_id);

-- name: AddParticipant :one
INSERT INTO event_attendees (event_id, attendee_credential_id, contract_address, is_attendee_accepted, first_name, last_name, email, bio, phone_number, address, academic_institution, academic_email)
VALUES (sqlc.arg(event_id), sqlc.arg(attendee_credential_id), sqlc.arg(contract_address), sqlc.arg(is_attendee_accepted), sqlc.narg(first_name), sqlc.narg(last_name), sqlc.narg(email), sqlc.narg(bio), sqlc.narg(phone_number), sqlc.narg(address), sqlc.narg(academic_institution), sqlc.narg(academic_email))
RETURNING *;

-- name: ListEventAttendeesByEventID :many
SELECT
    ea.id,
    ea.event_id,
    ea.attendee_credential_id,
    ea.contract_address,
    ea.is_attendee_accepted,
    ea.first_name,
    ea.last_name,
    ea.email,
    ea.phone_number,
    ea.academic_institution,
    ea.created_at,
    ea.updated_at,
    ac.wallet_address
FROM event_attendees ea
INNER JOIN authentication_credentials ac ON ea.attendee_credential_id = ac.id
WHERE ea.event_id = sqlc.arg(event_id)
ORDER BY ea.created_at DESC;