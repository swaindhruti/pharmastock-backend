package ui

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"github.com/swaindhruti/pharmastock-backend/internal/auth"
	"github.com/swaindhruti/pharmastock-backend/internal/inventory"
	"github.com/swaindhruti/pharmastock-backend/internal/medicine"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
	"github.com/swaindhruti/pharmastock-backend/internal/upload"
)

type Handler struct {
	renderer     *TemplateRenderer
	stockistSvc  stockist.Service
	retailerSvc  retailer.Service
	medicineSvc  medicine.Service
	inventorySvc inventory.Service
	authSvc      auth.Service
	uploadSvc    upload.Service
}

func NewHandler(
	renderer *TemplateRenderer,
	stockistSvc stockist.Service,
	retailerSvc retailer.Service,
	medicineSvc medicine.Service,
	inventorySvc inventory.Service,
	authSvc auth.Service,
	uploadSvc upload.Service,
) *Handler {
	return &Handler{
		renderer:     renderer,
		stockistSvc:  stockistSvc,
		retailerSvc:  retailerSvc,
		medicineSvc:  medicineSvc,
		inventorySvc: inventorySvc,
		authSvc:      authSvc,
		uploadSvc:    uploadSvc,
	}
}

func (h *Handler) render(c *echo.Context, template string, data map[string]any) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().Header().Set("Cache-Control", "no-cache, private")
	return c.Render(http.StatusOK, template, data)
}

func (h *Handler) renderCached(c *echo.Context, template string, data map[string]any) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().Header().Set("Cache-Control", "private, max-age=30")
	return c.Render(http.StatusOK, template, data)
}

func (h *Handler) renderError(c *echo.Context, msg string) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().Header().Set("Cache-Control", "no-cache, private")
	return c.HTML(http.StatusOK, `<div class="alert alert-error">`+msg+`</div>`)
}

func (h *Handler) LoginPage(c *echo.Context) error {
	return h.render(c, "login", map[string]any{"title": "Login"})
}

func (h *Handler) Login(c *echo.Context) error {
	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")

	if password == "" {
		return h.render(c, "login", map[string]any{"title": "Login", "error": "password is required"})
	}
	if email == "" && username == "" {
		return h.render(c, "login", map[string]any{"title": "Login", "error": "email or username is required"})
	}

	resp, err := h.authSvc.Login(c.Request().Context(), email, username, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return h.render(c, "login", map[string]any{"title": "Login", "error": "invalid credentials"})
		}
		return h.render(c, "login", map[string]any{"title": "Login", "error": "login failed"})
	}

	return h.render(c, "login", map[string]any{
		"title":    "Login",
		"token":    resp.Token,
		"user_id":  resp.UserID,
		"role":     resp.Role,
		"ref_id":   resp.ReferenceID,
		"loggedIn": true,
	})
}

func (h *Handler) Dashboard(c *echo.Context) error {
	list, _ := h.stockistSvc.ListStockists(c.Request().Context(), 1, 1)
	rList, _ := h.retailerSvc.ListRetailers(c.Request().Context(), 1, 1)

	return h.render(c, "dashboard", map[string]any{
		"title":         "Dashboard",
		"stockistCount": list.Total,
		"retailerCount": rList.Total,
	})
}

func (h *Handler) StockistsPage(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, err := h.stockistSvc.ListStockists(c.Request().Context(), page, 20)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return h.renderError(c, "failed to load stockists")
		}
		return h.render(c, "stockists", map[string]any{"title": "Stockists", "error": "failed to load stockists"})
	}

	isHX := c.Request().Header.Get("HX-Request") == "true"
	if isHX {
		return h.renderCached(c, "stockists_list", map[string]any{
			"title": "Stockists",
			"items": result.Items,
			"total": result.Total,
			"page":  result.Page,
			"limit": result.Limit,
			"pages": result.TotalPages,
			"isHX":  true,
		})
	}
	return h.render(c, "stockists", map[string]any{
		"title":   "Stockists",
		"items":   result.Items,
		"total":   result.Total,
		"page":    result.Page,
		"limit":   result.Limit,
		"pages":   result.TotalPages,
		"isHX":    false,
	})
}

func (h *Handler) StockistCreate(c *echo.Context) error {
	req := &stockist.CreateStockistRequest{
		OwnerName:    c.FormValue("name"),
		BusinessName: c.FormValue("business_name"),
		Email:        c.FormValue("email"),
		Phone:        c.FormValue("phone"),
		Country:      c.FormValue("country"),
		State:        c.FormValue("state"),
		City:         c.FormValue("city"),
		PinCode:      c.FormValue("pin_code"),
		Address:      c.FormValue("address"),
		GSTNumber:    c.FormValue("gst_number"),
	}
	s := req.ToDomain()
	if err := h.stockistSvc.CreateStockist(c.Request().Context(), s); err != nil {
		result, _ := h.stockistSvc.ListStockists(c.Request().Context(), 1, 20)
		return h.render(c, "stockists_list", map[string]any{
			"items": result.Items, "total": result.Total,
			"page": 1, "limit": 20, "pages": result.TotalPages, "error": err.Error(),
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, _ := h.stockistSvc.ListStockists(c.Request().Context(), page, 20)
	return h.renderCached(c, "stockists_list", map[string]any{
		"items":  result.Items,
		"total":  result.Total,
		"page":   result.Page,
		"limit":  result.Limit,
		"pages":  result.TotalPages,
		"isHX":   true,
	})
}

func (h *Handler) StockistEditForm(c *echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return h.render(c, "stockist_form", map[string]any{"isEdit": false})
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}
	s, err := h.stockistSvc.GetStockistByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "stockist not found")
	}
	return h.render(c, "stockist_form", map[string]any{
		"stockist": s,
		"isEdit":   true,
	})
}

func (h *Handler) StockistUpdate(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}

	req := &stockist.UpdateStockistRequest{
		OwnerName:    c.FormValue("name"),
		BusinessName: c.FormValue("business_name"),
		Email:        c.FormValue("email"),
		Phone:        c.FormValue("phone"),
		Country:      c.FormValue("country"),
		State:        c.FormValue("state"),
		City:         c.FormValue("city"),
		PinCode:      c.FormValue("pin_code"),
		Address:      c.FormValue("address"),
		GSTNumber:    c.FormValue("gst_number"),
	}
	s := req.ToDomain(id)
	if err := h.stockistSvc.UpdateStockist(c.Request().Context(), s); err != nil {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		result, _ := h.stockistSvc.ListStockists(c.Request().Context(), page, 20)
		return h.render(c, "stockists_list", map[string]any{
			"items": result.Items, "total": result.Total,
			"page": result.Page, "limit": result.Limit, "pages": result.TotalPages, "error": err.Error(),
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, _ := h.stockistSvc.ListStockists(c.Request().Context(), page, 20)
	return h.renderCached(c, "stockists_list", map[string]any{
		"items": result.Items, "total": result.Total,
		"page": result.Page, "limit": result.Limit, "pages": result.TotalPages,
		"isHX": true,
	})
}

func (h *Handler) StockistDelete(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}
	h.stockistSvc.DeleteStockist(c.Request().Context(), id)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, _ := h.stockistSvc.ListStockists(c.Request().Context(), page, 20)
	return h.renderCached(c, "stockists_list", map[string]any{
		"items": result.Items, "total": result.Total,
		"page": result.Page, "limit": result.Limit, "pages": result.TotalPages,
		"isHX": true,
	})
}

func (h *Handler) RetailersPage(c *echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, err := h.retailerSvc.ListRetailers(c.Request().Context(), page, 20)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return h.renderError(c, "failed to load retailers")
		}
		return h.render(c, "retailers", map[string]any{"title": "Retailers", "error": "failed to load retailers"})
	}

	isHX := c.Request().Header.Get("HX-Request") == "true"
	if isHX {
		return h.renderCached(c, "retailers_list", map[string]any{
			"title": "Retailers",
			"items": result.Items,
			"total": result.Total,
			"page":  result.Page,
			"limit": result.Limit,
			"pages": result.TotalPages,
			"isHX":  true,
		})
	}
	return h.render(c, "retailers", map[string]any{
		"title":   "Retailers",
		"items":   result.Items,
		"total":   result.Total,
		"page":    result.Page,
		"limit":   result.Limit,
		"pages":   result.TotalPages,
		"isHX":    false,
	})
}

func (h *Handler) RetailerCreate(c *echo.Context) error {
	req := &retailer.CreateRetailerRequest{
		OwnerName:    c.FormValue("name"),
		BusinessName: c.FormValue("business_name"),
		Email:        c.FormValue("email"),
		Phone:        c.FormValue("phone"),
		Country:      c.FormValue("country"),
		State:        c.FormValue("state"),
		City:         c.FormValue("city"),
		PinCode:      c.FormValue("pin_code"),
		Address:      c.FormValue("address"),
		GSTNumber:    c.FormValue("gst_number"),
	}
	r := req.ToDomain()
	if err := h.retailerSvc.CreateRetailer(c.Request().Context(), r); err != nil {
		result, _ := h.retailerSvc.ListRetailers(c.Request().Context(), 1, 20)
		return h.render(c, "retailers_list", map[string]any{
			"items": result.Items, "total": result.Total,
			"page": 1, "limit": 20, "pages": result.TotalPages, "error": err.Error(),
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, _ := h.retailerSvc.ListRetailers(c.Request().Context(), page, 20)
	return h.renderCached(c, "retailers_list", map[string]any{
		"items": result.Items, "total": result.Total,
		"page": result.Page, "limit": result.Limit, "pages": result.TotalPages,
		"isHX": true,
	})
}

func (h *Handler) RetailerEditForm(c *echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return h.render(c, "retailer_form", map[string]any{"isEdit": false})
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}
	r, err := h.retailerSvc.GetRetailerByID(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusNotFound, "retailer not found")
	}
	return h.render(c, "retailer_form", map[string]any{
		"retailer": r,
		"isEdit":   true,
	})
}

func (h *Handler) RetailerUpdate(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}

	req := &retailer.UpdateRetailerRequest{
		OwnerName:    c.FormValue("name"),
		BusinessName: c.FormValue("business_name"),
		Email:        c.FormValue("email"),
		Phone:        c.FormValue("phone"),
		Country:      c.FormValue("country"),
		State:        c.FormValue("state"),
		City:         c.FormValue("city"),
		PinCode:      c.FormValue("pin_code"),
		Address:      c.FormValue("address"),
		GSTNumber:    c.FormValue("gst_number"),
	}
	r := req.ToDomain(id)
	if err := h.retailerSvc.UpdateRetailer(c.Request().Context(), r); err != nil {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page < 1 {
			page = 1
		}
		result, _ := h.retailerSvc.ListRetailers(c.Request().Context(), page, 20)
		return h.render(c, "retailers_list", map[string]any{
			"items": result.Items, "total": result.Total,
			"page": result.Page, "limit": result.Limit, "pages": result.TotalPages, "error": err.Error(),
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, _ := h.retailerSvc.ListRetailers(c.Request().Context(), page, 20)
	return h.renderCached(c, "retailers_list", map[string]any{
		"items": result.Items, "total": result.Total,
		"page": result.Page, "limit": result.Limit, "pages": result.TotalPages,
		"isHX": true,
	})
}

func (h *Handler) RetailerDelete(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}
	h.retailerSvc.DeleteRetailer(c.Request().Context(), id)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	result, _ := h.retailerSvc.ListRetailers(c.Request().Context(), page, 20)
	return h.renderCached(c, "retailers_list", map[string]any{
		"items": result.Items, "total": result.Total,
		"page": result.Page, "limit": result.Limit, "pages": result.TotalPages,
		"isHX": true,
	})
}

func (h *Handler) MedicinesPage(c *echo.Context) error {
	query := c.QueryParam("q")
	var items []*medicine.Medicine
	if query != "" {
		items, _ = h.medicineSvc.SearchMedicines(c.Request().Context(), query, 20, 0)
	}

	isHX := c.Request().Header.Get("HX-Request") == "true"
	tmpl := "medicines"
	if isHX {
		tmpl = "medicines_results"
	}
	return h.render(c, tmpl, map[string]any{
		"title":   "Medicines",
		"query":   query,
		"items":   items,
		"isHX":    isHX,
	})
}

func (h *Handler) InventoryPage(c *echo.Context) error {
	medicineIDStr := c.QueryParam("medicine_id")
	var stockists []*inventory.StockistInfo
	var medID int64
	if medicineIDStr != "" {
		medID, _ = strconv.ParseInt(medicineIDStr, 10, 64)
		if medID > 0 {
			stockists, _ = h.inventorySvc.FindStockistsByMedicineID(c.Request().Context(), medID)
		}
	}

	isHX := c.Request().Header.Get("HX-Request") == "true"
	tmpl := "inventory"
	if isHX {
		tmpl = "inventory_results"
	}
	return h.render(c, tmpl, map[string]any{
		"title":      "Inventory",
		"medicineID": medID,
		"stockists":  stockists,
		"isHX":       isHX,
	})
}

func (h *Handler) UploadPage(c *echo.Context) error {
	return h.render(c, "upload", map[string]any{
		"title": "Upload",
	})
}

func (h *Handler) UploadFile(c *echo.Context) error {
	stockistIDStr := c.FormValue("stockist_id")
	if stockistIDStr == "" {
		if c.Request().Header.Get("HX-Request") == "true" {
			return h.renderError(c, "stockist_id is required")
		}
		return h.render(c, "upload", map[string]any{"title": "Upload", "error": "stockist_id is required"})
	}
	stockistID, err := strconv.ParseInt(stockistIDStr, 10, 64)
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return h.renderError(c, "invalid stockist_id")
		}
		return h.render(c, "upload", map[string]any{"title": "Upload", "error": "invalid stockist_id"})
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		if c.Request().Header.Get("HX-Request") == "true" {
			return h.renderError(c, "file is required")
		}
		return h.render(c, "upload", map[string]any{"title": "Upload", "error": "file is required"})
	}
	defer file.Close()

	jobID, err := h.uploadSvc.ProcessUpload(c.Request().Context(), stockistID, file, header)
	if err != nil {
		msg := err.Error()
		if errors.Is(err, upload.ErrInvalidFileType) {
			msg = err.Error()
		} else {
			msg = "upload failed"
		}
		if c.Request().Header.Get("HX-Request") == "true" {
			return h.renderError(c, msg)
		}
		return h.render(c, "upload", map[string]any{"title": "Upload", "error": msg})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return h.render(c, "upload_card", map[string]any{
			"jobID":   jobID,
			"success": true,
		})
	}

	return h.render(c, "upload", map[string]any{
		"title":   "Upload",
		"jobID":   jobID,
		"success": true,
	})
}
