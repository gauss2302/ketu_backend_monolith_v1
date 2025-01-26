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

	// Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, &cfg.JWT)
	placeService := service.NewPlaceService(placeRepo)

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