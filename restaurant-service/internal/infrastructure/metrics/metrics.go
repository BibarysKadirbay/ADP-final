package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	GRPCRequests *prometheus.CounterVec
	GRPCDuration *prometheus.HistogramVec
	DBDuration   *prometheus.HistogramVec
	CacheEvents  *prometheus.CounterVec
}

func New() *Metrics {
	m := &Metrics{
		GRPCRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "restaurant_grpc_requests_total",
			Help: "Total gRPC requests.",
		}, []string{"method", "code"}),
		GRPCDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "restaurant_grpc_request_duration_seconds",
			Help:    "gRPC request latency.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method"}),
		DBDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "restaurant_db_query_duration_seconds",
			Help:    "Database query latency.",
			Buckets: prometheus.DefBuckets,
		}, []string{"operation"}),
		CacheEvents: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "restaurant_cache_events_total",
			Help: "Cache hit and miss counts.",
		}, []string{"result"}),
	}
	prometheus.MustRegister(m.GRPCRequests, m.GRPCDuration, m.DBDuration, m.CacheEvents)
	return m
}

func (m *Metrics) ObserveDB(operation string, started time.Time) {
	if m != nil {
		m.DBDuration.WithLabelValues(operation).Observe(time.Since(started).Seconds())
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}
