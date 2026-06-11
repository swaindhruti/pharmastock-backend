package retailer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateRetailer(ctx context.Context, retailer *Retailer) error
	GetRetailerByEmail(ctx context.Context, email string) (*Retailer, error)
	UpdateRetailer(ctx context.Context, retailer *Retailer) error
	DeleteRetailer(ctx context.Context, id int64) error
	ListRetailers(ctx context.Context, limit, offset int) ([]*Retailer, error)
	CountRetailers(ctx context.Context) (int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRetailer(ctx context.Context, retailer *Retailer) error {

	query := `INSERT INTO retailers (owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			  RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		retailer.OwnerName, retailer.BusinessName, retailer.Email, retailer.Phone,
		retailer.Country, retailer.State, retailer.City, retailer.PinCode,
		retailer.Address, retailer.GSTNumber,
	).Scan(&retailer.RetailerID, &retailer.CreatedAt, &retailer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create retailer: %w", err)
	}

	return nil
}

func (r *repository) GetRetailerByEmail(ctx context.Context, email string) (*Retailer, error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number, created_at, updated_at
			  FROM retailers WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)

	retailer := &Retailer{}
	err := row.Scan(&retailer.RetailerID, &retailer.OwnerName, &retailer.BusinessName, &retailer.Email,
		&retailer.Phone, &retailer.Country, &retailer.State, &retailer.City, &retailer.PinCode,
		&retailer.Address, &retailer.GSTNumber, &retailer.CreatedAt, &retailer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("retailer not found: %w", ErrNotFound)
	}

	return retailer, nil
}

func (r *repository) UpdateRetailer(ctx context.Context, retailer *Retailer) error {

	query := `UPDATE retailers SET owner_name = $1, business_name = $2, phone = $3, country = $4, state = $5,
			  city = $6, pin_code = $7, address = $8, gst_number = $9, updated_at = NOW() WHERE id = $10
			  RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		retailer.OwnerName, retailer.BusinessName, retailer.Phone,
		retailer.Country, retailer.State, retailer.City, retailer.PinCode,
		retailer.Address, retailer.GSTNumber, retailer.RetailerID,
	).Scan(&retailer.CreatedAt, &retailer.UpdatedAt)

	if err != nil {
		return fmt.Errorf("retailer not found: %w", ErrNotFound)
	}

	return nil
}

func (r *repository) DeleteRetailer(ctx context.Context, id int64) error {

	query := `DELETE FROM retailers WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete retailer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) ListRetailers(ctx context.Context, limit, offset int) ([]*Retailer, error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number, created_at, updated_at
			  FROM retailers ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list retailers: %w", err)
	}
	defer rows.Close()

	var retailers []*Retailer
	for rows.Next() {
		retailer := &Retailer{}
		err := rows.Scan(&retailer.RetailerID, &retailer.OwnerName, &retailer.BusinessName, &retailer.Email,
			&retailer.Phone, &retailer.Country, &retailer.State, &retailer.City, &retailer.PinCode,
			&retailer.Address, &retailer.GSTNumber, &retailer.CreatedAt, &retailer.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan retailer: %w", err)
		}

		retailers = append(retailers, retailer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating retailers: %w", err)
	}

	return retailers, nil
}

func (r *repository) CountRetailers(ctx context.Context) (int64, error) {

	query := `SELECT COUNT(*) FROM retailers`

	var total int64
	err := r.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count retailers: %w", err)
	}

	return total, nil
}
