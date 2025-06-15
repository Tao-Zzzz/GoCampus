package etcd

import (
	"context"
	"fmt"
	"time" 

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdClient wraps the etcd client.
type EtcdClient struct {
	client *clientv3.Client
	logger *logger.Logger
	cfg    *config.Config
}

// NewEtcdClient initializes an etcd client.
func NewEtcdClient(cfg *config.Config, log *logger.Logger) (*EtcdClient, error) {
	if !cfg.Etcd.Enabled {
		return &EtcdClient{logger: log, cfg: cfg}, nil
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Etcd.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &EtcdClient{
		client: client,
		logger: log,
		cfg:    cfg,
	}, nil
}

// PutConfig stores a configuration key-value pair in etcd.
func (c *EtcdClient) PutConfig(ctx context.Context, key, value string) error {
	if !c.cfg.Etcd.Enabled {
		c.logger.Info(ctx).Msg("etcd is disabled, skipping config put")
		return nil
	}

	_, err := c.client.Put(ctx, key, value)
	if err != nil {
		c.logger.Error(ctx).Err(err).Msgf("Failed to put config key %s", key)
		return fmt.Errorf("failed to put config: %w", err)
	}

	c.logger.Info(ctx).Msgf("Stored config key %s in etcd", key)
	return nil
}

// Close closes the etcd client.
func (c *EtcdClient) Close() error {
	if !c.cfg.Etcd.Enabled || c.client == nil {
		return nil
	}
	return c.client.Close()
}