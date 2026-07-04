package stockist

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrDuplicateEmail = errors.New("stockist with this email already exists")
	ErrNotFound       = errors.New("stockist not found")
)

type PaginatedStockists struct {
	Items      []*Stockist `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

type Service interface {
	CreateStockist(ctx context.Context, stockist *Stockist) error
	GetStockistByID(ctx context.Context, id int64) (*Stockist, error)
	GetStockistByEmail(ctx context.Context, email string) (*Stockist, error)
	UpdateStockist(ctx context.Context, stockist *Stockist) error
	DeleteStockist(ctx context.Context, id int64) error
	ListStockists(ctx context.Context, page, limit int) (*PaginatedStockists, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateStockist(ctx context.Context, stockist *Stockist) error {
	existing, err := s.repo.GetStockistByEmail(ctx, stockist.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("failed to check existing stockist: %w", err)
	}
	if existing != nil {
		return ErrDuplicateEmail
	}

	if err := s.repo.CreateStockist(ctx, stockist); err != nil {
		return fmt.Errorf("failed to create stockist: %w", err)
	}
	return nil
}

func (s *service) GetStockistByID(ctx context.Context, id int64) (*Stockist, error) {
	stockist, err := s.repo.GetStockistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("stockist not found: %w", err)
	}
	return stockist, nil
}

func (s *service) GetStockistByEmail(ctx context.Context, email string) (*Stockist, error) {
	stockist, err := s.repo.GetStockistByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("stockist not found: %w", err)
	}
	return stockist, nil
}

func (s *service) UpdateStockist(ctx context.Context, stockist *Stockist) error {
	if err := s.repo.UpdateStockist(ctx, stockist); err != nil {
		return fmt.Errorf("failed to update stockist: %w", err)
	}
	return nil
}

func (s *service) DeleteStockist(ctx context.Context, id int64) error {
	if err := s.repo.DeleteStockist(ctx, id); err != nil {
		return fmt.Errorf("failed to delete stockist: %w", err)
	}
	return nil
}

func (s *service) ListStockists(ctx context.Context, page, limit int) (*PaginatedStockists, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	stockists, err := s.repo.ListStockists(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stockists: %w", err)
	}

	total, err := s.repo.CountStockists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count stockists: %w", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &PaginatedStockists{
		Items:      stockists,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
