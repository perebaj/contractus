// Package main grather all the pieces and start the service.
package main

import (
	"fmt"
	"os"

	"github.com/birdie-ai/contractus/postgres"
)

// Config have the core configuration for the service.
type Config struct {
	PORT     string
	Postgres postgres.Config
}

func main() {
	fmt.Println("Hello, playground")
	// cfg := Config{
	// 	PORT: getEnvWithDefault("PORT", "8080"),
	// 	Postgres: postgres.Config{
	// 		URL:             os.Getenv("CONTRACTUS_POSTGRES_URL"),
	// 		MaxOpenConns:    10,
	// 		MaxIdleConns:    5,
	// 		ConnMaxIdleTime: 1 * time.Minute,
	// 	},
	// }

	// db, err := postgres.OpenDB(cfg.Postgres)
	// if err != nil {
	// 	slog.Error("Failed to connect to database", "error", err)
	// 	syscall.Exit(1)
	// }
	// storage := postgres.NewTransactionStorage(db)

	// r := chi.NewRouter()
	// r.Group(func(r chi.Router) {
	// 	api.RegisterHandler(r, storage)
	// })

	// svc := &http.Server{
	// 	Addr:         ":" + cfg.PORT,
	// 	Handler:      r,
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 5 * time.Second,
	// }
	// slog.Info("Starting server", "addr", cfg.PORT)
	// err = svc.ListenAndServe()
	// if err != nil {
	// 	slog.Error("Failed to start server", "error", err)
	// 	syscall.Exit(1)
	// }
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
