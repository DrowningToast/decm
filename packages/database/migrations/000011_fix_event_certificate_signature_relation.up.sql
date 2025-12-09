-- Fix Event Certificate Signature relation
-- Change from linking to individual certificates to linking to certificate config
-- This makes sense because signatures are the same for all participants in an event

-- Step 1: Add new column for event_certificate_config_id (nullable initially)
ALTER TABLE event_certificate_signatures
ADD COLUMN event_certificate_config_id UUID REFERENCES event_certificate_configs(id) ON DELETE CASCADE;

-- Step 2: Migrate existing data by finding the config_id from the certificate's event_id
-- We join through event_certificates -> events -> event_certificate_configs
UPDATE event_certificate_signatures ecs
SET event_certificate_config_id = ecc.id
FROM event_certificates ec
JOIN event_certificate_configs ecc ON ec.event_id = ecc.event_id
WHERE ecs.event_certificate_id = ec.id;

-- Step 3: Make the new column NOT NULL (only if all rows were successfully migrated)
ALTER TABLE event_certificate_signatures
ALTER COLUMN event_certificate_config_id SET NOT NULL;

-- Step 4: Drop the old foreign key constraint and column
ALTER TABLE event_certificate_signatures
DROP CONSTRAINT IF EXISTS event_certificate_signatures_event_certificate_id_fkey;

ALTER TABLE event_certificate_signatures
DROP COLUMN event_certificate_id;

-- Step 5: Create index on the new foreign key column
CREATE INDEX idx_event_certificate_signatures_config_id ON event_certificate_signatures(event_certificate_config_id);

-- Step 6: Create unique constraint to ensure one signature per issuer per config
-- This prevents duplicate signatures for the same issuer/config combination
CREATE UNIQUE INDEX idx_event_certificate_signatures_config_issuer_unique 
ON event_certificate_signatures(event_certificate_config_id, issuer_credential_id);
