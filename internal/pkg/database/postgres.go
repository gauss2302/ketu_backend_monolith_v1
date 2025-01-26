package database

import (
	"fmt"
	configs "ketu_backend_monolith_v1/internal/config"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 25
	connMaxLifetime = 15 // minutes
	maxIdleConns    = 25
	connMaxIdleTime = 10 // minutes
	maxRetries      = 10
	retryDelay      = 5 * time.Second
)

func NewPostgresDB(cfg *configs.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Add retry logic with exponential backoff
	var db *sqlx.DB
	var err error
	
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		
		retryTime := retryDelay * time.Duration(i+1)
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...", 
			i+1, maxRetries, err, retryTime)
		time.Sleep(retryTime)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Minute * connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Minute * connMaxIdleTime)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying database connection: %v", err)
	}

	log.Println("Successfully connected to database")

	// Run migrations
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return db, nil
}
