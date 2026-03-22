CREATE TABLE certificate_share (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_certificate_id UUID NOT NULL REFERENCES event_certificates(id) ON DELETE CASCADE,
    active BOOLEAN NOT NULL DEFAULT false,

    handle TEXT NOT NULL,
    password TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_certificate_share_event_certificate_id ON certificate_share(event_certificate_id);
CREATE INDEX idx_certificate_share_handle ON certificate_share(handle);
CREATE INDEX idx_certificate_share_active ON certificate_share(active);