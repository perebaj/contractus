// Package main grather all the pieces and start the service.
package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"log/slog"

	"github.com/birdie-ai/contractus/api"
	"github.com/birdie-ai/contractus/postgres"
	"github.com/go-chi/chi/v5"
)

// Config have the core configuration for the service.
type Config struct {
	PORT     string
	LogLevel string
	LogType  string // json or text
	Postgres postgres.Config
}

func main() {
	cfg := Config{
		PORT: getEnvWithDefault("PORT", "8080"),
		Postgres: postgres.Config{
			URL:             os.Getenv("CONTRACTUS_POSTGRES_URL"),
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxIdleTime: 1 * time.Minute,
		},
		LogLevel: getEnvWithDefault("LOG_LEVEL", "INFO"),
		LogType:  getEnvWithDefault("LOG_TYPE", "json"),
	}

	err := setUpLog(cfg)
	if err != nil {
		slog.Error("Failed to set up logger", "error", err)
		syscall.Exit(1)
	}

	db, err := postgres.OpenDB(cfg.Postgres)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		syscall.Exit(1)
	}
	storage := postgres.NewStorage(db)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		api.RegisterHandler(r, storage)
	})

	svc := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	slog.Info("Starting server", "addr", cfg.PORT)
	err = svc.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		syscall.Exit(1)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// setUpLog initialize the logger.
func setUpLog(cfg Config) error {
	var level slog.Level
	switch cfg.LogLevel {
	case "INFO":
		level = slog.LevelInfo
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	var logger *slog.Logger
	if cfg.LogType == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
	} else if cfg.LogType == "text" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
	} else {
		return fmt.Errorf("invalid log type: %s", cfg.LogType)
	}

	slog.SetDefault(logger)
	return nil
}
