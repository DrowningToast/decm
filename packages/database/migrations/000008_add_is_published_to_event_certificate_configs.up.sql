-- Add is_published column to event_certificate_configs table
ALTER TABLE event_certificate_configs
ADD COLUMN is_published BOOLEAN DEFAULT FALSE NOT NULL;

-- Add index for faster querying of published configs
CREATE INDEX idx_event_certificate_configs_is_published ON event_certificate_configs(is_published);




