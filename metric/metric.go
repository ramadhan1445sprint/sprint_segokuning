package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ServerMetric struct {
	ReqDurationHist *prometheus.HistogramVec
}

func NewServerMetric() *ServerMetric {
	reqDurationHist := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "segokuning_http_request_duration_seconds",
		Help: "http request duration in seconds",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"method", "path", "status"})

	return &ServerMetric{
		ReqDurationHist: reqDurationHist,
	}
}