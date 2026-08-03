package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMetricsMiddlewareAndHandler(t *testing.T) {
	e := echo.New()
	e.Use(MetricsMiddleware)
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/metrics", MetricsHandler())

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	body, _ := io.ReadAll(rec.Body)
	assert.Contains(t, string(body), "http_requests_total")
	assert.Contains(t, string(body), "http_request_duration_seconds")
}

func TestMetricsMiddlewareRecordsErrorAs500(t *testing.T) {
	e := echo.New()
	e.Use(MetricsMiddleware)
	e.GET("/boom", func(c echo.Context) error {
		return c.NoContent(http.StatusInternalServerError)
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
