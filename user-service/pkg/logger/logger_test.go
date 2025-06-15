package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

func TestNewLogger(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "debug",
		},
	}

	logger := NewLogger(cfg)
	
	if logger.logger.GetLevel() != zerolog.DebugLevel {
		t.Errorf("Expected log level to be Debug, got %v", logger.logger.GetLevel())
	}

	// Test logging with context
	ctx := context.Background()
	var buf bytes.Buffer
	logger.logger = zerolog.New(&buf).With().Timestamp().Str("service", cfg.Service.Name).Logger()

	logger.Info(ctx).Msg("test message")
	var logOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logOutput); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logOutput["level"] != "info" {
		t.Errorf("Expected log level info, got %v", logOutput["level"])
	}
	if logOutput["service"] != "test-service" {
		t.Errorf("Expected service name test-service, got %v", logOutput["service"])
	}
	if logOutput["message"] != "test message" {
		t.Errorf("Expected message test message, got %v", logOutput["message"])
	}
}

func TestLogger_WithContext(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "info",
		},
	}
	logger := NewLogger(cfg)

	// Create a mock span context
	traceID := trace.TraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	spanID := trace.SpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8})
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)

	var buf bytes.Buffer
	logger.logger = zerolog.New(&buf).With().Timestamp().Str("service", cfg.Service.Name).Logger()

	logger.Info(ctx).Msg("test message")
	var logOutput map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logOutput); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logOutput["trace_id"] != traceID.String() {
		t.Errorf("Expected trace_id %v, got %v", traceID.String(), logOutput["trace_id"])
	}
	if logOutput["span_id"] != spanID.String() {
		t.Errorf("Expected span_id %v, got %v", spanID.String(), logOutput["span_id"])
	}
}