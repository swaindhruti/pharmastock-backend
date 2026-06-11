package retailer

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/swaindhruti/pharmastock-backend/internal/middleware"
)

func RegisterRoutes(group *echo.Group, h *Handler) {
	group.Use(middleware.RateLimit(100, 5*time.Minute))

	group.POST("", h.CreateRetailer)
	group.GET("/:email", h.GetRetailerByEmail)
	group.PUT("/:id", h.UpdateRetailer)
	group.DELETE("/:id", h.DeleteRetailer)
	group.GET("", h.ListRetailers)
}
