package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Service  ServiceConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Consul   ConsulConfig
	Etcd     EtcdConfig
	Tracing  TracingConfig
	Metrics  MetricsConfig
}
// MetricsConfig holds metrics settings.
type MetricsConfig struct {
	Port string `mapstructure:"port"`
}

type ServiceConfig struct {
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
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
	Enabled   bool `mapstructure:"enabled"`
	Endpoints []string `mapstructure:"endpoints"`
}

// TracingConfig holds tracing settings.
type TracingConfig struct {
    JaegerEnabled  bool   `mapstructure:"jaeger_enabled"`
    JaegerEndpoint string `mapstructure:"jaeger_endpoint"`
    OTLPEnabled    bool   `mapstructure:"otlp_enabled"`
    OTLPEndpoint   string `mapstructure:"otlp_endpoint"`
}

// LoadConfig initializes and returns the application configuration.
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("service.name", "user-service")
	v.SetDefault("service.port", 8080)
	v.SetDefault("service.log_level", "info")
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

	v.SetDefault("tracing.jaeger_enabled", false)
	v.SetDefault("tracing.jaeger_endpoint", "http://localhost:14268/api/traces")
	v.SetDefault("tracing.otlp_enabled", false)
	v.SetDefault("tracing.otlp_endpoint", "localhost:4317")

	v.SetDefault("metrics.port", "9090")
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