package main

import (
	"github.com/sha1sof/authRESTAPI/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)
	log.Info("starting logger",
		slog.String("env", cfg.Env))

	//TODO: init storage

	//TODO: init server
}

// setupLogger initializes and returns a logger based on the environment.
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("invalid env")
	}
	return logger
}
