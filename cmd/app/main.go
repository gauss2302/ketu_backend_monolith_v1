package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/repository/postgres"
	"ketu_backend_monolith_v1/internal/service"
	"log"
)

func main() {
	// Initialize app
	cfg, db := initializeApp()
	defer db.Close()

	// Setup dependencies
	handlers, middleware := setupDependencies(cfg, db)

	// Setup Fiber app and routes
	app := setupRouter(handlers, middleware)

	// Start server
	startServer(app, cfg)
}

func initializeApp() (*configs.Config, *sqlx.DB) {
	// Load configuration
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return cfg, db
}

type handlers struct {
	user *http.UserHandler
	auth *http.AuthHandler
}

type middlewares struct {
	auth *middleware.AuthMiddleware
}

func setupDependencies(cfg *configs.Config, db *sqlx.DB) (*handlers, *middlewares) {
	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserUsecase(userRepo)
	authService := service.NewAuthService(userRepo, &cfg.JWT)
	log.Printf("Auth service initialized: %v", authService != nil)

	// Print for debugging
	log.Printf("AuthService initialized: %v", authService != nil)

	// Initialize handlers
	handlers := &handlers{
		user: http.NewUserHandler(userService),
		auth: http.NewAuthHandler(authService),
	}

	// Print for debugging
	log.Printf("AuthHandler initialized: %v", handlers.auth != nil)

	// Initialize middleware
	middlewares := &middlewares{
		auth: middleware.NewAuthMiddleware(cfg.JWT),
	}

	return handlers, middlewares
}

func setupRouter(h *handlers, m *middlewares) *fiber.App {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// Add simple middleware to log all requests
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Incoming request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// IMPORTANT: Register auth routes BEFORE user routes
	// Auth routes (public)
	auth := v1.Group("/auth")
	auth.Post("/register", middleware.ValidateBody(&dto.RegisterRequest{}), h.auth.Register)
	auth.Post("/login", middleware.ValidateBody(&dto.LoginRequest{}), h.auth.Login)

	// User routes
	users := v1.Group("/users")

	// Public user routes
	users.Post("/", middleware.ValidateBody(&dto.CreateUserInput{}), h.user.Create)

	// Protected user routes
	users.Use(m.auth.AuthRequired())
	users.Get("/", h.user.GetAll)
	users.Get("/:id", h.user.GetByID)
	users.Put("/:id", middleware.ValidateBody(&dto.UpdateUserInput{}), h.user.Update)
	users.Delete("/:id", h.user.Delete)

	// Print all registered routes for debugging
	for _, route := range app.GetRoutes() {
		log.Printf("Route registered: %s %s", route.Method, route.Path)
	}

	return app
}

func startServer(app *fiber.App, cfg *configs.Config) {
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}
