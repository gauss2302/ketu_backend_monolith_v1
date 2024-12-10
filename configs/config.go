package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type (
	Config struct {
		Server ServerConfig   `mapstructure:"server"`
		DB     PostgresConfig `mapstructure:"postgres"`
		JWT    JWTConfig      `mapstructure:"jwt"`
		Logger LoggerConfig   `mapstructure:"logger"`
	}

	ServerConfig struct {
		Port         string        `mapstructure:"port"`
		Host         string        `mapstructure:"host"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	JWTConfig struct {
		AccessSecret  string        `mapstructure:"accessSecret"`
		RefreshSecret string        `mapstructure:"refreshSecret"`
		AccessTTL     time.Duration `mapstructure:"accessTTL"`
		RefreshTTL    time.Duration `mapstructure:"refreshTTL"`
	}

	LoggerConfig struct {
		Level string `mapstructure:"level"`
	}
)

func LoadConfig(path string) (*Config, error) {
	// Check environment
	env := os.Getenv("APP_ENV")

	if env == "production" {
		return loadProdConfig()
	}

	return loadDevConfig(path)
}

func loadDevConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return &config, nil
}

func loadProdConfig() (*Config, error) {
	setDefaults()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.readTimeout", "10s")
	viper.SetDefault("server.writeTimeout", "10s")

	// Database defaults
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.dbname", "myapp")
	viper.SetDefault("postgres.sslmode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.accessTTL", "15m")
	viper.SetDefault("jwt.refreshTTL", "720h") // 30 days

	// Logger defaults
	viper.SetDefault("logger.level", "info")
}
