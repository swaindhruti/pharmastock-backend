package upload

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

func (h *Handler) UploadFile(c *echo.Context) error {
	stockistIDStr := c.FormValue("stockist_id")
	if stockistIDStr == "" {
		return common.APIErrorResponse(c, http.StatusBadRequest, "stockist_id is required")
	}

	stockistID, err := strconv.ParseInt(stockistIDStr, 10, 64)
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid stockist_id")
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "file is required")
	}
	defer file.Close()

	jobID, err := h.service.ProcessUpload(c.Request().Context(), stockistID, file, header)
	if err != nil {
		if errors.Is(err, ErrInvalidFileType) {
			return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to process upload")
	}

	return common.APISuccessResponse(c, http.StatusCreated, map[string]int64{"job_id": jobID})
}
