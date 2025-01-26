package app

import (
	configs "ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func InitializeApp() (*configs.Config, *sqlx.DB) {
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

func initDB(cfg *configs.Config) (*sqlx.DB, error) {
	db, err := database.NewPostgresDB(&cfg.DB)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type appDependencies struct {
	handlers    *handlers
	middlewares *middlewares
}

func initDependencies(cfg *configs.Config, db *sqlx.DB) *appDependencies {
	h, m := setupDependencies(cfg, db)
	return &appDependencies{
		handlers:    h,
		middlewares: m,
	}
}

func initServer(deps *appDependencies) *fiber.App {
	return setupRouter(deps.handlers, deps.middlewares)
} 