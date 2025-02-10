package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds the application configurations
type Config struct {
	DB     PostgresConfig `mapstructure:"postgres"`
	Server ServerConfig   `mapstructure:"server"` // Add ServerConfig
	JWT    JWTConfig      `mapstructure:"jwt"`
	Redis  RedisConfig   `mapstructure:"redis"`
}

// PostgresConfig holds PostgreSQL configurations
type PostgresConfig struct {
	URL string `mapstructure:"url"`
}

// ServerConfig holds server configurations
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// JWTConfig holds JWT configurations
type JWTConfig struct {
	AccessSecret  string        `mapstructure:"accesssecret"`
	RefreshSecret string        `mapstructure:"refreshsecret"`
	AccessTTL     time.Duration `mapstructure:"accessttl"`
	RefreshTTL    time.Duration `mapstructure:"refreshttl"`
}

// RedisConfig holds Redis configurations
type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LoadConfig initializes and returns the application configuration
func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: no .env file found: %v", err)
	}

	// Viper configuration
	viper.AutomaticEnv() // Automatically read environment variables

	// Replace dots with underscores in environment variables
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Explicitly bind environment variables to configuration keys
	bindEnvs()

	// Set default values
	viper.SetDefault("server.host", "0.0.0.0")  // Allow external connections
	viper.SetDefault("server.port", "8090")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	// Update logging to use the new URL field
	log.Printf("Database configuration loaded with URL schema: %s", "postgres://****:****@" + config.DB.URL)

	return &config, nil
}

// bindEnvs binds each configuration key to its corresponding environment variable
func bindEnvs() {
	envBindings := map[string]string{
		"postgres.url": "DATABASE_URL",
		"server.port":       "APP_SERVER_PORT", // Bind server port
		"server.host":       "APP_SERVER_HOST", // Bind server host
		"jwt.accessSecret":  "APP_JWT_ACCESSSECRET",
		"jwt.refreshSecret": "APP_JWT_REFRESHSECRET",
		"jwt.accessTTL":     "APP_JWT_ACCESSTTL",
		"jwt.refreshTTL":    "APP_JWT_REFRESHTTL",
		"redis.url":      "REDIS_URL",
		"redis.password": "REDIS_PASSWORD",
		"redis.db":       "REDIS_DB",
	}

	for configKey, envVar := range envBindings {
		if err := viper.BindEnv(configKey, envVar); err != nil {
			log.Fatalf("Error binding %s: %v", configKey, err)
		}
	}
}