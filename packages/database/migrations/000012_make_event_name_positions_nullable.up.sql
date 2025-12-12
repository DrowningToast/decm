-- Make event_name_pos_x and event_name_pos_y nullable in event_certificate_configs table
-- This allows certificates to be created without event name templates

ALTER TABLE event_certificate_configs
ALTER COLUMN event_name_pos_x DROP NOT NULL,
ALTER COLUMN event_name_pos_y DROP NOT NULL;



