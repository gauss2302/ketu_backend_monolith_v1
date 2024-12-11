package main

import (
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/handler/dto"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/repository/postgres"
	"ketu_backend_monolith_v1/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Initialize app
	cfg, db := initializeApp()
	defer db.Close()

	// Add debug logging
	log.Printf("App initialized with config: %+v", cfg)

	// Setup dependencies
	handlers, middleware := setupDependencies(cfg, db)
	log.Printf("Dependencies setup complete")

	// Setup Fiber app and routes
	app := setupRouter(handlers, middleware)
	log.Printf("Router setup complete")

	// Add a very basic route directly in main to test
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Print all routes before starting
	log.Printf("Registered routes:")
	for _, route := range app.GetRoutes() {
		log.Printf("%s %s", route.Method, route.Path)
	}

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
	// Add this logging
	log.Printf("Setting up dependencies...")

	userRepo := postgres.NewUserRepository(db)
	log.Printf("UserRepo created: %v", userRepo != nil)

	userService := service.NewUserUsecase(userRepo)
	log.Printf("UserService created: %v", userService != nil)

	authService := service.NewAuthService(userRepo, &cfg.JWT)
	log.Printf("AuthService created: %v", authService != nil)

	handlers := &handlers{
		user: http.NewUserHandler(userService),
		auth: http.NewAuthHandler(authService),
	}
	log.Printf("Handlers created - auth handler: %v", handlers.auth != nil)

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

	log.Printf("Setting up routes with handlers: %+v", h)
	log.Printf("Auth handler nil check: %v", h.auth != nil)

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

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Test route works!")
	})

	return app
}

func startServer(app *fiber.App, cfg *configs.Config) {
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}
