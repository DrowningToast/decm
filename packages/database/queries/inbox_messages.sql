-- name: CreateInboxMessage :one
INSERT INTO inbox_messages (
    sender_credential_id,
    receiver_credential_id,
    receiver_email,
    message_type,
    message_content,
    fallback_message_content,
    is_read
) VALUES (
    sqlc.narg(sender_credential_id),
    sqlc.narg(receiver_credential_id),
    sqlc.narg(receiver_email),
    sqlc.arg(message_type),
    sqlc.narg(message_content),
    sqlc.narg(fallback_message_content),
    sqlc.arg(is_read)
)
RETURNING *;

-- name: GetInboxMessageByID :one
SELECT * FROM inbox_messages WHERE id = sqlc.arg(id);

-- name: GetInboxMessagesByReceiverEmail :many
SELECT * FROM inbox_messages 
WHERE receiver_email = sqlc.arg(receiver_email)
ORDER BY created_at DESC;

-- name: UpdateInboxMessageReadStatus :one
UPDATE inbox_messages 
SET is_read = sqlc.arg(is_read),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING *;
