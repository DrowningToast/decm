-- name: CreateEvent :one
INSERT INTO events (
    chain_id,
    contact_number,
    contact_address,
    owner_credential_id,
    banner_storage_key,
    icon_storage_key,
    title,
    short_description,
    long_description,
    start_date,
    end_date,
    location,
    google_map_query,
    max_attendees,
    is_public,
    is_booking_request_required,
    is_verified,
    is_ticket_transferable,
    event_status
) VALUES (
    sqlc.arg(chain_id),
    sqlc.arg(contact_number),
    sqlc.arg(contact_address),
    sqlc.arg(owner_credential_id),
    sqlc.arg(banner_storage_key),
    sqlc.arg(icon_storage_key),
    sqlc.arg(title),
    sqlc.arg(short_description),
    sqlc.arg(long_description),
    sqlc.arg(start_date),
    sqlc.arg(end_date),
    sqlc.arg(location),
    sqlc.arg(google_map_query),
    sqlc.arg(max_attendees),
    sqlc.arg(is_public),
    sqlc.arg(is_booking_request_required),
    sqlc.arg(is_verified),
    sqlc.arg(is_ticket_transferable),
    sqlc.arg(event_status)
) RETURNING *;

-- name: GetEventById :one
SELECT 
    id,
    event_type,
    chain_id,
    contact_number,
    contact_address,
    owner_credential_id,
    banner_storage_key,
    icon_storage_key,
    title,
    short_description,
    long_description,
    start_date,
    end_date,
    location,
    google_map_query,
    max_attendees,
    is_public,
    is_booking_request_required,
    is_verified,
    is_ticket_transferable,
    created_at,
    updated_at,
    event_status,
    COALESCE(
        (SELECT COUNT(event_attendees.id) 
         FROM event_attendees 
         WHERE event_attendees.event_id = events.id 
           AND event_attendees.is_attendee_accepted::INTEGER = 1
        ), 
    0)::INTEGER AS attendees_count
FROM events
WHERE events.id = sqlc.arg(id);

-- name: GetEventViewModelById :one
SELECT 
    events.*,
    event_registration_configs.*,
    event_contracts.*,
    COALESCE(
        (SELECT COUNT(event_attendees.id) 
         FROM event_attendees 
         WHERE event_attendees.event_id = events.id 
           AND event_attendees.is_attendee_accepted::INTEGER = 1
        ), 
    0)::INTEGER AS attendees_count
FROM events
INNER JOIN event_registration_configs 
    ON events.id = event_registration_configs.event_id
INNER JOIN event_contracts 
    ON events.id = event_contracts.event_id
WHERE events.id = sqlc.arg(id);

-- name: UpdateEvent :one
UPDATE events
SET 
    title = sqlc.arg(title),
    short_description = sqlc.arg(short_description),
    long_description = sqlc.arg(long_description),
    event_type = sqlc.arg(event_type),
    start_date = sqlc.arg(start_date),
    end_date = sqlc.arg(end_date),
    location = sqlc.arg(location),
    google_map_query = sqlc.arg(google_map_query),
    max_attendees = sqlc.arg(max_attendees),
    contact_number = sqlc.arg(contact_number),
    contact_address = sqlc.arg(contact_address),
    banner_storage_key = sqlc.arg(banner_storage_key),
    icon_storage_key = sqlc.arg(icon_storage_key),
    is_public = sqlc.arg(is_public),
    is_booking_request_required = sqlc.arg(is_booking_request_required),
    is_verified = sqlc.arg(is_verified),
    is_ticket_transferable = sqlc.arg(is_ticket_transferable),
    event_status = sqlc.arg(event_status)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteEvent :one
UPDATE events
SET
    deleted_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListEventsByOwner :many
SELECT 
    id,
    chain_id,
    contact_number,
    contact_address,
    owner_credential_id,
    banner_storage_key,
    icon_storage_key,
    title,
    short_description,
    long_description,
    start_date,
    end_date,
    location,
    google_map_query,
    max_attendees,
    is_public,
    is_booking_request_required,
    is_verified,
    is_ticket_transferable,
    created_at,
    updated_at,
    event_status
FROM events
WHERE owner_credential_id = sqlc.arg(owner_credential_id)
ORDER BY created_at DESC;

-- name: ListPublicEvents :many
SELECT 
    id,
    chain_id,
    contact_number,
    contact_address,
    owner_credential_id,
    banner_storage_key,
    icon_storage_key,
    title,
    short_description,
    long_description,
    start_date,
    end_date,
    location,
    google_map_query,
    max_attendees,
    is_public,
    is_booking_request_required,
    is_verified,
    is_ticket_transferable,
    created_at,
    updated_at,
    event_status
FROM events
WHERE is_public = 1
ORDER BY start_date ASC;

-- name: ListEventsByOwnerCredentialID :many
SELECT * 
FROM events 
WHERE owner_credential_id = sqlc.arg(owner_credential_id) AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: ListEvents :many
SELECT * 
FROM events 
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: ListEventsByEventAttendeeCredentialID :many
SELECT * 
FROM events 
INNER JOIN event_attendees ON events.id = event_attendees.event_id
WHERE event_attendees.attendee_credential_id = sqlc.arg(event_attendee_credential_id)
AND events.deleted_at IS NULL
AND event_attendees.is_attendee_accepted = 1
AND events.event_status = 'active'
ORDER BY events.created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);