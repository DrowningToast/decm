-- Add mark_as_expired_at column to user_signature table
-- This column is used by the background worker to mark signatures whose
-- blockchain deadline has passed, preventing infinite retry loops.
-- Distinct from broadcasted_at: a signature can be expired without being
-- successfully broadcasted to the blockchain.
ALTER TABLE user_signature
    ADD COLUMN mark_as_expired_at TIMESTAMPTZ;

COMMENT ON COLUMN user_signature.mark_as_expired_at IS 'Timestamp when the worker marked this signature as expired (deadline passed without successful broadcast). NULL means not expired.';

CREATE INDEX idx_user_signature_mark_as_expired_at
    ON user_signature(mark_as_expired_at)
    WHERE mark_as_expired_at IS NOT NULL;
