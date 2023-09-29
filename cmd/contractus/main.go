// Package main grather all the pieces and start the service.
package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/perebaj/contractus/api"
	"github.com/perebaj/contractus/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config have the core configuration for the service.
type Config struct {
	PORT     string
	LogLevel string
	LogType  string // json or text
	Postgres postgres.Config
	Auth     api.Auth
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
		Auth: api.Auth{
			ClientID:     os.Getenv("CONTRACTUS_GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("CONTRACTUS_GOOGLE_CLIENT_SECRET"),
			Domain:       os.Getenv("CONTRACTUS_DOMAIN"), // Example: http://localhost:8080
			RedirectURL:  os.Getenv("CONTRACTUS_DOMAIN") + "/callback",
			JWTSecretKey: os.Getenv("CONTRACTUS_JWT_SECRET_KEY"),
			AccessType:   getEnvWithDefault("CONTRACTUS_ACCESS_TYPE", "online"),
		},
	}

	err := setUpLog(cfg)
	if err != nil {
		slog.Error("Failed to set up logger", "error", err)
		syscall.Exit(1)
	}

	if cfg.Auth.ClientID == "" || cfg.Auth.ClientSecret == "" || cfg.Auth.RedirectURL == "" {
		slog.Error("missing Google OAuth2 configuration")
		syscall.Exit(1)
	}

	db, err := postgres.OpenDB(cfg.Postgres)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		syscall.Exit(1)
	}
	storage := postgres.NewStorage(db)

	googleOAuthConfig := oauth2.Config{
		ClientID:     cfg.Auth.ClientID,
		ClientSecret: cfg.Auth.ClientSecret,
		RedirectURL:  cfg.Auth.RedirectURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}

	cfg.Auth.GoogleOAuthConfig = &googleOAuthConfig

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		// Authenticated routes
		r.Use(jwtauth.Verifier(cfg.Auth.JWTAuth()))
		r.Use(jwtauth.Authenticator)

		api.RegisterTransactionsHandler(r, storage)
	})

	r.Group(func(r chi.Router) {
		// Public routes
		api.RegisterAuthHandler(r, cfg.Auth)
		api.RegisterSwaggerHandler(r)
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
