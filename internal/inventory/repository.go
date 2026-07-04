package inventory

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	BulkCreate(ctx context.Context, entries []Entry) error
	FindStockistsByMedicineID(ctx context.Context, medicineID int64) ([]*StockistInfo, error)
}

type repository struct {
	db *pgxpool.Pool
}

type Entry struct {
	StockistID int64
	MedicineID int64
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) BulkCreate(ctx context.Context, entries []Entry) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO inventories (stockist_id, medicine_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	for _, entry := range entries {
		_, err := tx.Exec(ctx, query, entry.StockistID, entry.MedicineID)
		if err != nil {
			return fmt.Errorf("failed to insert inventory entry: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) FindStockistsByMedicineID(ctx context.Context, medicineID int64) ([]*StockistInfo, error) {
	query := `SELECT s.id, s.owner_name, s.business_name, s.city, s.state
			  FROM inventories i
			  JOIN stockists s ON s.id = i.stockist_id
			  WHERE i.medicine_id = $1
			  ORDER BY s.business_name`

	rows, err := r.db.Query(ctx, query, medicineID)
	if err != nil {
		return nil, fmt.Errorf("failed to find stockists by medicine: %w", err)
	}
	defer rows.Close()

	var stockists []*StockistInfo
	for rows.Next() {
		info := &StockistInfo{}
		err := rows.Scan(&info.StockistID, &info.OwnerName, &info.BusinessName, &info.City, &info.State)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stockist info: %w", err)
		}
		stockists = append(stockists, info)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stockist rows: %w", err)
	}

	return stockists, nil
}
