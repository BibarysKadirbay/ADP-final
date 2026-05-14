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
			Name: "delivery_grpc_requests_total",
			Help: "Total delivery gRPC requests.",
		}, []string{"method", "code"}),
		GRPCDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "delivery_grpc_request_duration_seconds",
			Help:    "Delivery gRPC request duration.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method"}),
		DBDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "delivery_db_query_duration_seconds",
			Help:    "Delivery database query duration.",
			Buckets: prometheus.DefBuckets,
		}, []string{"query"}),
		CacheEvents: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "delivery_cache_events_total",
			Help: "Delivery cache hit and miss counters.",
		}, []string{"result"}),
	}
	prometheus.MustRegister(m.GRPCRequests, m.GRPCDuration, m.DBDuration, m.CacheEvents)
	return m
}

func (m *Metrics) ObserveDB(query string, started time.Time) {
	m.DBDuration.WithLabelValues(query).Observe(time.Since(started).Seconds())
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}
