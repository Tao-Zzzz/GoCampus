package etcd

import (
	"context"
	"testing"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
)

func TestNewEtcdClient(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "info",
		},
		Etcd: config.EtcdConfig{
			Enabled:   false,
			Endpoints: []string{"localhost:2379"},
		},
	}
	log := logger.NewLogger(cfg)

	client, err := NewEtcdClient(cfg, log)
	if err != nil {
		t.Fatalf("NewEtcdClient() error = %v", err)
	}
	if client.client != nil {
		t.Errorf("Expected nil client when etcd is disabled, got non-nil")
	}
}

func TestEtcdClient_PutConfig(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "info",
		},
		Etcd: config.EtcdConfig{
			Enabled:   false,
			Endpoints: []string{"localhost:2379"},
		},
	}
	log := logger.NewLogger(cfg)
	client, _ := NewEtcdClient(cfg, log)

	err := client.PutConfig(context.Background(), "test/key", "test-value")
	if err != nil {
		t.Errorf("PutConfig() error = %v, want nil", err)
	}
}