package main

import (
	config "auth-service/internal/config"
	controller "auth-service/internal/controller"
	db "auth-service/internal/infra/db"
	redis "auth-service/internal/infra/redis"
	middlewares "auth-service/internal/middleware"
	repository "auth-service/internal/repository"
	service "auth-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	log.Print("Environment: ", cfg.Environment)

	// Connect to DB
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	if err := redis.InitRedis(cfg.RedisAddress, cfg.RedisPassword); err != nil {
		log.Fatal("❌ Failed to connect to redis:", err)
	}
	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)

	// Create Gin instance
	r := gin.Default()

	// Apply middleware
	r.Use(middlewares.ErrorHandler())

	// Register routes
	api := r.Group("/api")
	{
		authController.RegisterRoutes(api) // All user routes require auth
	}

	// Run server
	r.Run(":" + cfg.AppPort)
}
