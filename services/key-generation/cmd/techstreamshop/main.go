package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/storage/postgresql"

	"log/slog"
	"os"

	"github.com/timeb30/techstreamshop/services/key-generation/internal/config"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/http-server/handlers/key/generate"
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
	storage, err := postgresql.New(cfg.StorageURL)
	if err != nil {
		log.Error("failed no init storage", "e", err)
		os.Exit(1)
	}
	//_ = storage
	//id, e := storage.SaveKey(1223, "test key", time.Now(), time.Now().Add(time.Hour))
	//if e != nil {
	//	log.Error("failed save key", "e", e)
	//	os.Exit(1)
	//}
	//log.Info("saved key", slog.Int64("id", id))
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/key", generate.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
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
