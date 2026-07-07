CREATE TABLE trackings (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE,
    tracking_number VARCHAR(50) NOT NULL UNIQUE,
    current_status VARCHAR(50) NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tracking_events (
    id UUID PRIMARY KEY,
    tracking_id UUID REFERENCES trackings(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    location VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tracking_outbox (
    id UUID PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_trackings_order_id ON trackings(order_id);
CREATE INDEX idx_trackings_number ON trackings(tracking_number);
CREATE INDEX idx_events_tracking_id ON tracking_events(tracking_id);
