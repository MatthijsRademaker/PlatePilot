package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/platepilot/backend/internal/common/config"
	"github.com/platepilot/backend/internal/mealplanner/domain"
	"github.com/platepilot/backend/internal/mealplanner/events"
	"github.com/platepilot/backend/internal/mealplanner/handler"
	pb "github.com/platepilot/backend/internal/mealplanner/pb"
	"github.com/platepilot/backend/internal/mealplanner/repository"
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

	slog.Info("starting mealplanner-api",
		"environment", cfg.Environment,
		"grpc_address", cfg.MealPlanner.GRPCAddress,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection
	pool, err := pgxpool.New(ctx, cfg.Database.MealPlannerDB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	// Initialize repository
	repo := repository.NewRepository(pool)

	// Initialize domain planner
	planner := domain.NewPlanner(repo)

	// Initialize gRPC handler
	grpcHandler := handler.NewGRPCHandler(planner, repo, logger)

	// Initialize event consumer (optional - only if RabbitMQ is configured)
	var consumer *events.Consumer
	if cfg.RabbitMQ.URL != "" {
		consumer, err = events.NewConsumer(events.ConsumerConfig{
			URL:          cfg.RabbitMQ.URL,
			ExchangeName: cfg.RabbitMQ.ExchangeName,
			QueueName:    "mealplanner.recipe-events",
			RoutingKey:   "recipe.#",
		}, repo, logger)
		if err != nil {
			slog.Warn("failed to create event consumer - continuing without event sync",
				"error", err,
			)
		} else {
			if err := consumer.Start(ctx); err != nil {
				slog.Error("failed to start event consumer", "error", err)
				consumer.Close()
				consumer = nil
			} else {
				slog.Info("event consumer started")
			}
		}
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterMealPlannerServiceServer(grpcServer, grpcHandler)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for grpcurl/grpcui
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.MealPlanner.GRPCAddress)
	if err != nil {
		slog.Error("failed to listen", "address", cfg.MealPlanner.GRPCAddress, "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server listening", "address", cfg.MealPlanner.GRPCAddress)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	slog.Info("mealplanner-api started successfully")

	// Wait for interrupt signal
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()

	slog.Info("shutting down mealplanner-api...")

	// Graceful shutdown
	cancel() // Cancel main context

	// Stop gRPC server
	grpcServer.GracefulStop()

	// Close event consumer
	if consumer != nil {
		if err := consumer.Close(); err != nil {
			slog.Error("failed to close event consumer", "error", err)
		}
	}

	slog.Info("mealplanner-api stopped")
}
