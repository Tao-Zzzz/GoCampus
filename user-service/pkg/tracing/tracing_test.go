package tracing

import (
    "context"
    "testing"

    "github.com/Tao-Zzzz/GoCampus/user-service/config"
    "go.opentelemetry.io/otel"
)

func TestInitTracer(t *testing.T) {
    cfg := &config.Config{
        Service: config.ServiceConfig{
            Name:     "test-service",
            Port:     8080,
            LogLevel: "info",
        },
        Tracing: config.TracingConfig{
            JaegerEnabled: false,
            OTLPEnabled:   false,
        },
    }

    ctx := context.Background()
    shutdown, err := InitTracer(ctx, cfg)
    if err != nil {
        t.Fatalf("InitTracer() error = %v", err)
    }
    defer shutdown(ctx)

    tracer := otel.Tracer("test-tracer")
    _, span := tracer.Start(ctx, "test-span")
    span.End()

    if !span.SpanContext().IsValid() {
        t.Errorf("Expected valid span context, got invalid")
    }
}