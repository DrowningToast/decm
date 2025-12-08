-- Add accepted_at field to event_registration_invitations table
-- NULL: Not accepted, timestamp: When invitation was accepted
ALTER TABLE event_registration_invitations
ADD COLUMN accepted_at TIMESTAMPTZ;

-- Create index for the accepted_at field for efficient queries
CREATE INDEX idx_event_registration_invitations_accepted_at ON event_registration_invitations(accepted_at);


