-- inbox message type table, index by id
CREATE TABLE inbox_message_types (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_inbox_message_types_id ON inbox_message_types(id);

INSERT INTO inbox_message_types (id, name) VALUES (0, 'general');
INSERT INTO inbox_message_types (id, name) VALUES (1, 'event_registration_invitation');
INSERT INTO inbox_message_types (id, name) VALUES (2, 'event_certificate_invitation');

-- inbox message, index by id, sender_credential_id, receiver_credential_id
CREATE TABLE inbox_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    sender_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    receiver_credential_id UUID REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    receiver_email TEXT,
    receiver_wallet_address TEXT,

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

-- event registration invitation
CREATE TABLE event_registration_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    inbox_message_id UUID NOT NULL REFERENCES inbox_messages(id) ON DELETE CASCADE,
    valid_until TIMESTAMPTZ,
    -- ASCI Case Insensitive 
    code TEXT, 

    -- PII: First name (encrypted with AES-GCM at repository layer)
    first_name TEXT,
    -- PII: Last name (encrypted with AES-GCM at repository layer)
    last_name TEXT,
    -- PII: Email (encrypted with AES-GCM at repository layer)
    email TEXT,
    -- PII: Phone number (encrypted with AES-GCM at repository layer)
    phone_number TEXT,
    -- PII: Academic institution (encrypted with AES-GCM at repository layer)
    academic_institution TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ
);

CREATE INDEX idx_event_registration_invitations_event_id ON event_registration_invitations(event_id);
CREATE INDEX idx_event_registration_invitations_id ON event_registration_invitations(id);
