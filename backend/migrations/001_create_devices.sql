CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,

    device_id VARCHAR(100) NOT NULL UNIQUE,

    device_name VARCHAR(150) NOT NULL,

    serial_number VARCHAR(150) NOT NULL UNIQUE,

    os_version VARCHAR(100),

    ip_address INET,

    location VARCHAR(255),

    status VARCHAR(20) NOT NULL DEFAULT 'OFFLINE',

    last_seen TIMESTAMPTZ,

    last_online_at TIMESTAMPTZ,

    last_offline_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT devices_status_check CHECK (status IN ('ONLINE', 'OFFLINE'))
);

CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);

CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON devices(last_seen);