-- Remove index
DROP INDEX IF EXISTS idx_event_certificate_configs_is_published;

-- Remove is_published column from event_certificate_configs table
ALTER TABLE event_certificate_configs
DROP COLUMN IF EXISTS is_published;



