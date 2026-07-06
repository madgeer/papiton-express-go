CREATE TYPE movement_type AS ENUM (
    'WAREHOUSE_IN',
    'WAREHOUSE_OUT'
);

CREATE TABLE warehouses (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE warehouse_routes (
    id UUID PRIMARY KEY,
    origin_warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
    destination_city VARCHAR(100) NOT NULL,
    next_warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_origin_dest UNIQUE (origin_warehouse_id, destination_city)
);

CREATE TABLE warehouse_movements (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    warehouse_id UUID REFERENCES warehouses(id) ON DELETE CASCADE,
    type movement_type NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE warehouse_outbox (
    id UUID PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_movements_order_id ON warehouse_movements(order_id);
CREATE INDEX idx_movements_warehouse_id ON warehouse_movements(warehouse_id);
