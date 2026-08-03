package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware_GeneratesWhenAbsent(t *testing.T) {
	e := echo.New()
	e.Use(RequestIDMiddleware())
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-Request-Id"))
}

func TestRequestIDMiddleware_ReusesIncoming(t *testing.T) {
	e := echo.New()
	e.Use(RequestIDMiddleware())
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "trace-123")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, "trace-123", rec.Header().Get("X-Request-Id"))
}

func TestSetupDB_InvalidDSN(t *testing.T) {
	_, err := setupDB("not-a-valid-dsn")
	assert.Error(t, err)
}
