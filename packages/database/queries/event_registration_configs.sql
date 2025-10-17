-- name: CreateEventRegistrationConfig :one
INSERT INTO event_registration_configs (
    event_id,
    final_call_for_registration,
    registration_password,
    first_name_requirement_status,
    last_name_requirement_status,
    email_requirement_status,
    bio_requirement_status,
    phone_number_requirement_status,
    address_requirement_status,
    academic_institution_requirement_status,
    academic_email_requirement_status
) VALUES (
    sqlc.arg('event_id'),
    sqlc.arg('final_call_for_registration'),
    sqlc.arg('registration_password'),
    sqlc.arg('first_name_requirement_status'),
    sqlc.arg('last_name_requirement_status'),
    sqlc.arg('email_requirement_status'),
    sqlc.arg('bio_requirement_status'),
    sqlc.arg('phone_number_requirement_status'),
    sqlc.arg('address_requirement_status'),
    sqlc.arg('academic_institution_requirement_status'),
    sqlc.arg('academic_email_requirement_status')
) RETURNING *;

-- name: GetEventRegistrationConfigByEventID :one
SELECT * FROM event_registration_configs WHERE event_id = sqlc.arg('event_id');

-- name: UpdateEventRegistrationConfig :one
UPDATE event_registration_configs
SET 
    final_call_for_registration = sqlc.arg('final_call_for_registration'),
    registration_password = sqlc.arg('registration_password'),
    first_name_requirement_status = sqlc.arg('first_name_requirement_status'),
    last_name_requirement_status = sqlc.arg('last_name_requirement_status'),
    email_requirement_status = sqlc.arg('email_requirement_status'),
    bio_requirement_status = sqlc.arg('bio_requirement_status'),
    phone_number_requirement_status = sqlc.arg('phone_number_requirement_status'),
    address_requirement_status = sqlc.arg('address_requirement_status'),
    academic_institution_requirement_status = sqlc.arg('academic_institution_requirement_status'),
    academic_email_requirement_status = sqlc.arg('academic_email_requirement_status'),
    updated_at = NOW()
WHERE event_id = sqlc.arg('event_id')
RETURNING *;

-- name: DeleteEventRegistrationConfig :exec
DELETE FROM event_registration_configs WHERE event_id = sqlc.arg('event_id');
