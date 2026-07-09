# System Design

## Overview

PharmaStock is a **feature-first modular monolith** — each business domain is self-contained with its own handler, service, repository, model, DTO, and routes linked via dependency injection. This avoids microservice overhead while keeping clear boundaries for future extraction.

---

## System Architecture

```mermaid
graph TB
    subgraph External["External"]
        Browser["Browser / HTMX Client"]
    end

    subgraph Server["Echo v5 Server"]
        Router["Router<br/>(/api/v1/* , / , /health)"]
        MW["Middleware Stack<br/>RequestID → Logger → Recovery"]

        subgraph Modules["Feature Modules"]
            Auth["auth<br/>Login · JWT · RBAC"]
            Stockist["stockist<br/>Distributor CRUD"]
            Retailer["retailer<br/>Pharmacy CRUD"]
            Medicine["medicine<br/>Catalog · Parsers"]
            Inventory["inventory<br/>Stockist-Medicine Join"]
            Upload["upload<br/>File Upload · Job Create"]
            Health["health<br/>DB Ping"]
            UI["ui<br/>Go Templates · HTMX · Alpine.js"]
            Job["job<br/>Background Processing"]
        end

        Renderer["Template Renderer<br/>(html/template per-page clones)"]
    end

    subgraph Worker["Background Worker"]
        W["Job Processor<br/>polls every 10s"]
    end

    subgraph Storage["Storage Layer"]
        PG[("PostgreSQL<br/>(pgx connection pool)")]
        FS[("File System<br/>./uploads/")]
    end

    Browser --> Router
    Router --> MW
    MW --> Auth & Stockist & Retailer & Medicine & Inventory & Upload & Health

    Router --> UI
    UI --> Renderer
    UI --> Auth & Stockist & Retailer & Medicine & Inventory

    Auth --> PG
    Stockist --> PG
    Retailer --> PG
    Medicine --> PG
    Inventory --> PG
    Upload --> PG
    Upload --> FS
    Job --> PG
    Health --> PG

    W --> PG
    W --> FS

    style Worker fill:#f9f,stroke:#333,stroke-width:1px
    style Storage fill:#bbf,stroke:#333,stroke-width:1px
```

### What runs where

| Process | Binary | Command | Purpose |
|---|---|---|---|
| **API Server** | `cmd/api/main.go` | `go run ./cmd/api` | Serves HTTP (API + UI), handles file uploads, creates jobs |
| **Worker** | `cmd/worker/main.go` | `go run ./cmd/worker` | Polls DB for pending jobs, processes files asynchronously |

---

## Domain Model

### Entity-Relationship Diagram

```mermaid
erDiagram
    users {
        bigint id PK
        varchar email UK
        varchar username UK
        varchar password_hash
        varchar role "admin | stockist | retailer"
        bigint reference_id "FK → stockists.id or retailers.id"
        timestamp created_at
        timestamp updated_at
    }

    stockists {
        bigint id PK
        varchar owner_name
        varchar business_name
        varchar email UK
        varchar phone
        varchar country
        varchar state
        varchar city
        varchar pin_code
        varchar address
        varchar gst_number
        timestamp created_at
        timestamp updated_at
    }

    retailers {
        bigint id PK
        varchar owner_name
        varchar business_name
        varchar email UK
        varchar phone
        varchar country
        varchar state
        varchar city
        varchar pin_code
        varchar address
        varchar gst_number
        timestamp created_at
        timestamp updated_at
    }

    medicines {
        bigint id PK
        varchar name UK
        timestamp created_at
    }

    inventories {
        bigint stockist_id PK, FK
        bigint medicine_id PK, FK
        timestamp created_at
    }

    jobs {
        bigint id PK
        bigint stockist_id FK
        varchar job_status "pending | processing | completed | failed"
        text file_path
        text error_message
        timestamp created_at
        timestamp started_at
        timestamp completed_at
    }

    stockists ||--o{ inventories : "has"
    medicines ||--o{ inventories : "appears in"
    stockists ||--o{ jobs : "owns"
    users }o--|| stockists : "reference_id (optional, polymorphic)"
    users }o--|| retailers : "reference_id (optional, polymorphic)"
```

### Entity Roles

| Entity | What it is | Key Relationships |
|---|---|---|
| **Stockist** | Distributor who supplies medicines | owns Inventory records, owns Jobs |
| **Retailer** | Pharmacy who buys medicines | standalone (no direct FKs yet) |
| **Medicine** | A single medicine in the global catalog | appears in many Inventories |
| **Inventory** | Join record — "this stockist carries this medicine" | links Stockist ↔ Medicine |
| **Job** | Processing record for a file upload | belongs to a Stockist |
| **User** | Login credential with role | optionally linked to Stockist or Retailer via `reference_id` |

---

## User Roles & Permissions

| Role | Created By | Routes |
|---|---|---|
| **admin** | Seeded from env vars on startup | Everything |
| **stockist** | Admin creates via `POST /auth/admin/stockists` | Medicine search, inventory lookup, file upload |
| **retailer** | Self-registers via `POST /auth/register` | Medicine search, inventory lookup |

---

## Key Flows

### 1. Authentication (login → JWT → protected request)

```mermaid
sequenceDiagram
    actor C as Client
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant DB as PostgreSQL

    C->>H: POST /auth/login<br/>{email or username, password}
    H->>S: Login(dto)
    alt email is provided
        S->>R: GetUserByEmail(email)
    else username is provided
        S->>R: GetUserByUsername(username)
    end
    R->>DB: SELECT * FROM users WHERE ...
    DB-->>R: user row
    R-->>S: User model
    S->>S: bcrypt.Compare(password_hash, password)
    S->>S: jwt.Generate(claims, 24h expiry)
    S-->>H: token, user_id, role, reference_id
    H-->>C: 200 OK<br/>{token, user_id, role, reference_id}

    Note over C,DB: Subsequent requests
    C->>H: GET /api/v1/medicines?q=para<br/>Authorization: Bearer <token>
    H->>H: jwt.Validate(token) → claims
    H->>H: c.Set("user_id"), c.Set("user_role")
    H->>S: SearchMedicines(q)
    S->>R: SearchMedicines(q)
    R->>DB: SELECT * FROM medicines WHERE name % 'para'
    DB-->>R: matching medicines
    R-->>S: []Medicine
    S-->>H: []Medicine
    H-->>C: 200 OK {data: {items: [...]}}
```

### 2. Inventory Upload & Background Processing

```mermaid
flowchart TB
    START(["Stockist uploads file"]) --> POST["POST /api/v1/upload<br/>(multipart: file + stockist_id)"]
    POST --> VALIDATE{"File extension valid?<br/>(.csv or .pdf)"}

    VALIDATE -->|Invalid| ERR["400 Bad Request"]
    VALIDATE -->|Valid| SAVE["Save file to ./uploads/"]
    SAVE --> CREATE_JOB["Create Job<br/>status = pending"]
    CREATE_JOB --> RESP["200 OK<br/>{job_id, status: pending}"]
    RESP --> DONE(["Done (async from here)"])

    subgraph Worker_Loop["Background Worker Loop (every 10s)"]
        direction TB
        RESET["ResetStaleJobs()<br/>→ reset jobs stuck in 'processing' >5min back to 'pending'"]
        FETCH["Fetch pending jobs (max 5)"]
        DECIDE{"Any jobs?"}
        MARK["Mark job → 'processing'<br/>(set started_at)"]
        PARSE["Parse file<br/>CSV → csv.NewReader<br/>PDF → ledongthuc/pdf"]
        BATCH["BatchInsert medicines<br/>(ON CONFLICT DO NOTHING)"]
        LOOKUP["GetMedicinesByNames<br/>→ get ID map"]
        BULK["BulkCreate inventories<br/>(stockist_id, medicine_id)<br/>(ON CONFLICT DO NOTHING)"]
        COMPLETE["Mark job → 'completed'<br/>(set completed_at)"]
        FAIL["Mark job → 'failed'<br/>(set error_message)"]
    end

    RESET --> FETCH --> DECIDE
    DECIDE -->|No| RESET
    DECIDE -->|Yes| MARK --> PARSE --> BATCH --> LOOKUP --> BULK
    BULK --> COMPLETE --> RESET
    BULK --> FAIL --> RESET

    style Worker_Loop fill:#f9f,stroke:#333,stroke-width:1px
```

### 3. Request Lifecycle (API vs UI)

```mermaid
flowchart TB
    subgraph API_Flow["API Request (/api/v1/*)"]
        direction TB
        A1["HTTP Request"] --> A2["RequestID"]
        A2 --> A3["Logger"]
        A3 --> A4["Recovery"]
        A4 --> A5{"Route requires auth?"}
        A5 -->|No: /auth/login, /health| A6["Handler"]
        A5 -->|Yes| A7["AuthRequired<br/>(validate JWT, extract claims)"]
        A7 --> A8["RequireRole<br/>(admin / stockist / retailer)"]
        A8 --> A6
        A6 --> A9["Service Layer<br/>(business logic)"]
        A9 --> A10["Repository Layer<br/>(raw SQL via pgx)"]
        A10 --> A11["PostgreSQL"]
        A11 --> A12["JSON Response<br/>{success, data/error}"]
    end

    subgraph UI_Flow["UI Request (/)"]
        direction TB
        B1["Browser request"] --> B2["RequestID"]
        B2 --> B3["Logger"]
        B3 --> B4["Recovery"]
        B4 --> B5["UI Handler<br/>(calls internal services)"]
        B5 --> B6["TemplateRenderer.Render()<br/>(per-page clone)"]
        B6 --> B7{"HTMX partial request<br/>(HX-Request header)?"}
        B7 -->|Yes| B8["Render named partial<br/>(e.g. 'stockists_list')"]
        B7 -->|No| B9["Execute 'layout' template<br/>(finds page's own 'content')"]
        B8 --> B10["HTML fragment"]
        B9 --> B11["Full HTML page<br/>+ Alpine.js + HTMX"]
    end
```

---

## Design Decisions

| Decision | Rationale |
|---|---|
| **Modular monolith** | Clear domain boundaries without microservice overhead |
| **Raw SQL / pgx** | Full query control, no ORM hidden N+1 |
| **DTOs ≠ domain models** | API contract changes don't affect internal logic; domain models have no JSON tags |
| **Sentinel errors** | `ErrNotFound` → 404, `ErrDuplicateEmail` → 409; handled via `errors.Is()` |
| **pg_trgm GIN index** | Fuzzy medicine name search with trigram similarity |
| **ON CONFLICT DO NOTHING** | Idempotent inserts — no duplicate error handling needed |
| **Polling-based worker** | Simple, no external queue; reset stale jobs every cycle |
| **Per-page template clones** | Prevents `{{define "content"}}` collisions; each page has its own template set |
| **HTMX + Alpine.js** | Server-rendered HTML, minimal JS, no SPA build step |

---

## Database Schema

### Migration Sequence

| # | Migration | Adds |
|---|---|---|
| 000001 | `create_stockists` | `stockists` table |
| 000002 | `create_inventory_jobs` | `jobs` table |
| 000003 | `create_retailers` | `retailers` table |
| 000004 | `create_medicines` | `medicines` table + `pg_trgm` extension + GIN index |
| 000005 | `create_inventories` | `inventories` table (composite PK) |
| 000006 | `create_users` | `users` table + role CHECK constraint |

### SQL Definitions

<details>
<summary>users</summary>

```sql
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    username      VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('admin','stockist','retailer')),
    reference_id  BIGINT,                           -- FK to stockists.id or retailers.id
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
</details>

<details>
<summary>stockists</summary>

```sql
CREATE TABLE stockists (
    id            BIGSERIAL PRIMARY KEY,
    owner_name    VARCHAR(255) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    phone         VARCHAR(20)  NOT NULL,
    country       VARCHAR(100) NOT NULL,
    state         VARCHAR(100) NOT NULL,
    city          VARCHAR(100) NOT NULL,
    pin_code      VARCHAR(20)  NOT NULL,
    address       VARCHAR(255) NOT NULL,
    gst_number    VARCHAR(50)  NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
</details>

<details>
<summary>retailers</summary>

```sql
CREATE TABLE retailers (
    id            BIGSERIAL PRIMARY KEY,
    owner_name    VARCHAR(255) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    phone         VARCHAR(20)  NOT NULL,
    country       VARCHAR(100) NOT NULL,
    state         VARCHAR(100) NOT NULL,
    city          VARCHAR(100) NOT NULL,
    pin_code      VARCHAR(20)  NOT NULL,
    address       VARCHAR(255) NOT NULL,
    gst_number    VARCHAR(50)  NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```
</details>

<details>
<summary>medicines</summary>

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE medicines (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_medicines_name_trgm ON medicines USING GIN (name gin_trgm_ops);
CREATE INDEX idx_medicines_name ON medicines (name);
```
</details>

<details>
<summary>inventories</summary>

```sql
CREATE TABLE inventories (
    stockist_id  BIGINT NOT NULL REFERENCES stockists(id) ON DELETE CASCADE,
    medicine_id  BIGINT NOT NULL REFERENCES medicines(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stockist_id, medicine_id)
);

CREATE INDEX idx_inventories_medicine_id ON inventories (medicine_id);
```
</details>

<details>
<summary>jobs (inventory_jobs)</summary>

```sql
CREATE TABLE jobs (
    id            BIGSERIAL PRIMARY KEY,
    stockist_id   BIGINT NOT NULL REFERENCES stockists(id) ON DELETE CASCADE,
    job_status    VARCHAR(20) NOT NULL CHECK (job_status IN ('pending','processing','completed','failed')),
    file_path     TEXT NOT NULL,
    error_message TEXT,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at    TIMESTAMP,
    completed_at  TIMESTAMP
);

CREATE INDEX idx_jobs_status ON jobs(job_status);
CREATE INDEX idx_jobs_stockist_id ON jobs(stockist_id);
```
</details>
