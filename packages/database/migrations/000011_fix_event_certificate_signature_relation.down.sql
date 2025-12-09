-- Revert Event Certificate Signature relation fix
-- Note: This migration cannot fully restore the original data relationships
-- as we cannot determine which specific certificate each signature belonged to

-- Step 1: Add back the old event_certificate_id column (nullable)
ALTER TABLE event_certificate_signatures
ADD COLUMN event_certificate_id UUID REFERENCES event_certificates(id) ON DELETE CASCADE;

-- Step 2: Attempt to restore data by picking the first certificate for each config
-- This is a best-effort restoration and may not match the original data exactly
UPDATE event_certificate_signatures ecs
SET event_certificate_id = (
    SELECT ec.id
    FROM event_certificates ec
    JOIN event_certificate_configs ecc ON ec.event_id = ecc.event_id
    WHERE ecc.id = ecs.event_certificate_config_id
    LIMIT 1
);

-- Step 3: Drop unique constraint
DROP INDEX IF EXISTS idx_event_certificate_signatures_config_issuer_unique;

-- Step 4: Drop index on config_id
DROP INDEX IF EXISTS idx_event_certificate_signatures_config_id;

-- Step 5: Drop foreign key constraint on config_id
ALTER TABLE event_certificate_signatures
DROP CONSTRAINT IF EXISTS event_certificate_signatures_event_certificate_config_id_fkey;

-- Step 6: Drop the config_id column
ALTER TABLE event_certificate_signatures
DROP COLUMN event_certificate_config_id;

-- Step 7: Make event_certificate_id NOT NULL (if we have data)
-- Note: This may fail if some signatures couldn't be matched to certificates
ALTER TABLE event_certificate_signatures
ALTER COLUMN event_certificate_id SET NOT NULL;
