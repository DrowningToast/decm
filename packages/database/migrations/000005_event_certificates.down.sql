-- Drop tables in reverse order of creation to respect foreign key constraints

-- Drop child table with dependencies first
DROP TABLE IF EXISTS event_certificate_signatures CASCADE;

-- Drop parent table second
DROP TABLE IF EXISTS event_certificates CASCADE;

-- Remove the column that was added to event_certificate_configs
ALTER TABLE event_certificate_configs DROP COLUMN IF EXISTS event_name_pos_x;
ALTER TABLE event_certificate_configs DROP COLUMN IF EXISTS certificate_title_pos_x;
ALTER TABLE event_certificate_configs DROP COLUMN IF EXISTS certificate_title_pos_y;
ALTER TABLE event_certificate_configs DROP COLUMN IF EXISTS certificate_subtitle_pos_x;
ALTER TABLE event_certificate_configs DROP COLUMN IF EXISTS certificate_subtitle_pos_y;

ALTER TABLE event_issuers DROP COLUMN IF EXISTS sign_message_digest;
ALTER TABLE event_issuers DROP COLUMN IF EXISTS deleted_at;