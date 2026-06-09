package worker

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/swaindhruti/pharmastock-backend/internal/config"
	"github.com/swaindhruti/pharmastock-backend/internal/database"
	"github.com/swaindhruti/pharmastock-backend/internal/job"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	jobRepo := job.NewRepository(db.Pool)
	jobProcessor := job.NewProcessor()
	jobService := job.NewService(jobRepo, jobProcessor)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Worker started, processing jobs every 10 seconds")

	for {
		select {
		case <-ticker.C:
			err := jobService.ProcessPendingJobs(context.Background(), 5)
			if err != nil {
				log.Printf("Error processing jobs: %v", err)
			}
		case sig := <-sigChan:
			log.Printf("Received signal: %v, shutting down worker", sig)
			return
		}
	}
}
