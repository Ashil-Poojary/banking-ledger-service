# Banking Ledger Service

## Introduction

The Banking Ledger Service manages banking transactions and account information. It uses PostgreSQL, MongoDB, Redis, and RabbitMQ.

## Installation & Setup

### Prerequisites

- Git
- Docker and Docker Compose (Recommended)
- Go (If not using Docker)

### 1. Clone the Repository

git clone https://github.com/ashil-poojary/banking-ledger-service.gitcd banking-ledger-service

### 2. Set Up Environment Variables

Copy and edit the `.env` file:

cp .env.example .env
**.env**:

POSTGRES_HOST= # PostgreSQL hostname (e.g., 'localhost' or 'postgres_db' in Docker)POSTGRES_PORT=5432POSTGRES_USER= # PostgreSQL usernamePOSTGRES_PASSWORD= # PostgreSQL passwordPOSTGRES_DB= # PostgreSQL database nameMIGRATE_DB=falseMONGO_URI= # MongoDB connection stringREDIS_HOST= # Redis hostnameREDIS_PORT=6379REDIS_PASSWORD= # Redis password (if applicable)RABBITMQ_HOST= # RabbitMQ hostnameRABBITMQ_PORT=5672RABBITMQ_USER= # RabbitMQ usernameRABBITMQ_PASSWORD= # RabbitMQ password

### 3. Start Services with Docker (Recommended)

docker-compose up --build

### 4. Start the API Server Manually (Alternative)

go run cmd/api/main.go

### 5. Start the Worker Service (RabbitMQ Consumer)

go run cmd/worker/main.go

### Checking the Logs

- **Docker:** `docker-compose logs -f`
- **Manual:** Terminal output

## API Endpoints

Available at `http://localhost:8080`.

## Troubleshooting

- Ensure dependencies are running.
- `Check .env file.`
- For Docker issues:

  ```
  docker-compose down
  docker-compose up --build --force-recreate
  ```
