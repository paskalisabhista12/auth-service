package main

import (
	"auth-service/internal/config"
	"auth-service/internal/controller"
	"auth-service/internal/infra/db"
	"auth-service/internal/infra/redis"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

func main() {
	// Initialize slog logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	slog.SetDefault(logger)

	// Load config
	cfg := config.LoadConfig()
	slog.Info("config loaded",
		"environment", cfg.Environment)

	// Connect to DB
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		slog.Error("failed to connect to database",
			"error", err,
		)
		os.Exit(1)
	}

	// Connect to Redis
	if err := redis.InitRedis(cfg.RedisAddress, cfg.RedisPassword); err != nil {
		slog.Error("failed to connect to redis",
			"error", err,
		)
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	roleRepo := repository.NewRoleRepository(db.DB)
	endpointRepo := repository.NewEndpointRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, endpointRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)

	// Create Gin instance
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		middlewares.SlogLogger(),
		middlewares.ErrorHandler(),
	)

	// Register routes
	api := r.Group("/api")
	{
		authController.RegisterRoutes(api)
	}
	// Run server
	slog.Info("server started", "port", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		slog.Error("failed to start server",
			"error", err,
		)
		os.Exit(1)
	}
}
