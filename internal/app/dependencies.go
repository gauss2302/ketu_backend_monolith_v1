package app

import (
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
	"ketu_backend_monolith_v1/internal/repository/postgres"
	"ketu_backend_monolith_v1/internal/service"
	"log"

	"github.com/jmoiron/sqlx"
)

type handlers struct {
	user       *http.UserHandler
	auth       *http.AuthHandler
	restaurant *http.RestaurantHandler
}

type middlewares struct {
	auth *middleware.AuthMiddleware
}

func setupDependencies(cfg *configs.Config, db *sqlx.DB) (*handlers, *middlewares) {
	log.Printf("Setting up dependencies...")

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	restaurantRepo := postgres.NewRestaurantRepository(db)

	// Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, &cfg.JWT)
	restaurantService := service.NewRestaurantService(restaurantRepo)

	// Handlers
	handlers := &handlers{
		user:       http.NewUserHandler(userService),
		auth:       http.NewAuthHandler(authService),
		restaurant: http.NewRestaurantHandler(restaurantService),
	}

	// Middleware
	middlewares := &middlewares{
		auth: middleware.NewAuthMiddleware(cfg.JWT),
	}

	return handlers, middlewares
} 