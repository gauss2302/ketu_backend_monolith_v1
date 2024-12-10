package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	config "ketu_backend_monolith_v1/configs"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
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
	defer func(db *sqlx.DB) {
		err := database.Close(db)
		if err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}(db)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize usecases
	userUseCase := service.NewUserUsecase(userRepo)

	// Initialize handlers
	userHandler := http.NewUserHandler(userUseCase)

	// Setup Fiber app
	app := fiber.New()

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT)

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//Public routes
	users := v1.Group("/users")
	users.Post("/", middleware.ValidateBody(&dto.CreateUserInput{}), userHandler.Create)

	// Protected routes
	users.Use(authMiddleware.AuthRequired())
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", middleware.ValidateBody(&dto.UpdateUserInput{}), userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(app.Listen(serverAddr))
}
