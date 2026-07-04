package inventory

import (
	"context"
	"fmt"
)

type Service interface {
	BulkCreate(ctx context.Context, entries []Entry) error
	FindStockistsByMedicineID(ctx context.Context, medicineID int64) ([]*StockistInfo, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) BulkCreate(ctx context.Context, entries []Entry) error {
	if len(entries) == 0 {
		return nil
	}

	if err := s.repo.BulkCreate(ctx, entries); err != nil {
		return fmt.Errorf("failed to bulk create inventory entries: %w", err)
	}

	return nil
}

func (s *service) FindStockistsByMedicineID(ctx context.Context, medicineID int64) ([]*StockistInfo, error) {
	stockists, err := s.repo.FindStockistsByMedicineID(ctx, medicineID)
	if err != nil {
		return nil, fmt.Errorf("failed to find stockists: %w", err)
	}

	return stockists, nil
}
