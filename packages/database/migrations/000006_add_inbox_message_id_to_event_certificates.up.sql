-- Add inbox_message_id to event_certificates table
ALTER TABLE event_certificates
ADD COLUMN inbox_message_id UUID REFERENCES inbox_messages(id) ON DELETE SET NULL;

-- Create index for the foreign key
CREATE INDEX idx_event_certificates_inbox_message_id ON event_certificates(inbox_message_id);

