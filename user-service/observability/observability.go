package observability

import (
	"context"
	"fmt"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/metrics"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/tracing"
)

// Observability holds observability components.
type Observability struct {
	Logger        *logger.Logger
	Metrics       *metrics.Metrics
	TracerShutdown func(context.Context) error
}

// InitObservability initializes logging, tracing, and metrics.
func InitObservability(ctx context.Context, cfg *config.Config) (*Observability, error) {
	// Initialize logger
	log := logger.NewLogger(cfg)

	// Initialize tracing
	tracerShutdown, err := tracing.InitTracer(ctx, cfg)
	if err != nil {
		log.Error(ctx).Err(err).Msg("Failed to initialize tracer")
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize metrics
	met := metrics.NewMetrics(cfg)

	// Start metrics server in a goroutine
	go func() {
		if err := metrics.StartMetricsServer(cfg); err != nil {
			log.Error(ctx).Err(err).Msg("Failed to start metrics server")
		}
	}()

	log.Info(ctx).Msg("Observability initialized")
	return &Observability{
		Logger:        log,
		Metrics:       met,
		TracerShutdown: tracerShutdown,
	}, nil
}

// Shutdown cleans up observability resources.
func (o *Observability) Shutdown(ctx context.Context) error {
	if err := o.TracerShutdown(ctx); err != nil {
		o.Logger.Error(ctx).Err(err).Msg("Failed to shutdown tracer")
		return fmt.Errorf("failed to shutdown tracer: %w", err)
	}
	o.Logger.Info(ctx).Msg("Observability shutdown complete")
	return nil
}