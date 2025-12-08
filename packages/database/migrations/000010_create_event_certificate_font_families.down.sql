-- Drop indexes on foreign key columns in event_certificate_configs
DROP INDEX IF EXISTS idx_event_certificate_configs_event_name_font_family_id;
DROP INDEX IF EXISTS idx_event_certificate_configs_name_font_family_id;
DROP INDEX IF EXISTS idx_event_certificate_configs_academic_institution_font_family_id;
DROP INDEX IF EXISTS idx_event_certificate_configs_certificate_title_font_family_id;
DROP INDEX IF EXISTS idx_event_certificate_configs_certificate_subtitle_font_family_id;

-- Re-add old VARCHAR font_family columns to event_certificate_configs
ALTER TABLE event_certificate_configs
ADD COLUMN event_name_font_family VARCHAR(100),
ADD COLUMN name_font_family VARCHAR(100),
ADD COLUMN academic_institution_font_family VARCHAR(100),
ADD COLUMN certificate_title_font_family VARCHAR(100),
ADD COLUMN certificate_subtitle_font_family VARCHAR(100);

-- Drop foreign key columns
ALTER TABLE event_certificate_configs
DROP COLUMN IF EXISTS event_name_font_family_id,
DROP COLUMN IF EXISTS name_font_family_id,
DROP COLUMN IF EXISTS academic_institution_font_family_id,
DROP COLUMN IF EXISTS certificate_title_font_family_id,
DROP COLUMN IF EXISTS certificate_subtitle_font_family_id;

-- Drop indexes on event_certificate_font_families
DROP INDEX IF EXISTS idx_event_certificate_font_families_deleted_at;
DROP INDEX IF EXISTS idx_event_certificate_font_families_is_default;

-- Drop event_certificate_font_families table
DROP TABLE IF EXISTS event_certificate_font_families;
