package observability

import (
	"context"
	"testing"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
)

func TestInitObservability(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "info",
		},
		Metrics: config.MetricsConfig{
			Port: "9091",
		},
		Tracing: config.TracingConfig{
			JaegerEnabled: false,
			OTLPEnabled:   false,
		},
	}

	ctx := context.Background()
	obs, err := InitObservability(ctx, cfg)
	if err != nil {
		t.Fatalf("InitObservability() error = %v", err)
	}

	if obs.Logger == nil {
		t.Errorf("Expected non-nil logger")
	}
	if obs.Metrics == nil {
		t.Errorf("Expected non-nil metrics")
	}
	if obs.TracerShutdown == nil {
		t.Errorf("Expected non-nil tracer shutdown function")
	}

	// Test shutdown
	err = obs.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}