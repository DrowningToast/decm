-- name: CreateInboxMessage :one
INSERT INTO inbox_messages (
    sender_credential_id,
    receiver_credential_id,
    receiver_email,
    receiver_wallet_address,
    message_type,
    message_content,
    fallback_message_content,
    is_read
) VALUES (
    sqlc.narg(sender_credential_id),
    sqlc.narg(receiver_credential_id),
    sqlc.narg(receiver_email),
    sqlc.narg(receiver_wallet_address),
    sqlc.arg(message_type),
    sqlc.narg(message_content),
    sqlc.narg(fallback_message_content),
    sqlc.arg(is_read)
)
RETURNING *;

-- name: GetInboxMessageByID :one
SELECT * FROM inbox_messages WHERE id = sqlc.arg(id);

-- name: GetInboxMessagesByCredentialID :many
SELECT * FROM inbox_messages
WHERE receiver_credential_id = sqlc.arg(receiver_credential_id)
OR receiver_email = sqlc.narg(receiver_email)
OR LOWER(receiver_wallet_address) = sqlc.narg(receiver_wallet_address)
ORDER BY created_at DESC;

-- name: GetUnreadInboxMessageCountByCredentialID :one
SELECT COUNT(*) FROM inbox_messages
WHERE (receiver_credential_id = sqlc.arg(receiver_credential_id)
OR receiver_email = sqlc.narg(receiver_email)
OR LOWER(receiver_wallet_address) = sqlc.narg(receiver_wallet_address))
AND is_read = 0;

-- name: GetInboxMessagesByReceiverEmail :many
SELECT * FROM inbox_messages 
WHERE receiver_email = sqlc.arg(receiver_email)
ORDER BY created_at DESC;

-- name: GetInboxMessagesByReceiverWalletAddress :many
SELECT * FROM inbox_messages
WHERE LOWER(receiver_wallet_address) = sqlc.arg(receiver_wallet_address)
ORDER BY created_at DESC;

-- name: GetInboxMessagesBySenderCredentialID :many
SELECT * FROM inbox_messages 
WHERE sender_credential_id = sqlc.arg(sender_credential_id)
ORDER BY created_at DESC;

-- name: UpdateInboxMessageReadStatus :one
UPDATE inbox_messages 
SET is_read = sqlc.arg(is_read),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateInboxMessageReadStatusAll :many
UPDATE inbox_messages
SET is_read = 1,
    updated_at = NOW()
WHERE receiver_credential_id = sqlc.arg(receiver_credential_id)
OR receiver_email = sqlc.narg(receiver_email)
OR LOWER(receiver_wallet_address) = sqlc.narg(receiver_wallet_address)
RETURNING *;