package router

import (
	"github.com/labstack/echo/v5"
	"github.com/swaindhruti/pharmastock-backend/internal/auth"
	"github.com/swaindhruti/pharmastock-backend/internal/health"
	"github.com/swaindhruti/pharmastock-backend/internal/inventory"
	"github.com/swaindhruti/pharmastock-backend/internal/medicine"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
	"github.com/swaindhruti/pharmastock-backend/internal/upload"
)

type Handlers struct {
	Auth      *auth.Handler
	Health    *health.Handler
	Stockist  *stockist.Handler
	Retailer  *retailer.Handler
	Medicine  *medicine.Handler
	Inventory *inventory.Handler
	Upload    *upload.Handler
}

func RegisterRoutes(e *echo.Echo, handlers *Handlers, jwtSecret string) {
	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Public endpoints
	v1.GET("/health", handlers.Health.HealthCheck)

	// Auth endpoints (login/register are public, admin routes are protected)
	authGroup := v1.Group("/auth")
	auth.RegisterRoutes(authGroup, handlers.Auth, jwtSecret)

	// Protected: any authenticated user
	authMW := auth.AuthRequired(jwtSecret)

	// Medicine — any authenticated user
	medicineGroup := v1.Group("/medicines")
	medicineGroup.Use(authMW)
	medicine.RegisterRoutes(medicineGroup, handlers.Medicine)

	// Inventory — any authenticated user
	inventoryGroup := v1.Group("/inventory")
	inventoryGroup.Use(authMW)
	inventory.RegisterRoutes(inventoryGroup, handlers.Inventory)

	// Stockist — admin only
	stockistGroup := v1.Group("/stockists")
	stockistGroup.Use(authMW, auth.RequireRole("admin"))
	stockist.RegisterRoutes(stockistGroup, handlers.Stockist)

	// Retailer — admin only (self-registration is via /auth/register)
	retailerGroup := v1.Group("/retailers")
	retailerGroup.Use(authMW, auth.RequireRole("admin"))
	retailer.RegisterRoutes(retailerGroup, handlers.Retailer)

	// Upload — stockist only
	uploadGroup := v1.Group("/upload")
	uploadGroup.Use(authMW, auth.RequireRole("stockist"))
	upload.RegisterRoutes(uploadGroup, handlers.Upload)
}
