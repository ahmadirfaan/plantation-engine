package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmadirfaan/plantation-engine/generated"
	"github.com/ahmadirfaan/plantation-engine/handler"
	"github.com/ahmadirfaan/plantation-engine/repository"
	"github.com/ahmadirfaan/plantation-engine/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const (
	shutdownTimeout = 10 * time.Second
	serverAddress   = ":1323"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := setupDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to setup database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	e := NewServerApp(db)
	e.Logger.SetOutput(os.Stdout)

	srv := &http.Server{
		Addr:    serverAddress,
		Handler: e,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server starting", "addr", serverAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
}

func setupDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}

func NewServerApp(db *sql.DB) *echo.Echo {
	estateRepository := repository.NewEstateRepository(db)
	blockRepository := repository.NewBlockRepository(db)
	treeRepository := repository.NewTreeRepository(db)
	statsEstateRepository := repository.NewStatsEstateRepository(db)

	statsEstateUseCase := usecase.NewStatsEstateUseCase(statsEstateRepository, estateRepository)
	estateUseCase := usecase.NewEstateUseCase(estateRepository)
	treeUseCase := usecase.NewTreeUseCase(estateRepository, treeRepository, blockRepository, statsEstateUseCase)

	handlerServer := handler.NewServer(handler.NewServerOptions{
		EstateUseCase:      estateUseCase,
		TreeUseCase:        treeUseCase,
		StatsEstateUseCase: statsEstateUseCase,
	})

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(RequestIDMiddleware())
	e.Use(middleware.Logger())
	e.Use(handler.MetricsMiddleware)

	generated.RegisterHandlers(e, handlerServer)
	e.GET("/metrics", handler.MetricsHandler())

	return e
}

// RequestIDMiddleware assigns a request id to every request and echoes it back.
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Header.Get("X-Request-Id")
			if requestId == "" {
				requestId = uuid.New().String()
			}
			c.Response().Header().Set("X-Request-Id", requestId)
			c.Set("request_id", requestId)
			return next(c)
		}
	}
}
