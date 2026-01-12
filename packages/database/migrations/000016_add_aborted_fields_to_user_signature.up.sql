-- Add aborted_at and aborted_reason columns to user_signature table.
-- When the blockchain submission worker determines that a queued signature
-- can never succeed (e.g. the participant never joined the event before
-- claiming a certificate), it marks the record as aborted rather than
-- deleting it, so the reason is preserved for auditing.
ALTER TABLE user_signature
    ADD COLUMN aborted_at TIMESTAMPTZ,
    ADD COLUMN aborted_reason SMALLINT;

COMMENT ON COLUMN user_signature.aborted_at IS 'Timestamp when the worker permanently aborted this signature (e.g. participant never joined). NULL means not aborted.';
COMMENT ON COLUMN user_signature.aborted_reason IS 'Integer reason code for the abort (see UserSignatureAbortReason in Go code). NULL if not aborted.';

CREATE INDEX idx_user_signature_aborted_at
    ON user_signature(aborted_at)
    WHERE aborted_at IS NOT NULL;
