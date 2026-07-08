package ui

import (
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	e.GET("/login", h.LoginPage)
	e.POST("/login", h.Login)
	e.GET("/", h.Dashboard)

	e.GET("/stockists", h.StockistsPage)
	e.POST("/stockists", h.StockistCreate)
	e.GET("/stockists/create", h.StockistEditForm)
	e.GET("/stockists/:id/edit", h.StockistEditForm)
	e.POST("/stockists/:id/update", h.StockistUpdate)
	e.PUT("/stockists/:id", h.StockistUpdate)
	e.DELETE("/stockists/:id", h.StockistDelete)

	e.GET("/retailers", h.RetailersPage)
	e.POST("/retailers", h.RetailerCreate)
	e.GET("/retailers/create", h.RetailerEditForm)
	e.GET("/retailers/:id/edit", h.RetailerEditForm)
	e.POST("/retailers/:id/update", h.RetailerUpdate)
	e.PUT("/retailers/:id", h.RetailerUpdate)
	e.DELETE("/retailers/:id", h.RetailerDelete)

	e.GET("/medicines", h.MedicinesPage)
	e.GET("/inventory", h.InventoryPage)
	e.GET("/upload", h.UploadPage)
	e.POST("/upload", h.UploadFile)


}
