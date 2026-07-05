-- Create Custom Types for Enums
CREATE TYPE order_status AS ENUM (
    'WAITING_PAYMENT',
    'PAID',
    'COURIER_ASSIGNED',
    'PICKED_UP',
    'WAREHOUSE_IN',
    'IN_DELIVERY',
    'DELIVERED',
    'COMPLETED'
);

CREATE TYPE address_type AS ENUM (
    'SENDER',
    'RECEIVER'
);

-- Table: service_types
CREATE TABLE service_types (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    estimated_day INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Table: tariffs
CREATE TABLE tariffs (
    id UUID PRIMARY KEY,
    origin_city VARCHAR(100) NOT NULL,
    destination_city VARCHAR(100) NOT NULL,
    service_type_id UUID NOT NULL,
    price_per_kg DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tariffs_service_type FOREIGN KEY (service_type_id) REFERENCES service_types(id) ON DELETE CASCADE
);

-- Table: orders
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    service_type_id UUID NOT NULL,
    tracking_number VARCHAR(100) UNIQUE NOT NULL,
    weight DECIMAL(10, 2) NOT NULL,
    length DECIMAL(10, 2) NOT NULL,
    width DECIMAL(10, 2) NOT NULL,
    height DECIMAL(10, 2) NOT NULL,
    shipping_cost DECIMAL(12, 2) NOT NULL,
    insurance_fee DECIMAL(12, 2) NOT NULL,
    total_price DECIMAL(12, 2) NOT NULL,
    status order_status NOT NULL DEFAULT 'WAITING_PAYMENT',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_orders_service_type FOREIGN KEY (service_type_id) REFERENCES service_types(id)
);

-- Table: addresses
CREATE TABLE addresses (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    type address_type NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_addresses_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Table: order_outbox
CREATE TABLE order_outbox (
    id UUID PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Create Indexes for optimization
CREATE INDEX idx_addresses_order_id ON addresses(order_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_order_outbox_status ON order_outbox(status);
