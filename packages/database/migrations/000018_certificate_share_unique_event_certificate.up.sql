ALTER TABLE certificate_share
    ADD CONSTRAINT uq_certificate_share_event_certificate_id UNIQUE (event_certificate_id);
