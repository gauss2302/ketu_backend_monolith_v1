package app

import (
	"log"

	"ketu_backend_monolith_v1/internal/dto"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupRouter(h *handlers, m *middlewares) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler(),
	})

	// Configure CORS
// Configure CORS with specific origin instead of wildcard
app.Use(cors.New(cors.Config{
	AllowOrigins: "http://localhost:3000",
	AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	AllowCredentials: true,
}))

	app.Use(middleware.RequestLogger())

	setupRoutes(app, h, m)
	return app
}

func setupRoutes(app *fiber.App, h *handlers, m *middlewares) {
	app.Get("/health", http.NewHealthHandler().Handle)

	api := app.Group("/api/v1")

	// Public routes (no auth required)
	// Owner auth routes
	ownerAuth := api.Group("/owner/auth")
	ownerAuth.Post("/register", middleware.ValidateBody(&dto.OwnerRegisterRequestDTO{}), h.ownerAuth.Register)
	ownerAuth.Post("/login", middleware.ValidateBody(&dto.OwnerLoginRequestDTO{}), h.ownerAuth.Login)

	// User auth routes
	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateBody(&dto.RegisterRequestDTO{}), h.auth.Register)
	auth.Post("/login", middleware.ValidateBody(&dto.LoginRequestDTO{}), h.auth.Login)

	// Protected routes group
	protected := api.Group("")
	protected.Use(m.auth.Authenticate())

	// Protected auth routes (under protected group)
	protected.Group("/auth").Post("/refresh", h.auth.RefreshToken)
	protected.Group("/auth").Post("/logout", h.auth.Logout)
	protected.Group("/owner/auth").Post("/refresh", h.ownerAuth.RefreshToken)
	protected.Group("/owner/auth").Post("/logout", h.ownerAuth.Logout)

	// Other protected routes
	setupProtectedRoutes(protected, h)
}

func setupProtectedRoutes(protected fiber.Router, h *handlers) {
	users := protected.Group("/users")
	users.Get("/", h.user.GetAll)
	users.Get("/:id", h.user.GetByID)
	users.Put("/:id", middleware.ValidateBody(&dto.UserUpdateDTO{}), h.user.Update)
	users.Delete("/:id", h.user.Delete)

	// Restaurant routes (require authentication)
	restaurants := protected.Group("/restaurants")
	restaurants.Post("/", middleware.ValidateBody(&dto.CreateRestaurantDTO{}), h.restaurant.Create)
	restaurants.Get("/", h.restaurant.List)
	restaurants.Get("/my", h.restaurant.ListByOwner)
	restaurants.Get("/:id", h.restaurant.GetByID)
	restaurants.Put("/:id", middleware.ValidateBody(&dto.UpdateRestaurantDTO{}), h.restaurant.Update)
	restaurants.Delete("/:id", h.restaurant.Delete)
}

func errorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		log.Printf("Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}