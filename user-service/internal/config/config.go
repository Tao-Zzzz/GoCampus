package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Postgres struct {
		DSN string `envconfig:"POSTGRES_DSN" default:"postgres://user:pass@localhost:5432/userdb?sslmode=disable"`
	}
	Redis struct {
		Addr     string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
		Password string `envconfig:"REDIS_PASSWORD"`
	}
	JWT struct {
		Secret string `envconfig:"JWT_SECRET" default:"my-secret-key"`
	}
	NATS struct {
		URL string `envconfig:"NATS_URL" default:"nats://localhost:4222"`
	}
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}