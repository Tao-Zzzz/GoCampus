package consul

import (
	"context"
	"fmt"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	"github.com/hashicorp/consul/api"
)

// ConsulClient wraps the Consul API client.
type ConsulClient struct {
	client *api.Client
	logger *logger.Logger
	cfg    *config.Config
}

// NewConsulClient initializes a Consul client.
func NewConsulClient(cfg *config.Config, log *logger.Logger) (*ConsulClient, error) {
	if !cfg.Consul.Enabled {
		return &ConsulClient{logger: log, cfg: cfg}, nil
	}

	config := api.DefaultConfig()
	config.Address = cfg.Consul.Address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulClient{
		client: client,
		logger: log,
		cfg:    cfg,
	}, nil
}

// RegisterService registers the service with Consul.
func (c *ConsulClient) RegisterService(ctx context.Context) error {
	if !c.cfg.Consul.Enabled {
		c.logger.Info(ctx).Msg("Consul is disabled, skipping service registration")
		return nil
	}

	registration := &api.AgentServiceRegistration{
		ID:      c.cfg.Consul.ServiceID,
		Name:    c.cfg.Service.Name,
		Port:    c.cfg.Service.Port,
		Address: "localhost",
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://localhost:%d/health", c.cfg.Service.Port),
			Interval: "10s",
			Timeout:  "3s",
		},
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to register service with Consul")
		return fmt.Errorf("failed to register service: %w", err)
	}

	c.logger.Info(ctx).Msg("Service registered with Consul")
	return nil
}

// DeregisterService deregisters the service from Consul.
func (c *ConsulClient) DeregisterService(ctx context.Context) error {
	if !c.cfg.Consul.Enabled {
		return nil
	}

	err := c.client.Agent().ServiceDeregister(c.cfg.Consul.ServiceID)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msg("Failed to deregister service from Consul")
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	c.logger.Info(ctx).Msg("Service deregistered from Consul")
	return nil
}