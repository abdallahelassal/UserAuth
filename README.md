# UserAuth

A lightweight Go user authentication service with PostgreSQL persistence, Docker support, and a simple signup endpoint.

## Features

- User signup with bcrypt password hashing
- PostgreSQL database storage
- Docker Compose for local environment setup
- Configurable via `.env`
- Health-checked database dependency and optional PGAdmin

## Requirements

- Go 1.25+
- Docker & Docker Compose (for containerized setup)
- PostgreSQL (if running without Docker)

## Environment

Create a `.env` file at the project root with values like:

```env
APP_PORT=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
ACCESS_TOKEN_SECRET=
REFRESH_TOKEN_SECRET=
ENVIROMENT=
```

> The application loads `.env` from the project root or parent directories.

## Run Locally

1. Install dependencies:

```bash
go mod download
```

2. Start PostgreSQL separately or use Docker Compose.

3. Run the service:

```bash
go run ./cmd
```

The server listens on `:8000` by default.

## Run with Docker Compose

Start the application and database together:

```bash
docker compose up --build
```

Services:

- `app` => Go API on `http://localhost:8000`
- `db` => PostgreSQL on `5432`
- `pgadmin` => PGAdmin on `http://localhost:8081`

## API Endpoints

- `GET /ping`
  - Health check
  - Response: `{"message":"pong"}`

- `POST /signup`
  - Create a new user
  - Request JSON:

```json
{
  "user_name": "johndoe",
  "email": "john@example.com",
  "password": "Password123"
}
```

- Response on success: `201 Created` with `{"message":"User created successfully"}`

## Notes

- Passwords are encrypted with `bcrypt` before storage.
- Database migration files are stored in the `migration/` folder.
- The current implementation focuses on signup and database connectivity.

