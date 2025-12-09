-- name: GetAllEventCertificateFontFamilies :many
SELECT * FROM event_certificate_font_families
WHERE deleted_at IS NULL
ORDER BY
    CASE WHEN is_default THEN 0 ELSE 1 END,
    font_family_name;

-- name: GetEventCertificateFontFamilyByID :one
SELECT * FROM event_certificate_font_families
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: GetEventCertificateFontFamilyByName :one
SELECT * FROM event_certificate_font_families
WHERE font_family_name = sqlc.arg('font_family_name') AND deleted_at IS NULL;

-- name: GetDefaultEventCertificateFontFamily :one
SELECT * FROM event_certificate_font_families
WHERE is_default = TRUE AND deleted_at IS NULL
LIMIT 1;

-- name: CreateEventCertificateFontFamily :one
INSERT INTO event_certificate_font_families (
    font_family_name,
    css_font_name,
    is_default,
    available_font_weights,
    is_support_italic
) VALUES (
    sqlc.arg('font_family_name'),
    sqlc.arg('css_font_name'),
    COALESCE(sqlc.arg('is_default'), FALSE),
    sqlc.arg('available_font_weights'),
    COALESCE(sqlc.arg('is_support_italic'), FALSE)
) RETURNING *;

-- name: UpdateEventCertificateFontFamily :one
UPDATE event_certificate_font_families
SET
    font_family_name = COALESCE(sqlc.narg('font_family_name'), font_family_name),
    css_font_name = COALESCE(sqlc.narg('css_font_name'), css_font_name),
    is_default = COALESCE(sqlc.narg('is_default'), is_default),
    available_font_weights = COALESCE(sqlc.narg('available_font_weights'), available_font_weights),
    is_support_italic = COALESCE(sqlc.narg('is_support_italic'), is_support_italic),
    updated_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteEventCertificateFontFamily :exec
UPDATE event_certificate_font_families
SET deleted_at = NOW()
WHERE id = sqlc.arg('id') AND deleted_at IS NULL;

-- name: DeleteEventCertificateFontFamily :exec
DELETE FROM event_certificate_font_families
WHERE id = sqlc.arg('id');
