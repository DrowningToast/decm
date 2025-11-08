-- ALTER TABLE event_certificate_configs
ALTER TABLE event_certificate_configs
ADD COLUMN certificate_title_pos_x FLOAT NOT NULL
ADD COLUMN certificate_title_pos_y FLOAT NOT NULL
ADD COLUMN certificate_subtitle_pos_x FLOAT NOT NULL
ADD COLUMN certificate_subtitle_pos_y FLOAT NOT NULL;

ALTER TABLE event_issuers
RENAME COLUMN sign_message TO sign_message_digest;
ALTER TABLE event_issuers
ADD COLUMN deleted_at TIMESTAMPTZ;

-- Create after issuer issued and host confirmed again.
CREATE TABLE event_certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    receiver_credential_id UUID REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    receiver_email TEXT,

    -- Attributes (Fill in credentials if available)
    -- PII encrypted
    name TEXT,
    -- PII encrypted
    academic_institution TEXT,
    certificate_title TEXT,
    certificate_subtitle TEXT,

    event_contract_address TEXT NOT NULL,
    event_certificate_address TEXT,
    certificate_token_id TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
)

-- Create at the same time of "event_certificates" row
CREATE TABLE event_certificate_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_certificate_id UUID NOT NULL REFERENCES event_certificates(id) ON DELETE CASCADE,
    issuer_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    -- signature that issuer has signed (of all participants) (same for all participants)
    issuer_signature TEXT,
    -- signature that host has signed (of all participants) (same for all participants)
    host_signature TEXT NOT NULL,
    sign_message_digest TEXT,
)