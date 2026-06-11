package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/swaindhruti/pharmastock-backend/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Shutdown()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
