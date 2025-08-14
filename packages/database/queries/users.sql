-- Users queries

-- name: CreateUser :one
INSERT INTO users (
    email, 
    username, 
    first_name, 
    last_name, 
    academic_institution, 
    academic_email,
    bio,
    profile_picture_url,
    wallet_address
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByWalletAddress :one
SELECT * FROM users WHERE wallet_address = $1;

-- name: UpdateUser :one
UPDATE users SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    bio = COALESCE($4, bio),
    profile_picture_url = COALESCE($5, profile_picture_url),
    academic_institution = COALESCE($6, academic_institution),
    academic_email = COALESCE($7, academic_email),
    wallet_address = COALESCE($8, wallet_address),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateUserVerificationStatus :one
UPDATE users SET 
    verification_status = $2,
    ldap_verified = $3,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users SET last_login_at = NOW() WHERE id = $1;

-- name: SearchUsers :many
SELECT * FROM users 
WHERE search_vector @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(search_vector, plainto_tsquery('english', $1)) DESC
LIMIT $2 OFFSET $3;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
