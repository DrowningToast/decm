-- Add font family and font weight columns for each text template in event_certificate_configs
ALTER TABLE event_certificate_configs
ADD COLUMN event_name_font_family VARCHAR(100),
ADD COLUMN event_name_font_weight INTEGER,
ADD COLUMN name_font_family VARCHAR(100),
ADD COLUMN name_font_weight INTEGER,
ADD COLUMN academic_institution_font_family VARCHAR(100),
ADD COLUMN academic_institution_font_weight INTEGER,
ADD COLUMN certificate_title_font_family VARCHAR(100),
ADD COLUMN certificate_title_font_weight INTEGER,
ADD COLUMN certificate_subtitle_font_family VARCHAR(100),
ADD COLUMN certificate_subtitle_font_weight INTEGER;







