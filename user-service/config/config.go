package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Service ServiceConfig
	Database DatabaseConfig
	JWT JWTConfig
	Consul ConsulConfig
	Etcd EtcdConfig
}

type ServiceConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	DurationHours int `mapstructure:"duration_hours"`
}

type ConsulConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Address   string `mapstructure:"address"`
	ServiceID string `mapstructure:"service_id"`
}


// EtcdConfig holds etcd settings.
type EtcdConfig struct {
	Enabled   bool
	Endpoints []string
}

// LoadConfig initializes and returns the application configuration.
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("service.name", "user-service")
	v.SetDefault("service.port", 8080)
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.dbname", "users")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("jwt.secret", "secret-key")
	v.SetDefault("jwt.duration_hours", 24)
	v.SetDefault("consul.enabled", false)
	v.SetDefault("consul.address", "localhost:8500")
	v.SetDefault("consul.service_id", "user-service-1")
	v.SetDefault("etcd.enabled", false)
	v.SetDefault("etcd.endpoints", []string{"localhost:2379"})

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Enable hot-reloading
	v.WatchConfig()

	// Load configuration
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// GetDSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}