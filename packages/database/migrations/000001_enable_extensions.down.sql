-- Remove PostgreSQL extensions for DECM platform
-- NOTE: Be careful when dropping extensions in production as they may be used by other databases

-- Remove comments first
COMMENT ON EXTENSION "uuid-ossp" IS NULL;
COMMENT ON EXTENSION "pg_trgm" IS NULL;
COMMENT ON EXTENSION "citext" IS NULL;

-- Drop extensions (only if no dependencies exist)
-- Note: These commands will fail if tables/functions depend on these extensions
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "uuid-ossp";
