-- Remove index
DROP INDEX IF EXISTS idx_event_registration_invitations_accepted_at;

-- Remove accepted_at column from event_registration_invitations table
ALTER TABLE event_registration_invitations
DROP COLUMN IF EXISTS accepted_at;


