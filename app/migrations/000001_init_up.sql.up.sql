CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE product_type AS ENUM (
  'keripik pangsit'
);

CREATE TYPE product_flavor AS ENUM (
  'jagung bakar',
  'rumput laut',
  'original',
  'jagung manis',
  'keju asin',
  'keju manis',
  'pedas'
);

CREATE TYPE product_size AS ENUM (
  'small',
  'medium',
  'large'
);

CREATE TABLE products (
    name TEXT NOT NULL PRIMARY KEY,
    type product_type NOT NULL,
    flavor product_flavor NOT NULL,
    size product_size NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    manufacture_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers (
    name TEXT NOT NULL PRIMARY KEY,
    points INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_name TEXT NOT NULL,
    product_name TEXT NOT NULL,
    product_size product_size NOT NULL,
    product_flavor product_flavor NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE point_redemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_name TEXT NOT NULL,
    product_name TEXT NOT NULL,
    product_size product_size NOT NULL,
    product_flavor product_flavor NOT NULL,
    quantity INTEGER NOT NULL,
    point_required INTEGER NOT NULL,
    redeemed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
