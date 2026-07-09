![PharmaStock Banner](pharmastock.png)

# PharmaStock Backend

B2B pharmaceutical stock management platform connecting **Stockists** (distributors) and **Retailers** (pharmacies).

Stockists upload inventory files, the system processes them into a searchable catalog. Retailers discover stockists with the medicines they need.

**Status**: Development | **Go**: 1.24+ | **License**: MIT

## Index

- [Architecture](docs/ARCHITECTURE.md)
- [System Design](docs/SYSTEM_DESIGN.md)
- [API Reference (OpenAPI)](docs/openapi.yaml)

---

## How It Works

```mermaid
flowchart LR
    subgraph Actors[" "]
        STK["<b>Stockist</b> (Distributor)"]
        RTL["<b>Retailer</b> (Pharmacy)"]
        ADM["<b>Admin</b>"]
    end

    subgraph System["PharmaStock Backend"]
        API["API Server<br/>Echo v5"]
        UI["Browser UI<br/>HTMX + Alpine.js"]
        W["Worker<br/>poll every 10s"]
        DB[("PostgreSQL")]
    end

    STK -->|"1. Upload .csv/.pdf"| API
    API -->|"2. Save file + create Job"| DB
    W -->|"3. Poll pending jobs"| DB
    W -->|"4. Parse file → seed medicines<br/>→ link inventory"| DB

    RTL -->|"5. Search medicines"| UI
    ADM -->|"Manage stockists & retailers"| UI
    UI -->|"Query catalog, lookup stockists"| API
    API -->|"6. Return results"| DB

    style Actors fill:#e6f3e6,stroke:#2d6a2d
    style System fill:#e6f0ff,stroke:#2d4a8a
```

---

## Quick Start

### Prerequisites

- Go 1.24+, PostgreSQL 13+ (or Docker), [golang-migrate CLI](https://github.com/golang-migrate/migrate)

```mermaid
flowchart LR
    GIT["git clone"] --> DOCKER["docker compose up -d"]
    DOCKER --> ENV["cp .env.example .env"]
    ENV --> MIGRATE["migrate -path migrations up"]
    MIGRATE --> API["go run ./cmd/api"]
    MIGRATE --> WORKER["go run ./cmd/worker"]
```

```bash
git clone https://github.com/swaindhruti/pharmastock-backend && cd pharmastock-backend

# Start PostgreSQL
docker compose up -d

# Configure environment
cp .env.example .env

# Install migrate if needed
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/pharmastock-db?sslmode=disable" up

# Start API server (terminal 1)
go run ./cmd/api

# Start background worker (terminal 2)
go run ./cmd/worker
```

### Environment Variables

| Variable | Default | Required | Purpose |
|---|---|---|---|
| `APP_PORT` | `8080` | No | HTTP listen port |
| `APP_ENV` | `development` | No | Log format (`production` → JSON) |
| `DB_HOST` | `localhost` | No | PostgreSQL host |
| `DB_PORT` | `5432` | No | PostgreSQL port |
| `DB_USER` | `postgres` | No | Database user |
| `DB_PASSWORD` | `postgres` | No | Database password |
| `DB_NAME` | `pharmastock-db` | No | Database name |
| `DB_SSL_MODE` | `disable` | No | SSL mode |
| `JWT_SECRET` | — | **Yes** | JWT signing key |
| `UPLOAD_DIR` | `./uploads` | No | File upload directory |
| `ADMIN_USERNAME` | `admin` | No | Default admin username |
| `ADMIN_PASSWORD` | — | **Yes** | Admin password |
| `ADMIN_EMAIL` | — | **Yes** | Admin email |

---

## Test UI

The project includes a browser-based testing interface served at the **root domain**. No API client needed.

| Route | Description |
|---|---|
| `GET /` | Dashboard with stockist/retailer counts |
| `GET /login` | Login form (email or username + password) |
| `GET /stockists` | Stockist list with create/edit/delete |
| `GET /retailers` | Retailer list with create/edit/delete |
| `GET /medicines` | Medicine search (fuzzy) |
| `GET /inventory` | Look up stockists by medicine ID |
| `GET /upload` | Upload inventory file (.csv/.pdf) |

Open **http://localhost:8080** and login with the seeded admin credentials.

---

## Documentation

| File | What's Inside |
|---|---|
| `docs/ARCHITECTURE.md` | Project structure, module pattern, middleware pipeline, template rendering, auth & JWT, worker, graceful shutdown |
| `docs/SYSTEM_DESIGN.md` | System architecture, domain model (ERD), user roles, auth flow, upload flow, request lifecycle, design decisions, database schema |
| `docs/openapi.yaml` | Full OpenAPI 3.0 specification for all `/api/v1/*` endpoints |
