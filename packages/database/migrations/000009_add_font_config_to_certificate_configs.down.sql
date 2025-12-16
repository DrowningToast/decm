-- Remove font family and font weight columns from event_certificate_configs
ALTER TABLE event_certificate_configs
DROP COLUMN event_name_font_family,
DROP COLUMN event_name_font_weight,
DROP COLUMN name_font_family,
DROP COLUMN name_font_weight,
DROP COLUMN academic_institution_font_family,
DROP COLUMN academic_institution_font_weight,
DROP COLUMN certificate_title_font_family,
DROP COLUMN certificate_title_font_weight,
DROP COLUMN certificate_subtitle_font_family,
DROP COLUMN certificate_subtitle_font_weight;















