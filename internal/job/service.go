package job

import "context"

type Service interface {
	CreateJob(ctx context.Context, stockistID int64, filePath string) (int64, error)
	ProcessPendingJobs(ctx context.Context, limit int) error
}

type service struct {
	repo      Repository
	processor Processor
}

func NewService(repo Repository, processor Processor) Service {
	return &service{repo: repo, processor: processor}
}

func (s *service) CreateJob(ctx context.Context, stockistID int64, filePath string) (int64, error) {
	return s.repo.CreateJob(ctx, stockistID, filePath)
}

func (s *service) ProcessPendingJobs(ctx context.Context, limit int) error {
	jobs, err := s.repo.GetPendingJobs(ctx, limit)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		err := s.repo.MarkJobProcessing(ctx, job.ID)
		if err != nil {
			continue
		}

		err = s.processor.Process(ctx, job)
		if err != nil {
			s.repo.MarkJobFailed(ctx, job.ID, err.Error())
		} else {
			s.repo.MarkJobCompleted(ctx, job.ID)
		}
	}

	return nil
}
