package medicine

import "github.com/labstack/echo/v5"

func RegisterRoutes(group *echo.Group, h *Handler) {
	group.GET("", h.SearchMedicines)
}
