# LessonHub

LessonHub is a SaaS platform for private tutors: manage students, schedule lessons,
track attendance, payments and lesson subscriptions — all in one place.

> **Status:** early development. The application foundation is in place
> (configuration, database connection, HTTP server with graceful shutdown);
> business features are coming next.

## Tech stack

- **Go** — backend, standard library `net/http`, `log/slog`
- **PostgreSQL 17** — database, accessed via [pgx](https://github.com/jackc/pgx)
- **Docker Compose** — local database

## Requirements

- [Go](https://go.dev/dl/) 1.26+
- [Docker](https://www.docker.com/) with Docker Compose

## Quick start

1. Clone the repository:

   ```bash
   git clone https://github.com/Keyrozenxz/lessonhub.git
   cd lessonhub
   ```

2. Create a local environment file from the template:

   ```bash
   cp .env.example .env
   ```

   On Windows PowerShell: `Copy-Item .env.example .env`

3. Start PostgreSQL:

   ```bash
   docker compose up -d
   ```

4. Run the API:

   ```bash
   go run ./cmd/api
   ```

5. Check that it works:

   ```bash
   curl http://localhost:8080/health
   ```

   Expected response: `{"status":"ok"}`

Press `Ctrl+C` to stop the server gracefully.

## Configuration

The application is configured via environment variables. For local development
they are loaded from the `.env` file; in other environments, set them directly.
All variables are required.

| Variable      | Description               | Example        |
|---------------|---------------------------|----------------|
| `DB_HOST`     | PostgreSQL host           | `localhost`    |
| `DB_PORT`     | PostgreSQL port           | `5433`         |
| `DB_USER`     | PostgreSQL user           | `lessonhub`    |
| `DB_PASSWORD` | PostgreSQL password       | `dev_password` |
| `DB_NAME`     | PostgreSQL database name  | `lessonhub`    |
| `HTTP_ADDR`   | HTTP server listen address | `:8080`       |

> **Why port 5433?** The Docker container exposes PostgreSQL on host port `5433`
> instead of the default `5432` to avoid conflicts with a locally installed PostgreSQL.

## Project structure

```
cmd/api/                 application entry point: wiring and startup
internal/config/         loading configuration from environment
internal/database/       PostgreSQL connection pool
internal/transport/http/ HTTP server and routing
```