package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	GRPCRequests *prometheus.CounterVec
	GRPCDuration *prometheus.HistogramVec
}

func New(serviceName string) *Metrics {

	m := &Metrics{

		GRPCRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "Total gRPC requests",
				ConstLabels: prometheus.Labels{
					"service": serviceName,
				},
			},
			[]string{"method", "code"},
		),

		GRPCDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "grpc_request_duration_seconds",
				Help: "gRPC request duration seconds",
				ConstLabels: prometheus.Labels{
					"service": serviceName,
				},
			},
			[]string{"method"},
		),
	}

	prometheus.MustRegister(
		m.GRPCRequests,
		m.GRPCDuration,
	)

	return m
}
