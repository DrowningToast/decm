ALTER TABLE certificate_share
    DROP CONSTRAINT IF EXISTS uq_certificate_share_event_certificate_id;
