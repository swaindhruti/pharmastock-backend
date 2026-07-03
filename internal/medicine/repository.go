package medicine

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	BatchInsertMedicines(ctx context.Context, medicines []*Medicine) error
	SearchMedicines(ctx context.Context, query string, limit, offset int) ([]*Medicine, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) BatchInsertMedicines(ctx context.Context, medicines []*Medicine) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO medicines (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`

	for _, medicine := range medicines {
		_, err := tx.Exec(ctx, query, medicine.Name)
		if err != nil {
			return fmt.Errorf("failed to insert medicine: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) SearchMedicines(ctx context.Context, query string, limit, offset int) ([]*Medicine, error) {
	sqlQuery := `SELECT id, name, created_at
				 FROM medicines
				 WHERE name ILIKE $1
				 ORDER BY name
				 LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, sqlQuery, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search medicines: %w", err)
	}
	defer rows.Close()

	var medicines []*Medicine
	for rows.Next() {
		medicine := &Medicine{}
		err := rows.Scan(&medicine.ID, &medicine.Name, &medicine.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan medicine: %w", err)
		}
		medicines = append(medicines, medicine)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return medicines, nil
}
