-- inbox message type table, index by id
CREATE TABLE inbox_message_types (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_inbox_message_types_id ON inbox_message_types(id);

INSERT INTO inbox_message_types (id, name) VALUES (1, 'event_registration_invitation');

-- inbox message, index by id, sender_credential_id, receiver_credential_id
CREATE TABLE inbox_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    receiver_credential_id UUID REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    receiver_email TEXT NOT NULL,
    message_type INTEGER NOT NULL REFERENCES inbox_message_types(id) ON DELETE CASCADE,
    -- message_content is a JSONB object that contains translations of the message content
    message_content JSONB,
    -- fallback message content is a JSONB object that contains the fallback message content
    fallback_message_content TEXT,
    is_read INTEGER DEFAULT 0,  
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    hidden_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- trusted identities table, index by id
CREATE TABLE trusted_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,

    expected_identifiers TEXT[] NOT NULL,
    expected_issuer_addresses TEXT[] NOT NULL,

    -- 0: Not Self-Issued, 1: Self-Issued
    -- self-issued means the identity is issued by this system itself
    is_self_issued INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_trusted_identities_id ON trusted_identities(id);
CREATE INDEX idx_trusted_identities_name ON trusted_identities(name);

-- event registration requirement type table, index by id
CREATE TABLE event_registration_requirement_types (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_registration_requirement_types_id ON event_registration_requirement_types(id);

INSERT INTO event_registration_requirement_types (id, name) VALUES (1, 'password-protected');
INSERT INTO event_registration_requirement_types (id, name) VALUES (2, 'invite-only');

-- event registration requirement table, index by event_id
CREATE TABLE event_registration_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    requirement_type_id INTEGER NOT NULL REFERENCES event_registration_requirement_types(id) ON DELETE CASCADE,

    password TEXT, 
    max_attempts_in_period INTEGER DEFAULT 3,
    period_duration_seconds INTEGER DEFAULT 3600,

    is_identity_verification_required INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_registration_requirements_event_id ON event_registration_requirements(event_id);
CREATE INDEX idx_event_registration_requirements_id ON event_registration_requirements(id);

-- event registration requirement table <-> mapped with trusted_identities table
CREATE TABLE event_registration_requirements_trusted_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_registration_requirement_id UUID NOT NULL REFERENCES event_registration_requirements(id) ON DELETE CASCADE,
    trusted_identity_id UUID NOT NULL REFERENCES trusted_identities(id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_registration_requirements_trusted_identities_event_registration_requirement_id ON event_registration_requirements_trusted_identities(event_registration_requirement_id);
CREATE INDEX idx_event_registration_requirements_trusted_identities_trusted_identity_id ON event_registration_requirements_trusted_identities(trusted_identity_id);

-- event registration invitation
CREATE TABLE event_registration_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    inbox_message_id UUID NOT NULL REFERENCES inbox_messages(id) ON DELETE CASCADE,
    valid_until TIMESTAMPTZ,
    -- ASCI Case Insensitive 
    code TEXT, 

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ
);

CREATE INDEX idx_event_registration_invitations_event_id ON event_registration_invitations(event_id);
CREATE INDEX idx_event_registration_invitations_id ON event_registration_invitations(id);

-- event registration email referral table, index by id, event_id
CREATE TABLE event_registration_email_referrals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    
    valid_until TIMESTAMPTZ,
    is_used INTEGER DEFAULT 0,

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

CREATE INDEX idx_event_registration_email_referrals_id ON event_registration_email_referrals(id);
CREATE INDEX idx_event_registration_email_referrals_event_id ON event_registration_email_referrals(event_id);

