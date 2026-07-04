package stockist

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateStockist(ctx context.Context, stockist *Stockist) error
	GetStockistByID(ctx context.Context, id int64) (*Stockist, error)
	GetStockistByEmail(ctx context.Context, email string) (*Stockist, error)
	UpdateStockist(ctx context.Context, stockist *Stockist) error
	DeleteStockist(ctx context.Context, id int64) error
	ListStockists(ctx context.Context, limit, offset int) ([]*Stockist, error)
	CountStockists(ctx context.Context) (int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateStockist(ctx context.Context, stockist *Stockist) error {

	query := `INSERT INTO stockists (owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err := r.db.QueryRow(ctx, query, stockist.OwnerName, stockist.BusinessName, stockist.Email, stockist.Phone,
		stockist.Country, stockist.State, stockist.City, stockist.PinCode, stockist.Address, stockist.GSTNumber).Scan(&stockist.StockistID)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetStockistByID(ctx context.Context, id int64) (*Stockist, error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number
			  FROM stockists WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	stockist := &Stockist{}
	err := row.Scan(&stockist.StockistID, &stockist.OwnerName, &stockist.BusinessName, &stockist.Email,
		&stockist.Phone, &stockist.Country, &stockist.State, &stockist.City, &stockist.PinCode,
		&stockist.Address, &stockist.GSTNumber)

	if err != nil {
		return nil, fmt.Errorf("stockist not found: %w", ErrNotFound)
	}

	return stockist, nil
}

func (r *repository) GetStockistByEmail(ctx context.Context, email string) (*Stockist, error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number
			  FROM stockists WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)

	stockist := &Stockist{}
	err := row.Scan(&stockist.StockistID, &stockist.OwnerName, &stockist.BusinessName, &stockist.Email,
		&stockist.Phone, &stockist.Country, &stockist.State, &stockist.City, &stockist.PinCode,
		&stockist.Address, &stockist.GSTNumber)

	if err != nil {
		return nil, fmt.Errorf("stockist not found: %w", ErrNotFound)
	}

	return stockist, nil
}

func (r *repository) UpdateStockist(ctx context.Context, stockist *Stockist) error {

	query := `UPDATE stockists SET owner_name = $1, business_name = $2, phone = $3, country = $4, state = $5,
			  city = $6, pin_code = $7, address = $8, gst_number = $9 WHERE id = $10`

	result, err := r.db.Exec(ctx, query, stockist.OwnerName, stockist.BusinessName, stockist.Phone,
		stockist.Country, stockist.State, stockist.City, stockist.PinCode,
		stockist.Address, stockist.GSTNumber, stockist.StockistID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) DeleteStockist(ctx context.Context, id int64) error {

	query := `DELETE FROM stockists WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) ListStockists(ctx context.Context, limit, offset int) ([]*Stockist, error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number
			  FROM stockists ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list stockists: %w", err)
	}
	defer rows.Close()

	var stockists []*Stockist
	for rows.Next() {
		stockist := &Stockist{}
		err := rows.Scan(&stockist.StockistID, &stockist.OwnerName, &stockist.BusinessName, &stockist.Email,
			&stockist.Phone, &stockist.Country, &stockist.State, &stockist.City, &stockist.PinCode,
			&stockist.Address, &stockist.GSTNumber)

		if err != nil {
			return nil, fmt.Errorf("failed to scan stockist: %w", err)
		}

		stockists = append(stockists, stockist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stockists: %w", err)
	}

	return stockists, nil
}

func (r *repository) CountStockists(ctx context.Context) (int64, error) {

	query := `SELECT COUNT(*) FROM stockists`

	var total int64
	err := r.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to count stockists: %w", err)
	}

	return total, nil
}
