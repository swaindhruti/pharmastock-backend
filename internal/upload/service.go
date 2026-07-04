package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrInvalidFileType = errors.New("only .csv and .pdf files are allowed")

type JobService interface {
	CreateJob(ctx context.Context, stockistID int64, filePath string) (int64, error)
}

type Service interface {
	ProcessUpload(ctx context.Context, stockistID int64, file multipart.File, header *multipart.FileHeader) (int64, error)
}

type service struct {
	jobSvc    JobService
	uploadDir string
}

func NewService(jobSvc JobService, uploadDir string) Service {
	return &service{jobSvc: jobSvc, uploadDir: uploadDir}
}

func (s *service) ProcessUpload(ctx context.Context, stockistID int64, file multipart.File, header *multipart.FileHeader) (int64, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".csv" && ext != ".pdf" {
		return 0, ErrInvalidFileType
	}

	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create upload directory: %w", err)
	}

	fileName := fmt.Sprintf("%d_%d%s", stockistID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.uploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return 0, fmt.Errorf("failed to save file: %w", err)
	}

	return s.jobSvc.CreateJob(ctx, stockistID, filePath)
}
