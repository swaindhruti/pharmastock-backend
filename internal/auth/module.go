package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaindhruti/pharmastock-backend/internal/retailer"
	"github.com/swaindhruti/pharmastock-backend/internal/stockist"
)

func NewModule(db *pgxpool.Pool, jwtSecret string, stockistSvc stockist.Service, retailerSvc retailer.Service) *Handler {
	repo := NewRepository(db)
	svc := NewService(repo, jwtSecret)
	handler := NewHandler(svc, stockistSvc, retailerSvc)
	return handler
}
