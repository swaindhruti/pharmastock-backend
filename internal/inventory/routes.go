package inventory

import "github.com/labstack/echo/v5"

func RegisterRoutes(group *echo.Group, h *Handler) {
	group.GET("/stockists", h.FindStockistsByMedicineID)
}
