package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	config "ketu_backend_monolith_v1/configs"
	"log"
	"time"
)

const (
	maxOpenConns    = 25
	connMaxLifetime = 15 // minutes
	maxIdleConns    = 25
	connMaxIdleTime = 10 // minutes
)

func NewPostgresDB(cfg *config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Add retry logic
	var db *sqlx.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres after %d attempts: %v", maxRetries, err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Minute * connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(time.Second * connMaxIdleTime)

	return db, nil
}
