-- Profiles CRUD queries with PII encryption

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
    email_hash,
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
    CASE 
        WHEN sqlc.narg(profile_picture_url) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(profile_picture_url), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_first_name_public),
    CASE 
        WHEN sqlc.narg(first_name) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(first_name), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_last_name_public),
    CASE 
        WHEN sqlc.narg(last_name) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(last_name), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_email_public),
    CASE 
        WHEN sqlc.narg(email) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(email), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    CASE 
        WHEN sqlc.narg(email) IS NOT NULL 
        THEN encode(digest(sqlc.narg(email), 'sha256'), 'hex')
        ELSE NULL 
    END,
    sqlc.arg(is_bio_public),
    CASE 
        WHEN sqlc.narg(bio) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(bio), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_phone_number_public),
    CASE 
        WHEN sqlc.narg(phone_number) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(phone_number), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_address_public),
    CASE 
        WHEN sqlc.narg(address) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(address), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_academic_institution_public),
    CASE 
        WHEN sqlc.narg(academic_institution) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(academic_institution), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text,
    sqlc.arg(is_academic_email_public),
    CASE 
        WHEN sqlc.narg(academic_email) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(academic_email), sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text
) RETURNING 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at;

-- name: GetProfileByID :one
SELECT 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at
FROM profiles 
WHERE id = sqlc.arg(id);

-- name: GetProfileByAuthCredentialID :one
SELECT 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at
FROM profiles 
WHERE authentication_credential_id = sqlc.arg(authentication_credential_id);

-- name: GetProfileByEmail :one
-- Note: Now uses email_hash for efficient searching instead of decrypting every row
SELECT 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at
FROM profiles 
WHERE email_hash = encode(digest(sqlc.arg(email_search), 'sha256'), 'hex');

-- name: ListProfiles :many
SELECT 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at
FROM profiles 
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: UpdateProfile :one
UPDATE profiles SET 
    is_profile_picture_public = COALESCE(sqlc.narg(is_profile_picture_public), is_profile_picture_public),
    profile_picture_url = CASE 
        WHEN sqlc.narg(profile_picture_url) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(profile_picture_url), sqlc.arg(encryption_key)::varchar)
        ELSE profile_picture_url
    END::text,
    is_first_name_public = COALESCE(sqlc.narg(is_first_name_public), is_first_name_public),
    first_name = CASE 
        WHEN sqlc.narg(first_name) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(first_name), sqlc.arg(encryption_key)::varchar)
        ELSE first_name
    END::text,
    is_last_name_public = COALESCE(sqlc.narg(is_last_name_public), is_last_name_public),
    last_name = CASE 
        WHEN sqlc.narg(last_name) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(last_name), sqlc.arg(encryption_key)::varchar)
        ELSE last_name
    END::text,
    is_email_public = COALESCE(sqlc.narg(is_email_public), is_email_public),
    email = CASE 
        WHEN sqlc.narg(email) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(email), sqlc.arg(encryption_key)::varchar)
        ELSE email
    END::text,
    email_hash = CASE 
        WHEN sqlc.narg(email) IS NOT NULL 
        THEN encode(digest(sqlc.narg(email), 'sha256'), 'hex')
        ELSE email_hash
    END,
    is_bio_public = COALESCE(sqlc.narg(is_bio_public), is_bio_public),
    bio = CASE 
        WHEN sqlc.narg(bio) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(bio), sqlc.arg(encryption_key)::varchar)
        ELSE bio
    END::text,
    is_phone_number_public = COALESCE(sqlc.narg(is_phone_number_public), is_phone_number_public),
    phone_number = CASE 
        WHEN sqlc.narg(phone_number) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(phone_number), sqlc.arg(encryption_key)::varchar)
        ELSE phone_number
    END::text,
    is_address_public = COALESCE(sqlc.narg(is_address_public), is_address_public),
    address = CASE 
        WHEN sqlc.narg(address) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(address), sqlc.arg(encryption_key)::varchar)
        ELSE address
    END::text,
    is_academic_institution_public = COALESCE(sqlc.narg(is_academic_institution_public), is_academic_institution_public),
    academic_institution = CASE 
        WHEN sqlc.narg(academic_institution) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(academic_institution), sqlc.arg(encryption_key)::varchar)
        ELSE academic_institution
    END::text,
    is_academic_email_public = COALESCE(sqlc.narg(is_academic_email_public), is_academic_email_public),
    academic_email = CASE 
        WHEN sqlc.narg(academic_email) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(academic_email), sqlc.arg(encryption_key)::varchar)
        ELSE academic_email
    END::text,
    updated_at = NOW()
WHERE id = sqlc.arg(id) 
RETURNING 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at;

-- name: UpdateProfileByAuthenticationCredentialId :one
UPDATE profiles SET 
    is_profile_picture_public = COALESCE(sqlc.narg(is_profile_picture_public), is_profile_picture_public),
    profile_picture_url = CASE 
        WHEN sqlc.narg(profile_picture_url) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(profile_picture_url), sqlc.arg(encryption_key)::varchar)
        ELSE profile_picture_url
    END::text,
    is_first_name_public = COALESCE(sqlc.narg(is_first_name_public), is_first_name_public),
    first_name = CASE 
        WHEN sqlc.narg(first_name) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(first_name), sqlc.arg(encryption_key)::varchar)
        ELSE first_name
    END::text,
    is_last_name_public = COALESCE(sqlc.narg(is_last_name_public), is_last_name_public),
    last_name = CASE 
        WHEN sqlc.narg(last_name) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(last_name), sqlc.arg(encryption_key)::varchar)
        ELSE last_name
    END::text,
    is_email_public = COALESCE(sqlc.narg(is_email_public), is_email_public),
    email = CASE 
        WHEN sqlc.narg(email) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(email), sqlc.arg(encryption_key)::varchar)
        ELSE email
    END::text,
    email_hash = CASE 
        WHEN sqlc.narg(email) IS NOT NULL 
        THEN encode(digest(sqlc.narg(email), 'sha256'), 'hex')
        ELSE email_hash
    END,
    is_bio_public = COALESCE(sqlc.narg(is_bio_public), is_bio_public),
    bio = CASE 
        WHEN sqlc.narg(bio) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(bio), sqlc.arg(encryption_key)::varchar)
        ELSE bio
    END::text,
    is_phone_number_public = COALESCE(sqlc.narg(is_phone_number_public), is_phone_number_public),
    phone_number = CASE 
        WHEN sqlc.narg(phone_number) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(phone_number), sqlc.arg(encryption_key)::varchar)
        ELSE phone_number
    END::text,
    is_address_public = COALESCE(sqlc.narg(is_address_public), is_address_public),
    address = CASE 
        WHEN sqlc.narg(address) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(address), sqlc.arg(encryption_key)::varchar)
        ELSE address
    END::text,
    is_academic_institution_public = COALESCE(sqlc.narg(is_academic_institution_public), is_academic_institution_public),
    academic_institution = CASE 
        WHEN sqlc.narg(academic_institution) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(academic_institution), sqlc.arg(encryption_key)::varchar)
        ELSE academic_institution
    END::text,
    is_academic_email_public = COALESCE(sqlc.narg(is_academic_email_public), is_academic_email_public),
    academic_email = CASE 
        WHEN sqlc.narg(academic_email) IS NOT NULL 
        THEN pgp_sym_encrypt(sqlc.narg(academic_email), sqlc.arg(encryption_key)::varchar)
        ELSE academic_email
    END::text,
    updated_at = NOW()
WHERE authentication_credential_id = sqlc.arg(authentication_credential_id) 
RETURNING 
    id,
    authentication_credential_id,
    is_profile_picture_public,
    CASE 
        WHEN profile_picture_url IS NOT NULL 
        THEN pgp_sym_decrypt(profile_picture_url::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as profile_picture_url,
    is_first_name_public,
    CASE 
        WHEN first_name IS NOT NULL 
        THEN pgp_sym_decrypt(first_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as first_name,
    is_last_name_public,
    CASE 
        WHEN last_name IS NOT NULL 
        THEN pgp_sym_decrypt(last_name::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as last_name,
    is_email_public,
    CASE 
        WHEN email IS NOT NULL 
        THEN pgp_sym_decrypt(email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as email,
    is_bio_public,
    CASE 
        WHEN bio IS NOT NULL 
        THEN pgp_sym_decrypt(bio::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as bio,
    is_phone_number_public,
    CASE 
        WHEN phone_number IS NOT NULL 
        THEN pgp_sym_decrypt(phone_number::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as phone_number,
    is_address_public,
    CASE 
        WHEN address IS NOT NULL 
        THEN pgp_sym_decrypt(address::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as address,
    is_academic_institution_public,
    CASE 
        WHEN academic_institution IS NOT NULL 
        THEN pgp_sym_decrypt(academic_institution::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_institution,
    is_academic_email_public,
    CASE 
        WHEN academic_email IS NOT NULL 
        THEN pgp_sym_decrypt(academic_email::bytea, sqlc.arg(encryption_key)::varchar)
        ELSE NULL 
    END::text as academic_email,
    created_at,
    updated_at;

-- name: CountProfiles :one
SELECT COUNT(*) FROM profiles;

-- name: DeleteProfile :exec
DELETE FROM profiles WHERE id = sqlc.arg(id);

-- name: DeleteProfileByAuthCredentialID :exec
DELETE FROM profiles WHERE authentication_credential_id = sqlc.arg(authentication_credential_id);
