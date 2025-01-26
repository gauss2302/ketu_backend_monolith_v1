package configs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DB     PostgresConfig `mapstructure:"postgres"`
	Server ServerConfig   `mapstructure:"server"`
	JWT    JWTConfig      `mapstructure:"jwt"`
	// Migration MigrationConfig `mapstructure:"migration"`
}

// type MigrationConfig struct {
// 	Path string `mapstructure:"path"`
// }

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"accessSecret"`
	RefreshSecret string        `mapstructure:"refreshSecret"`
	AccessTTL     time.Duration `mapstructure:"accessTTL"`
	RefreshTTL    time.Duration `mapstructure:"refreshTTL"`
}

func LoadConfig(path string) (*Config, error) {
	env := getEnvironment()
	log.Printf("Loading configuration for environment: %s", env)
	
	// Load environment-specific .env file
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		// Try fallback to default .env
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: no .env file found: %v", err)
		}
	}

	// Configure Viper
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// Enable environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults based on environment
	setDefaults(env)

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
		log.Printf("No config file found, using environment variables and defaults")
	} else {
		log.Printf("Loaded config file successfully")
	}

	// Load environment variables
	loadFromEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	// Log configuration (excluding sensitive data)
	log.Printf("Server configuration - Host: %s, Port: %s", config.Server.Host, config.Server.Port)
	log.Printf("Database configuration - Host: %s, Port: %s, Database: %s", config.DB.Host, config.DB.Port, config.DB.DBName)

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getEnvironment() string {
	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "development" // Default to development
	}
	return env
}

func loadFromEnv() {
	// Server
	viper.BindEnv("server.port", "APP_SERVER_PORT")
	viper.BindEnv("server.host", "APP_SERVER_HOST")

	// Database
	viper.BindEnv("postgres.host", "APP_POSTGRES_HOST")
	viper.BindEnv("postgres.port", "APP_POSTGRES_PORT")
	viper.BindEnv("postgres.username", "APP_POSTGRES_USERNAME")
	viper.BindEnv("postgres.password", "APP_POSTGRES_PASSWORD")
	viper.BindEnv("postgres.dbname", "APP_POSTGRES_DBNAME")
	viper.BindEnv("postgres.sslmode", "APP_POSTGRES_SSLMODE")

	// JWT
	viper.BindEnv("jwt.accessSecret", "APP_JWT_ACCESSSECRET")
	viper.BindEnv("jwt.refreshSecret", "APP_JWT_REFRESHSECRET")
	viper.BindEnv("jwt.accessTTL", "APP_JWT_ACCESSTTL")
	viper.BindEnv("jwt.refreshTTL", "APP_JWT_REFRESHTTL")
}

func validateConfig(config *Config) error {
	// Validate database configuration
	if config.DB.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.DB.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if config.DB.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if config.DB.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if config.DB.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate JWT configuration
	if config.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}
	if config.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}

	return nil
}

func setDefaults(env string) {
	// Common defaults
	viper.SetDefault("server.port", "8090")
	
	// Environment-specific defaults
	switch env {
	case "production":
		viper.SetDefault("server.host", "0.0.0.0")
		viper.SetDefault("postgres.host", "postgres")
		viper.SetDefault("logger.level", "info")
	default: // development
		viper.SetDefault("server.host", "localhost")
		viper.SetDefault("postgres.host", "localhost")
		viper.SetDefault("logger.level", "debug")
	}

	// Database defaults
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.sslmode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.accessTTL", "15m")
	viper.SetDefault("jwt.refreshTTL", "720h")
}
