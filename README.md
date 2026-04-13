# Scalable Coupon System

A Coupon Management System designed to handle high concurrency loads, specifically "Flash Sale" scenarios, with a focus on data integrity and stock consistency.

## Architecture

This project implements **Clean Architecture** principles to ensure modular, testable, and well-organized code:
- **Domain**: Business entities and repository contracts.
- **Usecase**: Core business logic and transaction orchestration.
- **Repository**: Data access implementation (PostgreSQL).
- **Delivery**: HTTP layer (Fiber framework).

## Key Features

- **Concurrency Safety**: Uses *Pessimistic Locking* (`SELECT FOR UPDATE`) to guarantee that coupon stock is never over-claimed.
- **Atomicity**: The entire claim process is wrapped in a single database transaction.
- **Integrated Infrastructure**: Ready to run with Docker and Docker Compose.

## Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Fiber v2](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL 15](https://www.postgresql.org/)
- **Infrastructure**: [Docker](https://www.docker.com/)

## Getting Started

### 1. Prerequisites
- Docker & Docker Compose installed on your machine.

### 2. Running the Application
Check your `.env` file and run Docker Compose:
```bash
docker compose up -d --build
```
The application will be available at `http://localhost:8080`.

## Testing

### Running Unit & Integration Tests
Ensure a local PostgreSQL database is running or use Docker, then run:
```bash
go test -v ./...
```

### Flash Sale Simulation (Concurrency Test)
There is a specific test to simulate 50 concurrent requests for a coupon with limited stock:
```bash
go test -v internal/usecase/coupon_concurrency_test.go internal/usecase/coupon_usecase_impl.go internal/usecase/coupon_usecase.go
```

## API Endpoints

- `POST /api/coupons`: Create a new coupon.
- `GET /api/coupons/:name`: Check coupon status/stock.
- `POST /api/coupons/claim`: Claim a coupon for a user.
