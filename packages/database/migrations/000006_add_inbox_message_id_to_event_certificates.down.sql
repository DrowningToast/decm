-- Remove inbox_message_id from event_certificates table

-- Drop index first
DROP INDEX IF EXISTS idx_event_certificates_inbox_message_id;

-- Drop the column
ALTER TABLE event_certificates DROP COLUMN IF EXISTS inbox_message_id;

