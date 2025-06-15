package tracing

import (
    "context"

    "github.com/Tao-Zzzz/GoCampus/user-service/config"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

// InitTracer initializes OpenTelemetry tracing and returns a shutdown function.
func InitTracer(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
    var tp *sdktrace.TracerProvider

    if cfg.Tracing.JaegerEnabled {
        exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Tracing.JaegerEndpoint)))
        if err != nil {
            return nil, err
        }
        res, err := resource.New(ctx, resource.WithAttributes(
            semconv.ServiceNameKey.String(cfg.Service.Name),
        ))
        if err != nil {
            return nil, err
        }
        tp = sdktrace.NewTracerProvider(
            sdktrace.WithBatcher(exporter),
            sdktrace.WithResource(res),
        )
    } else if cfg.Tracing.OTLPEnabled {
        exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(cfg.Tracing.OTLPEndpoint), otlptracegrpc.WithInsecure())
        if err != nil {
            return nil, err
        }
        res, err := resource.New(ctx, resource.WithAttributes(
            semconv.ServiceNameKey.String(cfg.Service.Name),
        ))
        if err != nil {
            return nil, err
        }
        tp = sdktrace.NewTracerProvider(
            sdktrace.WithBatcher(exporter),
            sdktrace.WithResource(res),
        )
    } else {
        // No-op tracer provider if tracing is disabled
        tp = sdktrace.NewTracerProvider()
    }

    otel.SetTracerProvider(tp)

    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))

    return func(ctx context.Context) error {
        return tp.Shutdown(ctx)
    }, nil
}