-- DECM Core Tables Migration
-- Create core tables for the DECM platform
-- Prerequisites: Extensions must be enabled (see 000001_enable_extensions.up.sql)

-- authentication credentials table, index by id and email
CREATE TABLE authentication_credentials (
    id SERIAL PRIMARY KEY,
    -- 0: Bring your own private key
    -- 1: System managed encrypted private key
    solution_status INTEGER NOT NULL,

    -- hashed of hashed password using Argon2id
    hashed_password VARCHAR(255),
    -- encrypted by a hash of password using AES-256-GCM
    encrypted_private_key VARCHAR(255),
    public_key VARCHAR(255) NOT NULL UNIQUE,

    google_connector_ref VARCHAR(255),
    github_connector_ref VARCHAR(255),

    is_verified_organizer INTEGER NOT NULL,
    is_verified_student INTEGER NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_authentication_credentials_id ON authentication_credentials(id);
CREATE INDEX idx_authentication_credentials_public_key ON authentication_credentials(public_key);

-- Student profile table, index by id and email
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    authentication_credential_id INTEGER NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,

    is_profile_picture_public INTEGER NOT NULL,
    profile_picture_url VARCHAR(255),

    is_first_name_public INTEGER NOT NULL,
    first_name VARCHAR(128),
    is_last_name_public INTEGER NOT NULL,
    last_name VARCHAR(128),
    is_email_public INTEGER NOT NULL,
    email VARCHAR(255) UNIQUE,

    is_bio_public INTEGER NOT NULL,
    bio VARCHAR(255),
    is_phone_number_public INTEGER NOT NULL,
    phone_number VARCHAR(255),
    is_address_public INTEGER NOT NULL,
    address VARCHAR(255),
    is_academic_institution_public INTEGER NOT NULL,
    academic_institution VARCHAR(255),
    is_academic_email_public INTEGER NOT NULL,
    academic_email VARCHAR(255),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_profiles_authentication_credential_id ON profiles(authentication_credential_id);
CREATE INDEX idx_profiles_id ON profiles(id);

-- events table, index by id and owner_credential_id
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    
    chain_id INTEGER NOT NULL,
    contact_address VARCHAR(255) NOT NULL,
    owner_credential_id INTEGER NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,

    title VARCHAR(255) NOT NULL,
    short_description VARCHAR(255) NOT NULL,
    long_description TEXT,

    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,

    location VARCHAR(255) NOT NULL,
    google_map_query VARCHAR(255) NOT NULL,
    max_attendees INTEGER NOT NULL,

    is_public INTEGER NOT NULL,
    is_booking_request_required INTEGER NOT NULL,
    is_verified INTEGER NOT NULL,
    is_ticket_transferable INTEGER NOT NULL,

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
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    attendee_credential_id INTEGER NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    
    contact_address VARCHAR(255) NOT NULL,

    is_attendee_accepted INTEGER NOT NULL,
    
    -- 0: Not Provided, 1: Provided
    first_name_provided INTEGER NOT NULL,
    last_name_provided INTEGER NOT NULL,
    email_provided INTEGER NOT NULL,
    bio_provided INTEGER NOT NULL,
    phone_number_provided INTEGER NOT NULL,
    address_provided INTEGER NOT NULL,
    academic_institution_provided INTEGER NOT NULL,
    academic_email_provided INTEGER NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_attendees_event_id ON event_attendees(event_id);
CREATE INDEX idx_event_attendees_attendee_credential_id ON event_attendees(attendee_credential_id);

-- event_certificates table, index by event_id and credential_id
CREATE TABLE event_certificates (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    credential_id INTEGER NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    
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