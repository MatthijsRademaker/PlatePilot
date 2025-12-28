package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/platepilot/backend/internal/common/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Set up structured logging
	logLevel := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("starting recipe-api",
		"environment", cfg.Environment,
		"grpc_address", cfg.RecipeAPI.GRPCAddress,
		"http_address", cfg.RecipeAPI.HTTPAddress,
	)

	// TODO: Initialize database connection
	// TODO: Initialize gRPC server
	// TODO: Initialize event publisher
	// TODO: Start servers

	slog.Info("recipe-api started successfully")

	// Wait for interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down recipe-api...")

	// TODO: Graceful shutdown

	slog.Info("recipe-api stopped")
}
