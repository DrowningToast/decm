-- DECM Event Config Tables Rollback Migration

-- Drop triggers in reverse order of creation
DROP TRIGGER IF EXISTS update_event_contracts_updated_at ON event_contracts;
DROP TRIGGER IF EXISTS update_event_issuers_updated_at ON event_issuers;
DROP TRIGGER IF EXISTS update_event_certificate_configs_updated_at ON event_certificate_configs;
DROP TRIGGER IF EXISTS update_event_registration_configs_updated_at ON event_registration_configs;

-- Drop tables in reverse order of creation (to respect foreign key constraints)
DROP TABLE IF EXISTS event_contracts;
DROP TABLE IF EXISTS event_issuers;
DROP TABLE IF EXISTS event_certificate_configs;
DROP TABLE IF EXISTS event_registration_configs;