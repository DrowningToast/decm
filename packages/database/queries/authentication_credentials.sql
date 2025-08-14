-- Authentication Credentials CRUD queries

-- name: CreateAuthenticationCredential :one
INSERT INTO authentication_credentials (
    solution_status,
    password,
    private_key,
    public_key,
    google_connector_ref,
    github_connector_ref,
    is_verified_organizer,
    is_verified_student
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAuthenticationCredentialByID :one
SELECT * FROM authentication_credentials WHERE id = $1;

-- name: GetAuthenticationCredentialByPublicKey :one
SELECT * FROM authentication_credentials WHERE public_key = $1;

-- name: ListAuthenticationCredentials :many
SELECT * FROM authentication_credentials 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAuthenticationCredential :one
UPDATE authentication_credentials SET 
    solution_status = COALESCE($2, solution_status),
    password = COALESCE($3, password),
    private_key = COALESCE($4, private_key),
    public_key = COALESCE($5, public_key),
    google_connector_ref = COALESCE($6, google_connector_ref),
    github_connector_ref = COALESCE($7, github_connector_ref),
    is_verified_organizer = COALESCE($8, is_verified_organizer),
    is_verified_student = COALESCE($9, is_verified_student),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateAuthenticationCredentialPassword :one
UPDATE authentication_credentials SET 
    password = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateAuthenticationCredentialKeys :one
UPDATE authentication_credentials SET 
    private_key = $2,
    public_key = $3,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateVerificationStatus :one
UPDATE authentication_credentials SET 
    is_verified_organizer = $2,
    is_verified_student = $3,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: SetGoogleConnector :one
UPDATE authentication_credentials SET 
    google_connector_ref = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: SetGithubConnector :one
UPDATE authentication_credentials SET 
    github_connector_ref = $2,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: RemoveGoogleConnector :one
UPDATE authentication_credentials SET 
    google_connector_ref = NULL,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: RemoveGithubConnector :one
UPDATE authentication_credentials SET 
    github_connector_ref = NULL,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetCredentialsByVerificationStatus :many
SELECT * FROM authentication_credentials 
WHERE ($3::int IS NULL OR is_verified_organizer = $3)
  AND ($4::int IS NULL OR is_verified_student = $4)
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetCredentialsBySolutionStatus :many
SELECT * FROM authentication_credentials 
WHERE solution_status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

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
DELETE FROM authentication_credentials WHERE id = $1;

-- name: SoftDeleteAuthenticationCredential :one
-- Note: This would require adding a deleted_at column in future migration
-- For now, we can use a status update approach
UPDATE authentication_credentials SET 
    password = NULL,
    private_key = NULL,
    updated_at = NOW()
WHERE id = $1 RETURNING *;
