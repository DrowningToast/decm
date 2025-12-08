-- Create event_certificate_font_families table
CREATE TABLE event_certificate_font_families (
    id SERIAL PRIMARY KEY,
    font_family_name VARCHAR(100) NOT NULL UNIQUE,
    css_font_name VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    available_font_weights TEXT NOT NULL,
    is_support_italic BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Insert default font families
INSERT INTO event_certificate_font_families (font_family_name, css_font_name, is_default, available_font_weights, is_support_italic) VALUES
    ('Inter', 'Inter', TRUE, '100,200,300,400,500,600,700,800,900', TRUE),
    ('Noto Sans Thai', 'Noto Sans Thai', FALSE, '100,200,300,400,500,600,700,800,900', TRUE),
    ('Prompt', 'Prompt', FALSE, '100,200,300,400,500,600,700,800,900', TRUE),
    ('TH Sarabun New', 'Sarabun', FALSE, '100,200,300,400,500,600,700,800', TRUE),
    ('Kanit', 'Kanit', FALSE, '100,200,300,400,500,600,700,800,900', TRUE),
    ('Arial', 'Arial', FALSE, '400,700', TRUE),
    ('Tahoma', 'Tahoma', FALSE, '400,700', TRUE);

-- Create indexes
CREATE INDEX idx_event_certificate_font_families_deleted_at ON event_certificate_font_families(deleted_at);
CREATE INDEX idx_event_certificate_font_families_is_default ON event_certificate_font_families(is_default);

-- Alter event_certificate_configs table to add foreign key constraints
-- First, add new columns for font family IDs (nullable)
ALTER TABLE event_certificate_configs
ADD COLUMN event_name_font_family_id INTEGER REFERENCES event_certificate_font_families(id) ON DELETE SET NULL,
ADD COLUMN name_font_family_id INTEGER REFERENCES event_certificate_font_families(id) ON DELETE SET NULL,
ADD COLUMN academic_institution_font_family_id INTEGER REFERENCES event_certificate_font_families(id) ON DELETE SET NULL,
ADD COLUMN certificate_title_font_family_id INTEGER REFERENCES event_certificate_font_families(id) ON DELETE SET NULL,
ADD COLUMN certificate_subtitle_font_family_id INTEGER REFERENCES event_certificate_font_families(id) ON DELETE SET NULL;

-- Drop old font_family columns (they were VARCHAR)
ALTER TABLE event_certificate_configs
DROP COLUMN IF EXISTS event_name_font_family,
DROP COLUMN IF EXISTS name_font_family,
DROP COLUMN IF EXISTS academic_institution_font_family,
DROP COLUMN IF EXISTS certificate_title_font_family,
DROP COLUMN IF EXISTS certificate_subtitle_font_family;

-- Create indexes on the foreign key columns
CREATE INDEX idx_event_certificate_configs_event_name_font_family_id ON event_certificate_configs(event_name_font_family_id);
CREATE INDEX idx_event_certificate_configs_name_font_family_id ON event_certificate_configs(name_font_family_id);
CREATE INDEX idx_event_certificate_configs_academic_institution_font_family_id ON event_certificate_configs(academic_institution_font_family_id);
CREATE INDEX idx_event_certificate_configs_certificate_title_font_family_id ON event_certificate_configs(certificate_title_font_family_id);
CREATE INDEX idx_event_certificate_configs_certificate_subtitle_font_family_id ON event_certificate_configs(certificate_subtitle_font_family_id);
