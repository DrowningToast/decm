-- Enable required PostgreSQL extensions for DECM platform

-- UUID generation for unique IDs (required for all entity primary keys)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Full-text search capabilities (required for user and event search)
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Case-insensitive text operations (required for emails and usernames)
CREATE EXTENSION IF NOT EXISTS "citext";

-- JSON operations for array handling (PostgreSQL < 14 compatibility)
-- CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Add helpful comments for future reference
COMMENT ON EXTENSION "uuid-ossp" IS 'UUID generation functions for primary keys and unique identifiers';
COMMENT ON EXTENSION "pg_trgm" IS 'Trigram matching for full-text search on users and events';
COMMENT ON EXTENSION "citext" IS 'Case-insensitive text data type for emails and usernames';

-- Note: PII encryption is done at application layer using AES-GCM in Go code, not at database level
