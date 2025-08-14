-- Credentials queries

-- name: CreateCredential :one
INSERT INTO credentials (
    recipient_id,
    issuer_id,
    title,
    description,
    credential_type,
    event_id,
    blockchain_tx_hash,
    contract_address,
    token_id,
    qr_code_hash,
    verification_url,
    skills,
    criteria,
    image_url,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: GetCredentialByID :one
SELECT c.*, 
       r.username as recipient_username, r.first_name as recipient_first_name, r.last_name as recipient_last_name,
       i.username as issuer_username, i.first_name as issuer_first_name, i.last_name as issuer_last_name,
       e.title as event_title, e.start_date as event_date
FROM credentials c
JOIN users r ON c.recipient_id = r.id
JOIN users i ON c.issuer_id = i.id
LEFT JOIN events e ON c.event_id = e.id
WHERE c.id = $1;

-- name: GetCredentialsByRecipient :many
SELECT c.*, 
       i.username as issuer_username, i.first_name as issuer_first_name, i.last_name as issuer_last_name,
       e.title as event_title, e.start_date as event_date
FROM credentials c
JOIN users i ON c.issuer_id = i.id
LEFT JOIN events e ON c.event_id = e.id
WHERE c.recipient_id = $1
  AND ($4::text IS NULL OR c.credential_type = $4)
  AND ($5::text IS NULL OR c.status = $5)
ORDER BY c.issued_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCredentialsByIssuer :many
SELECT c.*, 
       r.username as recipient_username, r.first_name as recipient_first_name, r.last_name as recipient_last_name,
       e.title as event_title, e.start_date as event_date
FROM credentials c
JOIN users r ON c.recipient_id = r.id
LEFT JOIN events e ON c.event_id = e.id
WHERE c.issuer_id = $1
ORDER BY c.issued_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCredentialsByEvent :many
SELECT c.*, 
       r.username as recipient_username, r.first_name as recipient_first_name, r.last_name as recipient_last_name,
       i.username as issuer_username, i.first_name as issuer_first_name, i.last_name as issuer_last_name
FROM credentials c
JOIN users r ON c.recipient_id = r.id
JOIN users i ON c.issuer_id = i.id
WHERE c.event_id = $1
ORDER BY c.issued_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCredentialByQRCode :one
SELECT c.*, 
       r.username as recipient_username, r.first_name as recipient_first_name, r.last_name as recipient_last_name,
       i.username as issuer_username, i.first_name as issuer_first_name, i.last_name as issuer_last_name,
       e.title as event_title, e.start_date as event_date
FROM credentials c
JOIN users r ON c.recipient_id = r.id
JOIN users i ON c.issuer_id = i.id
LEFT JOIN events e ON c.event_id = e.id
WHERE c.qr_code_hash = $1;

-- name: UpdateCredentialStatus :one
UPDATE credentials SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateCredentialBlockchainInfo :one
UPDATE credentials SET 
    blockchain_tx_hash = $2,
    contract_address = COALESCE($3, contract_address),
    token_id = COALESCE($4, token_id),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: RevokeCredential :one
UPDATE credentials SET 
    status = 'revoked',
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetExpiredCredentials :many
SELECT * FROM credentials 
WHERE expires_at < NOW() AND status = 'active'
ORDER BY expires_at ASC;

-- name: GetCredentialStats :one
SELECT 
    COUNT(*) as total_credentials,
    COUNT(*) FILTER (WHERE credential_type = 'certificate') as certificates,
    COUNT(*) FILTER (WHERE credential_type = 'badge') as badges,
    COUNT(*) FILTER (WHERE credential_type = 'achievement') as achievements,
    COUNT(*) FILTER (WHERE status = 'active') as active_credentials,
    COUNT(*) FILTER (WHERE status = 'revoked') as revoked_credentials
FROM credentials
WHERE ($1::uuid IS NULL OR recipient_id = $1)
  AND ($2::uuid IS NULL OR issuer_id = $2);

-- name: DeleteCredential :exec
DELETE FROM credentials WHERE id = $1;
