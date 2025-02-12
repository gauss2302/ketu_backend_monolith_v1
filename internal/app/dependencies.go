package app

import (
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/handler/http"
	"ketu_backend_monolith_v1/internal/handler/middleware"
	"ketu_backend_monolith_v1/internal/repository/postgres"
	"ketu_backend_monolith_v1/internal/service"
	"log"

	"ketu_backend_monolith_v1/internal/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

type handlers struct {
	user       *http.UserHandler
	auth       *http.AuthHandler
	ownerAuth  *http.OwnerAuthHandler
	restaurant *http.RestaurantHandler
}

type middlewares struct {
	auth *middleware.AuthMiddleware
}

func setupDependencies(cfg *configs.Config, repos *postgres.Repositories, redisClient *redis.Client) (*handlers, *middlewares) {
	log.Printf("Setting up dependencies...")

	// Services
	userService := service.NewUserService(repos.User)
	authService := service.NewAuthService(repos.User, redisClient, &cfg.JWT)
	ownerAuthService := service.NewOwnerAuthService(repos.Owner, redisClient, &cfg.JWT)
	restaurantService := service.NewRestaurantService(repos.Restaurant)

	// Handlers
	handlers := &handlers{
		user:       http.NewUserHandler(userService),
		auth:       http.NewAuthHandler(authService),
		ownerAuth:  http.NewOwnerAuthHandler(ownerAuthService),
		restaurant: http.NewRestaurantHandler(restaurantService),
	}

	// Middleware
	middlewares := &middlewares{
		auth: middleware.NewAuthMiddleware(cfg.JWT),
	}

	return handlers, middlewares
}

type appDependencies struct {
	handlers    *handlers
	middlewares *middlewares
}

func initDependencies(cfg *configs.Config, repos *postgres.Repositories, redisClient *redis.Client) *appDependencies {
	h, m := setupDependencies(cfg, repos, redisClient)
	return &appDependencies{
		handlers:    h,
		middlewares: m,
	}
}

func initServer(deps *appDependencies) *fiber.App {
	return setupRouter(deps.handlers, deps.middlewares)
}
