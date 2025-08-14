-- Portfolios queries

-- name: CreatePortfolio :one
INSERT INTO portfolios (
    owner_id,
    title,
    description,
    is_public,
    share_expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPortfolioByID :one
SELECT p.*, u.username as owner_username, u.first_name as owner_first_name, u.last_name as owner_last_name
FROM portfolios p
JOIN users u ON p.owner_id = u.id
WHERE p.id = $1;

-- name: GetPortfolioByShareToken :one
SELECT p.*, u.username as owner_username, u.first_name as owner_first_name, u.last_name as owner_last_name
FROM portfolios p
JOIN users u ON p.owner_id = u.id
WHERE p.share_token = $1
  AND (p.share_expires_at IS NULL OR p.share_expires_at > NOW());

-- name: GetPortfoliosByOwner :many
SELECT * FROM portfolios 
WHERE owner_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePortfolio :one
UPDATE portfolios SET 
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    is_public = COALESCE($4, is_public),
    share_expires_at = COALESCE($5, share_expires_at),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: RegenerateShareToken :one
UPDATE portfolios SET 
    share_token = uuid_generate_v4(),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetPublicPortfolios :many
SELECT p.*, u.username as owner_username, u.first_name as owner_first_name, u.last_name as owner_last_name
FROM portfolios p
JOIN users u ON p.owner_id = u.id
WHERE p.is_public = true
ORDER BY p.updated_at DESC
LIMIT $1 OFFSET $2;

-- name: DeletePortfolio :exec
DELETE FROM portfolios WHERE id = $1;

-- Portfolio Items queries

-- name: CreatePortfolioItem :one
INSERT INTO portfolio_items (
    portfolio_id,
    item_type,
    credential_id,
    event_id,
    custom_title,
    custom_description,
    custom_date,
    custom_image_url,
    display_order
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetPortfolioItems :many
SELECT 
    pi.*,
    c.title as credential_title, c.description as credential_description, c.credential_type, c.issued_at as credential_date, c.image_url as credential_image,
    e.title as event_title, e.description as event_description, e.start_date as event_date, e.image_url as event_image,
    ci.username as credential_issuer_username, ci.first_name as credential_issuer_first_name, ci.last_name as credential_issuer_last_name
FROM portfolio_items pi
LEFT JOIN credentials c ON pi.credential_id = c.id
LEFT JOIN events e ON pi.event_id = e.id
LEFT JOIN users ci ON c.issuer_id = ci.id
WHERE pi.portfolio_id = $1
ORDER BY pi.display_order ASC, pi.created_at DESC;

-- name: UpdatePortfolioItemOrder :exec
UPDATE portfolio_items SET 
    display_order = $2,
    created_at = NOW()  -- Using created_at as updated timestamp for portfolio items
WHERE id = $1;

-- name: DeletePortfolioItem :exec
DELETE FROM portfolio_items WHERE id = $1;

-- name: GetPortfolioItemCount :one
SELECT COUNT(*) FROM portfolio_items WHERE portfolio_id = $1;
