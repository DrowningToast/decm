-- NFT Tickets queries

-- name: CreateNFTTicket :one
INSERT INTO nft_tickets (
    event_id,
    holder_id,
    token_id,
    contract_address,
    blockchain_tx_hash,
    ticket_type,
    seat_number,
    qr_code_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetNFTTicketByID :one
SELECT nt.*, e.title as event_title, e.start_date as event_start_date, e.location as event_location,
       u.username as holder_username, u.first_name as holder_first_name, u.last_name as holder_last_name
FROM nft_tickets nt
JOIN events e ON nt.event_id = e.id
JOIN users u ON nt.holder_id = u.id
WHERE nt.id = $1;

-- name: GetNFTTicketByTokenID :one
SELECT nt.*, e.title as event_title, e.start_date as event_start_date, e.location as event_location,
       u.username as holder_username, u.first_name as holder_first_name, u.last_name as holder_last_name
FROM nft_tickets nt
JOIN events e ON nt.event_id = e.id
JOIN users u ON nt.holder_id = u.id
WHERE nt.contract_address = $1 AND nt.token_id = $2;

-- name: GetTicketsByUser :many
SELECT nt.*, e.title as event_title, e.start_date as event_start_date, e.location as event_location, e.status as event_status
FROM nft_tickets nt
JOIN events e ON nt.event_id = e.id
WHERE nt.holder_id = $1
ORDER BY e.start_date DESC
LIMIT $2 OFFSET $3;

-- name: GetTicketsByEvent :many
SELECT nt.*, u.username as holder_username, u.first_name as holder_first_name, u.last_name as holder_last_name
FROM nft_tickets nt
JOIN users u ON nt.holder_id = u.id
WHERE nt.event_id = $1
ORDER BY nt.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CheckInTicket :one
UPDATE nft_tickets SET 
    checked_in = true,
    check_in_time = NOW(),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetTicketByQRCode :one
SELECT nt.*, e.title as event_title, e.start_date as event_start_date, e.location as event_location,
       u.username as holder_username, u.first_name as holder_first_name, u.last_name as holder_last_name
FROM nft_tickets nt
JOIN events e ON nt.event_id = e.id
JOIN users u ON nt.holder_id = u.id
WHERE nt.qr_code_hash = $1;

-- name: UpdateNFTTicketBlockchainInfo :one
UPDATE nft_tickets SET 
    blockchain_tx_hash = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetEventCheckedInCount :one
SELECT COUNT(*) FROM nft_tickets WHERE event_id = $1 AND checked_in = true;

-- name: DeleteNFTTicket :exec
DELETE FROM nft_tickets WHERE id = $1;
