package main

import (
	"context"
	"flag"
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
	"github.com/platepilot/backend/internal/common/vector"
	"github.com/platepilot/backend/internal/recipe/events"
	"github.com/platepilot/backend/internal/recipe/handler"
	pb "github.com/platepilot/backend/internal/recipe/pb"
	"github.com/platepilot/backend/internal/recipe/repository"
	"github.com/platepilot/backend/internal/recipe/seed"
)

func main() {
	// Parse command line flags
	seedFile := flag.String("seed", "", "Path to seed file (e.g., recipes.json)")
	seedOnly := flag.Bool("seed-only", false, "Exit after seeding (use with -seed)")
	flag.Parse()

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
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection
	pool, err := pgxpool.New(ctx, cfg.Database.RecipeDB)
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

	// Initialize vector generator
	vectorGen := vector.NewHashGenerator()

	// Initialize event publisher (optional - only if RabbitMQ is configured)
	var publisher *events.Publisher
	if cfg.RabbitMQ.URL != "" {
		publisher, err = events.NewPublisher(events.PublisherConfig{
			URL:          cfg.RabbitMQ.URL,
			ExchangeName: cfg.RabbitMQ.ExchangeName,
		}, logger)
		if err != nil {
			slog.Warn("failed to create event publisher - continuing without event publishing",
				"error", err,
			)
		} else {
			slog.Info("event publisher initialized")
		}
	}

	// Run seeder if seed file is specified
	if *seedFile != "" {
		seeder := seed.NewSeeder(repo, vectorGen, publisher, logger)
		if err := seeder.SeedFromFile(ctx, *seedFile); err != nil {
			slog.Error("failed to seed database", "error", err)
			os.Exit(1)
		}

		// Exit after seeding if --seed-only is specified
		if *seedOnly {
			slog.Info("seed-only mode: exiting after successful seeding")
			if publisher != nil {
				_ = publisher.Close()
			}
			os.Exit(0)
		}
	}

	// Initialize gRPC handler
	// Explicitly handle nil publisher to avoid interface-wrapping-nil issue
	var eventPublisher handler.EventPublisher
	if publisher != nil {
		eventPublisher = publisher
	}
	grpcHandler := handler.NewGRPCHandler(repo, vectorGen, eventPublisher, logger)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterRecipeServiceServer(grpcServer, grpcHandler)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection for grpcurl/grpcui
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", cfg.RecipeAPI.GRPCAddress)
	if err != nil {
		slog.Error("failed to listen", "address", cfg.RecipeAPI.GRPCAddress, "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server listening", "address", cfg.RecipeAPI.GRPCAddress)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	slog.Info("recipe-api started successfully")

	// Wait for interrupt signal
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()

	slog.Info("shutting down recipe-api...")

	// Graceful shutdown
	cancel() // Cancel main context

	// Stop gRPC server
	grpcServer.GracefulStop()

	// Close event publisher
	if publisher != nil {
		if err := publisher.Close(); err != nil {
			slog.Error("failed to close event publisher", "error", err)
		}
	}

	slog.Info("recipe-api stopped")
}
