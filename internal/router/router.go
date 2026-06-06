package router

import (
	"github.com/labstack/echo/v5"

	"github.com/swaindhruti/pharmastock-backend/internal/health"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
)

type Handlers struct {
	Health   *health.Handler
	Stockist *stockist.Handler
}

func RegisterRoutes(e *echo.Echo, handlers *Handlers) {
	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Health Check Endpoint
	v1.GET("/health", handlers.Health.HealthCheck)

	// Stockist Endpoints
	stockistGroup := v1.Group("/stockists")
	stockist.RegisterRoutes(stockistGroup, handlers.Stockist)
}
