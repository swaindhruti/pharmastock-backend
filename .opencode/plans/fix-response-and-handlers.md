# Fix: common/response.go + Migrate All Handlers

## Files to modify (7 files)

### 1. `internal/common/response.go` — Implement response helpers

```go
package common

import "github.com/labstack/echo/v5"

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func APISuccessResponse(c *echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

func APISuccessMessage(c *echo.Context, statusCode int, message string, data any) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func APIErrorResponse(c *echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}
```

### 2. `internal/health/health.go` — Migrate to response helpers

Replace the two `c.JSON(...)` branches:
- Healthy (200): `common.APISuccessResponse(c, http.StatusOK, map[string]any{...})`
- Unhealthy (503): `common.APIErrorResponse(c, http.StatusServiceUnavailable, "database is down")`

Add import for `"github.com/swaindhruti/pharmastock-backend/internal/common"`

### 3. `internal/stockist/handler.go` — Full rewrite

Fix bugs:
- **GetStockistByEmail**: change `c.QueryParam("email")` → `c.PathParam("email")`
- **UpdateStockist**: add parsing of `id` from `c.PathParam("id")`, set `stockist.StockistID` before calling service. Validate `id` is non-zero.
- **DeleteStockist**: change `c.QueryParam("id")` → `c.PathParam("id")`, remove `email` query param requirement (use path param `id` only; email can be additional verification or removed)

Migrate all responses to `common.*Response`:
- `CreateStockist` (201): `common.APISuccessResponse(c, http.StatusCreated, stockist)`
- `GetStockistByEmail` (200): `common.APISuccessResponse(c, http.StatusOK, stockist)`
- `UpdateStockist` (200): `common.APISuccessResponse(c, http.StatusOK, stockist)`
- `DeleteStockist` (200): `common.APISuccessMessage(c, http.StatusOK, "stockist deleted successfully", nil)`
- `ListStockists` (200): `common.APISuccessResponse(c, http.StatusOK, stockists)`

Error handling: wrap raw errors with `common.APIErrorResponse` — never let DB errors reach the client.

Add `"net/http"` and `"github.com/swaindhruti/pharmastock-backend/internal/common"` imports.

### 4. `internal/stockist/repository.go` — Fix UpdateStockist

Current (broken):
```go
stockistExisting, err := r.GetStockistByEmail(ctx, stockist.Email)
if err != nil { return err }
if stockistExisting != nil { return errors.New("stockist with this email already exists") }
// UPDATE using email as WHERE
```

Fix:
```go
query := `UPDATE stockists SET owner_name=$1, business_name=$2, phone=$3, country=$4, state=$5,
          city=$6, pin_code=$7, address=$8, gst_number=$9 WHERE id=$10`
result, err := r.db.Exec(ctx, query,
    stockist.OwnerName, stockist.BusinessName, stockist.Phone,
    stockist.Country, stockist.State, stockist.City, stockist.PinCode,
    stockist.Address, stockist.GSTNumber, stockist.StockistID)
if err != nil { return err }
if result.RowsAffected() == 0 { return errors.New("stockist not found") }
return nil
```

### 5. `internal/stockist/model.go` — Fix typo

`BuisnessName` → `BusinessName` (field name, JSON tag, and validate tag)

### 6. Propagate `BusinessName` rename

- `internal/stockist/repository.go`: all `.BuisnessName` → `.BusinessName`
- Also fix `phone_number` → `phone` in the SQL INSERT query (line 28 currently uses `phone` which matches DB column, but migration says `phone_number` — keep as `phone` to match Go struct; migration will be fixed separately)

### 7. `internal/middleware/rate_limit.go` — Use common response

Replace inline `map[string]any` in `ErrorHandler` and `DenyHandler`:
- `common.APIErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")`
- `common.APIErrorResponse(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")`

Add import for common package.
