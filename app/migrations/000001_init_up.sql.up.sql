CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
    name TEXT NOT NULL PRIMARY KEY,
    type TEXT NOT NULL,
    flavor TEXT NOT NULL,
    size TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    manufacture_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers (
    name TEXT NOT NULL PRIMARY KEY,
    points INTEGER DEFAULT 0
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_name TEXT NOT NULL REFERENCES customers(name),
    product_name TEXT NOT NULL REFERENCES products(name),
    product_size TEXT NOT NULL,
    product_flavor TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE point_redemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_name TEXT NOT NULL REFERENCES customers(name),
    product_name TEXT NOT NULL REFERENCES products(name),
    quantity INTEGER NOT NULL,
    redeemed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);
