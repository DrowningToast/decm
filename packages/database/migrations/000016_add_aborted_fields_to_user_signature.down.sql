DROP INDEX IF EXISTS idx_user_signature_aborted_at;

ALTER TABLE user_signature
    DROP COLUMN IF EXISTS aborted_at,
    DROP COLUMN IF EXISTS aborted_reason;
