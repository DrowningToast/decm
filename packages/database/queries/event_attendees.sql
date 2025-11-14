-- name: GetEventAttendeeByEventIDAndCredentialID :one
SELECT *
FROM event_attendees
WHERE event_id = sqlc.arg(event_id)
AND attendee_credential_id = sqlc.arg(attendee_credential_id);

-- name: AddParticipant :one
INSERT INTO event_attendees (event_id, attendee_credential_id, contract_address, is_attendee_accepted, first_name, last_name, email, bio, phone_number, address, academic_institution, academic_email)
VALUES (sqlc.arg(event_id), sqlc.arg(attendee_credential_id), sqlc.arg(contract_address), sqlc.arg(is_attendee_accepted), sqlc.narg(first_name), sqlc.narg(last_name), sqlc.narg(email), sqlc.narg(bio), sqlc.narg(phone_number), sqlc.narg(address), sqlc.narg(academic_institution), sqlc.narg(academic_email))
RETURNING *;