package app

import (
	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupRouter(h *handlers, m *middlewares) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler(),
	})

	// Middleware
	app.Use(cors.New())

	// Routes
	setupRoutes(app, h, m)

	return app
}

func setupRoutes(app *fiber.App, h *handlers, m *middlewares) {
	// Health check
	app.Get("/health", http.NewHealthHandler().Handle)

	// API routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateBody(&dto.RegisterRequestDTO{}), h.auth.Register)
	auth.Post("/login", middleware.ValidateBody(&dto.LoginRequestDTO{}), h.auth.Login)
	auth.Post("/refresh", h.auth.RefreshToken)

	// Protected routes
	protected := api.Group("")
	protected.Use(m.auth.Authenticate())
	setupProtectedRoutes(protected, h)
}

func setupProtectedRoutes(protected fiber.Router, h *handlers) {
	// User routes
	users := protected.Group("/users")
	users.Get("/", h.user.GetAll)
	users.Get("/:id", h.user.GetByID)
	users.Put("/:id", middleware.ValidateBody(&dto.UserUpdateDTO{}), h.user.Update)
	users.Delete("/:id", h.user.Delete)

	// Place routes
	places := protected.Group("/places")
	places.Post("/", middleware.ValidateBody(&dto.PlaceCreateDTO{}), h.place.Create)
	places.Get("/", h.place.List)
	places.Get("/:id", h.place.GetByID)
	places.Put("/:id", middleware.ValidateBody(&dto.PlaceUpdateDTO{}), h.place.Update)
	places.Delete("/:id", h.place.Delete)
}

func errorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		log.Printf("Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
} 