-- Evaluations queries

-- name: CreateEvaluation :one
INSERT INTO evaluations (
    event_id,
    evaluator_id,
    rating,
    feedback,
    content_rating,
    organization_rating,
    venue_rating,
    is_anonymous
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetEvaluationByID :one
SELECT e.*, 
       u.username as evaluator_username, u.first_name as evaluator_first_name, u.last_name as evaluator_last_name,
       ev.title as event_title, ev.start_date as event_date
FROM evaluations e
LEFT JOIN users u ON e.evaluator_id = u.id AND e.is_anonymous = false
JOIN events ev ON e.event_id = ev.id
WHERE e.id = $1;

-- name: GetEvaluationsByEvent :many
SELECT e.*, 
       u.username as evaluator_username, u.first_name as evaluator_first_name, u.last_name as evaluator_last_name
FROM evaluations e
LEFT JOIN users u ON e.evaluator_id = u.id AND e.is_anonymous = false
WHERE e.event_id = $1
ORDER BY e.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetEvaluationsByUser :many
SELECT e.*, 
       ev.title as event_title, ev.start_date as event_date
FROM evaluations e
JOIN events ev ON e.event_id = ev.id
WHERE e.evaluator_id = $1
ORDER BY e.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetEventEvaluationStats :one
SELECT 
    COUNT(*) as total_evaluations,
    AVG(rating)::numeric(3,2) as average_rating,
    AVG(content_rating)::numeric(3,2) as average_content_rating,
    AVG(organization_rating)::numeric(3,2) as average_organization_rating,
    AVG(venue_rating)::numeric(3,2) as average_venue_rating,
    COUNT(*) FILTER (WHERE rating = 1) as rating_1_count,
    COUNT(*) FILTER (WHERE rating = 2) as rating_2_count,
    COUNT(*) FILTER (WHERE rating = 3) as rating_3_count,
    COUNT(*) FILTER (WHERE rating = 4) as rating_4_count,
    COUNT(*) FILTER (WHERE rating = 5) as rating_5_count
FROM evaluations 
WHERE event_id = $1;

-- name: GetOrganizerEvaluationStats :one
SELECT 
    COUNT(DISTINCT e.event_id) as total_events_evaluated,
    COUNT(*) as total_evaluations,
    AVG(ev.rating)::numeric(3,2) as average_rating,
    AVG(ev.content_rating)::numeric(3,2) as average_content_rating,
    AVG(ev.organization_rating)::numeric(3,2) as average_organization_rating,
    AVG(ev.venue_rating)::numeric(3,2) as average_venue_rating
FROM events e
JOIN evaluations ev ON e.id = ev.event_id
WHERE e.organizer_id = $1;

-- name: UpdateEvaluation :one
UPDATE evaluations SET 
    rating = COALESCE($2, rating),
    feedback = COALESCE($3, feedback),
    content_rating = COALESCE($4, content_rating),
    organization_rating = COALESCE($5, organization_rating),
    venue_rating = COALESCE($6, venue_rating),
    is_anonymous = COALESCE($7, is_anonymous),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: CheckUserEvaluationExists :one
SELECT EXISTS(
    SELECT 1 FROM evaluations 
    WHERE event_id = $1 AND evaluator_id = $2
);

-- name: GetTopRatedEvents :many
SELECT 
    e.*,
    u.username as organizer_username, u.first_name as organizer_first_name, u.last_name as organizer_last_name,
    AVG(ev.rating)::numeric(3,2) as average_rating,
    COUNT(ev.id) as evaluation_count
FROM events e
JOIN users u ON e.organizer_id = u.id
JOIN evaluations ev ON e.id = ev.event_id
WHERE e.status = 'completed'
GROUP BY e.id, u.username, u.first_name, u.last_name
HAVING COUNT(ev.id) >= $3  -- Minimum number of evaluations
ORDER BY AVG(ev.rating) DESC, COUNT(ev.id) DESC
LIMIT $1 OFFSET $2;

-- name: DeleteEvaluation :exec
DELETE FROM evaluations WHERE id = $1;
