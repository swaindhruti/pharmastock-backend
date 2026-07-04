package auth

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(group *echo.Group, h *Handler, jwtSecret string) {
	group.POST("/login", h.Login)
	group.POST("/register", h.RegisterRetailer)

	admin := group.Group("/admin")
	admin.Use(AuthRequired(jwtSecret), RequireRole("admin"))
	admin.POST("/stockists", h.AdminCreateStockist)
}
