package consul

import (
	"context"
	"testing"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
)

func TestNewConsulClient(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "info",
		},
		Consul: config.ConsulConfig{
			Enabled:   false,
			Address:   "localhost:8500",
			ServiceID: "test-service-1",
		},
	}
	log := logger.NewLogger(cfg)

	client, err := NewConsulClient(cfg, log)
	if err != nil {
		t.Fatalf("NewConsulClient() error = %v", err)
	}
	if client.client != nil {
		t.Errorf("Expected nil client when Consul is disabled, got non-nil")
	}
}

func TestConsulClient_RegisterService(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "info",
		},
		Consul: config.ConsulConfig{
			Enabled:   false,
			Address:   "localhost:8500",
			ServiceID: "test-service-1",
		},
	}
	log := logger.NewLogger(cfg)
	client, _ := NewConsulClient(cfg, log)

	err := client.RegisterService(context.Background())
	if err != nil {
		t.Errorf("RegisterService() error = %v, want nil", err)
	}
}