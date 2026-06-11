package retailer

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

func (h *Handler) CreateRetailer(c *echo.Context) error {

	var req CreateRetailerRequest

	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	retailer := req.ToDomain()

	if err := h.service.CreateRetailer(c.Request().Context(), retailer); err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return common.APIErrorResponse(c, http.StatusConflict, err.Error())
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to create retailer")
	}

	return common.APISuccessResponse(c, http.StatusCreated, NewRetailerResponse(retailer))
}

func (h *Handler) GetRetailerByEmail(c *echo.Context) error {

	email := c.Param("email")
	if email == "" {
		return common.APIErrorResponse(c, http.StatusBadRequest, "email is required")
	}

	retailer, err := h.service.GetRetailerByEmail(c.Request().Context(), email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return common.APIErrorResponse(c, http.StatusNotFound, "retailer not found")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to get retailer")
	}

	return common.APISuccessResponse(c, http.StatusOK, NewRetailerResponse(retailer))
}

func (h *Handler) UpdateRetailer(c *echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid retailer id")
	}

	var req UpdateRetailerRequest

	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	retailer := req.ToDomain(id)

	if err := h.service.UpdateRetailer(c.Request().Context(), retailer); err != nil {
		if errors.Is(err, ErrNotFound) {
			return common.APIErrorResponse(c, http.StatusNotFound, "retailer not found")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to update retailer")
	}

	return common.APISuccessResponse(c, http.StatusOK, NewRetailerResponse(retailer))
}

func (h *Handler) DeleteRetailer(c *echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid retailer id")
	}

	if err := h.service.DeleteRetailer(c.Request().Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return common.APIErrorResponse(c, http.StatusNotFound, "retailer not found")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to delete retailer")
	}

	return common.APISuccessMessage(c, http.StatusOK, "retailer deleted successfully", nil)
}

func (h *Handler) ListRetailers(c *echo.Context) error {

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	result, err := h.service.ListRetailers(c.Request().Context(), page, limit)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to list retailers")
	}

	return common.APISuccessResponse(c, http.StatusOK, NewPaginatedRetailerResponse(result))
}
