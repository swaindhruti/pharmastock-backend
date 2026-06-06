package main

import (
	"log"

	"github.com/swaindhruti/pharmastock-backend/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Shutdown()

	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		log.Printf("Server started on port %s", application.Config.AppPort)
	}

}
