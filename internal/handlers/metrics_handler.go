package handlers

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// MetricsHandler espone le metriche per Prometheus.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
