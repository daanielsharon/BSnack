CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE product_type AS ENUM (
  'Keripik Pangsit'
);

CREATE TYPE product_flavor AS ENUM (
  'Jagung Bakar',
  'Rumput Laut',
  'Original',
  'Jagung Manis',
  'Keju Asin',
  'Keju Manis',
  'Pedas'
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
    points INTEGER DEFAULT 0
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

CustomerName  string    `gorm:"column:customer_name;type:varchar(255);not null"`
	ProductSize   string    `gorm:"column:product_size;type:varchar(255);not null"`
	ProductType   string    `gorm:"column:product_type;type:varchar(255);not null"`
	ProductFlavor string    `gorm:"column:product_flavor;type:varchar(255);not null"`
	Quantity      int       `gorm:"column:quantity;type:int;not null"`
	PointRequired int       `gorm:"column:point_required;type:int;not null"`
	RedeemedAt    time.Time `gorm:"column:redemmed_at;type:timestamp;not null"`

CREATE TABLE point_redemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_name TEXT NOT NULL,
    product_name  NOT NULL REFERENCES products(name),
    product_size product_size NOT NULL,
    product_flavor product_flavor NOT NULL,
    quantity INTEGER NOT NULL,
    point_required INTEGER NOT NULL,
    redeemed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
);
