package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/swaindhruti/pharmastock-backend/internal/common"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
)

type Handler struct {
	authSvc      Service
	stockistSvc  stockist.Service
	retailerSvc  retailer.Service
	validate     *validator.Validate
}

func NewHandler(authSvc Service, stockistSvc stockist.Service, retailerSvc retailer.Service) *Handler {
	return &Handler{
		authSvc:     authSvc,
		stockistSvc: stockistSvc,
		retailerSvc: retailerSvc,
		validate:    validator.New(),
	}
}

func (h *Handler) SeedAdmin(ctx context.Context, username, password, email string) error {
	return h.authSvc.SeedAdmin(ctx, username, password, email)
}

func (h *Handler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if req.Email == "" && req.Username == "" {
		return common.APIErrorResponse(c, http.StatusBadRequest, "email or username is required")
	}

	if err := h.validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	resp, err := h.authSvc.Login(c.Request().Context(), req.Email, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return common.APIErrorResponse(c, http.StatusUnauthorized, "invalid email or username")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "login failed")
	}

	return common.APISuccessResponse(c, http.StatusOK, resp)
}

func (h *Handler) RegisterRetailer(c *echo.Context) error {
	var req RegisterRetailerRequest
	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	retailerReq := &retailer.CreateRetailerRequest{
		OwnerName:    req.OwnerName,
		BusinessName: req.BusinessName,
		Email:        req.Email,
		Phone:        req.Phone,
		Country:      req.Country,
		State:        req.State,
		City:         req.City,
		PinCode:      req.PinCode,
		Address:      req.Address,
		GSTNumber:    req.GSTNumber,
	}

	r := retailerReq.ToDomain()
	if err := h.retailerSvc.CreateRetailer(c.Request().Context(), r); err != nil {
		if errors.Is(err, retailer.ErrDuplicateEmail) {
			return common.APIErrorResponse(c, http.StatusConflict, "retailer with this email already exists")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to create retailer")
	}

	resp, err := h.authSvc.CreateUser(c.Request().Context(), req.Email, req.Username, req.Password, "retailer", r.RetailerID)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return common.APIErrorResponse(c, http.StatusConflict, "user with this email already exists")
		}
		if errors.Is(err, ErrDuplicateUsername) {
			return common.APIErrorResponse(c, http.StatusConflict, "username already taken")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to create user")
	}

	return common.APISuccessResponse(c, http.StatusCreated, resp)
}

func (h *Handler) AdminCreateStockist(c *echo.Context) error {
	var req CreateStockistUserRequest
	if err := c.Bind(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return common.APIErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	stockistReq := &stockist.CreateStockistRequest{
		OwnerName:    req.OwnerName,
		BusinessName: req.BusinessName,
		Email:        req.Email,
		Phone:        req.Phone,
		Country:      req.Country,
		State:        req.State,
		City:         req.City,
		PinCode:      req.PinCode,
		Address:      req.Address,
		GSTNumber:    req.GSTNumber,
	}

	s := stockistReq.ToDomain()
	if err := h.stockistSvc.CreateStockist(c.Request().Context(), s); err != nil {
		if errors.Is(err, stockist.ErrDuplicateEmail) {
			return common.APIErrorResponse(c, http.StatusConflict, "stockist with this email already exists")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to create stockist")
	}

	resp, err := h.authSvc.CreateUser(c.Request().Context(), req.Email, req.Username, req.Password, "stockist", s.StockistID)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			return common.APIErrorResponse(c, http.StatusConflict, "user with this email already exists")
		}
		if errors.Is(err, ErrDuplicateUsername) {
			return common.APIErrorResponse(c, http.StatusConflict, "username already taken")
		}
		return common.APIErrorResponse(c, http.StatusInternalServerError, "failed to create user")
	}

	return common.APISuccessResponse(c, http.StatusCreated, resp)
}
