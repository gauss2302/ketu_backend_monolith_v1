package main

import (
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/dto"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/repository/postgres"

	"ketu_backend_monolith_v1/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app := setupRouter(handlers, middleware, db)
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
	user  *http.UserHandler
	auth  *http.AuthHandler
	place *http.PlaceHandler
}

type middlewares struct {
	auth *middleware.AuthMiddleware
}

func setupDependencies(cfg *configs.Config, db *sqlx.DB) (*handlers, *middlewares) {
	log.Printf("Setting up dependencies...")

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	placeRepo := postgres.NewPlaceRepository(db)
	log.Printf("Repositories created - UserRepo: %v, PlaceRepo: %v", userRepo != nil, placeRepo != nil)

	// Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, &cfg.JWT)
	placeService := service.NewPlaceService(placeRepo)
	log.Printf("Services created - UserService: %v, AuthService: %v, PlaceService: %v",
		userService != nil, authService != nil, placeService != nil)

	// Handlers
	handlers := &handlers{
		user:  http.NewUserHandler(userService),
		auth:  http.NewAuthHandler(authService),
		place: http.NewPlaceHandler(placeService),
	}

	// Middleware
	middlewares := &middlewares{
		auth: middleware.NewAuthMiddleware(cfg.JWT),
	}

	return handlers, middlewares
}

func setupRouter(h *handlers, m *middlewares, db *sqlx.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Configure based on your needs
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Request logging middleware
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Incoming request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/register", middleware.ValidateBody(&dto.RegisterRequestDTO{}), h.auth.Register)
	auth.Post("/login", middleware.ValidateBody(&dto.LoginRequestDTO{}), h.auth.Login)

	// User routes
	users := v1.Group("/users")
	users.Post("/", middleware.ValidateBody(&dto.UserCreateDTO{}), h.user.Create)

	// Protected user routes
	users.Use(m.auth.AuthRequired())
	users.Get("/", h.user.GetAll)
	users.Get("/:id", h.user.GetByID)
	users.Put("/:id", middleware.ValidateBody(&dto.UserUpdateDTO{}), h.user.Update)
	users.Delete("/:id", h.user.Delete)

	// Place routes
	places := v1.Group("/places")

	// Public place routes
	places.Get("/", h.place.List)         // List all places
	places.Get("/search", h.place.Search) // Search places
	places.Get("/:id", h.place.GetByID)   // Get place by ID

	// Protected place routes
	places.Post("/", middleware.ValidateBody(&dto.PlaceCreateDTO{}), h.place.Create)
	places.Put("/:id", middleware.ValidateBody(&dto.PlaceUpdateDTO{}), h.place.Update)
	places.Delete("/:id", h.place.Delete) // Delete place

	// Print all registered routes for debugging
	for _, route := range app.GetRoutes() {
		log.Printf("Route registered: %s %s", route.Method, route.Path)
	}

	return app
}

func startServer(app *fiber.App, cfg *configs.Config) {
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	log.Printf("Server starting on %s", serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatal(err)
	}
}
