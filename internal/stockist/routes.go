package stockist

import (
	"github.com/labstack/echo/v5"
	"time"

	"github.com/swaindhruti/pharmastock-backend/internal/middleware"
)

func RegisterRoutes(group *echo.Group, h *Handler) {
	group.Use(middleware.RateLimit(100, 5*time.Minute))

	group.POST("", h.CreateStockist)
	group.GET("/:id", h.GetStockistByID)
	group.PUT("/:id", h.UpdateStockist)
	group.DELETE("/:id", h.DeleteStockist)
	group.GET("", h.ListStockists)
}
