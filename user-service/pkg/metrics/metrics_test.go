package metrics

import (
    "testing"

    "github.com/Tao-Zzzz/GoCampus/user-service/config"
    "github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewMetrics(t *testing.T) {
    cfg := &config.Config{
        Service: config.ServiceConfig{
            Name:     "test-service",
            Port:     8080,
            LogLevel: "info",
        },
        Metrics: config.MetricsConfig{
            Port: "9091",
        },
    }

    metrics := NewMetrics(cfg)
    metrics.RequestDuration().WithLabelValues("test_method", "success").Observe(0.1)

    // 检查整个 HistogramVec 的采样数
    count := testutil.CollectAndCount(metrics.RequestDuration())
    if count == 0 {
        t.Errorf("Expected request duration metric to be recorded")
    }
}