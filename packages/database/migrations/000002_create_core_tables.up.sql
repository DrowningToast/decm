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

    -- PII: OAuth connector references (encrypted)
    google_connector_ref BYTEA,
    -- Hash for efficient searching of google_connector_ref
    google_connector_ref_hash VARCHAR(64),
    -- PII: OAuth connector references (encrypted)
    github_connector_ref BYTEA,
    -- Hash for efficient searching of github_connector_ref
    github_connector_ref_hash VARCHAR(64),

    is_verified_organizer INTEGER NOT NULL,
    is_verified_student INTEGER NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_authentication_credentials_id ON authentication_credentials(id);
CREATE INDEX idx_authentication_credentials_wallet_address ON authentication_credentials(wallet_address);
CREATE INDEX idx_authentication_credentials_google_hash ON authentication_credentials(google_connector_ref_hash);
CREATE INDEX idx_authentication_credentials_github_hash ON authentication_credentials(github_connector_ref_hash);

-- Student profile table, index by id and email
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    authentication_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,

    is_profile_picture_public INTEGER NOT NULL,
    -- PII: Profile picture URL
    profile_picture_url BYTEA,

    is_first_name_public INTEGER NOT NULL,
    -- PII: First name
    first_name BYTEA,
    is_last_name_public INTEGER NOT NULL,
    -- PII: Last name
    last_name BYTEA,
    is_email_public INTEGER NOT NULL,
    -- PII: Email (encrypted, cannot have UNIQUE constraint due to salt)
    email BYTEA,
    -- Hash for efficient searching and uniqueness constraint on email
    email_hash VARCHAR(64) UNIQUE,

    is_bio_public INTEGER NOT NULL,
    -- PII: Bio
    bio BYTEA,
    is_phone_number_public INTEGER NOT NULL,
    -- PII: Phone number
    phone_number BYTEA,
    is_address_public INTEGER NOT NULL,
    -- PII: Address
    address BYTEA,
    is_academic_institution_public INTEGER NOT NULL,
    -- PII: Academic institution
    academic_institution BYTEA,
    is_academic_email_public INTEGER NOT NULL,
    -- PII: Academic email
    academic_email BYTEA,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_profiles_authentication_credential_id ON profiles(authentication_credential_id);
CREATE INDEX idx_profiles_id ON profiles(id);
CREATE INDEX idx_profiles_email_hash ON profiles(email_hash);

-- events table, index by id and owner_credential_id
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    chain_id INTEGER NOT NULL,
    contact_number VARCHAR(255) NOT NULL,
    contact_address VARCHAR(255) NOT NULL,
    owner_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,

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

    -- 0: Not Required, 1: Required, 2: Optional
    first_name_requirement_status INTEGER NOT NULL,
    last_name_requirement_status INTEGER NOT NULL,
    email_requirement_status INTEGER NOT NULL,
    bio_requirement_status INTEGER NOT NULL,
    phone_number_requirement_status INTEGER NOT NULL,
    address_requirement_status INTEGER NOT NULL,
    academic_institution_requirement_status INTEGER NOT NULL,
    academic_email_requirement_status INTEGER NOT NULL,

    -- base event certificate image references
    base_event_certificate_url VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
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
    -- PII: First name
    first_name BYTEA,
    -- PII: Last name
    last_name BYTEA,
    -- PII: Email
    email BYTEA,
    -- PII: Bio
    bio BYTEA,
    -- PII: Phone number    
    phone_number BYTEA,
    -- PII: Address
    address BYTEA,
    -- PII: Academic institution
    academic_institution BYTEA,
    -- PII: Academic email
    academic_email BYTEA,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_attendees_event_id ON event_attendees(event_id);
CREATE INDEX idx_event_attendees_attendee_credential_id ON event_attendees(attendee_credential_id);

-- event_certificates table, index by event_id and credential_id
CREATE TABLE event_certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    
    -- TODO: Add an entity for handling the issuance of the certificate
    
    -- 0: Not published, 1: Published
    is_published INTEGER NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_certificates_event_id ON event_certificates(event_id);
CREATE INDEX idx_event_certificates_credential_id ON event_certificates(credential_id);

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

CREATE TRIGGER update_event_certificates_updated_at 
    BEFORE UPDATE ON event_certificates 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();