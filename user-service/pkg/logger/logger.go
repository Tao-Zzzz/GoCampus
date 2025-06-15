package logger

import (
	"context"
	"os"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Logger wraps a zerolog.Logger with trace correlation.
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new Logger instance.
func NewLogger(cfg *config.Config) *Logger {
    var level zerolog.Level
    switch cfg.Service.LogLevel {
    case "debug":
        level = zerolog.DebugLevel
    case "info":
        level = zerolog.InfoLevel
    case "warn":
        level = zerolog.WarnLevel
    case "error":
        level = zerolog.ErrorLevel
    default:
        level = zerolog.InfoLevel
    }
    zerolog.SetGlobalLevel(level)
    logger := zerolog.New(os.Stdout).
        Level(level). // 这里设置实例级别
        With().
        Timestamp().
        Str("service", cfg.Service.Name).
        Logger()
    return &Logger{logger: logger}
}

// WithContext adds trace ID and span ID from the context to the logger.
func (l *Logger) WithContext(ctx context.Context) *zerolog.Logger {
    log := l.logger
    spanCtx := trace.SpanContextFromContext(ctx)
    if spanCtx.HasTraceID() && spanCtx.HasSpanID() {
        log = log.With().
            Str("trace_id", spanCtx.TraceID().String()).
            Str("span_id", spanCtx.SpanID().String()).
            Logger()
    }
    return &log
}

// Debug logs a debug message.
func (l *Logger) Debug(ctx context.Context) *zerolog.Event {
	return l.WithContext(ctx).Debug()
}

// Info logs an info message.
func (l *Logger) Info(ctx context.Context) *zerolog.Event {
	return l.WithContext(ctx).Info()
}

// Warn logs a warning message.
func (l *Logger) Warn(ctx context.Context) *zerolog.Event {
	return l.WithContext(ctx).Warn()
}

// Error logs an error message.
func (l *Logger) Error(ctx context.Context) *zerolog.Event {
	return l.WithContext(ctx).Error()
}