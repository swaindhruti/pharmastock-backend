package health

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/swaindhruti/pharmastock-backend/internal/common"
	"github.com/swaindhruti/pharmastock-backend/internal/database"
)

type Handler struct {
	DB *database.PostgresDB
}

func NewHandler(db *database.PostgresDB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) HealthCheck(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.DB.Health(ctx); err != nil {
		return common.APIErrorResponse(c, http.StatusServiceUnavailable, "database is down")
	}

	return common.APISuccessResponse(c, http.StatusOK, map[string]any{
		"status": "healthy",
		"checks": map[string]string{
			"api":      "up",
			"database": "up",
		},
	})
}
