package stockist

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/swaindhruti/pharmastock-backend/internal/common"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateStockist(c *echo.Context) error {

	var req CreateStockistRequest

	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	stockist := req.ToDomain()

	if err := h.service.CreateStockist(c.Request().Context(), stockist); err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return common.APIErrorResponse(c, http.StatusConflict, err.Error())
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to create stockist")
	}

	return common.APISuccessResponse(c, http.StatusCreated, NewStockistResponse(stockist))
}

func (h *Handler) GetStockistByID(c *echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid stockist id")
	}

	stockist, err := h.service.GetStockistByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return common.APIErrorResponse(c, http.StatusNotFound, "stockist not found")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to get stockist")
	}

	return common.APISuccessResponse(c, http.StatusOK, NewStockistResponse(stockist))
}

func (h *Handler) UpdateStockist(c *echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid stockist id")
	}

	var req UpdateStockistRequest

	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	stockist := req.ToDomain(id)

	if err := h.service.UpdateStockist(c.Request().Context(), stockist); err != nil {
		if errors.Is(err, ErrNotFound) {
			return common.APIErrorResponse(c, http.StatusNotFound, "stockist not found")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to update stockist")
	}

	return common.APISuccessResponse(c, http.StatusOK, NewStockistResponse(stockist))
}

func (h *Handler) DeleteStockist(c *echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid stockist id")
	}

	if err := h.service.DeleteStockist(c.Request().Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return common.APIErrorResponse(c, http.StatusNotFound, "stockist not found")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to delete stockist")
	}

	return common.APISuccessMessage(c, http.StatusOK, "stockist deleted successfully", nil)
}

func (h *Handler) ListStockists(c *echo.Context) error {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	result, err := h.service.ListStockists(c.Request().Context(), page, limit)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to list stockists")
	}

	return common.APISuccessResponse(c, http.StatusOK, NewPaginatedStockistResponse(result))
}
