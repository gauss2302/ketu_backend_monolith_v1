package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DB     PostgresConfig `mapstructure:"postgres"`
	Server ServerConfig   `mapstructure:"server"`
	JWT    JWTConfig      `mapstructure:"jwt"`
}

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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8090")
	viper.SetDefault("server.host", "localhost")

	// Database defaults
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", "5430")
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.dbname", "myapp")
	viper.SetDefault("postgres.sslmode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.accessTTL", "15m")
	viper.SetDefault("jwt.refreshTTL", "720h")
}
