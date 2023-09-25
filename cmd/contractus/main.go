// contractus runs the service.
// For details on how to configure it just run:
//
//	contractus --help
package main

import (
	"net/http"
	"os"
	"syscall"
	"time"

	"log/slog"
)

type Config struct {
	PORT string
}

func main() {
	cfg := Config{
		PORT: getEnvWithDefault("PORT", "8080"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/jojo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Jojo is aweasome!"))
	})

	svc := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	slog.Info("Starting server", "addr", cfg.PORT)
	err := svc.ListenAndServe()
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
