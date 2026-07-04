# PharmaStock Backend

B2B pharmaceutical stock management platform connecting **Stockists** (distributors) and **Retailers** (pharmacies).

Stockists upload inventory files, the system processes them into a searchable catalog. Retailers discover stockists with the medicines they need.

**Status**: Development | **Go**: 1.26 | **License**: MIT

## Index

- [Architecture & System Design](docs/ARCHITECTURE.md)
- [System Design](docs/SYSTEM_DESIGN.md)
- [API Reference (OpenAPI)](docs/openapi.yaml)

---

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 13+ (or Docker)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

### Setup

```bash
git clone https://github.com/swaindhruti/pharmastock-backend && cd pharmastock-backend

# Start PostgreSQL
docker-compose up -d

# Configure environment
cp .env.example .env

# Run migrations
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/pharmastock-db?sslmode=disable" up

# Start API server
go run cmd/api/main.go

# Start background worker (separate terminal)
go run cmd/worker/main.go
```

### Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `APP_PORT` | `8080` | HTTP listen port |
| `APP_ENV` | `development` | Log format (`production` → JSON) |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `pharmastock-db` | Database name |
| `DB_SSL_MODE` | `disable` | SSL mode for DB connection |
| `JWT_SECRET` | — | JWT signing key (required) |
| `UPLOAD_DIR` | `./uploads` | File upload directory |
| `ADMIN_USERNAME` | `admin` | Default admin username |
| `ADMIN_PASSWORD` | — | Admin password (required) |
| `ADMIN_EMAIL` | — | Admin email (required) |

### Important URLs

| File | Description |
|---|---|
| `docs/ARCHITECTURE.md` | Detailed architecture, middleware pipeline, auth flow, worker design |
| `docs/SYSTEM_DESIGN.md` | System design, decision records, module map, data flow |
| `docs/openapi.yaml` | Full OpenAPI 3.0 spec for all endpoints |
