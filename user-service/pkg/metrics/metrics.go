package metrics

import (
	"net/http"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds Prometheus metrics collectors.
type Metrics struct {
	requestDuration *prometheus.HistogramVec
}

// NewMetrics initializes Prometheus metrics.
func NewMetrics(cfg *config.Config) *Metrics {
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_service_request_duration_seconds",
			Help:    "Duration of user service requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
	prometheus.MustRegister(requestDuration)
	return &Metrics{requestDuration: requestDuration}
}

// RequestDuration returns the request duration histogram.
func (m *Metrics) RequestDuration() *prometheus.HistogramVec {
	return m.requestDuration
}

// StartMetricsServer starts an HTTP server for Prometheus metrics.
func StartMetricsServer(cfg *config.Config) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	addr := ":" + cfg.Metrics.Port
	return http.ListenAndServe(addr, mux)
}