-- name: GetEventAttendeeByEventIDAndCredentialID :one
SELECT *
FROM event_attendees
WHERE event_id = sqlc.arg(event_id)
AND attendee_credential_id = sqlc.arg(attendee_credential_id);