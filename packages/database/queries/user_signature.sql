-- name: CreateUserSignature :one
INSERT INTO user_signature (
    authentication_credential_id,
    sign_message,
    signature,
    deadline_block,
    estimated_deadline,
    broadcasted_at
) VALUES (
    sqlc.arg(authentication_credential_id),
    sqlc.arg(sign_message),
    sqlc.arg(signature),
    sqlc.narg(deadline_block),
    sqlc.narg(estimated_deadline),
    sqlc.narg(broadcasted_at)
) RETURNING *;

-- name: GetUserSignatureByID :one
SELECT * FROM user_signature WHERE id = sqlc.arg(id);

-- name: GetUserSignaturesByCredentialID :many
SELECT * FROM user_signature
WHERE authentication_credential_id = sqlc.arg(authentication_credential_id)
ORDER BY created_at DESC;

-- name: UpdateUserSignatureBroadcastedAt :one
UPDATE user_signature
SET broadcasted_at = sqlc.arg(broadcasted_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserSignatureMarkAsExpiredAt :one
UPDATE user_signature
SET mark_as_expired_at = sqlc.arg(mark_as_expired_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateUserSignatureAbortedAt :one
UPDATE user_signature
SET aborted_at = sqlc.arg(aborted_at),
    aborted_reason = sqlc.arg(aborted_reason)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: GetPendingUserSignatures :many
SELECT
    us.*,
    ac.wallet_address
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
WHERE us.broadcasted_at IS NULL
ORDER BY us.created_at DESC;

-- name: GetStaleUserSignatures :many
SELECT
    us.*,
    ac.wallet_address
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
WHERE us.broadcasted_at IS NULL
  AND us.created_at < sqlc.arg(created_before)
ORDER BY us.created_at ASC;

-- name: GetBroadcastedUserSignatures :many
SELECT
    us.*,
    ac.wallet_address
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
WHERE us.broadcasted_at IS NOT NULL
ORDER BY us.broadcasted_at DESC;

-- name: GetUserSignaturesByDeadlineBlockRange :many
SELECT * FROM user_signature
WHERE deadline_block BETWEEN sqlc.arg(min_block) AND sqlc.arg(max_block)
ORDER BY deadline_block ASC;

-- name: GetUserSignaturesExpiringBefore :many
SELECT
    us.*,
    ac.wallet_address
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
WHERE us.estimated_deadline IS NOT NULL
  AND us.estimated_deadline < sqlc.arg(deadline)
  AND us.broadcasted_at IS NULL
ORDER BY us.estimated_deadline ASC;

-- name: GetEventJoinSignatures :many
SELECT
    us.id,
    us.authentication_credential_id,
    us.sign_message,
    us.signature,
    us.broadcasted_at,
    us.deadline_block,
    us.estimated_deadline,
    us.mark_as_expired_at,
    us.aborted_at,
    us.created_at,
    ea.event_id,
    ea.attendee_credential_id,
    e.title as event_title,
    ac.wallet_address
FROM user_signature us
JOIN event_attendees ea ON us.id = ea.user_signature_id
JOIN events e ON ea.event_id = e.id
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
ORDER BY us.created_at DESC;

-- name: GetPendingEventJoinSignatures :many
SELECT
    us.id,
    us.authentication_credential_id,
    us.sign_message,
    us.signature,
    us.deadline_block,
    us.estimated_deadline,
    us.created_at,
    ea.event_id,
    ac.wallet_address
FROM user_signature us
JOIN event_attendees ea ON us.id = ea.user_signature_id
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
WHERE us.broadcasted_at IS NULL
  AND us.mark_as_expired_at IS NULL
  AND us.aborted_at IS NULL
ORDER BY us.created_at ASC;

-- name: GetCertificateClaimSignatures :many
SELECT
    us.id,
    us.authentication_credential_id,
    us.sign_message,
    us.signature,
    us.broadcasted_at,
    us.deadline_block,
    us.estimated_deadline,
    us.mark_as_expired_at,
    us.aborted_at,
    us.created_at,
    ec.id as certificate_id,
    ec.event_id,
    ec.certificate_token_id,
    e.title as event_title,
    ac.wallet_address
FROM user_signature us
JOIN event_certificates ec ON us.id = ec.user_claim_signature_id
JOIN events e ON ec.event_id = e.id
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
ORDER BY us.created_at DESC;

-- name: GetPendingCertificateClaimSignatures :many
SELECT
    us.id,
    us.authentication_credential_id,
    us.sign_message,
    us.signature,
    us.deadline_block,
    us.estimated_deadline,
    us.created_at,
    ec.id as certificate_id,
    ec.event_id,
    ac.wallet_address
FROM user_signature us
JOIN event_certificates ec ON us.id = ec.user_claim_signature_id
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
WHERE us.broadcasted_at IS NULL
  AND us.mark_as_expired_at IS NULL
  AND us.aborted_at IS NULL
  AND ec.certificate_token_id IS NULL
ORDER BY us.created_at ASC;

-- name: GetOrphanedUserSignatures :many
SELECT
    us.id,
    us.authentication_credential_id,
    us.created_at,
    us.broadcasted_at,
    ac.wallet_address
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
LEFT JOIN event_attendees ea ON us.id = ea.user_signature_id
LEFT JOIN event_certificates ec ON us.id = ec.user_claim_signature_id
WHERE us.broadcasted_at IS NULL
  AND us.created_at < sqlc.arg(created_before)
  AND ea.id IS NULL
  AND ec.id IS NULL
ORDER BY us.created_at ASC;

-- name: GetUserSignatureWithUsageDetails :one
SELECT
    us.id,
    us.authentication_credential_id,
    us.sign_message,
    us.signature,
    us.broadcasted_at,
    us.deadline_block,
    us.estimated_deadline,
    us.created_at,
    us.updated_at,
    ac.wallet_address,
    ea.event_id as joined_event_id,
    ea.created_at as joined_at,
    ec.id as claimed_certificate_id,
    ec.certificate_token_id,
    ec.created_at as certificate_claimed_at
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
LEFT JOIN event_attendees ea ON us.id = ea.user_signature_id
LEFT JOIN event_certificates ec ON us.id = ec.user_claim_signature_id
WHERE us.id = sqlc.arg(id);

-- name: GetRecentUserSignatureActivity :many
SELECT
    us.id,
    us.authentication_credential_id,
    us.created_at,
    us.broadcasted_at,
    us.deadline_block,
    us.estimated_deadline,
    ac.wallet_address,
    ea.event_id as event_join_id,
    ec.id as cert_claim_id
FROM user_signature us
JOIN authentication_credentials ac ON us.authentication_credential_id = ac.id
LEFT JOIN event_attendees ea ON us.id = ea.user_signature_id
LEFT JOIN event_certificates ec ON us.id = ec.user_claim_signature_id
WHERE us.created_at >= sqlc.arg(since)
ORDER BY us.created_at DESC;

-- name: DeleteUserSignature :exec
DELETE FROM user_signature WHERE id = sqlc.arg(id);

-- name: CountUserSignaturesByCredentialID :one
SELECT COUNT(*) FROM user_signature
WHERE authentication_credential_id = sqlc.arg(authentication_credential_id);

-- name: CountPendingUserSignatures :one
SELECT COUNT(*) FROM user_signature
WHERE broadcasted_at IS NULL;

-- name: CountBroadcastedUserSignatures :one
SELECT COUNT(*) FROM user_signature
WHERE broadcasted_at IS NOT NULL;
