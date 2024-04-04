package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paniccaaa/test-task-biocad/internal/config"
	"github.com/paniccaaa/test-task-biocad/internal/lib/logger"
	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing"
	"github.com/paniccaaa/test-task-biocad/internal/router"
	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

func main() {
	config, configDB := config.MustLoad()

	log := logger.SetupLogger(config.Env)

	log.Info(
		"starting test-task-biocad",
		slog.String("env", config.Env),
	)

	log.Debug("debug messages are enabled")

	storage, err := postgres.NewPostgres(configDB)
	if err != nil {
		log.Error("failed to init storage: %s", err)
		os.Exit(1)
	}

	defer storage.Close()

	router := router.InitRouter(log, storage)
	log.Info("starting server", slog.String("address", config.Address))

	srv := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server: %w", err)
		}
	}()

	go func() {
		if err := parsing.Start(config, log, storage); err != nil {
			log.Error("failed to start scanning dir input_files: %w", err)
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server: %w", err)
		return
	}

	log.Info("server stopped")
}
