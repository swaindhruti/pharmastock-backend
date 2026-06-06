package stockist

import "github.com/jackc/pgx/v5/pgxpool"

func NewModule(db *pgxpool.Pool) *Handler {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return handler
}
