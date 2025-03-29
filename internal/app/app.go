package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	type result struct {
		db    *database.DB
		redis *redis.Client
		err   error
		which string
	}

	initChan := make(chan result, 2)

	go func() {
		db, err := database.NewPostgresDB(cfg.DB.URL)
		initChan <- result{db: db, err: err, which: "database"}
	}()

	go func() {
		redisClient, err := redis.NewRedisClient(&cfg.Redis)
		initChan <- result{redis: redisClient, err: err, which: "redis"}
	}()

	var db *database.DB
	var redisClient *redis.Client

	for i := 0; i < 2; i++ {
		res := <-initChan

		if res.err != nil {
			return nil, fmt.Errorf("failed to init %s: %v", res.which, res.err)
		}

		if res.which == "database" {
			db = res.db
		} else {
			redisClient = res.redis
		}
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
	var wg sync.WaitGroup

	errChan := make(chan error, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := a.server.Shutdown(); err != nil {
			errChan <- fmt.Errorf("Error shutting down Http server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.db.Close(); err != nil {
			errChan <- fmt.Errorf("Error closing db connection: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.redis.Close(); err != nil {
			errChan <- fmt.Errorf("Error closing redis connection: %v", err)
		}
	}()

	wg.Wait()
	close(errChan)

	var errors []error

	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		var errMsg string
		for _, err := range errors {
			errMsg += err.Error() + "; "
		}
		return fmt.Errorf("Shutdown errors: %s", errMsg)
	}

	return nil
}

func (a *App) GetFiberApp() *fiber.App {
	return a.server
}
