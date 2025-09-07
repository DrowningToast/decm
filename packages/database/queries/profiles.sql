-- Profiles CRUD queries

-- name: CreateProfile :one
INSERT INTO profiles (
    authentication_credential_id,
    is_profile_picture_public,
    profile_picture_url,
    is_first_name_public,
    first_name,
    is_last_name_public,
    last_name,
    is_email_public,
    email,
    is_bio_public,
    bio,
    is_phone_number_public,
    phone_number,
    is_address_public,
    address,
    is_academic_institution_public,
    academic_institution,
    is_academic_email_public,
    academic_email
) VALUES (
    $1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;

-- name: GetProfileByID :one
SELECT * FROM profiles WHERE id = $1::uuid;

-- name: GetProfileByAuthCredentialID :one
SELECT * FROM profiles WHERE authentication_credential_id = $1::uuid;

-- name: GetProfileByEmail :one
SELECT * FROM profiles WHERE email = $1;

-- name: ListProfiles :many
SELECT * FROM profiles 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProfile :one
UPDATE profiles SET 
    is_profile_picture_public = COALESCE($2, is_profile_picture_public),
    profile_picture_url = COALESCE($3, profile_picture_url),
    is_first_name_public = COALESCE($4, is_first_name_public),
    first_name = COALESCE($5, first_name),
    is_last_name_public = COALESCE($6, is_last_name_public),
    last_name = COALESCE($7, last_name),
    is_email_public = COALESCE($8, is_email_public),
    email = COALESCE($9, email),
    is_bio_public = COALESCE($10, is_bio_public),
    bio = COALESCE($11, bio),
    is_phone_number_public = COALESCE($12, is_phone_number_public),
    phone_number = COALESCE($13, phone_number),
    is_address_public = COALESCE($14, is_address_public),
    address = COALESCE($15, address),
    is_academic_institution_public = COALESCE($16, is_academic_institution_public),
    academic_institution = COALESCE($17, academic_institution),
    is_academic_email_public = COALESCE($18, is_academic_email_public),
    academic_email = COALESCE($19, academic_email),
    updated_at = NOW()
WHERE id = $1::uuid RETURNING *;

-- name: UpdateProfilePersonalInfo :one
UPDATE profiles SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    email = COALESCE($4, email),
    phone_number = COALESCE($5, phone_number),
    address = COALESCE($6, address),
    bio = COALESCE($7, bio),
    updated_at = NOW()
WHERE id = $1::uuid RETURNING *;

-- name: UpdateProfileAcademicInfo :one
UPDATE profiles SET 
    academic_institution = COALESCE($2, academic_institution),
    academic_email = COALESCE($3, academic_email),
    updated_at = NOW()
WHERE id = $1::uuid RETURNING *;

-- name: UpdateProfilePicture :one
UPDATE profiles SET 
    profile_picture_url = $2,
    updated_at = NOW()
WHERE id = $1::uuid RETURNING *;

-- name: UpdateProfilePrivacySettings :one
UPDATE profiles SET 
    is_profile_picture_public = COALESCE($2, is_profile_picture_public),
    is_first_name_public = COALESCE($3, is_first_name_public),
    is_last_name_public = COALESCE($4, is_last_name_public),
    is_email_public = COALESCE($5, is_email_public),
    is_bio_public = COALESCE($6, is_bio_public),
    is_phone_number_public = COALESCE($7, is_phone_number_public),
    is_address_public = COALESCE($8, is_address_public),
    is_academic_institution_public = COALESCE($9, is_academic_institution_public),
    is_academic_email_public = COALESCE($10, is_academic_email_public),
    updated_at = NOW()
WHERE id = $1::uuid RETURNING *;

-- name: SearchProfilesByName :many
SELECT * FROM profiles 
WHERE (first_name ILIKE '%' || $1 || '%' OR last_name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchProfilesByEmail :many
SELECT * FROM profiles 
WHERE email ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchProfilesByAcademicInstitution :many
SELECT * FROM profiles 
WHERE academic_institution ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPublicProfileInfo :one
SELECT 
    id,
    authentication_credential_id,
    CASE WHEN is_profile_picture_public = 1 THEN profile_picture_url ELSE NULL END as profile_picture_url,
    CASE WHEN is_first_name_public = 1 THEN first_name ELSE NULL END as first_name,
    CASE WHEN is_last_name_public = 1 THEN last_name ELSE NULL END as last_name,
    CASE WHEN is_email_public = 1 THEN email ELSE NULL END as email,
    CASE WHEN is_bio_public = 1 THEN bio ELSE NULL END as bio,
    CASE WHEN is_phone_number_public = 1 THEN phone_number ELSE NULL END as phone_number,
    CASE WHEN is_address_public = 1 THEN address ELSE NULL END as address,
    CASE WHEN is_academic_institution_public = 1 THEN academic_institution ELSE NULL END as academic_institution,
    CASE WHEN is_academic_email_public = 1 THEN academic_email ELSE NULL END as academic_email,
    created_at,
    updated_at
FROM profiles 
WHERE id = $1::uuid;

-- name: CountProfiles :one
SELECT COUNT(*) FROM profiles;

-- name: DeleteProfile :exec
DELETE FROM profiles WHERE id = $1::uuid;