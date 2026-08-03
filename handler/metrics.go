package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

// MetricsMiddleware records request counts and latency per route.
func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		status := c.Response().Status
		if err != nil {
			status = http.StatusInternalServerError
		}
		httpRequestsTotal.WithLabelValues(c.Request().Method, c.Path(), strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request().Method, c.Path()).Observe(time.Since(start).Seconds())
		return err
	}
}

// MetricsHandler exposes Prometheus metrics at /metrics.
func MetricsHandler() echo.HandlerFunc {
	return echo.WrapHandler(promhttp.Handler())
}
