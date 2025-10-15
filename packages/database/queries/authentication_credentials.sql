-- Authentication Credentials CRUD queries
-- Note: Encryption is handled at the repository layer using AES-GCM

-- name: CreateAuthenticationCredential :one
INSERT INTO authentication_credentials (
    solution_status,
    hashed_password,
    encrypted_private_key,
    wallet_address,
    google_connector_ref,
    github_connector_ref,
    is_verified_organizer,
    is_verified_student
) VALUES (
    sqlc.arg(solution_status),
    sqlc.narg(hashed_password),
    sqlc.narg(encrypted_private_key),
    sqlc.arg(wallet_address),
    sqlc.narg(google_connector_ref),
    sqlc.narg(github_connector_ref),
    sqlc.arg(is_verified_organizer),
    sqlc.arg(is_verified_student)
) RETURNING *;

-- name: GetAuthenticationCredentialById :one
SELECT * FROM authentication_credentials 
WHERE id = sqlc.arg(id);

-- name: GetAuthenticationCredentialByWalletAddress :one
SELECT * FROM authentication_credentials 
WHERE wallet_address = sqlc.arg(wallet_address);

-- name: GetAuthenticationCredentialByGoogleConnectorRef :one
SELECT * FROM authentication_credentials 
WHERE google_connector_ref = sqlc.arg(google_connector_ref);

-- name: ListAuthenticationCredentials :many
SELECT * FROM authentication_credentials 
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: UpdateAuthenticationCredential :one
UPDATE authentication_credentials SET 
    solution_status = COALESCE(sqlc.narg(solution_status), solution_status),
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    encrypted_private_key = COALESCE(sqlc.narg(encrypted_private_key), encrypted_private_key),
    wallet_address = COALESCE(sqlc.narg(wallet_address), wallet_address),
    google_connector_ref = COALESCE(sqlc.narg(google_connector_ref), google_connector_ref),
    github_connector_ref = COALESCE(sqlc.narg(github_connector_ref), github_connector_ref),
    is_verified_organizer = COALESCE(sqlc.narg(is_verified_organizer), is_verified_organizer),
    is_verified_student = COALESCE(sqlc.narg(is_verified_student), is_verified_student),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: UpdateAuthenticationCredentialPassword :one
UPDATE authentication_credentials SET 
    hashed_password = sqlc.narg(hashed_password),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: UpdateAuthenticationCredentialKeys :one
UPDATE authentication_credentials SET 
    encrypted_private_key = sqlc.narg(encrypted_private_key),
    wallet_address = sqlc.narg(wallet_address),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: UpdateVerificationStatus :one
UPDATE authentication_credentials SET 
    is_verified_organizer = sqlc.arg(is_verified_organizer),
    is_verified_student = sqlc.arg(is_verified_student),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: SetGoogleConnector :one
UPDATE authentication_credentials SET 
    google_connector_ref = sqlc.arg(google_connector_ref),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: SetGithubConnector :one
UPDATE authentication_credentials SET 
    github_connector_ref = sqlc.arg(github_connector_ref),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: RemoveGoogleConnector :one
UPDATE authentication_credentials SET 
    google_connector_ref = NULL,
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: RemoveGithubConnector :one
UPDATE authentication_credentials SET 
    github_connector_ref = NULL,
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: GetCredentialsByVerificationStatus :many
SELECT * FROM authentication_credentials 
WHERE (sqlc.narg(is_verified_organizer) IS NULL OR is_verified_organizer = sqlc.narg(is_verified_organizer))
  AND (sqlc.narg(is_verified_student) IS NULL OR is_verified_student = sqlc.narg(is_verified_student))
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: GetCredentialsBySolutionStatus :many
SELECT * FROM authentication_credentials 
WHERE solution_status = sqlc.arg(solution_status)
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: CountAuthenticationCredentials :one
SELECT COUNT(*) FROM authentication_credentials;

-- name: CountCredentialsByVerificationStatus :one
SELECT 
    COUNT(*) as total_credentials,
    COUNT(*) FILTER (WHERE is_verified_organizer = 1) as verified_organizers,
    COUNT(*) FILTER (WHERE is_verified_student = 1) as verified_students,
    COUNT(*) FILTER (WHERE solution_status = 0) as byok_credentials,
    COUNT(*) FILTER (WHERE solution_status = 1) as system_managed_credentials
FROM authentication_credentials;

-- name: DeleteAuthenticationCredential :exec
DELETE FROM authentication_credentials WHERE id = sqlc.arg(id);

-- name: SoftDeleteAuthenticationCredential :one
-- Note: This would require adding a deleted_at column in future migration
-- For now, we can use a status update approach
UPDATE authentication_credentials SET 
    hashed_password = NULL,
    encrypted_private_key = NULL,
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;
