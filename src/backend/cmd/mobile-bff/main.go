// @title           PlatePilot API
// @version         1.0
// @description     Intelligent meal planning and recipe management API

// @contact.name   PlatePilot Support

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /v1

// @schemes http https

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/platepilot/backend/internal/bff/auth"
	"github.com/platepilot/backend/internal/bff/client"
	"github.com/platepilot/backend/internal/bff/handler"
	bffmiddleware "github.com/platepilot/backend/internal/bff/middleware"
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

	ctx := context.Background()

	slog.Info("starting mobile-bff",
		"environment", cfg.Environment,
		"http_address", cfg.BFF.HTTPAddress,
		"recipe_api", cfg.BFF.RecipeAPIAddress,
		"mealplan_api", cfg.BFF.MealPlanAddress,
	)

	// Create database pool for auth
	dbPool, err := newDBPool(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Create gRPC clients
	recipeClient, err := client.NewRecipeClient(cfg.BFF.RecipeAPIAddress, logger)
	if err != nil {
		slog.Error("failed to create recipe client", "error", err)
		os.Exit(1)
	}
	defer recipeClient.Close()

	mealPlannerClient, err := client.NewMealPlannerClient(cfg.BFF.MealPlanAddress, logger)
	if err != nil {
		slog.Error("failed to create mealplanner client", "error", err)
		os.Exit(1)
	}
	defer mealPlannerClient.Close()

	// Create auth services
	authRepo := auth.NewRepository(dbPool)
	tokenService := auth.NewTokenService(cfg.Auth.JWTSecret, cfg.Auth.Issuer, cfg.Auth.AccessTokenTTL)
	authService := auth.NewService(authRepo, tokenService, cfg.Auth.RefreshTokenTTL)

	// Create handlers
	recipeHandler := handler.NewRecipeHandler(recipeClient, logger)
	mealPlanHandler := handler.NewMealPlanHandler(mealPlannerClient, logger)
	authHandler := handler.NewAuthHandler(authService, logger)

	// Set up router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.BFF.Timeout))
	r.Use(corsMiddleware(cfg.BFF.CORSAllowedOrigins))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check downstream service health
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/logout", authHandler.Logout)
		})

		r.Group(func(r chi.Router) {
			r.Use(bffmiddleware.AuthMiddleware(tokenService))

			r.Route("/recipe", func(r chi.Router) {
				r.Get("/{id}", recipeHandler.GetByID)
				r.Get("/", recipeHandler.List)
				r.Get("/similar", recipeHandler.GetSimilar)
				r.Post("/", recipeHandler.Create)
				r.Put("/{id}", recipeHandler.Update)
				r.Delete("/{id}", recipeHandler.Delete)
				r.Get("/cuisines", recipeHandler.GetCuisines)
				r.Post("/cuisines", recipeHandler.CreateCuisine)
			})
			r.Route("/mealplan", func(r chi.Router) {
				r.Get("/week", mealPlanHandler.GetWeek)
				r.Put("/week", mealPlanHandler.UpsertWeek)
				r.Post("/suggest", mealPlanHandler.Suggest)
			})
		})
	})

	// Create server
	srv := &http.Server{
		Addr:         cfg.BFF.HTTPAddress,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("server listening", "address", cfg.BFF.HTTPAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("mobile-bff started successfully")

	// Wait for interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down mobile-bff...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("mobile-bff stopped")
}

func newDBPool(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.RecipeDB)
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns > 0 {
		poolCfg.MaxConns = int32(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		poolCfg.MinConns = int32(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLife > 0 {
		poolCfg.MaxConnLifetime = cfg.ConnMaxLife
	}

	return pgxpool.NewWithConfig(ctx, poolCfg)
}

// corsMiddleware returns a CORS middleware
func corsMiddleware(allowedOrigins []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "300")
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
