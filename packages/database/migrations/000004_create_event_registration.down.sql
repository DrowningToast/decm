-- Drop tables in reverse order to respect foreign key constraints

-- Drop event registration email referrals table and indexes
DROP INDEX IF EXISTS idx_event_registration_email_referrals_event_id;
DROP INDEX IF EXISTS idx_event_registration_email_referrals_id;
DROP TABLE IF EXISTS event_registration_email_referrals CASCADE;

-- Drop event registration invitations table and indexes
DROP INDEX IF EXISTS idx_event_registration_invitations_id;
DROP INDEX IF EXISTS idx_event_registration_invitations_event_id;
DROP TABLE IF EXISTS event_registration_invitations CASCADE;

-- Drop event registration requirements trusted identities table and indexes
DROP INDEX IF EXISTS idx_event_registration_requirements_trusted_identities_trusted_identity_id;
DROP INDEX IF EXISTS idx_event_registration_requirements_trusted_identities_event_registration_requirement_id;
DROP TABLE IF EXISTS event_registration_requirements_trusted_identities CASCADE;

-- Drop event registration requirements table and indexes
DROP INDEX IF EXISTS idx_event_registration_requirements_id;
DROP INDEX IF EXISTS idx_event_registration_requirements_event_id;
DROP TABLE IF EXISTS event_registration_requirements CASCADE;

-- Drop event registration requirement types table and indexes
DROP INDEX IF EXISTS idx_event_registration_requirement_types_id;
DROP TABLE IF EXISTS event_registration_requirement_types CASCADE;

-- Drop trusted identities table and indexes
DROP INDEX IF EXISTS idx_trusted_identities_name;
DROP INDEX IF EXISTS idx_trusted_identities_id;
DROP TABLE IF EXISTS trusted_identities CASCADE;

-- Drop inbox messages table
DROP TABLE IF EXISTS inbox_messages CASCADE;

-- Drop inbox message types table and indexes
DROP INDEX IF EXISTS idx_inbox_message_types_id;
DROP TABLE IF EXISTS inbox_message_types CASCADE;

