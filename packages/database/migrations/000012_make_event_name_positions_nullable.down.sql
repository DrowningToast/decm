-- Revert event_name_pos_x and event_name_pos_y to NOT NULL in event_certificate_configs table

-- Set any NULL values to 0,0 to satisfy NOT NULL constraint before reverting
UPDATE event_certificate_configs
SET event_name_pos_x = 0, event_name_pos_y = 0
WHERE event_name_pos_x IS NULL OR event_name_pos_y IS NULL;

-- Alter columns to require NOT NULL
ALTER TABLE event_certificate_configs
ALTER COLUMN event_name_pos_x SET NOT NULL,
ALTER COLUMN event_name_pos_y SET NOT NULL;



