# System Design

## Overview

PharmaStock is a **feature-first modular monolith** — each business domain is self-contained with its own handler, service, repository, model, DTO, and routes. Modules communicate through shared interfaces and dependency injection.

This avoids premature microservice complexity while keeping clear boundaries for future extraction.

---

## Domain Model

### Core Entities

```
User ──┬── admin     (seeded from env)
       ├── stockist  (FK → stockists.id via reference_id)
       └── retailer  (FK → retailers.id via reference_id)

Stockist ────< Inventory >──── Medicine

Retailer (independent, no direct relation to other entities)

Job ─── belongs to Stockist (stockist_id FK)
```

### Entity Relationships

| Entity | Relationships |
|---|---|
| **User** | Login credential, linked to Stockist/Retailer via `reference_id` |
| **Stockist** | Distributor — owns inventory records and jobs |
| **Retailer** | Pharmacy — independent entity for future orders |
| **Medicine** | Global catalog — unique name constraint |
| **Inventory** | Join table (stockist_id, medicine_id) — which stockist carries which medicine |
| **Job** | File upload processing — created by upload, processed by worker |

---

## Module Map

```
internal/
├── auth/        Auth, JWT, RBAC middleware
├── stockist/    Distributor CRUD
├── retailer/    Pharmacy CRUD
├── medicine/    Global medicine catalog + file parsers
├── inventory/   Stockist-medicine join table
├── job/         Background job processing (incl. ResetStaleJobs)
├── upload/      File upload handler + service
├── ui/          Browser testing interface (HTMX + Alpine.js, Go templates)
├── health/      Health check endpoint
├── middleware/  Global middleware pipeline
├── router/      API route registration hub
├── common/      Shared response helpers, sentinel errors
├── config/      Env-based configuration
├── database/    PostgreSQL connection pool
└── app/         DI wiring, bootstrap, Renderer registration
```

---

## Data Flow

### Inventory Upload & Processing

```
Stockist
  │ POST /api/v1/upload (multipart: file + stockist_id)
  ▼
Upload Handler → validates file (.csv/.pdf)
  │
  ▼
Upload Service → saves file to disk → creates Job (status: pending)
  │
  ▼ (async)
Worker (polls every 10s)
  │ 1. ResetStaleJobs (resets jobs stuck in "processing" >5min)
  │ 2. fetches "pending" jobs (max 5)
  ▼
Job Processor
  │ 1. Parse file (CSV or PDF)
  │ 2. BatchInsert new medicines (ON CONFLICT DO NOTHING)
  │ 3. GetMedicinesByNames → get ID map
  │ 4. BulkCreate inventory entries (ON CONFLICT DO NOTHING)
  ▼
Job marked "completed" (or "failed" with error message)
```

### Authentication Flow

```
Client                          Server
  │ POST /auth/login              │
  │ { email|username, password }  │
  ▼                               ▼
Handler → validate → Service
  │                              ┌─────────────────┐
  ├── email != "" ───────────────► GetUserByEmail   │
  │                              └─────────────────┘
  │                              ┌─────────────────────┐
  └── username != "" ───────────► GetUserByUsername     │
                                 └─────────────────────┘
  ▼
checkPassword (bcrypt)
  ▼
generateJWT (HS256, 24h expiry)
  ▼
{ token, user_id, role, reference_id }
```

### Request Lifecycle

```
HTTP Request
    │
    ▼
Middleware: RequestID → Logger → Recovery
    │
    ▼
Group Middleware: AuthRequired → RequireRole (per group) [API only]
    │
    ▼
Handler (bind DTO, manual validate, call service)
    │
    ▼
Service (business logic, call repository)
    │
    ▼
Repository (raw SQL via pgx)
    │
    ▼
PostgreSQL
    │
    ▼
Response ← common.APISuccessResponse / common.APIErrorResponse
```

### UI Request Flow (Browser)

```
Browser →  / (root routes)
    │
    ▼
Middleware: RequestID → Logger → Recovery
    │
    ▼
UI Handler → calls internal services → renders Go templates
    │
    ▼
TemplateRenderer.Render (per-page clone)
    │
    ▼
HTML response with HTMX + Alpine.js
    │
    ▼
HTMX makes AJAX calls → partial template renders (lists, forms)
```

---

## Key Design Decisions

| Decision | Rationale |
|---|---|
| **Modular monolith** | Clear domain boundaries without microservice overhead |
| **Feature modules** | Each domain is a self-contained package (model, repo, service, handler, routes, module.go) |
| **Clean Architecture** | Handler → Service → Repository (dependencies point inward) |
| **Raw SQL / pgx** | Full control over queries, no ORM overhead or hidden N+1 |
| **DTOs separate from domain** | API contract changes don't affect internal logic; domain models have no JSON/validation tags |
| **Sentinel errors** | `ErrNotFound`, `ErrDuplicateEmail` — handlers discriminate via `errors.Is()` for correct HTTP codes |
| **Pointer context (`*echo.Context`)** | Echo v5 uses struct context; all handlers/middleware use pointer semantics |
| **pg_trgm GIN index on medicines** | Fuzzy medicine name search with trigram similarity |
| **Composite PK on inventories** | `PRIMARY KEY (stockist_id, medicine_id)` prevents duplicates |
| **ON CONFLICT DO NOTHING** | Idempotent inserts for medicines and inventories |
| **Polling-based worker** | Simple, no external queue dependency; scales to moderate load |
| **Local file storage** | Files saved to `./uploads/`; cloud storage migration when scaling |
| **bcrypt cost 10** | Default cost — balances security and performance |
| **Per-page template clones** | Prevents `{{define "content"}}` collisions across page templates |
| **HTMX + Alpine.js UI** | Server-rendered HTML with minimal client JS; no SPA framework needed |

---

## Database Schema

### Tables

#### users

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin','stockist','retailer')),
    reference_id BIGINT,                   -- FK to stockists.id or retailers.id
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### stockists

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
    gst_number VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### retailers

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

#### medicines

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE medicines (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_medicines_name_trgm ON medicines USING GIN (name gin_trgm_ops);
CREATE INDEX idx_medicines_name ON medicines (name);
```

#### inventories

```sql
CREATE TABLE inventories (
    stockist_id BIGINT NOT NULL REFERENCES stockists(id) ON DELETE CASCADE,
    medicine_id BIGINT NOT NULL REFERENCES medicines(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stockist_id, medicine_id)
);

CREATE INDEX idx_inventories_medicine_id ON inventories (medicine_id);
```

#### jobs (inventory_jobs)

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

### Migration Sequence

| # | Name | Key Change |
|---|---|---|
| 000001 | `create_stockists` | Base stockist table |
| 000002 | `create_inventory_jobs` | Job processing table |
| 000003 | `create_retailers` | Retailer table with timestamps |
| 000004 | `create_medicines` | Medicine catalog with pg_trgm |
| 000005 | `create_inventories` | Stockist-medicine join table |
| 000006 | `create_users` | Auth users table with roles |
