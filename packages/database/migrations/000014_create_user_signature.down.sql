-- Drop the foreign key columns from dependent tables first
ALTER TABLE event_attendees
DROP COLUMN IF EXISTS user_signature_id;

ALTER TABLE event_certificates
DROP COLUMN IF EXISTS user_claim_signature_id;

-- Drop user_signature table
DROP TABLE IF EXISTS user_signature;
