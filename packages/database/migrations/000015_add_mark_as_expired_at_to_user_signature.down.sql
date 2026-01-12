DROP INDEX IF EXISTS idx_user_signature_mark_as_expired_at;

ALTER TABLE user_signature
    DROP COLUMN IF EXISTS mark_as_expired_at;
