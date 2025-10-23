-- DECM Event Config Tables Rollback Migration

-- Drop triggers in reverse order of creation
DROP TRIGGER IF EXISTS update_event_contracts_updated_at ON event_contracts;
DROP TRIGGER IF EXISTS update_event_issuers_updated_at ON event_issuers;
DROP TRIGGER IF EXISTS update_event_certificate_configs_updated_at ON event_certificate_configs;
DROP TRIGGER IF EXISTS update_event_registration_configs_updated_at ON event_registration_configs;

-- Drop indexes
DROP INDEX IF EXISTS idx_event_contracts_id;
DROP INDEX IF EXISTS idx_event_contracts_event_id;
DROP INDEX IF EXISTS idx_event_issuers_id;
DROP INDEX IF EXISTS idx_event_issuers_issuer_credential_id;
DROP INDEX IF EXISTS idx_event_issuers_event_id;
DROP INDEX IF EXISTS idx_event_certificate_configs_id;
DROP INDEX IF EXISTS idx_event_certificate_configs_event_id;
DROP INDEX IF EXISTS idx_event_registration_configs_id;
DROP INDEX IF EXISTS idx_event_registration_configs_event_id;

-- Drop tables in reverse order of creation (to respect foreign key constraints)
DROP TABLE IF EXISTS event_contracts CASCADE;
DROP TABLE IF EXISTS event_issuers CASCADE;
DROP TABLE IF EXISTS event_certificate_configs CASCADE;
DROP TABLE IF EXISTS event_registration_configs CASCADE;