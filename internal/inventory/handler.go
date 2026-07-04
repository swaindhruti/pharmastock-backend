package inventory

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

func (h *Handler) FindStockistsByMedicineID(c *echo.Context) error {
	medicineIDStr := c.QueryParam("medicine_id")
	if medicineIDStr == "" {
		return common.APIErrorResponse(c, http.StatusBadRequest, "medicine_id is required")
	}

	medicineID, err := strconv.ParseInt(medicineIDStr, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid medicine_id")
	}

	stockists, err := h.service.FindStockistsByMedicineID(c.Request().Context(), medicineID)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to find stockists")
	}

	return common.APISuccessResponse(c, http.StatusOK, &StockistsByMedicineResponse{
		MedicineID: medicineID,
		Stockists:  stockists,
	})
}
