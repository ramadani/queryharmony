
-- +migrate Up
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone_number VARCHAR(20),
    address TEXT,
    registration_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    partner_id INT REFERENCES partners(id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX idx_customers_id_partner_id ON customers(id, partner_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_customers_id_partner_id;
DROP TABLE IF EXISTS customers;
