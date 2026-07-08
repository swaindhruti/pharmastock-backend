package ui

import (
	"github.com/swaindhruti/pharmastock-backend/internal/auth"
	"github.com/swaindhruti/pharmastock-backend/internal/inventory"
	"github.com/swaindhruti/pharmastock-backend/internal/medicine"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
	"github.com/swaindhruti/pharmastock-backend/internal/upload"
)

type Module struct {
	Handler *Handler
}

func NewModule(
	stockistSvc stockist.Service,
	retailerSvc retailer.Service,
	medicineSvc medicine.Service,
	inventorySvc inventory.Service,
	authSvc auth.Service,
	uploadSvc upload.Service,
) *Module {
	renderer := NewTemplateRenderer("internal/ui/templates/*.gohtml")
	h := NewHandler(renderer, stockistSvc, retailerSvc, medicineSvc, inventorySvc, authSvc, uploadSvc)
	return &Module{Handler: h}
}
