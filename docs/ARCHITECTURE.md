# Architecture

## Project Structure

```
pharmastock-backend/
├── cmd/
│   ├── api/main.go               # API server entry point
│   └── worker/main.go            # Background job worker entry point
│
├── internal/
│   ├── app/app.go                # Bootstrap, DI wiring, graceful shutdown
│   │
│   ├── common/response.go        # APISuccessResponse / APIErrorResponse helpers
│   │
│   ├── config/config.go          # Env-based config (Viper-like manual parsing)
│   │
│   ├── database/postgres.go      # pgxpool init, connection pooling (min 2, max 10)
│   │
│   ├── middleware/
│   │   ├── request_id.go         # X-Request-ID generation
│   │   ├── logger.go             # Zap structured request logging
│   │   ├── recovery.go           # Panic recovery → 500
│   │   └── rate_limit.go         # IP-based: 100 req / 5 min
│   │
│   ├── health/health.go          # GET /health — API + DB status
│   │
│   ├── auth/                     # Authentication & Authorization
│   │   ├── model.go              # User, Claims, request/response DTOs
│   │   ├── repository.go         # User CRUD (Create, GetByEmail, GetByUsername, AdminExists)
│   │   ├── password.go           # bcrypt hash + verify
│   │   ├── jwt.go                # HS256 generate + validate (24h expiry)
│   │   ├── service.go            # Login, CreateUser, SeedAdmin
│   │   ├── handler.go            # Login, RegisterRetailer, AdminCreateStockist
│   │   ├── middleware.go         # AuthRequired, RequireRole, context helpers
│   │   ├── routes.go             # /auth/* route registration
│   │   └── module.go             # DI wiring → returns Module{Handler, Service}
│   │
│   ├── stockist/                 # Distributor module
│   │   ├── model.go              # Domain model
│   │   ├── dto.go                # Request/response DTOs with validation
│   │   ├── validator.go          # Struct validator instance
│   │   ├── handler.go            # CRUD HTTP handlers
│   │   ├── service.go            # Business logic, duplicate check, pagination
│   │   ├── repository.go         # SQL CRUD
│   │   ├── routes.go             # Route registration
│   │   └── module.go             # Returns Module{Handler, Service}
│   │
│   ├── retailer/                 # Pharmacy module
│   │   ├── model.go
│   │   ├── dto.go
│   │   ├── validator.go
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── routes.go
│   │   └── module.go             # Returns Module{Handler, Service}
│   │
│   ├── medicine/                 # Medicine catalog
│   │   ├── model.go              # Medicine (id, name, created_at)
│   │   ├── repository.go         # BatchInsert, GetMedicinesByNames, SearchMedicines (pg_trgm)
│   │   ├── parser.go             # ParseCSV, ParsePDF (ledongthuc/pdf)
│   │   ├── service.go            # Search, batch seed logic
│   │   ├── handler.go            # GET /medicines?q=...
│   │   ├── routes.go
│   │   └── module.go             # Returns Module{Handler, Service}
│   │
│   ├── inventory/                # Stockist-medicine join
│   │   ├── model.go              # Inventory (stockist_id, medicine_id, created_at)
│   │   ├── repository.go         # BulkCreate, FindStockistsByMedicineID
│   │   ├── service.go
│   │   ├── handler.go            # GET /inventory/stockists?medicine_id=X
│   │   ├── routes.go
│   │   └── module.go             # Returns Module{Handler, Service}
│   │
│   ├── job/                      # Background job processing
│   │   ├── model.go              # Job with status enum
│   │   ├── repository.go         # CreateJob, FetchPendingJobs, UpdateJobStatus, ResetStaleJobs
│   │   ├── service.go            # CreateJob, ProcessPendingJobs (resets stale jobs first)
│   │   └── processor.go          # Parse file → seed medicines → link inventory
│   │
│   ├── upload/                   # File upload
│   │   ├── handler.go            # POST /upload (multipart)
│   │   ├── service.go            # Validate file type, save to disk, create job
│   │   └── routes.go
│   │
│   ├── ui/                       # Browser testing interface (HTMX + Alpine.js)
│   │   ├── handler.go            # Page renderers, form handlers for all CRUD
│   │   ├── renderer.go           # Per-page isolated template engine (clones)
│   │   ├── module.go             # DI wiring → returns Module{Handler, Renderer}
│   │   ├── routes.go             # Root-level routes (/, /login, /stockists, …)
│   │   └── templates/
│   │       ├── layout.gohtml     # HTML shell with nav, Alpine.js, HTMX
│   │       ├── partials.gohtml   # Shared partials (lists, forms)
│   │       ├── login.gohtml
│   │       ├── dashboard.gohtml
│   │       ├── stockists.gohtml
│   │       ├── retailers.gohtml
│   │       ├── medicines.gohtml
│   │       ├── inventory.gohtml
│   │       └── upload.gohtml
│   │
│   └── router/router.go          # Route registration hub, middleware per group
│
├── migrations/                   # 6 migration pairs (up/down)
├── docs/                         # Documentation
│   ├── SYSTEM_DESIGN.md
│   ├── ARCHITECTURE.md
│   └── openapi.yaml
│
├── docker-compose.yml
├── .env.example
├── go.mod
└── go.sum
```

---

## Module Pattern

Every feature module follows this structure:

```
module/
├── model.go        # Domain types (pure data, no tags)
├── dto.go          # Request/response DTOs with validation tags
├── repository.go   # Interface + implementation (raw SQL)
├── service.go      # Interface + implementation (business logic)
├── handler.go      # HTTP handler methods
├── routes.go       # Route registration on an echo.Group
├── validator.go    # Shared validator instance (if needed)
└── module.go       # DI wiring → returns Module{Handler, Service}
```

All modules return a `Module` struct:

```go
type Module struct {
    Handler *Handler
    Service Service
}
```

The UI module is an exception — it returns:

```go
type Module struct {
    Handler  *Handler          # Page/form HTTP handlers
    Renderer *TemplateRenderer # Echo v5 Renderer (registered as e.Renderer)
}
```

---

## Middleware Pipeline

### Global Middleware (applied to all routes)

```
RequestID (outermost) → Logger → Recovery
```

| Middleware | Order | Purpose |
|---|---|---|
| **RequestID** | 1st | Injects/tracks `X-Request-ID` header |
| **Logger** | 2nd | Logs method, path, status, latency, client IP, user-agent via Zap |
| **Recovery** | 3rd | Catches panics, returns 500 instead of crashing |

### Per-Group Middleware (applied to API route groups)

| Group | Middleware |
|---|---|
| `/api/v1/auth/login` | — (public) |
| `/api/v1/auth/register` | — (public) |
| `/api/v1/health` | — (public) |
| `/api/v1/auth/admin/*` | `AuthRequired` + `RequireRole("admin")` |
| `/api/v1/medicines` | `AuthRequired` |
| `/api/v1/inventory` | `AuthRequired` |
| `/api/v1/stockists` | `AuthRequired` + `RequireRole("admin")` |
| `/api/v1/retailers` | `AuthRequired` + `RequireRole("admin")` |
| `/api/v1/upload` | `AuthRequired` + `RequireRole("stockist")` |

### UI Routes (root-level, browser-facing)

All UI routes are **public** (no auth middleware). The login form stores the JWT in `localStorage` for subsequent API calls via HTMX. Route groups use the same handlers as the API module where applicable.

---

## Template Rendering

The UI module uses Go's `html/template` with **per-page template isolation** to avoid name collisions:

```
Shared Base (layout + partials)
  ├── layout.gohtml         → {{define "layout"}} ... {{block "content" .}}{{end}} ... {{end}}
  └── partials.gohtml       → {{define "stockists_list"}} ... , {{define "stockist_form"}} ...

Per-Page Clone (shared + page file)
  ├── login.gohtml          → cloned from base, page's {{define "content"}} isolated
  ├── dashboard.gohtml
  ├── stockists.gohtml      → references {{template "stockists_list" .}} from partials
  └── ...
```

- **Page requests** → execute `"layout"` from the page's clone (finds its own `"content"`)
- **HTMX partial requests** → execute the named partial from the shared base set

---

## Auth Flow

### JWT Token Structure (HS256)

```json
{
  "user_id": 1,
  "email": "admin@example.com",
  "role": "admin",
  "reference_id": 0
}
```

Claims are extracted in `AuthRequired` middleware and set in `echo.Context`:
- `user_id` — `auth.GetUserID(c)`
- `user_role` — `auth.GetUserRole(c)`
- `reference_id` — `auth.GetReferenceID(c)` (FK to stockists.id or retailers.id)

### Role-Based Access

| Role | Routes Accessible |
|---|---|
| `admin` | Everything (stockist CRUD, retailer CRUD, admin stockist creation) |
| `stockist` | Medicine search, inventory lookup, file upload |
| `retailer` | Medicine search, inventory lookup |

### User Creation

| User Type | Created By | Endpoint |
|---|---|---|
| **Admin** | Seed on startup | `AUTH_ADMIN_*` env vars (uses `username` OR `email`) |
| **Stockist** | Admin | `POST /auth/admin/stockists` |
| **Retailer** | Self-registration | `POST /auth/register` |

---

## Background Worker

### Architecture

```
┌────────────────────────────────────────────────────┐
│ cmd/worker/main.go                                 │
│   └── job.Processor                                │
│         ├── ParseCSV / ParsePDF                    │
│         ├── medicineRepo.BatchInsert                │
│         ├── medicineRepo.GetMedicinesByNames        │
│         └── inventoryRepo.BulkCreate                │
└────────────────────────────────────────────────────┘
        │ Polls every 10s
        ▼
┌────────────────────────────────────────────────────┐
│ PostgreSQL (jobs table with status = 'pending')    │
└────────────────────────────────────────────────────┘
```

### Job State Machine

```
pending ──► processing ──► completed
               │
               ▼
            failed
```

### Processing Cycle

1. `ResetStaleJobs` — jobs stuck in `processing` for >5 minutes are reset back to `pending`
2. Fetch up to 5 `pending` jobs ordered by `created_at ASC`
3. For each job:
   - Mark as `processing` (set `started_at`)
   - Parse file based on extension (`.csv` → CSV parser, `.pdf` → PDF parser)
   - `BatchInsert` all medicine names from file (ON CONFLICT DO NOTHING)
   - `GetMedicinesByNames` to get ID map
   - `BulkCreate` inventory entries (ON CONFLICT DO NOTHING)
   - Mark as `completed` (set `completed_at`)
   - On error: mark as `failed` (set `error_message`)

---

## Error Handling

### Sentinel Errors

| Error | HTTP Status | When |
|---|---|---|
| `ErrNotFound` | `404 Not Found` | Resource missing |
| `ErrDuplicateEmail` | `409 Conflict` | Duplicate email on create |
| `ErrInvalidCredentials` | `401 Unauthorized` | Bad login |
| `ErrDuplicateUsername` | `409 Conflict` | Duplicate username |
| Rate limit | `429 Too Many Requests` | >100 req/5min per IP |

### Response Envelope

**Success:**
```json
{ "success": true, "data": { ... } }
```

**Error:**
```json
{ "success": false, "error": "descriptive message" }
```

**Paginated:**
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

---

## Configuration

All configuration is environment-driven. The `config.LoadConfig()` function reads from `os.Getenv` with sensible defaults.

| Variable | Required | Default | Purpose |
|---|---|---|---|
| `APP_PORT` | No | `8080` | HTTP listen port |
| `APP_ENV` | No | `development` | Log format |
| `DB_HOST` | No | `localhost` | PostgreSQL host |
| `DB_PORT` | No | `5432` | PostgreSQL port |
| `DB_USER` | No | `postgres` | Database user |
| `DB_PASSWORD` | No | `postgres` | Database password |
| `DB_NAME` | No | `pharmastock-db` | Database name |
| `DB_SSL_MODE` | No | `disable` | SSL mode |
| `JWT_SECRET` | **Yes** | — | JWT signing key |
| `UPLOAD_DIR` | No | `./uploads` | File upload directory |
| `ADMIN_USERNAME` | No | `admin` | Default admin username |
| `ADMIN_PASSWORD` | **Yes** | — | Admin password |
| `ADMIN_EMAIL` | **Yes** | — | Admin email |

Pool settings (hardcoded): min 2, max 10 connections, 1h max lifetime, 30m max idle time.

---

## Graceful Shutdown

```go
func (a *App) Start(ctx context.Context) error {
    sc := echo.StartConfig{
        Address:         ":" + a.Config.AppPort,
        GracefulTimeout: 10 * time.Second,
    }
    return sc.Start(ctx, a.Echo)
}
```

On SIGINT/SIGTERM:
1. Echo stops accepting new connections
2. Existing requests have 10s to complete
3. Logger is flushed
4. Connection pool is closed

---

## Performance Considerations

- **pgx connection pooling** — min 2, max 10 connections with 1h max lifetime
- **pg_trgm GIN index** — fast fuzzy text search on medicine names
- **Rate limiting** — 100 requests per 5 min per IP prevents abuse
- **Batch inserts** — `BatchInsert` and `BulkCreate` use single round-trip SQL
- **ON CONFLICT DO NOTHING** — idempotent, no error handling needed for duplicates
- **Polling interval** — 10s is tunable; suitable for moderate upload volumes
- **No N+1 queries** — all lookups fetch complete result sets
- **Template isolation** — per-page template clones prevent `{{define}}` name collisions without runtime overhead
