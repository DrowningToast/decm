-- DECM Core Tables Migration Rollback
-- Drop all tables and related objects created in 000002_create_core_tables.up.sql

-- Drop triggers first
DROP TRIGGER IF EXISTS update_event_certificates_updated_at ON event_certificates;
DROP TRIGGER IF EXISTS update_event_attendees_updated_at ON event_attendees;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP TRIGGER IF EXISTS update_authentication_credentials_updated_at ON authentication_credentials;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes first (although they'll be dropped with the tables)
DROP INDEX IF EXISTS idx_event_certificates_credential_id;
DROP INDEX IF EXISTS idx_event_certificates_event_id;
DROP INDEX IF EXISTS idx_event_attendees_attendee_credential_id;
DROP INDEX IF EXISTS idx_event_attendees_event_id;
DROP INDEX IF EXISTS idx_events_id;
DROP INDEX IF EXISTS idx_events_owner_credential_id;
DROP INDEX IF EXISTS idx_profiles_id;
DROP INDEX IF EXISTS idx_profiles_authentication_credential_id;
DROP INDEX IF EXISTS idx_authentication_credentials_public_key;
DROP INDEX IF EXISTS idx_authentication_credentials_id;

-- Drop tables in reverse order of dependencies
-- (child tables first, then parent tables)

-- Drop event_certificates (depends on events and authentication_credentials)
DROP TABLE IF EXISTS event_certificates;

-- Drop event_attendees (depends on events and authentication_credentials)  
DROP TABLE IF EXISTS event_attendees;

-- Drop events (depends on authentication_credentials)
DROP TABLE IF EXISTS events;

-- Drop profiles (depends on authentication_credentials)
DROP TABLE IF EXISTS profiles;

-- Drop authentication_credentials (base table with no dependencies)
DROP TABLE IF EXISTS authentication_credentials;