-- DECM Core Tables Migration
-- Create core tables for the DECM platform
-- Prerequisites: Extensions must be enabled (see 000001_enable_extensions.up.sql)

-- authentication credentials table, index by id and email
CREATE TABLE authentication_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- 0: Bring your own private key
    -- 1: System managed encrypted private key
    solution_status INTEGER NOT NULL,

    -- hashed of hashed password using Argon2id
    hashed_password VARCHAR(255),
    -- encrypted by a hash of password using AES-256-GCM
    encrypted_private_key BYTEA,
    wallet_address VARCHAR(255) NOT NULL UNIQUE,

    -- PII: OAuth connector references (encrypted with AES-GCM at repository layer)
    google_connector_ref TEXT,
    -- PII: OAuth connector references (encrypted)
    github_connector_ref TEXT,
    
    is_verified_organizer INTEGER NOT NULL,
    is_verified_issuer INTEGER NOT NULL,
    is_verified_student INTEGER NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_authentication_credentials_id ON authentication_credentials(id);
CREATE INDEX idx_authentication_credentials_wallet_address ON authentication_credentials(wallet_address);

-- Student profile table, index by id and email
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    authentication_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,

    is_profile_picture_public INTEGER NOT NULL,
    -- PII: Profile picture URL (encrypted with AES-GCM at repository layer)
    profile_picture_url TEXT,

    is_first_name_public INTEGER NOT NULL,
    -- PII: First name (encrypted with AES-GCM at repository layer)
    first_name TEXT,
    is_last_name_public INTEGER NOT NULL,
    -- PII: Last name (encrypted with AES-GCM at repository layer)
    last_name TEXT,
    is_email_public INTEGER NOT NULL,
    -- PII: Email (encrypted with AES-GCM at repository layer)
    email TEXT,

    is_bio_public INTEGER NOT NULL,
    -- PII: Bio (encrypted with AES-GCM at repository layer)
    bio TEXT,
    is_phone_number_public INTEGER NOT NULL,
    -- PII: Phone number (encrypted with AES-GCM at repository layer)
    phone_number TEXT,
    is_address_public INTEGER NOT NULL,
    -- PII: Address (encrypted with AES-GCM at repository layer)
    address TEXT,
    is_academic_institution_public INTEGER NOT NULL,
    -- PII: Academic institution (encrypted with AES-GCM at repository layer)
    academic_institution TEXT,
    is_academic_email_public INTEGER NOT NULL,
    -- PII: Academic email (encrypted with AES-GCM at repository layer)
    academic_email TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_profiles_authentication_credential_id ON profiles(authentication_credential_id);
CREATE INDEX idx_profiles_id ON profiles(id);

-- events table, index by id and owner_credential_id
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    chain_id INTEGER NOT NULL,
    contact_number VARCHAR(255) NOT NULL,
    contact_address VARCHAR(255) NOT NULL,
    owner_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,

    banner_storage_key VARCHAR(255) NOT NULL,
    icon_storage_key VARCHAR(255) NOT NULL,

    title VARCHAR(255) NOT NULL,
    short_description VARCHAR(255) NOT NULL,
    long_description TEXT,

    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,

    location VARCHAR(255) NOT NULL,
    google_map_query VARCHAR(255) NOT NULL,
    max_attendees INTEGER NOT NULL,

    is_public INTEGER DEFAULT 0,
    is_booking_request_required INTEGER DEFAULT 0,
    is_verified INTEGER DEFAULT 0,
    is_ticket_transferable INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_events_owner_credential_id ON events(owner_credential_id);
CREATE INDEX idx_events_id ON events(id);

-- event_attendees table, index by event_id and credential_id
CREATE TABLE event_attendees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    attendee_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    
    contact_address VARCHAR(255) NOT NULL,

    is_attendee_accepted INTEGER NOT NULL,
    
    -- 0: Not Provided, 1: Provided
    -- PII: First name (encrypted with AES-GCM at repository layer)
    first_name TEXT,
    -- PII: Last name (encrypted with AES-GCM at repository layer)
    last_name TEXT,
    -- PII: Email (encrypted with AES-GCM at repository layer)
    email TEXT,
    -- PII: Bio (encrypted with AES-GCM at repository layer)
    bio TEXT,
    -- PII: Phone number (encrypted with AES-GCM at repository layer)
    phone_number TEXT,
    -- PII: Address (encrypted with AES-GCM at repository layer)
    address TEXT,
    -- PII: Academic institution (encrypted with AES-GCM at repository layer)
    academic_institution TEXT,
    -- PII: Academic email (encrypted with AES-GCM at repository layer)
    academic_email TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_attendees_event_id ON event_attendees(event_id);
CREATE INDEX idx_event_attendees_attendee_credential_id ON event_attendees(attendee_credential_id);

-- event_certificates table, index by event_id and credential_id
-- CREATE TABLE event_certificates (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
--     credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    
--     -- TODO: Add an entity for handling the issuance of the certificate
    
--     -- 0: Not published, 1: Published
--     is_published INTEGER NOT NULL,
    
--     created_at TIMESTAMPTZ DEFAULT NOW(),
--     updated_at TIMESTAMPTZ DEFAULT NOW()
-- );

-- CREATE INDEX idx_event_certificates_event_id ON event_certificates(event_id);
-- CREATE INDEX idx_event_certificates_credential_id ON event_certificates(credential_id);

-- Create function to update updated_at column automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic updated_at timestamp updates
CREATE TRIGGER update_authentication_credentials_updated_at 
    BEFORE UPDATE ON authentication_credentials 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profiles_updated_at 
    BEFORE UPDATE ON profiles 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_event_attendees_updated_at 
    BEFORE UPDATE ON event_attendees 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- CREATE TRIGGER update_event_certificates_updated_at 
--     BEFORE UPDATE ON event_certificates 
--     FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();