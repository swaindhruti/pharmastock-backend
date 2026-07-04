package retailer

import "github.com/jackc/pgx/v5/pgxpool"

type Module struct {
	Handler *Handler
	Service Service
}

func NewModule(db *pgxpool.Pool) *Module {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)
	return &Module{Handler: handler, Service: svc}
}
