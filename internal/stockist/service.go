package stockist

import "context"

type Service interface {
	CreateStockist(ctx context.Context, stockist *Stockist) error
	GetStockistByEmail(ctx context.Context, email string) (*Stockist, error)
	UpdateStockist(ctx context.Context, stockist *Stockist) error
	DeleteStockist(ctx context.Context, email string, id int64) error
	ListStockists(ctx context.Context) (stockists []*Stockist, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateStockist(ctx context.Context, stockist *Stockist) error {
	return s.repo.CreateStockist(ctx, stockist)
}

func (s *service) GetStockistByEmail(ctx context.Context, email string) (*Stockist, error) {
	return s.repo.GetStockistByEmail(ctx, email)
}

func (s *service) UpdateStockist(ctx context.Context, stockist *Stockist) error {
	return s.repo.UpdateStockist(ctx, stockist)
}

func (s *service) DeleteStockist(ctx context.Context, email string, id int64) error {
	return s.repo.DeleteStockist(ctx, email, id)
}

func (s *service) ListStockists(ctx context.Context) (stockists []*Stockist, err error) {
	return s.repo.ListStockists(ctx)
}
