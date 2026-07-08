package retailer

import (
	"context"
	"errors"
	"fmt"
)

var ErrDuplicateEmail = errors.New("retailer with this email already exists")
var ErrNotFound = errors.New("retailer not found")

type PaginatedRetailers struct {
	Items      []*Retailer `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

type Service interface {
	CreateRetailer(ctx context.Context, retailer *Retailer) error
	GetRetailerByEmail(ctx context.Context, email string) (*Retailer, error)
	GetRetailerByID(ctx context.Context, id int64) (*Retailer, error)
	UpdateRetailer(ctx context.Context, retailer *Retailer) error
	DeleteRetailer(ctx context.Context, id int64) error
	ListRetailers(ctx context.Context, page, limit int) (*PaginatedRetailers, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateRetailer(ctx context.Context, retailer *Retailer) error {
	existing, err := s.repo.GetRetailerByEmail(ctx, retailer.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("failed to check existing retailer: %w", err)
	}
	if existing != nil {
		return ErrDuplicateEmail
	}

	if err := s.repo.CreateRetailer(ctx, retailer); err != nil {
		return fmt.Errorf("failed to create retailer: %w", err)
	}
	return nil
}

func (s *service) GetRetailerByEmail(ctx context.Context, email string) (*Retailer, error) {
	retailer, err := s.repo.GetRetailerByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("retailer not found: %w", err)
	}
	return retailer, nil
}

func (s *service) GetRetailerByID(ctx context.Context, id int64) (*Retailer, error) {
	retailer, err := s.repo.GetRetailerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retailer not found: %w", err)
	}
	return retailer, nil
}

func (s *service) UpdateRetailer(ctx context.Context, retailer *Retailer) error {
	if err := s.repo.UpdateRetailer(ctx, retailer); err != nil {
		return fmt.Errorf("failed to update retailer: %w", err)
	}
	return nil
}

func (s *service) DeleteRetailer(ctx context.Context, id int64) error {
	if err := s.repo.DeleteRetailer(ctx, id); err != nil {
		return fmt.Errorf("failed to delete retailer: %w", err)
	}
	return nil
}

func (s *service) ListRetailers(ctx context.Context, page, limit int) (*PaginatedRetailers, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	retailers, err := s.repo.ListRetailers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list retailers: %w", err)
	}

	total, err := s.repo.CountRetailers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count retailers: %w", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &PaginatedRetailers{
		Items:      retailers,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
