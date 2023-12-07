# Ledger API

The Ledger API is designed for managing ledger accounts, providing functionality for depositing funds, transferring funds between accounts, checking account balances, and retrieving the transfer history.

## Getting Started

To get started with the Ledger API, follow the steps below:

### Prerequisites

- [Docker](https://www.docker.com/get-started) - Make sure you have Docker installed to run the required services.
- [Go](https://golang.org/dl/) - Install Go for building and running the API.
- [PostgreSQL](https://www.postgresql.org/download/) - Install PostgreSQL as the database for the API.
- [Redis](https://redis.io/download) - Install Redis for caching purposes.

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/julioc98/ledger.git
    cd ledger
    ```

2. Run the database migrations:

    ```bash
    make migrate-up
    ```

3. Start the API:

    ```bash
    make run
    ```

### Usage

The API will be accessible at [http://localhost:3000](http://localhost:3000). Explore the API using Swagger documentation available at [http://localhost:3000/swagger](http://localhost:3000/swagger).

## Makefile Commands

- `make docker-run`: Run the API in a Dockerized environment.
- `make migrate-up`: Run PostgreSQL database migrations.
- `make run`: Build and run the API.
- `make docker-tests-up`: Set up Dockerized test environment.
- `make migrate-tests`: Run migrations for the test database.
- `make test`: Run the API tests.

## API Endpoints

### Create Deposit

```http
POST /api/v1/account/{account}/deposit
```

Deposit funds into a ledger account.

### Create Transfer

```http
POST /api/v1/account/{account}/transfers
```

Transfer funds from one ledger account to another.

### Get Account Balance

```http
GET /api/v1/account/{account}/balance
```

Retrieve the balance of a ledger account.

### Get Transfer History

```http
GET /api/v1/account/{account}/transfers
```

Retrieve the transfer history of a ledger account.

## Configuration

The API can be configured using environment variables.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the Ledger API.
