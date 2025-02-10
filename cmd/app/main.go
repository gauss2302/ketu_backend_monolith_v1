package main

// @title           Ketu Backend API
// @version         1.0
// @description     Restaurant Management System API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8090
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization

import (
	"ketu_backend_monolith_v1/internal/app"
	"log"

	"github.com/gofiber/swagger"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Swagger documentation
	application.GetFiberApp().Get("/swagger/*", swagger.HandlerDefault)

	if err := application.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
