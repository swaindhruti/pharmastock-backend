package stockist

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateStockist(c *echo.Context) error {

	var stockist Stockist

	if err := c.Bind(&stockist); err != nil {
		return err
	}

	if err := ValidateStockist(&stockist); err != nil {
		return err
	}

	return h.service.CreateStockist(c.Request().Context(), &stockist)
}

func (h *Handler) GetStockistByEmail(c *echo.Context) error {

	email := c.QueryParam("email")
	if email == "" {
		return echo.NewHTTPError(400, "email query parameter is required")
	}

	stockist, err := h.service.GetStockistByEmail(c.Request().Context(), email)
	if err != nil {
		return echo.NewHTTPError(404, "stockist not found")
	}

	return c.JSON(200, stockist)
}

func (h *Handler) UpdateStockist(c *echo.Context) error {

	var stockist Stockist

	if err := c.Bind(&stockist); err != nil {
		return err
	}

	if err := ValidateStockist(&stockist); err != nil {
		return err
	}

	return h.service.UpdateStockist(c.Request().Context(), &stockist)
}

func (h *Handler) DeleteStockist(c *echo.Context) error {

	email := c.QueryParam("email")

	if email == "" {
		return echo.NewHTTPError(400, "email query parameter is required")
	}

	idParam := c.QueryParam("id")
	if idParam == "" {
		return echo.NewHTTPError(400, "id query parameter is required")
	}

	var id int64
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		return echo.NewHTTPError(400, "invalid id query parameter")
	}

	return h.service.DeleteStockist(c.Request().Context(), email, id)
}

func (h *Handler) ListStockists(c *echo.Context) error {

	stockists, err := h.service.ListStockists(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "failed to list stockists")
	}

	return c.JSON(200, stockists)
}
