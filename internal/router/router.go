package router

import (
	"github.com/labstack/echo/v5"

	"github.com/swaindhruti/pharmastock-backend/internal/health"
)

type Handlers struct {
	Health *health.Handler
}

func RegisterRoutes(e *echo.Echo, handlers *Handlers) {
	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", handlers.Health.HealthCheck)
}
