-- Events queries

-- name: CreateEvent :one
INSERT INTO events (
    title,
    description,
    event_type,
    start_date,
    end_date,
    location,
    max_attendees,
    organizer_id,
    organization,
    nft_contract_address,
    ticket_price,
    tags,
    image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetEventByID :one
SELECT e.*, u.username as organizer_username, u.first_name as organizer_first_name, u.last_name as organizer_last_name
FROM events e
JOIN users u ON e.organizer_id = u.id
WHERE e.id = $1;

-- name: ListEvents :many
SELECT e.*, u.username as organizer_username, u.first_name as organizer_first_name, u.last_name as organizer_last_name
FROM events e
JOIN users u ON e.organizer_id = u.id
WHERE ($3::text IS NULL OR e.event_type = $3)
  AND ($4::text IS NULL OR e.status = $4)
  AND ($5::timestamptz IS NULL OR e.start_date >= $5)
  AND ($6::timestamptz IS NULL OR e.end_date <= $6)
ORDER BY e.start_date ASC
LIMIT $1 OFFSET $2;

-- name: GetEventsByOrganizer :many
SELECT * FROM events 
WHERE organizer_id = $1
ORDER BY start_date DESC
LIMIT $2 OFFSET $3;

-- name: SearchEvents :many
SELECT e.*, u.username as organizer_username, u.first_name as organizer_first_name, u.last_name as organizer_last_name
FROM events e
JOIN users u ON e.organizer_id = u.id
WHERE e.search_vector @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(e.search_vector, plainto_tsquery('english', $1)) DESC
LIMIT $2 OFFSET $3;

-- name: UpdateEvent :one
UPDATE events SET 
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    event_type = COALESCE($4, event_type),
    start_date = COALESCE($5, start_date),
    end_date = COALESCE($6, end_date),
    location = COALESCE($7, location),
    max_attendees = COALESCE($8, max_attendees),
    organization = COALESCE($9, organization),
    nft_contract_address = COALESCE($10, nft_contract_address),
    ticket_price = COALESCE($11, ticket_price),
    tags = COALESCE($12, tags),
    image_url = COALESCE($13, image_url),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateEventStatus :one
UPDATE events SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetUpcomingEvents :many
SELECT e.*, u.username as organizer_username, u.first_name as organizer_first_name, u.last_name as organizer_last_name
FROM events e
JOIN users u ON e.organizer_id = u.id
WHERE e.start_date > NOW() 
  AND e.status = 'published'
ORDER BY e.start_date ASC
LIMIT $1 OFFSET $2;

-- name: GetEventAttendeeCount :one
SELECT COUNT(*) FROM nft_tickets WHERE event_id = $1;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;
