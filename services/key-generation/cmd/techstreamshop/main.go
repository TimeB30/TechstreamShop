package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/timeb30/techstreamshop/services/key-generation/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv done
	// TODO: init logger: slog done
	// TODO: init storage: pgx done
	// TODO: init router: chi
	// TODO: run server:
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := setupLogger(cfg.Env)
	log.Info("Starting techstreamshop", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")
	//storage, err := postgresql.New(cfg.StorageURL)
	//if err != nil {
	//	log.Error("failed no init storage", "err", err)
	//	os.Exit(1)
	//}
	//_ = storage
	//id, err := storage.SaveKey(1223, "test key", time.Now(), time.Now().Add(time.Hour))
	//if err != nil {
	//	log.Error("failed save key", "err", err)
	//	os.Exit(1)
	//}
	//log.Info("saved url", slog.Int64("id", id))
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
