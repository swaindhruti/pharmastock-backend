package stockist

import "github.com/labstack/echo/v5"

func RegisterRoutes(group *echo.Group, h *Handler) {
	group.POST("", h.CreateStockist)
	group.GET("/:email", h.GetStockistByEmail)
	group.PUT("/:id", h.UpdateStockist)
	group.DELETE("/:id", h.DeleteStockist)
	group.GET("", h.ListStockists)
}
