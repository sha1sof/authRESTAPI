package main

import (
	"fmt"
	"github.com/sha1sof/authRESTAPI/internal/config"
	"github.com/sha1sof/authRESTAPI/internal/storage"
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

	db, err := storage.New(cfg, log, cfg.Auth.Cost, cfg.Auth.Secret, cfg.Auth.TimeD)
	if err != nil {
		log.Error("failed to connect to database")
		os.Exit(1)
	}

	_, err = db.RegisterUser("example@gmail.com", "123")
	if err != nil {
		log.Error("failed to register user")
	}

	token, _, err := db.LoginUser("example@gmail.com", "123")
	if err != nil {
		log.Error("failed to login")
	}
	fmt.Println(token)

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
