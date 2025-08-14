package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"auth-service/models"
	"github.com/joho/godotenv"
)

func Connect() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get env variable
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Database URL:", dsn)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to database")

	db.AutoMigrate(&models.User{})
}
