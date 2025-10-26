-- DECM Event Config Tables Migration

CREATE TABLE event_registration_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    final_call_for_registration TIMESTAMPTZ,
    registration_password TEXT,

    -- 0: Not Required, 1: Required, 2: Optional
    first_name_requirement_status INTEGER DEFAULT 0 ,
    last_name_requirement_status INTEGER DEFAULT 0,
    email_requirement_status INTEGER DEFAULT 0,
    bio_requirement_status INTEGER DEFAULT 0,
    phone_number_requirement_status INTEGER DEFAULT 0,
    address_requirement_status INTEGER DEFAULT 0, 
    academic_institution_requirement_status INTEGER DEFAULT 0,
    academic_email_requirement_status INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_registration_configs_event_id ON event_registration_configs(event_id);
CREATE INDEX idx_event_registration_configs_id ON event_registration_configs(id);

CREATE TABLE event_certificate_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    -- base event certificate image references
    base_certificate_storage_key VARCHAR(255) NOT NULL,

    -- element positions
    event_name_pos_x FLOAT NOT NULL,
    event_name_pos_y FLOAT NOT NULL,
    name_pos_x FLOAT NOT NULL,
    name_pos_y FLOAT NOT NULL,
    academic_institution_pos_x FLOAT,
    academic_institution_pos_y FLOAT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_certificate_configs_event_id ON event_certificate_configs(event_id);
CREATE INDEX idx_event_certificate_configs_id ON event_certificate_configs(id);

CREATE TABLE event_issuers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    issuer_credential_id UUID NOT NULL REFERENCES authentication_credentials(id) ON DELETE CASCADE,
    -- 0: Not Signed, 1: Signed
    is_signed INTEGER NOT NULL DEFAULT 0,
    signature TEXT,
    sign_message TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_issuers_event_id ON event_issuers(event_id);
CREATE INDEX idx_event_issuers_issuer_credential_id ON event_issuers(issuer_credential_id);
CREATE INDEX idx_event_issuers_id ON event_issuers(id);

CREATE TABLE event_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    access_manager_contract_address VARCHAR(255) NOT NULL,
    event_contract_address VARCHAR(255) NOT NULL,
    ticket_contract_address VARCHAR(255),
    certificate_contract_address VARCHAR(255),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_contracts_event_id ON event_contracts(event_id);
CREATE INDEX idx_event_contracts_id ON event_contracts(id);

-- Create triggers for automatic updated_at timestamp updates
CREATE TRIGGER update_event_registration_configs_updated_at 
    BEFORE UPDATE ON event_registration_configs 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_event_certificate_configs_updated_at 
    BEFORE UPDATE ON event_certificate_configs 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_event_issuers_updated_at 
    BEFORE UPDATE ON event_issuers 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_event_contracts_updated_at 
    BEFORE UPDATE ON event_contracts 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
