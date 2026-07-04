package medicine

import "github.com/jackc/pgx/v5/pgxpool"

func NewModule(db *pgxpool.Pool) *Handler {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)
	return handler
}
