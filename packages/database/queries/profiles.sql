-- Profiles CRUD queries
-- Note: PII encryption is handled at the repository layer using AES-GCM

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
    sqlc.arg(authentication_credential_id),
    sqlc.arg(is_profile_picture_public),
    sqlc.narg(profile_picture_url),
    sqlc.arg(is_first_name_public),
    sqlc.narg(first_name),
    sqlc.arg(is_last_name_public),
    sqlc.narg(last_name),
    sqlc.arg(is_email_public),
    sqlc.narg(email),
    sqlc.arg(is_bio_public),
    sqlc.narg(bio),
    sqlc.arg(is_phone_number_public),
    sqlc.narg(phone_number),
    sqlc.arg(is_address_public),
    sqlc.narg(address),
    sqlc.arg(is_academic_institution_public),
    sqlc.narg(academic_institution),
    sqlc.arg(is_academic_email_public),
    sqlc.narg(academic_email)
) RETURNING *;

-- name: GetProfileByID :one
SELECT * FROM profiles 
WHERE id = sqlc.arg(id);

-- name: GetProfileByAuthCredentialID :one
SELECT * FROM profiles 
WHERE authentication_credential_id = sqlc.arg(authentication_credential_id);

-- name: GetProfileByEmail :one
-- Note: Searches encrypted email field directly (linear scan)
SELECT * FROM profiles 
WHERE email = sqlc.arg(email);

-- name: ListProfiles :many
SELECT * FROM profiles 
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: UpdateProfile :one
UPDATE profiles SET 
    is_profile_picture_public = COALESCE(sqlc.narg(is_profile_picture_public), is_profile_picture_public),
    profile_picture_url = COALESCE(sqlc.narg(profile_picture_url), profile_picture_url),
    is_first_name_public = COALESCE(sqlc.narg(is_first_name_public), is_first_name_public),
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    is_last_name_public = COALESCE(sqlc.narg(is_last_name_public), is_last_name_public),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    is_email_public = COALESCE(sqlc.narg(is_email_public), is_email_public),
    email = COALESCE(sqlc.narg(email), email),
    is_bio_public = COALESCE(sqlc.narg(is_bio_public), is_bio_public),
    bio = COALESCE(sqlc.narg(bio), bio),
    is_phone_number_public = COALESCE(sqlc.narg(is_phone_number_public), is_phone_number_public),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    is_address_public = COALESCE(sqlc.narg(is_address_public), is_address_public),
    address = COALESCE(sqlc.narg(address), address),
    is_academic_institution_public = COALESCE(sqlc.narg(is_academic_institution_public), is_academic_institution_public),
    academic_institution = COALESCE(sqlc.narg(academic_institution), academic_institution),
    is_academic_email_public = COALESCE(sqlc.narg(is_academic_email_public), is_academic_email_public),
    academic_email = COALESCE(sqlc.narg(academic_email), academic_email),
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING *;

-- name: UpdateProfileByAuthenticationCredentialId :one
UPDATE profiles SET 
    is_profile_picture_public = COALESCE(sqlc.narg(is_profile_picture_public), is_profile_picture_public),
    profile_picture_url = COALESCE(sqlc.narg(profile_picture_url), profile_picture_url),
    is_first_name_public = COALESCE(sqlc.narg(is_first_name_public), is_first_name_public),
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    is_last_name_public = COALESCE(sqlc.narg(is_last_name_public), is_last_name_public),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    is_email_public = COALESCE(sqlc.narg(is_email_public), is_email_public),
    email = COALESCE(sqlc.narg(email), email),
    is_bio_public = COALESCE(sqlc.narg(is_bio_public), is_bio_public),
    bio = COALESCE(sqlc.narg(bio), bio),
    is_phone_number_public = COALESCE(sqlc.narg(is_phone_number_public), is_phone_number_public),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    is_address_public = COALESCE(sqlc.narg(is_address_public), is_address_public),
    address = COALESCE(sqlc.narg(address), address),
    is_academic_institution_public = COALESCE(sqlc.narg(is_academic_institution_public), is_academic_institution_public),
    academic_institution = COALESCE(sqlc.narg(academic_institution), academic_institution),
    is_academic_email_public = COALESCE(sqlc.narg(is_academic_email_public), is_academic_email_public),
    academic_email = COALESCE(sqlc.narg(academic_email), academic_email),
    updated_at = NOW()
WHERE authentication_credential_id = sqlc.arg(authentication_credential_id) 
RETURNING *;

-- name: CountProfiles :one
SELECT COUNT(*) FROM profiles;

-- name: DeleteProfile :exec
DELETE FROM profiles WHERE id = sqlc.arg(id);

-- name: DeleteProfileByAuthCredentialID :exec
DELETE FROM profiles WHERE authentication_credential_id = sqlc.arg(authentication_credential_id);

-- name: ListVerifiedIssuerProfiles :many
SELECT * FROM profiles 
INNER JOIN authentication_credentials ON profiles.authentication_credential_id = authentication_credentials.id
WHERE authentication_credentials.is_verified_issuer = 1
ORDER BY profiles.created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: ListIssuerProfiles :many
SELECT profiles.* FROM profiles 
INNER JOIN authentication_credentials ON profiles.authentication_credential_id = authentication_credentials.id
WHERE authentication_credentials.is_verified_issuer = 1
ORDER BY profiles.created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: GetProfileAndCredentialWithCredentialId :one
SELECT
 authentication_credentials.ID as authentication_credential_id,
 authentication_credentials.solution_status as solution_status,
 authentication_credentials.hashed_password as hashed_password,
 authentication_credentials.encrypted_private_key as encrypted_private_key,
 authentication_credentials.wallet_address as wallet_address,
 authentication_credentials.google_connector_ref as google_connector_ref,
 authentication_credentials.github_connector_ref as github_connector_ref,
 authentication_credentials.is_verified_organizer as is_verified_organizer,
 authentication_credentials.is_verified_issuer as is_verified_issuer,
 authentication_credentials.is_verified_student as is_verified_student,
 authentication_credentials.created_at as authentication_credential_created_at,
 authentication_credentials.updated_at as authentication_credential_updated_at,
 profiles.ID as profile_id,
 profiles.authentication_credential_id as profile_authentication_credential_id,
 profiles.is_profile_picture_public as profile_is_profile_picture_public,
 profiles.profile_picture_url as profile_profile_picture_url,
 profiles.is_first_name_public as profile_is_first_name_public,
 profiles.first_name as profile_first_name,
 profiles.is_last_name_public as profile_is_last_name_public,
 profiles.last_name as profile_last_name,
 profiles.is_email_public as profile_is_email_public,
 profiles.email as profile_email,
 profiles.is_bio_public as profile_is_bio_public,
 profiles.bio as profile_bio,
 profiles.is_phone_number_public as profile_is_phone_number_public,
 profiles.phone_number as profile_phone_number,
 profiles.is_address_public as profile_is_address_public,
 profiles.address as profile_address,
 profiles.is_academic_institution_public as profile_is_academic_institution_public,
 profiles.academic_institution as profile_academic_institution,
 profiles.is_academic_email_public as profile_is_academic_email_public,
 profiles.academic_email as profile_academic_email,
 profiles.created_at as profile_created_at,
 profiles.updated_at as profile_updated_at
 FROM profiles 
INNER JOIN authentication_credentials ON profiles.authentication_credential_id = authentication_credentials.id
WHERE profiles.authentication_credential_id = sqlc.arg(authentication_credential_id);