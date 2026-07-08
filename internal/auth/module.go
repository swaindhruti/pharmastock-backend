package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
)

type Module struct {
	Handler *Handler
	Service Service
}

func NewModule(db *pgxpool.Pool, jwtSecret string, stockistSvc stockist.Service, retailerSvc retailer.Service) *Module {
	repo := NewRepository(db)
	svc := NewService(repo, jwtSecret)
	handler := NewHandler(svc, stockistSvc, retailerSvc)
	return &Module{Handler: handler, Service: svc}
}
