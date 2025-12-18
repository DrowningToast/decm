-- Create system_status_schedules table for tracking uptime and downtime
CREATE TABLE system_status_schedules (
    id SERIAL PRIMARY KEY,
    order_id SERIAL NOT NULL,
    start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    planned_end_time TIMESTAMPTZ,
    status INTEGER NOT NULL CHECK (status IN (0, 1)),
    is_planned BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Add comment on table
COMMENT ON TABLE system_status_schedules IS 'Tracks system uptime and downtime schedules';

-- Add comments on columns
COMMENT ON COLUMN system_status_schedules.status IS '0 = maintenance, 1 = operating';
COMMENT ON COLUMN system_status_schedules.is_planned IS 'Indicates if this is a planned maintenance/status change';
COMMENT ON COLUMN system_status_schedules.start_time IS 'When this status takes immediate effect';
COMMENT ON COLUMN system_status_schedules.planned_end_time IS 'Estimated end time for display purposes only';
COMMENT ON COLUMN system_status_schedules.order_id IS 'Sequential ID for business/display reference';

-- Create indexes for efficient querying
CREATE INDEX idx_system_status_schedules_start_time ON system_status_schedules(start_time);
CREATE INDEX idx_system_status_schedules_status ON system_status_schedules(status);
CREATE INDEX idx_system_status_schedules_deleted_at ON system_status_schedules(deleted_at);
CREATE UNIQUE INDEX idx_system_status_schedules_order_id ON system_status_schedules(order_id);
CREATE INDEX idx_system_status_schedules_updated_at ON system_status_schedules(updated_at);

-- Create index for finding current active status
CREATE INDEX idx_system_status_schedules_active ON system_status_schedules(start_time, deleted_at) WHERE deleted_at IS NULL;

-- Insert initial operating status (planned and available)
INSERT INTO system_status_schedules (status, is_planned)
SELECT 1, TRUE
WHERE NOT EXISTS (SELECT 1 FROM system_status_schedules);
