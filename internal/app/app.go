package app

import (
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type App struct {
	config *configs.Config
	db     *sqlx.DB
	server *fiber.App
}

func New(configPath string) (*App, error) {
	// Initialize config
	config, err := configs.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Initialize database
	db, err := initDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %v", err)
	}

	// Initialize dependencies
	deps := initDependencies(config, db)

	// Initialize server
	server := initServer(deps)

	return &App{
		config: config,
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run() error {
	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		_ = a.Shutdown()
	}()

	addr := fmt.Sprintf("%s:%s", a.config.Server.Host, a.config.Server.Port)
	log.Printf("Server starting on %s", addr)
	return a.server.Listen(addr)
}

func (a *App) Shutdown() error {
	if err := a.server.Shutdown(); err != nil {
		return fmt.Errorf("error shutting down server: %v", err)
	}
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("error closing database: %v", err)
	}
	return nil
}

// GetFiberApp returns the Fiber app instance
func (a *App) GetFiberApp() *fiber.App {
	return a.server
} 