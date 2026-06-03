package health

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
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

	dbStatus := "up"

	if err := h.DB.Health(ctx); err != nil {
		dbStatus = "down"

		return c.JSON(
			http.StatusServiceUnavailable,
			map[string]any{
				"status": "unhealthy",
				"checks": map[string]string{
					"api":      "up",
					"database": dbStatus,
				},
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"status": "healthy",
			"checks": map[string]string{
				"api":      "up",
				"database": dbStatus,
			},
		},
	)
}
