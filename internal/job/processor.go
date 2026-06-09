package job

import (
	"context"
	"time"
)

type Processor interface {
	Process(ctx context.Context, job *Job) error
}

type processor struct{}

func NewProcessor() Processor {
	return &processor{}
}

func (p *processor) Process(ctx context.Context, job *Job) error {
	// Simulate job processing time
	time.Sleep(5 * time.Second)

	// Here you would add the actual logic to process the job, e.g., read the file, update the database, etc.

	return nil
}
