package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	config "ketu_backend_monolith_v1/configs"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/repository/postgres"
	"ketu_backend_monolith_v1/internal/service"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close(db)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize usecases
	userUseCase := service.NewUserUsecase(userRepo)

	// Initialize handlers
	userHandler := http.NewUserHandler(userUseCase)

	// Setup Fiber app
	app := fiber.New()

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes
	users := v1.Group("/users")
	users.Post("/", userHandler.Create)
	//users.Get("/:id", userHandler.GetByID)
	//users.Put("/:id", userHandler.Update)
	//users.Delete("/:id", userHandler.Delete)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(serverAddr))
}
