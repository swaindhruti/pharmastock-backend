package job

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swaindhruti/pharmastock-backend/internal/common"
)

type Repository interface {
	CreateJob(ctx context.Context, stockistID int64, filePath string) (int64, error)
	GetPendingJobs(ctx context.Context, limit int) ([]*Job, error)
	MarkJobProcessing(ctx context.Context, jobID int64) error
	MarkJobCompleted(ctx context.Context, jobID int64) error
	MarkJobFailed(ctx context.Context, jobID int64, errorMessage string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateJob(ctx context.Context, stockistID int64, filePath string) (int64, error) {

	query := `INSERT INTO jobs (stockist_id, job_status, file_path) VALUES ($1, 'pending', $2) RETURNING id`

	var jobID int64

	err := r.db.QueryRow(ctx, query, stockistID, filePath).Scan(&jobID)

	if err != nil {
		return 0, common.NewInternalError(fmt.Sprintf("failed to create job: %v", err))
	}

	return jobID, nil
}

func (r *repository) GetPendingJobs(ctx context.Context, limit int) ([]*Job, error) {

	query := `SELECT id, stockist_id, job_status, file_path, error_message, created_at, started_at, completed_at FROM jobs WHERE job_status = 'pending' ORDER BY created_at LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, common.NewInternalError(fmt.Sprintf("failed to get pending jobs: %v", err))
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		job := &Job{}

		err := rows.Scan(&job.ID, &job.StockistID, &job.JobStatus, &job.FilePath,
			&job.ErrorMessage, &job.CreatedAt, &job.StartedAt, &job.CompletedAt)

		if err != nil {
			return nil, common.NewInternalError(fmt.Sprintf("failed to scan job: %v", err))
		}

		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, common.NewInternalError(fmt.Sprintf("error iterating over job rows: %v", err))
	}

	return jobs, nil
}

func (r *repository) MarkJobProcessing(ctx context.Context, jobID int64) error {

	query := `UPDATE jobs SET job_status = 'processing', started_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(ctx, query, jobID)

	if err != nil {
		return common.NewInternalError(fmt.Sprintf("failed to mark job as processing: %v", err))
	}

	if result.RowsAffected() == 0 {
		return common.NewNotFoundError(common.ErrJobNotFound)
	}

	return nil
}

func (r *repository) MarkJobCompleted(ctx context.Context, jobID int64) error {

	query := `UPDATE jobs SET job_status = 'completed', completed_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(ctx, query, jobID)
	if err != nil {
		return common.NewInternalError(fmt.Sprintf("failed to mark job as completed: %v", err))
	}

	if result.RowsAffected() == 0 {
		return common.NewNotFoundError(common.ErrJobNotFound)
	}

	return nil
}

func (r *repository) MarkJobFailed(ctx context.Context, jobID int64, errorMessage string) error {

	query := `UPDATE jobs SET job_status = 'failed', error_message = $2, completed_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(ctx, query, jobID, errorMessage)

	if err != nil {
		return common.NewInternalError(fmt.Sprintf("failed to mark job as failed: %v", err))
	}

	if result.RowsAffected() == 0 {
		return common.NewNotFoundError(common.ErrJobNotFound)
	}

	return nil
}
