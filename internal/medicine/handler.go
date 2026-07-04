package medicine

import (
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

func (h *Handler) SearchMedicines(c *echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return common.APIErrorResponse(c, http.StatusBadRequest, "search query is required")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 50 {
		limit = 20
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	medicines, err := h.service.SearchMedicines(c.Request().Context(), query, limit, offset)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to search medicines")
	}

	return common.APISuccessResponse(c, http.StatusOK, medicines)
}
