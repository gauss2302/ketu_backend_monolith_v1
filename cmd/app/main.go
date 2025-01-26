package main

import (
	"ketu_backend_monolith_v1/internal/app"
	"log"
)

func main() {
	application, err := app.New("configs")
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
