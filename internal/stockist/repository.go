package stockist

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateStockist(ctx context.Context, stockist *Stockist) error
	GetStockistByEmail(ctx context.Context, email string) (*Stockist, error)
	UpdateStockist(ctx context.Context, stockist *Stockist) error
	DeleteStockist(ctx context.Context, email string, id int64) error
	ListStockists(ctx context.Context) (stockists []*Stockist, err error)
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

	err := r.db.QueryRow(ctx, query, stockist.OwnerName, stockist.BuisnessName, stockist.Email, stockist.Phone,
		stockist.Country, stockist.State, stockist.City, stockist.PinCode, stockist.Address, stockist.GSTNumber).Scan(&stockist.StockistID)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetStockistByEmail(ctx context.Context, email string) (*Stockist, error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number
			  FROM stockists WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)

	stockist := &Stockist{}
	err := row.Scan(&stockist.StockistID, &stockist.OwnerName, &stockist.BuisnessName, &stockist.Email,
		&stockist.Phone, &stockist.Country, &stockist.State, &stockist.City, &stockist.PinCode,
		&stockist.Address, &stockist.GSTNumber)

	if err != nil {
		return nil, errors.New("stockist not found")
	}

	return stockist, nil
}

func (r *repository) UpdateStockist(ctx context.Context, stockist *Stockist) error {

	stockistExisting, err := r.GetStockistByEmail(ctx, stockist.Email)
	if err != nil {
		return err
	}

	if stockistExisting != nil {
		return errors.New("stockist with this email already exists")
	}

	query := `UPDATE stockists SET owner_name = $1, business_name = $2, phone = $3, country = $4, state = $5,
			  city = $6, pin_code = $7, address = $8, gst_number = $9 WHERE email = $10`

	_, err = r.db.Exec(ctx, query, stockist.OwnerName, stockist.BuisnessName, stockist.Phone,
		stockist.Country, stockist.State, stockist.City, stockist.PinCode,
		stockist.Address, stockist.GSTNumber, stockist.Email)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteStockist(ctx context.Context, email string, id int64) error {

	stockistExisting, err := r.GetStockistByEmail(ctx, email)
	if err != nil {
		return err
	}

	if stockistExisting == nil {
		return errors.New("stockist not found")
	}

	query := `DELETE FROM stockists WHERE email = $1 AND id = $2`

	_, err = r.db.Exec(ctx, query, email, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) ListStockists(ctx context.Context) (stockists []*Stockist, err error) {

	query := `SELECT id, owner_name, business_name, email, phone, country, state, city, pin_code, address, gst_number
			  FROM stockists`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		stockist := &Stockist{}
		err := rows.Scan(&stockist.StockistID, &stockist.OwnerName, &stockist.BuisnessName, &stockist.Email,
			&stockist.Phone, &stockist.Country, &stockist.State, &stockist.City, &stockist.PinCode,
			&stockist.Address, &stockist.GSTNumber)

		if err != nil {
			return nil, err
		}

		stockists = append(stockists, stockist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stockists, nil
}
