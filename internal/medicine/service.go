package medicine

import (
	"context"
	"fmt"
)

type Service interface {
	SearchMedicines(ctx context.Context, query string, limit, offset int) ([]*Medicine, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SearchMedicines(ctx context.Context, query string, limit, offset int) ([]*Medicine, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	medicines, err := s.repo.SearchMedicines(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search medicines: %w", err)
	}

	return medicines, nil
}
