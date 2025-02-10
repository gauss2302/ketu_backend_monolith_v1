package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ketu_backend_monolith_v1/internal/config"
	"ketu_backend_monolith_v1/internal/pkg/database"
	"ketu_backend_monolith_v1/internal/pkg/redis"
	"ketu_backend_monolith_v1/internal/repository/postgres"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	config *config.Config
	db     *database.DB
	redis  *redis.Client
	repos  *postgres.Repositories
	server *fiber.App
}

func New() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.DB.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %v", err)
	}

	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("failed to init Redis: %v", err)
	}

	repos := postgres.NewRepositories(db)
	server := initServer(initDependencies(cfg, repos, redisClient))

	return &App{
		config: cfg,
		db:     db,
		redis:  redisClient,
		repos:  repos,
		server: server,
	}, nil
}

func (a *App) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := a.Shutdown(); err != nil {
			log.Fatalf("Failed to shutdown server: %v", err)
		}
	}()

	// Assuming the server configuration is directly under config
	addr := fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port)
	log.Printf("Server starting on %s", addr)
	return a.server.Listen(addr)
}

func (a *App) Shutdown() error {
	if err := a.server.Shutdown(); err != nil {
		return fmt.Errorf("error shutting down HTTP server: %v", err)
	}
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %v", err)
	}
	if err := a.redis.Close(); err != nil {
		return fmt.Errorf("error closing redis connection: %v", err)
	}
	return nil
}

func (a *App) GetFiberApp() *fiber.App {
	return a.server
}
