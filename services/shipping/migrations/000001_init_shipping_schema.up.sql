CREATE TYPE courier_status AS ENUM (
    'AVAILABLE',
    'ON_DELIVERY',
    'OFFLINE'
);

CREATE TYPE shipment_status AS ENUM (
    'PENDING_PICKUP',
    'PICKED_UP',
    'ON_TRANSIT',
    'OUT_FOR_DELIVERY',
    'DELIVERED',
    'FAILED'
);

CREATE TABLE couriers (
    id         UUID PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    phone      VARCHAR(20) NOT NULL,
    status     courier_status NOT NULL DEFAULT 'AVAILABLE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shipments (
    id         UUID PRIMARY KEY,
    order_id   UUID NOT NULL,
    courier_id UUID REFERENCES couriers(id) ON DELETE SET NULL,
    status     shipment_status NOT NULL DEFAULT 'PENDING_PICKUP',
    notes      TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shipping_outbox (
    id             UUID PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id   UUID NOT NULL,
    event_type     VARCHAR(100) NOT NULL,
    payload        JSONB NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);