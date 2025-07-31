# BSnack API

A high-performance API system for managing product sales and point-based rewards, built for weekly sales of crispy snacks with various flavors and sizes.

## Features

- **Product Management**
  - Add new products with details (name, size, flavor, price, manufacture date)
  - View products filtered by manufacture date
  - Automatic stock deduction on purchases

- **Customer Rewards**
  - Earn points on purchases (1 point per Rp 1,000 spent)
  - Point-based product redemption system
  - Track customer transaction history

- **Transaction Processing**
  - Record product purchases
  - Handle point redemptions
  - Maintain inventory levels

## Tech Stack

- **Backend**: Go 1.21+
- **Database**: PostgreSQL 14+
- **Cache**: Redis 7+
- **Containerization**: Docker & Docker Compose

## Getting Started

### Environment Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Configure your environment variables in `.env`:
   ```env
   # Server
   PORT=8080
   ENV=development
   TIMEZONE=Asia/Jakarta
   
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=bsnack
   DB_PASSWORD=bsnack
   DB_NAME=bsnack
   DB_SSLMODE=disable
   
   # Redis
   REDIS_HOST=localhost
   REDIS_PORT=6379
   ```

### Development Commands

- Start API server only:
  ```bash
  ./run.sh
  ```

### Project Structure

```
app/
├── cmd/           # Application entry points
├── internal/      # Private application code
│   ├── api/       # API versioning
│   ├── customer/  # Customer domain
│   ├── product/   # Product domain
│   └── transaction/ # Transaction domain
├── pkg/           # Reusable packages
└── migrations/    # Database migrations
```

## Postman Collection

Import the Postman collection from either `Bsnack.postman_collection.json` or `Bsnack_v2.postman_collection.json`.
