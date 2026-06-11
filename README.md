# PharmaStock Backend

A B2B pharmaceutical stock management platform connecting **Stockists** (distributors) and **Retailers** (pharmacies).

Stockists upload inventory files, the system processes them into a searchable catalog. Retailers can discover stockists with the medicines they need.

**Status**: Development | **Go Version**: 1.26

---

## Tech Stack

### Backend
- **Go** — Type-safe, compiled language
- **Echo v5** — Minimal, high-performance HTTP framework with middleware
- **pgx/v5** — Native PostgreSQL driver with connection pooling
- **go-playground/validator** — Struct tag validation for request DTOs

### Database
- **PostgreSQL** (13+) — Primary data store
- **Raw SQL** — No ORM; explicit query control

### Logging
- **Zap** — Structured, leveled logging (JSON in production, console in development)

### Infrastructure
- **Docker** / **Docker Compose** — Containerized PostgreSQL

---

## Project Structure

```
pharmaX-server/
├── cmd/
│   ├── api/main.go              # API server entry point
│   └── worker/main.go           # Background job worker entry point
│
├── internal/
│   ├── app/
│   │   └── app.go               # App bootstrap, logger init, DI wiring
│   │
│   ├── common/
│   │   └── response.go          # Standardized JSON response helpers
│   │
│   ├── config/
│   │   └── config.go            # Env-based configuration
│   │
│   ├── database/
│   │   └── postgres.go          # pgxpool connection setup
│   │
│   ├── middleware/
│   │   ├── logger.go            # Structured request logging (zap)
│   │   ├── rate_limit.go        # IP-based rate limiting
│   │   ├── recovery.go          # Panic recovery
│   │   └── request_id.go        # Request ID injection
│   │
│   ├── health/
│   │   └── health.go            # Health check endpoint
│   │
│   ├── stockist/                # Stockist (distributor) module
│   │   ├── model.go             # Domain model
│   │   ├── dto.go               # Request/response DTOs with validation
│   │   ├── validator.go         # Struct validator instance
│   │   ├── handler.go           # HTTP handlers
│   │   ├── service.go           # Business logic + pagination
│   │   ├── repository.go        # DB operations
│   │   ├── module.go            # DI wiring
│   │   └── routes.go            # Route registration
│   │
│   ├── retailer/                # Retailer (pharmacy) module
│   │   ├── model.go             # Domain model
│   │   ├── dto.go               # Request/response DTOs with validation
│   │   ├── validator.go         # Struct validator instance
│   │   ├── handler.go           # HTTP handlers
│   │   ├── service.go           # Business logic + pagination
│   │   ├── repository.go        # DB operations
│   │   ├── module.go            # DI wiring
│   │   └── routes.go            # Route registration
│   │
│   ├── job/                     # Inventory job module
│   │   ├── model.go             # Job domain model
│   │   ├── service.go           # Job processing logic
│   │   ├── repository.go        # Job DB operations
│   │   └── processor.go         # Job processor stub
│   │
│   └── router/
│       └── router.go            # Route registration hub
│
├── migrations/
│   ├── 000001_create_stockists.up.sql
│   ├── 000001_create_stockists.down.sql
│   ├── 000002_create_inventory_jobs.up.sql
│   ├── 000002_create_inventory_jobs.down.sql
│   ├── 000003_create_retailers.up.sql
│   └── 000003_create_retailers.down.sql
│
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Architecture

### Request Flow

```
HTTP Request
    ↓
Middleware: RequestID → Logger → Recovery → RateLimit
    ↓
Handler  (bind DTO, validate, map to domain model)
    ↓
Service  (business logic, duplicate check, pagination)
    ↓
Repository (SQL queries, data access)
    ↓
PostgreSQL
```

### Layer Responsibilities

| Layer | Role |
|---|---|
| **DTO** | Request/response shapes; validation tags; `ToDomain()` mapping |
| **Handler** | HTTP parsing, DTO binding, status codes, response serialization |
| **Service** | Business rules (duplicate email check), pagination, error wrapping |
| **Repository** | SQL CRUD, connection management, sentinel errors |

### Key Decisions

- **DTOs separate from domain models** — API contract changes don't affect internal logic
- **Domain models are pure data** — no JSON or validation tags, used only internally
- **Sentinel errors** (`ErrNotFound`, `ErrDuplicateEmail`) — handlers discriminate via `errors.Is()` for correct HTTP status codes
- **No ORM** — raw SQL with `$1, $2` parameterized queries
- **Feature modules** — each domain (stockist, retailer, job) is self-contained with its own handler/service/repository

---

## Features

### Current

- **Stockist CRUD** — Create, read (by email), update, delete, list (paginated) pharmaceutical distributors
- **Retailer CRUD** — Create, read (by email), update, delete, list (paginated) pharmacies
- **Duplicate email detection** — Returns 409 Conflict on duplicate email for both stockists and retailers
- **Pagination** — `?page=1&limit=20` on list endpoints (default page=1, limit=20, max=100)
- **Background Job Worker** — Polls pending inventory jobs every 10s, processes up to 5 per cycle
- **Health Check** — Reports API + database status
- **Rate Limiting** — 100 requests per 5 minutes per IP on all business endpoints
- **Request Tracing** — Unique `X-Request-ID` on every response
- **Structured Logging** — Zap logger with method, path, status, latency, client IP, user agent; logs at error level on failures
- **Graceful Shutdown** — SIGINT/SIGTERM drains connections within 10s timeout
- **Configurable Log Format** — `APP_ENV=production` → JSON, otherwise human-readable console

### In Progress

- Authentication & Authorization (JWT)
- Inventory file upload module
- PDF extraction / medicine catalog
- Retailer medicine search

---

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 13+ (or Docker)

### Local Development

```bash
# 1. Clone
git clone <repo>
cd pharmaX-server

# 2. Start PostgreSQL (Docker)
docker-compose up -d

# 3. Copy and configure environment
cp .env.example .env
# Edit .env if needed (defaults work with docker-compose)

# 4. Run migrations
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/pharmastock-db?sslmode=disable" up

# 5. Start the API server
go run cmd/api/main.go

# 6. (optional) Start the background worker
go run cmd/worker/main.go
```

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `8080` | HTTP server port |
| `APP_ENV` | `development` | `production` for JSON logs |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `pharmastock-db` | Database name |
| `DB_SSL_MODE` | `disable` | SSL mode |

---

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Health

```http
GET /api/v1/health
```

Response `200`:
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "checks": {
      "api": "up",
      "database": "up"
    }
  }
}
```

### Stockists

#### Create Stockist

```http
POST /api/v1/stockists
Content-Type: application/json

{
  "name": "John Doe",
  "business_name": "Healthcare Supplies Ltd",
  "email": "john@example.com",
  "phone": "9876543210",
  "country": "India",
  "state": "Maharashtra",
  "city": "Mumbai",
  "pin_code": 400001,
  "address": "123 Main Street, Mumbai",
  "gst_number": "27AXXXXX0001A1Z5"
}
```

Response `201`:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "business_name": "Healthcare Supplies Ltd",
    "email": "john@example.com",
    "phone": "9876543210",
    "country": "India",
    "state": "Maharashtra",
    "city": "Mumbai",
    "pin_code": 400001,
    "address": "123 Main Street, Mumbai",
    "gst_number": "27AXXXXX0001A1Z5"
  }
}
```

#### Get Stockist by Email

```http
GET /api/v1/stockists/john@example.com
```

Response `200`: Same shape as create response.

#### Update Stockist

```http
PUT /api/v1/stockists/1
Content-Type: application/json

{
  "name": "Jane Doe",
  "business_name": "Healthcare Supplies Ltd",
  "email": "jane@example.com",
  "phone": "9876543210",
  "country": "India",
  "state": "Maharashtra",
  "city": "Mumbai",
  "pin_code": 400001,
  "address": "456 Oak Street, Mumbai",
  "gst_number": "27AXXXXX0001A1Z5"
}
```

Response `200`: Updated stockist object.

#### Delete Stockist

```http
DELETE /api/v1/stockists/1
```

Response `200`:
```json
{
  "success": true,
  "message": "stockist deleted successfully",
  "data": null
}
```

#### List Stockists (Paginated)

```http
GET /api/v1/stockists?page=1&limit=20
```

Response `200`:
```json
{
  "success": true,
  "data": {
    "items": [ ... ],
    "total": 42,
    "page": 1,
    "limit": 20,
    "total_pages": 3
  }
}
```

### Retailers

All retailer endpoints mirror the stockist endpoints under `/api/v1/retailers`.

| Method | Route | Description |
|---|---|---|
| `POST` | `/api/v1/retailers` | Create retailer |
| `GET` | `/api/v1/retailers/:email` | Get retailer by email |
| `PUT` | `/api/v1/retailers/:id` | Update retailer |
| `DELETE` | `/api/v1/retailers/:id` | Delete retailer |
| `GET` | `/api/v1/retailers?page=&limit=` | List retailers (paginated) |

Request/response shapes are identical to stockists, with the addition of `created_at` and `updated_at` timestamps on the response.

---

## Database Schema

### stockists

```sql
CREATE TABLE stockists (
    id BIGSERIAL PRIMARY KEY,
    owner_name VARCHAR(255) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    pin_code VARCHAR(20) NOT NULL,
    address VARCHAR(255) NOT NULL,
    gst_number VARCHAR(50) NOT NULL
);
```

### retailers

```sql
CREATE TABLE retailers (
    id BIGSERIAL PRIMARY KEY,
    owner_name VARCHAR(255) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    pin_code VARCHAR(20) NOT NULL,
    address VARCHAR(255) NOT NULL,
    gst_number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### jobs

```sql
CREATE TABLE jobs (
    id BIGSERIAL PRIMARY KEY,
    stockist_id BIGINT NOT NULL REFERENCES stockists(id) ON DELETE CASCADE,
    job_status VARCHAR(20) NOT NULL CHECK (job_status IN ('pending','processing','completed','failed')),
    file_path TEXT NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_jobs_status ON jobs(job_status);
CREATE INDEX idx_jobs_stockist_id ON jobs(stockist_id);
```

---

## Error Responses

All endpoints return a consistent envelope on errors:

```json
{
  "success": false,
  "error": "descriptive message"
}
```

| Code | When |
|---|---|
| `400` | Invalid request body or validation failure |
| `404` | Resource not found |
| `409` | Duplicate email on create |
| `429` | Rate limit exceeded |
| `500` | Internal server error |

---

## Middleware Pipeline

| Middleware | Runs | Description |
|---|---|---|
| **RequestID** | 1st (outermost) | Generates `X-Request-ID` for every request |
| **Logger** | 2nd | Logs method, path, status, latency, IP, user-agent via Zap |
| **Recovery** | 3rd | Catches panics, returns 500 instead of crash |
| **RateLimit** | Per-group | 100 requests / 5 min per IP; returns 429 |

---

## Background Worker

The worker (`cmd/worker/main.go`) is a separate binary that:

1. Polls the database every 10 seconds for `pending` jobs
2. Processes up to 5 jobs per cycle
3. Marks jobs `processing` → `completed` or `failed`
4. Shuts down gracefully on SIGINT/SIGTERM

```bash
go run cmd/worker/main.go
```

---

## Configuration Reference

| Variable | Default | Purpose |
|---|---|---|
| `APP_PORT` | `8080` | HTTP listen port |
| `APP_ENV` | `development` | Log format (`production` → JSON) |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `pharmastock-db` | Database name |
| `DB_SSL_MODE` | `disable` | SSL mode for DB connection |

Pool settings (hardcoded): min 2, max 10 connections, 1h max lifetime, 30m max idle.
